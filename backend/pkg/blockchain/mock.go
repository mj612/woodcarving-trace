package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MockBlockchain 模拟区块链存储
type MockBlockchain struct {
	materials map[string]*RawMaterial
	products  map[string]*Product
	records   map[string][]TransferRecord
	history   map[string][]HistoryRecord
	mu        sync.RWMutex
}

var mockChain *MockBlockchain

// RawMaterial 原料结构
type RawMaterial struct {
	MaterialID   string    `json:"materialId"`
	WoodType     string    `json:"woodType"`
	Origin       string    `json:"origin"`
	HarvestCert  string    `json:"harvestCert"`
	SupplierID   string    `json:"supplierId"`
	SupplierName string    `json:"supplierName"`
	Quantity     float64   `json:"quantity"`
	Quality      string    `json:"quality"`
	CertHash     string    `json:"certHash"`
	Status       string    `json:"status"`
	CurrentOwner string    `json:"currentOwner"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// Product 产品结构
type Product struct {
	ProductID    string    `json:"productId"`
	ProductName  string    `json:"productName"`
	MaterialID   string    `json:"materialId"`
	ArtisanID    string    `json:"artisanId"`
	ArtisanName  string    `json:"artisanName"`
	Dimensions   string    `json:"dimensions"`
	Weight       float64   `json:"weight"`
	CraftDesc    string    `json:"craftDesc"`
	DesignHash   string    `json:"designHash"`
	ImageHash    string    `json:"imageHash"`
	CarveTime    int       `json:"carveTime"`
	Status       string    `json:"status"`
	CurrentOwner string    `json:"currentOwner"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TransferRecord 流转记录
type TransferRecord struct {
	AssetID      string    `json:"assetId"`
	AssetType    string    `json:"assetType"`
	FromOwner    string    `json:"fromOwner"`
	ToOwner      string    `json:"toOwner"`
	ToName       string    `json:"toName"`
	TransferType string    `json:"transferType"`
	Location     string    `json:"location"`
	Remarks      string    `json:"remarks"`
	Timestamp    time.Time `json:"timestamp"`
	TxID         string    `json:"txId"`
}

// HistoryRecord 历史记录
type HistoryRecord struct {
	TxID      string                 `json:"txId"`
	Timestamp time.Time              `json:"timestamp"`
	Action    string                 `json:"action"`
	Data      map[string]interface{} `json:"data"`
}

// StorageRecord 仓储记录
type StorageRecord struct {
	ProductID     string    `json:"productId"`
	OperationType string    `json:"operationType"`
	WarehouseID   string    `json:"warehouseId"`
	Location      string    `json:"location"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	OperatorID    string    `json:"operatorId"`
	OperatorName  string    `json:"operatorName"`
	Remarks       string    `json:"remarks"`
	Timestamp     time.Time `json:"timestamp"`
	TxID          string    `json:"txId"`
}

// SalesRecord 销售记录
type SalesRecord struct {
	ProductID    string    `json:"productId"`
	OrderID      string    `json:"orderId"`
	SellerID     string    `json:"sellerId"`
	SellerName   string    `json:"sellerName"`
	BuyerName    string    `json:"buyerName"`
	BuyerContact string    `json:"buyerContact"`
	Price        float64   `json:"price"`
	SaleDate     string    `json:"saleDate"`
	TrackingNo   string    `json:"trackingNo"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
	TxID         string    `json:"txId"`
}

// InitMockBlockchain 初始化模拟区块链
func InitMockBlockchain() {
	mockChain = &MockBlockchain{
		materials: make(map[string]*RawMaterial),
		products:  make(map[string]*Product),
		records:   make(map[string][]TransferRecord),
		history:   make(map[string][]HistoryRecord),
	}
}

// generateTxID 生成交易ID
func generateTxID() string {
	data := fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}

// CreateRawMaterial 创建原料
func CreateRawMaterial(materialID, woodType, origin, harvestCert, supplierID, supplierName string,
	quantity float64, quality, certHash string) (string, error) {

	if mockChain == nil {
		InitMockBlockchain()
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	txID := generateTxID()

	material := &RawMaterial{
		MaterialID:   materialID,
		WoodType:     woodType,
		Origin:       origin,
		HarvestCert:  harvestCert,
		SupplierID:   supplierID,
		SupplierName: supplierName,
		Quantity:     quantity,
		Quality:      quality,
		CertHash:     certHash,
		Status:       "available",
		CurrentOwner: supplierID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChain.materials[materialID] = material

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "create",
		Data: map[string]interface{}{
			"materialId": materialID,
			"woodType":   woodType,
			"origin":     origin,
		},
	}
	mockChain.history[materialID] = append(mockChain.history[materialID], history)

	return txID, nil
}

// CreateProduct 创建产品
func CreateProduct(productID, productName, materialID, artisanID, artisanName, dimensions string,
	weight float64, craftDesc, designHash, imageHash string, carveTime int) (string, error) {

	if mockChain == nil {
		InitMockBlockchain()
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	txID := generateTxID()

	product := &Product{
		ProductID:    productID,
		ProductName:  productName,
		MaterialID:   materialID,
		ArtisanID:    artisanID,
		ArtisanName:  artisanName,
		Dimensions:   dimensions,
		Weight:       weight,
		CraftDesc:    craftDesc,
		DesignHash:   designHash,
		ImageHash:    imageHash,
		CarveTime:    carveTime,
		Status:       "produced",
		CurrentOwner: artisanID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockChain.products[productID] = product

	// 更新原料状态
	if material, exists := mockChain.materials[materialID]; exists {
		material.Status = "used"
		material.UpdatedAt = time.Now()
	}

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "create",
		Data: map[string]interface{}{
			"productId":   productID,
			"productName": productName,
			"materialId":  materialID,
		},
	}
	mockChain.history[productID] = append(mockChain.history[productID], history)

	return txID, nil
}

// TransferMaterial 转移原料
func TransferMaterial(materialID, toOwner, toName, location, remarks string) (string, error) {
	if mockChain == nil {
		return "", fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	material, exists := mockChain.materials[materialID]
	if !exists {
		return "", fmt.Errorf("原料不存在")
	}

	txID := generateTxID()

	// 记录流转
	record := TransferRecord{
		AssetID:      materialID,
		AssetType:    "material",
		FromOwner:    material.CurrentOwner,
		ToOwner:      toOwner,
		ToName:       toName,
		TransferType: "transfer",
		Location:     location,
		Remarks:      remarks,
		Timestamp:    time.Now(),
		TxID:         txID,
	}
	mockChain.records[materialID] = append(mockChain.records[materialID], record)

	// 更新所有者
	material.CurrentOwner = toOwner
	material.Status = "transferred"
	material.UpdatedAt = time.Now()

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "transfer",
		Data: map[string]interface{}{
			"materialId": materialID,
			"toOwner":    toOwner,
			"toName":     toName,
		},
	}
	mockChain.history[materialID] = append(mockChain.history[materialID], history)

	return txID, nil
}

// TransferProduct 转移产品
func TransferProduct(productID, toOwner, toName, transferType, location, remarks string) (string, error) {
	if mockChain == nil {
		return "", fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	product, exists := mockChain.products[productID]
	if !exists {
		return "", fmt.Errorf("产品不存在")
	}

	txID := generateTxID()

	// 记录流转
	record := TransferRecord{
		AssetID:      productID,
		AssetType:    "product",
		FromOwner:    product.CurrentOwner,
		ToOwner:      toOwner,
		ToName:       toName,
		TransferType: transferType,
		Location:     location,
		Remarks:      remarks,
		Timestamp:    time.Now(),
		TxID:         txID,
	}
	mockChain.records[productID] = append(mockChain.records[productID], record)

	// 更新所有者和状态
	product.CurrentOwner = toOwner
	switch transferType {
	case "store_in":
		product.Status = "in_storage"
	case "store_out":
		product.Status = "out_storage"
	case "sell":
		product.Status = "sold"
	}
	product.UpdatedAt = time.Now()

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "transfer",
		Data: map[string]interface{}{
			"productId":    productID,
			"toOwner":      toOwner,
			"toName":       toName,
			"transferType": transferType,
		},
	}
	mockChain.history[productID] = append(mockChain.history[productID], history)

	return txID, nil
}

// RecordStorage 记录仓储信息
func RecordStorage(productID, operationType, warehouseID, location string,
	temperature, humidity float64, operatorID, operatorName, remarks string) (string, error) {

	if mockChain == nil {
		return "", fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	product, exists := mockChain.products[productID]
	if !exists {
		return "", fmt.Errorf("产品不存在")
	}

	txID := generateTxID()

	// 更新产品状态
	if operationType == "in" {
		product.Status = "in_storage"
	} else if operationType == "out" {
		product.Status = "out_storage"
	}
	product.UpdatedAt = time.Now()

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "storage",
		Data: map[string]interface{}{
			"productId":     productID,
			"operationType": operationType,
			"warehouseId":   warehouseID,
			"location":      location,
			"temperature":   temperature,
			"humidity":      humidity,
			"operatorName":  operatorName,
		},
	}
	mockChain.history[productID] = append(mockChain.history[productID], history)

	return txID, nil
}

// RecordSales 记录销售信息
func RecordSales(productID, orderID, sellerID, sellerName, buyerName, buyerContact string,
	price float64, saleDate, trackingNo, status string) (string, error) {

	if mockChain == nil {
		return "", fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.Lock()
	defer mockChain.mu.Unlock()

	product, exists := mockChain.products[productID]
	if !exists {
		return "", fmt.Errorf("产品不存在")
	}

	txID := generateTxID()

	// 更新产品状态
	product.Status = "sold"
	product.UpdatedAt = time.Now()

	// 记录历史
	history := HistoryRecord{
		TxID:      txID,
		Timestamp: time.Now(),
		Action:    "sales",
		Data: map[string]interface{}{
			"productId":    productID,
			"orderId":      orderID,
			"sellerName":   sellerName,
			"buyerName":    buyerName,
			"buyerContact": buyerContact,
			"price":        price,
			"saleDate":     saleDate,
			"trackingNo":   trackingNo,
			"status":       status,
		},
	}
	mockChain.history[productID] = append(mockChain.history[productID], history)

	return txID, nil
}

// GetMaterialByID 查询原料
func GetMaterialByID(materialID string) (map[string]interface{}, error) {
	if mockChain == nil {
		return nil, fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.RLock()
	defer mockChain.mu.RUnlock()

	material, exists := mockChain.materials[materialID]
	if !exists {
		return nil, fmt.Errorf("原料不存在")
	}

	data, _ := json.Marshal(material)
	var result map[string]interface{}
	json.Unmarshal(data, &result)

	return result, nil
}

// GetProductByID 查询产品
func GetProductByID(productID string) (map[string]interface{}, error) {
	if mockChain == nil {
		return nil, fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.RLock()
	defer mockChain.mu.RUnlock()

	product, exists := mockChain.products[productID]
	if !exists {
		return nil, fmt.Errorf("产品不存在")
	}

	data, _ := json.Marshal(product)
	var result map[string]interface{}
	json.Unmarshal(data, &result)

	return result, nil
}

// GetCompleteTrace 获取完整溯源信息
func GetCompleteTrace(productID string) (map[string]interface{}, error) {
	if mockChain == nil {
		return nil, fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.RLock()
	defer mockChain.mu.RUnlock()

	product, exists := mockChain.products[productID]
	if !exists {
		return nil, fmt.Errorf("产品不存在")
	}

	// 获取原料信息
	var material *RawMaterial
	if mat, exists := mockChain.materials[product.MaterialID]; exists {
		material = mat
	}

	// 获取流转记录
	records := mockChain.records[productID]

	// 获取历史记录
	history := mockChain.history[productID]

	trace := map[string]interface{}{
		"product":  product,
		"material": material,
		"records":  records,
		"history":  history,
	}

	return trace, nil
}

// QueryHistory 查询资产历史
func QueryHistory(assetID string) ([]map[string]interface{}, error) {
	if mockChain == nil {
		return nil, fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.RLock()
	defer mockChain.mu.RUnlock()

	history := mockChain.history[assetID]
	if history == nil {
		return []map[string]interface{}{}, nil
	}

	var result []map[string]interface{}
	for _, h := range history {
		data, _ := json.Marshal(h)
		var item map[string]interface{}
		json.Unmarshal(data, &item)
		result = append(result, item)
	}

	return result, nil
}

// GetTransferRecords 获取流转记录
func GetTransferRecords(assetID string) ([]map[string]interface{}, error) {
	if mockChain == nil {
		return nil, fmt.Errorf("区块链未初始化")
	}

	mockChain.mu.RLock()
	defer mockChain.mu.RUnlock()

	records := mockChain.records[assetID]
	if records == nil {
		return []map[string]interface{}{}, nil
	}

	var result []map[string]interface{}
	for _, r := range records {
		data, _ := json.Marshal(r)
		var item map[string]interface{}
		json.Unmarshal(data, &item)
		result = append(result, item)
	}

	return result, nil
}

// InitFabric Fabric初始化占位函数
// 当需要使用Fabric时，将 fabric_impl.go.bak 重命名为 fabric.go
func InitFabric() error {
	return fmt.Errorf("Fabric功能未启用，请部署Fabric网络并重命名 fabric_impl.go.bak 为 fabric.go")
}
