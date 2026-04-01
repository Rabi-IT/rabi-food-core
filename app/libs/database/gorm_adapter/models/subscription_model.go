package models

import (
	"rabi-food-core/domain/payment"
	"rabi-food-core/domain/subscription"
	"time"

	"gorm.io/datatypes"
)

const (
	UniqueIndex__SubscriptionDelivery = "subscription_delivery_unique"
)

type SubscriptionItem struct {
	ItemID    string `json:"itemId"`
	Name      string `json:"name"`
	Quantity  uint   `json:"quantity"`
	UnitPrice uint   `json:"unitPrice"`
	Subtotal  uint   `json:"subtotal"`
	Discount  uint   `json:"discount"`
	Total     uint   `json:"total"`
}

type DeliveryDay struct {
	Weekday uint8 `json:"weekday"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday

	StartHour uint8 `json:"startHour"` // format 0-23 representing the hour of the day
	EndHour   uint8 `json:"endHour"`   // format 0-23 representing the hour of the day
}

type Subscription struct {
	ID string `gorm:"type:uuid;primaryKey"`

	// RootSubscriptionID links subscription revisions when a subscription is modified.
	// When a subscription is replaced by a new configuration, the new
	// subscription keeps the same RootSubscriptionID to preserve history.
	RootSubscriptionID *string `gorm:"type:uuid;index"`

	TenantID string `gorm:"type:uuid;index;not null"`
	Tenant   Tenant `gorm:"not null"`
	UserID   string `gorm:"type:uuid;index;not null"`
	User     User   `gorm:"not null"`

	// Configuration

	// DeliveryDays stores the detailed delivery schedule configuration,
	// including weekdays and time windows. This is used to calculate
	// future deliveries.
	DeliveryDays datatypes.JSON `gorm:"type:jsonb"`

	// DeliveryWeekdaysMask is a bitmask representing the weekdays when
	// deliveries may occur. It allows efficient filtering without parsing JSON.
	//
	// Example:
	// Monday = 1 << 1
	// Wednesday = 1 << 3
	DeliveryWeekdaysMask uint8          `gorm:"not null;default:0;index"`
	Items                datatypes.JSON `gorm:"type:jsonb"`
	Notes                string         `gorm:"type:text"`

	// Cycles represent prepaid delivery units.
	//
	// Each cycle corresponds to one delivery.
	// When a delivery is completed, one cycle is consumed.
	//
	// Cycles also determine the remaining monetary value of a subscription
	// when performing plan upgrades or downgrades.

	// TotalCycles represents the number of delivery cycles purchased
	// when the subscription was created or last modified.
	TotalCycles uint `gorm:"not null"`

	// RemainingCycles represents how many delivery cycles are still available.
	// Each successful delivery consumes one cycle.
	RemainingCycles uint `gorm:"not null;index"`

	// CycleDiscount represents the percentage discount applied to the
	// subscription cycle price when calculating the total payment amount.
	// The value ranges from 0 to 100.
	//
	// Example:
	// CyclePrice = 1000
	// CycleDiscount = 10
	// Final price per cycle = 900
	CycleDiscount uint `gorm:"not null;default:0"`

	// CutoffOffsetMinutes defines how many minutes before the scheduled
	// delivery time modifications are no longer allowed.
	CutoffOffsetMinutes uint16 `gorm:"not null;default:0"`
	AutoRenew           bool   `gorm:"not null;default:true"`

	// MaxAttemptsPerOrder defines how many delivery attempts are allowed
	// before the delivery is considered failed.
	MaxAttemptsPerOrder uint8 `gorm:"not null;default:1"`

	// Payment

	ItemsTotal    uint           `gorm:"not null;default:0"`
	ItemsDiscount uint           `gorm:"not null;default:0"`
	PaymentAmount uint           `gorm:"not null"`
	PaymentStatus payment.Status `gorm:"not null;type:varchar(20);index"`

	// ExternalPaymentID references the payment transaction created
	// in the external payment provider. It must be unique to ensure
	// idempotent payment processing.
	ExternalPaymentID string `gorm:"not null;uniqueIndex"`

	Status subscription.Status `gorm:"type:varchar(20);index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type SubscriptionDelivery struct {
	ID string `gorm:"type:uuid;primaryKey"`

	SubscriptionID string `gorm:"type:uuid;uniqueIndex:subscription_delivery_unique;index;not null"`
	Subscription   Subscription

	// Scheduling

	ScheduledDate time.Time `gorm:"uniqueIndex:subscription_delivery_unique;index"`
	StartHour     uint8     `gorm:"not null"`
	EndHour       uint8     `gorm:"not null"`
	CutoffAt      time.Time `gorm:"not null"`

	// State

	Status subscription.DeliveryStatus `gorm:"type:varchar(20);index"`

	// Attempts

	// DeliveryAttempts tracks how many times the system attempted
	// to fulfill this delivery.
	DeliveryAttempts uint8

	// MaxDeliveryAttempts defines how many delivery attempts are allowed before
	// the delivery is considered failed. This is a snapshot from the parent subscription
	// at the time of scheduling to ensure consistent behavior even if the subscription
	// configuration changes later.
	MaxDeliveryAttempts uint8

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (SubscriptionDelivery) TableName() string {
	return "subscription_deliveries"
}

type SubscriptionConfig struct {
	TenantID string `gorm:"type:uuid;primaryKey;not null"`
	Tenant   Tenant `gorm:"not null"`

	IsOpen              bool `gorm:"not null"`
	MaxAttemptsPerOrder uint8
	DiscountRules       datatypes.JSON
	CutoffOffsetMinutes uint16
	UpdatedAt           time.Time
}

func (SubscriptionConfig) TableName() string {
	return "subscription_configs"
}
