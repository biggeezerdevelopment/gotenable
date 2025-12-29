// Example: Using the exports API for bulk data retrieval
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tenable/gotenable/pkg/tio"
)

func main() {
	// Create client
	client, err := tio.New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Example 1: Export assets
	fmt.Println("Starting asset export...")
	exportUUID, err := client.Exports.AssetsExport(ctx, &tio.ExportAssetsRequest{
		ChunkSize: 1000,
	})
	if err != nil {
		log.Fatalf("Failed to start asset export: %v", err)
	}
	fmt.Printf("Asset export started: %s\n", exportUUID)

	// Wait for export to complete
	status, err := client.Exports.WaitForExport(ctx, "assets", exportUUID, 5*time.Second)
	if err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	fmt.Printf("Export complete! %d chunks available\n", len(status.ChunksAvailable))

	// Download and process chunks
	for _, chunkID := range status.ChunksAvailable {
		reader, err := client.Exports.AssetsExportChunk(ctx, exportUUID, chunkID)
		if err != nil {
			log.Printf("Failed to download chunk %d: %v", chunkID, err)
			continue
		}

		var assets []tio.ExportedAsset
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&assets); err != nil {
			log.Printf("Failed to decode chunk %d: %v", chunkID, err)
			continue
		}

		fmt.Printf("Chunk %d: %d assets\n", chunkID, len(assets))
		for _, asset := range assets[:min(5, len(assets))] {
			fmt.Printf("  - %s: %v\n", asset.ID, asset.IPv4)
		}
		if len(assets) > 5 {
			fmt.Printf("  ... and %d more\n", len(assets)-5)
		}
	}

	// Example 2: Use the iterator for simpler access
	fmt.Println("\nUsing asset iterator...")
	iterator := client.Exports.AssetsIterator(ctx, &tio.ExportAssetsRequest{
		ChunkSize: 100,
	})

	count := 0
	for iterator.Next() {
		asset := iterator.Item()
		count++
		if count <= 5 {
			fmt.Printf("Asset: %s\n", asset.ID)
		}
	}
	if err := iterator.Err(); err != nil {
		log.Fatalf("Iterator error: %v", err)
	}
	fmt.Printf("Total assets processed: %d\n", count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	if os.Getenv("TIO_ACCESS_KEY") == "" || os.Getenv("TIO_SECRET_KEY") == "" {
		fmt.Println("Note: Set TIO_ACCESS_KEY and TIO_SECRET_KEY environment variables")
	}
}
