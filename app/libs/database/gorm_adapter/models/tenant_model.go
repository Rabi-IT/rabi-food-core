package models

// Tenant represents the tenant model in the database.
type Tenant struct {
	ID   string `gorm:"type:uuid"`
	Name string `gorm:"not null"`

	AppPlanID string `gorm:"type:uuid"`
	AppPlan   AppPlan
}

func (m Tenant) TableName() string {
	return "tenants"
}
