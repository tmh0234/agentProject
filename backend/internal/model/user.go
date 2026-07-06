package model

import (
	"time"

	"gorm.io/gorm"
)

type AdminUser struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Nickname     string         `gorm:"size:64;not null" json:"nickname"`
	PasswordHash string         `gorm:"size:128;not null" json:"-"`
	Status       int            `gorm:"not null" json:"status"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	Roles        []Role         `gorm:"many2many:admin_user_roles;" json:"roles,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type AdminUserRole struct {
	AdminUserID uint `gorm:"primaryKey"`
	RoleID      uint `gorm:"primaryKey"`
}
