package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type ProductHandler struct {
	products *service.ProductService
}

func NewProductHandler(products *service.ProductService) *ProductHandler {
	return &ProductHandler{products: products}
}

func (h *ProductHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.products.List(c.Query("keyword"), parseStatus(c), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := h.products.Get(id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req service.ProductInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.products.Create(req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.ProductInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.products.Update(id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.products.Delete(id); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *ProductHandler) CreateSKU(c *gin.Context) {
	productID, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.SKUInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.products.CreateSKU(productID, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *ProductHandler) UpdateSKU(c *gin.Context) {
	productID, ok := parseID(c, "id")
	if !ok {
		return
	}
	skuID, ok := parseID(c, "skuId")
	if !ok {
		return
	}
	var req service.SKUInput
	if !bindJSON(c, &req) {
		return
	}
	if err := h.products.UpdateSKU(productID, skuID, req); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *ProductHandler) DeleteSKU(c *gin.Context) {
	productID, ok := parseID(c, "id")
	if !ok {
		return
	}
	skuID, ok := parseID(c, "skuId")
	if !ok {
		return
	}
	if err := h.products.DeleteSKU(productID, skuID); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.ProductStatusInput
	if !bindJSON(c, &req) {
		return
	}
	if err := h.products.UpdateStatus(id, req.Status); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}
