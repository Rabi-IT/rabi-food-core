package logger

const (
	// Core context keys.

	// RequestID represents the request identifier context key.
	RequestID = "request_id"
	// IsBackoffice indicates if the request is from backoffice context key.
	IsBackoffice = "is_backoffice"

	// Domain context keys.

	AppPlanID  = "app_plan_id"
	OrderID    = "order_id"
	UserID     = "user_id"
	TenantID   = "tenant_id"
	PaymentID  = "payment_id"
	ProductID  = "product_id"
	CategoryID = "category_id"

	// HTTP context keys.

	// Path represents the HTTP path context key.
	Path = "path"
	// Method represents the HTTP method context key.
	Method = "method"
	// Query represents the HTTP query parameters context key.
	Query = "query"
)
