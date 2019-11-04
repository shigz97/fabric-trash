## 删除所有容器
docker rm $(docker ps -aq)

## 设置工作路径
export FABRIC_CFG_PATH=$GOPATH/src/github.com/hyperledger/fabric-samples/trash


## 环境清理
rm -fr config/*
rm -fr crypto-config/*

## 生成证书文件
../bin/cryptogen generate --config=./crypto-config.yaml

## 生成创世区块
../bin/configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block

## 生成通道的创世交易
../bin/configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/mychannel.tx -channelID mychannel
../bin/configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/assetschannel.tx -channelID assetschannel

## 生成组织关于通道的锚节点（主节点）交易
../bin/configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org0MSPanchors.tx -channelID mychannel -asOrg Org0MSP
../bin/configtxgen -profile TwoOrgChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

## 启动网络
docker-compose -f docker-compose.yaml up -d

## 进入CLI容器
docker exec -it cli bash

## 创建通道
peer channel create -o orderer.zjucst.com:7050 -c mychannel -f /etc/hyperledger/config/mychannel.tx

## 加入通道
peer channel join -b mychannel.block

## 设置主节点
peer channel update -o orderer.zjucst.com:7050 -c mychannel -f /etc/hyperledger/config/Org1MSPanchors.tx

## 链码安装
peer chaincode install -n assets -v 1.0 -l golang -p github.com/trash
peer chaincode install -n trash -v 1.0 -l golang -p github.com/trash


## 链码实例化

peer chaincode instantiate -o orderer.zjucst.com:7050 -C mychannel -n trash -l golang -v 1.0 -c '{"Args":["init"]}'


##选择CLI操作或者SDK操作
## 链码交互
peer chaincode invoke -C mychannel -n trash -c '{"Args":["RecyclerRegister", "rcy1", "rcy1"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["ProcessorRegister", "pro1", "pro1"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashEnroll", "t1","t1","c1","10","rcy1"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashEnroll", "t1","t1","c1","10","rcy1"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashEnroll", "t2","t2","c2","30","rcy1"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashTrans", "rcy1","pro1","t1","20"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashTrans", "rcy1","pro1","t2","10"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashProcess", "pro1","t1","burn","10"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashProcess", "pro1","t1","bury","10"]}'
peer chaincode invoke -C mychannel -n trash -c '{"Args":["TrashProcess", "pro1","t2","burn","10"]}'




## 链码查询
peer chaincode query -C mychannel -n trash -c '{"Args":["RecyclerQuery", "rcy1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["ProcessorQuery", "pro1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryRecyleHistory", "rcy1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryRecyleHistory", "rcy1","t1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryRecyleHistory", "rcy1","t2"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryTransHistory", "t1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryTransHistory", "t1","rcy1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryTransHistory", "t1","rcy1","pro1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryProcessHistory", "pro1"]}'
peer chaincode query -C mychannel -n trash -c '{"Args":["queryProcessHistory", "pro1","t1"]}'








