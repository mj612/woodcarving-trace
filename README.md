# 基于区块链的特色木雕工艺品溯源系统

## 项目简介

本项目是一个基于区块链技术的木雕工艺品全生命周期溯源系统，旨在实现从原料采购到消费者购买的全流程可信追溯。系统利用Hyperledger Fabric区块链确保数据不可篡改，通过Web应用为供应链各参与方提供便捷的数据录入和查询服务。

## 技术栈

### 前端
- **Vue.js 2.6+**: 渐进式JavaScript框架
- **Element UI**: 基于Vue的组件库
- **Axios**: HTTP客户端
- **Vue Router**: 路由管理
- **Vuex**: 状态管理

### 后端
- **Go 1.19+**: 主要开发语言
- **Gin**: Web框架
- **GORM**: ORM框架
- **MySQL 8.0**: 关系型数据库
- **Redis**: 缓存数据库
- **JWT**: 身份认证

### 区块链
- **Hyperledger Fabric 2.5**: 联盟链平台
- **Fabric Go SDK**: 区块链交互SDK
- **Go**: 智能合约开发语言

### 部署
- **Docker & Docker Compose**: 容器化部署
- **Nginx**: 反向代理（可选）

## 功能列表

### 1. 用户管理模块
- 用户注册与登录
- 角色分配（供应商、工匠、仓管、销售商、消费者、监管方）
- 基于RBAC的权限控制
- 个人资料管理

### 2. 原料管理模块（供应商）
- 新增原料批次
- 原料信息录入（种类、产地、采伐证等）
- 原料列表查询
- 原料转移给工匠

### 3. 生产管理模块（工匠）
- 领取原料
- 成品登记（关联原料ID）
- 雕刻过程记录
- 作品列表管理

### 4. 仓储管理模块（仓管）
- 产品入库登记
- 产品出库登记
- 库存盘点
- 环境监控数据记录

### 5. 销售管理模块（销售商）
- 产品上架
- 订单管理
- 物流跟踪
- 销售记录

### 6. 溯源查询模块（公开）
- 扫码查询产品信息
- 可视化流程展示
- 完整溯源链条查看

### 7. 管理后台（监管方）
- 区块链交易查看
- 数据审核
- 异常标记
- 系统统计

## 系统特点

1. **数据不可篡改**: 关键数据上链，利用区块链特性确保数据真实性
2. **全程可追溯**: 从原料到消费者，记录完整流转过程
3. **权限分级**: RBAC权限控制，不同角色不同权限
4. **高性能**: Redis缓存热点数据，提升查询效率
5. **易部署**: Docker容器化，一键启动
6. **用户友好**: 直观的Web界面，简单易用

## 项目结构

```
woodcarving-trace/
├── backend/                # Go后端服务
│   ├── api/               # API控制器
│   ├── config/            # 配置文件
│   ├── middleware/        # 中间件
│   ├── models/            # 数据模型
│   ├── pkg/               # 工具包
│   │   ├── blockchain/    # 区块链SDK封装
│   │   ├── cache/         # Redis缓存
│   │   └── utils/         # 通用工具
│   ├── routes/            # 路由定义
│   ├── main.go            # 程序入口
│   └── go.mod             # Go依赖
├── frontend/              # Vue前端应用
│   ├── public/            # 静态资源
│   ├── src/
│   │   ├── api/          # API封装
│   │   ├── assets/       # 资源文件
│   │   ├── components/   # 公共组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # Vuex状态
│   │   ├── views/        # 页面组件
│   │   ├── App.vue       # 根组件
│   │   └── main.js       # 入口文件
│   ├── package.json       # 依赖配置
│   └── vue.config.js      # Vue配置
├── chaincode/             # 智能合约
│   ├── main.go           # 合约入口
│   ├── contract.go       # 合约实现
│   └── go.mod            # 依赖管理
├── deploy/                # 部署文件
│   ├── docker-compose.yml # Docker编排
│   ├── fabric/           # Fabric网络配置
│   └── scripts/          # 启动脚本
└── docs/                  # 文档
    ├── design.md         # 设计文档
    └── deploy.md         # 部署文档

```

## 快速开始

详细部署步骤请参考 [部署文档](docs/deploy.md)

### 前置条件
- Go 1.19+
- Node.js 14+
- Docker & Docker Compose
- MySQL 8.0
- Redis

### 启动步骤

1. 启动Fabric网络
```bash
cd deploy
./network.sh up
```

2. 安装智能合约
```bash
./network.sh deployCC
```

3. 启动后端服务
```bash
cd backend
go mod download
go run main.go
```

4. 启动前端服务
```bash
cd frontend
npm install
npm run serve
```

5. 访问系统
- 前端: http://localhost:8080
- 后端API: http://localhost:3000
- 管理后台: http://localhost:8080/admin

## 许可证

MIT License

## 联系方式

如有问题，请通过以下方式联系：
- Email: 1125933779@qq.com


