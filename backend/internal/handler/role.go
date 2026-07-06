package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type RoleHandler struct {
	roles *service.RoleService
}

func NewRoleHandler(roles *service.RoleService) *RoleHandler {
	return &RoleHandler{roles: roles}
}

func (h *RoleHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	items, total, err := h.roles.List(c.Query("keyword"), parseStatus(c), page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: page, PageSize: pageSize})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req service.RoleInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.roles.Create(req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.RoleInput
	if !bindJSON(c, &req) {
		return
	}
	item, err := h.roles.Update(id, req)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.roles.Delete(id); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}

func (h *RoleHandler) PermissionDetail(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	detail, err := h.roles.PermissionDetail(id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, detail)
}

func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req service.UpdateRolePermissionsInput
	if !bindJSON(c, &req) {
		return
	}
	if err := h.roles.UpdatePermissions(id, req); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContentOK(c)
}
