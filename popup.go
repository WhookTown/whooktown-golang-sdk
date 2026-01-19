package whooktown

import (
	"context"
)

// PopupClient provides popup/labels control functionality
type PopupClient struct {
	http *httpClient
}

// PopupCommand represents a popup control command
type PopupCommand struct {
	LayoutID    string   `json:"layout_id"`
	Command     string   `json:"command"` // labels, detail, close, close_all
	BuildingIDs []string `json:"building_ids,omitempty"`
	Enabled     *bool    `json:"enabled,omitempty"`
}

// SendCommand sends a popup command
func (c *PopupClient) SendCommand(ctx context.Context, cmd *PopupCommand) error {
	return c.http.Post(ctx, "/ui/popup/command", cmd, nil)
}

// ShowLabels shows labels for all buildings
func (c *PopupClient) ShowLabels(ctx context.Context, layoutID string) error {
	enabled := true
	cmd := &PopupCommand{
		Command:  "labels",
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// HideLabels hides labels for all buildings
func (c *PopupClient) HideLabels(ctx context.Context, layoutID string) error {
	enabled := false
	cmd := &PopupCommand{
		Command:  "labels",
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// ToggleLabels toggles labels visibility
func (c *PopupClient) ToggleLabels(ctx context.Context, layoutID string, enabled bool) error {
	cmd := &PopupCommand{
		Command:  "labels",
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// ShowDetail shows detail popup for specific buildings
func (c *PopupClient) ShowDetail(ctx context.Context, layoutID string, buildingIDs []string) error {
	cmd := &PopupCommand{
		Command:     "detail",
		LayoutID:    layoutID,
		BuildingIDs: buildingIDs,
	}
	return c.SendCommand(ctx, cmd)
}

// CloseDetail closes detail popup for specific buildings
func (c *PopupClient) CloseDetail(ctx context.Context, layoutID string, buildingIDs []string) error {
	cmd := &PopupCommand{
		Command:     "close",
		LayoutID:    layoutID,
		BuildingIDs: buildingIDs,
	}
	return c.SendCommand(ctx, cmd)
}

// CloseAllDetails closes all detail popups
func (c *PopupClient) CloseAllDetails(ctx context.Context, layoutID string) error {
	cmd := &PopupCommand{
		Command:  "close_all",
		LayoutID: layoutID,
	}
	return c.SendCommand(ctx, cmd)
}
