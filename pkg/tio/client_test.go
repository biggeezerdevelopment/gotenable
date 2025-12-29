package tio

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "default client",
			opts: []Option{
				WithAPIKeys("access", "secret"),
			},
			wantErr: false,
		},
		{
			name: "client with custom URL",
			opts: []Option{
				WithURL("https://custom.tenable.com"),
				WithAPIKeys("access", "secret"),
			},
			wantErr: false,
		},
		{
			name: "client with all options",
			opts: []Option{
				WithURL("https://custom.tenable.com"),
				WithAPIKeys("access", "secret"),
				WithTimeout(60 * time.Second),
				WithRetries(3),
				WithBackoff(2 * time.Second),
				WithVendor("TestVendor"),
				WithProduct("TestProduct"),
				WithBuild("1.0.0"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("New() returned nil client")
			}
		})
	}
}

func TestClientEndpoints(t *testing.T) {
	client, err := New(
		WithURL("https://cloud.tenable.com"),
		WithAPIKeys("access", "secret"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Verify all endpoints are initialized
	if client.Scans == nil {
		t.Error("Scans endpoint is nil")
	}
	if client.Assets == nil {
		t.Error("Assets endpoint is nil")
	}
	if client.Agents == nil {
		t.Error("Agents endpoint is nil")
	}
	if client.Scanners == nil {
		t.Error("Scanners endpoint is nil")
	}
	if client.Users == nil {
		t.Error("Users endpoint is nil")
	}
	if client.Exports == nil {
		t.Error("Exports endpoint is nil")
	}
	if client.Tags == nil {
		t.Error("Tags endpoint is nil")
	}
	if client.Policies == nil {
		t.Error("Policies endpoint is nil")
	}
	if client.Folders == nil {
		t.Error("Folders endpoint is nil")
	}
	if client.Plugins == nil {
		t.Error("Plugins endpoint is nil")
	}
	if client.Filters == nil {
		t.Error("Filters endpoint is nil")
	}
	if client.Credentials == nil {
		t.Error("Credentials endpoint is nil")
	}
	if client.Session == nil {
		t.Error("Session endpoint is nil")
	}
	if client.Server == nil {
		t.Error("Server endpoint is nil")
	}
	if client.Networks == nil {
		t.Error("Networks endpoint is nil")
	}
	if client.Permissions == nil {
		t.Error("Permissions endpoint is nil")
	}
	if client.AccessControl == nil {
		t.Error("AccessControl endpoint is nil")
	}
	if client.AuditLog == nil {
		t.Error("AuditLog endpoint is nil")
	}
	if client.Editor == nil {
		t.Error("Editor endpoint is nil")
	}
	if client.Workbenches == nil {
		t.Error("Workbenches endpoint is nil")
	}
}

func TestDefaultURL(t *testing.T) {
	if DefaultURL != "https://cloud.tenable.com" {
		t.Errorf("DefaultURL = %v, want https://cloud.tenable.com", DefaultURL)
	}
}

func TestEnvBase(t *testing.T) {
	if EnvBase != "TIO" {
		t.Errorf("EnvBase = %v, want TIO", EnvBase)
	}
}

