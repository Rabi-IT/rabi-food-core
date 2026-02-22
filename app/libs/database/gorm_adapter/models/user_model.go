package models

import "rabi-food-core/domain"

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
	Role         domain.Role

	TenantID string `gorm:"type:uuid"`
	Tenant   Tenant
}

func (m User) TableName() string {
	return "users"
}
