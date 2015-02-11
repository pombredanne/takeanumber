package item_test

import (
	"testing"
	"github.com/toastdriven/takeanumber/item"
)

func TestItem(t *testing.T) {
	// Test initialization.
	i, err := item.New("test", 2)

	if err != nil {
		t.Error("Saw an error: ", err)
	}

	if i.Id == "" {
		t.Error("No Id automatically created")
	}

	if i.Body != "test" {
		t.Error("Body not set appropriately")
	}

	if i.InitialRetries != 2 {
		t.Error("Initial retries not 2, was: ", i.InitialRetries)
	}

	if i.RemainingRetries != 2 {
		t.Error("Remaining retries not 2, was: ", i.RemainingRetries)
	}

	if i.IsReserved() {
		t.Error("Should not be reserved by default")
	}

	// Test methods.
	i.Reserve()

	if !i.IsReserved() {
		t.Error("Reserving failed")
	}

	i.Release()

	if i.IsReserved() {
		t.Error("Releasing failed")
	}

	retryable := i.ShouldRetry()

	if !retryable {
		t.Error("Should still be retryable")
	}

	decremented := i.DecrRetries()

	if i.RemainingRetries == i.InitialRetries {
		t.Error("Retry counts still match")
	}

	if !decremented {
		t.Error("Decrementing retries failed")
	}

	retryable = i.ShouldRetry()

	if !retryable {
		t.Error("Should still be retryable")
	}

	decremented = i.DecrRetries()

	if !decremented {
		t.Error("Decrementing retries failed")
	}

	retryable = i.ShouldRetry()

	if retryable {
		t.Error("Should NOT still be retryable")
	}

	decremented = i.DecrRetries()

	if decremented {
		t.Error("Decrementing below zero should fail")
	}
}
