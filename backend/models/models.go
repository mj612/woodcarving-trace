package models

import (
	"fmt"
	"log"
	"time"
	"woodcarving-backend/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RealName  string    `gorm:"size:50" json:"realName"`
	Role      string    `gorm:"size:20;not null" json:"role"` // supplier, artisan, warehouse, seller, consumer, supervisor
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:100" json:"email"`
	Company   string    `gorm:"size:100" json:"company"`
	Address   string    `gorm:"size:200" json:"address"`
	Status    int       `gorm:"default:1" json:"status"` // 0:禁用 1:启用 2:待审核
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Material 原料元数据（不上链的部分）
type Material struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MaterialID  string    `gorm:"uniqueIndex;size:50;not null" json:"materialId"`
	UserID      uint      `gorm:"not null" json:"userId"`
	Images      string    `gorm:"type:text" json:"images"` // JSON数组
	CertFile    string    `gorm:"size:200" json:"certFile"`
	Description string    `gorm:"type:text" json:"description"`
	TxID        string    `gorm:"size:100" json:"txId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Product 产品元数据（不上链的部分）
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   string    `gorm:"uniqueIndex;size:50;not null" json:"productId"`
	UserID      uint      `gorm:"not null" json:"userId"`
	Images      string    `gorm:"type:text" json:"images"` // JSON数组
	DesignFile  string    `gorm:"size:200" json:"designFile"`
	Description string    `gorm:"type:text" json:"description"`
	QRCode      string    `gorm:"size:200" json:"qrCode"` // 二维码路径
	TxID        string    `gorm:"size:100" json:"txId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Transaction 交易记录（用于前端显示）
type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TxID      string    `gorm:"uniqueIndex;size:100;not null" json:"txId"`
	AssetID   string    `gorm:"size:50;not null;index" json:"assetId"`
	AssetType string    `gorm:"size:20" json:"assetType"` // material, product
	Action    string    `gorm:"size:50" json:"action"`
	UserID    uint      `json:"userId"`
	Detail    string    `gorm:"type:text" json:"detail"` // JSON格式详情
	CreatedAt time.Time `json:"createdAt"`
}

// InitDB 初始化数据库
func InitDB() {
	cfg := config.GlobalConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(&User{}, &Material{}, &Product{}, &Transaction{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Println("数据库初始化成功")
}
