package whooktown

import (
	"context"
)

// AudioClient provides audio control functionality
type AudioClient struct {
	http *httpClient
}

// AudioCommand represents an audio control command
type AudioCommand struct {
	Command     string `json:"command"` // play, stop, volume, mood, toggle
	LayoutID    string `json:"layout_id"`
	Mood        string `json:"mood,omitempty"` // calm, active, tension, critical, epic
	MusicVolume *int   `json:"music_volume,omitempty"` // 0-100
	SfxVolume   *int   `json:"sfx_volume,omitempty"`   // 0-100
	Enabled     *bool  `json:"enabled,omitempty"`
	AutoMood    *bool  `json:"auto_mood,omitempty"`
}

// SendCommand sends an audio command
func (c *AudioClient) SendCommand(ctx context.Context, cmd *AudioCommand) error {
	return c.http.Post(ctx, "/ui/audio/command", cmd, nil)
}

// Play starts audio playback
func (c *AudioClient) Play(ctx context.Context, layoutID string) error {
	cmd := &AudioCommand{
		Command:  "play",
		LayoutID: layoutID,
	}
	return c.SendCommand(ctx, cmd)
}

// Stop stops audio playback
func (c *AudioClient) Stop(ctx context.Context, layoutID string) error {
	cmd := &AudioCommand{
		Command:  "stop",
		LayoutID: layoutID,
	}
	return c.SendCommand(ctx, cmd)
}

// SetMood sets the audio mood
func (c *AudioClient) SetMood(ctx context.Context, layoutID string, mood Mood) error {
	cmd := &AudioCommand{
		Command:  "mood",
		LayoutID: layoutID,
		Mood:     string(mood),
	}
	return c.SendCommand(ctx, cmd)
}

// SetVolume sets the music and/or SFX volume
func (c *AudioClient) SetVolume(ctx context.Context, layoutID string, musicVolume, sfxVolume *int) error {
	cmd := &AudioCommand{
		Command:     "volume",
		LayoutID:    layoutID,
		MusicVolume: musicVolume,
		SfxVolume:   sfxVolume,
	}
	return c.SendCommand(ctx, cmd)
}

// SetMusicVolume sets only the music volume
func (c *AudioClient) SetMusicVolume(ctx context.Context, layoutID string, volume int) error {
	return c.SetVolume(ctx, layoutID, &volume, nil)
}

// SetSfxVolume sets only the SFX volume
func (c *AudioClient) SetSfxVolume(ctx context.Context, layoutID string, volume int) error {
	return c.SetVolume(ctx, layoutID, nil, &volume)
}

// Enable enables audio for a layout
func (c *AudioClient) Enable(ctx context.Context, layoutID string) error {
	enabled := true
	cmd := &AudioCommand{
		Command:  "toggle",
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// Disable disables audio for a layout
func (c *AudioClient) Disable(ctx context.Context, layoutID string) error {
	enabled := false
	cmd := &AudioCommand{
		Command:  "toggle",
		LayoutID: layoutID,
		Enabled:  &enabled,
	}
	return c.SendCommand(ctx, cmd)
}

// EnableAutoMood enables automatic mood selection
func (c *AudioClient) EnableAutoMood(ctx context.Context, layoutID string) error {
	autoMood := true
	cmd := &AudioCommand{
		Command:  "toggle",
		LayoutID: layoutID,
		AutoMood: &autoMood,
	}
	return c.SendCommand(ctx, cmd)
}

// DisableAutoMood disables automatic mood selection
func (c *AudioClient) DisableAutoMood(ctx context.Context, layoutID string) error {
	autoMood := false
	cmd := &AudioCommand{
		Command:  "toggle",
		LayoutID: layoutID,
		AutoMood: &autoMood,
	}
	return c.SendCommand(ctx, cmd)
}

// GetStates returns audio states for all layouts
func (c *AudioClient) GetStates(ctx context.Context) ([]AudioState, error) {
	var states []AudioState
	if err := c.http.Get(ctx, "/ui/audio", &states); err != nil {
		return nil, err
	}
	return states, nil
}
