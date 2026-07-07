package service_test

import (
	"testing"

	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/service"
	"github.com/mall-admin/backend/internal/testutil"
)

func TestCategoryListFiltersByKeywordAndStatus(t *testing.T) {
	db := testutil.NewTestDB(t)
	if err := db.Create(&[]model.ProductCategory{
		{Name: "手机数码", Status: 1, Sort: 1},
		{Name: "手机配件", Status: 0, Sort: 2},
		{Name: "家用电器", Status: 1, Sort: 3},
	}).Error; err != nil {
		t.Fatalf("seed categories: %v", err)
	}

	catalog := service.NewCatalogService(db)
	status := 1
	items, err := catalog.CategoryList("手机", &status)
	if err != nil {
		t.Fatalf("category list failed: %v", err)
	}
	if len(items) != 1 || items[0].Name != "手机数码" {
		t.Fatalf("unexpected filtered categories: %#v", items)
	}
}
