package main

import (
	"log"
	"woodcarving-backend/config"
	"woodcarving-backend/models"
	"woodcarving-backend/pkg/blockchain"
	"woodcarving-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	models.InitDB()

	// 初始化区块链
	if config.GlobalConfig.Fabric.Enabled {
		// 使用真实的 Fabric 区块链
		err := blockchain.InitFabric()
		if err != nil {
			log.Printf("Fabric初始化失败: %v，切换到模拟模式", err)
			blockchain.InitMockBlockchain()
			log.Println("区块链模块初始化成功（模拟模式）")
		} else {
			log.Println("区块链模块初始化成功（Fabric模式）")
		}
	} else {
		// 使用模拟区块链
		blockchain.InitMockBlockchain()
		log.Println("区块链模块初始化成功（模拟模式）")
	}

	// 设置Gin模式
	if config.GlobalConfig.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := routes.SetupRouter()

	// 启动服务
	addr := ":" + config.GlobalConfig.App.Port
	log.Printf("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
