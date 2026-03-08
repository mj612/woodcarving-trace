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

# 链码信息
CHAINCODE_NAME="woodcarving-trace"
CHAINCODE_VERSION="1.0"
CHANNEL_NAME="woodcarving-channel"
CHAINCODE_PATH="/opt/gopath/src/github.com/chaincode"

# 打包链码
package_chaincode() {
    print_info "打包链码..."
    
    # 删除现有包
    rm -f ${CHAINCODE_NAME}.tar.gz
    
    # 打包链码
    docker exec cli peer lifecycle chaincode package ${CHAINCODE_NAME}.tar.gz \
        --path ${CHAINCODE_PATH} \
        --lang golang \
        --label ${CHAINCODE_NAME}_${CHAINCODE_VERSION}
    
    if [ $? -ne 0 ]; then
        print_error "打包链码失败"
        exit 1
    fi
    
    print_info "链码打包成功"
}

# 安装链码到Org1
install_chaincode_org1() {
    print_info "在Org1上安装链码..."
    
    # 设置Org1环境变量
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/peers/peer0.org1.woodcarving.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/users/Admin@org1.woodcarving.com/msp
    export CORE_PEER_ADDRESS=peer0.org1.woodcarving.com:7051
    
    docker exec cli peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz
    
    if [ $? -ne 0 ]; then
        print_error "在Org1上安装链码失败"
        exit 1
    fi
    
    print_info "Org1链码安装成功"
}

# 安装链码到Org2
install_chaincode_org2() {
    print_info "在Org2上安装链码..."
    
    docker exec -e CORE_PEER_LOCALMSPID="Org2MSP" \
                -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
                -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/users/Admin@org2.woodcarving.com/msp \
                -e CORE_PEER_ADDRESS=peer0.org2.woodcarving.com:9051 \
                cli peer lifecycle chaincode install ${CHAINCODE_NAME}.tar.gz
    
    if [ $? -ne 0 ]; then
        print_error "在Org2上安装链码失败"
        exit 1
    fi
    
    print_info "Org2链码安装成功"
}

# 查询已安装的链码
query_installed() {
    print_info "查询已安装的链码..."
    
    # 查询Org1
    print_info "Org1已安装的链码："
    docker exec cli peer lifecycle chaincode queryinstalled
    
    # 查询Org2
    print_info "Org2已安装的链码："
    docker exec -e CORE_PEER_LOCALMSPID="Org2MSP" \
                -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
                -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/users/Admin@org2.woodcarving.com/msp \
                -e CORE_PEER_ADDRESS=peer0.org2.woodcarving.com:9051 \
                cli peer lifecycle chaincode queryinstalled
}

# 获取包ID
get_package_id() {
    print_info "获取链码包ID..."
    
    PACKAGE_ID=$(docker exec cli peer lifecycle chaincode queryinstalled | grep ${CHAINCODE_NAME}_${CHAINCODE_VERSION} | cut -d' ' -f3 | cut -d',' -f1)
    
    if [ -z "$PACKAGE_ID" ]; then
        print_error "无法获取包ID"
        exit 1
    fi
    
    print_info "包ID: $PACKAGE_ID"
    export PACKAGE_ID
}

# 批准链码定义 - Org1
approve_chaincode_org1() {
    print_info "Org1批准链码定义..."
    
    docker exec cli peer lifecycle chaincode approveformyorg \
        -o orderer.woodcarving.com:7050 \
        --channelID $CHANNEL_NAME \
        --name $CHAINCODE_NAME \
        --version $CHAINCODE_VERSION \
        --package-id $PACKAGE_ID \
        --sequence 1 \
        --tls \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
    
    if [ $? -ne 0 ]; then
        print_error "Org1批准链码定义失败"
        exit 1
    fi
    
    print_info "Org1批准成功"
}

# 批准链码定义 - Org2
approve_chaincode_org2() {
    print_info "Org2批准链码定义..."
    
    docker exec -e CORE_PEER_LOCALMSPID="Org2MSP" \
                -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
                -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/users/Admin@org2.woodcarving.com/msp \
                -e CORE_PEER_ADDRESS=peer0.org2.woodcarving.com:9051 \
                cli peer lifecycle chaincode approveformyorg \
                -o orderer.woodcarving.com:7050 \
                --channelID $CHANNEL_NAME \
                --name $CHAINCODE_NAME \
                --version $CHAINCODE_VERSION \
                --package-id $PACKAGE_ID \
                --sequence 1 \
                --tls \
                --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
    
    if [ $? -ne 0 ]; then
        print_error "Org2批准链码定义失败"
        exit 1
    fi
    
    print_info "Org2批准成功"
}

