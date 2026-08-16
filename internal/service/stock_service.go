package service

import (
	"fmt"
	"sort"
	"time"

	"inventory/internal/model"
	"inventory/pkg/idgen"
)

// CreateStockItem 为商品在指定仓库建立库存记录。
func (s *Service) CreateStockItem(input model.StockItem) (*model.StockItem, error) {
	if input.ProductID == "" {
		return nil, model.NewValidationError("product_id", "商品不能为空")
	}
	if input.WarehouseID == "" {
		return nil, model.NewValidationError("warehouse_id", "仓库不能为空")
	}
	if input.Quantity < 0 {
		return nil, model.NewValidationError("quantity", "库存数量不能为负数")
	}
	// 校验商品与仓库必须存在。
	if _, err := s.store.GetProduct(input.ProductID); err != nil {
		return nil, err
	}
	if _, err := s.store.GetWarehouse(input.WarehouseID); err != nil {
		return nil, err
	}

	item := &model.StockItem{
		ID:           idgen.Hex(),
		ProductID:    input.ProductID,
		WarehouseID:  input.WarehouseID,
		Quantity:     input.Quantity,
		LowThreshold: input.LowThreshold,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if item.LowThreshold <= 0 {
		item.LowThreshold = s.cfg.LowStockThreshold
	}
	if err := s.store.CreateStockItem(item); err != nil {
		return nil, err
	}
	s.log.Infof("建立库存项 product=%s warehouse=%s qty=%d", item.ProductID, item.WarehouseID, item.Quantity)
	return item, nil
}

// GetStockItem 按 ID 查询库存项。
func (s *Service) GetStockItem(id string) (*model.StockItem, error) {
	return s.store.GetStockItem(id)
}

// ListStockItems 列出库存项，支持筛选与分页。
func (s *Service) ListStockItems(filter model.StockItemFilter, page, size int) ([]*model.StockItem, int, error) {
	all := s.store.ListStockItems()
	matched := make([]*model.StockItem, 0, len(all))
	for _, item := range all {
		if filter.Match(item) {
			matched = append(matched, item)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.StockItem{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateStockThreshold 调整库存项的预警阈值。
func (s *Service) UpdateStockThreshold(id string, threshold int) (*model.StockItem, error) {
	if threshold <= 0 {
		return nil, model.NewValidationError("low_threshold", "阈值必须为正数")
	}
	item, err := s.store.GetStockItem(id)
	if err != nil {
		return nil, err
	}
	item.LowThreshold = threshold
	item.UpdatedAt = time.Now()
	if err := s.store.UpdateStockItem(item); err != nil {
		return nil, err
	}
	s.log.Infof("库存项 %s 阈值调整为 %d", id, threshold)
	return item, nil
}

// DeleteStockItem 删除库存项。
func (s *Service) DeleteStockItem(id string) error {
	if err := s.store.DeleteStockItem(id); err != nil {
		return err
	}
	s.log.Infof("删除库存项 %s", id)
	return nil
}

// StockIn 入库：累加库存并记录流水。
func (s *Service) StockIn(productID, warehouseID string, qty int, operator, remark string) (*model.StockMovement, error) {
	if qty <= 0 {
		return nil, model.NewValidationError("quantity", "入库数量必须为正数")
	}
	item, err := s.store.GetStockItemByProductWarehouse(productID, warehouseID)
	if err != nil {
		return nil, err
	}
	before := item.Quantity
	item.Add(qty)
	if err := s.store.UpdateStockItem(item); err != nil {
		return nil, err
	}
	mov := s.newMovement(productID, warehouseID, model.MovementIn, qty, before, item.Quantity, operator, remark)
	if err := s.store.CreateMovement(mov); err != nil {
		return nil, err
	}
	s.log.Infof("入库 product=%s warehouse=%s qty=%d -> %d", productID, warehouseID, qty, item.Quantity)
	return mov, nil
}

// StockOut 出库：校验库存充足后扣减并记录流水。
func (s *Service) StockOut(productID, warehouseID string, qty int, operator, remark string) (*model.StockMovement, error) {
	if qty <= 0 {
		return nil, model.NewValidationError("quantity", "出库数量必须为正数")
	}
	item, err := s.store.GetStockItemByProductWarehouse(productID, warehouseID)
	if err != nil {
		return nil, err
	}
	if item.Quantity < qty {
		return nil, model.NewValidationError("quantity", fmt.Sprintf("库存不足，当前仅剩 %d 件", item.Quantity))
	}
	before := item.Quantity
	item.Subtract(qty)
	if err := s.store.UpdateStockItem(item); err != nil {
		return nil, err
	}
	mov := s.newMovement(productID, warehouseID, model.MovementOut, qty, before, before, operator, remark)
	if err := s.store.CreateMovement(mov); err != nil {
		return nil, err
	}
	s.log.Infof("出库 product=%s warehouse=%s qty=%d -> %d", productID, warehouseID, qty, item.Quantity)
	return mov, nil
}

// newMovement 构造一条出入库流水。
func (s *Service) newMovement(productID, warehouseID, typ string, qty, before, after int, operator, remark string) *model.StockMovement {
	return &model.StockMovement{
		ID:          idgen.Hex(),
		ProductID:   productID,
		WarehouseID: warehouseID,
		Type:        typ,
		Quantity:    qty,
		BeforeQty:   before,
		AfterQty:    after,
		Operator:    operator,
		Remark:      remark,
		CreatedAt:   time.Now(),
	}
}
