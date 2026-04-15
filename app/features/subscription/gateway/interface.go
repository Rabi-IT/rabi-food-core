package subscription_gateway

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type SubscriptionGateway interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	UpsertConfig(ctx context.Context, tenantID string, input UpsertConfigInput) error
	GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error)
	List(ctx context.Context, filter ListFilter) ([]ListOutput, error)
	CreateDeliveries(ctx context.Context, inputs []CreateDeliveryInput) error
}

// SubscriptionItem holds a snapshot of a product at subscription creation time.
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
	Weekday   uint8 `json:"weekday"   validate:"min=0,max=6"`  // 0 = Sunday … 6 = Saturday
	StartHour uint8 `json:"startHour" validate:"min=0,max=23"` // 0–23
	EndHour   uint8 `json:"endHour"   validate:"min=1,max=23"` // 0–23
}

type CreateInput struct {
	UserID   string
	TenantID string
	Status   subscription.Status

	// RootID links subscription revisions. When a subscription is replaced
	// by a new configuration, the new subscription keeps the same RootID
	// to preserve history.
	RootID *string

	// Configuration

	// DeliveryDays stores the delivery schedule configuration,
	// including weekdays and time windows, used to calculate future deliveries.
	DeliveryDays []DeliveryDay
	Items        []SubscriptionItem
	Notes        string

	// Cycles represent prepaid delivery units.
	// Each cycle corresponds to one delivery.
	// When a delivery is completed, one cycle is consumed.

	// TotalCycles is the number of delivery cycles purchased
	// when the subscription was created or last modified.
	TotalCycles uint
	// RemainingCycles is how many delivery cycles are still available.
	RemainingCycles uint
	// CycleDiscount is the percentage (0–100) applied to the cycle price.
	CycleDiscount uint
	// CutoffOffsetMinutes defines how many minutes before the scheduled
	// delivery time modifications are no longer allowed.
	CutoffOffsetMinutes uint16
	AutoRenew           bool
	// MaxAttemptsPerOrder defines how many delivery attempts are allowed
	// before the delivery is considered failed.
	MaxAttemptsPerOrder uint8

	// Payment snapshot
	ItemsTotal    uint
	ItemsDiscount uint
	PaymentAmount uint
	PaymentStatus payment_status.Status
	// ExternalPaymentID references the payment transaction in the external
	// provider. Must be unique to ensure idempotent payment processing.
	ExternalPaymentID string
}

type GetByIDFilter struct {
	ID       string
	TenantID string
	UserID   string
}

type GetByIDOutput struct {
	ID       string              `json:"id"`
	TenantID string              `json:"tenantId"`
	UserID   string              `json:"userId"`
	Status   subscription.Status `json:"status"`

	DeliveryDays []DeliveryDay      `json:"deliveryDays"`
	Items        []SubscriptionItem `json:"items"`
	Notes        string             `json:"notes"`

	TotalCycles         uint   `json:"totalCycles"`
	RemainingCycles     uint   `json:"remainingCycles"`
	CycleDiscount       uint   `json:"cycleDiscount"`
	CutoffOffsetMinutes uint16 `json:"cutoffOffsetMinutes"`
	AutoRenew           bool   `json:"autoRenew"`

	ItemsTotal    uint `json:"itemsTotal"`
	ItemsDiscount uint `json:"itemsDiscount"`
	PaymentAmount uint `json:"paymentAmount"`

	CreatedAt time.Time `json:"createdAt"`
}

type PatchFilter struct {
	ID       string
	TenantID string
	Statuses []subscription.Status
}

type PatchValues struct {
	Status subscription.Status
}

type PaginateFilter struct {
	UserID   *string
	TenantID *string
	Status   subscription.Status
}

type PaginateData struct {
	ID              string              `json:"id"`
	TenantID        string              `json:"tenantId"`
	UserID          string              `json:"userId"`
	Status          subscription.Status `json:"status"`
	TotalCycles     int                 `json:"totalCycles"`
	RemainingCycles int                 `json:"remainingCycles"`
	PaymentAmount   uint                `json:"paymentAmount"`
	AutoRenew       bool                `json:"autoRenew"`
	CreatedAt       time.Time           `json:"createdAt"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type CreateDeliveryInput struct {
	SubscriptionID string
	ScheduledDate  time.Time
	StartHour      uint8
	EndHour        uint8
	CutoffAt       time.Time
	Status         subscription.DeliveryStatus
	MaxAttempts    uint8
}

type ListFilter struct {
	ID                 *string
	Status             subscription.Status
	RemainingCyclesGTE uint
	Weekday            time.Weekday
}

type ListOutput struct {
	ID                  string
	DeliveryDays        []DeliveryDay
	CutoffOffsetMinutes uint16
	MaxAttemptsPerOrder uint8
}

type DiscountRule struct {
	CyclesThreshold uint  `json:"cyclesThreshold" validate:"min=1"`
	Discount        uint8 `json:"discount"        validate:"min=1,max=100"` // percentage (1–100)
}

type UpsertConfigInput struct {
	TenantID            string `json:"-"`
	MaxAttemptsPerOrder uint8
	DiscountRules       []DiscountRule
	CutoffOffsetMinutes uint16
	IsOpen              bool
}

type GetConfigOutput struct {
	MaxAttemptsPerOrder uint8          `json:"maxAttemptsPerOrder"`
	DiscountRules       []DiscountRule `json:"discountRules"`
	CutoffOffsetMinutes uint16         `json:"cutoffOffsetMinutes"`
	IsOpen              bool           `json:"isOpen"`
}
