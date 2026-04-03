package model

import (
	"time"

	category_model "rabi-food-core/features/category/model"
	tenant_model "rabi-food-core/features/tenant/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product represents the product model in the database.
type Product struct {
	ID       string `gorm:"type:uuid"`
	TenantID string `gorm:"type:uuid;not null"`
	Tenant   tenant_model.Tenant

	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Photo       string

	CategoryID string `gorm:"type:uuid"`
	Category   category_model.Category

	DiscountRules datatypes.JSON `gorm:"type:jsonb;default:'[]'"`

	// Measurement and pricing details
	// Unit of measurement (e.g., "kg", "liter", "piece")
	Unit  string `gorm:"not null"`
	Price uint   `gorm:"not null"`

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m Product) TableName() string {
	return "products"
}
