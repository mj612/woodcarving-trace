package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TraceContract 溯源智能合约
type TraceContract struct {
	contractapi.Contract
}

// RawMaterial 原料结构
type RawMaterial struct {
	MaterialID    string    `json:"materialId"`    // 原料ID
	WoodType      string    `json:"woodType"`      // 木材种类
	Origin        string    `json:"origin"`        // 产地
	HarvestCert   string    `json:"harvestCert"`   // 采伐证编号
	SupplierID    string    `json:"supplierId"`    // 供应商ID
	SupplierName  string    `json:"supplierName"`  // 供应商名称
	Quantity      float64   `json:"quantity"`      // 数量（立方米）
	Quality       string    `json:"quality"`       // 质量等级
	CertHash      string    `json:"certHash"`      // 证书文件哈希
	Status        string    `json:"status"`        // 状态：available, transferred, used
	CurrentOwner  string    `json:"currentOwner"`  // 当前持有者ID
	CreatedAt     string    `json:"createdAt"`     // 创建时间
	UpdatedAt     string    `json:"updatedAt"`     // 更新时间
	AssetType     string    `json:"assetType"`     // 资产类型标识
}

// Product 成品结构
type Product struct {
	ProductID     string    `json:"productId"`     // 产品ID
	ProductName   string    `json:"productName"`   // 产品名称
	MaterialID    string    `json:"materialId"`    // 关联原料ID
	ArtisanID     string    `json:"artisanId"`     // 雕刻师ID
	ArtisanName   string    `json:"artisanName"`   // 雕刻师姓名
	Dimensions    string    `json:"dimensions"`    // 尺寸
	Weight        float64   `json:"weight"`        // 重量（kg）
	CraftDesc     string    `json:"craftDesc"`     // 工艺描述
	DesignHash    string    `json:"designHash"`    // 设计图哈希
	ImageHash     string    `json:"imageHash"`     // 成品图片哈希
	CarveTime     int       `json:"carveTime"`     // 雕刻耗时（小时）
	Status        string    `json:"status"`        // 状态：produced, in_storage, in_transit, sold
	CurrentOwner  string    `json:"currentOwner"`  // 当前持有者ID
	Location      string    `json:"location"`      // 当前位置
	CreatedAt     string    `json:"createdAt"`     // 创建时间
	UpdatedAt     string    `json:"updatedAt"`     // 更新时间
	AssetType     string    `json:"assetType"`     // 资产类型标识
}

// TransferRecord 流转记录
type TransferRecord struct {
	RecordID      string    `json:"recordId"`      // 记录ID
	AssetID       string    `json:"assetId"`       // 资产ID（原料或产品）
	AssetType     string    `json:"assetType"`     // 资产类型：material, product
	FromOwner     string    `json:"fromOwner"`     // 转出方ID
	FromName      string    `json:"fromName"`      // 转出方名称
	ToOwner       string    `json:"toOwner"`       // 接收方ID
	ToName        string    `json:"toName"`        // 接收方名称
	TransferType  string    `json:"transferType"`  // 转移类型：supply, produce, store_in, store_out, sell
	Timestamp     string    `json:"timestamp"`     // 时间戳
	Location      string    `json:"location"`      // 位置
	Remarks       string    `json:"remarks"`       // 备注
	TxID          string    `json:"txId"`          // 交易ID
}

// StorageRecord 仓储记录
type StorageRecord struct {
	RecordID      string    `json:"recordId"`      // 记录ID
	ProductID     string    `json:"productId"`     // 产品ID
	OperationType string    `json:"operationType"` // 操作类型：in, out
	WarehouseID   string    `json:"warehouseId"`   // 仓库ID
	Location      string    `json:"location"`      // 库位
	Temperature   float64   `json:"temperature"`   // 温度
	Humidity      float64   `json:"humidity"`      // 湿度
	OperatorID    string    `json:"operatorId"`    // 操作员ID
	OperatorName  string    `json:"operatorName"`  // 操作员姓名
	Timestamp     string    `json:"timestamp"`     // 时间戳
	Remarks       string    `json:"remarks"`       // 备注
}

