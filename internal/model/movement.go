package model

import "time"

// 出入库类型常量。
const (
	MovementIn  = "in"  // 入库
	MovementOut = "out" // 出库
)

// StockMovement 出入库流水，记录一次库存变动的完整快照。
type StockMovement struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	WarehouseID string    `json:"warehouse_id"`
	Type        string    `json:"type"`      // in / out
	Quantity    int       `json:"quantity"`  // 本次变动数量（正数）
	BeforeQty   int       `json:"before_qty"` // 变动前库存
	AfterQty    int       `json:"after_qty"`  // 变动后库存
	Operator    string    `json:"operator"`   // 操作人
	Remark      string    `json:"remark"`     // 备注
	CreatedAt   time.Time `json:"created_at"`
}

// MovementFilter 流水筛选条件。
type MovementFilter struct {
	ProductID   string
	WarehouseID string
	Type        string
}

// Match 判断流水是否命中筛选条件。
func (f MovementFilter) Match(m *StockMovement) bool {
	if f.ProductID != "" && m.ProductID != f.ProductID {
		return false
	}
	if f.WarehouseID != "" && m.WarehouseID != f.WarehouseID {
		return false
	}
	if f.Type != "" && m.Type != f.Type {
		return false
	}
	return true
}

// IsInbound 是否为入库流水。
func (m *StockMovement) IsInbound() bool {
	return m.Type == MovementIn
}
