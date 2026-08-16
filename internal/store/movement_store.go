package store

import "inventory/internal/model"

// CreateMovement 新增出入库流水。
func (s *MemoryStore) CreateMovement(m *model.StockMovement) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.movements[m.ID] = m
	return nil
}

// GetMovement 按 ID 查询流水。
func (s *MemoryStore) GetMovement(id string) (*model.StockMovement, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.movements[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// ListMovements 返回全部流水。
func (s *MemoryStore) ListMovements() []*model.StockMovement {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.StockMovement, 0, len(s.movements))
	for _, m := range s.movements {
		list = append(list, m)
	}
	return list
}

// DeleteMovement 按 ID 删除流水。
func (s *MemoryStore) DeleteMovement(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.movements[id]; !ok {
		return ErrNotFound
	}
	delete(s.movements, id)
	return nil
}