// SalesRecord 销售记录
type SalesRecord struct {
	RecordID      string    `json:"recordId"`      // 记录ID
	ProductID     string    `json:"productId"`     // 产品ID
	OrderID       string    `json:"orderId"`       // 订单ID
	SellerID      string    `json:"sellerId"`      // 销售商ID
	SellerName    string    `json:"sellerName"`    // 销售商名称
	BuyerName     string    `json:"buyerName"`     // 购买者姓名
	BuyerContact  string    `json:"buyerContact"`  // 购买者联系方式
	Price         float64   `json:"price"`         // 价格
	SaleDate      string    `json:"saleDate"`      // 销售日期
	TrackingNo    string    `json:"trackingNo"`    // 物流单号
	Status        string    `json:"status"`        // 状态：ordered, shipped, delivered
	Timestamp     string    `json:"timestamp"`     // 时间戳
}

// Init 初始化合约
func (t *TraceContract) Init(ctx contractapi.TransactionContextInterface) error {
	return nil
}

// CreateRawMaterial 创建原料记录
func (t *TraceContract) CreateRawMaterial(ctx contractapi.TransactionContextInterface, 
	materialID, woodType, origin, harvestCert, supplierID, supplierName string, 
	quantity float64, quality, certHash string) error {
	
	// 检查原料ID是否已存在
	exists, err := t.AssetExists(ctx, materialID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("原料ID %s 已存在", materialID)
	}

	timestamp := time.Now().Format(time.RFC3339)
	
	material := RawMaterial{
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
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
		AssetType:    "material",
	}

	materialJSON, err := json.Marshal(material)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(materialID, materialJSON)
	if err != nil {
		return fmt.Errorf("写入原料数据失败: %v", err)
	}

	// 记录流转
	txID := ctx.GetStub().GetTxID()
	record := TransferRecord{
		RecordID:     fmt.Sprintf("TR_%s_%s", materialID, timestamp),
		AssetID:      materialID,
		AssetType:    "material",
		FromOwner:    "",
		FromName:     "",
		ToOwner:      supplierID,
		ToName:       supplierName,
		TransferType: "supply",
		Timestamp:    timestamp,
		Location:     origin,
		Remarks:      "原料上链",
		TxID:         txID,
	}

	return t.createTransferRecord(ctx, record)
}

