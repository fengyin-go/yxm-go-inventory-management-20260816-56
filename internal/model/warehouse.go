package model

import (
	"strings"
	"time"
)

// 仓库状态常量。
const (
	WarehouseActive   = "active"   // 启用
	WarehouseInactive = "inactive" // 停用
)

// Warehouse 仓库实体。
type Warehouse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`     // 仓库编码
	Name      string    `json:"name"`     // 仓库名称
	Location  string    `json:"location"` // 仓库地址
	Manager   string    `json:"manager"`  // 负责人
	Status    string    `json:"status"`   // active / inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 规范化并校验仓库字段。
func (w *Warehouse) Validate() error {
	w.Code = strings.TrimSpace(w.Code)
	w.Name = strings.TrimSpace(w.Name)
	w.Location = strings.TrimSpace(w.Location)
	w.Manager = strings.TrimSpace(w.Manager)
	if w.Code == "" {
		return NewValidationError("code", "仓库编码不能为空")
	}
	if w.Name == "" {
		return NewValidationError("name", "仓库名称不能为空")
	}
	if w.Status == "" {
		w.Status = WarehouseActive
	}
	if w.Status != WarehouseActive && w.Status != WarehouseInactive {
		return NewValidationError("status", "仓库状态不合法，应为 active 或 inactive")
	}
	return nil
}
