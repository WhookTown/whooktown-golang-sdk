package whooktown

import (
	"context"

	"github.com/gofrs/uuid"
)

// UIClient provides access to the UI endpoint for layout management
type UIClient struct {
	http *httpClient
}

// CreateLayout creates or updates a layout
func (c *UIClient) CreateLayout(ctx context.Context, layout *Layout) (*LayoutDB, error) {
	var result LayoutDB
	if err := c.http.Post(ctx, "/ui/layout", layout, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateLayout is an alias for CreateLayout (upsert operation)
func (c *UIClient) UpdateLayout(ctx context.Context, layout *Layout) (*LayoutDB, error) {
	return c.CreateLayout(ctx, layout)
}

// DeleteLayout deletes a layout by ID
func (c *UIClient) DeleteLayout(ctx context.Context, layoutID uuid.UUID) error {
	return c.http.Delete(ctx, "/ui/layout/"+layoutID.String())
}

// GetQuota returns the current quota usage for the account
func (c *UIClient) GetQuota(ctx context.Context) (*QuotaInfo, error) {
	var quota QuotaInfo
	if err := c.http.Get(ctx, "/ui/quota", &quota); err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetArchivedLayouts returns archived layouts
func (c *UIClient) GetArchivedLayouts(ctx context.Context) ([]LayoutDB, error) {
	var layouts []LayoutDB
	if err := c.http.Get(ctx, "/ui/layout/archived", &layouts); err != nil {
		return nil, err
	}
	return layouts, nil
}

// RestoreLayout restores an archived layout
func (c *UIClient) RestoreLayout(ctx context.Context, layoutID uuid.UUID) error {
	return c.http.Post(ctx, "/ui/layout/"+layoutID.String()+"/restore", nil, nil)
}

// ListScenes returns connected threejs-scene instances
func (c *UIClient) ListScenes(ctx context.Context) ([]ConnectedScene, error) {
	var scenes []ConnectedScene
	if err := c.http.Get(ctx, "/ui/scenes", &scenes); err != nil {
		return nil, err
	}
	return scenes, nil
}

// SceneStateRequest represents a request to update scene state
type SceneStateRequest struct {
	LayoutID       string `json:"layout_id"`
	FlyoverEnabled *bool  `json:"flyover_enabled,omitempty"`
	FlyoverSpeed   *int   `json:"flyover_speed,omitempty"`
	ActivePathID   string `json:"active_path_id,omitempty"`
}

// UpdateSceneState updates the state of a connected scene
func (c *UIClient) UpdateSceneState(ctx context.Context, sceneID string, req *SceneStateRequest) error {
	return c.http.Post(ctx, "/ui/scene/"+sceneID+"/state", req, nil)
}
