package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anishnarang9/eth-indexer/internal/event"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}
type Sender struct {
	client *http.Client
}

func NewSender(client *http.Client) *Sender {
	return &Sender{client: client}
}
func (s *Sender) Send(ctx context.Context, e event.Event,
	destination string) (Result, error) {

	start := time.Now()
	if err := event.Validate(e); err != nil {
		return Result{}, fmt.Errorf("validate event: %w", err)
	}
	body, err := json.Marshal(e)

	if err != nil {
		return Result{}, fmt.Errorf("encode event: %w", err)
	}
	// Creats  a new POST request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		destination,
		bytes.NewReader(body),
	)
	if err != nil {
		return Result{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	result := Result{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
	return result, nil
}
