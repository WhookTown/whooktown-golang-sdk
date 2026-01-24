package whooktown

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// Status represents the status of a sensor/building
type Status string

const (
	StatusOnline   Status = "online"
	StatusOffline  Status = "offline"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
)

// Activity represents the activity level
type Activity string

const (
	ActivitySlow   Activity = "slow"
	ActivityNormal Activity = "normal"
	ActivityFast   Activity = "fast"
)

// Speed represents traffic speed levels
type Speed string

const (
	SpeedSlow   Speed = "slow"
	SpeedNormal Speed = "normal"
	SpeedFast   Speed = "fast"
)

// CameraMode represents the camera mode
type CameraMode string

const (
	CameraModeOrbit   CameraMode = "orbit"
	CameraModeFPS     CameraMode = "fps"
	CameraModeFlyover CameraMode = "flyover"
)

// Orientation represents compass directions
type Orientation string

const (
	OrientationN  Orientation = "N"
	OrientationNE Orientation = "NE"
	OrientationE  Orientation = "E"
	OrientationSE Orientation = "SE"
	OrientationS  Orientation = "S"
	OrientationSW Orientation = "SW"
	OrientationW  Orientation = "W"
	OrientationNW Orientation = "NW"
)

// BuildingType constants for all supported building types
const (
	BuildingWindmill      = "windmill"
	BuildingDataCenter    = "data_center"
	BuildingArcade        = "arcade"
	BuildingPyramid       = "pyramid"
	BuildingTowerA        = "tower_a"
	BuildingTowerB        = "tower_b"
	BuildingSupervisor    = "supervisor"
	BuildingBank          = "bank"
	BuildingMonitorTube   = "monitor_tube"
	BuildingBakery        = "bakery"
	BuildingHouseA        = "house_a"
	BuildingHouseB        = "house_b"
	BuildingHouseC        = "house_c"
	BuildingTree          = "tree"
	BuildingDisplayA      = "display_a"
	BuildingTrafficLight  = "traffic_light"
	BuildingFarmBuildingA = "farm_building_a"
	BuildingFarmBuildingB = "farm_building_b"
	BuildingFarmSilo      = "farm_silo"
	BuildingFarmFieldA    = "farm_field_a"
	BuildingFarmFieldB    = "farm_field_b"
	BuildingFarmCattleA   = "farm_cattle_a"
	BuildingGrass         = "grass"
	BuildingSpire         = "spire"
	BuildingLedFacade     = "led_facade"
	BuildingTwinTowers    = "twin_towers"
	BuildingDiamondTower  = "diamond_tower"
)

// TokenType represents the type of authentication token
type TokenType string

const (
	TokenTypeAdmin  TokenType = "admin"
	TokenTypeUser   TokenType = "user"
	TokenTypeViewer TokenType = "viewer"
	TokenTypeSensor TokenType = "sensor"
)

// Account represents a user account
type Account struct {
	ID         uuid.UUID  `json:"id,omitempty"`
	Email      string     `json:"email,omitempty"`
	Validated  bool       `json:"validated,omitempty"`
	Locked     bool       `json:"locked,omitempty"`
	LockedAt   *time.Time `json:"locked_at,omitempty"`
	LockReason string     `json:"lock_reason,omitempty"`
	CreatedAt  time.Time  `json:"created_at,omitempty"`
}

// Token represents an authentication token
type Token struct {
	Token          string            `json:"app_token,omitempty"`
	ValidationLink string            `json:"validation_link,omitempty"`
	Name           string            `json:"name,omitempty"`
	Type           string            `json:"type,omitempty"`
	Roles          map[string]string `json:"roles,omitempty"`
	AccountID      uuid.UUID         `json:"account_id,omitempty"`
	Account        *Account          `json:"account,omitempty"`
	CreatedAt      time.Time         `json:"created_at,omitempty"`
	UpdatedAt      time.Time         `json:"updated_at,omitempty"`
	ExpiredAt      time.Time         `json:"expired_at,omitempty"`
}

