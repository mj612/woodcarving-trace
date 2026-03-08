package routes

import (
	"woodcarving-backend/api"
	"woodcarving-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 添加中间件
	r.Use(middleware.CORSMiddleware())

	// 静态文件服务（用于访问上传的文件）
	r.Static("/uploads", "./uploads")

	// API v1路由组
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		public := v1.Group("")
		{
			// 用户相关
			public.POST("/register", api.Register)
			public.POST("/login", api.Login)

			// 溯源查询（公开）
			public.GET("/trace/:id", api.GetTrace)
			public.GET("/history/:id", api.GetHistory)
		}

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			auth.GET("/user/info", api.GetUserInfo)
			auth.PUT("/user/info", api.UpdateUserInfo)
			auth.PUT("/user/password", api.ChangePassword)

			// 统计数据
			auth.GET("/stats", api.GetStats)
			auth.GET("/activities/recent", api.GetRecentActivities)

			// 文件上传
			auth.POST("/upload", api.UploadFile)

			// 原料管理
			materials := auth.Group("/materials")
			{
				// 创建原料（仅供应商和监管者）
				materials.POST("", middleware.RoleMiddleware("supplier", "supervisor"), api.CreateMaterial)
				// 查看原料列表（供应商、工匠、监管者都可以查看）
				materials.GET("", middleware.RoleMiddleware("supplier", "artisan", "supervisor"), api.GetMaterialList)
				// 查看原料详情（供应商、工匠、监管者都可以查看）
				materials.GET("/:id", middleware.RoleMiddleware("supplier", "artisan", "supervisor"), api.GetMaterialDetail)
				// 转移原料（仅供应商和监管者）
				materials.POST("/transfer", middleware.RoleMiddleware("supplier", "supervisor"), api.TransferMaterial)
			}

			// 产品管理
			products := auth.Group("/products")
			{
				// 创建产品（仅工匠和监管者）
				products.POST("", middleware.RoleMiddleware("artisan", "supervisor"), api.CreateProduct)
				// 查看产品列表（所有角色都可以查看）
				products.GET("", api.GetProductList)
				// 查看产品详情（所有角色都可以查看）
				products.GET("/:id", api.GetProductDetail)
				// 转移产品（工匠、仓管、销售商、监管者）
				products.POST("/transfer", middleware.RoleMiddleware("artisan", "warehouse", "seller", "supervisor"), api.TransferProduct)
			}

			// 仓储管理（仓管）
			storage := auth.Group("/storage")
			storage.Use(middleware.RoleMiddleware("warehouse", "supervisor"))
			{
				storage.POST("/record", api.RecordStorage)
			}

			// 销售管理（销售商）
			sales := auth.Group("/sales")
			sales.Use(middleware.RoleMiddleware("seller", "supervisor"))
			{
				sales.POST("/record", api.RecordSales)
			}

			// 管理后台（监管方）
			admin := auth.Group("/admin")
			admin.Use(middleware.RoleMiddleware("supervisor"))
			{
				admin.GET("/users", api.GetUserList)
				admin.PUT("/users/status", api.UpdateUserStatus)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
