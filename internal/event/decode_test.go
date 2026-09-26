package event

import "testing"

func TestDecode(t *testing.T) {
	goodLine := []byte(`{"id":"evt_001","type":"order.created","occurred_at":"2026-09-23T14:00:00Z","payload":{"order_id":"order_123"}}`)
	badLine := []byte(`{"id":`)

	e, err := Decode(goodLine)

	if err != nil {
		t.Fatalf("valid JSON: unexpected error: %v", err)
	}
	if e.ID != "evt_001" {
		t.Errorf("ID = %q; want %q", e.ID, "evt_001")
	}

	if _, err := Decode(badLine); err == nil {
		t.Error("malformed JSON: expected an error")
	}
}
