package whooktown

import (
	"net/http"
	"os"
	"time"
)

// Environment represents the deployment environment
type Environment string

const (
	// EnvProduction is the production environment (default)
	EnvProduction Environment = "PROD"
	// EnvDevelopment is the development environment
	EnvDevelopment Environment = "DEV"
)

// Production URLs (default)
const (
	ProdAuthURL         = "https://auth.whook.town"
	ProdSensorURL       = "https://sensors.whook.town"
	ProdUIURL           = "https://api.whook.town"
	ProdWorkflowURL     = "https://api.whook.town"
	ProdBackofficeURL   = "https://admin.whook.town"
	ProdSSEURL          = "https://ws.whook.town"
	ProdSubscriptionURL = "https://subscription.whook.town"
	ProdAudioStreamURL  = "https://stream.whook.town"
)

// Development URLs
const (
	DevAuthURL         = "https://auth.dev.whook.town"
	DevSensorURL       = "https://sensors.dev.whook.town"
	DevUIURL           = "https://api.dev.whook.town"
	DevWorkflowURL     = "https://api.dev.whook.town"
	DevBackofficeURL   = "https://admin.dev.whook.town"
	DevSSEURL          = "https://ws.dev.whook.town"
	DevSubscriptionURL = "https://subscription.dev.whook.town"
	DevAudioStreamURL  = "https://stream.dev.whook.town"
)

// Config holds the client configuration
type Config struct {
	// Base URLs for services
	AuthURL         string
	SensorURL       string
	UIURL           string
	WorkflowURL     string
	BackofficeURL   string
	SSEURL          string
	SubscriptionURL string
	AudioStreamURL  string

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

// getEnvironmentFromEnv returns the environment based on WHOOKTOWN_ENV variable
func getEnvironmentFromEnv() Environment {
	env := os.Getenv("WHOOKTOWN_ENV")
	if env == "DEV" {
		return EnvDevelopment
	}
	return EnvProduction
}

// defaultConfig returns the default configuration based on WHOOKTOWN_ENV
func defaultConfig() Config {
	env := getEnvironmentFromEnv()
	return configForEnvironment(env)
}

// configForEnvironment returns the configuration for a specific environment
func configForEnvironment(env Environment) Config {
	cfg := Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryWait:  time.Second,
	}

	if env == EnvDevelopment {
		cfg.AuthURL = DevAuthURL
		cfg.SensorURL = DevSensorURL
		cfg.UIURL = DevUIURL
		cfg.WorkflowURL = DevWorkflowURL
		cfg.BackofficeURL = DevBackofficeURL
		cfg.SSEURL = DevSSEURL
		cfg.SubscriptionURL = DevSubscriptionURL
		cfg.AudioStreamURL = DevAudioStreamURL
	} else {
		cfg.AuthURL = ProdAuthURL
		cfg.SensorURL = ProdSensorURL
		cfg.UIURL = ProdUIURL
		cfg.WorkflowURL = ProdWorkflowURL
		cfg.BackofficeURL = ProdBackofficeURL
		cfg.SSEURL = ProdSSEURL
		cfg.SubscriptionURL = ProdSubscriptionURL
		cfg.AudioStreamURL = ProdAudioStreamURL
	}

	return cfg
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	// Configuration is valid by default, services will fail at request time if URLs are wrong
	return nil
}

// WithEnvironment sets all URLs based on the environment (PROD or DEV)
func WithEnvironment(env Environment) Option {
	return func(c *Config) {
		envCfg := configForEnvironment(env)
		c.AuthURL = envCfg.AuthURL
		c.SensorURL = envCfg.SensorURL
		c.UIURL = envCfg.UIURL
		c.WorkflowURL = envCfg.WorkflowURL
		c.BackofficeURL = envCfg.BackofficeURL
		c.SSEURL = envCfg.SSEURL
		c.SubscriptionURL = envCfg.SubscriptionURL
		c.AudioStreamURL = envCfg.AudioStreamURL
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

// WithSubscriptionURL sets the subscription service URL
func WithSubscriptionURL(url string) Option {
	return func(c *Config) {
		c.SubscriptionURL = url
	}
}

// WithAudioStreamURL sets the audio streaming service URL
func WithAudioStreamURL(url string) Option {
	return func(c *Config) {
		c.AudioStreamURL = url
	}
}

// WithBaseURL sets all service URLs from a single base URL (for custom deployments)
// Note: This is for custom deployments only. Use WithEnvironment for standard PROD/DEV.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.AuthURL = baseURL
		c.SensorURL = baseURL
		c.UIURL = baseURL
		c.WorkflowURL = baseURL
		c.BackofficeURL = baseURL
		c.SSEURL = baseURL
		c.SubscriptionURL = baseURL
		c.AudioStreamURL = baseURL
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
