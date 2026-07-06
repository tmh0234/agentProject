package service

import (
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type ProductService struct {
	repo *repository.Repository
}

type ProductInput struct {
	Name        string     `json:"name" binding:"required,max=128"`
	Code        string     `json:"code" binding:"required,max=64"`
	CategoryID  uint       `json:"categoryId" binding:"required"`
	BrandID     uint       `json:"brandId" binding:"required"`
	MainImage   string     `json:"mainImage" binding:"max=255"`
	Status      int        `json:"status" binding:"oneof=0 1"`
	Description string     `json:"description"`
	Skus        []SKUInput `json:"skus" binding:"required,min=1,dive"`
}

type SKUInput struct {
	Code       string            `json:"code" binding:"required,max=64"`
	Specs      map[string]string `json:"specs"`
	PriceCents int64             `json:"priceCents" binding:"required,min=1"`
	Stock      int               `json:"stock" binding:"min=0"`
	Status     int               `json:"status" binding:"oneof=0 1"`
}

type ProductStatusInput struct {
	Status int `json:"status" binding:"oneof=0 1"`
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{repo: repository.New(db)}
}

func (s *ProductService) List(keyword string, status *int, page, pageSize int) ([]model.Product, int64, error) {
	var items []model.Product
	query := s.repo.DB().Model(&model.Product{}).Preload("Category").Preload("Brand").Preload("Skus")
	if keyword != "" {
		query = query.Where("products.name LIKE ? OR products.code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("products.status = ?", *status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, Internal("查询商品失败", err)
	}
	if err := query.Order("products.id DESC").Scopes(repository.Paginate(page, pageSize)).Find(&items).Error; err != nil {
		return nil, 0, Internal("查询商品失败", err)
	}
	return items, total, nil
}

func (s *ProductService) Get(id uint) (model.Product, error) {
	var item model.Product
	if err := s.repo.DB().Preload("Category").Preload("Brand").Preload("Skus").First(&item, id).Error; err != nil {
		return model.Product{}, mapNotFound(err, "商品不存在")
	}
	return item, nil
}

func (s *ProductService) Create(input ProductInput) (model.Product, error) {
	if err := validateSKUInputs(input.Skus); err != nil {
		return model.Product{}, err
	}
	var product model.Product
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := ensureCategoryBrand(tx, input.CategoryID, input.BrandID); err != nil {
			return err
		}
		minPrice, maxPrice := priceRange(input.Skus)
		product = model.Product{
			Name:          input.Name,
			Code:          input.Code,
			CategoryID:    input.CategoryID,
			BrandID:       input.BrandID,
			MainImage:     input.MainImage,
			MinPriceCents: minPrice,
			MaxPriceCents: maxPrice,
			Status:        input.Status,
			Description:   input.Description,
		}
		if err := tx.Create(&product).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("商品编码已存在")
			}
			return Internal("创建商品失败", err)
		}
		for _, skuInput := range input.Skus {
			sku := skuFromInput(product.ID, skuInput)
			if err := tx.Create(&sku).Error; err != nil {
				if isDuplicate(err) {
					return Conflict("SKU 编码已存在")
				}
				return Internal("创建 SKU 失败", err)
			}
		}
		if input.Status == 1 {
			if err := validatePublish(tx, product.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return model.Product{}, err
	}
	return s.Get(product.ID)
}

func (s *ProductService) Update(id uint, input ProductInput) (model.Product, error) {
	if err := validateSKUInputs(input.Skus); err != nil {
		return model.Product{}, err
	}
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := ensureCategoryBrand(tx, input.CategoryID, input.BrandID); err != nil {
			return err
		}
		var product model.Product
		if err := tx.First(&product, id).Error; err != nil {
			return mapNotFound(err, "商品不存在")
		}
		minPrice, maxPrice := priceRange(input.Skus)
		product.Name = input.Name
		product.Code = input.Code
		product.CategoryID = input.CategoryID
		product.BrandID = input.BrandID
		product.MainImage = input.MainImage
		product.MinPriceCents = minPrice
		product.MaxPriceCents = maxPrice
		product.Status = input.Status
		product.Description = input.Description
		if err := tx.Save(&product).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("商品编码已存在")
			}
			return Internal("更新商品失败", err)
		}
		if err := tx.Unscoped().Where("product_id = ?", id).Delete(&model.ProductSKU{}).Error; err != nil {
			return Internal("清理 SKU 失败", err)
		}
		for _, skuInput := range input.Skus {
			sku := skuFromInput(id, skuInput)
			if err := tx.Create(&sku).Error; err != nil {
				if isDuplicate(err) {
					return Conflict("SKU 编码已存在")
				}
				return Internal("保存 SKU 失败", err)
			}
		}
		if input.Status == 1 {
			if err := validatePublish(tx, id); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return model.Product{}, err
	}
	return s.Get(id)
}

