package logger

const (
	// Core context keys.

	// Event represents the event name context key.
	EventName = "event"

	ActorID = "actor_id"

	// RequestID represents the request identifier context key.
	RequestID = "request_id"
	// IsBackoffice indicates if the request is from backoffice context key.
	IsBackoffice = "is_backoffice"

	// Domain context keys.

	// OrderID represents the order identifier context key.
	OrderID = "order_id"
	// UserID represents the user identifier context key.
	UserID = "user_id"
	// TenantID represents the tenant identifier context key.
	TenantID = "tenant_id"
	// PaymentID represents the payment identifier context key.
	PaymentID = "payment_id"
	// ProductID represents the product identifier context key.
	ProductID = "product_id"
	// CategoryID represents the category identifier context key.
	CategoryID = "category_id"
	// SubscriptionID represents the subscription identifier context key.
	SubscriptionID = "subscription_id"
	// CurrentStatus represents the current status context key.
	CurrentStatus = "current_status"
	// TargetStatus represents the target status context key.
	TargetStatus = "target_status"

	// HTTP context keys.

	// Path represents the HTTP path context key.
	Path = "path"
	// Method represents the HTTP method context key.
	Method = "method"
	// Query represents the HTTP query parameters context key.
	Query = "query"
)
