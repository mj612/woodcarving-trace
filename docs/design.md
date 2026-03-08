# 木雕工艺品溯源系统设计文档

## 1. 系统概述

本系统是基于区块链技术的木雕工艺品全生命周期溯源平台，通过Hyperledger Fabric实现数据不可篡改和全程可追溯。

### 1.1 系统目标
- 实现木雕工艺品从原料到消费者的全流程追溯
- 利用区块链确保溯源数据真实可信
- 为各参与方提供便捷的数据录入和查询服务
- 为消费者提供透明的产品信息

## 2. 系统架构

### 2.1 总体架构图
```
前端层(Vue.js) → 应用层(Gin) → 数据层(MySQL/Redis) → 区块链层(Fabric)
```

### 2.2 技术栈
- **前端**: Vue.js 2.6 + Element UI + Axios
- **后端**: Go + Gin + GORM + Redis
- **区块链**: Hyperledger Fabric 2.5
- **数据库**: MySQL 8.0

## 3. 数据库设计

### 3.1 users表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| username | VARCHAR(50) | 用户名 |
| password | VARCHAR(255) | 密码(加密) |
| role | VARCHAR(20) | 角色 |
| status | TINYINT | 状态 |

### 3.2 materials表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| material_id | VARCHAR(50) | 原料ID |
| user_id | INT | 供应商ID |
| tx_id | VARCHAR(100) | 交易ID |

### 3.3 products表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INT | 主键 |
| product_id | VARCHAR(50) | 产品ID |
| user_id | INT | 工匠ID |
| tx_id | VARCHAR(100) | 交易ID |

## 4. 智能合约设计

### 4.1 数据结构
- **RawMaterial**: 原料信息
- **Product**: 产品信息
- **TransferRecord**: 流转记录
- **StorageRecord**: 仓储记录
- **SalesRecord**: 销售记录

### 4.2 主要函数
- CreateRawMaterial: 创建原料
- CreateProduct: 创建产品
- TransferProduct: 转移产品
- GetCompleteTrace: 完整溯源

## 5. API设计

### 5.1 认证接口
- POST /api/v1/register - 注册
- POST /api/v1/login - 登录

### 5.2 原料管理
- POST /api/v1/materials - 创建原料
- GET /api/v1/materials - 原料列表
- GET /api/v1/materials/:id - 原料详情

### 5.3 产品管理
- POST /api/v1/products - 创建产品
- GET /api/v1/products - 产品列表
- GET /api/v1/products/:id - 产品详情

### 5.4 溯源查询
- GET /api/v1/trace/:id - 完整溯源(公开)

## 6. 业务流程

### 6.1 原料上链流程
1. 供应商登录系统
2. 填写原料信息
3. 上传证书文件
4. 提交后自动上链
5. 生成原料ID

### 6.2 产品生产流程
1. 工匠选择原料
2. 填写产品信息
3. 关联原料ID
4. 上传作品图片
5. 提交上链

### 6.3 溯源查询流程
1. 扫描产品二维码
2. 获取产品ID
3. 查询链上数据
4. 展示完整溯源链

## 7. 安全设计

### 7.1 身份认证
- JWT Token机制
- 密码BCrypt加密
- Token过期自动刷新

### 7.2 权限控制
- RBAC角色权限
- API接口鉴权
- 数据隔离

### 7.3 数据安全
- 链上数据不可篡改
- 文件哈希验证
- HTTPS传输

## 8. 性能优化

### 8.1 缓存策略
- Redis缓存热点数据
- 查询结果缓存5-10分钟
- 减少链上查询次数

### 8.2 数据库优化
- 合理建立索引
- 分页查询
- 读写分离(扩展)

