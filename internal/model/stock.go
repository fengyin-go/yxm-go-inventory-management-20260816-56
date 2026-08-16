package model

import "time"

// StockItem 表示某商品在某仓库的库存记录。
type StockItem struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	WarehouseID  string    `json:"warehouse_id"`
	Quantity     int       `json:"quantity"`      // 当前库存数量
	LowThreshold int       `json:"low_threshold"` // 该库存项的预警阈值
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsLow 判断库存是否低于预警阈值。
func (s *StockItem) IsLow() bool {
	return s.Quantity <= s.LowThreshold
}

// Add 增加库存数量，返回增加后的数量。
func (s *StockItem) Add(qty int) int {
	s.Quantity += qty
	s.UpdatedAt = time.Now()
	return s.Quantity
}

// Subtract 扣减库存数量，返回扣减后的数量。
// 调用方需保证扣减后不为负。
func (s *StockItem) Subtract(qty int) int {
	s.Quantity -= qty
	s.UpdatedAt = time.Now()
	return s.Quantity
}

// StockItemFilter 库存列表筛选条件。
type StockItemFilter struct {
	ProductID   string
	WarehouseID string
	OnlyLow     bool // 仅返回低库存项
}

// Match 判断库存项是否命中筛选条件。
func (f StockItemFilter) Match(s *StockItem) bool {
	if f.ProductID != "" && s.ProductID != f.ProductID {
		return false
	}
	if f.WarehouseID != "" && s.WarehouseID != f.WarehouseID {
		return false
	}
	if f.OnlyLow && !s.IsLow() {
		return false
	}
	return true
}