// CreateProduct 创建成品记录
func (t *TraceContract) CreateProduct(ctx contractapi.TransactionContextInterface,
	productID, productName, materialID, artisanID, artisanName, dimensions string,
	weight float64, craftDesc, designHash, imageHash string, carveTime int) error {
	
	// 检查产品ID是否已存在
	exists, err := t.AssetExists(ctx, productID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("产品ID %s 已存在", productID)
	}

	// 检查原料是否存在
	materialJSON, err := ctx.GetStub().GetState(materialID)
	if err != nil {
		return fmt.Errorf("读取原料数据失败: %v", err)
	}
	if materialJSON == nil {
		return fmt.Errorf("原料ID %s 不存在", materialID)
	}

	// 更新原料状态
	var material RawMaterial
	err = json.Unmarshal(materialJSON, &material)
	if err != nil {
		return err
	}
	
	material.Status = "used"
	material.UpdatedAt = time.Now().Format(time.RFC3339)
	
	materialJSON, _ = json.Marshal(material)
	ctx.GetStub().PutState(materialID, materialJSON)

	// 创建产品
	timestamp := time.Now().Format(time.RFC3339)
	
	product := Product{
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
		Location:     "工坊",
		CreatedAt:    timestamp,
		UpdatedAt:    timestamp,
		AssetType:    "product",
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(productID, productJSON)
	if err != nil {
		return fmt.Errorf("写入产品数据失败: %v", err)
	}

	// 记录流转
	txID := ctx.GetStub().GetTxID()
	record := TransferRecord{
		RecordID:     fmt.Sprintf("TR_%s_%s", productID, timestamp),
		AssetID:      productID,
		AssetType:    "product",
		FromOwner:    "",
		FromName:     "",
		ToOwner:      artisanID,
		ToName:       artisanName,
		TransferType: "produce",
		Timestamp:    timestamp,
		Location:     "工坊",
		Remarks:      fmt.Sprintf("使用原料 %s 生产完成", materialID),
		TxID:         txID,
	}

	return t.createTransferRecord(ctx, record)
}

// TransferMaterial 转移原料
func (t *TraceContract) TransferMaterial(ctx contractapi.TransactionContextInterface,
	materialID, toOwner, toName, location, remarks string) error {
	
	materialJSON, err := ctx.GetStub().GetState(materialID)
	if err != nil {
		return fmt.Errorf("读取原料数据失败: %v", err)
	}
	if materialJSON == nil {
		return fmt.Errorf("原料ID %s 不存在", materialID)
	}

	var material RawMaterial
	err = json.Unmarshal(materialJSON, &material)
	if err != nil {
		return err
	}

	fromOwner := material.CurrentOwner
	fromName := material.SupplierName

	// 更新原料状态
	material.CurrentOwner = toOwner
	material.Status = "transferred"
	material.UpdatedAt = time.Now().Format(time.RFC3339)

	materialJSON, _ = json.Marshal(material)
	ctx.GetStub().PutState(materialID, materialJSON)

	// 记录流转
	timestamp := time.Now().Format(time.RFC3339)
	txID := ctx.GetStub().GetTxID()
	record := TransferRecord{
		RecordID:     fmt.Sprintf("TR_%s_%s", materialID, timestamp),
		AssetID:      materialID,
		AssetType:    "material",
		FromOwner:    fromOwner,
		FromName:     fromName,
		ToOwner:      toOwner,
		ToName:       toName,
		TransferType: "supply",
		Timestamp:    timestamp,
		Location:     location,
		Remarks:      remarks,
		TxID:         txID,
	}

	return t.createTransferRecord(ctx, record)
}

// TransferProduct 转移产品
func (t *TraceContract) TransferProduct(ctx contractapi.TransactionContextInterface,
	productID, toOwner, toName, transferType, location, remarks string) error {
	
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return fmt.Errorf("读取产品数据失败: %v", err)
	}
	if productJSON == nil {
		return fmt.Errorf("产品ID %s 不存在", productID)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return err
	}

	fromOwner := product.CurrentOwner
	fromName := product.ArtisanName

	// 更新产品状态
	product.CurrentOwner = toOwner
	product.Location = location
	
	switch transferType {
	case "store_in":
		product.Status = "in_storage"
	case "store_out":
		product.Status = "in_transit"
	case "sell":
		product.Status = "sold"
	}
	
	product.UpdatedAt = time.Now().Format(time.RFC3339)

	productJSON, _ = json.Marshal(product)
	ctx.GetStub().PutState(productID, productJSON)

	// 记录流转
	timestamp := time.Now().Format(time.RFC3339)
	txID := ctx.GetStub().GetTxID()
	record := TransferRecord{
		RecordID:     fmt.Sprintf("TR_%s_%s", productID, timestamp),
		AssetID:      productID,
		AssetType:    "product",
		FromOwner:    fromOwner,
		FromName:     fromName,
		ToOwner:      toOwner,
		ToName:       toName,
		TransferType: transferType,
		Timestamp:    timestamp,
		Location:     location,
		Remarks:      remarks,
		TxID:         txID,
	}

	return t.createTransferRecord(ctx, record)
}

// RecordStorage 记录仓储信息
func (t *TraceContract) RecordStorage(ctx contractapi.TransactionContextInterface,
	productID, operationType, warehouseID, location string,
	temperature, humidity float64,
	operatorID, operatorName, remarks string) error {
	
	// 检查产品是否存在
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return fmt.Errorf("读取产品数据失败: %v", err)
	}
	if productJSON == nil {
		return fmt.Errorf("产品ID %s 不存在", productID)
	}

	timestamp := time.Now().Format(time.RFC3339)
	recordID := fmt.Sprintf("ST_%s_%s", productID, timestamp)

	record := StorageRecord{
		RecordID:      recordID,
		ProductID:     productID,
		OperationType: operationType,
		WarehouseID:   warehouseID,
		Location:      location,
		Temperature:   temperature,
		Humidity:      humidity,
		OperatorID:    operatorID,
		OperatorName:  operatorName,
		Timestamp:     timestamp,
		Remarks:       remarks,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(recordID, recordJSON)
}

// RecordSales 记录销售信息
func (t *TraceContract) RecordSales(ctx contractapi.TransactionContextInterface,
	productID, orderID, sellerID, sellerName, buyerName, buyerContact string,
	price float64, saleDate, trackingNo, status string) error {
	
	// 检查产品是否存在
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return fmt.Errorf("读取产品数据失败: %v", err)
	}
	if productJSON == nil {
		return fmt.Errorf("产品ID %s 不存在", productID)
	}

	timestamp := time.Now().Format(time.RFC3339)
	recordID := fmt.Sprintf("SL_%s_%s", productID, timestamp)

	record := SalesRecord{
		RecordID:     recordID,
		ProductID:    productID,
		OrderID:      orderID,
		SellerID:     sellerID,
		SellerName:   sellerName,
		BuyerName:    buyerName,
		BuyerContact: buyerContact,
		Price:        price,
		SaleDate:     saleDate,
		TrackingNo:   trackingNo,
		Status:       status,
		Timestamp:    timestamp,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(recordID, recordJSON)
}

