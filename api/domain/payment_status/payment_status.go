package payment_status

import "slices"

// Status tracks the financial lifecycle, independent from production/delivery.
type Status string

const (
	StatusPending    Status = "pending"    // not initiated or waiting user action
	StatusAuthorized Status = "authorized" // auth ok, yet to capture
	StatusPaid       Status = "paid"       // captured/settled
	StatusFailed     Status = "failed"     // attempt failed
	StatusCancelled  Status = "cancelled"  // voided before capture
	StatusRefunded   Status = "refunded"   // refunded after capture
)

var paymentStatusPrerequisites = map[Status][]Status{
	StatusPending:    {},
	StatusAuthorized: {StatusPending},
	StatusPaid:       {StatusAuthorized},
	StatusFailed:     {StatusPending},
	StatusCancelled:  {StatusAuthorized},
	StatusRefunded:   {StatusPaid},
}

func (p Status) GetPrerequisites() []Status {
	return paymentStatusPrerequisites[p]
}

func (p Status) CanTransitionTo(target Status) bool {
	prerequisites := target.GetPrerequisites()

	return slices.Contains(prerequisites, p)
}

func (p Status) String() string {
	return string(p)
}
