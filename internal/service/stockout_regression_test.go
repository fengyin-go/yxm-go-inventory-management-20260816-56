package service

import (
	"testing"

	"inventory/internal/model"
)

func TestStockOutRegression(t *testing.T) {
	s := newTestService()
	pid, wid, _ := s.seedInventory(t)

	p2, err := s.CreateProduct(model.Product{SKU: "SKU-EDGE", Name: "边界商品"})
	if err != nil {
		t.Fatalf("创建边界商品失败: %v", err)
	}
	w2, err := s.CreateWarehouse(model.Warehouse{Code: "WH-EDGE", Name: "边界仓"})
	if err != nil {
		t.Fatalf("创建边界仓失败: %v", err)
	}
	if _, err := s.CreateStockItem(model.StockItem{
		ProductID:    p2.ID,
		WarehouseID:  w2.ID,
		Quantity:     10,
		LowThreshold: 10,
	}); err != nil {
		t.Fatalf("创建边界库存失败: %v", err)
	}

	out, err := s.StockOut(pid, wid, 30, "op", "销售出库")
	if err != nil {
		t.Fatalf("出库失败: %v", err)
	}
	if out.BeforeQty != 100 || out.AfterQty != 70 {
		t.Fatalf("出库快照异常: before=%d after=%d", out.BeforeQty, out.AfterQty)
	}
	if _, err := s.StockIn(pid, wid, 20, "op", "补货入库"); err != nil {
		t.Fatalf("入库失败: %v", err)
	}

	movs, _, err := s.ListMovements(model.MovementFilter{}, 1, 20)
	if err != nil {
		t.Fatalf("查询流水失败: %v", err)
	}
	if len(movs) != 2 {
		t.Fatalf("期望 2 条流水，得到 %d", len(movs))
	}
	if movs[0].Type != model.MovementIn {
		t.Fatalf("最新流水应为入库，得到 %s", movs[0].Type)
	}

	stats, err := s.Stats()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if stats.TotalQuantity != 100 {
		t.Fatalf("期望库存总量 100，得到 %d", stats.TotalQuantity)
	}
	if stats.LowStockCount != 0 {
		t.Fatalf("期望低库存 0，得到 %d", stats.LowStockCount)
	}

	low, err := s.LowStockItems()
	if err != nil {
		t.Fatalf("查询低库存失败: %v", err)
	}
	if len(low) != 0 {
		t.Fatalf("期望低库存列表为空，得到 %d", len(low))
	}
}
