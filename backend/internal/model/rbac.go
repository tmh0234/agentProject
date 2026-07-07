package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"size:64;not null" json:"name"`
	Code        string          `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Status      int             `gorm:"not null" json:"status"`
	Description string          `gorm:"size:255" json:"description"`
	Menus       []Menu          `gorm:"many2many:role_menus;" json:"menus,omitempty"`
	Permissions []APIPermission `gorm:"many2many:role_api_permissions;" json:"apiPermissions,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

type Menu struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ParentID   uint           `gorm:"not null;default:0;index" json:"parentId"`
	Title      string         `gorm:"size:64;not null" json:"title"`
	RouteName  string         `gorm:"size:64" json:"routeName"`
	Path       string         `gorm:"size:128" json:"path"`
	Component  string         `gorm:"size:128" json:"component"`
	Icon       string         `gorm:"size:64" json:"icon"`
	Permission string         `gorm:"size:128" json:"permission"`
	Sort       int            `gorm:"not null;default:0" json:"sort"`
	Status     int            `gorm:"not null" json:"status"`
	Children   []Menu         `gorm:"-" json:"children,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type APIPermission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Method      string         `gorm:"size:16;not null;index:idx_api_method_path,unique" json:"method"`
	Path        string         `gorm:"size:255;not null;index:idx_api_method_path,unique" json:"path"`
	Code        string         `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type RoleMenu struct {
	RoleID uint `gorm:"primaryKey"`
	MenuID uint `gorm:"primaryKey"`
}

type RoleAPIPermission struct {
	RoleID          uint `gorm:"primaryKey"`
	APIPermissionID uint `gorm:"primaryKey"`
}
