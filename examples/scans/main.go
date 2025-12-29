// Example: Working with scans
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
	client, err := tio.New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// List all scans
	scans, err := client.Scans.List(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list scans: %v", err)
	}

	fmt.Printf("Found %d scans\n", len(scans))

	if len(scans) == 0 {
		fmt.Println("No scans found. Creating a new scan...")

		// Get templates
		templates, err := client.Policies.Templates(ctx)
		if err != nil {
			log.Fatalf("Failed to get templates: %v", err)
		}

		basicTemplate := templates["basic"]
		if basicTemplate == "" {
			log.Fatal("Basic template not found")
		}

		// Create a new scan
		scan, err := client.Scans.Create(ctx, &tio.ScanCreateRequest{
			UUID: basicTemplate,
			Settings: tio.ScanSettings{
				Name:        "Example Scan",
				Description: "Created by goTenable example",
				TextTargets: "127.0.0.1",
				Enabled:     false,
			},
		})
		if err != nil {
			log.Fatalf("Failed to create scan: %v", err)
		}

		fmt.Printf("Created scan: %s (ID: %d)\n", scan.Name, scan.ID)
		return
	}

	// Get details of the first scan
	scan := scans[0]
	fmt.Printf("\nScan Details for '%s' (ID: %d):\n", scan.Name, scan.ID)

	details, err := client.Scans.Details(ctx, scan.ID)
	if err != nil {
		log.Fatalf("Failed to get scan details: %v", err)
	}

	fmt.Printf("  Status: %s\n", details.Info.Status)
	fmt.Printf("  Policy: %s\n", details.Info.Policy)
	fmt.Printf("  Targets: %s\n", details.Info.Targets)
	fmt.Printf("  Host Count: %d\n", details.Info.HostCount)

	// Get scan history
	fmt.Println("\nScan History:")
	historyIter := client.Scans.History(ctx, scan.ID, 10, 0)
	for historyIter.Next() {
		h := historyIter.Item()
		fmt.Printf("  - %s: %s (ID: %d)\n", h.UUID, h.Status, h.HistoryID)
	}
	if err := historyIter.Err(); err != nil {
		log.Printf("Error getting history: %v", err)
	}

	// Get available timezones
	timezones, err := client.Scans.Timezones(ctx)
	if err != nil {
		log.Printf("Failed to get timezones: %v", err)
	} else {
		fmt.Printf("\nAvailable timezones: %d\n", len(timezones))
	}
}

func init() {
	if os.Getenv("TIO_ACCESS_KEY") == "" || os.Getenv("TIO_SECRET_KEY") == "" {
		fmt.Println("Note: Set TIO_ACCESS_KEY and TIO_SECRET_KEY environment variables")
	}
}
