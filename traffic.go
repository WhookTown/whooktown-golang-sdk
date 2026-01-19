package whooktown

import (
	"context"
)

// TrafficClient provides traffic control functionality
type TrafficClient struct {
	http *httpClient
}

// TrafficCommand represents a traffic control command
type TrafficCommand struct {
	LayoutID string `json:"layout_id"`
	Density  int    `json:"density"` // 0-100
	Speed    string `json:"speed"`   // slow, normal, fast
	Enabled  *bool  `json:"enabled,omitempty"`
}

// SendCommand sends a traffic command
func (c *TrafficClient) SendCommand(ctx context.Context, cmd *TrafficCommand) error {
	return c.http.Post(ctx, "/ui/traffic/command", cmd, nil)
}

// SetTraffic sets the traffic state for a layout
func (c *TrafficClient) SetTraffic(ctx context.Context, layoutID string, density int, speed Speed, enabled bool) error {
	cmd := &TrafficCommand{
		LayoutID: layoutID,
		Density:  density,
		Speed:    string(speed),
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// SetDensity sets only the traffic density
func (c *TrafficClient) SetDensity(ctx context.Context, layoutID string, density int) error {
	cmd := &TrafficCommand{
		LayoutID: layoutID,
		Density:  density,
	}
	return c.SendCommand(ctx, cmd)
}

// SetSpeed sets only the traffic speed
func (c *TrafficClient) SetSpeed(ctx context.Context, layoutID string, speed Speed) error {
	cmd := &TrafficCommand{
		LayoutID: layoutID,
		Speed:    string(speed),
	}
	return c.SendCommand(ctx, cmd)
}

// Enable enables traffic for a layout
func (c *TrafficClient) Enable(ctx context.Context, layoutID string) error {
	enabled := true
	cmd := &TrafficCommand{
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// Disable disables traffic for a layout
func (c *TrafficClient) Disable(ctx context.Context, layoutID string) error {
	enabled := false
	cmd := &TrafficCommand{
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// GetStates returns traffic states for all layouts
func (c *TrafficClient) GetStates(ctx context.Context) ([]TrafficState, error) {
	var states []TrafficState
	if err := c.http.Get(ctx, "/ui/traffic", &states); err != nil {
		return nil, err
	}
	return states, nil
}
