package service

import (
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type SystemService struct {
	repo *repository.Repository
}

func NewSystemService(db *gorm.DB) *SystemService {
	return &SystemService{repo: repository.New(db)}
}

func (s *SystemService) MenuTree() ([]model.Menu, error) {
	var menus []model.Menu
	if err := s.repo.DB().Order("parent_id ASC, sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, Internal("查询菜单失败", err)
	}
	return BuildMenuTree(menus), nil
}

func (s *SystemService) APIPermissions() ([]model.APIPermission, error) {
	var permissions []model.APIPermission
	if err := s.repo.DB().Order("method ASC, path ASC").Find(&permissions).Error; err != nil {
		return nil, Internal("查询接口权限失败", err)
	}
	return permissions, nil
}
