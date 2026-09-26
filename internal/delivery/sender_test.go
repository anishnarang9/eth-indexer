package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anishnarang9/eth-indexer/internal/event"
)

func testEvent() event.Event {
	return event.Event{
		ID:         "evt_001",
		Type:       "order.created",
		OccurredAt: "2026-09-23T14:00:00Z",
		Payload:    json.RawMessage(`{"order_id":"order_123"}`),
	}
}

func TestSender(t *testing.T) {

	valid := testEvent()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s; want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Header = %q; want application/json", got)
		}
		var received event.Event
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if received.ID != valid.ID {
			t.Errorf("ID = %q; want %q", received.ID, valid.ID)
		}
		if received.Type != valid.Type {
			t.Errorf("Type = %q; want %q", received.Type, valid.Type)
		}
		if received.OccurredAt != valid.OccurredAt {
			t.Errorf("OccurredAt = %q; want %q", received.OccurredAt, valid.OccurredAt)
		}
		if string(received.Payload) != string(valid.Payload) {
			t.Errorf("Payload = %s; want %s", received.Payload, valid.Payload)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	sender := NewSender(client)
	result, err := sender.Send(context.Background(), valid, server.URL)

	if err != nil {
		t.Fatalf("Send returned an error: %v", err)
	}
	if result.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d; want 204", result.StatusCode)
	}

}
func TestHTTPFailurePreservesStatus(t *testing.T) {

	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {

		t.Run(http.StatusText(status), func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			sender := NewSender(&http.Client{Timeout: time.Second})
			result, err := sender.Send(context.Background(), testEvent(), server.URL)

			if err != nil {
				t.Fatalf("Send: %v", err)
			}
			if result.StatusCode != status {
				t.Errorf("status = %d; want %d", result.StatusCode, status)
			}
		})
	}
}

func TestInvalidEventDoesNotSend(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	invalid := testEvent()
	invalid.ID = ""

	sender := NewSender(&http.Client{Timeout: time.Second})
	if _, err := sender.Send(context.Background(), invalid, server.URL); err == nil {
		t.Fatal("expected a validation error")
	}
	if got := requests.Load(); got != 0 {
		t.Errorf("receiver got %d requests; want 0", got)
	}
}
