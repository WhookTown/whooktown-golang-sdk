package whooktown

import (
	"context"
)

// AuthClient provides access to the authentication service
type AuthClient struct {
	http *httpClient
}

// SignupRequest represents a signup request
type SignupRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"` // admin, user, viewer, sensor
	Name  string `json:"name,omitempty"`
	AppID string `json:"app_id,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"`
	Name  string `json:"name,omitempty"`
	AppID string `json:"app_id,omitempty"`
}

// CreateTokenRequest represents a request to create a new token
type CreateTokenRequest struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

// SignupResponse represents the response from signup
type SignupResponse struct {
	AppToken string `json:"app_token"`
}

// LoginResponse represents the response from login
type LoginResponse struct {
	AppToken string `json:"app_token"`
}

// Signup creates a new account
func (c *AuthClient) Signup(ctx context.Context, req *SignupRequest) (*Token, error) {
	var resp SignupResponse
	if err := c.http.Post(ctx, "/auth/signup", req, &resp); err != nil {
		return nil, err
	}
	return &Token{Token: resp.AppToken}, nil
}

// Login logs into an existing account
func (c *AuthClient) Login(ctx context.Context, req *LoginRequest) (*Token, error) {
	var resp LoginResponse
	if err := c.http.Post(ctx, "/auth/login", req, &resp); err != nil {
		return nil, err
	}
	return &Token{Token: resp.AppToken}, nil
}

// Logout logs out and optionally revokes the current token
func (c *AuthClient) Logout(ctx context.Context, appID string) error {
	body := map[string]string{}
	if appID != "" {
		body["app_id"] = appID
	}
	return c.http.Post(ctx, "/auth/logout", body, nil)
}

// GetRoles returns available role types
func (c *AuthClient) GetRoles(ctx context.Context) (map[string]map[string]string, error) {
	var roles map[string]map[string]string
	if err := c.http.Get(ctx, "/auth/roles", &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// CheckToken validates a token and returns its details
func (c *AuthClient) CheckToken(ctx context.Context, token string) (*Token, error) {
	var t Token
	if err := c.http.Get(ctx, "/auth/check/"+token, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// ListTokens returns all tokens for the authenticated account
func (c *AuthClient) ListTokens(ctx context.Context) ([]Token, error) {
	var tokens []Token
	if err := c.http.Get(ctx, "/account/token", &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

// CreateToken creates a new token for the authenticated account
func (c *AuthClient) CreateToken(ctx context.Context, req *CreateTokenRequest) (*Token, error) {
	var t Token
	if err := c.http.Post(ctx, "/account/token", req, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// RevokeToken revokes a token
func (c *AuthClient) RevokeToken(ctx context.Context, token string) error {
	return c.http.Delete(ctx, "/account/token/"+token)
}

// DeleteAccount deletes the authenticated user's account
func (c *AuthClient) DeleteAccount(ctx context.Context) error {
	return c.http.Delete(ctx, "/account/delete")
}
