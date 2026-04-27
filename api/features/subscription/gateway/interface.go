package subscription_gateway

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/subscription"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type SubscriptionGateway interface {
	BulkCreate(ctx context.Context, inputs []CreateInput) ([]string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	UpsertConfig(ctx context.Context, tenantID string, input UpsertConfigInput) error
	GetConfig(ctx context.Context, tenantID string) (*GetConfigOutput, error)
	List(ctx context.Context, filter ListFilter) ([]ListOutput, error)
	CreateDeliveries(ctx context.Context, inputs []CreateDeliveryInput) error
}

type DeliveryDay struct {
	Weekday   uint8 `json:"weekday"   validate:"min=0,max=6"`  // 0 = Sunday … 6 = Saturday
	StartHour uint8 `json:"startHour" validate:"min=0,max=23"` // 0–23
	EndHour   uint8 `json:"endHour"   validate:"min=1,max=23"` // 0–23
}

type CreateInput struct {
	UserID              string
	TenantID            string
	Status              subscription.Status
	RootID              *string
	SubscriptionGroupID *string

	// Product (1:1 per subscription row)
	ProductID string
	Quantity  uint16
	UnitPrice uint

	// Configuration
	DeliveryDays        []DeliveryDay
	Notes               string
	TotalCycles         uint
	RemainingCycles     uint
	CycleDiscount       uint
	CutoffOffsetMinutes uint16
	AutoRenew           bool
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
	ID                  string              `db:"id"                   json:"id"`
	TenantID            string              `db:"tenant_id"            json:"tenantId"`
	UserID              string              `db:"user_id"              json:"userId"`
	Status              subscription.Status `db:"status"               json:"status"`
	SubscriptionGroupID *string             `db:"subscription_group_id" json:"subscriptionGroupId,omitempty"`

	ProductID string `db:"product_id" json:"productId"`
	Quantity  uint16 `db:"quantity"   json:"quantity"`
	UnitPrice uint   `db:"unit_price" json:"unitPrice"`

	DeliveryDays []DeliveryDay `db:"delivery_days" json:"deliveryDays"`
	Notes        string        `db:"notes"         json:"notes"`

	TotalCycles         uint   `db:"total_cycles"          json:"totalCycles"`
	RemainingCycles     uint   `db:"remaining_cycles"      json:"remainingCycles"`
	CycleDiscount       uint   `db:"cycle_discount"        json:"cycleDiscount"`
	CutoffOffsetMinutes uint16 `db:"cutoff_offset_minutes" json:"cutoffOffsetMinutes"`
	AutoRenew           bool   `db:"auto_renew"            json:"autoRenew"`

	ItemsTotal    uint `db:"items_total"    json:"itemsTotal"`
	ItemsDiscount uint `db:"items_discount" json:"itemsDiscount"`
	PaymentAmount uint `db:"payment_amount" json:"paymentAmount"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
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
	UserID   string
	TenantID string
	Status   subscription.Status
}

type PaginateData struct {
	ID              string              `json:"id"`
	TenantID        string              `json:"tenantId"`
	UserID          string              `json:"userId"`
	ProductID       string              `json:"productId"`
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
	ID                 string
	Status             subscription.Status
	RemainingCyclesGTE uint
	Weekday            time.Weekday
}

type ListOutput struct {
	ID                  string        `db:"id"`
	DeliveryDays        []DeliveryDay `db:"delivery_days"`
	CutoffOffsetMinutes uint16        `db:"cutoff_offset_minutes"`
	MaxAttemptsPerOrder uint8         `db:"max_attempts_per_order"`
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
	OrderLeadMinutes    uint16
	IsOpen              bool
}

type GetConfigOutput struct {
	MaxAttemptsPerOrder uint8          `db:"max_attempts_per_order" json:"maxAttemptsPerOrder"`
	DiscountRules       []DiscountRule `db:"discount_rules"         json:"discountRules"`
	CutoffOffsetMinutes uint16         `db:"cutoff_offset_minutes"  json:"cutoffOffsetMinutes"`
	OrderLeadMinutes    uint16         `db:"order_lead_minutes"     json:"orderLeadMinutes"`
	IsOpen              bool           `db:"is_open"                json:"isOpen"`
}
