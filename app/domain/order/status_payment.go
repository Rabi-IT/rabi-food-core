package order

import "slices"

// PaymentStatus tracks the financial lifecycle, independent from production/delivery.
type PaymentStatus string

const (
	PaymentPending    PaymentStatus = "pending"    // not initiated or waiting user action
	PaymentAuthorized PaymentStatus = "authorized" // auth ok, yet to capture
	PaymentPaid       PaymentStatus = "paid"       // captured/settled
	PaymentFailed     PaymentStatus = "failed"     // attempt failed
	PaymentCancelled  PaymentStatus = "cancelled"  // voided before capture
	PaymentRefunded   PaymentStatus = "refunded"   // refunded after capture
)

var paymentStatusPrerequisites = map[PaymentStatus][]PaymentStatus{
	PaymentPending:    {},
	PaymentAuthorized: {PaymentPending},
	PaymentPaid:       {PaymentAuthorized},
	PaymentFailed:     {PaymentPending},
	PaymentCancelled:  {PaymentAuthorized},
	PaymentRefunded:   {PaymentPaid},
}

func (p PaymentStatus) GetPrerequisites() []PaymentStatus {
	return paymentStatusPrerequisites[p]
}

func (p PaymentStatus) CanTransitionTo(target PaymentStatus) bool {
	prerequisites := target.GetPrerequisites()

	return slices.Contains(prerequisites, p)
}

func (p PaymentStatus) String() string {
	return string(p)
}
