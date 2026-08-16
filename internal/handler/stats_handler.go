package handler

import (
	"net/http"

	"inventory/pkg/httpx"
)

// registerStatsRoutes 注册统计相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats", s.stats)
	mux.HandleFunc("GET /api/stats/warehouses", s.warehouseSummaries)
}

// stats 全局统计：GET /api/stats
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.Stats()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

// warehouseSummaries 按仓库汇总：GET /api/stats/warehouses
func (s *Server) warehouseSummaries(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.WarehouseSummaries()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