func (s *ProductService) Delete(id uint) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", id).Delete(&model.ProductSKU{}).Error; err != nil {
			return Internal("删除 SKU 失败", err)
		}
		if err := tx.Delete(&model.Product{}, id).Error; err != nil {
			return Internal("删除商品失败", err)
		}
		return nil
	})
}

func (s *ProductService) CreateSKU(productID uint, input SKUInput) (model.ProductSKU, error) {
	if err := validateSKUInputs([]SKUInput{input}); err != nil {
		return model.ProductSKU{}, err
	}
	var sku model.ProductSKU
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var product model.Product
		if err := tx.First(&product, productID).Error; err != nil {
			return mapNotFound(err, "商品不存在")
		}
		sku = skuFromInput(productID, input)
		if err := tx.Create(&sku).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("SKU 编码已存在")
			}
			return Internal("创建 SKU 失败", err)
		}
		return syncProductPriceRange(tx, productID)
	})
	return sku, err
}

func (s *ProductService) UpdateSKU(productID, skuID uint, input SKUInput) error {
	if err := validateSKUInputs([]SKUInput{input}); err != nil {
		return err
	}
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var sku model.ProductSKU
		if err := tx.Where("product_id = ?", productID).First(&sku, skuID).Error; err != nil {
			return mapNotFound(err, "SKU 不存在")
		}
		sku.Code = input.Code
		sku.Specs = model.JSONMap(input.Specs)
		sku.PriceCents = input.PriceCents
		sku.Stock = input.Stock
		sku.Status = input.Status
		if sku.Specs == nil {
			sku.Specs = model.JSONMap{}
		}
		if err := tx.Save(&sku).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("SKU 编码已存在")
			}
			return Internal("更新 SKU 失败", err)
		}
		return syncProductPriceRange(tx, productID)
	})
}

func (s *ProductService) DeleteSKU(productID, skuID uint) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var product model.Product
		if err := tx.First(&product, productID).Error; err != nil {
			return mapNotFound(err, "商品不存在")
		}
		if product.Status == 1 {
			validAfterDelete, err := hasValidSKUAfterDelete(tx, productID, skuID)
			if err != nil {
				return err
			}
			if !validAfterDelete {
				return BadRequest("已上架商品至少需要保留一个启用且有库存的 SKU")
			}
		}
		result := tx.Where("product_id = ?", productID).Delete(&model.ProductSKU{}, skuID)
		if result.Error != nil {
			return Internal("删除 SKU 失败", result.Error)
		}
		if result.RowsAffected == 0 {
			return NotFound("SKU 不存在")
		}
		return syncProductPriceRange(tx, productID)
	})
}

func (s *ProductService) UpdateStatus(id uint, status int) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var product model.Product
		if err := tx.First(&product, id).Error; err != nil {
			return mapNotFound(err, "商品不存在")
		}
		if status == 1 {
			if err := validatePublish(tx, id); err != nil {
				return err
			}
		}
		if err := tx.Model(&product).Update("status", status).Error; err != nil {
			return Internal("更新商品状态失败", err)
		}
		return nil
	})
}

