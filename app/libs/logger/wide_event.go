package logger

import (
	"context"
	"errors"
	"sync"

	"github.com/Rabi-IT/rabi-food-core/libs/errs"
	"github.com/rs/zerolog"
)

const wideEventKey = "wideEventKey"

// WideEvent is a struct that captures all relevant information about an HTTP request and its context for logging purposes.
// It is designed to be used in the Logging middleware to emit structured logs with rich contextual data.
type WideEvent struct {
	mu sync.Mutex

	// Event is a custom field to describe the specific event or action being logged
	Event string

	Method     string
	Path       string
	Query      string
	StatusCode int
	LatencyMs  int64
	RequestID  string

	ActorID       string
	ActorTenantID string
	IsBackoffice  bool

	OrderID        string
	CategoryID     string
	ProductID      string
	SubscriptionID string
	UserID         string
	TenantID       string

	// Payment
	ExternalPaymentID         string
	OriginalExternalPaymentID string
	PaymentIdempotencyHit     *bool

	err       string
	errorCode string
}

func (w *WideEvent) SetErrorCode(code string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.errorCode = code
}

func (w *WideEvent) FinishRequest(statusCode int, err error, latencyMs int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.StatusCode = statusCode
	w.LatencyMs = latencyMs
	if err != nil {
		w.err = err.Error()
		var appError *errs.AppError
		if errors.As(err, &appError) {
			w.errorCode = appError.Code
		}
	}
}

// reset clears all fields of the WideEvent, preparing it for reuse from the pool.
func (w *WideEvent) reset() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Method = ""
	w.StatusCode = 0
	w.LatencyMs = 0
	w.RequestID = ""

	w.ActorID = ""
	w.ActorTenantID = ""
	w.IsBackoffice = false

	w.OrderID = ""
	w.CategoryID = ""
	w.ProductID = ""
	w.SubscriptionID = ""
	w.UserID = ""
	w.TenantID = ""

	w.ExternalPaymentID = ""
	w.OriginalExternalPaymentID = ""
	w.PaymentIdempotencyHit = nil

	w.err = ""
	w.errorCode = ""
}

func (w *WideEvent) MarshalZerologObject(e *zerolog.Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	e.
		Str("method", w.Method).
		Str("path", w.Path).
		Str("query", w.Query).
		Int("status", w.StatusCode).
		Int64("latency_ms", w.LatencyMs)

	if w.Event != "" {
		e.Str("event", w.Event)
	}

	if w.OrderID != "" {
		e.Str("order_id", w.OrderID)
	}
	if w.CategoryID != "" {
		e.Str("category_id", w.CategoryID)
	}
	if w.ProductID != "" {
		e.Str("product_id", w.ProductID)
	}
	if w.SubscriptionID != "" {
		e.Str("subscription_id", w.SubscriptionID)
	}
	if w.UserID != "" {
		e.Str("user_id", w.UserID)
	}
	if w.TenantID != "" {
		e.Str("tenant_id", w.TenantID)
	}

	if w.ExternalPaymentID != "" {
		e.Str("external_payment_id", w.ExternalPaymentID)
	}
	if w.OriginalExternalPaymentID != "" {
		e.Str("original_external_payment_id", w.OriginalExternalPaymentID)
	}
	if w.PaymentIdempotencyHit != nil {
		e.Bool("payment_idempotency_hit", *w.PaymentIdempotencyHit)
	}

	if w.err != "" {
		e.Str("error", w.err)
	}
	if w.errorCode != "" {
		e.Str("error_code", w.errorCode)
	}
}

var wideEventPool = sync.Pool{
	New: func() any {
		return &WideEvent{}
	},
}

// Acquire retrieves a WideEvent from the pool, resetting its fields before returning it for use.
func Acquire() *WideEvent {
	wideEvent := wideEventPool.Get().(*WideEvent)
	return wideEvent
}

// Release returns a WideEvent to the pool for reuse
func Release(w *WideEvent) {
	w.reset()
	wideEventPool.Put(w)
}

func GetWideEvent(ctx context.Context) *WideEvent {
	if we, ok := ctx.Value(wideEventKey).(*WideEvent); ok && we != nil {
		return we
	}

	l := Get(ctx)
	l.Error().Msg("WideEvent not found in context")

	return &WideEvent{}
}

func WithWideEvent(ctx context.Context, we *WideEvent) context.Context {
	return context.WithValue(ctx, wideEventKey, we)
}
