package whooktown

import (
	"net/http"
)

// Client is the main whooktown SDK client
type Client struct {
	config     Config
	httpClient *http.Client

	// Service clients
	Auth       *AuthClient
	Sensors    *SensorsClient
	UI         *UIClient
	Camera     *CameraClient
	Traffic    *TrafficClient
	Popup      *PopupClient
	Groups     *GroupsClient
	Workflow   *WorkflowClient
	Backoffice *BackofficeClient
}

// New creates a new whooktown client with the given options
func New(opts ...Option) (*Client, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	c := &Client{
		config:     cfg,
		httpClient: httpClient,
	}

	// Create HTTP clients for each service
	authHTTP := newHTTPClient(httpClient, cfg.AuthURL)
	authHTTP.SetToken(cfg.Token)

	sensorHTTP := newHTTPClient(httpClient, cfg.SensorURL)
	sensorHTTP.SetToken(cfg.Token)

	uiHTTP := newHTTPClient(httpClient, cfg.UIURL)
	uiHTTP.SetToken(cfg.Token)

	workflowHTTP := newHTTPClient(httpClient, cfg.WorkflowURL)
	workflowHTTP.SetToken(cfg.Token)

	backofficeHTTP := newHTTPClient(httpClient, cfg.BackofficeURL)
	backofficeHTTP.SetAdminToken(cfg.AdminSecret)

	// Initialize service clients
	c.Auth = &AuthClient{http: authHTTP}
	c.Sensors = &SensorsClient{http: sensorHTTP}
	c.UI = &UIClient{http: uiHTTP}
	c.Camera = &CameraClient{http: uiHTTP}
	c.Traffic = &TrafficClient{http: uiHTTP}
	c.Popup = &PopupClient{http: uiHTTP}
	c.Groups = &GroupsClient{http: uiHTTP}
	c.Workflow = &WorkflowClient{http: workflowHTTP}
	c.Backoffice = &BackofficeClient{http: backofficeHTTP}

	return c, nil
}

// SetToken updates the authentication token for all service clients
func (c *Client) SetToken(token string) {
	c.config.Token = token
	c.Auth.http.SetToken(token)
	c.Sensors.http.SetToken(token)
	c.UI.http.SetToken(token)
	c.Camera.http.SetToken(token)
	c.Traffic.http.SetToken(token)
	c.Popup.http.SetToken(token)
	c.Groups.http.SetToken(token)
	c.Workflow.http.SetToken(token)
}

// SetAdminSecret updates the admin secret for the backoffice client
func (c *Client) SetAdminSecret(secret string) {
	c.config.AdminSecret = secret
	c.Backoffice.http.SetAdminToken(secret)
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() Config {
	return c.config
}
