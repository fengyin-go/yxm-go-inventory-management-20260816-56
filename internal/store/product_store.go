package store

import "inventory/internal/model"

// CreateProduct 新增商品，SKU 重复时返回 ErrConflict。
func (s *MemoryStore) CreateProduct(p *model.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.products {
		if exist.SKU == p.SKU {
			return ErrConflict
		}
	}
	s.products[p.ID] = p
	return nil
}

// GetProduct 按 ID 查询商品。
func (s *MemoryStore) GetProduct(id string) (*model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.products[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// GetProductBySKU 按 SKU 查询商品。
func (s *MemoryStore) GetProductBySKU(sku string) (*model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.products {
		if p.SKU == sku {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// ListProducts 返回全部商品（无序，由服务层排序）。
func (s *MemoryStore) ListProducts() []*model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Product, 0, len(s.products))
	for _, p := range s.products {
		list = append(list, p)
	}
	return list
}

// UpdateProduct 覆盖保存商品，SKU 与其它商品重复时返回 ErrConflict。
func (s *MemoryStore) UpdateProduct(p *model.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.products {
		if exist.ID != p.ID && exist.SKU == p.SKU {
			return ErrConflict
		}
	}
	s.products[p.ID] = p
	return nil
}

// DeleteProduct 按 ID 删除商品。
func (s *MemoryStore) DeleteProduct(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; !ok {
		return ErrNotFound
	}
	delete(s.products, id)
	return nil
}