func validateSKUInputs(skus []SKUInput) error {
	if len(skus) == 0 {
		return BadRequest("至少需要一个 SKU")
	}
	seen := map[string]bool{}
	for _, sku := range skus {
		if sku.Code == "" {
			return BadRequest("SKU 编码不能为空")
		}
		if seen[sku.Code] {
			return Conflict("SKU 编码已存在")
		}
		seen[sku.Code] = true
		if sku.PriceCents <= 0 {
			return BadRequest("SKU 价格必须大于 0")
		}
		if sku.Stock < 0 {
			return BadRequest("SKU 库存不能小于 0")
		}
	}
	return nil
}

func priceRange(skus []SKUInput) (int64, int64) {
	minPrice := skus[0].PriceCents
	maxPrice := skus[0].PriceCents
	for _, sku := range skus {
		if sku.PriceCents < minPrice {
			minPrice = sku.PriceCents
		}
		if sku.PriceCents > maxPrice {
			maxPrice = sku.PriceCents
		}
	}
	return minPrice, maxPrice
}

func skuFromInput(productID uint, input SKUInput) model.ProductSKU {
	specs := model.JSONMap(input.Specs)
	if specs == nil {
		specs = model.JSONMap{}
	}
	return model.ProductSKU{
		ProductID:  productID,
		Code:       input.Code,
		Specs:      specs,
		PriceCents: input.PriceCents,
		Stock:      input.Stock,
		Status:     input.Status,
	}
}

func ensureCategoryBrand(tx *gorm.DB, categoryID, brandID uint) error {
	var categoryCount int64
	if err := tx.Model(&model.ProductCategory{}).Where("id = ? AND status = ?", categoryID, 1).Count(&categoryCount).Error; err != nil {
		return Internal("校验商品分类失败", err)
	}
	if categoryCount == 0 {
		return BadRequest("商品分类不存在或已禁用")
	}
	var brandCount int64
	if err := tx.Model(&model.Brand{}).Where("id = ? AND status = ?", brandID, 1).Count(&brandCount).Error; err != nil {
		return Internal("校验品牌失败", err)
	}
	if brandCount == 0 {
		return BadRequest("品牌不存在或已禁用")
	}
	return nil
}

func validatePublish(tx *gorm.DB, productID uint) error {
	var count int64
	if err := tx.Model(&model.ProductSKU{}).
		Where("product_id = ? AND status = ? AND stock > ?", productID, 1, 0).
		Count(&count).Error; err != nil {
		return Internal("校验上架条件失败", err)
	}
	if count == 0 {
		return BadRequest("上架前至少需要一个启用且有库存的 SKU")
	}
	return nil
}

func hasValidSKUAfterDelete(tx *gorm.DB, productID, skuID uint) (bool, error) {
	var count int64
	if err := tx.Model(&model.ProductSKU{}).
		Where("product_id = ? AND id <> ? AND status = ? AND stock > ?", productID, skuID, 1, 0).
		Count(&count).Error; err != nil {
		return false, Internal("校验 SKU 删除条件失败", err)
	}
	return count > 0, nil
}

func syncProductPriceRange(tx *gorm.DB, productID uint) error {
	var skus []model.ProductSKU
	if err := tx.Where("product_id = ?", productID).Find(&skus).Error; err != nil {
		return Internal("查询 SKU 失败", err)
	}
	if len(skus) == 0 {
		return tx.Model(&model.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
			"min_price_cents": 0,
			"max_price_cents": 0,
		}).Error
	}
	minPrice := skus[0].PriceCents
	maxPrice := skus[0].PriceCents
	for _, sku := range skus {
		if sku.PriceCents < minPrice {
			minPrice = sku.PriceCents
		}
		if sku.PriceCents > maxPrice {
			maxPrice = sku.PriceCents
		}
	}
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"min_price_cents": minPrice,
		"max_price_cents": maxPrice,
	}).Error; err != nil {
		return Internal("同步商品价格失败", err)
	}
	return nil
}
