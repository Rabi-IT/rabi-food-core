package models

import (
	"rabi-food-core/domain/order"
	"rabi-food-core/domain/payment"
	"time"

	"gorm.io/datatypes"
)

// Order represents the Order model in the database.
type Order struct {
	ID             string `gorm:"type:uuid"`
	TenantID       string `gorm:"type:uuid;not null"`
	Tenant         Tenant
	UserID         string `gorm:"type:uuid;not null"`
	User           User
	SubscriptionID *string `gorm:"type:uuid"`
	Subscription   *Subscription

	Code              string                  `gorm:"uniqueIndex;not null"`
	DeliveryStatus    order.DeliveryStatus    `gorm:"type:varchar(20);not null"`
	FulfillmentStatus order.FulfillmentStatus `gorm:"type:varchar(20);not null"`
	Notes             string                  `gorm:"type:text"`
	TotalPrice        uint                    `gorm:"not null"`

	PaymentStatus     payment.Status `gorm:"type:varchar(20);not null"`
	ExternalPaymentID *string        `gorm:"type:varchar(100);uniqueIndex"`
	PaidAt            *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	Items datatypes.JSON `gorm:"not null"`
}

func (Order) TableName() string {
	return "orders"
}
