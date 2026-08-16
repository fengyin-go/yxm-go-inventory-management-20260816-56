package model

import (
	"strings"
	"time"
)

// 商品状态常量。
const (
	ProductActive   = "active"   // 上架
	ProductInactive = "inactive" // 下架
)

// Product 商品实体。
type Product struct {
	ID          string    `json:"id"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Unit        string    `json:"unit"`        // 计量单位，如 件/个/箱
	Price       float64   `json:"price"`       // 参考单价（元）
	Description string    `json:"description"` // 商品描述
	Status      string    `json:"status"`      // active / inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductFilter 商品列表筛选条件。
type ProductFilter struct {
	Category string // 按分类精确匹配
	Status   string // 按状态匹配
	Keyword  string // 匹配名称或 SKU（忽略大小写）
}

// Match 判断商品是否命中筛选条件。
func (f ProductFilter) Match(p *Product) bool {
	if f.Category != "" && p.Category != f.Category {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		// 关键词匹配名称或 SKU，统一转为小写以忽略大小写。
		kw := strings.ToLower(f.Keyword)
		if !strings.Contains(strings.ToLower(p.Name), kw) &&
			!strings.Contains(strings.ToLower(p.SKU), kw) {
			return false
		}
	}
	return true
}

// Validate 规范化并校验商品字段。
func (p *Product) Validate() error {
	p.SKU = strings.TrimSpace(p.SKU)
	p.Name = strings.TrimSpace(p.Name)
	p.Category = strings.TrimSpace(p.Category)
	p.Unit = strings.TrimSpace(p.Unit)
	if p.SKU == "" {
		return NewValidationError("sku", "SKU 不能为空")
	}
	if p.Name == "" {
		return NewValidationError("name", "商品名称不能为空")
	}
	if p.Price < 0 {
		return NewValidationError("price", "价格不能为负数")
	}
	if p.Unit == "" {
		p.Unit = "件"
	}
	if p.Status == "" {
		p.Status = ProductActive
	}
	if p.Status != ProductActive && p.Status != ProductInactive {
		return NewValidationError("status", "商品状态不合法，应为 active 或 inactive")
	}
	return nil
}

// Disable 下架商品。
func (p *Product) Disable() {
	p.Status = ProductInactive
	p.UpdatedAt = time.Now()
}

// Enable 上架商品。
func (p *Product) Enable() {
	p.Status = ProductActive
	p.UpdatedAt = time.Now()
}
