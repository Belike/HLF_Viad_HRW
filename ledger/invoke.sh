echo "$1"
echo ##PEER0 ORG2 QUERY##
if [ $1 = "Org2" ]
then docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.hrw.com/users/Admin@org2.hrw.com/msp" -e CORE_PEER_ADDRESS=peer0.org2.hrw.com:7051 cli peer chaincode query -o orderer.hrw.com:7050 -C mychannel -n supply -c '{"function":"queryAllGoods","Args":["A"]}'
fi

