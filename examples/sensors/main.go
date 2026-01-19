// Example: Sending sensor data to whooktown
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	whooktown "github.com/fredericalix/whooktown-golang-sdk"
	"github.com/gofrs/uuid"
)

func main() {
	// Get configuration from environment
	token := os.Getenv("WHOOKTOWN_TOKEN")
	if token == "" {
		log.Fatal("WHOOKTOWN_TOKEN environment variable is required")
	}

	sensorIDStr := os.Getenv("WHOOKTOWN_SENSOR_ID")
	if sensorIDStr == "" {
		log.Fatal("WHOOKTOWN_SENSOR_ID environment variable is required")
	}

	sensorID, err := uuid.FromString(sensorIDStr)
	if err != nil {
		log.Fatalf("Invalid sensor ID: %v", err)
	}

	// Create client
	client, err := whooktown.New(
		whooktown.WithBaseURL("http://localhost"),
		whooktown.WithToken(token),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Send sensor data in a loop
	fmt.Println("Sending sensor data every 5 seconds... (Ctrl+C to stop)")

	for {
		// Generate random metrics
		cpuUsage := rand.Intn(100)
		ramUsage := rand.Intn(100)
		temp := 30 + rand.Intn(40)

		// Determine status based on metrics
		status := whooktown.StatusOnline
		activity := whooktown.ActivityNormal

		if cpuUsage > 90 || ramUsage > 90 || temp > 60 {
			status = whooktown.StatusCritical
			activity = whooktown.ActivityFast
		} else if cpuUsage > 70 || ramUsage > 70 || temp > 50 {
			status = whooktown.StatusWarning
			activity = whooktown.ActivityFast
		} else if cpuUsage < 30 && ramUsage < 30 {
			activity = whooktown.ActivitySlow
		}

		// Send sensor data
		err := client.Sensors.Send(ctx, &whooktown.SensorData{
			ID:          sensorID,
			Status:      status,
			Activity:    activity,
			CPUUsage:    cpuUsage,
			RAMUsage:    ramUsage,
			Temperature: temp,
		})

		if err != nil {
			if whooktown.IsUnauthorized(err) {
				log.Fatal("Token expired or invalid")
			}
			log.Printf("Failed to send sensor data: %v", err)
		} else {
			fmt.Printf("[%s] Sent: CPU=%d%%, RAM=%d%%, Temp=%d°C, Status=%s\n",
				time.Now().Format("15:04:05"), cpuUsage, ramUsage, temp, status)
		}

		time.Sleep(5 * time.Second)
	}
}
