package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"woodcarving-backend/models"
	"woodcarving-backend/pkg/blockchain"
	"woodcarving-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateMaterial 创建原料
func CreateMaterial(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")

	var req struct {
		WoodType    string  `json:"woodType" binding:"required"`
		Origin      string  `json:"origin" binding:"required"`
		HarvestCert string  `json:"harvestCert" binding:"required"`
		Quantity    float64 `json:"quantity" binding:"required"`
		Quality     string  `json:"quality" binding:"required"`
		Images      string  `json:"images"`
		CertFile    string  `json:"certFile"`
		Description string  `json:"description"`
	}

	// 打印接收到的原始数据
	bodyBytes, _ := c.GetRawData()
	println("接收到的数据:", string(bodyBytes))

	// 重新设置body以便后续读取
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		println("参数绑定失败:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误: " + err.Error()})
		return
	}

	// 生成原料ID
	materialID := utils.GenerateID("MAT")

	// 计算证书哈希
	certHash := utils.CalculateStringHash(req.HarvestCert + req.CertFile)

	// 上链
	txID, err := blockchain.CreateRawMaterial(
		materialID,
		req.WoodType,
		req.Origin,
		req.HarvestCert,
		username.(string),
		username.(string),
		req.Quantity,
		req.Quality,
		certHash,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上链失败: " + err.Error()})
		return
	}

	// 保存元数据到数据库
	material := models.Material{
		MaterialID:  materialID,
		UserID:      userID.(uint),
		Images:      req.Images,
		CertFile:    req.CertFile,
		Description: req.Description,
		TxID:        txID,
	}

	if err := models.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存数据失败"})
		return
	}

	// 记录交易
	tx := models.Transaction{
		TxID:      txID,
		AssetID:   materialID,
		AssetType: "material",
		Action:    "create",
		UserID:    userID.(uint),
		Detail:    toJSON(req),
	}
	models.DB.Create(&tx)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "原料创建成功",
		"data": gin.H{
			"materialId": materialID,
			"txId":       txID,
		},
	})
}

// GetMaterialList 获取原料列表
func GetMaterialList(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"))
	pageSize := parseInt(c.DefaultQuery("pageSize", "10"))

	var materials []models.Material
	query := models.DB.Model(&models.Material{})

	// 所有用户都可以查看所有原料
	// 用于协作和溯源查询

	var total int64
	query.Count(&total)

	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&materials)

	// 增强数据（添加链上信息）
	var result []map[string]interface{}
	for _, m := range materials {
		// 从链上查询
		chainData, err := blockchain.GetMaterialByID(m.MaterialID)
		if err != nil {
			chainData = map[string]interface{}{}
		}

		item := map[string]interface{}{
			"id":          m.ID,
			"materialId":  m.MaterialID,
			"images":      m.Images,
			"certFile":    m.CertFile,
			"description": m.Description,
			"txId":        m.TxID,
			"createdAt":   m.CreatedAt,
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

// GetMaterialDetail 获取原料详情
func GetMaterialDetail(c *gin.Context) {
	materialID := c.Param("id")

	// 从数据库查询
	var material models.Material
	if err := models.DB.Where("material_id = ?", materialID).First(&material).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "原料不存在"})
		return
	}

	// 从链上查询
	chainData, err := blockchain.GetMaterialByID(materialID)
	if err != nil {
		// 如果链上数据不存在，返回空对象而不是错误
		chainData = map[string]interface{}{
			"materialId": materialID,
			"status":     "unknown",
		}
	}

	// 查询流转记录
	records, _ := blockchain.GetTransferRecords(materialID)
	if records == nil {
		records = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"material":  material,
			"chainData": chainData,
			"records":   records,
		},
	})
}

// TransferMaterial 转移原料
func TransferMaterial(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		MaterialID string `json:"materialId" binding:"required"`
		ToUser     string `json:"toUser" binding:"required"`
		ToName     string `json:"toName" binding:"required"`
		Location   string `json:"location"`
		Remarks    string `json:"remarks"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 上链
	txID, err := blockchain.TransferMaterial(
		req.MaterialID,
		req.ToUser,
		req.ToName,
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
		AssetID:   req.MaterialID,
		AssetType: "material",
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

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
