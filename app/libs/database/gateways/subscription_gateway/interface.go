package subscription_gateway

import (
	"context"
	"rabi-food-core/domain/payment"
	"rabi-food-core/domain/subscription"
	"time"
)

type SubscriptionGateway interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	UpsertConfig(ctx context.Context, tenantID string, input UpsertConfigInput) (bool, error)
	GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error)
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
	RootID   *string

	// Configuration
	DeliveryDays []DeliveryDay
	Items        []SubscriptionItem
	Notes        string

	// Cycles
	TotalCycles         uint
	RemainingCycles     uint
	CycleDiscount       uint
	CutoffOffsetMinutes uint16
	AutoRenew           bool
	MaxAttemptsPerOrder uint8

	// Payment snapshot
	ItemsTotal        uint
	ItemsDiscount     uint
	PaymentAmount     uint
	PaymentStatus     payment.Status
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

type DiscountRule struct {
	CyclesThreshold uint  `json:"cyclesThreshold" validate:"min=1"`
	Discount        uint8 `json:"discount"        validate:"min=1,max=100"`
}

type UpsertConfigInput struct {
	MaxAttemptsPerOrder uint8
	DiscountRules       []DiscountRule
	CutoffOffsetMinutes uint16
}

type GetConfigOutput struct {
	MaxAttemptsPerOrder uint8          `json:"maxAttemptsPerOrder"`
	DiscountRules       []DiscountRule `json:"discountRules"`
	CutoffOffsetMinutes uint16         `json:"cutoffOffsetMinutes"`
}
