package store

import "inventory/internal/model"

// CreateWarehouse 新增仓库，编码重复时返回 ErrConflict。
func (s *MemoryStore) CreateWarehouse(w *model.Warehouse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.warehouses {
		if exist.Name == w.Name {
			return ErrConflict
		}
	}
	s.warehouses[w.ID] = w
	return nil
}

// GetWarehouse 按 ID 查询仓库。
func (s *MemoryStore) GetWarehouse(id string) (*model.Warehouse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.warehouses[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

// GetWarehouseByCode 按编码查询仓库。
func (s *MemoryStore) GetWarehouseByCode(code string) (*model.Warehouse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.warehouses {
		if w.Code == code {
			return w, nil
		}
	}
	return nil, ErrNotFound
}

// ListWarehouses 返回全部仓库。
func (s *MemoryStore) ListWarehouses() []*model.Warehouse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Warehouse, 0, len(s.warehouses))
	for _, w := range s.warehouses {
		list = append(list, w)
	}
	return list
}

// UpdateWarehouse 覆盖保存仓库。
func (s *MemoryStore) UpdateWarehouse(w *model.Warehouse) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.warehouses[w.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.warehouses {
		if exist.ID != w.ID && exist.Code == w.Code {
			return ErrConflict
		}
	}
	s.warehouses[w.ID] = w
	return nil
}

// DeleteWarehouse 按 ID 删除仓库。
func (s *MemoryStore) DeleteWarehouse(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.warehouses[id]; !ok {
		return ErrNotFound
	}
	delete(s.warehouses, id)
	return nil
}
