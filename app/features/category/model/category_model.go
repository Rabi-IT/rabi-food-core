package model

import tenant_model "github.com/Rabi-IT/rabi-food-core/features/tenant/model"

// Category represents the Category model in the database.
type Category struct {
	ID       string `gorm:"type:uuid"`
	TenantID string `gorm:"type:uuid;not null"`
	Tenant   tenant_model.Tenant

	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
}

func (m Category) TableName() string {
	return "categories"
}
