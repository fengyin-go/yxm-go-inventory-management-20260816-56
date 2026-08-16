package store

import (
	"sync"

	"inventory/internal/model"
)

// MemoryStore 基于内存的 Store 实现，使用读写锁保证并发安全。
type MemoryStore struct {
	mu         sync.RWMutex
	products   map[string]*model.Product
	warehouses map[string]*model.Warehouse
	stocks     map[string]*model.StockItem
	movements  map[string]*model.StockMovement
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		products:   make(map[string]*model.Product),
		warehouses: make(map[string]*model.Warehouse),
		stocks:     make(map[string]*model.StockItem),
		movements:  make(map[string]*model.StockMovement),
	}
}

// compile-time 校验 MemoryStore 实现了 Store 接口。
var _ Store = (*MemoryStore)(nil)
