package subscription

type DeliveryStatus string

const (
	DeliveryScheduled DeliveryStatus = "scheduled"
	DeliverySkipped   DeliveryStatus = "skipped"
	DeliveryLocked    DeliveryStatus = "locked"
	DeliveryDelivered DeliveryStatus = "delivered"
	DeliveryFailed    DeliveryStatus = "failed_delivery"
)
