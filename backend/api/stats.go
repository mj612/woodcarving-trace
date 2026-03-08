package api

import (
	"net/http"
	"woodcarving-backend/models"

	"github.com/gin-gonic/gin"
)

// GetStats 获取统计数据
func GetStats(c *gin.Context) {
	var materialCount int64
	var productCount int64
	var transactionCount int64
	var userCount int64

	// 统计原料数量
	models.DB.Model(&models.Material{}).Count(&materialCount)

	// 统计产品数量
	models.DB.Model(&models.Product{}).Count(&productCount)

	// 统计交易记录数量
	models.DB.Model(&models.Transaction{}).Count(&transactionCount)

	// 统计用户数量
	models.DB.Model(&models.User{}).Count(&userCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"materialCount":    materialCount,
			"productCount":     productCount,
			"transactionCount": transactionCount,
			"userCount":        userCount,
		},
	})
}

// GetRecentActivities 获取最近活动
func GetRecentActivities(c *gin.Context) {
	limit := parseInt(c.DefaultQuery("limit", "10"))

	var transactions []models.Transaction
	models.DB.Order("created_at DESC").
		Limit(limit).
		Find(&transactions)

	var activities []map[string]interface{}
	for _, tx := range transactions {
		// 查询用户信息
		var user models.User
		models.DB.First(&user, tx.UserID)

		activity := map[string]interface{}{
			"time":   tx.CreatedAt.Format("2006-01-02 15:04:05"),
			"action": getActionText(tx.Action, tx.AssetType, tx.AssetID),
			"user":   getUserName(user),
		}
		activities = append(activities, activity)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"activities": activities,
		},
	})
}

func getActionText(action, assetType, assetID string) string {
	actionMap := map[string]string{
		"create":   "创建",
		"transfer": "转移",
		"storage":  "仓储操作",
		"sales":    "销售",
	}

	typeMap := map[string]string{
		"material": "原料",
		"product":  "产品",
	}

	actionText := actionMap[action]
	if actionText == "" {
		actionText = action
	}

	typeText := typeMap[assetType]
	if typeText == "" {
		typeText = assetType
	}

	return actionText + typeText + "：" + assetID
}

func getUserName(user models.User) string {
	if user.RealName != "" {
		return user.RealName
	}
	return user.Username
}
