package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/service"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Page struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "成功", Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Body{Code: 0, Message: "成功", Data: data})
}

func NoContentOK(c *gin.Context) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "成功"})
}

func Fail(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := 500
	message := "服务器内部错误"
	if appErr, ok := service.AsAppError(err); ok {
		status = appErr.HTTPStatus
		code = appErr.Code
		message = appErr.Message
	}
	c.JSON(status, Body{Code: code, Message: message})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Body{Code: 400, Message: message})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Body{Code: 401, Message: "未登录或登录已过期"})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Body{Code: 403, Message: "无权限访问"})
}
