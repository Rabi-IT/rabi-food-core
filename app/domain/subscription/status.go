package subscription

import "slices"

// Status tracks the lifecycle of a subscription.
type Status string

const (
	StatusActive     Status = "active"
	StatusPaused     Status = "paused"
	StatusCancelled  Status = "cancelled"
	StatusExpired    Status = "expired"
	StatusSuperseded Status = "superseded"
)

var statusPrerequisites = map[Status][]Status{
	StatusActive:     {StatusPaused}, // resume
	StatusPaused:     {StatusActive}, // pause
	StatusCancelled:  {StatusActive, StatusPaused},
	StatusExpired:    {StatusActive, StatusPaused},
	StatusSuperseded: {StatusActive, StatusPaused},
}

func (s Status) GetPrerequisites() []Status {
	return statusPrerequisites[s]
}

func (s Status) CanTransitionTo(target Status) bool {
	return slices.Contains(target.GetPrerequisites(), s)
}

func (s Status) String() string {
	return string(s)
}
