package models

import (
	"time"

	"gorm.io/gorm"
)

// AppPlan represents the AppPlan model in the database.
type AppPlan struct {
	ID          string `gorm:"type:uuid"`
	Price       int    `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`

	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m AppPlan) TableName() string {
	return "app_plans"
}
