package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-service/internal/config"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrTokenExpired = errors.New("access token expired")
)

type Client struct {
	cfg        *config.Config
	logger     *zap.Logger
	httpClient *http.Client
	mu         sync.RWMutex
}

func New(cfg *config.Config, lg *zap.Logger) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return &Client{
		cfg:    cfg,
		logger: lg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Jar:     jar,
		},
	}
}

func (c *Client) Login(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	payload := map[string]string{
		"username": c.cfg.NodeAuthEmail,
		"password": c.cfg.NodeAuthPassword,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/auth/login", c.cfg.NodeAuthURL),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed: status=%d", resp.StatusCode)
	}

	c.logger.Info("authenticated against Node API")
	return nil
}

func (c *Client) AddAuthHeaders(req *http.Request) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	u, _ := url.Parse(c.cfg.Host)

	cookies := c.httpClient.Jar.Cookies(u)

	var accessToken, csrfToken, refreshToken string
	for _, ck := range cookies {
		switch ck.Name {
		case "accessToken":
			accessToken = ck.Value
		case "csrfToken":
			csrfToken = ck.Value
		case "refreshToken":
			refreshToken = ck.Value
		}
	}

	if accessToken == "" || csrfToken == "" || refreshToken == "" {
		c.logger.Warn("auth headers missing", zap.Any("cookies", cookies))
		return ErrTokenExpired
	}

	cookieHeader := fmt.Sprintf(
		"accessToken=%s; csrfToken=%s; refreshToken=%s",
		accessToken, csrfToken, refreshToken,
	)

	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-CSRF-Token", csrfToken)

	return nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := c.AddAuthHeaders(req); err != nil {
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
		if err := c.AddAuthHeaders(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest && containsTokenExpired(resp) {
		_ = resp.Body.Close()
		if err := c.Login(ctx); err != nil {
			return nil, err
		}
		req2 := req.Clone(ctx)
		if err := c.AddAuthHeaders(req2); err != nil {
			return nil, err
		}
		return c.httpClient.Do(req2)
	}
	return resp, nil
}

func containsTokenExpired(resp *http.Response) bool {
	var msg struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&msg)
	return msg.Error == "Token expired"
}
