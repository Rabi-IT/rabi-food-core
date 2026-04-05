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

	orderID        string
	categoryID     string
	productID      string
	subscriptionID string
	userID         string
	tenantID       string

	// Payment
	externalPaymentID         string
	originalExternalPaymentID string
	paymentIdempotencyHit     *bool

	err       string
	errorCode string
}

// --- Usecase setters (thread-safe) ---

func (w *WideEvent) SetOrderID(id string) {
	w.mu.Lock()
	w.orderID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetCategoryID(id string) {
	w.mu.Lock()
	w.categoryID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetProductID(id string) {
	w.mu.Lock()
	w.productID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetSubscriptionID(id string) {
	w.mu.Lock()
	w.subscriptionID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetUserID(id string) {
	w.mu.Lock()
	w.userID = id
	w.mu.Unlock()
}

func (w *WideEvent) SetTenantID(id string) {
	w.mu.Lock()
	w.tenantID = id
	w.mu.Unlock()
}

func (w *WideEvent) TrackOrderPayment(orderID, externalPaymentID string) {
	w.mu.Lock()
	w.orderID = orderID
	w.externalPaymentID = externalPaymentID
	w.mu.Unlock()
}

func (w *WideEvent) MarkPaymentIdempotent() {
	w.mu.Lock()
	v := true
	w.paymentIdempotencyHit = &v
	w.mu.Unlock()
}

func (w *WideEvent) TrackPaymentConflict(originalExternalPaymentID string) {
	w.mu.Lock()
	w.originalExternalPaymentID = originalExternalPaymentID
	w.mu.Unlock()
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

	w.Event = ""
	w.Method = ""
	w.Path = ""
	w.Query = ""
	w.StatusCode = 0
	w.LatencyMs = 0
	w.RequestID = ""

	w.ActorID = ""
	w.ActorTenantID = ""
	w.IsBackoffice = false

	w.orderID = ""
	w.categoryID = ""
	w.productID = ""
	w.subscriptionID = ""
	w.userID = ""
	w.tenantID = ""

	w.externalPaymentID = ""
	w.originalExternalPaymentID = ""
	w.paymentIdempotencyHit = nil

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

	if w.orderID != "" {
		e.Str("order_id", w.orderID)
	}
	if w.categoryID != "" {
		e.Str("category_id", w.categoryID)
	}
	if w.productID != "" {
		e.Str("product_id", w.productID)
	}
	if w.subscriptionID != "" {
		e.Str("subscription_id", w.subscriptionID)
	}
	if w.userID != "" {
		e.Str("user_id", w.userID)
	}
	if w.tenantID != "" {
		e.Str("tenant_id", w.tenantID)
	}

	if w.externalPaymentID != "" {
		e.Str("payment.external_id", w.externalPaymentID)
	}
	if w.originalExternalPaymentID != "" {
		e.Str("payment.original_external_id", w.originalExternalPaymentID)
	}
	if w.paymentIdempotencyHit != nil {
		e.Bool("payment.idempotency_hit", *w.paymentIdempotencyHit)
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
