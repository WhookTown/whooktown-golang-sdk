package whooktown

import (
	"context"

	"github.com/gofrs/uuid"
)

// CameraClient provides camera control functionality
type CameraClient struct {
	http *httpClient
}

// CameraCommand represents a camera command
type CameraCommand struct {
	Command      string   `json:"command"`                   // position, preset, mode, sequence, path
	LayoutID     string   `json:"layout_id"`
	Position     *Vector3 `json:"position,omitempty"`
	Rotation     *Vector3 `json:"rotation,omitempty"`
	FOV          float64  `json:"fov,omitempty"`
	Animate      bool     `json:"animate,omitempty"`
	Duration     float64  `json:"duration,omitempty"`
	PresetID     string   `json:"preset_id,omitempty"`
	Mode         string   `json:"mode,omitempty"`            // orbit, fps, flyover
	FlyoverSpeed float64  `json:"flyover_speed,omitempty"`
	Action       string   `json:"action,omitempty"`          // play, pause, stop
	PathID       string   `json:"path_id,omitempty"`
	SequenceID   string   `json:"sequence_id,omitempty"`
}

// SendCommand sends a camera command
func (c *CameraClient) SendCommand(ctx context.Context, cmd *CameraCommand) error {
	return c.http.Post(ctx, "/ui/camera/command", cmd, nil)
}

// SetPosition sets the camera position
func (c *CameraClient) SetPosition(ctx context.Context, layoutID string, position, rotation *Vector3, fov float64, animate bool, duration float64) error {
	cmd := &CameraCommand{
		Command:  "position",
		LayoutID: layoutID,
		Position: position,
		Rotation: rotation,
		FOV:      fov,
		Animate:  animate,
		Duration: duration,
	}
	return c.SendCommand(ctx, cmd)
}

// SetMode sets the camera mode
func (c *CameraClient) SetMode(ctx context.Context, layoutID string, mode CameraMode, flyoverSpeed float64) error {
	cmd := &CameraCommand{
		Command:      "mode",
		LayoutID:     layoutID,
		Mode:         string(mode),
		FlyoverSpeed: flyoverSpeed,
	}
	return c.SendCommand(ctx, cmd)
}

// GoToPreset moves camera to a preset position
func (c *CameraClient) GoToPreset(ctx context.Context, layoutID, presetID string, animate bool, duration float64) error {
	cmd := &CameraCommand{
		Command:  "preset",
		LayoutID: layoutID,
		PresetID: presetID,
		Animate:  animate,
		Duration: duration,
	}
	return c.SendCommand(ctx, cmd)
}

// PlayPath starts playing a camera path
func (c *CameraClient) PlayPath(ctx context.Context, layoutID, pathID string) error {
	cmd := &CameraCommand{
		Command:  "path",
		LayoutID: layoutID,
		PathID:   pathID,
		Action:   "play",
	}
	return c.SendCommand(ctx, cmd)
}

// PausePath pauses the current camera path
func (c *CameraClient) PausePath(ctx context.Context, layoutID string) error {
	cmd := &CameraCommand{
		Command:  "path",
		LayoutID: layoutID,
		Action:   "pause",
	}
	return c.SendCommand(ctx, cmd)
}

// StopPath stops the current camera path
func (c *CameraClient) StopPath(ctx context.Context, layoutID string) error {
	cmd := &CameraCommand{
		Command:  "path",
		LayoutID: layoutID,
		Action:   "stop",
	}
	return c.SendCommand(ctx, cmd)
}

// === Presets ===

// ListPresets returns camera presets for a layout
func (c *CameraClient) ListPresets(ctx context.Context, layoutID uuid.UUID) ([]CameraPreset, error) {
	var presets []CameraPreset
	if err := c.http.Get(ctx, "/ui/presets/"+layoutID.String(), &presets); err != nil {
		return nil, err
	}
	return presets, nil
}

// CreatePresetRequest represents a request to create a camera preset
type CreatePresetRequest struct {
	LayoutID  uuid.UUID `json:"layout_id"`
	Name      string    `json:"name"`
	PositionX float64   `json:"position_x"`
	PositionY float64   `json:"position_y"`
	PositionZ float64   `json:"position_z"`
	RotationX float64   `json:"rotation_x"`
	RotationY float64   `json:"rotation_y"`
	RotationZ float64   `json:"rotation_z"`
	FOV       float64   `json:"fov,omitempty"`
	Mode      string    `json:"mode,omitempty"`
}

// CreatePreset creates a new camera preset
func (c *CameraClient) CreatePreset(ctx context.Context, req *CreatePresetRequest) (*CameraPreset, error) {
	var preset CameraPreset
	if err := c.http.Post(ctx, "/ui/presets", req, &preset); err != nil {
		return nil, err
	}
	return &preset, nil
}

// UpdatePresetRequest represents a request to update a camera preset
type UpdatePresetRequest struct {
	Name      string  `json:"name,omitempty"`
	PositionX float64 `json:"position_x,omitempty"`
	PositionY float64 `json:"position_y,omitempty"`
	PositionZ float64 `json:"position_z,omitempty"`
	RotationX float64 `json:"rotation_x,omitempty"`
	RotationY float64 `json:"rotation_y,omitempty"`
	RotationZ float64 `json:"rotation_z,omitempty"`
	FOV       float64 `json:"fov,omitempty"`
	Mode      string  `json:"mode,omitempty"`
}

