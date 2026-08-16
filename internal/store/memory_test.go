package store

import (
	"testing"
	"time"

	"inventory/internal/model"
)

// 辅助：构造一个测试商品。
func testProduct(sku string) *model.Product {
	return &model.Product{
		ID:        sku + "-id",
		SKU:       sku,
		Name:      "商品-" + sku,
		Category:  "测试",
		Unit:      "件",
		Price:     10,
		Status:    model.ProductActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestMemoryStore_ProductCRUD(t *testing.T) {
	s := NewMemoryStore()
	p := testProduct("SKU-001")

	// 创建
	if err := s.CreateProduct(p); err != nil {
		t.Fatalf("创建商品失败: %v", err)
	}

	// 重复 SKU 应冲突
	if err := s.CreateProduct(testProduct("SKU-001")); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}

	// 查询
	got, err := s.GetProduct(p.ID)
	if err != nil {
		t.Fatalf("查询商品失败: %v", err)
	}
	if got.Name != p.Name {
		t.Fatalf("期望名称 %q，得到 %q", p.Name, got.Name)
	}

	// 按 SKU 查询
	bySKU, err := s.GetProductBySKU("SKU-001")
	if err != nil || bySKU.ID != p.ID {
		t.Fatalf("按 SKU 查询失败: %v, %v", bySKU, err)
	}

	// 列表
	if got := len(s.ListProducts()); got != 1 {
		t.Fatalf("期望 1 条商品，得到 %d", got)
	}

	// 更新
	p.Name = "新名称"
	if err := s.UpdateProduct(p); err != nil {
		t.Fatalf("更新失败: %v", err)
	}

	// 删除
	if err := s.DeleteProduct(p.ID); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
	if _, err := s.GetProduct(p.ID); err != ErrNotFound {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}

func TestMemoryStore_WarehouseCRUD(t *testing.T) {
	s := NewMemoryStore()
	w := &model.Warehouse{
		ID:        "w1",
		Code:      "WH-001",
		Name:      "一号仓",
		Status:    model.WarehouseActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.CreateWarehouse(w); err != nil {
		t.Fatalf("创建仓库失败: %v", err)
	}
	if err := s.CreateWarehouse(&model.Warehouse{ID: "w2", Code: "WH-001"}); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}

	got, err := s.GetWarehouse("w1")
	if err != nil || got.Name != "一号仓" {
		t.Fatalf("查询仓库失败: %v, %v", got, err)
	}

	byCode, err := s.GetWarehouseByCode("WH-001")
	if err != nil || byCode.ID != "w1" {
		t.Fatalf("按编码查询失败: %v, %v", byCode, err)
	}

	if err := s.DeleteWarehouse("w1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
	if _, err := s.GetWarehouse("w1"); err != ErrNotFound {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}

func TestMemoryStore_StockItem(t *testing.T) {
	s := NewMemoryStore()
	item := &model.StockItem{
		ID:           "s1",
		ProductID:    "p1",
		WarehouseID:  "w1",
		Quantity:     5,
		LowThreshold: 10,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.CreateStockItem(item); err != nil {
		t.Fatalf("创建库存失败: %v", err)
	}
	// 同商品+仓库重复
	dup := &model.StockItem{ID: "s2", ProductID: "p1", WarehouseID: "w1"}
	if err := s.CreateStockItem(dup); err != ErrConflict {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}

	got, err := s.GetStockItemByProductWarehouse("p1", "w1")
	if err != nil || got.ID != "s1" {
		t.Fatalf("按商品仓库查询失败: %v, %v", got, err)
	}

	got.Add(10)
	if err := s.UpdateStockItem(got); err != nil {
		t.Fatalf("更新失败: %v", err)
	}
	after, _ := s.GetStockItem("s1")
	if after.Quantity != 15 {
		t.Fatalf("期望数量 15，得到 %d", after.Quantity)
	}
}

func TestMemoryStore_Movement(t *testing.T) {
	s := NewMemoryStore()
	m := &model.StockMovement{
		ID:          "m1",
		ProductID:   "p1",
		WarehouseID: "w1",
		Type:        model.MovementIn,
		Quantity:    10,
		BeforeQty:   0,
		AfterQty:    10,
		CreatedAt:   time.Now(),
	}
	if err := s.CreateMovement(m); err != nil {
		t.Fatalf("创建流水失败: %v", err)
	}
	got, err := s.GetMovement("m1")
	if err != nil || got.Quantity != 10 {
		t.Fatalf("查询流水失败: %v, %v", got, err)
	}
	if len(s.ListMovements()) != 1 {
		t.Fatalf("期望 1 条流水")
	}
	if err := s.DeleteMovement("m1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}