// SensorData represents sensor payload
type SensorData struct {
	ID       uuid.UUID `json:"id"`
	Status   Status    `json:"status,omitempty"`
	Activity Activity  `json:"activity,omitempty"`

	// Building-specific fields (optional)
	Quantity        string `json:"quantity,omitempty"`         // Bank: none, low, medium, full
	Amount          int    `json:"amount,omitempty"`           // Bank: displayed amount
	Text1           string `json:"text1,omitempty"`            // DisplayA text
	Text2           string `json:"text2,omitempty"`            // DisplayA text
	Text3           string `json:"text3,omitempty"`            // DisplayA text
	TowerText       string `json:"towerText,omitempty"`        // TowerA LED text
	TowerBText      string `json:"towerBText,omitempty"`       // TowerB LED text
	RingCount       int    `json:"ringCount,omitempty"`        // DisplayA: 2 or 3
	DancerEnabled   *bool  `json:"dancerEnabled,omitempty"`    // Arcade
	MusicEnabled    *bool  `json:"musicEnabled,omitempty"`     // Arcade
	SignText        string `json:"signText,omitempty"`         // Arcade sign
	FaceRotation    *bool  `json:"faceRotationEnabled,omitempty"` // Supervisor
	CPUUsage        int    `json:"cpuUsage,omitempty"`         // DataCenter: 0-100
	RAMUsage        int    `json:"ramUsage,omitempty"`         // DataCenter: 0-100
	NetworkTraffic  int    `json:"networkTraffic,omitempty"`   // DataCenter: 0-100
	ActiveConns     int    `json:"activeConnections,omitempty"`// DataCenter
	Temperature     int    `json:"temperature,omitempty"`      // DataCenter: Celsius
	AlertLevel      string `json:"alertLevel,omitempty"`       // DataCenter: normal, warning, critical
	BandCount       int    `json:"bandCount,omitempty"`        // MonitorTube: 3-7
	Bands           []Band `json:"bands,omitempty"`            // MonitorTube bands

	// Extra fields for custom data
	Extra map[string]interface{} `json:"-"`
}

