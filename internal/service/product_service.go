package service

import (
	"sort"
	"time"

	"inventory/internal/model"
	"inventory/pkg/idgen"
)

// CreateProduct 创建商品。
func (s *Service) CreateProduct(input model.Product) (*model.Product, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	p := &model.Product{
		ID:          idgen.Hex(),
		SKU:         input.SKU,
		Name:        input.Name,
		Category:    input.Category,
		Unit:        input.Unit,
		Price:       input.Price,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.store.CreateProduct(p); err != nil {
		return nil, err
	}
	s.log.Infof("创建商品 %s(%s)", p.Name, p.SKU)
	return p, nil
}

// GetProduct 按 ID 查询商品。
func (s *Service) GetProduct(id string) (*model.Product, error) {
	return s.store.GetProduct(id)
}

// ListProducts 列出商品，支持筛选与分页，返回当前页数据与总数。
func (s *Service) ListProducts(filter model.ProductFilter, page, size int) ([]*model.Product, int, error) {
	all := s.store.ListProducts()
	matched := make([]*model.Product, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Product{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateProduct 更新商品可编辑字段（ID 与创建时间保持不变）。
func (s *Service) UpdateProduct(id string, input model.Product) (*model.Product, error) {
	exist, err := s.store.GetProduct(id)
	if err != nil {
		return nil, err
	}
	// 在副本上应用可编辑字段并校验，避免校验失败时污染已存记录。
	merged := *exist
	merged.Name = input.Name
	merged.Category = input.Category
	merged.Unit = input.Unit
	merged.Price = input.Price
	merged.Description = input.Description
	if input.Status != "" {
		merged.Status = input.Status
	}
	if input.SKU != "" {
		merged.SKU = input.SKU
	}
	if err := merged.Validate(); err != nil {
		return nil, err
	}
	merged.UpdatedAt = time.Now()
	if err := s.store.UpdateProduct(&merged); err != nil {
		return nil, err
	}
	s.log.Infof("更新商品 %s(%s)", merged.Name, merged.SKU)
	return &merged, nil
}

// DeleteProduct 删除商品。
func (s *Service) DeleteProduct(id string) error {
	if err := s.store.DeleteProduct(id); err != nil {
		return err
	}
	s.log.Infof("删除商品 %s", id)
	return nil
}

// SetProductStatus 上架/下架商品。
func (s *Service) SetProductStatus(id, status string) (*model.Product, error) {
	exist, err := s.store.GetProduct(id)
	if err != nil {
		return nil, err
	}
	// 在副本上校验状态合法性，避免非法值污染已存记录。
	merged := *exist
	merged.Status = status
	if err := merged.Validate(); err != nil {
		return nil, err
	}
	merged.UpdatedAt = time.Now()
	if err := s.store.UpdateProduct(&merged); err != nil {
		return nil, err
	}
	s.log.Infof("商品 %s 状态变更为 %s", id, status)
	return &merged, nil
}
