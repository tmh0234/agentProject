package service_test

import (
	"testing"

	"github.com/mall-admin/backend/internal/service"
	"github.com/mall-admin/backend/internal/testutil"
)

func TestProductCreateSyncsPriceRangeFromSKUs(t *testing.T) {
	db := testutil.NewTestDB(t)
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")

	products := service.NewProductService(db)
	product, err := products.Create(service.ProductInput{
		Name:       "旗舰手机",
		Code:       "SPU-001",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		MainImage:  "https://example.com/phone.png",
		Status:     0,
		Skus: []service.SKUInput{
			{Code: "SKU-001", PriceCents: 259900, Stock: 10, Status: 1, Specs: map[string]string{"颜色": "黑色"}},
			{Code: "SKU-002", PriceCents: 199900, Stock: 5, Status: 1, Specs: map[string]string{"颜色": "白色"}},
		},
	})
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}
	if product.MinPriceCents != 199900 || product.MaxPriceCents != 259900 {
		t.Fatalf("unexpected price range: min=%d max=%d", product.MinPriceCents, product.MaxPriceCents)
	}
}

func TestProductCreateRejectsDuplicateSKUCode(t *testing.T) {
	db := testutil.NewTestDB(t)
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")

	products := service.NewProductService(db)
	_, err := products.Create(service.ProductInput{
		Name:       "旗舰手机",
		Code:       "SPU-001",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Skus: []service.SKUInput{
			{Code: "SKU-001", PriceCents: 259900, Stock: 10, Status: 1},
			{Code: "SKU-001", PriceCents: 199900, Stock: 5, Status: 1},
		},
	})
	if err == nil {
		t.Fatal("expected duplicate sku error")
	}
	if !service.IsConflict(err) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestProductUpdateCanKeepExistingSKUCode(t *testing.T) {
	db := testutil.NewTestDB(t)
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")

	products := service.NewProductService(db)
	product, err := products.Create(service.ProductInput{
		Name:       "旗舰手机",
		Code:       "SPU-001",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Skus: []service.SKUInput{
			{Code: "SKU-001", PriceCents: 259900, Stock: 10, Status: 1},
		},
	})
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	updated, err := products.Update(product.ID, service.ProductInput{
		Name:       "旗舰手机 Pro",
		Code:       "SPU-001",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Skus: []service.SKUInput{
			{Code: "SKU-001", PriceCents: 299900, Stock: 12, Status: 1},
		},
	})
	if err != nil {
		t.Fatalf("update product failed: %v", err)
	}
	if updated.MinPriceCents != 299900 || len(updated.Skus) != 1 {
		t.Fatalf("unexpected updated product: %#v", updated)
	}
}

func TestProductPublishRequiresEnabledSKUWithValidStock(t *testing.T) {
	db := testutil.NewTestDB(t)
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")

	products := service.NewProductService(db)
	product, err := products.Create(service.ProductInput{
		Name:       "待上架商品",
		Code:       "SPU-002",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Status:     0,
		Skus: []service.SKUInput{
			{Code: "SKU-010", PriceCents: 1000, Stock: 0, Status: 1},
			{Code: "SKU-011", PriceCents: 1200, Stock: 4, Status: 0},
		},
	})
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	err = products.UpdateStatus(product.ID, 1)
	if err == nil {
		t.Fatal("expected publish validation error")
	}
	if !service.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}

	if err := products.UpdateSKU(product.ID, product.Skus[0].ID, service.SKUInput{
		Code:       "SKU-010",
		PriceCents: 1000,
		Stock:      8,
		Status:     1,
	}); err != nil {
		t.Fatalf("update sku failed: %v", err)
	}
	if err := products.UpdateStatus(product.ID, 1); err != nil {
		t.Fatalf("publish failed after valid sku: %v", err)
	}
}

func TestProductDeleteSKURejectsRemovingLastValidSKUFromPublishedProduct(t *testing.T) {
	db := testutil.NewTestDB(t)
	category := testutil.CreateCategory(t, db, "手机")
	brand := testutil.CreateBrand(t, db, "Acme")

	products := service.NewProductService(db)
	product, err := products.Create(service.ProductInput{
		Name:       "已上架商品",
		Code:       "SPU-003",
		CategoryID: category.ID,
		BrandID:    brand.ID,
		Status:     1,
		Skus: []service.SKUInput{
			{Code: "SKU-020", PriceCents: 1000, Stock: 8, Status: 1},
		},
	})
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	err = products.DeleteSKU(product.ID, product.Skus[0].ID)
	if err == nil {
		t.Fatal("expected delete sku validation error")
	}
	if !service.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}

	loaded, loadErr := products.Get(product.ID)
	if loadErr != nil {
		t.Fatalf("load product failed: %v", loadErr)
	}
	if len(loaded.Skus) != 1 {
		t.Fatalf("expected sku to remain after rejected delete, got %d", len(loaded.Skus))
	}
}
