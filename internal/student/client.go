package student

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-service/internal/auth"
	"go-service/internal/config"
	"net/http"

	"go.uber.org/zap"
)

var ErrStudentNotFound = errors.New("student not found")

type Client struct {
	auth   *auth.Client
	cfg    *config.Config
	logger *zap.Logger
}

func NewClient(a *auth.Client, cfg *config.Config, lg *zap.Logger) *Client {
	return &Client{auth: a, cfg: cfg, logger: lg}
}

func (c *Client) Get(ctx context.Context, id int64) (*Student, error) {
	url := fmt.Sprintf("%s/students/%d", c.cfg.NodeAPIURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create student request: %w", err)
	}

	resp, err := c.auth.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("request to node api failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var s Student
		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, fmt.Errorf("failed to decode student response: %w", err)
		}
		return &s, nil
	case http.StatusNotFound:
		return nil, ErrStudentNotFound
	default:
		return nil, fmt.Errorf("unexpected status code from node api: %d", resp.StatusCode)
	}
}
