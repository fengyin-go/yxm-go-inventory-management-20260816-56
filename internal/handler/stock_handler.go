package handler

import (
	"net/http"

	"inventory/internal/model"
	"inventory/pkg/httpx"
)

// registerStockRoutes 注册库存相关路由。
func (s *Server) registerStockRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/stocks", s.createStockItem)
	mux.HandleFunc("GET /api/stocks", s.listStockItems)
	mux.HandleFunc("GET /api/stocks/{id}", s.getStockItem)
	mux.HandleFunc("PATCH /api/stocks/{id}/threshold", s.updateStockThreshold)
	mux.HandleFunc("DELETE /api/stocks/{id}", s.deleteStockItem)

	mux.HandleFunc("POST /api/stocks/in", s.stockIn)
	mux.HandleFunc("POST /api/stocks/out", s.stockOut)
	mux.HandleFunc("GET /api/stocks/low", s.listLowStock)
}

// createStockItemRequest 创建库存项请求体。
type createStockItemRequest struct {
	ProductID    string `json:"product_id"`
	WarehouseID  string `json:"warehouse_id"`
	Quantity     int    `json:"quantity"`
	LowThreshold int    `json:"low_threshold"`
}

// createStockItem 建立库存记录：POST /api/stocks
func (s *Server) createStockItem(w http.ResponseWriter, r *http.Request) {
	var req createStockItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.CreateStockItem(model.StockItem{
		ProductID:    req.ProductID,
		WarehouseID:  req.WarehouseID,
		Quantity:     req.Quantity,
		LowThreshold: req.LowThreshold,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, item)
}

// listStockItems 库存列表：GET /api/stocks?product_id=&warehouse_id=&page=&size=
func (s *Server) listStockItems(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StockItemFilter{
		ProductID:   r.URL.Query().Get("product_id"),
		WarehouseID: r.URL.Query().Get("warehouse_id"),
	}
	items, total, err := s.svc.ListStockItems(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

// getStockItem 库存详情：GET /api/stocks/{id}
func (s *Server) getStockItem(w http.ResponseWriter, r *http.Request) {
	item, err := s.svc.GetStockItem(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

// updateStockThresholdRequest 阈值调整请求体。
type updateStockThresholdRequest struct {
	LowThreshold int `json:"low_threshold"`
}

// updateStockThreshold 调整预警阈值：PATCH /api/stocks/{id}/threshold
func (s *Server) updateStockThreshold(w http.ResponseWriter, r *http.Request) {
	var req updateStockThresholdRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	item, err := s.svc.UpdateStockThreshold(r.PathValue("id"), req.LowThreshold)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, item)
}

// deleteStockItem 删除库存项：DELETE /api/stocks/{id}
func (s *Server) deleteStockItem(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteStockItem(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

// stockMoveRequest 出入库请求体。
type stockMoveRequest struct {
	ProductID   string `json:"product_id"`
	WarehouseID string `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
	Operator    string `json:"operator"`
	Remark      string `json:"remark"`
}

// stockIn 入库：POST /api/stocks/in
func (s *Server) stockIn(w http.ResponseWriter, r *http.Request) {
	var req stockMoveRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mov, err := s.svc.StockIn(req.ProductID, req.WarehouseID, req.Quantity, req.Operator, req.Remark)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mov)
}

// stockOut 出库：POST /api/stocks/out
func (s *Server) stockOut(w http.ResponseWriter, r *http.Request) {
	var req stockMoveRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	mov, err := s.svc.StockOut(req.ProductID, req.WarehouseID, req.Quantity, req.Operator, req.Remark)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mov)
}

// listLowStock 低库存列表：GET /api/stocks/low
func (s *Server) listLowStock(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.LowStockItems()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
