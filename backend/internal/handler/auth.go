package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	profile, err := h.auth.Profile(c.GetUint("userID"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}

func (h *AuthHandler) Menus(c *gin.Context) {
	menus, err := h.auth.Menus(c.GetUint("userID"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, menus)
}
