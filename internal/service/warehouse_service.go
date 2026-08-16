package service

import (
	"sort"
	"time"

	"inventory/internal/model"
	"inventory/pkg/idgen"
)

// CreateWarehouse 创建仓库。
func (s *Service) CreateWarehouse(input model.Warehouse) (*model.Warehouse, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	w := &model.Warehouse{
		ID:        idgen.Hex(),
		Code:      input.Code,
		Name:      input.Name,
		Location:  input.Location,
		Manager:   input.Manager,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.store.CreateWarehouse(w); err != nil {
		return nil, err
	}
	s.log.Infof("创建仓库 %s(%s)", w.Name, w.Code)
	return w, nil
}

// GetWarehouse 按 ID 查询仓库。
func (s *Service) GetWarehouse(id string) (*model.Warehouse, error) {
	return s.store.GetWarehouse(id)
}

// ListWarehouses 列出全部仓库，按创建时间倒序。
func (s *Service) ListWarehouses() ([]*model.Warehouse, error) {
	list := s.store.ListWarehouses()
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})
	return list, nil
}

// UpdateWarehouse 更新仓库可编辑字段。
func (s *Service) UpdateWarehouse(id string, input model.Warehouse) (*model.Warehouse, error) {
	exist, err := s.store.GetWarehouse(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Location = input.Location
	exist.Manager = input.Manager
	if input.Code != "" {
		exist.Code = input.Code
	}
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateWarehouse(exist); err != nil {
		return nil, err
	}
	s.log.Infof("更新仓库 %s(%s)", exist.Name, exist.Code)
	return exist, nil
}

// DeleteWarehouse 删除仓库。
func (s *Service) DeleteWarehouse(id string) error {
	if err := s.store.DeleteWarehouse(id); err != nil {
		return err
	}
	s.log.Infof("删除仓库 %s", id)
	return nil
}
