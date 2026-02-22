package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents the product model in the database.
type Product struct {
	ID       string `gorm:"type:uuid"`
	TenantID string `gorm:"type:uuid;not null"`
	Tenant   Tenant

	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Photo       string

	CategoryID string `gorm:"type:uuid"`
	Category   Category

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
