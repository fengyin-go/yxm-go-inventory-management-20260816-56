package handler

import (
	"net/http"

	"inventory/internal/model"
	"inventory/pkg/httpx"
)

// registerMovementRoutes 注册出入库流水相关路由。
func (s *Server) registerMovementRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/movements", s.listMovements)
	mux.HandleFunc("GET /api/movements/{id}", s.getMovement)
	mux.HandleFunc("DELETE /api/movements/{id}", s.deleteMovement)
}

// listMovements 流水列表：GET /api/movements?product_id=&warehouse_id=&type=&page=&size=
func (s *Server) listMovements(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MovementFilter{
		ProductID:   r.URL.Query().Get("product_id"),
		WarehouseID: r.URL.Query().Get("warehouse_id"),
		Type:        r.URL.Query().Get("type"),
	}
	items, total, err := s.svc.ListMovements(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

// getMovement 流水详情：GET /api/movements/{id}
func (s *Server) getMovement(w http.ResponseWriter, r *http.Request) {
	mov, err := s.svc.GetMovement(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, mov)
}

// deleteMovement 删除流水：DELETE /api/movements/{id}
func (s *Server) deleteMovement(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteMovement(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
