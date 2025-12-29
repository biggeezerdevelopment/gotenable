package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// AssetsAPI handles asset-related operations.
type AssetsAPI struct {
	client *Client
}

// Asset represents an asset in Tenable.io.
type Asset struct {
	ID                    string    `json:"id"`
	HasAgent              bool      `json:"has_agent"`
	HasPluginResults      bool      `json:"has_plugin_results"`
	CreatedAt             time.Time `json:"created_at"`
	TerminatedAt          time.Time `json:"terminated_at,omitempty"`
	TerminatedBy          string    `json:"terminated_by,omitempty"`
	UpdatedAt             time.Time `json:"updated_at"`
	DeletedAt             time.Time `json:"deleted_at,omitempty"`
	DeletedBy             string    `json:"deleted_by,omitempty"`
	FirstSeen             time.Time `json:"first_seen"`
	LastSeen              time.Time `json:"last_seen"`
	FirstScanTime         time.Time `json:"first_scan_time,omitempty"`
	LastScanTime          time.Time `json:"last_scan_time,omitempty"`
	LastAuthenticatedScanDate time.Time `json:"last_authenticated_scan_date,omitempty"`
	LastLicensedScanDate  time.Time `json:"last_licensed_scan_date,omitempty"`
	LastScheduleID        string    `json:"last_schedule_id,omitempty"`
	AzureVMID             string    `json:"azure_vm_id,omitempty"`
	AzureResourceID       string    `json:"azure_resource_id,omitempty"`
	AWSEC2InstanceAMIID   string    `json:"aws_ec2_instance_ami_id,omitempty"`
	AWSEC2InstanceID      string    `json:"aws_ec2_instance_id,omitempty"`
	AgentUUID             string    `json:"agent_uuid,omitempty"`
	BiosUUID              string    `json:"bios_uuid,omitempty"`
	NetworkID             string    `json:"network_id,omitempty"`
	NetworkName           string    `json:"network_name,omitempty"`
	AWSEC2InstanceGroupName string  `json:"aws_ec2_instance_group_name,omitempty"`
	AWSEC2InstanceStateName string  `json:"aws_ec2_instance_state_name,omitempty"`
	AWSEC2InstanceType    string    `json:"aws_ec2_instance_type,omitempty"`
	AWSOwnerID            string    `json:"aws_owner_id,omitempty"`
	AWSAvailabilityZone   string    `json:"aws_availability_zone,omitempty"`
	AWSEC2ProductCode     string    `json:"aws_ec2_product_code,omitempty"`
	AWSSubnetID           string    `json:"aws_subnet_id,omitempty"`
	AWSVPCID              string    `json:"aws_vpc_id,omitempty"`
	AWSRegion             string    `json:"aws_region,omitempty"`
	MacAddress            []string  `json:"mac_address,omitempty"`
	McafeeEPOGUID         string    `json:"mcafee_epo_guid,omitempty"`
	McafeeEPOAgentGUID    string    `json:"mcafee_epo_agent_guid,omitempty"`
	NetbiosName           []string  `json:"netbios_name,omitempty"`
	OperatingSystem       []string  `json:"operating_system,omitempty"`
	SystemType            []string  `json:"system_type,omitempty"`
	TenableUUID           string    `json:"tenable_uuid,omitempty"`
	Hostname              []string  `json:"hostname,omitempty"`
	AgentName             []string  `json:"agent_name,omitempty"`
	FQDN                  []string  `json:"fqdn,omitempty"`
	IPv4                  []string  `json:"ipv4,omitempty"`
	IPv6                  []string  `json:"ipv6,omitempty"`
	SSHFingerprint        []string  `json:"ssh_fingerprint,omitempty"`
	QualysAssetID         string    `json:"qualys_asset_id,omitempty"`
	QualysHostID          string    `json:"qualys_host_id,omitempty"`
	ServiceNowSystemID    string    `json:"servicenow_sysid,omitempty"`
	InstalledSoftware     []string  `json:"installed_software,omitempty"`
	Sources               []AssetSource `json:"sources,omitempty"`
	Tags                  []AssetTag    `json:"tags,omitempty"`
	AcrScore              int       `json:"acr_score,omitempty"`
	AcrDrivers            []ACRDriver `json:"acr_drivers,omitempty"`
	ExposureScore         int       `json:"exposure_score,omitempty"`
	ScanFrequency         int       `json:"scan_frequency,omitempty"`
}