# 检查提交就绪状态
check_commit_readiness() {
    print_info "检查链码提交就绪状态..."
    
    docker exec cli peer lifecycle chaincode checkcommitreadiness \
        --channelID $CHANNEL_NAME \
        --name $CHAINCODE_NAME \
        --version $CHAINCODE_VERSION \
        --sequence 1 \
        --tls \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem \
        --output json
}

# 提交链码定义
commit_chaincode() {
    print_info "提交链码定义..."
    
    docker exec cli peer lifecycle chaincode commit \
        -o orderer.woodcarving.com:7050 \
        --channelID $CHANNEL_NAME \
        --name $CHAINCODE_NAME \
        --version $CHAINCODE_VERSION \
        --sequence 1 \
        --tls \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem \
        --peerAddresses peer0.org1.woodcarving.com:7051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/peers/peer0.org1.woodcarving.com/tls/ca.crt \
        --peerAddresses peer0.org2.woodcarving.com:9051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt
    
    if [ $? -ne 0 ]; then
        print_error "提交链码定义失败"
        exit 1
    fi
    
    print_info "链码定义提交成功"
}

# 查询已提交的链码
query_committed() {
    print_info "查询已提交的链码..."
    
    docker exec cli peer lifecycle chaincode querycommitted \
        --channelID $CHANNEL_NAME \
        --name $CHAINCODE_NAME \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem
}

# 初始化链码
init_chaincode() {
    print_info "初始化链码..."
    
    docker exec cli peer chaincode invoke \
        -o orderer.woodcarving.com:7050 \
        --tls \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        --peerAddresses peer0.org1.woodcarving.com:7051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/peers/peer0.org1.woodcarving.com/tls/ca.crt \
        --peerAddresses peer0.org2.woodcarving.com:9051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
        -c '{"function":"Init","Args":[]}'
    
    if [ $? -ne 0 ]; then
        print_error "初始化链码失败"
        exit 1
    fi
    
    print_info "链码初始化成功"
}

# 测试链码
test_chaincode() {
    print_info "测试链码功能..."
    
    # 测试创建原料
    print_info "测试创建原料..."
    docker exec cli peer chaincode invoke \
        -o orderer.woodcarving.com:7050 \
        --tls \
        --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/woodcarving.com/orderers/orderer.woodcarving.com/msp/tlscacerts/tlsca.woodcarving.com-cert.pem \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        --peerAddresses peer0.org1.woodcarving.com:7051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.woodcarving.com/peers/peer0.org1.woodcarving.com/tls/ca.crt \
        --peerAddresses peer0.org2.woodcarving.com:9051 \
        --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.woodcarving.com/peers/peer0.org2.woodcarving.com/tls/ca.crt \
        -c '{"function":"CreateRawMaterial","Args":["TEST001","红木","云南","YN-HC-001","supplier1","测试供应商",2.5,"优","hash123"]}'
    
    sleep 3
    
    # 测试查询原料
    print_info "测试查询原料..."
    docker exec cli peer chaincode query \
        -C $CHANNEL_NAME \
        -n $CHAINCODE_NAME \
        -c '{"function":"GetMaterialByID","Args":["TEST001"]}'
    
    print_info "链码测试完成"
}

# 主函数
main() {
    print_info "开始部署链码..."
    
    package_chaincode
    install_chaincode_org1
    install_chaincode_org2
    query_installed
    get_package_id
    approve_chaincode_org1
    approve_chaincode_org2
    check_commit_readiness
    commit_chaincode
    query_committed
    init_chaincode
    test_chaincode
    
    print_info "链码部署完成！"
    print_info "链码信息："
    print_info "  - 名称: $CHAINCODE_NAME"
    print_info "  - 版本: $CHAINCODE_VERSION"
    print_info "  - 通道: $CHANNEL_NAME"
    print_info "  - 包ID: $PACKAGE_ID"
    print_info ""
    print_info "现在可以配置后端连接到Fabric网络"
}

# 执行主函数
main