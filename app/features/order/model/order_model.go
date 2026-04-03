package model

import (
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/order"
	subscription_model "github.com/Rabi-IT/rabi-food-core/features/subscription/model"
	tenant_model "github.com/Rabi-IT/rabi-food-core/features/tenant/model"
	user_model "github.com/Rabi-IT/rabi-food-core/features/user/model"

	"gorm.io/datatypes"
)

// Order represents the Order model in the database.
type Order struct {
	ID             string `gorm:"type:uuid"`
	TenantID       string `gorm:"type:uuid;not null"`
	Tenant         tenant_model.Tenant
	UserID         string `gorm:"type:uuid;not null"`
	User           user_model.User
	SubscriptionID *string `gorm:"type:uuid"`
	Subscription   *subscription_model.Subscription

	Code              string                  `gorm:"uniqueIndex;not null"`
	DeliveryStatus    order.DeliveryStatus    `gorm:"type:varchar(20);not null"`
	FulfillmentStatus order.FulfillmentStatus `gorm:"type:varchar(20);not null"`
	Notes             string                  `gorm:"type:text"`
	TotalPrice        uint                    `gorm:"not null"`

	PaymentStatus     payment_status.Status `gorm:"type:varchar(20);not null"`
	ExternalPaymentID *string               `gorm:"type:varchar(100);uniqueIndex"`
	PaidAt            *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`

	Items datatypes.JSON `gorm:"not null"`
}

func (Order) TableName() string {
	return "orders"
}
