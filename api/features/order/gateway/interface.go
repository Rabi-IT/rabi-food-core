package gateway

import (
	"context"
	"time"

	"github.com/Rabi-IT/rabi-food-core/domain/payment_status"
	"github.com/Rabi-IT/rabi-food-core/features/order"
	"github.com/Rabi-IT/rabi-food-core/libs/database"
)

type OrderGateway interface {
	Create(ctx context.Context, input CreateInput) (string, error)
	GetByID(ctx context.Context, filter GetByIDFilter) (*GetByIDOutput, error)
	Patch(ctx context.Context, filter PatchFilter, values PatchValues) (bool, error)
	Paginate(ctx context.Context, filter PaginateFilter, paginate database.PaginateInput) (PaginateOutput, error)
	Delete(ctx context.Context, filter DeleteFilter) (bool, error)
}

type OrderItem struct {
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	Quantity    uint   `json:"quantity"`
	UnitPrice   uint   `json:"unitPrice"`
	Total       uint   `json:"total"`
}

type CreateInput struct {
	UserID            string
	TenantID          string
	Code              string
	PaymentStatus     payment_status.Status
	FulfillmentStatus order.FulfillmentStatus
	DeliveryStatus    order.DeliveryStatus
	Notes             string
	TotalPrice        uint
	Items             []OrderItem
}

type GetByIDFilter struct {
	ID       string `json:"id"`
	UserID   string `json:"-"`
	TenantID string `json:"-"`
}

type GetByIDOutput struct {
	ID                string                  `db:"id"                  json:"id"`
	TenantID          string                  `db:"tenant_id"           json:"tenantId"`
	Code              string                  `db:"code"                json:"code"`
	PaymentStatus     payment_status.Status   `db:"payment_status"      json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `db:"fulfillment_status"  json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `db:"delivery_status"     json:"deliveryStatus"`
	Notes             string                  `db:"notes"               json:"notes"`
	TotalPrice        uint                    `db:"total_price"         json:"totalPrice"`
	Items             []OrderItem             `db:"items"               json:"items"`
	PaidAt            *time.Time              `db:"paid_at"             json:"paidAt"`
	CreatedAt         time.Time               `db:"created_at"          json:"createdAt"`
	ExternalPaymentID *string                 `db:"external_payment_id" json:"externalPaymentId"`
}

type PatchFilter struct {
	ID            string
	TenantID      string
	PaymentStatus payment_status.Status
}

type PatchValues struct {
	PaidAt            *time.Time               `json:"paidAt"`
	Provider          *string                  `json:"provider"`
	PaymentStatus     *payment_status.Status   `json:"paymentStatus"`
	ExternalPaymentID *string                  `json:"externalPaymentId"`
	FulfillmentStatus *order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    *order.DeliveryStatus    `json:"deliveryStatus"`
}

type PaginateFilter struct {
	UserID            string                  `json:"userId"`
	TenantID          string                  `json:"tenantId"`
	PaymentStatus     payment_status.Status   `json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `json:"deliveryStatus"`
	CreatedAtFrom     time.Time               `json:"createdAtFrom"`
	CreatedAtTo       time.Time               `json:"createdAtTo"`
}

type PaginateData struct {
	ID                string                  `db:"id"                 json:"id"`
	TenantID          string                  `db:"tenant_id"          json:"tenantId"`
	Code              string                  `db:"code"               json:"code"`
	PaymentStatus     payment_status.Status   `db:"payment_status"     json:"paymentStatus"`
	FulfillmentStatus order.FulfillmentStatus `db:"fulfillment_status" json:"fulfillmentStatus"`
	DeliveryStatus    order.DeliveryStatus    `db:"delivery_status"    json:"deliveryStatus"`
	Notes             string                  `db:"notes"              json:"notes"`
	TotalPrice        uint                    `db:"total_price"        json:"totalPrice"`
	CreatedAt         time.Time               `db:"created_at"         json:"createdAt"`
}

type PaginateOutput struct {
	Data     []PaginateData `json:"data"`
	MaxPages int            `json:"maxPages"`
}

type DeleteFilter struct {
	ID       string `json:"id"`
	UserID   string `json:"-"`
	TenantID string `json:"-"`
}
