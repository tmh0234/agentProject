package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

func bindJSON(c *gin.Context, dst interface{}) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		response.BadRequest(c, "参数错误")
		return false
	}
	return true
}

func parseID(c *gin.Context, name string) (uint, bool) {
	id, err := service.ParseUintParam(c.Param(name))
	if err != nil {
		response.Fail(c, err)
		return 0, false
	}
	return id, true
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func parseStatus(c *gin.Context) *int {
	raw := c.Query("status")
	if raw == "" {
		return nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &value
}
