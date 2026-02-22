package order

import "slices"

// DeliveryStatus tracks the last-mile logistics lifecycle.
// If the order is pickup-at-store only, it can remain "pending" and jump to "delivered".
type DeliveryStatus string

var (
	DeliveryPending   DeliveryStatus = "pending"          // awaiting fulfillment=ready
	DeliveryAssigned  DeliveryStatus = "assigned"         // driver assigned / courier matched
	DeliveryPickedUp  DeliveryStatus = "picked_up"        // courier picked up at store
	DeliveryOutFor    DeliveryStatus = "out_for_delivery" // on the way
	DeliveryDelivered DeliveryStatus = "delivered"        // customer received
	DeliveryCancelled DeliveryStatus = "cancelled"        // cancelled before dispatch
	DeliveryRefunded  DeliveryStatus = "refunded"         // refunded after delivery/charge
)

var deliveryStatusPrerequisites = map[DeliveryStatus][]DeliveryStatus{
	DeliveryPending:   {},
	DeliveryAssigned:  {DeliveryPending},
	DeliveryPickedUp:  {DeliveryAssigned},
	DeliveryOutFor:    {DeliveryPickedUp},
	DeliveryDelivered: {DeliveryOutFor},
	DeliveryCancelled: {DeliveryPending, DeliveryAssigned},
	DeliveryRefunded:  {DeliveryDelivered},
}

func (d DeliveryStatus) GetPrerequisites() []DeliveryStatus {
	return deliveryStatusPrerequisites[d]
}

func (d DeliveryStatus) CanTransitionTo(target DeliveryStatus) bool {
	prerequisites := target.GetPrerequisites()

	return slices.Contains(prerequisites, d)
}

func (d DeliveryStatus) String() string {
	return string(d)
}
