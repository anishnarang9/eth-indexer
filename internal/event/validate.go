package event

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func Validate(e Event) error {
	if strings.TrimSpace(e.ID) == "" {
		return errors.New("missing event ID")
	}
	if e.Type != "order.created" {
		return errors.New("unsupported event type")
	}
	if _, err := time.Parse(time.RFC3339, e.OccurredAt); err != nil {
		return errors.New("invalid timestamp")
	}
	var payload map[string]json.RawMessage

	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return errors.New("payload must be a JSON object")
	}
	if payload == nil {
		return errors.New("payload must be a JSON object")
	}

	return nil
}
