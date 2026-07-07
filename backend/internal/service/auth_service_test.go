package service_test

import (
	"testing"

	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/service"
	"github.com/mall-admin/backend/internal/testutil"
)

func TestAuthServiceLoginValidatesPasswordAndReturnsProfileToken(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.SeedAdmin(t, db)

	auth := service.NewAuthService(db, "test-secret", 3600)

	result, err := auth.Login("admin", "Admin@123456")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected token")
	}
	if result.Profile.Username != "admin" {
		t.Fatalf("expected admin profile, got %q", result.Profile.Username)
	}

	claims, err := auth.ParseToken(result.Token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.UserID != result.Profile.ID {
		t.Fatalf("expected token user id %d, got %d", result.Profile.ID, claims.UserID)
	}
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	db := testutil.NewTestDB(t)
	testutil.SeedAdmin(t, db)

	auth := service.NewAuthService(db, "test-secret", 3600)

	_, err := auth.Login("admin", "wrong-password")
	if err == nil {
		t.Fatal("expected login error")
	}
	if !service.IsUnauthorized(err) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestRBACAllowsAndDeniesAPIByRole(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "商品管理员", "product_admin")
	permission := testutil.CreateAPIPermission(t, db, "GET", "/api/v1/products", "product:list")
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachPermissionToRole(t, db, role.ID, permission.ID)

	rbac := service.NewRBACService(db)

	allowed, err := rbac.CanAccess(admin.ID, "GET", "/api/v1/products")
	if err != nil {
		t.Fatalf("rbac check failed: %v", err)
	}
	if !allowed {
		t.Fatal("expected access to be allowed")
	}

	denied, err := rbac.CanAccess(admin.ID, "DELETE", "/api/v1/products/1")
	if err != nil {
		t.Fatalf("rbac check failed: %v", err)
	}
	if denied {
		t.Fatal("expected access to be denied")
	}
}

func TestRBACDeniesDisabledAdminWithExistingRole(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "商品管理员", "product_admin")
	permission := testutil.CreateAPIPermission(t, db, "GET", "/api/v1/products", "product:list")
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachPermissionToRole(t, db, role.ID, permission.ID)
	if err := db.Model(&admin).Update("status", 0).Error; err != nil {
		t.Fatalf("disable admin: %v", err)
	}

	rbac := service.NewRBACService(db)

	allowed, err := rbac.CanAccess(admin.ID, "GET", "/api/v1/products")
	if err != nil {
		t.Fatalf("rbac check failed: %v", err)
	}
	if allowed {
		t.Fatal("expected disabled admin to be denied")
	}
}

func TestProfileReturnsRolesAndMenus(t *testing.T) {
	db := testutil.NewTestDB(t)
	admin := testutil.SeedAdmin(t, db)
	role := testutil.CreateRole(t, db, "系统管理员", "system_admin")
	menu := testutil.CreateMenu(t, db, model.Menu{Title: "商品管理", RouteName: "Product", Permission: "product"})
	testutil.AttachRoleToAdmin(t, db, admin.ID, role.ID)
	testutil.AttachMenuToRole(t, db, role.ID, menu.ID)

	auth := service.NewAuthService(db, "test-secret", 3600)

	profile, err := auth.Profile(admin.ID)
	if err != nil {
		t.Fatalf("profile failed: %v", err)
	}
	if len(profile.Roles) != 1 || profile.Roles[0].Code != "system_admin" {
		t.Fatalf("unexpected roles: %#v", profile.Roles)
	}
	menus, err := auth.Menus(admin.ID)
	if err != nil {
		t.Fatalf("menus failed: %v", err)
	}
	if len(menus) != 1 || menus[0].Title != "商品管理" || menus[0].RouteName != "Product" {
		t.Fatalf("unexpected menus: %#v", menus)
	}
}
