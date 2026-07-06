package service_test

import (
	"testing"

	"github.com/mall-admin/backend/internal/service"
	"github.com/mall-admin/backend/internal/testutil"
)

func TestRolePermissionUpdateIsTransactional(t *testing.T) {
	db := testutil.NewTestDB(t)
	role := testutil.CreateRole(t, db, "运营", "operator")
	menu := testutil.CreateMenu(t, db, testutil.MenuFixture("商品列表", "product:list"))
	permission := testutil.CreateAPIPermission(t, db, "GET", "/api/v1/products", "product:list")
	existing := testutil.CreateAPIPermission(t, db, "GET", "/api/v1/brands", "brand:list")
	testutil.AttachPermissionToRole(t, db, role.ID, existing.ID)

	roles := service.NewRoleService(db)

	err := roles.UpdatePermissions(role.ID, service.UpdateRolePermissionsInput{
		MenuIDs:          []uint{menu.ID},
		APIPermissionIDs: []uint{permission.ID, 999999},
	})
	if err == nil {
		t.Fatal("expected invalid permission error")
	}

	detail, detailErr := roles.PermissionDetail(role.ID)
	if detailErr != nil {
		t.Fatalf("permission detail failed: %v", detailErr)
	}
	if len(detail.APIPermissionIDs) != 1 || detail.APIPermissionIDs[0] != existing.ID {
		t.Fatalf("expected original permissions to remain after rollback, got %#v", detail.APIPermissionIDs)
	}
	if len(detail.MenuIDs) != 0 {
		t.Fatalf("expected menu permissions to remain empty after rollback, got %#v", detail.MenuIDs)
	}

	err = roles.UpdatePermissions(role.ID, service.UpdateRolePermissionsInput{
		MenuIDs:          []uint{menu.ID},
		APIPermissionIDs: []uint{permission.ID},
	})
	if err != nil {
		t.Fatalf("update permissions failed: %v", err)
	}

	detail, detailErr = roles.PermissionDetail(role.ID)
	if detailErr != nil {
		t.Fatalf("permission detail failed: %v", detailErr)
	}
	if len(detail.MenuIDs) != 1 || detail.MenuIDs[0] != menu.ID {
		t.Fatalf("unexpected menu ids: %#v", detail.MenuIDs)
	}
	if len(detail.APIPermissionIDs) != 1 || detail.APIPermissionIDs[0] != permission.ID {
		t.Fatalf("unexpected api permission ids: %#v", detail.APIPermissionIDs)
	}
}
