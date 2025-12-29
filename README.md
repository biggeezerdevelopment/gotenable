# goTenable

A Go SDK for Tenable APIs, providing idiomatic Go interfaces to Tenable's security products.

## Installation

```bash
go get github.com/biggeezerdevelopment/gotenable
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/biggeezerdevelopment/gotenable/pkg/tio"
)

func main() {
    // Create a new TenableIO client
    client, err := tio.New(
        tio.WithAPIKeys("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // List scans
    scans, err := client.Scans.List(nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, scan := range scans {
        fmt.Printf("Scan: %s (ID: %d)\n", scan.Name, scan.ID)
    }
}
```

## Environment Variables

The SDK supports configuration via environment variables:

- `TIO_ACCESS_KEY`: API access key for Tenable.io
- `TIO_SECRET_KEY`: API secret key for Tenable.io
- `TIO_URL`: Base URL (defaults to `https://cloud.tenable.com`)

## Features

- Full support for Tenable.io (Vulnerability Management) APIs
- Automatic retry with exponential backoff
- Pagination iterators for large result sets
- Context support for cancellation and timeouts
- Comprehensive error handling

## Supported APIs

### Tenable.io (Vulnerability Management)

- Access Control
- Agents & Agent Groups
- Assets
- Audit Log
- Credentials
- Exclusions
- Exports
- Files
- Filters
- Folders
- Groups
- Networks
- Permissions
- Plugins
- Policies
- Scanners & Scanner Groups
- Scans
- Server
- Session
- Tags
- Users
- Workbenches

## License

MIT License - see LICENSE file for details.
