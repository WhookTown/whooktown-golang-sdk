# whooktown-golang-sdk

Official Go SDK for [whooktown](https://whooktown.com) - the IT infrastructure visualization platform as a 3D city.

## Installation

```bash
go get github.com/fredericalix/whooktown-golang-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"

    whooktown "github.com/fredericalix/whooktown-golang-sdk"
    "github.com/gofrs/uuid"
)

func main() {
    // Create client (uses PROD by default, set WHOOKTOWN_ENV=DEV for development)
    client, err := whooktown.New(
        whooktown.WithToken(os.Getenv("WHOOKTOWN_TOKEN")),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Send sensor data
    err = client.Sensors.Send(ctx, &whooktown.SensorData{
        ID:       uuid.Must(uuid.NewV4()),
        Status:   whooktown.StatusOnline,
        Activity: whooktown.ActivityNormal,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Configuration

### Environments (PROD / DEV)

By default, the SDK connects to **production** servers (`whook.town`). To use development servers, set the `WHOOKTOWN_ENV` environment variable:

```bash
# Production (default)
export WHOOKTOWN_TOKEN="your-token"
go run main.go

# Development
export WHOOKTOWN_ENV=DEV
export WHOOKTOWN_TOKEN="your-token"
go run main.go
```

**Production URLs** (default):
| Service | URL |
|---------|-----|
| Auth | https://auth.whook.town |
| Sensors | https://sensors.whook.town |
| API (UI/Workflow) | https://api.whook.town |
| SSE/WebSocket | https://ws.whook.town |
| Backoffice | https://admin.whook.town |
| Subscription | https://subscription.whook.town |

**Development URLs** (`WHOOKTOWN_ENV=DEV`):
| Service | URL |
|---------|-----|
| Auth | https://auth.dev.whook.town |
| Sensors | https://sensors.dev.whook.town |
| API (UI/Workflow) | https://api.dev.whook.town |
| SSE/WebSocket | https://ws.dev.whook.town |
| Backoffice | https://admin.dev.whook.town |
| Subscription | https://subscription.dev.whook.town |

You can also set the environment programmatically:

```go
client, err := whooktown.New(
    whooktown.WithToken("your-token"),
    whooktown.WithEnvironment(whooktown.EnvDevelopment), // Override WHOOKTOWN_ENV
)
```

### Available Options

```go
client, err := whooktown.New(
    // Authentication token
    whooktown.WithToken("your-token"),

    // Admin secret for backoffice
    whooktown.WithAdminSecret("your-secret"),

    // Environment (PROD or DEV) - overrides WHOOKTOWN_ENV
    whooktown.WithEnvironment(whooktown.EnvProduction),

    // HTTP timeout
    whooktown.WithTimeout(30 * time.Second),

    // Retry configuration
    whooktown.WithRetry(3, time.Second),

    // Custom HTTP client
    whooktown.WithHTTPClient(customClient),

    // Debug mode
    whooktown.WithDebug(true),
)
```

### Individual Service URLs

For custom deployments, you can override individual service URLs:

```go
client, err := whooktown.New(
    whooktown.WithToken("your-token"),
    whooktown.WithAuthURL("https://auth.custom.example.com"),
    whooktown.WithSensorURL("https://sensors.custom.example.com"),
    whooktown.WithUIURL("https://api.custom.example.com"),
    whooktown.WithWorkflowURL("https://api.custom.example.com"),
    whooktown.WithBackofficeURL("https://admin.custom.example.com"),
    whooktown.WithSSEURL("https://ws.custom.example.com"),
    whooktown.WithSubscriptionURL("https://subscription.custom.example.com"),
)
```

## Available Clients

### Auth Client

Authentication and token management.

```go
// Sign up
token, err := client.Auth.Signup(ctx, &whooktown.SignupRequest{
    Email: "user@example.com",
    Type:  "user",
    Name:  "My App",
})

// Login
token, err := client.Auth.Login(ctx, &whooktown.LoginRequest{
    Email: "user@example.com",
    Type:  "user",
})

// Validate a token
tokenInfo, err := client.Auth.CheckToken(ctx, "the-token")

// List account tokens
tokens, err := client.Auth.ListTokens(ctx)

// Create a new token
newToken, err := client.Auth.CreateToken(ctx, &whooktown.CreateTokenRequest{
    Name: "API Token",
    Type: "sensor",
})

// Revoke a token
err = client.Auth.RevokeToken(ctx, "token-to-revoke")
```

### Sensors Client

Send sensor data.

```go
// Send sensor data
err := client.Sensors.Send(ctx, &whooktown.SensorData{
    ID:       sensorID,
    Status:   whooktown.StatusOnline,
    Activity: whooktown.ActivityFast,
})

// With building-specific fields
err := client.Sensors.Send(ctx, &whooktown.SensorData{
    ID:          dataCenterID,
    Status:      whooktown.StatusOnline,
    CPUUsage:    75,
    RAMUsage:    60,
    Temperature: 42,
})

// Send raw data
err := client.Sensors.SendRaw(ctx, map[string]interface{}{
    "id":       "sensor-uuid",
    "status":   "online",
    "activity": "normal",
    "custom":   "value",
})
```

### UI Client

Layout management.

```go
// Create a layout
layout, err := client.UI.CreateLayout(ctx, &whooktown.Layout{
    Name: "My City",
    Grid: whooktown.Grid{Width: 10, Height: 10},
    Buildings: []whooktown.Building{
        {
            ID:       uuid.Must(uuid.NewV4()),
            Type:     whooktown.BuildingWindmill,
            Location: whooktown.Location{X: 5, Y: 5},
        },
    },
})

// Delete a layout
err = client.UI.DeleteLayout(ctx, layoutID)

// Get quota information
quota, err := client.UI.GetQuota(ctx)

// List connected scenes
scenes, err := client.UI.ListScenes(ctx)
```

### Camera Client

Camera control.

```go
// Change camera mode
err := client.Camera.SetMode(ctx, layoutID, whooktown.CameraModeOrbit, 0)

// Set camera position
err := client.Camera.SetPosition(ctx, layoutID,
    &whooktown.Vector3{X: 10, Y: 5, Z: 10},  // position
    &whooktown.Vector3{X: 0, Y: -45, Z: 0},  // rotation
    60,    // FOV
    true,  // animate
    2.0,   // duration
)

// Go to a preset
err := client.Camera.GoToPreset(ctx, layoutID, presetID, true, 1.5)

// Play a camera path
err := client.Camera.PlayPath(ctx, layoutID, pathID)
err := client.Camera.PausePath(ctx, layoutID)
err := client.Camera.StopPath(ctx, layoutID)

// Manage presets
presets, err := client.Camera.ListPresets(ctx, layoutID)
preset, err := client.Camera.CreatePreset(ctx, &whooktown.CreatePresetRequest{
    LayoutID:  layoutID,
    Name:      "Main View",
    PositionX: 10, PositionY: 5, PositionZ: 10,
    RotationX: 0, RotationY: -45, RotationZ: 0,
})
err = client.Camera.DeletePreset(ctx, presetID)

// Manage paths
paths, err := client.Camera.ListPaths(ctx, layoutID)
path, err := client.Camera.CreatePath(ctx, &whooktown.CreatePathRequest{
    LayoutID: layoutID,
    Name:     "City Tour",
    Loop:     true,
})
path, err = client.Camera.AddCheckpoint(ctx, pathID, &whooktown.AddCheckpointRequest{
    GridX:       5,
    GridY:       5,
    Orientation: "N",
    Altitude:    50,
})
```

### Traffic Client

Traffic control.

```go
// Configure traffic
err := client.Traffic.SetTraffic(ctx, layoutID,
    50,                     // density (0-100)
    whooktown.SpeedNormal,  // speed
    true,                   // enabled
)

// Change only density
err := client.Traffic.SetDensity(ctx, layoutID, 75)

// Change only speed
err := client.Traffic.SetSpeed(ctx, layoutID, whooktown.SpeedFast)

// Enable/disable
err := client.Traffic.Enable(ctx, layoutID)
err := client.Traffic.Disable(ctx, layoutID)

// Get traffic states
states, err := client.Traffic.GetStates(ctx)
```

### Popup Client

Popup and label control.

```go
// Labels
err := client.Popup.ShowLabels(ctx, layoutID)
err := client.Popup.HideLabels(ctx, layoutID)

// Building details
err := client.Popup.ShowDetail(ctx, layoutID, []string{buildingID1, buildingID2})
err := client.Popup.CloseDetail(ctx, layoutID, []string{buildingID1})
err := client.Popup.CloseAllDetails(ctx, layoutID)
```

### Groups Client

Asset group management.

```go
// List groups
groups, err := client.Groups.ListGroups(ctx, layoutID)

// Create a group
group, err := client.Groups.CreateGroup(ctx, &whooktown.CreateGroupRequest{
    LayoutID: layoutID,
    Name:     "Web Servers",
})

// Add/remove members
group, err = client.Groups.AddMember(ctx, groupID, buildingID)
group, err = client.Groups.RemoveMember(ctx, groupID, buildingID)

// Delete a group
err = client.Groups.DeleteGroup(ctx, groupID)
```

### Workflow Client

Workflow management.

```go
// List workflows
workflows, err := client.Workflow.List(ctx)

// Create a workflow
workflow, err := client.Workflow.Create(ctx, &whooktown.CreateWorkflowRequest{
    Name: "My Workflow",
    Graph: map[string]*whooktown.FlowNode{
        "input1": whooktown.NewInputNode("input1", "sensor-id-1"),
        "input2": whooktown.NewInputNode("input2", "sensor-id-2"),
        "and1":   whooktown.NewAndNode("and1", []string{"input1", "input2"}),
        "output": whooktown.NewOutputNode("output", "result-sensor", []string{"and1"}),
    },
})

// Enable/disable
err = client.Workflow.Enable(ctx, workflowID)
err = client.Workflow.Disable(ctx, workflowID)

// Delete
err = client.Workflow.Delete(ctx, workflowID)

// Get available operations
operations, err := client.Workflow.GetOperations(ctx)
```

### Backoffice Client (Admin)

Administration.

```go
// Create client with admin secret
client, _ := whooktown.New(
    whooktown.WithAdminSecret("admin-secret"),
)

// Statistics
stats, err := client.Backoffice.GetStats(ctx)

// Account management
accounts, err := client.Backoffice.ListAccounts(ctx)
account, err := client.Backoffice.CreateAccount(ctx, &whooktown.CreateAccountRequest{
    Email: "new@example.com",
    Type:  "user",
})
err = client.Backoffice.LockAccount(ctx, accountID, "Lock reason")
err = client.Backoffice.UnlockAccount(ctx, accountID)
err = client.Backoffice.DeleteAccount(ctx, accountID)

// Token management
tokens, err := client.Backoffice.ListAccountTokens(ctx, accountID)
token, err := client.Backoffice.CreateAccountToken(ctx, accountID, &whooktown.CreateAccountTokenRequest{
    Type: "sensor",
    Name: "API Token",
})
err = client.Backoffice.DeleteToken(ctx, tokenString)

// Subscription management
stats, err := client.Backoffice.GetSubscriptionStats(ctx)
plans, err := client.Backoffice.ListPlans(ctx)
sub, err := client.Backoffice.UpdateAccountSubscription(ctx, accountID, "premium")
```

## Error Handling

```go
err := client.Sensors.Send(ctx, data)
if err != nil {
    if whooktown.IsUnauthorized(err) {
        // Invalid or expired token
        log.Println("Please reconnect")
    } else if whooktown.IsQuotaExceeded(err) {
        // Quota exceeded
        log.Println("Limit reached, please upgrade your plan")
    } else if whooktown.IsNotFound(err) {
        // Resource not found
        log.Println("Resource not found")
    } else {
        log.Printf("Error: %v", err)
    }
}
```

## Status Types

```go
whooktown.StatusOnline   // "online" - Service operational
whooktown.StatusOffline  // "offline" - Service stopped
whooktown.StatusWarning  // "warning" - Service degraded
whooktown.StatusCritical // "critical" - Service in error state
```

## Activity Types

```go
whooktown.ActivitySlow   // "slow"
whooktown.ActivityNormal // "normal"
whooktown.ActivityFast   // "fast"
```

## Building Types

```go
whooktown.BuildingWindmill      // Wind turbine
whooktown.BuildingDataCenter    // Data center
whooktown.BuildingArcade        // Arcade
whooktown.BuildingPyramid       // Pyramid
whooktown.BuildingTowerA        // Communication tower A
whooktown.BuildingTowerB        // Communication tower B
whooktown.BuildingSupervisor    // Monitoring station
whooktown.BuildingBank          // Bank
whooktown.BuildingMonitorTube   // Monitor tube
whooktown.BuildingBakery        // Bakery
whooktown.BuildingHouseA        // Housing A
whooktown.BuildingHouseB        // Housing B
whooktown.BuildingHouseC        // Housing C
whooktown.BuildingTree          // Holographic tree
whooktown.BuildingDisplayA      // Display screen
whooktown.BuildingTrafficLight  // Traffic light
// ... and more
```

## License

MIT License - see [LICENSE](LICENSE) for details.
