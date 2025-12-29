// Example: Basic usage of the goTenable TenableIO client
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/biggeezerdevelopment/gotenable/pkg/tio"
)

func main() {
	// Create a new TenableIO client using environment variables
	// Set TIO_ACCESS_KEY and TIO_SECRET_KEY in your environment
	client, err := tio.New(
		tio.WithVendor("Example Company"),
		tio.WithProduct("Basic Example"),
		tio.WithBuild("1.0.0"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check server status
	status, err := client.Server.Status(ctx)
	if err != nil {
		log.Fatalf("Failed to get server status: %v", err)
	}
	fmt.Printf("Server Status: %s (code: %d)\n", status.Status, status.Code)

	// List scans
	scans, err := client.Scans.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list scans: %v", err)
	}

	fmt.Printf("\nFound %d scans:\n", len(scans))
	for _, scan := range scans {
		fmt.Printf("  - %s (ID: %d, Status: %s)\n", scan.Name, scan.ID, scan.Status)
	}

	// List folders
	folders, err := client.Folders.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list folders: %v", err)
	}

	fmt.Printf("\nFound %d folders:\n", len(folders))
	for _, folder := range folders {
		fmt.Printf("  - %s (ID: %d, Type: %s)\n", folder.Name, folder.ID, folder.Type)
	}
}

func init() {
	// Ensure required environment variables are set
	if os.Getenv("TIO_ACCESS_KEY") == "" || os.Getenv("TIO_SECRET_KEY") == "" {
		fmt.Println("Note: Set TIO_ACCESS_KEY and TIO_SECRET_KEY environment variables to run this example")
	}
}
