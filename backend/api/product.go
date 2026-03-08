package api

import (
	"net/http"
	"woodcarving-backend/models"
	"woodcarving-backend/pkg/blockchain"
	"woodcarving-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateProduct 创建产品
func CreateProduct(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req struct {
		ProductName string  `json:"productName" binding:"required"`
		MaterialID  string  `json:"materialId" binding:"required"`
		Dimensions  string  `json:"dimensions" binding:"required"`
		Weight      float64 `json:"weight" binding:"required"`
		CraftDesc   string  `json:"craftDesc"`
		CarveTime   int     `json:"carveTime"`
		Images      string  `json:"images"`
		DesignFile  string  `json:"designFile"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 生成产品ID
	productID := utils.GenerateID("PRD")

	// 计算哈希
	designHash := utils.CalculateStringHash(req.DesignFile)
	imageHash := utils.CalculateStringHash(req.Images)

	// 上链
	txID, err := blockchain.CreateProduct(
		productID,
		req.ProductName,
		req.MaterialID,
		username.(string),
		username.(string),
		req.Dimensions,
		req.Weight,
		req.CraftDesc,
		designHash,
		imageHash,
		req.CarveTime,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上链失败: " + err.Error()})
		return
	}

	// 保存元数据到数据库
	product := models.Product{
		ProductID:   productID,
		UserID:      userID.(uint),
		Images:      req.Images,
		DesignFile:  req.DesignFile,
		Description: req.Description,
		TxID:        txID,
	}

	if err := models.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存数据失败"})
		return
	}

	// 记录交易
	tx := models.Transaction{
		TxID:      txID,
		AssetID:   productID,
		AssetType: "product",
		Action:    "create",
		UserID:    userID.(uint),
		Detail:    toJSON(req),
	}
	models.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "产品创建成功",
		"data": gin.H{
			"productId": productID,
			"txId":      txID,
		},
	})
}

// GetProductList 获取产品列表
func GetProductList(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	var products []models.Product
	query := models.DB.Model(&models.Product{})

	// 所有用户都可以查看所有产品
	// 管理员、仓管、销售商可以看到所有产品
	// 其他角色也可以看到所有产品（用于协作）

	var total int64
	query.Count(&total)

	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&products)

	// 增强数据（添加链上信息）
	var result []map[string]interface{}
	for _, p := range products {
		// 从链上查询
		chainData, err := blockchain.GetProductByID(p.ProductID)
		if err != nil {
			chainData = map[string]interface{}{}
		}

		// 如果指定了状态过滤，检查状态
		if status != "" {
			if statusVal, ok := chainData["status"].(string); ok && statusVal != status {
				continue
			}
		}

		item := map[string]interface{}{
			"id":          p.ID,
			"productId":   p.ProductID,
			"images":      p.Images,
			"designFile":  p.DesignFile,
			"description": p.Description,
			"qrCode":      p.QRCode,
			"txId":        p.TxID,
			"createdAt":   p.CreatedAt,
			"chainData":   chainData,
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  result,
			"total": total,
		},
	})
}

// GetProductDetail 获取产品详情
func GetProductDetail(c *gin.Context) {
	productID := c.Param("id")

	// 从数据库查询
	var product models.Product
	if err := models.DB.Where("product_id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "产品不存在"})
		return
	}

	// 从链上查询
	chainData, err := blockchain.GetProductByID(productID)
	if err != nil {
		// 如果链上数据不存在，返回空对象而不是错误
		chainData = map[string]interface{}{
			"productId": productID,
			"status":    "unknown",
		}
	}

	// 查询流转记录
	records, _ := blockchain.GetTransferRecords(productID)
	if records == nil {
		records = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"product":   product,
			"chainData": chainData,
			"records":   records,
		},
	})
}

// TransferProduct 转移产品
func TransferProduct(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		ProductID    string `json:"productId" binding:"required"`
		ToUser       string `json:"toUser" binding:"required"`
		ToName       string `json:"toName" binding:"required"`
		TransferType string `json:"transferType" binding:"required"` // store_in, store_out, sell
		Location     string `json:"location"`
		Remarks      string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 上链
	txID, err := blockchain.TransferProduct(
		req.ProductID,
		req.ToUser,
		req.ToName,
		req.TransferType,
		req.Location,
		req.Remarks,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "转移失败: " + err.Error()})
		return
	}

	// 记录交易
	tx := models.Transaction{
		TxID:      txID,
		AssetID:   req.ProductID,
		AssetType: "product",
		Action:    "transfer",
		UserID:    userID.(uint),
		Detail:    toJSON(req),
	}
	models.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "转移成功",
		"data": gin.H{
			"txId": txID,
		},
	})
}

// RecordStorage 记录仓储信息
func RecordStorage(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req struct {
		ProductID     string  `json:"productId" binding:"required"`
		OperationType string  `json:"operationType" binding:"required"` // in, out
		WarehouseID   string  `json:"warehouseId" binding:"required"`
		Location      string  `json:"location"`
		Temperature   float64 `json:"temperature"`
		Humidity      float64 `json:"humidity"`
		Remarks       string  `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 上链
	txID, err := blockchain.RecordStorage(
		req.ProductID,
		req.OperationType,
		req.WarehouseID,
		req.Location,
		req.Temperature,
		req.Humidity,
		username.(string),
		username.(string),
		req.Remarks,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "记录失败: " + err.Error()})
		return
	}

	// 记录交易
	tx := models.Transaction{
		TxID:      txID,
		AssetID:   req.ProductID,
		AssetType: "product",
		Action:    "storage",
		UserID:    userID.(uint),
		Detail:    toJSON(req),
	}
	models.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "记录成功",
		"data": gin.H{
			"txId": txID,
		},
	})
}

// RecordSales 记录销售信息
func RecordSales(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req struct {
		ProductID    string  `json:"productId" binding:"required"`
		OrderID      string  `json:"orderId" binding:"required"`
		BuyerName    string  `json:"buyerName" binding:"required"`
		BuyerContact string  `json:"buyerContact"`
		Price        float64 `json:"price" binding:"required"`
		SaleDate     string  `json:"saleDate"`
		TrackingNo   string  `json:"trackingNo"`
		Status       string  `json:"status"` // ordered, shipped, delivered
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 上链
	txID, err := blockchain.RecordSales(
		req.ProductID,
		req.OrderID,
		username.(string),
		username.(string),
		req.BuyerName,
		req.BuyerContact,
		req.Price,
		req.SaleDate,
		req.TrackingNo,
		req.Status,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "记录失败: " + err.Error()})
		return
	}

	// 记录交易
	tx := models.Transaction{
		TxID:      txID,
		AssetID:   req.ProductID,
		AssetType: "product",
		Action:    "sales",
		UserID:    userID.(uint),
		Detail:    toJSON(req),
	}
	models.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "记录成功",
		"data": gin.H{
			"txId": txID,
		},
	})
}
