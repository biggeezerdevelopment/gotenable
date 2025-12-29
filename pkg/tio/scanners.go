package tio

import (
	"context"
	"fmt"
	"time"
)

// ScannersAPI handles scanner-related operations.
type ScannersAPI struct {
	client *Client
}

// Scanner represents a Nessus scanner.
type Scanner struct {
	ID                   int       `json:"id"`
	UUID                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	Type                 string    `json:"type"`
	Status               string    `json:"status"`
	ScanCount            int       `json:"scan_count"`
	EngineVersion        string    `json:"engine_version"`
	Platform             string    `json:"platform"`
	LoadedPluginSet      string    `json:"loaded_plugin_set"`
	RegistrationCode     string    `json:"registration_code,omitempty"`
	Owner                string    `json:"owner"`
	OwnerID              int       `json:"owner_id"`
	OwnerName            string    `json:"owner_name"`
	OwnerUUID            string    `json:"owner_uuid"`
	Key                  string    `json:"key,omitempty"`
	License              *License  `json:"license,omitempty"`
	NetworkName          string    `json:"network_name,omitempty"`
	Linked               int       `json:"linked"`
	Pool                 bool      `json:"pool"`
	AgentCount           int       `json:"agent_count,omitempty"`
	SupportsRemoteLogs   bool      `json:"supports_remote_logs"`
	SupportsWebapp       bool      `json:"supports_webapp"`
	Shared               int       `json:"shared"`
	UserPermissions      int       `json:"user_permissions"`
	Timestamp            int64     `json:"timestamp"`
	CreationDate         int64     `json:"creation_date"`
	LastModificationDate int64     `json:"last_modification_date"`
	LastConnect          time.Time `json:"last_connect,omitempty"`
	UIBuild              string    `json:"ui_build,omitempty"`
	UIVersion            string    `json:"ui_version,omitempty"`
}

// License represents scanner license information.
type License struct {
	Type       string `json:"type"`
	IPS        int    `json:"ips"`
	Agents     int    `json:"agents"`
	Scanners   int    `json:"scanners"`
	Expiration int64  `json:"expiration_date"`
}

// ScannerDetails contains detailed scanner information.
type ScannerDetails struct {
	Scanner
	AWSAvailabilityZone string   `json:"aws_availability_zone,omitempty"`
	AWSInstanceID       string   `json:"aws_instance_id,omitempty"`
	AWSInstanceType     string   `json:"aws_instance_type,omitempty"`
	AWSUpdateInterval   int      `json:"aws_update_interval,omitempty"`
	AWSVPCID            string   `json:"aws_vpc_id,omitempty"`
	Hostname            string   `json:"hostname,omitempty"`
	IPAddresses         []string `json:"ip_addresses,omitempty"`
	NumHosts            int      `json:"num_hosts,omitempty"`
	NumSessions         int      `json:"num_sessions,omitempty"`
	NumScans            int      `json:"num_scans,omitempty"`
	NumTCPSessions      int      `json:"num_tcp_sessions,omitempty"`
	RemoteUUID          string   `json:"remote_uuid,omitempty"`
	ScannerBootTime     int64    `json:"scanner_boottime,omitempty"`
}

// List retrieves all scanners.
func (s *ScannersAPI) List(ctx context.Context) ([]Scanner, error) {
	var result struct {
		Scanners []Scanner `json:"scanners"`
	}

	_, err := s.client.Get(ctx, "scanners", &result)
	if err != nil {
		return nil, err
	}

	return result.Scanners, nil
}

