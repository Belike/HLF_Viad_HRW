export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mychannel

echo "Aufräumen"
docker rm -f $(docker ps -aq)
rm -rf crypto-config
rm -rf channel-artifacts
sudo rm -rf chaincode/hyperledger

echo "Cryptogen wird erstellt"
cryptogen generate --config=./crypto-config.yaml
if [ ! -d channel-artifacts ] ; then
	mkdir channel-artifacts
fi
echo "Configtxgen"
sleep 2
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
echo "Warte 5 Sekunden bis AnchorPeersUpdate"
sleep 5
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
sleep 2
echo "Docker Compose um Netz zu starten"
CHANNEL_NAME=$CHANNEL_NAME docker-compose -f docker-compose.yml up -d
sleep 10
echo "Peer0.Org1.HRW.COM Definition gestartet"
docker exec -it cli peer channel create -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/channel.tx
docker exec -it cli peer channel join -b mychannel.block
docker exec -it cli peer channel update -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx

echo "Peer0.Org2.HRW.COM Definition gestartet"
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051" -it cli peer channel join -b mychannel.block
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e "CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051" -it cli peer channel update -o orderer.hrw.com:7050 -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx

