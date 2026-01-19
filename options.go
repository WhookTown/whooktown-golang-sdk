package whooktown

import (
	"net/http"
	"time"
)

// Config holds the client configuration
type Config struct {
	// Base URLs for services
	AuthURL       string
	SensorURL     string
	UIURL         string
	WorkflowURL   string
	BackofficeURL string
	SSEURL        string

	// Authentication
	Token       string // Bearer token for user authentication
	AdminSecret string // For backoffice API (X-Admin-Token header)

	// HTTP settings
	Timeout    time.Duration
	MaxRetries int
	RetryWait  time.Duration
	HTTPClient *http.Client

	// Debug
	Debug bool
}

// Option configures the client
type Option func(*Config)

// defaultConfig returns the default configuration
func defaultConfig() Config {
	return Config{
		AuthURL:       "http://localhost:8981",
		SensorURL:     "http://localhost:8081",
		UIURL:         "http://localhost:8083",
		WorkflowURL:   "http://localhost:8084",
		BackofficeURL: "http://localhost:8086",
		SSEURL:        "http://localhost:8082",
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryWait:     time.Second,
	}
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	// Configuration is valid by default, services will fail at request time if URLs are wrong
	return nil
}

// WithBaseURL sets all service URLs from a single base URL
// Example: WithBaseURL("http://api.whooktown.com") sets all services to that host
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.AuthURL = baseURL + ":8981"
		c.SensorURL = baseURL + ":8081"
		c.UIURL = baseURL + ":8083"
		c.WorkflowURL = baseURL + ":8084"
		c.BackofficeURL = baseURL + ":8086"
		c.SSEURL = baseURL + ":8082"
	}
}

// WithToken sets the Bearer authentication token
func WithToken(token string) Option {
	return func(c *Config) {
		c.Token = token
	}
}

// WithAdminSecret sets the admin secret for backoffice API (X-Admin-Token header)
func WithAdminSecret(secret string) Option {
	return func(c *Config) {
		c.AdminSecret = secret
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithRetry configures retry behavior for failed requests
func WithRetry(maxRetries int, retryWait time.Duration) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
		c.RetryWait = retryWait
	}
}

// WithDebug enables debug logging
func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.Debug = debug
	}
}

// WithAuthURL sets the auth service URL
func WithAuthURL(url string) Option {
	return func(c *Config) {
		c.AuthURL = url
	}
}

// WithSensorURL sets the sensor endpoint URL
func WithSensorURL(url string) Option {
	return func(c *Config) {
		c.SensorURL = url
	}
}

// WithUIURL sets the UI endpoint URL
func WithUIURL(url string) Option {
	return func(c *Config) {
		c.UIURL = url
	}
}

// WithWorkflowURL sets the workflow engine URL
func WithWorkflowURL(url string) Option {
	return func(c *Config) {
		c.WorkflowURL = url
	}
}

// WithBackofficeURL sets the backoffice API URL
func WithBackofficeURL(url string) Option {
	return func(c *Config) {
		c.BackofficeURL = url
	}
}

// WithSSEURL sets the SSE endpoint URL (for WebSocket)
func WithSSEURL(url string) Option {
	return func(c *Config) {
		c.SSEURL = url
	}
}

// WithServices configures individual service URLs
func WithServices(auth, sensor, ui, workflow, backoffice, sse string) Option {
	return func(c *Config) {
		if auth != "" {
			c.AuthURL = auth
		}
		if sensor != "" {
			c.SensorURL = sensor
		}
		if ui != "" {
			c.UIURL = ui
		}
		if workflow != "" {
			c.WorkflowURL = workflow
		}
		if backoffice != "" {
			c.BackofficeURL = backoffice
		}
		if sse != "" {
			c.SSEURL = sse
		}
	}
}
