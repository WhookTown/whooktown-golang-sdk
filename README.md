# whooktown-golang-sdk

SDK Golang officiel pour [whooktown](https://whooktown.com) - la plateforme de visualisation d'infrastructure IT en ville 3D.

## Installation

```bash
go get github.com/fredericalix/whooktown-golang-sdk
```

## Demarrage rapide

```go
package main

import (
    "context"
    "log"

    whooktown "github.com/fredericalix/whooktown-golang-sdk"
    "github.com/gofrs/uuid"
)

func main() {
    // Creer le client
    client, err := whooktown.New(
        whooktown.WithBaseURL("http://localhost"),
        whooktown.WithToken("votre-token"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Envoyer des donnees de capteur
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

### Options disponibles

```go
client, err := whooktown.New(
    // URL de base pour tous les services
    whooktown.WithBaseURL("http://api.whooktown.com"),

    // Token d'authentification
    whooktown.WithToken("votre-token"),

    // Secret admin pour le backoffice
    whooktown.WithAdminSecret("votre-secret"),

    // Timeout HTTP
    whooktown.WithTimeout(30 * time.Second),

    // Configuration du retry
    whooktown.WithRetry(3, time.Second),

    // Client HTTP personnalise
    whooktown.WithHTTPClient(customClient),

    // Mode debug
    whooktown.WithDebug(true),
)
```

### URLs de services individuelles

```go
client, err := whooktown.New(
    whooktown.WithAuthURL("http://auth.example.com:8981"),
    whooktown.WithSensorURL("http://sensors.example.com:8081"),
    whooktown.WithUIURL("http://ui.example.com:8083"),
    whooktown.WithWorkflowURL("http://workflow.example.com:8084"),
    whooktown.WithBackofficeURL("http://backoffice.example.com:8086"),
)
```

## Clients disponibles

### Auth Client

Gestion de l'authentification et des tokens.

```go
// Inscription
token, err := client.Auth.Signup(ctx, &whooktown.SignupRequest{
    Email: "user@example.com",
    Type:  "user",
    Name:  "Mon App",
})

// Connexion
token, err := client.Auth.Login(ctx, &whooktown.LoginRequest{
    Email: "user@example.com",
    Type:  "user",
})

// Valider un token
tokenInfo, err := client.Auth.CheckToken(ctx, "le-token")

// Lister les tokens du compte
tokens, err := client.Auth.ListTokens(ctx)

// Creer un nouveau token
newToken, err := client.Auth.CreateToken(ctx, &whooktown.CreateTokenRequest{
    Name: "API Token",
    Type: "sensor",
})

// Revoquer un token
err = client.Auth.RevokeToken(ctx, "token-a-revoquer")
```

### Sensors Client

Envoi de donnees de capteurs.

```go
// Envoyer des donnees de capteur
err := client.Sensors.Send(ctx, &whooktown.SensorData{
    ID:       sensorID,
    Status:   whooktown.StatusOnline,
    Activity: whooktown.ActivityFast,
})

// Avec des champs specifiques au type de batiment
err := client.Sensors.Send(ctx, &whooktown.SensorData{
    ID:       dataCenterID,
    Status:   whooktown.StatusOnline,
    CPUUsage: 75,
    RAMUsage: 60,
    Temperature: 42,
})

// Envoyer des donnees brutes
err := client.Sensors.SendRaw(ctx, map[string]interface{}{
    "id":       "sensor-uuid",
    "status":   "online",
    "activity": "normal",
    "custom":   "value",
})
```

### UI Client

Gestion des layouts.

```go
// Creer un layout
layout, err := client.UI.CreateLayout(ctx, &whooktown.Layout{
    Name: "Ma Ville",
    Grid: whooktown.Grid{Width: 10, Height: 10},
    Buildings: []whooktown.Building{
        {
            ID:       uuid.Must(uuid.NewV4()),
            Type:     whooktown.BuildingWindmill,
            Location: whooktown.Location{X: 5, Y: 5},
        },
    },
})

// Supprimer un layout
err = client.UI.DeleteLayout(ctx, layoutID)

// Obtenir les quotas
quota, err := client.UI.GetQuota(ctx)

// Lister les scenes connectees
scenes, err := client.UI.ListScenes(ctx)
```

### Camera Client

Controle de la camera.

```go
// Changer le mode camera
err := client.Camera.SetMode(ctx, layoutID, whooktown.CameraModeOrbit, 0)

// Aller a une position
err := client.Camera.SetPosition(ctx, layoutID,
    &whooktown.Vector3{X: 10, Y: 5, Z: 10},  // position
    &whooktown.Vector3{X: 0, Y: -45, Z: 0},  // rotation
    60,    // FOV
    true,  // animate
    2.0,   // duration
)

// Utiliser un preset
err := client.Camera.GoToPreset(ctx, layoutID, presetID, true, 1.5)

// Jouer un chemin de camera
err := client.Camera.PlayPath(ctx, layoutID, pathID)
err := client.Camera.PausePath(ctx, layoutID)
err := client.Camera.StopPath(ctx, layoutID)

// Gerer les presets
presets, err := client.Camera.ListPresets(ctx, layoutID)
preset, err := client.Camera.CreatePreset(ctx, &whooktown.CreatePresetRequest{
    LayoutID:  layoutID,
    Name:      "Vue principale",
    PositionX: 10, PositionY: 5, PositionZ: 10,
    RotationX: 0, RotationY: -45, RotationZ: 0,
})
err = client.Camera.DeletePreset(ctx, presetID)

// Gerer les chemins
paths, err := client.Camera.ListPaths(ctx, layoutID)
path, err := client.Camera.CreatePath(ctx, &whooktown.CreatePathRequest{
    LayoutID: layoutID,
    Name:     "Tour de ville",
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

Controle du trafic.

```go
// Configurer le trafic
err := client.Traffic.SetTraffic(ctx, layoutID,
    50,                     // density (0-100)
    whooktown.SpeedNormal,  // speed
    true,                   // enabled
)

// Changer uniquement la densite
err := client.Traffic.SetDensity(ctx, layoutID, 75)

// Changer uniquement la vitesse
err := client.Traffic.SetSpeed(ctx, layoutID, whooktown.SpeedFast)

// Activer/desactiver
err := client.Traffic.Enable(ctx, layoutID)
err := client.Traffic.Disable(ctx, layoutID)

// Obtenir l'etat du trafic
states, err := client.Traffic.GetStates(ctx)
```

### Audio Client

Controle audio.

```go
// Demarrer/arreter l'audio
err := client.Audio.Play(ctx, layoutID)
err := client.Audio.Stop(ctx, layoutID)

// Changer le mood
err := client.Audio.SetMood(ctx, layoutID, whooktown.MoodTension)

// Regler le volume (0-100)
err := client.Audio.SetMusicVolume(ctx, layoutID, 80)
err := client.Audio.SetSfxVolume(ctx, layoutID, 60)

// Activer/desactiver
err := client.Audio.Enable(ctx, layoutID)
err := client.Audio.Disable(ctx, layoutID)

// Auto-mood
err := client.Audio.EnableAutoMood(ctx, layoutID)
err := client.Audio.DisableAutoMood(ctx, layoutID)
```

### Popup Client

Controle des popups et labels.

```go
// Labels
err := client.Popup.ShowLabels(ctx, layoutID)
err := client.Popup.HideLabels(ctx, layoutID)

// Details de batiments
err := client.Popup.ShowDetail(ctx, layoutID, []string{buildingID1, buildingID2})
err := client.Popup.CloseDetail(ctx, layoutID, []string{buildingID1})
err := client.Popup.CloseAllDetails(ctx, layoutID)
```

### Groups Client

Gestion des groupes d'assets.

```go
// Lister les groupes
groups, err := client.Groups.ListGroups(ctx, layoutID)

// Creer un groupe
group, err := client.Groups.CreateGroup(ctx, &whooktown.CreateGroupRequest{
    LayoutID: layoutID,
    Name:     "Serveurs Web",
})

// Ajouter/retirer des membres
group, err = client.Groups.AddMember(ctx, groupID, buildingID)
group, err = client.Groups.RemoveMember(ctx, groupID, buildingID)

// Supprimer un groupe
err = client.Groups.DeleteGroup(ctx, groupID)
```

### Workflow Client

Gestion des workflows.

```go
// Lister les workflows
workflows, err := client.Workflow.List(ctx)

// Creer un workflow
workflow, err := client.Workflow.Create(ctx, &whooktown.CreateWorkflowRequest{
    Name: "Mon Workflow",
    Graph: map[string]*whooktown.FlowNode{
        "input1": whooktown.NewInputNode("input1", "sensor-id-1"),
        "input2": whooktown.NewInputNode("input2", "sensor-id-2"),
        "and1":   whooktown.NewAndNode("and1", []string{"input1", "input2"}),
        "output": whooktown.NewOutputNode("output", "result-sensor", []string{"and1"}),
    },
})

// Activer/desactiver
err = client.Workflow.Enable(ctx, workflowID)
err = client.Workflow.Disable(ctx, workflowID)

// Supprimer
err = client.Workflow.Delete(ctx, workflowID)

// Obtenir les operations disponibles
operations, err := client.Workflow.GetOperations(ctx)
```

### Backoffice Client (Admin)

Administration.

```go
// Creer le client avec le secret admin
client, _ := whooktown.New(
    whooktown.WithAdminSecret("admin-secret"),
)

// Statistiques
stats, err := client.Backoffice.GetStats(ctx)

// Gestion des comptes
accounts, err := client.Backoffice.ListAccounts(ctx)
account, err := client.Backoffice.CreateAccount(ctx, &whooktown.CreateAccountRequest{
    Email: "new@example.com",
    Type:  "user",
})
err = client.Backoffice.LockAccount(ctx, accountID, "Raison du blocage")
err = client.Backoffice.UnlockAccount(ctx, accountID)
err = client.Backoffice.DeleteAccount(ctx, accountID)

// Gestion des tokens
tokens, err := client.Backoffice.ListAccountTokens(ctx, accountID)
token, err := client.Backoffice.CreateAccountToken(ctx, accountID, &whooktown.CreateAccountTokenRequest{
    Type: "sensor",
    Name: "API Token",
})
err = client.Backoffice.DeleteToken(ctx, tokenString)

// Gestion des subscriptions
stats, err := client.Backoffice.GetSubscriptionStats(ctx)
plans, err := client.Backoffice.ListPlans(ctx)
sub, err := client.Backoffice.UpdateAccountSubscription(ctx, accountID, "premium")
```

## Gestion des erreurs

```go
err := client.Sensors.Send(ctx, data)
if err != nil {
    if whooktown.IsUnauthorized(err) {
        // Token invalide ou expire
        log.Println("Veuillez vous reconnecter")
    } else if whooktown.IsQuotaExceeded(err) {
        // Quota depasse
        log.Println("Limite atteinte, mettez a jour votre plan")
    } else if whooktown.IsNotFound(err) {
        // Ressource non trouvee
        log.Println("Ressource introuvable")
    } else {
        log.Printf("Erreur: %v", err)
    }
}
```

## Types de statuts

```go
whooktown.StatusOnline   // "online" - Service operationnel
whooktown.StatusOffline  // "offline" - Service arrete
whooktown.StatusWarning  // "warning" - Service degrade
whooktown.StatusCritical // "critical" - Service en erreur
```

## Types d'activite

```go
whooktown.ActivitySlow   // "slow"
whooktown.ActivityNormal // "normal"
whooktown.ActivityFast   // "fast"
```

## Moods audio

```go
whooktown.MoodCalm     // "calm" - Atmosphere calme
whooktown.MoodActive   // "active" - Atmosphere active
whooktown.MoodTension  // "tension" - Atmosphere tendue
whooktown.MoodCritical // "critical" - Atmosphere critique
whooktown.MoodEpic     // "epic" - Atmosphere epique
```

## Types de batiments

```go
whooktown.BuildingWindmill      // Turbine eolienne
whooktown.BuildingDataCenter    // Centre de donnees
whooktown.BuildingArcade        // Salle de jeux
whooktown.BuildingPyramid       // Pyramide
whooktown.BuildingTowerA        // Tour de communication A
whooktown.BuildingTowerB        // Tour de communication B
whooktown.BuildingSupervisor    // Station de surveillance
whooktown.BuildingBank          // Banque
whooktown.BuildingMonitorTube   // Tube de monitoring
whooktown.BuildingBakery        // Boulangerie
whooktown.BuildingHouseA        // Habitation A
whooktown.BuildingHouseB        // Habitation B
whooktown.BuildingHouseC        // Habitation C
whooktown.BuildingTree          // Arbre holographique
whooktown.BuildingDisplayA      // Ecran d'affichage
whooktown.BuildingTrafficLight  // Feu de circulation
// ... et plus encore
```

## License

MIT License - voir [LICENSE](LICENSE) pour plus de details.
