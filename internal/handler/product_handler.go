package handler

import (
	"net/http"

	"inventory/internal/model"
	"inventory/pkg/httpx"
)

// registerProductRoutes 注册商品相关路由。
func (s *Server) registerProductRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/products", s.createProduct)
	mux.HandleFunc("GET /api/products", s.listProducts)
	mux.HandleFunc("GET /api/products/{id}", s.getProduct)
	mux.HandleFunc("PUT /api/products/{id}", s.updateProduct)
	mux.HandleFunc("DELETE /api/products/{id}", s.deleteProduct)
	mux.HandleFunc("PATCH /api/products/{id}/status", s.setProductStatus)
}

// createProductRequest 创建商品请求体。
type createProductRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// createProduct 创建商品：POST /api/products
func (s *Server) createProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateProduct(model.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Category:    req.Category,
		Unit:        req.Unit,
		Price:       req.Price,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

// listProducts 商品列表：GET /api/products?category=&status=&keyword=&page=&size=
func (s *Server) listProducts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ProductFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListProducts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

// getProduct 商品详情：GET /api/products/{id}
func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetProduct(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

// updateProductRequest 更新商品请求体。
type updateProductRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
}

// updateProduct 更新商品：PUT /api/products/{id}
func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request) {
	var req updateProductRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateProduct(r.PathValue("id"), model.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Category:    req.Category,
		Unit:        req.Unit,
		Price:       req.Price,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

// deleteProduct 删除商品：DELETE /api/products/{id}
func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteProduct(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": "删除成功"})
}

// setProductStatusRequest 状态变更请求体。
type setProductStatusRequest struct {
	Status string `json:"status"`
}

// setProductStatus 上架/下架：PATCH /api/products/{id}/status
func (s *Server) setProductStatus(w http.ResponseWriter, r *http.Request) {
	var req setProductStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.SetProductStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}
