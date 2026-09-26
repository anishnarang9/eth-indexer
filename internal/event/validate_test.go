package event

import (
	"encoding/json"
	"testing"
)

func TestValidity(t *testing.T) {

	valid := Event{
		ID:         "evt_001",
		Type:       "order.created",
		OccurredAt: "2026-09-23T14:00:00Z",
		Payload:    json.RawMessage(`{"order_id":"order_123"}`),
	}

	missingID := Event{
		ID:         "",
		Type:       "order.created",
		OccurredAt: "2026-09-23T14:00:00Z",
		Payload:    json.RawMessage(`{"order_id":"order_123"}`),
	}

	if err := Validate(valid); err != nil {
		t.Errorf("valid event: unexpected error: %v", err)
	}

	if err := Validate(missingID); err == nil {
		t.Error("event with missing ID: expected an error")
	}
}
