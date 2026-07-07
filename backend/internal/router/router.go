package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/handler"
	"github.com/mall-admin/backend/internal/middleware"
	"github.com/mall-admin/backend/internal/response"
	"github.com/mall-admin/backend/internal/service"
	"gorm.io/gorm"
)

type Options struct {
	DB               *gorm.DB
	JWTSecret        string
	JWTExpireSeconds int64
}

func New(opts Options) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	authService := service.NewAuthService(opts.DB, opts.JWTSecret, opts.JWTExpireSeconds)
	rbacService := service.NewRBACService(opts.DB)
	adminUserService := service.NewAdminUserService(opts.DB)
	roleService := service.NewRoleService(opts.DB)
	systemService := service.NewSystemService(opts.DB)
	catalogService := service.NewCatalogService(opts.DB)
	productService := service.NewProductService(opts.DB)

	authHandler := handler.NewAuthHandler(authService)
	adminUserHandler := handler.NewAdminUserHandler(adminUserService)
	roleHandler := handler.NewRoleHandler(roleService)
	systemHandler := handler.NewSystemHandler(systemService)
	catalogHandler := handler.NewCatalogHandler(catalogService)
	productHandler := handler.NewProductHandler(productService)

	api := engine.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Body{Code: 0, Message: "成功", Data: gin.H{"status": "ok"}})
	})
	api.POST("/auth/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.Auth(authService))
	protected.GET("/auth/profile", authHandler.Profile)
	protected.GET("/auth/menus", authHandler.Menus)

	secured := protected.Group("")
	secured.Use(middleware.RBAC(rbacService))
	secured.GET("/admin-users", adminUserHandler.List)
	secured.POST("/admin-users", adminUserHandler.Create)
	secured.PUT("/admin-users/:id", adminUserHandler.Update)
	secured.PATCH("/admin-users/:id/password", adminUserHandler.ResetPassword)
	secured.DELETE("/admin-users/:id", adminUserHandler.Delete)

	secured.GET("/roles", roleHandler.List)
	secured.POST("/roles", roleHandler.Create)
	secured.PUT("/roles/:id", roleHandler.Update)
	secured.DELETE("/roles/:id", roleHandler.Delete)
	secured.GET("/roles/:id/permissions", roleHandler.PermissionDetail)
	secured.PUT("/roles/:id/permissions", roleHandler.UpdatePermissions)

	secured.GET("/menus/tree", systemHandler.MenuTree)
	secured.GET("/api-permissions", systemHandler.APIPermissions)

	secured.GET("/product-categories", catalogHandler.CategoryList)
	secured.POST("/product-categories", catalogHandler.CreateCategory)
	secured.PUT("/product-categories/:id", catalogHandler.UpdateCategory)
	secured.DELETE("/product-categories/:id", catalogHandler.DeleteCategory)

	secured.GET("/brands", catalogHandler.BrandList)
	secured.POST("/brands", catalogHandler.CreateBrand)
	secured.PUT("/brands/:id", catalogHandler.UpdateBrand)
	secured.DELETE("/brands/:id", catalogHandler.DeleteBrand)

	secured.GET("/products", productHandler.List)
	secured.GET("/products/:id", productHandler.Get)
	secured.POST("/products", productHandler.Create)
	secured.PUT("/products/:id", productHandler.Update)
	secured.DELETE("/products/:id", productHandler.Delete)
	secured.POST("/products/:id/skus", productHandler.CreateSKU)
	secured.PUT("/products/:id/skus/:skuId", productHandler.UpdateSKU)
	secured.DELETE("/products/:id/skus/:skuId", productHandler.DeleteSKU)
	secured.PATCH("/products/:id/status", productHandler.UpdateStatus)

	return engine
}
