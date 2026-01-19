// Example: Basic usage of whooktown SDK
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
)

func main() {
	// Get token from environment
	token := os.Getenv("WHOOKTOWN_TOKEN")
	if token == "" {
		log.Fatal("WHOOKTOWN_TOKEN environment variable is required")
	}

	// Create client (uses PROD by default, set WHOOKTOWN_ENV=DEV for development)
	client, err := whooktown.New(
		whooktown.WithToken(token),
		whooktown.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Check token validity
	tokenInfo, err := client.Auth.CheckToken(ctx, token)
	if err != nil {
		log.Fatalf("Invalid token: %v", err)
	}
	fmt.Printf("Authenticated as account: %s\n", tokenInfo.AccountID)

	// Get quota info
	quota, err := client.UI.GetQuota(ctx)
	if err != nil {
		log.Fatalf("Failed to get quota: %v", err)
	}
	fmt.Printf("Plan: %s, Layouts: %d/%d\n", quota.Plan, quota.Layouts.Used, quota.Layouts.Max)

	// List connected scenes
	scenes, err := client.UI.ListScenes(ctx)
	if err != nil {
		log.Fatalf("Failed to list scenes: %v", err)
	}
	fmt.Printf("Connected scenes: %d\n", len(scenes))
	for _, scene := range scenes {
		fmt.Printf("  - Scene %s (layout: %s)\n", scene.SceneID, scene.LayoutID)
	}

	fmt.Println("Done!")
}
