package service

import (
	"errors"
	"testing"

	"inventory/internal/model"
	"inventory/internal/store"
)

func TestProductRegression(t *testing.T) {
	s := newTestService()

	a, err := s.CreateProduct(model.Product{SKU: "SKU-A", Name: "Apple", Price: 10})
	if err != nil {
		t.Fatalf("创建商品 A 失败: %v", err)
	}
	b, err := s.CreateProduct(model.Product{SKU: "SKU-B", Name: "Banana", Price: 20})
	if err != nil {
		t.Fatalf("创建商品 B 失败: %v", err)
	}

	t.Run("关键词大小写不敏感搜索", func(t *testing.T) {
		items, total, err := s.ListProducts(model.ProductFilter{Keyword: "apple"}, 1, 10)
		if err != nil {
			t.Fatalf("搜索失败: %v", err)
		}
		if total != 1 || len(items) != 1 || items[0].ID != a.ID {
			t.Fatalf("期望搜到商品 A，得到 total=%d items=%v", total, items)
		}
	})

	t.Run("重复 SKU 更新应拒绝", func(t *testing.T) {
		if _, err := s.UpdateProduct(b.ID, model.Product{SKU: "SKU-A", Name: "Banana", Price: 20}); !errors.Is(err, store.ErrConflict) {
			t.Fatalf("期望 ErrConflict，得到 %v", err)
		}
	})

	t.Run("空名称更新应拒绝", func(t *testing.T) {
		if _, err := s.UpdateProduct(a.ID, model.Product{SKU: "SKU-A", Name: "", Price: 10}); !model.IsValidationError(err) {
			t.Fatalf("期望校验错误，得到 %v", err)
		}
	})

	t.Run("非法状态应拒绝", func(t *testing.T) {
		if _, err := s.SetProductStatus(a.ID, "banana"); !model.IsValidationError(err) {
			t.Fatalf("期望校验错误，得到 %v", err)
		}
	})

	t.Run("不存在的商品不能建库存", func(t *testing.T) {
		if _, err := s.CreateStockItem(model.StockItem{ProductID: "missing", WarehouseID: "missing", Quantity: 1}); !errors.Is(err, store.ErrNotFound) {
			t.Fatalf("期望 ErrNotFound，得到 %v", err)
		}
	})
}
