package whooktown

import (
	"context"
)

// SensorsClient provides access to the sensor endpoint
type SensorsClient struct {
	http *httpClient
}

// Send sends sensor data to whooktown
func (c *SensorsClient) Send(ctx context.Context, data *SensorData) error {
	return c.http.Post(ctx, "/sensors", data, nil)
}

// SendRaw sends raw sensor data (as a map) to whooktown
func (c *SensorsClient) SendRaw(ctx context.Context, data map[string]interface{}) error {
	return c.http.Post(ctx, "/sensors", data, nil)
}

// SendMultiple sends multiple sensor data points
func (c *SensorsClient) SendMultiple(ctx context.Context, data []*SensorData) error {
	for _, d := range data {
		if err := c.Send(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

// Health checks the sensor endpoint health
func (c *SensorsClient) Health(ctx context.Context) error {
	return c.http.Get(ctx, "/sensors/_health", nil)
}

// SetCameraMode sets the camera mode for a layout via sensor endpoint
func (c *SensorsClient) SetCameraMode(ctx context.Context, layoutID string, mode CameraMode, flyoverSpeed float64) error {
	body := map[string]interface{}{
		"layout_id": layoutID,
		"mode":      string(mode),
	}
	if flyoverSpeed > 0 {
		body["flyover_speed"] = flyoverSpeed
	}
	return c.http.Post(ctx, "/camera", body, nil)
}

// GetCameraStates returns camera states for all layouts
func (c *SensorsClient) GetCameraStates(ctx context.Context) ([]map[string]interface{}, error) {
	var states []map[string]interface{}
	if err := c.http.Get(ctx, "/camera", &states); err != nil {
		return nil, err
	}
	return states, nil
}

// SetTrafficState sets the traffic state for a layout via sensor endpoint
func (c *SensorsClient) SetTrafficState(ctx context.Context, layoutID string, density int, speed Speed, enabled bool) error {
	body := map[string]interface{}{
		"layout_id": layoutID,
		"density":   density,
		"speed":     string(speed),
		"enabled":   enabled,
	}
	return c.http.Post(ctx, "/traffic", body, nil)
}

// GetTrafficStates returns traffic states for all layouts
func (c *SensorsClient) GetTrafficStates(ctx context.Context) ([]TrafficState, error) {
	var states []TrafficState
	if err := c.http.Get(ctx, "/traffic", &states); err != nil {
		return nil, err
	}
	return states, nil
}
