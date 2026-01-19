// Example: Creating a workflow with whooktown SDK
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

	// Create client
	client, err := whooktown.New(
		whooktown.WithBaseURL("http://localhost"),
		whooktown.WithToken(token),
		whooktown.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// List available operations
	fmt.Println("Available workflow operations:")
	operations, err := client.Workflow.GetOperations(ctx)
	if err != nil {
		log.Fatalf("Failed to get operations: %v", err)
	}
	for name, op := range operations {
		fmt.Printf("  - %s: %s\n", name, op.Description)
	}
	fmt.Println()

	// Create a simple workflow
	// This workflow:
	// 1. Takes input from two sensors (server1, server2)
	// 2. ANDs their status together
	// 3. Outputs the result to a third sensor (cluster-status)
	fmt.Println("Creating workflow...")

	workflow, err := client.Workflow.Create(ctx, &whooktown.CreateWorkflowRequest{
		Name: "Cluster Health Monitor",
		Graph: map[string]*whooktown.FlowNode{
			// Input nodes - read status from sensors
			"server1_in": whooktown.NewInputNode("server1_in", "server-1-uuid"),
			"server2_in": whooktown.NewInputNode("server2_in", "server-2-uuid"),

			// Logic node - AND the statuses
			"cluster_and": whooktown.NewAndNode("cluster_and", []string{"server1_in", "server2_in"}),

			// Output node - write result to cluster status
			"cluster_out": whooktown.NewOutputNode("cluster_out", "cluster-status-uuid", []string{"cluster_and"}),
		},
		Enabled: true,
	})

	if err != nil {
		log.Fatalf("Failed to create workflow: %v", err)
	}

	fmt.Printf("Created workflow: %s (ID: %s)\n", workflow.Name, workflow.ID)
	fmt.Printf("  Enabled: %v\n", workflow.Enabled)

	// List all workflows
	fmt.Println("\nAll workflows:")
	workflows, err := client.Workflow.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list workflows: %v", err)
	}
	for _, w := range workflows {
		enabled := "disabled"
		if w.Enabled {
			enabled = "enabled"
		}
		fmt.Printf("  - %s (%s) [%s]\n", w.Name, w.ID, enabled)
	}

	// Example: Disable then re-enable the workflow
	fmt.Printf("\nDisabling workflow %s...\n", workflow.ID)
	err = client.Workflow.Disable(ctx, workflow.ID)
	if err != nil {
		log.Printf("Failed to disable workflow: %v", err)
	} else {
		fmt.Println("Workflow disabled")
	}

	fmt.Printf("Re-enabling workflow %s...\n", workflow.ID)
	err = client.Workflow.Enable(ctx, workflow.ID)
	if err != nil {
		log.Printf("Failed to enable workflow: %v", err)
	} else {
		fmt.Println("Workflow enabled")
	}

	fmt.Println("\nDone!")
}
