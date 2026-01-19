package whooktown

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// httpClient wraps http.Client with common functionality
type httpClient struct {
	client     *http.Client
	baseURL    string
	token      string
	adminToken string
	debug      bool
	maxRetries int
	retryWait  time.Duration
}

// newHTTPClient creates a new HTTP client wrapper
func newHTTPClient(client *http.Client, baseURL string) *httpClient {
	return &httpClient{
		client:     client,
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		maxRetries: 3,
		retryWait:  time.Second,
	}
}

// SetToken sets the Bearer token for authentication
func (c *httpClient) SetToken(token string) {
	c.token = token
}

// SetAdminToken sets the X-Admin-Token header value
func (c *httpClient) SetAdminToken(token string) {
	c.adminToken = token
}

// Get performs a GET request
func (c *httpClient) Get(ctx context.Context, path string, result interface{}) error {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (c *httpClient) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request
func (c *httpClient) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPut, path, body, result)
}

// Patch performs a PATCH request
func (c *httpClient) Patch(ctx context.Context, path string, body, result interface{}) error {
	return c.doRequest(ctx, http.MethodPatch, path, body, result)
}

// Delete performs a DELETE request
func (c *httpClient) Delete(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil)
}

// doRequest performs an HTTP request with retry logic
func (c *httpClient) doRequest(ctx context.Context, method, path string, body, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return &Error{
					Code:    ErrTimeout,
					Message: "request cancelled",
					Cause:   ctx.Err(),
				}
			case <-time.After(c.retryWait * time.Duration(attempt)):
			}
		}

		err := c.executeRequest(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		// Don't retry on client errors (4xx)
		if e, ok := err.(*Error); ok {
			if e.StatusCode >= 400 && e.StatusCode < 500 {
				return err
			}
		}

		lastErr = err
	}

	return lastErr
}

// executeRequest performs a single HTTP request
func (c *httpClient) executeRequest(ctx context.Context, method, path string, body, result interface{}) error {
	// Build URL
	reqURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return &Error{
			Code:    ErrValidation,
			Message: fmt.Sprintf("invalid path: %s", path),
			Cause:   err,
		}
	}

	// Prepare body
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return &Error{
				Code:    ErrValidation,
				Message: "failed to marshal request body",
				Cause:   err,
			}
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return &Error{
			Code:    ErrNetworkError,
			Message: "failed to create request",
			Cause:   err,
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if c.adminToken != "" {
		req.Header.Set("X-Admin-Token", c.adminToken)
	}

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return &Error{
			Code:    ErrNetworkError,
			Message: "request failed",
			Cause:   err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{
			Code:    ErrNetworkError,
			Message: "failed to read response body",
			Cause:   err,
		}
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		return parseHTTPError(resp.StatusCode, respBody)
	}

	// Parse response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return &Error{
				Code:    ErrInternalServer,
				Message: "failed to parse response",
				Cause:   err,
			}
		}
	}

	return nil
}

// parseHTTPError converts HTTP response to SDK error
func parseHTTPError(statusCode int, body []byte) error {
	e := &Error{
		StatusCode: statusCode,
	}

	// Try to parse JSON error response
	var errResp struct {
		Message string                 `json:"message"`
		Error   string                 `json:"error"`
		Code    string                 `json:"code"`
		Details map[string]interface{} `json:"details"`
	}

	if json.Unmarshal(body, &errResp) == nil {
		if errResp.Message != "" {
			e.Message = errResp.Message
		} else if errResp.Error != "" {
			e.Message = errResp.Error
		}

		// Check for quota exceeded error
		if errResp.Code == "QUOTA_EXCEEDED" || errResp.Code == "ASSET_QUOTA_EXCEEDED" || errResp.Code == "LAYOUT_QUOTA_EXCEEDED" {
			qe := &QuotaError{
				Code:       ErrQuotaExceeded,
				Message:    e.Message,
				StatusCode: statusCode,
			}
			if errResp.Details != nil {
				if plan, ok := errResp.Details["plan"].(string); ok {
					qe.Plan = plan
				}
				if current, ok := errResp.Details["current"].(float64); ok {
					qe.Current = int(current)
				}
				if limit, ok := errResp.Details["limit"].(float64); ok {
					qe.Limit = int(limit)
				}
				if typ, ok := errResp.Details["type"].(string); ok {
					qe.QuotaType = typ
				}
			}
			return qe
		}

		e.Details = errResp.Details
	} else if len(body) > 0 {
		e.Message = string(body)
	}

	// Set error code based on status code
	switch statusCode {
	case http.StatusUnauthorized:
		e.Code = ErrUnauthorized
		if e.Message == "" {
			e.Message = "unauthorized"
		}
	case http.StatusForbidden:
		e.Code = ErrForbidden
		if e.Message == "" {
			e.Message = "forbidden"
		}
	case http.StatusNotFound:
		e.Code = ErrNotFound
		if e.Message == "" {
			e.Message = "not found"
		}
	case http.StatusBadRequest:
		e.Code = ErrBadRequest
		if e.Message == "" {
			e.Message = "bad request"
		}
	case http.StatusPaymentRequired:
		e.Code = ErrQuotaExceeded
		if e.Message == "" {
			e.Message = "quota exceeded"
		}
	default:
		e.Code = ErrInternalServer
		if e.Message == "" {
			e.Message = fmt.Sprintf("server error: %d", statusCode)
		}
	}

	return e
}
