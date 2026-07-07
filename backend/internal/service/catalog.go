package service

import (
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type CatalogService struct {
	repo *repository.Repository
}

type CategoryInput struct {
	ParentID uint   `json:"parentId"`
	Name     string `json:"name" binding:"required,max=64"`
	Status   int    `json:"status" binding:"oneof=0 1"`
	Sort     int    `json:"sort"`
}

type BrandInput struct {
	Name   string `json:"name" binding:"required,max=64"`
	Logo   string `json:"logo" binding:"max=255"`
	Status int    `json:"status" binding:"oneof=0 1"`
	Sort   int    `json:"sort"`
}

func NewCatalogService(db *gorm.DB) *CatalogService {
	return &CatalogService{repo: repository.New(db)}
}

func (s *CatalogService) CategoryList(keyword string, status *int) ([]model.ProductCategory, error) {
	var items []model.ProductCategory
	query := s.repo.DB().Model(&model.ProductCategory{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Order("parent_id ASC, sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, Internal("查询商品分类失败", err)
	}
	return BuildCategoryTree(items), nil
}

func (s *CatalogService) CreateCategory(input CategoryInput) (model.ProductCategory, error) {
	item := model.ProductCategory{ParentID: input.ParentID, Name: input.Name, Status: input.Status, Sort: input.Sort}
	if err := s.repo.DB().Create(&item).Error; err != nil {
		if isDuplicate(err) {
			return model.ProductCategory{}, Conflict("分类名称已存在")
		}
		return model.ProductCategory{}, Internal("创建商品分类失败", err)
	}
	return item, nil
}

func (s *CatalogService) UpdateCategory(id uint, input CategoryInput) (model.ProductCategory, error) {
	var item model.ProductCategory
	if err := s.repo.DB().First(&item, id).Error; err != nil {
		return model.ProductCategory{}, mapNotFound(err, "商品分类不存在")
	}
	item.ParentID = input.ParentID
	item.Name = input.Name
	item.Status = input.Status
	item.Sort = input.Sort
	if err := s.repo.DB().Save(&item).Error; err != nil {
		if isDuplicate(err) {
			return model.ProductCategory{}, Conflict("分类名称已存在")
		}
		return model.ProductCategory{}, Internal("更新商品分类失败", err)
	}
	return item, nil
}

func (s *CatalogService) DeleteCategory(id uint) error {
	if err := s.repo.DB().Delete(&model.ProductCategory{}, id).Error; err != nil {
		return Internal("删除商品分类失败", err)
	}
	return nil
}

func (s *CatalogService) BrandList(keyword string, status *int, page, pageSize int) ([]model.Brand, int64, error) {
	var items []model.Brand
	query := s.repo.DB().Model(&model.Brand{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, Internal("查询品牌失败", err)
	}
	if err := query.Order("sort ASC, id DESC").Scopes(repository.Paginate(page, pageSize)).Find(&items).Error; err != nil {
		return nil, 0, Internal("查询品牌失败", err)
	}
	return items, total, nil
}

func (s *CatalogService) CreateBrand(input BrandInput) (model.Brand, error) {
	item := model.Brand{Name: input.Name, Logo: input.Logo, Status: input.Status, Sort: input.Sort}
	if err := s.repo.DB().Create(&item).Error; err != nil {
		if isDuplicate(err) {
			return model.Brand{}, Conflict("品牌名称已存在")
		}
		return model.Brand{}, Internal("创建品牌失败", err)
	}
	return item, nil
}

func (s *CatalogService) UpdateBrand(id uint, input BrandInput) (model.Brand, error) {
	var item model.Brand
	if err := s.repo.DB().First(&item, id).Error; err != nil {
		return model.Brand{}, mapNotFound(err, "品牌不存在")
	}
	item.Name = input.Name
	item.Logo = input.Logo
	item.Status = input.Status
	item.Sort = input.Sort
	if err := s.repo.DB().Save(&item).Error; err != nil {
		if isDuplicate(err) {
			return model.Brand{}, Conflict("品牌名称已存在")
		}
		return model.Brand{}, Internal("更新品牌失败", err)
	}
	return item, nil
}

func (s *CatalogService) DeleteBrand(id uint) error {
	if err := s.repo.DB().Delete(&model.Brand{}, id).Error; err != nil {
		return Internal("删除品牌失败", err)
	}
	return nil
}

func BuildCategoryTree(items []model.ProductCategory) []model.ProductCategory {
	byParent := map[uint][]model.ProductCategory{}
	for _, item := range items {
		item.Children = nil
		byParent[item.ParentID] = append(byParent[item.ParentID], item)
	}
	var attach func(uint) []model.ProductCategory
	attach = func(parentID uint) []model.ProductCategory {
		children := byParent[parentID]
		for i := range children {
			children[i].Children = attach(children[i].ID)
		}
		return children
	}
	return attach(0)
}