// UpdatePreset updates a camera preset
func (c *CameraClient) UpdatePreset(ctx context.Context, presetID uuid.UUID, req *UpdatePresetRequest) (*CameraPreset, error) {
	var preset CameraPreset
	if err := c.http.Put(ctx, "/ui/presets/"+presetID.String(), req, &preset); err != nil {
		return nil, err
	}
	return &preset, nil
}

// DeletePreset deletes a camera preset
func (c *CameraClient) DeletePreset(ctx context.Context, presetID uuid.UUID) error {
	return c.http.Delete(ctx, "/ui/presets/"+presetID.String())
}

// SetDefaultPreset sets a preset as the default for its layout
func (c *CameraClient) SetDefaultPreset(ctx context.Context, presetID uuid.UUID) error {
	return c.http.Post(ctx, "/ui/presets/"+presetID.String()+"/default", nil, nil)
}

// === Paths ===

// ListPaths returns camera paths for a layout
func (c *CameraClient) ListPaths(ctx context.Context, layoutID uuid.UUID) ([]CameraPath, error) {
	var paths []CameraPath
	if err := c.http.Get(ctx, "/ui/paths/"+layoutID.String(), &paths); err != nil {
		return nil, err
	}
	return paths, nil
}

// GetPath returns a single camera path
func (c *CameraClient) GetPath(ctx context.Context, layoutID, pathID uuid.UUID) (*CameraPath, error) {
	var path CameraPath
	if err := c.http.Get(ctx, "/ui/paths/"+layoutID.String()+"/"+pathID.String(), &path); err != nil {
		return nil, err
	}
	return &path, nil
}

// CreatePathRequest represents a request to create a camera path
type CreatePathRequest struct {
	LayoutID    uuid.UUID `json:"layout_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Loop        bool      `json:"loop,omitempty"`
}

// CreatePath creates a new camera path
func (c *CameraClient) CreatePath(ctx context.Context, req *CreatePathRequest) (*CameraPath, error) {
	var path CameraPath
	if err := c.http.Post(ctx, "/ui/paths", req, &path); err != nil {
		return nil, err
	}
	return &path, nil
}

// UpdatePathRequest represents a request to update a camera path
type UpdatePathRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Loop        *bool  `json:"loop,omitempty"`
}

// UpdatePath updates a camera path
func (c *CameraClient) UpdatePath(ctx context.Context, pathID uuid.UUID, req *UpdatePathRequest) (*CameraPath, error) {
	var path CameraPath
	if err := c.http.Put(ctx, "/ui/paths/"+pathID.String(), req, &path); err != nil {
		return nil, err
	}
	return &path, nil
}

// DeletePath deletes a camera path
func (c *CameraClient) DeletePath(ctx context.Context, pathID uuid.UUID) error {
	return c.http.Delete(ctx, "/ui/paths/"+pathID.String())
}

// AddCheckpointRequest represents a request to add a checkpoint to a path
type AddCheckpointRequest struct {
	GridX              int     `json:"grid_x"`
	GridY              int     `json:"grid_y"`
	Orientation        string  `json:"orientation,omitempty"` // N, NE, E, SE, S, SW, W, NW
	Altitude           int     `json:"altitude,omitempty"`    // 0-100
	Tilt               int     `json:"tilt,omitempty"`        // -90 to +90
	Zoom               int     `json:"zoom,omitempty"`        // 30-120 (FOV)
	TransitionDuration float64 `json:"transition_duration,omitempty"`
	HoldDuration       float64 `json:"hold_duration,omitempty"`
}

// AddCheckpoint adds a checkpoint to a camera path
func (c *CameraClient) AddCheckpoint(ctx context.Context, pathID uuid.UUID, req *AddCheckpointRequest) (*CameraPath, error) {
	var path CameraPath
	if err := c.http.Post(ctx, "/ui/paths/"+pathID.String()+"/checkpoints", req, &path); err != nil {
		return nil, err
	}
	return &path, nil
}

// UpdateCheckpoint updates a checkpoint on a camera path
func (c *CameraClient) UpdateCheckpoint(ctx context.Context, pathID, checkpointID uuid.UUID, req *AddCheckpointRequest) (*CameraPath, error) {
	var path CameraPath
	if err := c.http.Put(ctx, "/ui/paths/"+pathID.String()+"/checkpoints/"+checkpointID.String(), req, &path); err != nil {
		return nil, err
	}
	return &path, nil
}

// DeleteCheckpoint deletes a checkpoint from a camera path
func (c *CameraClient) DeleteCheckpoint(ctx context.Context, pathID, checkpointID uuid.UUID) error {
	return c.http.Delete(ctx, "/ui/paths/"+pathID.String()+"/checkpoints/"+checkpointID.String())
}

// ReorderCheckpoints reorders checkpoints on a camera path
func (c *CameraClient) ReorderCheckpoints(ctx context.Context, pathID uuid.UUID, checkpointIDs []uuid.UUID) (*CameraPath, error) {
	body := map[string][]uuid.UUID{
		"checkpoint_ids": checkpointIDs,
	}
	var path CameraPath
	if err := c.http.Put(ctx, "/ui/paths/"+pathID.String()+"/checkpoints/reorder", body, &path); err != nil {
		return nil, err
	}
	return &path, nil
}
