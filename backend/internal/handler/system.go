package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type SystemHandler struct {
	system *service.SystemService
}

func NewSystemHandler(system *service.SystemService) *SystemHandler {
	return &SystemHandler{system: system}
}

func (h *SystemHandler) MenuTree(c *gin.Context) {
	menus, err := h.system.MenuTree()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, menus)
}

func (h *SystemHandler) APIPermissions(c *gin.Context) {
	permissions, err := h.system.APIPermissions()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, permissions)
}
