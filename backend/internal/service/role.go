package service

import (
	"errors"

	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type RoleService struct {
	repo *repository.Repository
}

type RoleInput struct {
	Name        string `json:"name" binding:"required,max=64"`
	Code        string `json:"code" binding:"required,max=64"`
	Status      int    `json:"status" binding:"oneof=0 1"`
	Description string `json:"description" binding:"max=255"`
}

type UpdateRolePermissionsInput struct {
	MenuIDs          []uint `json:"menuIds"`
	APIPermissionIDs []uint `json:"apiPermissionIds"`
}

type RolePermissionDetail struct {
	RoleID           uint   `json:"roleId"`
	MenuIDs          []uint `json:"menuIds"`
	APIPermissionIDs []uint `json:"apiPermissionIds"`
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{repo: repository.New(db)}
}

func (s *RoleService) List(keyword string, status *int, page, pageSize int) ([]model.Role, int64, error) {
	var items []model.Role
	query := s.repo.DB().Model(&model.Role{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, Internal("查询角色失败", err)
	}
	if err := query.Order("id DESC").Scopes(repository.Paginate(page, pageSize)).Find(&items).Error; err != nil {
		return nil, 0, Internal("查询角色失败", err)
	}
	return items, total, nil
}

func (s *RoleService) Create(input RoleInput) (model.Role, error) {
	role := model.Role{Name: input.Name, Code: input.Code, Status: input.Status, Description: input.Description}
	if role.Status == 0 {
		role.Status = input.Status
	}
	if err := s.repo.DB().Create(&role).Error; err != nil {
		if isDuplicate(err) {
			return model.Role{}, Conflict("角色编码已存在")
		}
		return model.Role{}, Internal("创建角色失败", err)
	}
	return role, nil
}

func (s *RoleService) Update(id uint, input RoleInput) (model.Role, error) {
	var role model.Role
	if err := s.repo.DB().First(&role, id).Error; err != nil {
		return model.Role{}, mapNotFound(err, "角色不存在")
	}
	role.Name = input.Name
	role.Code = input.Code
	role.Status = input.Status
	role.Description = input.Description
	if err := s.repo.DB().Save(&role).Error; err != nil {
		if isDuplicate(err) {
			return model.Role{}, Conflict("角色编码已存在")
		}
		return model.Role{}, Internal("更新角色失败", err)
	}
	return role, nil
}

func (s *RoleService) Delete(id uint) error {
	if err := s.repo.DB().Delete(&model.Role{}, id).Error; err != nil {
		return Internal("删除角色失败", err)
	}
	return nil
}

func (s *RoleService) PermissionDetail(roleID uint) (RolePermissionDetail, error) {
	var role model.Role
	if err := s.repo.DB().First(&role, roleID).Error; err != nil {
		return RolePermissionDetail{}, mapNotFound(err, "角色不存在")
	}
	detail := RolePermissionDetail{RoleID: roleID}
	if err := s.repo.DB().Model(&model.RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &detail.MenuIDs).Error; err != nil {
		return RolePermissionDetail{}, Internal("查询菜单权限失败", err)
	}
	if err := s.repo.DB().Model(&model.RoleAPIPermission{}).Where("role_id = ?", roleID).Pluck("api_permission_id", &detail.APIPermissionIDs).Error; err != nil {
		return RolePermissionDetail{}, Internal("查询接口权限失败", err)
	}
	return detail, nil
}

func (s *RoleService) UpdatePermissions(roleID uint, input UpdateRolePermissionsInput) error {
	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return mapNotFound(err, "角色不存在")
		}
		if len(input.MenuIDs) > 0 {
			var count int64
			if err := tx.Model(&model.Menu{}).Where("id IN ?", input.MenuIDs).Count(&count).Error; err != nil {
				return Internal("校验菜单权限失败", err)
			}
			if int(count) != len(uniqueUInts(input.MenuIDs)) {
				return BadRequest("菜单权限不存在")
			}
		}
		if len(input.APIPermissionIDs) > 0 {
			var count int64
			if err := tx.Model(&model.APIPermission{}).Where("id IN ?", input.APIPermissionIDs).Count(&count).Error; err != nil {
				return Internal("校验接口权限失败", err)
			}
			if int(count) != len(uniqueUInts(input.APIPermissionIDs)) {
				return BadRequest("接口权限不存在")
			}
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return Internal("清理菜单权限失败", err)
		}
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleAPIPermission{}).Error; err != nil {
			return Internal("清理接口权限失败", err)
		}
		for _, id := range uniqueUInts(input.MenuIDs) {
			if err := tx.Create(&model.RoleMenu{RoleID: roleID, MenuID: id}).Error; err != nil {
				return Internal("保存菜单权限失败", err)
			}
		}
		for _, id := range uniqueUInts(input.APIPermissionIDs) {
			if err := tx.Create(&model.RoleAPIPermission{RoleID: roleID, APIPermissionID: id}).Error; err != nil {
				return Internal("保存接口权限失败", err)
			}
		}
		return nil
	})
}

func mapNotFound(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(message)
	}
	return Internal(message, err)
}

func uniqueUInts(values []uint) []uint {
	seen := map[uint]bool{}
	result := make([]uint, 0, len(values))
	for _, value := range values {
		if value == 0 || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
