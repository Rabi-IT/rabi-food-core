package model

import (
	"github.com/Rabi-IT/rabi-food-core/domain/auth"
	tenant_model "github.com/Rabi-IT/rabi-food-core/features/tenant/model"
)

// User represents the user model in the database.
type User struct {
	ID           string `gorm:"type:uuid"`
	SocialID     string
	Street       string
	Complement   string
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	Photo        string
	TaxID        string `gorm:"not null"`
	Phone        string `gorm:"not null"`
	City         string
	State        string
	ZIP          string `gorm:"column:zip"`
	Neighborhood string
	Role         auth.Role

	TenantID string `gorm:"type:uuid"`
	Tenant   tenant_model.Tenant
}

func (m User) TableName() string {
	return "users"
}