// Get retrieves a specific scanner.
func (s *ScannersAPI) Get(ctx context.Context, scannerID int) (*ScannerDetails, error) {
	var result ScannerDetails
	_, err := s.client.Get(ctx, fmt.Sprintf("scanners/%d", scannerID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a scanner.
func (s *ScannersAPI) Delete(ctx context.Context, scannerID int) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("scanners/%d", scannerID))
	return err
}

// Edit updates a scanner's settings.
func (s *ScannersAPI) Edit(ctx context.Context, scannerID int, settings map[string]interface{}) error {
	_, err := s.client.Put(ctx, fmt.Sprintf("scanners/%d", scannerID), settings, nil)
	return err
}

// ControlScans controls scan operations on a scanner.
func (s *ScannersAPI) ControlScans(ctx context.Context, scannerID int, action string) error {
	payload := map[string]string{"scan_control": action}
	_, err := s.client.Post(ctx, fmt.Sprintf("scanners/%d/scan-control", scannerID), payload, nil)
	return err
}

// ToggleLink enables or disables a scanner link.
func (s *ScannersAPI) ToggleLink(ctx context.Context, scannerID int, link bool) error {
	payload := map[string]int{}
	if link {
		payload["link"] = 1
	} else {
		payload["link"] = 0
	}
	_, err := s.client.Put(ctx, fmt.Sprintf("scanners/%d", scannerID), payload, nil)
	return err
}

// GetAWSTargets retrieves AWS targets for a scanner.
func (s *ScannersAPI) GetAWSTargets(ctx context.Context, scannerID int) ([]string, error) {
	var result struct {
		Targets []string `json:"targets"`
	}
	_, err := s.client.Get(ctx, fmt.Sprintf("scanners/%d/aws-targets", scannerID), &result)
	if err != nil {
		return nil, err
	}
	return result.Targets, nil
}

// AllowedScanners retrieves the list of scanners the user can use.
func (s *ScannersAPI) AllowedScanners(ctx context.Context) ([]Scanner, error) {
	scanners, err := s.List(ctx)
	if err != nil {
		return nil, err
	}

	var allowed []Scanner
	for _, scanner := range scanners {
		// Filter to only linked and type "local" or "pool" scanners
		if scanner.Status == "on" || scanner.Pool {
			allowed = append(allowed, scanner)
		}
	}
	return allowed, nil
}

// ScannerGroupsAPI handles scanner group operations.
type ScannerGroupsAPI struct {
	client *Client
}

// ScannerGroup represents a scanner group.
type ScannerGroup struct {
	ID                   int         `json:"id"`
	UUID                 string      `json:"uuid"`
	Name                 string      `json:"name"`
	Type                 string      `json:"type"`
	OwnerID              int         `json:"owner_id"`
	OwnerUUID            string      `json:"owner_uuid"`
	OwnerName            string      `json:"owner_name"`
	Owner                string      `json:"owner"`
	DefaultPermissions   int         `json:"default_permissions"`
	UserPermissions      int         `json:"user_permissions"`
	Shared               int         `json:"shared"`
	ScanCount            int         `json:"scan_count"`
	CreationDate         int64       `json:"creation_date"`
	LastModificationDate int64       `json:"last_modification_date"`
	Timestamp            int64       `json:"timestamp"`
	ScannerCount         int         `json:"scanner_count"`
	NetworkName          string      `json:"network_name,omitempty"`
	ScannerID            int         `json:"scanner_id,omitempty"`
	ScannerUUID          string      `json:"scanner_uuid,omitempty"`
	Token                string      `json:"token,omitempty"`
	Scanners             []Scanner   `json:"scanners,omitempty"`
	Routes               []ScanRoute `json:"routes,omitempty"`
}

// ScanRoute represents a scanner group routing rule.
type ScanRoute struct {
	Route string `json:"route"`
}

// List retrieves all scanner groups.
func (g *ScannerGroupsAPI) List(ctx context.Context) ([]ScannerGroup, error) {
	var result struct {
		ScannerPools []ScannerGroup `json:"scanner_pools"`
	}

	_, err := g.client.Get(ctx, "scanner-groups", &result)
	if err != nil {
		return nil, err
	}

	return result.ScannerPools, nil
}

// Create creates a new scanner group.
func (g *ScannerGroupsAPI) Create(ctx context.Context, name, groupType string) (*ScannerGroup, error) {
	payload := map[string]string{
		"name": name,
		"type": groupType,
	}

	var result ScannerGroup
	_, err := g.client.Post(ctx, "scanner-groups", payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a scanner group.
func (g *ScannerGroupsAPI) Delete(ctx context.Context, groupID int) error {
	_, err := g.client.Delete(ctx, fmt.Sprintf("scanner-groups/%d", groupID))
	return err
}

// Details retrieves details about a scanner group.
func (g *ScannerGroupsAPI) Details(ctx context.Context, groupID int) (*ScannerGroup, error) {
	var result ScannerGroup
	_, err := g.client.Get(ctx, fmt.Sprintf("scanner-groups/%d", groupID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates a scanner group.
func (g *ScannerGroupsAPI) Edit(ctx context.Context, groupID int, name string) error {
	payload := map[string]string{"name": name}
	_, err := g.client.Put(ctx, fmt.Sprintf("scanner-groups/%d", groupID), payload, nil)
	return err
}

// AddScanner adds a scanner to a group.
func (g *ScannerGroupsAPI) AddScanner(ctx context.Context, groupID, scannerID int) error {
	_, err := g.client.Post(ctx, fmt.Sprintf("scanner-groups/%d/scanners/%d", groupID, scannerID), nil, nil)
	return err
}

// RemoveScanner removes a scanner from a group.
func (g *ScannerGroupsAPI) RemoveScanner(ctx context.Context, groupID, scannerID int) error {
	_, err := g.client.Delete(ctx, fmt.Sprintf("scanner-groups/%d/scanners/%d", groupID, scannerID))
	return err
}

// ListScanners lists scanners in a group.
func (g *ScannerGroupsAPI) ListScanners(ctx context.Context, groupID int) ([]Scanner, error) {
	var result struct {
		Scanners []Scanner `json:"scanners"`
	}
	_, err := g.client.Get(ctx, fmt.Sprintf("scanner-groups/%d/scanners", groupID), &result)
	if err != nil {
		return nil, err
	}
	return result.Scanners, nil
}

// AddRoute adds a route to a scanner group.
func (g *ScannerGroupsAPI) AddRoute(ctx context.Context, groupID int, route string) error {
	payload := map[string]string{"route": route}
	_, err := g.client.Post(ctx, fmt.Sprintf("scanner-groups/%d/routes", groupID), payload, nil)
	return err
}

// DeleteRoute removes a route from a scanner group.
func (g *ScannerGroupsAPI) DeleteRoute(ctx context.Context, groupID int, route string) error {
	_, err := g.client.Delete(ctx, fmt.Sprintf("scanner-groups/%d/routes/%s", groupID, route))
	return err
}

// ListRoutes lists routes in a scanner group.
func (g *ScannerGroupsAPI) ListRoutes(ctx context.Context, groupID int) ([]ScanRoute, error) {
	var result struct {
		Routes []ScanRoute `json:"routes"`
	}
	_, err := g.client.Get(ctx, fmt.Sprintf("scanner-groups/%d/routes", groupID), &result)
	if err != nil {
		return nil, err
	}
	return result.Routes, nil
}
