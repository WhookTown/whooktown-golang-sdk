package whooktown

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

// BackofficeClient provides admin functionality via the backoffice API
type BackofficeClient struct {
	http *httpClient
}

// === Health & Stats ===

// Health checks the backoffice API health
func (c *BackofficeClient) Health(ctx context.Context) error {
	return c.http.Get(ctx, "/api/health", nil)
}

// GetStats returns dashboard statistics
func (c *BackofficeClient) GetStats(ctx context.Context) (*Stats, error) {
	var stats Stats
	if err := c.http.Get(ctx, "/api/stats", &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// === Account Management ===

// ListAccounts returns all accounts
func (c *BackofficeClient) ListAccounts(ctx context.Context) ([]Account, error) {
	var accounts []Account
	if err := c.http.Get(ctx, "/api/accounts", &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccount returns an account by ID
func (c *BackofficeClient) GetAccount(ctx context.Context, accountID uuid.UUID) (*Account, error) {
	var account Account
	if err := c.http.Get(ctx, "/api/accounts/"+accountID.String(), &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// CreateAccountRequest represents a request to create an account
type CreateAccountRequest struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
	Type  string `json:"type"` // admin, user, viewer, sensor
}

// CreateAccount creates a new account
func (c *BackofficeClient) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*Account, error) {
	var account Account
	if err := c.http.Post(ctx, "/api/accounts", req, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// UpdateAccountRequest represents a request to update an account
type UpdateAccountRequest struct {
	Email     string `json:"email,omitempty"`
	Validated *bool  `json:"validated,omitempty"`
}

// UpdateAccount updates an account
func (c *BackofficeClient) UpdateAccount(ctx context.Context, accountID uuid.UUID, req *UpdateAccountRequest) (*Account, error) {
	var account Account
	if err := c.http.Put(ctx, "/api/accounts/"+accountID.String(), req, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// DeleteAccount deletes an account
func (c *BackofficeClient) DeleteAccount(ctx context.Context, accountID uuid.UUID) error {
	return c.http.Delete(ctx, "/api/accounts/"+accountID.String())
}

// LockAccount locks an account with an optional reason
func (c *BackofficeClient) LockAccount(ctx context.Context, accountID uuid.UUID, reason string) error {
	body := map[string]string{}
	if reason != "" {
		body["reason"] = reason
	}
	return c.http.Put(ctx, "/api/accounts/"+accountID.String()+"/lock", body, nil)
}

// UnlockAccount unlocks an account
func (c *BackofficeClient) UnlockAccount(ctx context.Context, accountID uuid.UUID) error {
	return c.http.Put(ctx, "/api/accounts/"+accountID.String()+"/unlock", nil, nil)
}

// === Token Management ===

// ListAccountTokens returns all tokens for an account
func (c *BackofficeClient) ListAccountTokens(ctx context.Context, accountID uuid.UUID) ([]Token, error) {
	var tokens []Token
	if err := c.http.Get(ctx, "/api/accounts/"+accountID.String()+"/tokens", &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

// CreateAccountTokenRequest represents a request to create a token for an account
type CreateAccountTokenRequest struct {
	Type       string        `json:"type"`
	Name       string        `json:"name,omitempty"`
	Expiration time.Duration `json:"expiration,omitempty"`
}

// CreateAccountToken creates a new token for an account
func (c *BackofficeClient) CreateAccountToken(ctx context.Context, accountID uuid.UUID, req *CreateAccountTokenRequest) (*Token, error) {
	body := map[string]interface{}{
		"type": req.Type,
	}
	if req.Name != "" {
		body["name"] = req.Name
	}
	if req.Expiration > 0 {
		body["expiration"] = req.Expiration.String()
	}
	var token Token
	if err := c.http.Post(ctx, "/api/accounts/"+accountID.String()+"/tokens", body, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteToken revokes a token
func (c *BackofficeClient) DeleteToken(ctx context.Context, token string) error {
	return c.http.Delete(ctx, "/api/tokens/"+token)
}

// === Layout Management ===

// ListAccountLayouts returns all layouts for an account
func (c *BackofficeClient) ListAccountLayouts(ctx context.Context, accountID uuid.UUID) ([]LayoutDB, error) {
	var layouts []LayoutDB
	if err := c.http.Get(ctx, "/api/accounts/"+accountID.String()+"/layouts", &layouts); err != nil {
		return nil, err
	}
	return layouts, nil
}

// DeleteAccountLayout deletes a layout for an account
func (c *BackofficeClient) DeleteAccountLayout(ctx context.Context, accountID, layoutID uuid.UUID) error {
	return c.http.Delete(ctx, "/api/accounts/"+accountID.String()+"/layouts/"+layoutID.String())
}

// === Subscription Management ===

// GetSubscriptionStats returns subscription statistics
func (c *BackofficeClient) GetSubscriptionStats(ctx context.Context) (*SubscriptionStats, error) {
	var stats SubscriptionStats
	if err := c.http.Get(ctx, "/api/subscriptions/stats", &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// ListSubscriptions returns all subscriptions
func (c *BackofficeClient) ListSubscriptions(ctx context.Context) ([]Subscription, error) {
	var subscriptions []Subscription
	if err := c.http.Get(ctx, "/api/subscriptions", &subscriptions); err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// ListPlans returns all available subscription plans
func (c *BackofficeClient) ListPlans(ctx context.Context) ([]Plan, error) {
	var plans []Plan
	if err := c.http.Get(ctx, "/api/subscriptions/plans", &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

// GetAccountSubscription returns the subscription for an account
func (c *BackofficeClient) GetAccountSubscription(ctx context.Context, accountID uuid.UUID) (*Subscription, error) {
	var subscription Subscription
	if err := c.http.Get(ctx, "/api/accounts/"+accountID.String()+"/subscription", &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}

// UpdateAccountSubscription updates the subscription for an account
func (c *BackofficeClient) UpdateAccountSubscription(ctx context.Context, accountID uuid.UUID, planID string) (*Subscription, error) {
	body := map[string]string{
		"plan_id": planID,
	}
	var subscription Subscription
	if err := c.http.Put(ctx, "/api/accounts/"+accountID.String()+"/subscription", body, &subscription); err != nil {
		return nil, err
	}
	return &subscription, nil
}

// === Asset Types ===

// AssetTypeConfig represents an asset type configuration
type AssetTypeConfig struct {
	TypeName    string `json:"type_name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description,omitempty"`
}

// ListAssetTypes returns all asset type configurations
func (c *BackofficeClient) ListAssetTypes(ctx context.Context) ([]AssetTypeConfig, error) {
	var types []AssetTypeConfig
	if err := c.http.Get(ctx, "/api/asset-types", &types); err != nil {
		return nil, err
	}
	return types, nil
}

// UpdateAssetType updates an asset type configuration
func (c *BackofficeClient) UpdateAssetType(ctx context.Context, typeName string, enabled bool) (*AssetTypeConfig, error) {
	body := map[string]bool{
		"enabled": enabled,
	}
	var config AssetTypeConfig
	if err := c.http.Put(ctx, "/api/asset-types/"+typeName, body, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
