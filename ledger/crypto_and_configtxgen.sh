export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mychannel
export PATH=$PATH:~/HLF_Viad_HRW/ledger/bin
echo "#################################################"
echo "##################  Aufräumen  ##################"
docker rm -f $(docker ps -aq)
docker rmi $(docker images)
echo ""
echo "########### Images werden neu geladen ###########"
./get_images.sh
echo "########## Lösche altes Cryptomaterial ##########"
rm -rf crypto-config
rm -rf channel-artifacts
sudo rm -rf chaincode/hyperledger
echo "############ Cryptogen wird erstellt ############"
cryptogen generate --config=./crypto-config.yaml
echo "############ Configtx wird erstellt #############"
if [ ! -d channel-artifacts ] ; then
	mkdir channel-artifacts
fi
sleep 2
echo "### Genesis und ChannelTx werden erstellt... ###"
configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
sleep 5
echo "############## AnchorPeersUpdate ##############"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org3MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org3MSP
sleep 2
echo "######### Netzwerk wird gestartet... #########"
CHANNEL_NAME=$CHANNEL_NAME docker-compose -f docker-compose.yml up -d
sleep 10
echo "## Peer0.Org1.HRW.COM Definition gestartet ##"
docker exec -it cli peer channel create -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/channel.tx
echo "######### Channel wird gestartet ############"
docker exec -it cli peer channel join -b mychannel.block
docker exec -it cli peer channel update -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx
echo "## Peer0.Org2.HRW.COM Definition gestartet ## "
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051" -it cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051" -it cli peer channel update -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx
echo "## Peer0.Org3.HRW.COM Definition gestartet ## "
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.hrw.com/users/Admin@org3.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org3.hrw.com:7051" -it cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.hrw.com/users/Admin@org3.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org3.hrw.com:7051" -it cli peer channel update -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/Org3MSPanchors.tx
echo "### Netz gestartet, warte 5 sekunden.... ###"
sleep 5
echo "#### PEER0 ORG1 INSTALL CHAINCODE Supply ###"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hrw.com/users/Admin@org1.hrw.com/msp" cli peer chaincode install -n supply -v 1.0 -p github.com/supply
echo "#### PEER0 ORG2 INSTALL CHAINCODE Supply ###"
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051 cli peer chaincode install -n supply -v 1.0 -p github.com/supply
echo "#### PEER0 ORG3 INSTALL CHAINCODE Supply ###"
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.hrw.com/users/Admin@org3.hrw.com/msp" -e CORE_PEER_ADDRESS=peer0.org3.hrw.com:7051 cli peer chaincode install -n supply -v 1.0 -p github.com/supply
sleep 20
echo "##### PEER0 ORG1 INSTANTIATE CHAINCODE #####"
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.hrw.com/users/Admin@org1.hrw.com/msp"  cli peer chaincode instantiate -o orderer.hrw.com:7050 -C mychannel -n supply -v 1.0 -c '{"Args":[""]}' -P "OR ('Org1MSP.member','Org2MSP.member')"
