package store

import "inventory/internal/model"

// CreateStockItem 新增库存项，同一商品+仓库重复时返回 ErrConflict。
func (s *MemoryStore) CreateStockItem(item *model.StockItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.stocks {
		if exist.ProductID == item.ProductID && exist.WarehouseID == item.WarehouseID {
			return ErrConflict
		}
	}
	s.stocks[item.ID] = item
	return nil
}

// GetStockItem 按 ID 查询库存项。
func (s *MemoryStore) GetStockItem(id string) (*model.StockItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.stocks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return item, nil
}

// GetStockItemByProductWarehouse 按商品+仓库定位库存项。
func (s *MemoryStore) GetStockItemByProductWarehouse(productID, warehouseID string) (*model.StockItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.stocks {
		if item.ProductID == productID && item.WarehouseID == warehouseID {
			return item, nil
		}
	}
	return nil, ErrNotFound
}

// ListStockItems 返回全部库存项。
func (s *MemoryStore) ListStockItems() []*model.StockItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.StockItem, 0, len(s.stocks))
	for _, item := range s.stocks {
		list = append(list, item)
	}
	return list
}

// UpdateStockItem 覆盖保存库存项。
func (s *MemoryStore) UpdateStockItem(item *model.StockItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stocks[item.ID]; !ok {
		return ErrNotFound
	}
	s.stocks[item.ID] = item
	return nil
}

// DeleteStockItem 按 ID 删除库存项。
func (s *MemoryStore) DeleteStockItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stocks[id]; !ok {
		return ErrNotFound
	}
	delete(s.stocks, id)
	return nil
}
