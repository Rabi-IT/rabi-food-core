package models

// Category represents the Category model in the database.
type Category struct {
	ID       string `gorm:"type:uuid"`
	TenantID string `gorm:"type:uuid;not null"`
	Tenant   Tenant

	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
}

func (m Category) TableName() string {
	return "categories"
}
