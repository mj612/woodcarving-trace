# 木雕工艺品溯源系统部署文档

## 1. 环境要求

### 1.1 基础环境
- **操作系统**: Ubuntu 20.04+ / CentOS 7+ / macOS
- **Go**: 1.19或更高版本
- **Node.js**: 14.x或更高版本
- **Docker**: 20.10+
- **Docker Compose**: 1.29+
- **MySQL**: 8.0+
- **Redis**: 6.0+

### 1.2 Hyperledger Fabric环境
- Fabric 2.5
- Fabric CA 1.5
- 至少4GB内存
- 至少20GB磁盘空间

## 2. 前置准备

### 2.1 安装Go
```bash
# 下载Go
wget https://go.dev/dl/go1.19.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 2.2 安装Node.js
```bash
# 使用nvm安装
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 14
nvm use 14

# 验证安装
node -v
npm -v
```

### 2.3 安装Docker
```bash
# Ubuntu系统
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 添加当前用户到docker组
sudo usermod -aG docker $USER

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

### 2.4 安装MySQL
```bash
# Ubuntu系统
sudo apt update
sudo apt install mysql-server

# 启动MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation

# 创建数据库
mysql -u root -p
CREATE DATABASE woodcarving_trace CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'woodcarving'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON woodcarving_trace.* TO 'woodcarving'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 2.5 安装Redis
```bash
# Ubuntu系统
sudo apt install redis-server

# 启动Redis
sudo systemctl start redis-server
sudo systemctl enable redis-server

# 验证
redis-cli ping
```

## 3. Hyperledger Fabric网络部署

### 3.1 下载Fabric工具和镜像
```bash
cd ~/
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.0 1.5.5

# 添加到PATH
echo 'export PATH=$PATH:$HOME/fabric-samples/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3.2 启动Fabric测试网络
```bash
cd ~/fabric-samples/test-network

# 启动网络
./network.sh up createChannel -c mychannel -ca

# 设置环境变量
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

### 3.3 部署智能合约
```bash
# 将智能合约代码复制到chaincode目录
cp -r /path/to/woodcarving-trace/chaincode ~/fabric-samples/chaincode/woodcarving

# 打包链码
peer lifecycle chaincode package woodcarving.tar.gz --path ~/fabric-samples/chaincode/woodcarving --lang golang --label woodcarving_1.0

# 安装链码到Org1
peer lifecycle chaincode install woodcarving.tar.gz

# 记录PackageID
peer lifecycle chaincode queryinstalled

# 设置PackageID环境变量（替换为实际值）
export CC_PACKAGE_ID=woodcarving_1.0:xxxxxxxxxxxx

# 批准链码（Org1）
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name woodcarving --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# 切换到Org2
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# 安装链码到Org2
peer lifecycle chaincode install woodcarving.tar.gz

# 批准链码（Org2）
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name woodcarving --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# 提交链码定义
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name woodcarving --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# 验证链码
peer lifecycle chaincode querycommitted --channelID mychannel --name woodcarving
```

## 4. 后端服务部署

### 4.1 配置连接信息
```bash
cd /path/to/woodcarving-trace/backend

# 修改配置文件
vim config/config.yaml

# 配置数据库连接
database:
  host: "localhost"
  port: "3306"
  user: "woodcarving"
  password: "your_password"
  dbname: "woodcarving_trace"
  
# 配置Fabric连接
fabric:
  config_path: "./config/connection-profile.yaml"
  channel_name: "mychannel"
  chaincode_name: "woodcarving"
  org_name: "Org1"
  user_name: "Admin"
```

### 4.2 复制Fabric配置文件
```bash
# 从test-network复制连接配置
cp ~/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.yaml ./config/connection-profile.yaml

# 创建wallet目录并复制身份证书
mkdir -p wallet
cp -r ~/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp wallet/Admin
```

### 4.3 安装依赖并启动
```bash
# 下载依赖
go mod download
go mod tidy

# 初始化Fabric连接（可选）
# 在代码中添加初始化逻辑或单独运行初始化脚本

# 启动服务
go run main.go

# 或编译后运行
go build -o woodcarving-backend
./woodcarving-backend
```

### 4.4 使用Supervisor管理（生产环境推荐）
```bash
# 安装supervisor
sudo apt install supervisor

