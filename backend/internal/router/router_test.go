package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/router"
	"github.com/mall-admin/backend/internal/service"
	"github.com/mall-admin/backend/internal/testutil"
	"gorm.io/gorm"
)

func TestLoginResponseIncludesProfile(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.SeedAdmin(t, db)
	engine := newTestRouter(db)

	w := performJSON(engine, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"username": "admin",
		"password": "Admin@123456",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Code int `json:"code"`
		Data struct {
			Token   string `json:"token"`
			Profile struct {
				Username string `json:"username"`
			} `json:"profile"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Token == "" || body.Data.Profile.Username != "admin" {
		t.Fatalf("unexpected login response: %#v", body.Data)
	}
}

func TestDisabledAdminOldJWTCannotAccessSecuredAPI(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "商品管理员", "product_admin")
	permission := testutil.CreateAPIPermission(t, db, "GET", "/api/v1/products", "product:list")
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachPermissionToRole(t, db, role.ID, permission.ID)
	auth := service.NewAuthService(db, "test-secret", 3600)
	token, err := auth.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	if err := db.Model(&admin).Update("status", 0).Error; err != nil {
		t.Fatalf("disable admin: %v", err)
	}
	engine := newTestRouter(db)

	w := performJSON(engine, http.MethodGet, "/api/v1/products", token, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCategoryListEndpointAppliesKeywordAndStatusFilters(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin, token := seedAdminWithPermission(t, db, "GET", "/api/v1/product-categories", "category:list")
	_ = admin
	if err := db.Create(&[]model.ProductCategory{
		{Name: "手机数码", Status: 1, Sort: 1},
		{Name: "手机配件", Status: 0, Sort: 2},
		{Name: "家用电器", Status: 1, Sort: 3},
	}).Error; err != nil {
		t.Fatalf("seed categories: %v", err)
	}
	engine := newTestRouter(db)

	w := performJSON(engine, http.MethodGet, "/api/v1/product-categories?keyword=手机&status=1", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data []model.ProductCategory `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Name != "手机数码" {
		t.Fatalf("unexpected categories: %#v", body.Data)
	}
}

func TestDeleteLastValidSKUEndpointReturnsBusinessError(t *testing.T) {
	db := testutil.NewTestDB(t)
	_, token := seedAdminWithPermission(t, db, "DELETE", "/api/v1/products/:id/skus/:skuId", "product-sku:delete")
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")
	products := service.NewProductService(db)
	product, err := products.Create(service.ProductInput{
		Name:       "已上架商品",
		Code:       "SPU-API-001",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Status:     1,
		Skus: []service.SKUInput{
			{Code: "SKU-API-001", PriceCents: 1000, Stock: 8, Status: 1},
		},
	})
	if err != nil {
		t.Fatalf("create product: %v", err)
	}
	engine := newTestRouter(db)

	path := "/api/v1/products/" + uintText(product.ID) + "/skus/" + uintText(product.Skus[0].ID)
	w := performJSON(engine, http.MethodDelete, path, token, nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthMenusResponseIncludesRouteName(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "系统管理员", "system_admin")
	menu := testutil.CreateMenu(t, db, model.Menu{Title: "系统管理", RouteName: "System", Status: 1})
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachMenuToRole(t, db, role.ID, menu.ID)
	auth := service.NewAuthService(db, "test-secret", 3600)
	token, err := auth.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	engine := newTestRouter(db)

	w := performJSON(engine, http.MethodGet, "/api/v1/auth/menus", token, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data []struct {
			Title     string `json:"title"`
			RouteName string `json:"routeName"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].RouteName != "System" {
		t.Fatalf("unexpected menus response: %#v", body.Data)
	}
}

func newTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.New(router.Options{DB: db, JWTSecret: "test-secret", JWTExpireSeconds: 3600})
}

func performJSON(engine *gin.Engine, method, path, token string, body interface{}) *httptest.ResponseRecorder {
	var payload bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&payload).Encode(body)
	}
	req := httptest.NewRequest(method, path, &payload)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	return w
}

func seedAdminWithPermission(t *testing.T, db *gorm.DB, method, path, code string) (model.AdminUser, string) {
	t.Helper()
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "测试角色", "test_"+code)
	permission := testutil.CreateAPIPermission(t, db, method, path, code)
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachPermissionToRole(t, db, role.ID, permission.ID)
	auth := service.NewAuthService(db, "test-secret", 3600)
	token, err := auth.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return admin, token
}

func uintText(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}
