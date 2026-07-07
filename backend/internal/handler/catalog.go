package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type CatalogHandler struct {
	catalog *service.CatalogService
}

func NewCatalogHandler(catalog *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalog: catalog}
}

func (h *CatalogHandler) CategoryList(c *gin.Context) {
	items, err := h.catalog.CategoryList(c.Query("keyword"), parseStatus(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var req service.CategoryInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.catalog.CreateCategory(req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *CatalogHandler) UpdateCategory(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.CategoryInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.catalog.UpdateCategory(id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *CatalogHandler) DeleteCategory(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.catalog.DeleteCategory(id); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *CatalogHandler) BrandList(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.catalog.BrandList(c.Query("keyword"), parseStatus(c), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *CatalogHandler) CreateBrand(c *gin.Context) {
	var req service.BrandInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.catalog.CreateBrand(req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *CatalogHandler) UpdateBrand(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.BrandInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.catalog.UpdateBrand(id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *CatalogHandler) DeleteBrand(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.catalog.DeleteBrand(id); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}