// AssetSource represents the source of an asset.
type AssetSource struct {
	Name      string    `json:"name"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// AssetTag represents a tag on an asset.
type AssetTag struct {
	TagUUID      string    `json:"tag_uuid"`
	TagKey       string    `json:"tag_key"`
	TagValue     string    `json:"tag_value"`
	AddedBy      string    `json:"added_by"`
	AddedAt      time.Time `json:"added_at"`
}

// ACRDriver represents an ACR driver.
type ACRDriver struct {
	DriverName  string `json:"driver_name"`
	DriverValue string `json:"driver_value"`
}

// AssetListOptions contains options for listing assets.
type AssetListOptions struct {
	DateRange int    // Number of days to look back
	Filter    string // Filter expression
}

// List retrieves a list of assets.
func (a *AssetsAPI) List(ctx context.Context, opts *AssetListOptions) *base.Iterator[Asset] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}
		if opts != nil {
			if opts.DateRange > 0 {
				params["date_range"] = strconv.Itoa(opts.DateRange)
			}
			if opts.Filter != "" {
				params["filter"] = opts.Filter
			}
		}

		var result struct {
			Assets []Asset `json:"assets"`
			Total  int     `json:"total"`
		}

		_, err := a.client.GetWithParams(ctx, "assets", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Assets)
		return data, &base.PaginationInfo{
			Total:  result.Total,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]Asset, error) {
		var items []Asset
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// Get retrieves a specific asset by UUID.
func (a *AssetsAPI) Get(ctx context.Context, assetUUID string) (*Asset, error) {
	var result Asset
	_, err := a.client.Get(ctx, fmt.Sprintf("assets/%s", assetUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an asset.
func (a *AssetsAPI) Delete(ctx context.Context, assetUUID string) error {
	_, err := a.client.Delete(ctx, fmt.Sprintf("assets/%s", assetUUID))
	return err
}

// AssetInfo contains summary information about an asset.
type AssetInfo struct {
	ID                    string    `json:"id"`
	HasAgent              bool      `json:"has_agent"`
	HasPluginResults      bool      `json:"has_plugin_results"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	FirstSeen             time.Time `json:"first_seen"`
	LastSeen              time.Time `json:"last_seen"`
	IPv4                  []string  `json:"ipv4,omitempty"`
	IPv6                  []string  `json:"ipv6,omitempty"`
	FQDN                  []string  `json:"fqdn,omitempty"`
	Hostname              []string  `json:"hostname,omitempty"`
	NetbiosName           []string  `json:"netbios_name,omitempty"`
	OperatingSystem       []string  `json:"operating_system,omitempty"`
	MacAddress            []string  `json:"mac_address,omitempty"`
	AgentName             []string  `json:"agent_name,omitempty"`
}

// Info retrieves summary information about an asset.
func (a *AssetsAPI) Info(ctx context.Context, assetUUID string) (*AssetInfo, error) {
	var result AssetInfo
	_, err := a.client.Get(ctx, fmt.Sprintf("assets/%s/info", assetUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AssetVulnerability represents a vulnerability on an asset.
type AssetVulnerability struct {
	PluginID           int       `json:"plugin_id"`
	PluginName         string    `json:"plugin_name"`
	PluginFamily       string    `json:"plugin_family"`
	Severity           int       `json:"severity"`
	SeverityIndex      int       `json:"severity_index"`
	VPRScore           float64   `json:"vpr_score,omitempty"`
	State              string    `json:"state"`
	Count              int       `json:"count"`
	FirstFound         time.Time `json:"first_found"`
	LastFound          time.Time `json:"last_found"`
	LastFixed          time.Time `json:"last_fixed,omitempty"`
	AcceptedCount      int       `json:"accepted_count,omitempty"`
	RecastedCount      int       `json:"recasted_count,omitempty"`
	CountsTotal        int       `json:"counts_by_severity_total,omitempty"`
	CVSSBaseScore      float64   `json:"cvss_base_score,omitempty"`
	CVSSTemporalScore  float64   `json:"cvss_temporal_score,omitempty"`
}

// Vulnerabilities retrieves vulnerabilities for an asset.
func (a *AssetsAPI) Vulnerabilities(ctx context.Context, assetUUID string) *base.Iterator[AssetVulnerability] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			Vulnerabilities []AssetVulnerability `json:"vulnerabilities"`
			Total           int                  `json:"total"`
		}

		_, err := a.client.GetWithParams(ctx, fmt.Sprintf("assets/%s/vulnerabilities", assetUUID), params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Vulnerabilities)
		return data, &base.PaginationInfo{
			Total:  result.Total,
			Limit:  limit,
			Offset: offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]AssetVulnerability, error) {
		var items []AssetVulnerability
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// BulkDeleteRequest represents a request to bulk delete assets.
type BulkDeleteRequest struct {
	Query   *BulkDeleteQuery `json:"query,omitempty"`
	HardDelete bool          `json:"hard_delete,omitempty"`
}

// BulkDeleteQuery represents the query for bulk delete.
type BulkDeleteQuery struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// BulkDelete deletes multiple assets based on a query.
func (a *AssetsAPI) BulkDelete(ctx context.Context, req *BulkDeleteRequest) error {
	_, err := a.client.Post(ctx, "assets/bulk-jobs/delete", req, nil)
	return err
}

// AssignTags assigns tags to assets.
func (a *AssetsAPI) AssignTags(ctx context.Context, assetUUIDs []string, tagUUIDs []string) error {
	payload := map[string]interface{}{
		"action": "add",
		"assets": assetUUIDs,
		"tags":   tagUUIDs,
	}
	_, err := a.client.Post(ctx, "tags/assets/assignments", payload, nil)
	return err
}

// UnassignTags removes tags from assets.
func (a *AssetsAPI) UnassignTags(ctx context.Context, assetUUIDs []string, tagUUIDs []string) error {
	payload := map[string]interface{}{
		"action": "remove",
		"assets": assetUUIDs,
		"tags":   tagUUIDs,
	}
	_, err := a.client.Post(ctx, "tags/assets/assignments", payload, nil)
	return err
}

// MoveToNetwork moves assets to a different network.
func (a *AssetsAPI) MoveToNetwork(ctx context.Context, sourceNetworkID, destNetworkID string, assetUUIDs []string) error {
	payload := map[string]interface{}{
		"source":      sourceNetworkID,
		"destination": destNetworkID,
		"targets":     assetUUIDs,
	}
	_, err := a.client.Post(ctx, "assets/bulk-jobs/move-to-network", payload, nil)
	return err
}

