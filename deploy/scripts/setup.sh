#!/bin/bash

# 设置环境变量
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否运行
check_docker() {
    print_info "检查Docker状态..."
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker未运行，请启动Docker"
        exit 1
    fi
    print_info "Docker运行正常"
}

# 下载Fabric二进制文件
download_binaries() {
    print_info "下载Hyperledger Fabric二进制文件..."
    
    if [ ! -d "../bin" ]; then
        curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.4.0 1.5.0
        mv fabric-samples/bin ../
        mv fabric-samples/config ../
        rm -rf fabric-samples
    else
        print_info "二进制文件已存在，跳过下载"
    fi
}

# 生成加密材料
generate_crypto() {
    print_info "生成加密材料..."
    
    if [ -d "crypto-config" ]; then
        print_warn "删除现有加密材料..."
        rm -rf crypto-config
    fi
    
    ../bin/cryptogen generate --config=./crypto-config.yaml
    
    if [ $? -ne 0 ]; then
        print_error "生成加密材料失败"
        exit 1
    fi
    
    print_info "加密材料生成成功"
}

# 生成创世区块和通道配置
generate_artifacts() {
    print_info "生成创世区块和通道配置..."
    
    if [ -d "channel-artifacts" ]; then
        rm -rf channel-artifacts
    fi
    mkdir channel-artifacts
    
    # 生成创世区块
    print_info "生成创世区块..."
    ../bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./channel-artifacts/genesis.block
    
    if [ $? -ne 0 ]; then
        print_error "生成创世区块失败"
        exit 1
    fi
    
    # 生成通道配置交易
    print_info "生成通道配置交易..."
    ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID woodcarving-channel
    
    if [ $? -ne 0 ]; then
        print_error "生成通道配置交易失败"
        exit 1
    fi
    
    # 生成锚节点配置
    print_info "生成Org1锚节点配置..."
    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID woodcarving-channel -asOrg Org1MSP
    
    print_info "生成Org2锚节点配置..."
    ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID woodcarving-channel -asOrg Org2MSP
    
    print_info "配置文件生成成功"
}

# 启动网络
start_network() {
    print_info "启动Fabric网络..."
    
    # 停止现有容器
    docker-compose down
    
    # 清理Docker卷
    docker volume prune -f
    
    # 启动网络
    docker-compose up -d
    
    if [ $? -ne 0 ]; then
        print_error "启动网络失败"
        exit 1
    fi
    
    print_info "等待网络启动..."
    sleep 10
    
    # 检查容器状态
    print_info "检查容器状态..."
    docker-compose ps
}

# 创建通道
create_channel() {
    print_info "创建通道..."
    
    # 设置环境变量
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/peers/peer0.org1.woodcarving.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/users/Admin@org1.woodcarving.com/msp
    export CORE_PEER_ADDRESS=peer0.org1.woodcarving.com:7051
    
    # 创建通道
    docker exec cli peer channel create -o orderer.woodcarving.com:7050 -c woodcarving-channel -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
    
    if [ $? -ne 0 ]; then
        print_error "创建通道失败"
        exit 1
    fi
    
    print_info "通道创建成功"
}

# 加入通道
join_channel() {
    print_info "节点加入通道..."
    
    # Org1 Peer0 加入通道
    print_info "Org1 Peer0 加入通道..."
    docker exec cli peer channel join -b woodcarving-channel.block
    
    # 切换到Org2
    print_info "Org2 Peer0 加入通道..."
    docker exec -e CORE_PEER_LOCALMSPID="Org2MSP" \
                -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
                -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/users/Admin@org2.woodcarving.com/msp \
                -e CORE_PEER_ADDRESS=peer0.org2.woodcarving.com:9051 \
                cli peer channel join -b woodcarving-channel.block
    
    print_info "所有节点已加入通道"
}

# 更新锚节点
update_anchor_peers() {
    print_info "更新锚节点..."
    
    # 更新Org1锚节点
    docker exec cli peer channel update -o orderer.woodcarving.com:7050 -c woodcarving-channel -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
    
    # 更新Org2锚节点
    docker exec -e CORE_PEER_LOCALMSPID="Org2MSP" \
                -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
                -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/users/Admin@org2.woodcarving.com/msp \
                -e CORE_PEER_ADDRESS=peer0.org2.woodcarving.com:9051 \
                cli peer channel update -o orderer.woodcarving.com:7050 -c woodcarving-channel -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
    
    print_info "锚节点更新完成"
}

# 主函数
main() {
    print_info "开始部署Hyperledger Fabric网络..."
    
    check_docker
    download_binaries
    generate_crypto
    generate_artifacts
    start_network
    create_channel
    join_channel
    update_anchor_peers
    
    print_info "Fabric网络部署完成！"
    print_info "网络组件："
    print_info "  - Orderer: orderer.woodcarving.com:7050"
    print_info "  - Org1 Peer0: peer0.org1.woodcarving.com:7051"
    print_info "  - Org2 Peer0: peer0.org2.woodcarving.com:9051"
    print_info "  - 通道名称: woodcarving-channel"
    print_info ""
    print_info "下一步：部署链码"
    print_info "运行: ./deploy-chaincode.sh"
}

# 执行主函数
main