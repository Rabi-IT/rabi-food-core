package order

import "slices"

// FulfillmentStatus tracks the preparation lifecycle inside the fulfillment.
type FulfillmentStatus string

var (
	FulfillmentPending      FulfillmentStatus = "pending"       // order created, not accepted yet
	FulfillmentConfirmed    FulfillmentStatus = "confirmed"     // accepted by the fulfillment
	FulfillmentInProduction FulfillmentStatus = "in_production" // being prepared
	FulfillmentCancelled    FulfillmentStatus = "cancelled"     // cancelled
	FulfillmentReady        FulfillmentStatus = "ready"         // finished; ready for pickup/dispatch
)

var fulfillmentStatusPrerequisites = map[FulfillmentStatus][]FulfillmentStatus{
	FulfillmentPending:      {},
	FulfillmentConfirmed:    {FulfillmentPending},
	FulfillmentInProduction: {FulfillmentConfirmed},
	FulfillmentCancelled:    {FulfillmentPending, FulfillmentConfirmed, FulfillmentInProduction},
	FulfillmentReady:        {FulfillmentInProduction},
}

func (f FulfillmentStatus) GetPrerequisites() []FulfillmentStatus {
	return fulfillmentStatusPrerequisites[f]
}

func (f FulfillmentStatus) CanTransitionTo(target FulfillmentStatus) bool {
	prerequisites := target.GetPrerequisites()

	return slices.Contains(prerequisites, f)
}

func (f FulfillmentStatus) String() string {
	return string(f)
}
