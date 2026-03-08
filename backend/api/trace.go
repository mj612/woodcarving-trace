package api

import (
	"net/http"
	"woodcarving-backend/pkg/blockchain"

	"github.com/gin-gonic/gin"
)

// GetTrace 获取完整溯源信息（公开接口）
func GetTrace(c *gin.Context) {
	productID := c.Param("id")

	// 从链上查询完整溯源
	trace, err := blockchain.GetCompleteTrace(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询溯源信息失败: " + err.Error(),
		})
		return
	}

	// 转换为前端需要的格式
	product := trace["product"]
	material := trace["material"]
	history := trace["history"]

	// 构建前端需要的数据结构
	result := gin.H{
		"productId":   productID,
		"productName": "",
		"createdAt":   "",
		"status":      "unknown",
		"history":     []gin.H{},
	}

	// 填充产品信息
	if product != nil {
		if productMap, ok := product.(map[string]interface{}); ok {
			if name, ok := productMap["productName"].(string); ok {
				result["productName"] = name
			}
			if createdAt, ok := productMap["createdAt"].(string); ok {
				result["createdAt"] = createdAt
			}
			if status, ok := productMap["status"].(string); ok {
				result["status"] = status
			}
		}
	}

	// 填充原料信息
	if material != nil {
		if materialMap, ok := material.(map[string]interface{}); ok {
			result["materialInfo"] = gin.H{
				"materialId": materialMap["materialId"],
				"woodType":   materialMap["woodType"],
				"origin":     materialMap["origin"],
				"supplier":   materialMap["supplierName"],
			}
		}
	}

	// 填充历史记录
	if history != nil {
		if historyList, ok := history.([]interface{}); ok {
			historyRecords := []gin.H{}
			for _, h := range historyList {
				if hMap, ok := h.(map[string]interface{}); ok {
					action := "操作"
					description := ""

					// 根据action类型生成描述
					if actionStr, ok := hMap["action"].(string); ok {
						switch actionStr {
						case "create":
							action = "创建产品"
							description = "产品已创建并上链"
						case "transfer":
							action = "流转记录"
							if data, ok := hMap["data"].(map[string]interface{}); ok {
								if transferType, ok := data["transferType"].(string); ok {
									switch transferType {
									case "store_in":
										action = "入库"
										description = "产品已入库"
									case "store_out":
										action = "出库"
										description = "产品已出库"
									case "sell":
										action = "销售"
										description = "产品已售出"
									}
								}
							}
						case "storage":
							action = "仓储记录"
							if data, ok := hMap["data"].(map[string]interface{}); ok {
								if opType, ok := data["operationType"].(string); ok {
									if opType == "in" {
										description = "产品入库存储"
									} else {
										description = "产品出库"
									}
								}
							}
						case "sales":
							action = "销售记录"
							description = "产品已售出"
						}
					}

					record := gin.H{
						"action":      action,
						"description": description,
						"time":        hMap["timestamp"],
						"operator":    "系统",
					}

					// 尝试获取操作人信息
					if data, ok := hMap["data"].(map[string]interface{}); ok {
						if operatorName, ok := data["operatorName"].(string); ok && operatorName != "" {
							record["operator"] = operatorName
						} else if sellerName, ok := data["sellerName"].(string); ok && sellerName != "" {
							record["operator"] = sellerName
						}
					}

					historyRecords = append(historyRecords, record)
				}
			}
			result["history"] = historyRecords
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetHistory 获取资产历史（公开接口）
func GetHistory(c *gin.Context) {
	assetID := c.Param("id")

	history, err := blockchain.QueryHistory(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": history,
	})
}
