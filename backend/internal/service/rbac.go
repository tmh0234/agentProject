package service

import (
	"strings"

	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type RBACService struct {
	repo *repository.Repository
}

func NewRBACService(db *gorm.DB) *RBACService {
	return &RBACService{repo: repository.New(db)}
}

func (s *RBACService) CanAccess(userID uint, method, path string) (bool, error) {
	method = strings.ToUpper(method)
	var count int64
	err := s.repo.DB().Model(&model.APIPermission{}).
		Joins("JOIN role_api_permissions ON role_api_permissions.api_permission_id = api_permissions.id").
		Joins("JOIN admin_user_roles ON admin_user_roles.role_id = role_api_permissions.role_id").
		Joins("JOIN admin_users ON admin_users.id = admin_user_roles.admin_user_id AND admin_users.deleted_at IS NULL").
		Joins("JOIN roles ON roles.id = admin_user_roles.role_id AND roles.deleted_at IS NULL").
		Where("admin_user_roles.admin_user_id = ? AND admin_users.status = ? AND roles.status = ? AND api_permissions.method = ? AND api_permissions.path = ?", userID, 1, 1, method, path).
		Count(&count).Error
	if err != nil {
		return false, Internal("权限校验失败", err)
	}
	return count > 0, nil
}
