package service

import (
	"errors"
	"testing"

	"inventory/internal/config"
	"inventory/internal/model"
	"inventory/internal/store"
	"inventory/pkg/logger"
)

// newTestService 构造用于测试的 Service 实例。
func newTestService() *Service {
	cfg := &config.Config{LowStockThreshold: 10, MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

// 辅助：准备一个商品 + 仓库 + 库存项。
func (s *Service) seedInventory(t *testing.T) (productID, warehouseID, stockID string) {
	t.Helper()
	p, err := s.CreateProduct(model.Product{SKU: "SKU-1", Name: "商品", Category: "测试", Price: 10})
	if err != nil {
		t.Fatalf("创建商品失败: %v", err)
	}
	w, err := s.CreateWarehouse(model.Warehouse{Code: "WH-1", Name: "仓库"})
	if err != nil {
		t.Fatalf("创建仓库失败: %v", err)
	}
	item, err := s.CreateStockItem(model.StockItem{
		ProductID:   p.ID,
		WarehouseID: w.ID,
		Quantity:    100,
	})
	if err != nil {
		t.Fatalf("创建库存失败: %v", err)
	}
	return p.ID, w.ID, item.ID
}

func TestService_CreateProduct(t *testing.T) {
	s := newTestService()

	p, err := s.CreateProduct(model.Product{SKU: "SKU-1", Name: "商品A", Price: 20})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if p.ID == "" || p.Status != model.ProductActive {
		t.Fatalf("商品创建结果异常: %+v", p)
	}

	// 空 SKU 应校验失败
	if _, err := s.CreateProduct(model.Product{Name: "无SKU"}); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}

	// 重复 SKU 应冲突
	if _, err := s.CreateProduct(model.Product{SKU: "SKU-1", Name: "商品B"}); !errors.Is(err, store.ErrConflict) {
		t.Fatalf("期望 ErrConflict，得到 %v", err)
	}
}

func TestService_ListProductsPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 25; i++ {
		sku := "SKU-" + string(rune('A'+i))
		if _, err := s.CreateProduct(model.Product{SKU: sku, Name: "商品", Category: "测试"}); err != nil {
			t.Fatalf("创建失败: %v", err)
		}
	}

	items, total, err := s.ListProducts(model.ProductFilter{}, 2, 10)
	if err != nil {
		t.Fatalf("列表失败: %v", err)
	}
	if total != 25 {
		t.Fatalf("期望总数 25，得到 %d", total)
	}
	if len(items) != 10 {
		t.Fatalf("期望第二页 10 条，得到 %d", len(items))
	}

	// 关键词筛选
	items, total, _ = s.ListProducts(model.ProductFilter{Keyword: "SKU-A"}, 1, 10)
	if total != 1 {
		t.Fatalf("期望筛选后 1 条，得到 %d", total)
	}
}

func TestService_StockInOut(t *testing.T) {
	s := newTestService()
	pid, wid, _ := s.seedInventory(t)

	// 入库 50
	mov, err := s.StockIn(pid, wid, 50, "op", "采购入库")
	if err != nil {
		t.Fatalf("入库失败: %v", err)
	}
	if mov.BeforeQty != 100 || mov.AfterQty != 150 {
		t.Fatalf("入库快照异常: before=%d after=%d", mov.BeforeQty, mov.AfterQty)
	}

	// 出库 30
	mov, err = s.StockOut(pid, wid, 30, "op", "销售出库")
	if err != nil {
		t.Fatalf("出库失败: %v", err)
	}
	if mov.AfterQty != 120 {
		t.Fatalf("期望出库后 120，得到 %d", mov.AfterQty)
	}

	// 超量出库应失败
	if _, err := s.StockOut(pid, wid, 9999, "op", ""); !model.IsValidationError(err) {
		t.Fatalf("期望校验错误，得到 %v", err)
	}
}

func TestService_LowStock(t *testing.T) {
	s := newTestService()
	pid, wid, _ := s.seedInventory(t)

	// 初始 100 > 阈值 10，不应低库存
	low, _ := s.LowStockItems()
	if len(low) != 0 {
		t.Fatalf("期望无低库存，得到 %d", len(low))
	}

	// 出库到 5 件，应触发低库存
	if _, err := s.StockOut(pid, wid, 95, "op", ""); err != nil {
		t.Fatalf("出库失败: %v", err)
	}
	low, _ = s.LowStockItems()
	if len(low) != 1 {
		t.Fatalf("期望 1 条低库存，得到 %d", len(low))
	}

	stats, _ := s.Stats()
	if stats.LowStockCount != 1 || stats.TotalQuantity != 5 {
		t.Fatalf("统计异常: %+v", stats)
	}
}

func TestService_StockItemNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.StockIn("nope", "nope", 10, "", ""); !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}
