package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type AdminUserHandler struct {
	users *service.AdminUserService
}

func NewAdminUserHandler(users *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{users: users}
}

func (h *AdminUserHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.users.List(c.Query("keyword"), parseStatus(c), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req service.AdminUserInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.users.Create(req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.AdminUserInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.users.Update(id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.ResetPasswordInput
	if !bindJSON(c, &req) {
		return
	}
	if err := h.users.ResetPassword(id, req.Password); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.users.Delete(id); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}
