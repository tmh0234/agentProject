package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
)

func Auth(auth *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c)
			c.Abort()
			return
		}
		claims, err := auth.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Fail(c, err)
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func RBAC(rbac *service.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		allowed, err := rbac.CanAccess(c.GetUint("userID"), c.Request.Method, path)
		if err != nil {
			response.Fail(c, err)
			c.Abort()
			return
		}
		if !allowed {
			response.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
