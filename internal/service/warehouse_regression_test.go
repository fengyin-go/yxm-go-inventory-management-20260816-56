package service

import (
	"testing"
	"time"

	"inventory/internal/model"
)

func TestWarehouseRegression(t *testing.T) {
	s := newTestService()

	w1, err := s.CreateWarehouse(model.Warehouse{Code: "WH-1", Name: "一号仓"})
	if err != nil {
		t.Fatalf("创建一号仓失败: %v", err)
	}
	time.Sleep(5 * time.Millisecond)
	w2, err := s.CreateWarehouse(model.Warehouse{Code: "WH-2", Name: "二号仓"})
	if err != nil {
		t.Fatalf("创建二号仓失败: %v", err)
	}

	p, err := s.CreateProduct(model.Product{SKU: "SKU-W", Name: "库存商品"})
	if err != nil {
		t.Fatalf("创建商品失败: %v", err)
	}
	if _, err := s.CreateStockItem(model.StockItem{ProductID: p.ID, WarehouseID: w1.ID, Quantity: 5, LowThreshold: 10}); err != nil {
		t.Fatalf("创建低库存失败: %v", err)
	}
	if _, err := s.CreateStockItem(model.StockItem{ProductID: p.ID, WarehouseID: w2.ID, Quantity: 20, LowThreshold: 10}); err != nil {
		t.Fatalf("创建正常库存失败: %v", err)
	}

	t.Run("仓库列表按创建时间倒序", func(t *testing.T) {
		list, err := s.ListWarehouses()
		if err != nil {
			t.Fatalf("仓库列表失败: %v", err)
		}
		if len(list) != 2 || list[0].ID != w2.ID {
			t.Fatalf("期望二号仓排第一，得到 %+v", list)
		}
	})

	t.Run("仓库汇总名称正确", func(t *testing.T) {
		sums, err := s.WarehouseSummaries()
		if err != nil {
			t.Fatalf("仓库汇总失败: %v", err)
		}
		byID := map[string]string{}
		lowByID := map[string]int{}
		for _, sum := range sums {
			byID[sum.WarehouseID] = sum.WarehouseName
			lowByID[sum.WarehouseID] = sum.LowStockCount
		}
		if byID[w1.ID] != "一号仓" || lowByID[w1.ID] != 1 {
			t.Fatalf("一号仓汇总异常: %+v", sums)
		}
		if byID[w2.ID] != "二号仓" || lowByID[w2.ID] != 0 {
			t.Fatalf("二号仓汇总异常: %+v", sums)
		}
	})

	t.Run("重复仓库编码应拒绝", func(t *testing.T) {
		if _, err := s.CreateWarehouse(model.Warehouse{Code: "WH-1", Name: "三号仓"}); err == nil {
			t.Fatalf("期望重复编码被拒绝，实际创建成功")
		}
	})

	t.Run("非法仓库状态应拒绝", func(t *testing.T) {
		if _, err := s.CreateWarehouse(model.Warehouse{Code: "WH-3", Name: "四号仓", Status: "banana"}); !model.IsValidationError(err) {
			t.Fatalf("期望校验错误，得到 %v", err)
		}
	})
}
