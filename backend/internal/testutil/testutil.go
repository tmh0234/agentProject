package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", name)), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(model.Models()...); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	return db
}

func SeedAdmin(t *testing.T, db *gorm.DB) model.AdminUser {
	t.Helper()
	hash, err := service.HashPassword("Admin@123456")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	admin := model.AdminUser{
		Username:     "admin",
		Nickname:     "超级管理员",
		PasswordHash: hash,
		Status:       1,
	}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	return admin
}

func CreateRole(t *testing.T, db *gorm.DB, name, code string) model.Role {
	t.Helper()
	role := model.Role{Name: name, Code: code, Status: 1}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("create role: %v", err)
	}
	return role
}

func MenuFixture(title, permission string) model.Menu {
	return model.Menu{Title: title, Permission: permission, Status: 1}
}

func CreateMenu(t *testing.T, db *gorm.DB, menu model.Menu) model.Menu {
	t.Helper()
	if menu.Status == 0 {
		menu.Status = 1
	}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("create menu: %v", err)
	}
	return menu
}

func CreateAPIPermission(t *testing.T, db *gorm.DB, method, path, code string) model.APIPermission {
	t.Helper()
	permission := model.APIPermission{Method: method, Path: path, Code: code, Description: code}
	if err := db.Create(&permission).Error; err != nil {
		t.Fatalf("create api permission: %v", err)
	}
	return permission
}

func AttachRoleToAdmin(t *testing.T, db *gorm.DB, adminID, roleID uint) {
	t.Helper()
	if err := db.Create(&model.AdminUserRole{AdminUserID: adminID, RoleID: roleID}).Error; err != nil {
		t.Fatalf("attach role to admin: %v", err)
	}
}

func AttachPermissionToRole(t *testing.T, db *gorm.DB, roleID, permissionID uint) {
	t.Helper()
	if err := db.Create(&model.RoleAPIPermission{RoleID: roleID, APIPermissionID: permissionID}).Error; err != nil {
		t.Fatalf("attach permission to role: %v", err)
	}
}

func AttachMenuToRole(t *testing.T, db *gorm.DB, roleID, menuID uint) {
	t.Helper()
	if err := db.Create(&model.RoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
		t.Fatalf("attach menu to role: %v", err)
	}
}

func CreateCategory(t *testing.T, db *gorm.DB, name string) model.ProductCategory {
	t.Helper()
	category := model.ProductCategory{Name: name, Status: 1}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	return category
}

func CreateBrand(t *testing.T, db *gorm.DB, name string) model.Brand {
	t.Helper()
	brand := model.Brand{Name: name, Status: 1}
	if err := db.Create(&brand).Error; err != nil {
		t.Fatalf("create brand: %v", err)
	}
	return brand
}
