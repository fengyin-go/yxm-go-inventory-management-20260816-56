package handler

import (
	"net/http"

	"inventory/internal/model"
	"inventory/pkg/httpx"
)

// registerWarehouseRoutes 注册仓库相关路由。
func (s *Server) registerWarehouseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/warehouses", s.createWarehouse)
	mux.HandleFunc("GET /api/warehouses", s.listWarehouses)
	mux.HandleFunc("GET /api/warehouses/{id}", s.getWarehouse)
	mux.HandleFunc("PUT /api/warehouses/{id}", s.updateWarehouse)
	mux.HandleFunc("DELETE /api/warehouses/{id}", s.deleteWarehouse)
}

// createWarehouseRequest 创建仓库请求体。
type createWarehouseRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Manager  string `json:"manager"`
}

// createWarehouse 创建仓库：POST /api/warehouses
func (s *Server) createWarehouse(w http.ResponseWriter, r *http.Request) {
	var req createWarehouseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.CreateWarehouse(model.Warehouse{
		Code:     req.Code,
		Name:     req.Name,
		Location: req.Location,
		Manager:  req.Manager,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wh)
}

// listWarehouses 仓库列表：GET /api/warehouses
func (s *Server) listWarehouses(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ListWarehouses()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

// getWarehouse 仓库详情：GET /api/warehouses/{id}
func (s *Server) getWarehouse(w http.ResponseWriter, r *http.Request) {
	wh, err := s.svc.GetWarehouse(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

// updateWarehouseRequest 更新仓库请求体。
type updateWarehouseRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Manager  string `json:"manager"`
	Status   string `json:"status"`
}

// updateWarehouse 更新仓库：PUT /api/warehouses/{id}
func (s *Server) updateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req updateWarehouseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.UpdateWarehouse(r.PathValue("id"), model.Warehouse{
		Code:     req.Code,
		Name:     req.Name,
		Location: req.Location,
		Manager:  req.Manager,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

// deleteWarehouse 删除仓库：DELETE /api/warehouses/{id}
func (s *Server) deleteWarehouse(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteWarehouse(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}
