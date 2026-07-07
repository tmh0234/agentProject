package service

import (
	"github.com/mall-admin/backend/internal/model"
	"github.com/mall-admin/backend/internal/repository"
	"gorm.io/gorm"
)

type AdminUserService struct {
	repo *repository.Repository
}

type AdminUserInput struct {
	Username string `json:"username" binding:"required,max=64"`
	Nickname string `json:"nickname" binding:"required,max=64"`
	Password string `json:"password" binding:"omitempty,min=8,max=64"`
	Status   int    `json:"status" binding:"oneof=0 1"`
	RoleIDs  []uint `json:"roleIds"`
}

type ResetPasswordInput struct {
	Password string `json:"password" binding:"required,min=8,max=64"`
}

func NewAdminUserService(db *gorm.DB) *AdminUserService {
	return &AdminUserService{repo: repository.New(db)}
}

func (s *AdminUserService) List(keyword string, status *int, page, pageSize int) ([]model.AdminUser, int64, error) {
	var items []model.AdminUser
	query := s.repo.DB().Model(&model.AdminUser{}).Preload("Roles")
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, Internal("查询管理员失败", err)
	}
	if err := query.Order("id DESC").Scopes(repository.Paginate(page, pageSize)).Find(&items).Error; err != nil {
		return nil, 0, Internal("查询管理员失败", err)
	}
	return items, total, nil
}

func (s *AdminUserService) Create(input AdminUserInput) (model.AdminUser, error) {
	if input.Password == "" {
		return model.AdminUser{}, BadRequest("密码不能为空")
	}
	hash, err := HashPassword(input.Password)
	if err != nil {
		return model.AdminUser{}, Internal("密码加密失败", err)
	}
	user := model.AdminUser{Username: input.Username, Nickname: input.Nickname, PasswordHash: hash, Status: input.Status}
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("用户名已存在")
			}
			return Internal("创建管理员失败", err)
		}
		return replaceUserRoles(tx, user.ID, input.RoleIDs)
	})
	if err != nil {
		return model.AdminUser{}, err
	}
	_ = s.repo.DB().Preload("Roles").First(&user, user.ID).Error
	return user, nil
}

func (s *AdminUserService) Update(id uint, input AdminUserInput) (model.AdminUser, error) {
	var user model.AdminUser
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return mapNotFound(err, "管理员不存在")
		}
		user.Username = input.Username
		user.Nickname = input.Nickname
		user.Status = input.Status
		if err := tx.Save(&user).Error; err != nil {
			if isDuplicate(err) {
				return Conflict("用户名已存在")
			}
			return Internal("更新管理员失败", err)
		}
		return replaceUserRoles(tx, user.ID, input.RoleIDs)
	})
	if err != nil {
		return model.AdminUser{}, err
	}
	_ = s.repo.DB().Preload("Roles").First(&user, id).Error
	return user, nil
}

func (s *AdminUserService) ResetPassword(id uint, password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return Internal("密码加密失败", err)
	}
	result := s.repo.DB().Model(&model.AdminUser{}).Where("id = ?", id).Update("password_hash", hash)
	if result.Error != nil {
		return Internal("重置密码失败", result.Error)
	}
	if result.RowsAffected == 0 {
		return NotFound("管理员不存在")
	}
	return nil
}

func (s *AdminUserService) Delete(id uint) error {
	if err := s.repo.DB().Delete(&model.AdminUser{}, id).Error; err != nil {
		return Internal("删除管理员失败", err)
	}
	return nil
}

func replaceUserRoles(tx *gorm.DB, userID uint, roleIDs []uint) error {
	roleIDs = uniqueUInts(roleIDs)
	if len(roleIDs) > 0 {
		var count int64
		if err := tx.Model(&model.Role{}).Where("id IN ?", roleIDs).Count(&count).Error; err != nil {
			return Internal("校验角色失败", err)
		}
		if int(count) != len(roleIDs) {
			return BadRequest("角色不存在")
		}
	}
	if err := tx.Where("admin_user_id = ?", userID).Delete(&model.AdminUserRole{}).Error; err != nil {
		return Internal("清理管理员角色失败", err)
	}
	for _, roleID := range roleIDs {
		if err := tx.Create(&model.AdminUserRole{AdminUserID: userID, RoleID: roleID}).Error; err != nil {
			return Internal("保存管理员角色失败", err)
		}
	}
	return nil
}
