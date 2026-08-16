package service

import (
	"sort"

	"inventory/internal/model"
)

// GetMovement 按 ID 查询流水。
func (s *Service) GetMovement(id string) (*model.StockMovement, error) {
	return s.store.GetMovement(id)
}

// ListMovements 列出流水，支持筛选与分页，按创建时间倒序。
func (s *Service) ListMovements(filter model.MovementFilter, page, size int) ([]*model.StockMovement, int, error) {
	all := s.store.ListMovements()
	matched := make([]*model.StockMovement, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.StockMovement{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteMovement 删除流水。
func (s *Service) DeleteMovement(id string) error {
	if err := s.store.DeleteMovement(id); err != nil {
		return err
	}
	s.log.Infof("删除流水 %s", id)
	return nil
}
