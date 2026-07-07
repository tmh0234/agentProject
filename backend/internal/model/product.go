package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	ParentID  uint              `gorm:"not null;default:0;index" json:"parentId"`
	Name      string            `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Status    int               `gorm:"not null" json:"status"`
	Sort      int               `gorm:"not null;default:0" json:"sort"`
	Children  []ProductCategory `gorm:"-" json:"children,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	DeletedAt gorm.DeletedAt    `gorm:"index" json:"-"`
}

type Brand struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Logo      string         `gorm:"size:255" json:"logo"`
	Status    int            `gorm:"not null" json:"status"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	Name          string          `gorm:"size:128;not null" json:"name"`
	Code          string          `gorm:"size:64;not null;uniqueIndex" json:"code"`
	CategoryID    uint            `gorm:"not null;index" json:"categoryId"`
	BrandID       uint            `gorm:"not null;index" json:"brandId"`
	MainImage     string          `gorm:"size:255" json:"mainImage"`
	MinPriceCents int64           `gorm:"not null;default:0" json:"minPriceCents"`
	MaxPriceCents int64           `gorm:"not null;default:0" json:"maxPriceCents"`
	Status        int             `gorm:"not null" json:"status"`
	Description   string          `gorm:"type:text" json:"description"`
	Category      ProductCategory `json:"category,omitempty"`
	Brand         Brand           `json:"brand,omitempty"`
	Skus          []ProductSKU    `json:"skus,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

type ProductSKU struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProductID  uint           `gorm:"not null;index" json:"productId"`
	Code       string         `gorm:"size:64;not null;uniqueIndex" json:"code"`
	Specs      JSONMap        `gorm:"type:json;not null" json:"specs"`
	PriceCents int64          `gorm:"not null" json:"priceCents"`
	Stock      int            `gorm:"not null;default:0" json:"stock"`
	Status     int            `gorm:"not null" json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func Models() []interface{} {
	return []interface{}{
		&AdminUser{},
		&Role{},
		&Menu{},
		&APIPermission{},
		&AdminUserRole{},
		&RoleMenu{},
		&RoleAPIPermission{},
		&ProductCategory{},
		&Brand{},
		&Product{},
		&ProductSKU{},
	}
}