// Band represents a monitor tube band
type Band struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// MarshalJSON implements custom JSON marshaling to flatten Extra fields
func (s *SensorData) MarshalJSON() ([]byte, error) {
	type Alias SensorData
	data, err := json.Marshal((*Alias)(s))
	if err != nil {
		return nil, err
	}

	if len(s.Extra) == 0 {
		return data, nil
	}

	// Merge Extra fields into the JSON
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	for k, v := range s.Extra {
		if _, exists := m[k]; !exists {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

// Layout represents a city layout
type Layout struct {
	ID        uuid.UUID       `json:"id,omitempty"`
	Name      string          `json:"name"`
	Grid      Grid            `json:"grid"`
	Buildings []Building      `json:"buildings"`
	Roads     json.RawMessage `json:"roads,omitempty"`
}

// LayoutDB represents a persisted layout
type LayoutDB struct {
	AccountID     uuid.UUID       `json:"account_id"`
	LayoutID      uuid.UUID       `json:"layout_id"`
	ReceivedAt    time.Time       `json:"received_at"`
	Data          json.RawMessage `json:"data"`
	Archived      bool            `json:"archived"`
	ArchivedAt    *time.Time      `json:"archived_at,omitempty"`
	ArchiveReason string          `json:"archive_reason,omitempty"`
}

// Grid represents the city grid dimensions
type Grid struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Building represents a building in the layout
type Building struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name,omitempty"`
	Type        string    `json:"type"`
	Roles       []string  `json:"roles,omitempty"`
	Location    Location  `json:"location"`
	Orientation string    `json:"orientation,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}

// Location represents building position on grid
type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Vector3 represents a 3D position or rotation
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// CameraPreset represents a saved camera position
type CameraPreset struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	LayoutID  uuid.UUID `json:"layout_id"`
	Name      string    `json:"name"`
	PositionX float64   `json:"position_x"`
	PositionY float64   `json:"position_y"`
	PositionZ float64   `json:"position_z"`
	RotationX float64   `json:"rotation_x"`
	RotationY float64   `json:"rotation_y"`
	RotationZ float64   `json:"rotation_z"`
	FOV       float64   `json:"fov"`
	Mode      string    `json:"mode"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CameraPath represents a camera path with checkpoints
type CameraPath struct {
	ID          uuid.UUID              `json:"id"`
	AccountID   uuid.UUID              `json:"account_id"`
	LayoutID    uuid.UUID              `json:"layout_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Loop        bool                   `json:"loop"`
	Checkpoints []CameraPathCheckpoint `json:"checkpoints"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// CameraPathCheckpoint represents a waypoint on the camera path
type CameraPathCheckpoint struct {
	ID                 uuid.UUID `json:"id"`
	PathID             uuid.UUID `json:"path_id"`
	GridX              int       `json:"grid_x"`
	GridY              int       `json:"grid_y"`
	OrderIndex         int       `json:"order_index"`
	Orientation        string    `json:"orientation"`
	Altitude           int       `json:"altitude"`
	Tilt               int       `json:"tilt"`
	Zoom               int       `json:"zoom"`
	TransitionDuration float64   `json:"transition_duration"`
	HoldDuration       float64   `json:"hold_duration"`
	CreatedAt          time.Time `json:"created_at"`
}

// CameraSequence represents a deprecated camera sequence
type CameraSequence struct {
	ID          uuid.UUID                `json:"id"`
	AccountID   uuid.UUID                `json:"account_id"`
	LayoutID    uuid.UUID                `json:"layout_id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description,omitempty"`
	Loop        bool                     `json:"loop"`
	Keyframes   []CameraSequenceKeyframe `json:"keyframes"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

// CameraSequenceKeyframe represents a keyframe in a camera sequence
type CameraSequenceKeyframe struct {
	ID                 uuid.UUID `json:"id"`
	SequenceID         uuid.UUID `json:"sequence_id"`
	PresetID           uuid.UUID `json:"preset_id"`
	OrderIndex         int       `json:"order_index"`
	TransitionDuration float64   `json:"transition_duration"`
	HoldDuration       float64   `json:"hold_duration"`
	CreatedAt          time.Time `json:"created_at"`
}

// AssetGroup represents a group of buildings
type AssetGroup struct {
	ID        uuid.UUID   `json:"id"`
	AccountID uuid.UUID   `json:"account_id"`
	LayoutID  uuid.UUID   `json:"layout_id"`
	Name      string      `json:"name"`
	Members   []uuid.UUID `json:"members"`
	CreatedAt time.Time   `json:"created_at"`
}

// ConnectedScene represents a connected threejs-scene instance
type ConnectedScene struct {
	SceneID       string    `json:"scene_id"`
	LayoutID      uuid.UUID `json:"layout_id"`
	ConnectedAt   time.Time `json:"connected_at"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

// TrafficState represents traffic control state
type TrafficState struct {
	LayoutID      string `json:"layout_id"`
	Density       int    `json:"density"`
	Speed         string `json:"speed"`
	Enabled       bool   `json:"enabled"`
	LabelsVisible bool   `json:"labels_visible"`
}

// QuotaInfo represents account quota information
type QuotaInfo struct {
	Plan   string `json:"plan"`
	Status string `json:"status"`
	Layouts struct {
		Used     int `json:"used"`
		Max      int `json:"max"`
		Archived int `json:"archived"`
	} `json:"layouts"`
	AssetsPerLayout struct {
		Max int `json:"max"`
	} `json:"assets_per_layout"`
}

// Workflow represents a workflow definition
type Workflow struct {
	AccountID uuid.UUID       `json:"account_id"`
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Worker    string          `json:"worker"`
	Version   string          `json:"version"`
	Graph     json.RawMessage `json:"graph"`
	Enabled   bool            `json:"enabled"`
	CreatedAt time.Time       `json:"created_at"`
}

// FlowNode represents a node in the workflow graph
type FlowNode struct {
	ID         string   `json:"id"`
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	Operator   string   `json:"operator"`
	Inputs     []string `json:"inputs,omitempty"`
	Values     []string `json:"values,omitempty"`
	Condition  []string `json:"condition,omitempty"`
	Latch      bool     `json:"latch,omitempty"`
	LatchValue string   `json:"latchValue,omitempty"`

	// Control node fields
	LayoutID    string `json:"layout_id,omitempty"`
	Density     int    `json:"density,omitempty"`
	Speed       string `json:"speed,omitempty"`
	Enabled     *bool  `json:"enabled,omitempty"`
	Command     string `json:"command,omitempty"`
	Mood        string `json:"mood,omitempty"`
	MusicVolume int    `json:"music_volume,omitempty"`
	PathID      string `json:"path_id,omitempty"`
	Action      string `json:"action,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	OutputField string `json:"output_field,omitempty"`
	OutputValue string `json:"output_value,omitempty"`
}

// Operation describes a workflow operation
type Operation struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MaxLenInput int    `json:"max_len_input,omitempty"`
	OutputType  string `json:"output_type,omitempty"`
	InputsType  string `json:"inputs_type,omitempty"`
}

// Subscription represents an account subscription
type Subscription struct {
	ID              uuid.UUID  `json:"id"`
	AccountID       uuid.UUID  `json:"account_id"`
	PlanID          string     `json:"plan_id"`
	Status          string     `json:"status"`
	StripeSubID     string     `json:"stripe_subscription_id,omitempty"`
	StripeCustID    string     `json:"stripe_customer_id,omitempty"`
	CurrentPeriodStart time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   time.Time `json:"current_period_end,omitempty"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Plan represents a subscription plan
type Plan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MaxAssets   int    `json:"max_assets"`
	MaxLayouts  int    `json:"max_layouts"`
	StripePriceID string `json:"stripe_price_id,omitempty"`
}

// Stats represents dashboard statistics
type Stats struct {
	TotalAccounts    int `json:"total_accounts"`
	ActiveAccounts   int `json:"active_accounts"`
	LockedAccounts   int `json:"locked_accounts"`
	TotalTokens      int `json:"total_tokens"`
	TotalLayouts     int `json:"total_layouts"`
}

// SubscriptionStats represents subscription statistics
type SubscriptionStats struct {
	TotalSubscriptions int            `json:"total_subscriptions"`
	ByPlan             map[string]int `json:"by_plan"`
	ByStatus           map[string]int `json:"by_status"`
	MRR                float64        `json:"mrr"`
}
