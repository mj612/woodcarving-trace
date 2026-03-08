#!/bin/bash

# 简单的网络管理脚本
# 入口与 README 中的使用方式保持一致：
#   ./network.sh up        启动并初始化 Fabric 网络（含证书、通道等）
#   ./network.sh deployCC  部署并测试链码
#   ./network.sh down      关闭网络并清理数据

set -e

ROOT_DIR=$(cd "$(dirname "$0")"; pwd)

case "$1" in
  up)
    echo "===> 启动 Hyperledger Fabric 网络（生成证书、创世块、通道等）..."
    (cd "$ROOT_DIR/scripts" && ./setup.sh)
    echo "===> Fabric 网络启动完成。"
    ;;

  deployCC)
    echo "===> 部署并测试链码..."
    (cd "$ROOT_DIR/scripts" && ./deploy-chaincode.sh)
    echo "===> 链码部署与测试完成。"
    ;;

  down)
    echo "===> 关闭 Fabric 网络并清理数据..."
    cd "$ROOT_DIR"

    # 关闭并清理容器
    docker-compose down -v || true

    # 删除生成的加密材料和通道文件
    rm -rf crypto-config channel-artifacts

    echo "===> 网络已关闭，生成数据已清理。"
    ;;

  *)
    echo "Usage: $0 {up|deployCC|down}"
    exit 1
    ;;
esac

