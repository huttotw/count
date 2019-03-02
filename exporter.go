package count

import (
	"context"
	"time"
)

type Writer interface {
	// Write will be called for each key that has been counted in the last window.
	// The given time is the time that the tick happened, and can be different
	// than the current time. Your writer should handle context appropriately and
	// respond to cancellations, timeouts, and deadlines. If the writer fails
	// to persist the data, you should return an error, signaling that it is not
	// safe to reset the counter for that key.
	Write(ctx context.Context, key string, value uint64, current time.Time) error
}