# 创建配置文件
sudo vim /etc/supervisor/conf.d/woodcarving-backend.conf

# 配置内容
[program:woodcarving-backend]
command=/path/to/woodcarving-backend/woodcarving-backend
directory=/path/to/woodcarving-backend
autostart=true
autorestart=true
user=ubuntu
stdout_logfile=/var/log/woodcarving/backend.log
stderr_logfile=/var/log/woodcarving/backend_error.log

# 重启supervisor
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start woodcarving-backend
```

## 5. 前端应用部署

### 5.1 开发环境运行
```bash
cd /path/to/woodcarving-trace/frontend

# 安装依赖
npm install

# 配置API地址
# 创建.env文件
echo "VUE_APP_API_BASE_URL=http://localhost:3000/api/v1" > .env.local

# 启动开发服务器
npm run serve
```

### 5.2 生产环境部署
```bash
# 构建生产版本
npm run build

# 构建完成后，dist目录包含静态文件
# 使用Nginx部署
sudo apt install nginx

# 创建Nginx配置
sudo vim /etc/nginx/sites-available/woodcarving-trace

# 配置内容
server {
    listen 80;
    server_name your-domain.com;
    
    root /path/to/woodcarving-trace/frontend/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# 启用配置
sudo ln -s /etc/nginx/sites-available/woodcarving-trace /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 6. 测试数据

### 6.1 创建测试用户
```bash
# 使用API或直接插入数据库
# 密码: test123456 (BCrypt加密后)
mysql -u woodcarving -p woodcarving_trace << EOF
INSERT INTO users (username, password, real_name, role, status) VALUES
('supplier1', '$2a$10$...', '张供应商', 'supplier', 1),
('artisan1', '$2a$10$...', '李工匠', 'artisan', 1),
('warehouse1', '$2a$10$...', '王仓管', 'warehouse', 1),
('seller1', '$2a$10$...', '赵销售', 'seller', 1),
('admin', '$2a$10$...', '管理员', 'supervisor', 1);
EOF
```

### 6.2 测试流程
1. 登录系统（supplier1）
2. 创建原料记录
3. 切换到artisan1账号
4. 使用原料创建产品
5. 切换到warehouse1账号
6. 记录入库信息
7. 记录出库信息
8. 切换到seller1账号
9. 记录销售信息
10. 访问溯源页面查看完整链条

## 7. 常见问题

### 7.1 Fabric网络启动失败
```bash
# 清理并重启
cd ~/fabric-samples/test-network
./network.sh down
docker system prune -a
./network.sh up createChannel -c mychannel -ca
```

### 7.2 链码安装失败
- 检查Go模块路径
- 确保链码代码无语法错误
- 检查网络连接

### 7.3 后端连接Fabric失败
- 检查connection-profile.yaml路径
- 确认证书文件路径正确
- 检查Fabric网络是否运行

### 7.4 数据库连接失败
- 检查MySQL服务状态
- 确认数据库用户权限
- 检查配置文件中的连接参数

## 8. 监控与维护

### 8.1 日志查看
```bash
# 后端日志
tail -f /var/log/woodcarving/backend.log

# Fabric日志
docker logs -f peer0.org1.example.com
docker logs -f orderer.example.com
```

### 8.2 性能监控
- 使用Prometheus + Grafana监控系统性能
- 监控数据库连接数
- 监控Redis内存使用
- 监控Fabric节点状态

### 8.3 备份策略
- 定期备份MySQL数据库
- 备份Fabric账本数据
- 备份上传文件
- 备份配置文件

## 9. 扩展建议

### 9.1 高可用部署
- 使用多个后端服务实例 + 负载均衡
- MySQL主从复制
- Redis哨兵模式
- Fabric多节点部署

### 9.2 性能优化
- CDN加速静态资源
- 数据库读写分离
- 增加缓存层
- 异步处理耗时任务

### 9.3 安全加固
- 启用HTTPS
- 实施API限流
- 增强密码策略
- 定期安全审计

---

**注意**: 本文档为部署指南，生产环境部署时请根据实际情况调整配置，并做好安全加固工作。
