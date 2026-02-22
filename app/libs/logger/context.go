package logger

const (
	// Core
	RequestID    = "request_id"
	IsBackoffice = "is_backoffice"

	// Domain
	OrderID    = "order_id"
	UserID     = "user_id"
	TenantID   = "tenant_id"
	PaymentID  = "payment_id"
	ProductID  = "product_id"
	CategoryID = "category_id"

	// HTTP
	Path   = "path"
	Method = "method"
	Query  = "query"
)
