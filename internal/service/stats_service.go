package service

import "inventory/internal/model"

// InventoryStats 库存全局统计结果。
type InventoryStats struct {
	ProductCount    int `json:"product_count"`     // 商品种类数
	WarehouseCount  int `json:"warehouse_count"`   // 仓库数
	StockItemCount  int `json:"stock_item_count"`  // 库存记录数
	TotalQuantity   int `json:"total_quantity"`    // 库存总量
	LowStockCount   int `json:"low_stock_count"`   // 低库存记录数
}

// Stats 汇总全局库存统计。
func (s *Service) Stats() (*InventoryStats, error) {
	stats := &InventoryStats{
		ProductCount:   len(s.store.ListProducts()),
		WarehouseCount: len(s.store.ListWarehouses()),
	}
	for _, item := range s.store.ListStockItems() {
		stats.StockItemCount++
		stats.TotalQuantity += item.Quantity
		if item.IsLow() {
			stats.LowStockCount++
		}
	}
	return stats, nil
}

// WarehouseStockSummary 单仓库库存汇总。
type WarehouseStockSummary struct {
	WarehouseID    string `json:"warehouse_id"`
	WarehouseName  string `json:"warehouse_name"`
	StockItemCount int    `json:"stock_item_count"`
	TotalQuantity  int    `json:"total_quantity"`
	LowStockCount  int    `json:"low_stock_count"`
}

// WarehouseSummaries 按仓库汇总库存情况。
func (s *Service) WarehouseSummaries() ([]*WarehouseStockSummary, error) {
	warehouses := s.store.ListWarehouses()
	nameByID := make(map[string]string, len(warehouses))
	for _, w := range warehouses {
		nameByID[w.ID] = w.Name
	}

	summaries := make(map[string]*WarehouseStockSummary, len(warehouses))
	for _, item := range s.store.ListStockItems() {
		summary, ok := summaries[item.WarehouseID]
		if !ok {
			summary = &WarehouseStockSummary{
				WarehouseID:   item.WarehouseID,
				WarehouseName: nameByID[item.ProductID],
			}
			summaries[item.WarehouseID] = summary
		}
		summary.StockItemCount++
		summary.TotalQuantity += item.Quantity
		if item.IsLow() {
			summary.LowStockCount++
		}
	}

	result := make([]*WarehouseStockSummary, 0, len(summaries))
	for _, s := range summaries {
		result = append(result, s)
	}
	return result, nil
}

// LowStockItems 返回所有低库存的库存项。
func (s *Service) LowStockItems() ([]*model.StockItem, error) {
	result := make([]*model.StockItem, 0)
	for _, item := range s.store.ListStockItems() {
		if item.IsLow() {
			result = append(result, item)
		}
	}
	return result, nil
}