// GetMaterialByID 根据ID查询原料
func (t *TraceContract) GetMaterialByID(ctx contractapi.TransactionContextInterface, materialID string) (*RawMaterial, error) {
	materialJSON, err := ctx.GetStub().GetState(materialID)
	if err != nil {
		return nil, fmt.Errorf("读取原料数据失败: %v", err)
	}
	if materialJSON == nil {
		return nil, fmt.Errorf("原料ID %s 不存在", materialID)
	}

	var material RawMaterial
	err = json.Unmarshal(materialJSON, &material)
	if err != nil {
		return nil, err
	}

	return &material, nil
}

// GetProductByID 根据ID查询产品
func (t *TraceContract) GetProductByID(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return nil, fmt.Errorf("读取产品数据失败: %v", err)
	}
	if productJSON == nil {
		return nil, fmt.Errorf("产品ID %s 不存在", productID)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// QueryHistory 查询资产历史
func (t *TraceContract) QueryHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]map[string]interface{}, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var history []map[string]interface{}

	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record map[string]interface{}
		record = make(map[string]interface{})
		record["txId"] = response.TxId
		record["timestamp"] = response.Timestamp.AsTime().Format(time.RFC3339)
		record["isDelete"] = response.IsDelete

		if !response.IsDelete {
			var asset interface{}
			json.Unmarshal(response.Value, &asset)
			record["value"] = asset
		}

		history = append(history, record)
	}

	return history, nil
}

// GetTransferRecords 获取资产的所有流转记录
func (t *TraceContract) GetTransferRecords(ctx contractapi.TransactionContextInterface, assetID string) ([]*TransferRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"assetId":"%s"}}`, assetID)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*TransferRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record TransferRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// GetCompleteTrace 获取产品完整溯源信息
func (t *TraceContract) GetCompleteTrace(ctx contractapi.TransactionContextInterface, productID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取产品信息
	product, err := t.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	result["product"] = product

	// 获取原料信息
	if product.MaterialID != "" {
		material, err := t.GetMaterialByID(ctx, product.MaterialID)
		if err == nil {
			result["material"] = material
		}
	}

	// 获取产品流转记录
	productRecords, err := t.GetTransferRecords(ctx, productID)
	if err == nil {
		result["productTransfers"] = productRecords
	}

	// 获取原料流转记录
	if product.MaterialID != "" {
		materialRecords, err := t.GetTransferRecords(ctx, product.MaterialID)
		if err == nil {
			result["materialTransfers"] = materialRecords
		}
	}

	// 获取仓储记录
	storageRecords, err := t.getStorageRecords(ctx, productID)
	if err == nil {
		result["storageRecords"] = storageRecords
	}

	// 获取销售记录
	salesRecords, err := t.getSalesRecords(ctx, productID)
	if err == nil {
		result["salesRecords"] = salesRecords
	}

	return result, nil
}

// AssetExists 检查资产是否存在
func (t *TraceContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("读取数据失败: %v", err)
	}

	return assetJSON != nil, nil
}

// createTransferRecord 创建流转记录（内部方法）
func (t *TraceContract) createTransferRecord(ctx contractapi.TransactionContextInterface, record TransferRecord) error {
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(record.RecordID, recordJSON)
}

// getStorageRecords 获取仓储记录（内部方法）
func (t *TraceContract) getStorageRecords(ctx contractapi.TransactionContextInterface, productID string) ([]*StorageRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"productId":"%s"}}`, productID)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*StorageRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record StorageRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}
		
		// 只添加仓储记录
		if record.ProductID == productID {
			records = append(records, &record)
		}
	}

	return records, nil
}

// getSalesRecords 获取销售记录（内部方法）
func (t *TraceContract) getSalesRecords(ctx contractapi.TransactionContextInterface, productID string) ([]*SalesRecord, error) {
	queryString := fmt.Sprintf(`{"selector":{"productId":"%s"}}`, productID)
	
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*SalesRecord
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var record SalesRecord
		err = json.Unmarshal(queryResponse.Value, &record)
		if err != nil {
			continue
		}
		
		if record.ProductID == productID {
			records = append(records, &record)
		}
	}

	return records, nil
}
