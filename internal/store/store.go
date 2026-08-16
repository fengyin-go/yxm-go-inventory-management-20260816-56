// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"inventory/internal/model"
)

// 数据访问层常见错误。
var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或发生状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// 商品
	CreateProduct(p *model.Product) error
	GetProduct(id string) (*model.Product, error)
	GetProductBySKU(sku string) (*model.Product, error)
	ListProducts() []*model.Product
	UpdateProduct(p *model.Product) error
	DeleteProduct(id string) error

	// 仓库
	CreateWarehouse(w *model.Warehouse) error
	GetWarehouse(id string) (*model.Warehouse, error)
	GetWarehouseByCode(code string) (*model.Warehouse, error)
	ListWarehouses() []*model.Warehouse
	UpdateWarehouse(w *model.Warehouse) error
	DeleteWarehouse(id string) error

	// 库存
	CreateStockItem(s *model.StockItem) error
	GetStockItem(id string) (*model.StockItem, error)
	GetStockItemByProductWarehouse(productID, warehouseID string) (*model.StockItem, error)
	ListStockItems() []*model.StockItem
	UpdateStockItem(s *model.StockItem) error
	DeleteStockItem(id string) error

	// 出入库流水
	CreateMovement(m *model.StockMovement) error
	GetMovement(id string) (*model.StockMovement, error)
	ListMovements() []*model.StockMovement
	DeleteMovement(id string) error
}
