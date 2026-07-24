package gin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextCopyConcurrency(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())

	// Create a cancelable context to simulate HTTP server request lifecycle
	parentCtx, cancel := context.WithCancel(context.Background())
	parentCtx = context.WithValue(parentCtx, "key", "value")
	req, _ := http.NewRequestWithContext(parentCtx, "GET", "/", nil)
	c.Request = req

	// Copy the context for async execution
	copied := c.Copy()

	// Cancel the original request context (simulating request completion)
	cancel()

	// Verify the copied context's request context is NOT cancelled
	select {
	case <-copied.Request.Context().Done():
		t.Error("Copied context was cancelled when the parent context was cancelled")
	default:
		// Success
	}

	// Verify values are still accessible in the copied context
	if copied.Request.Context().Value("key") != "value" {
		t.Error("Copied context lost values from the original request context")
	}
}
