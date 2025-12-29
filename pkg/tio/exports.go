package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// ExportsAPI handles export operations for assets and vulnerabilities.
type ExportsAPI struct {
	client *Client
}

// ExportAssetsRequest represents a request to export assets.
type ExportAssetsRequest struct {
	ChunkSize        int                    `json:"chunk_size,omitempty"`
	Filters          map[string]interface{} `json:"filters,omitempty"`
	IncludeUnlicensed bool                  `json:"include_unlicensed,omitempty"`
}

// ExportVulnsRequest represents a request to export vulnerabilities.
type ExportVulnsRequest struct {
	NumAssets         int                    `json:"num_assets,omitempty"`
	Filters           map[string]interface{} `json:"filters,omitempty"`
	IncludeUnlicensed bool                   `json:"include_unlicensed,omitempty"`
}

// ExportComplianceRequest represents a request to export compliance data.
type ExportComplianceRequest struct {
	NumAssets int                    `json:"num_assets,omitempty"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
}

// ExportResponse represents the response from initiating an export.
type ExportResponse struct {
	ExportUUID string `json:"export_uuid"`
}

// ExportStatus represents the status of an export.
type ExportStatus struct {
	Status           string   `json:"status"`
	ChunksAvailable  []int    `json:"chunks_available"`
	ChunksCancelled  []int    `json:"chunks_cancelled,omitempty"`
	ChunksFailed     []int    `json:"chunks_failed,omitempty"`
	TotalChunks      int      `json:"total_chunks,omitempty"`
	ChunksFinished   int      `json:"finished_chunks,omitempty"`
	NumAssetsPerChunk int     `json:"num_assets_per_chunk,omitempty"`
	Created          int64    `json:"created,omitempty"`
	UUID             string   `json:"uuid,omitempty"`
	EmptyChunksCount int      `json:"empty_chunks_count,omitempty"`
	Filters          map[string]interface{} `json:"filters,omitempty"`
}

// AssetsExport initiates an assets export.
func (e *ExportsAPI) AssetsExport(ctx context.Context, req *ExportAssetsRequest) (string, error) {
	if req == nil {
		req = &ExportAssetsRequest{}
	}
	if req.ChunkSize == 0 {
		req.ChunkSize = 1000
	}

	var result ExportResponse
	_, err := e.client.Post(ctx, "assets/export", req, &result)
	if err != nil {
		return "", err
	}
	return result.ExportUUID, nil
}

// VulnsExport initiates a vulnerabilities export.
func (e *ExportsAPI) VulnsExport(ctx context.Context, req *ExportVulnsRequest) (string, error) {
	if req == nil {
		req = &ExportVulnsRequest{}
	}
	if req.NumAssets == 0 {
		req.NumAssets = 500
	}

	var result ExportResponse
	_, err := e.client.Post(ctx, "vulns/export", req, &result)
	if err != nil {
		return "", err
	}
	return result.ExportUUID, nil
}

// ComplianceExport initiates a compliance export.
func (e *ExportsAPI) ComplianceExport(ctx context.Context, req *ExportComplianceRequest) (string, error) {
	if req == nil {
		req = &ExportComplianceRequest{}
	}
	if req.NumAssets == 0 {
		req.NumAssets = 500
	}

	var result ExportResponse
	_, err := e.client.Post(ctx, "compliance/export", req, &result)
	if err != nil {
		return "", err
	}
	return result.ExportUUID, nil
}

// AssetsExportStatus retrieves the status of an assets export.
func (e *ExportsAPI) AssetsExportStatus(ctx context.Context, exportUUID string) (*ExportStatus, error) {
	var result ExportStatus
	_, err := e.client.Get(ctx, fmt.Sprintf("assets/export/%s/status", exportUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// VulnsExportStatus retrieves the status of a vulnerabilities export.
func (e *ExportsAPI) VulnsExportStatus(ctx context.Context, exportUUID string) (*ExportStatus, error) {
	var result ExportStatus
	_, err := e.client.Get(ctx, fmt.Sprintf("vulns/export/%s/status", exportUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ComplianceExportStatus retrieves the status of a compliance export.
func (e *ExportsAPI) ComplianceExportStatus(ctx context.Context, exportUUID string) (*ExportStatus, error) {
	var result ExportStatus
	_, err := e.client.Get(ctx, fmt.Sprintf("compliance/export/%s/status", exportUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AssetsExportChunk downloads an assets export chunk.
func (e *ExportsAPI) AssetsExportChunk(ctx context.Context, exportUUID string, chunkID int) (io.Reader, error) {
	data, err := e.client.Download(ctx, fmt.Sprintf("assets/export/%s/chunks/%d", exportUUID, chunkID))
	if err != nil {
		return nil, err
	}
	return &bytesReader{data: data}, nil
}

// VulnsExportChunk downloads a vulnerabilities export chunk.
func (e *ExportsAPI) VulnsExportChunk(ctx context.Context, exportUUID string, chunkID int) (io.Reader, error) {
	data, err := e.client.Download(ctx, fmt.Sprintf("vulns/export/%s/chunks/%d", exportUUID, chunkID))
	if err != nil {
		return nil, err
	}
	return &bytesReader{data: data}, nil
}

// ComplianceExportChunk downloads a compliance export chunk.
func (e *ExportsAPI) ComplianceExportChunk(ctx context.Context, exportUUID string, chunkID int) (io.Reader, error) {
	data, err := e.client.Download(ctx, fmt.Sprintf("compliance/export/%s/chunks/%d", exportUUID, chunkID))
	if err != nil {
		return nil, err
	}
	return &bytesReader{data: data}, nil
}

// CancelAssetsExport cancels an assets export.
func (e *ExportsAPI) CancelAssetsExport(ctx context.Context, exportUUID string) error {
	_, err := e.client.Post(ctx, fmt.Sprintf("assets/export/%s/cancel", exportUUID), nil, nil)
	return err
}

// CancelVulnsExport cancels a vulnerabilities export.
func (e *ExportsAPI) CancelVulnsExport(ctx context.Context, exportUUID string) error {
	_, err := e.client.Post(ctx, fmt.Sprintf("vulns/export/%s/cancel", exportUUID), nil, nil)
	return err
}

// CancelComplianceExport cancels a compliance export.
func (e *ExportsAPI) CancelComplianceExport(ctx context.Context, exportUUID string) error {
	_, err := e.client.Post(ctx, fmt.Sprintf("compliance/export/%s/cancel", exportUUID), nil, nil)
	return err
}

// ExportedAsset represents an exported asset.
type ExportedAsset struct {
	ID                     string    `json:"id"`
	HasAgent               bool      `json:"has_agent"`
	HasPluginResults       bool      `json:"has_plugin_results"`
	CreatedAt              time.Time `json:"created_at"`
	TerminatedAt           time.Time `json:"terminated_at,omitempty"`
	TerminatedBy           string    `json:"terminated_by,omitempty"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              time.Time `json:"deleted_at,omitempty"`
	DeletedBy              string    `json:"deleted_by,omitempty"`
	FirstSeen              time.Time `json:"first_seen"`
	LastSeen               time.Time `json:"last_seen"`
	FirstScanTime          time.Time `json:"first_scan_time,omitempty"`
	LastScanTime           time.Time `json:"last_scan_time,omitempty"`
	LastAuthenticatedScanDate time.Time `json:"last_authenticated_scan_date,omitempty"`
	LastLicensedScanDate   time.Time `json:"last_licensed_scan_date,omitempty"`
	LastScheduleID         string    `json:"last_schedule_id,omitempty"`
	Sources                []AssetSource `json:"sources,omitempty"`
	Tags                   []AssetTag    `json:"tags,omitempty"`
	NetworkInterfaces      []NetworkInterface `json:"network_interfaces,omitempty"`
	FQDN                   []string  `json:"fqdn,omitempty"`
	IPv4                   []string  `json:"ipv4,omitempty"`
	IPv6                   []string  `json:"ipv6,omitempty"`
	MacAddress             []string  `json:"mac_address,omitempty"`
	NetbiosName            []string  `json:"netbios_name,omitempty"`
	OperatingSystem        []string  `json:"operating_system,omitempty"`
	AgentName              []string  `json:"agent_name,omitempty"`
	Hostname               []string  `json:"hostname,omitempty"`
}

// NetworkInterface represents a network interface on an asset.
type NetworkInterface struct {
	Name       string   `json:"name"`
	Virtual    bool     `json:"virtual,omitempty"`
	Aliased    bool     `json:"aliased,omitempty"`
	FQDN       []string `json:"fqdn,omitempty"`
	IPv4       []string `json:"ipv4,omitempty"`
	IPv6       []string `json:"ipv6,omitempty"`
	MacAddress []string `json:"mac_address,omitempty"`
}

// ExportedVuln represents an exported vulnerability.
type ExportedVuln struct {
	Asset                Asset                `json:"asset"`
	Output               string               `json:"output"`
	Plugin               VulnPlugin           `json:"plugin"`
	Port                 VulnPort             `json:"port"`
	Scan                 VulnScan             `json:"scan"`
	Severity             string               `json:"severity"`
	SeverityID           int                  `json:"severity_id"`
	SeverityDefaultID    int                  `json:"severity_default_id"`
	SeverityModificationType string           `json:"severity_modification_type,omitempty"`
	FirstFound           time.Time            `json:"first_found"`
	LastFound            time.Time            `json:"last_found"`
	LastFixed            time.Time            `json:"last_fixed,omitempty"`
	State                string               `json:"state"`
	IndexedAt            time.Time            `json:"indexed_at,omitempty"`
}

// VulnPlugin contains plugin information for a vulnerability.
type VulnPlugin struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Family                string   `json:"family"`
	FamilyID              int      `json:"family_id"`
	HasPatch              bool     `json:"has_patch"`
	Type                  string   `json:"type"`
	Version               string   `json:"version"`
	RiskFactor            string   `json:"risk_factor"`
	Solution              string   `json:"solution,omitempty"`
	Synopsis              string   `json:"synopsis,omitempty"`
	SeeAlso               []string `json:"see_also,omitempty"`
	CVE                   []string `json:"cve,omitempty"`
	BID                   []int    `json:"bid,omitempty"`
	XREF                  []string `json:"xref,omitempty"`
	VPRScore              float64  `json:"vpr_score,omitempty"`
	CVSSBaseScore         float64  `json:"cvss_base_score,omitempty"`
	CVSSTemporalScore     float64  `json:"cvss_temporal_score,omitempty"`
	CVSS3BaseScore        float64  `json:"cvss3_base_score,omitempty"`
	CVSS3TemporalScore    float64  `json:"cvss3_temporal_score,omitempty"`
}

// VulnPort contains port information for a vulnerability.
type VulnPort struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service,omitempty"`
}

// VulnScan contains scan information for a vulnerability.
type VulnScan struct {
	CompletedAt  time.Time `json:"completed_at"`
	ScheduleUUID string    `json:"schedule_uuid"`
	StartedAt    time.Time `json:"started_at"`
	UUID         string    `json:"uuid"`
}

// WaitForExport waits for an export to complete and returns all data.
func (e *ExportsAPI) WaitForExport(ctx context.Context, exportType, exportUUID string, pollInterval time.Duration) (*ExportStatus, error) {
	if pollInterval == 0 {
		pollInterval = 5 * time.Second
	}

	for {
		var status *ExportStatus
		var err error

		switch exportType {
		case "assets":
			status, err = e.AssetsExportStatus(ctx, exportUUID)
		case "vulns":
			status, err = e.VulnsExportStatus(ctx, exportUUID)
		case "compliance":
			status, err = e.ComplianceExportStatus(ctx, exportUUID)
		default:
			return nil, &base.ValidationError{Field: "exportType", Message: "invalid export type"}
		}

		if err != nil {
			return nil, err
		}

		switch status.Status {
		case "FINISHED":
			return status, nil
		case "CANCELLED":
			return nil, &base.ExportError{ExportType: exportType, UUID: exportUUID, Message: "export cancelled"}
		case "ERROR":
			return nil, &base.ExportError{ExportType: exportType, UUID: exportUUID, Message: "export failed"}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

// AssetsIterator returns an iterator over exported assets.
func (e *ExportsAPI) AssetsIterator(ctx context.Context, req *ExportAssetsRequest) *base.Iterator[ExportedAsset] {
	var exportUUID string
	var status *ExportStatus
	var currentChunk int
	var chunkData []ExportedAsset

	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		// Initialize export if needed
		if exportUUID == "" {
			uuid, err := e.AssetsExport(ctx, req)
			if err != nil {
				return nil, nil, err
			}
			exportUUID = uuid
		}

		// Wait for export and get status
		if status == nil {
			s, err := e.WaitForExport(ctx, "assets", exportUUID, 5*time.Second)
			if err != nil {
				return nil, nil, err
			}
			status = s
		}

		// Check if we have more chunks
		if currentChunk >= len(status.ChunksAvailable) {
			return json.RawMessage("[]"), &base.PaginationInfo{Total: 0}, nil
		}

		// Download chunk
		chunkID := status.ChunksAvailable[currentChunk]
		reader, err := e.AssetsExportChunk(ctx, exportUUID, chunkID)
		if err != nil {
			return nil, nil, err
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}

		if err := json.Unmarshal(data, &chunkData); err != nil {
			return nil, nil, err
		}

		currentChunk++
		return data, &base.PaginationInfo{
			Total: len(status.ChunksAvailable),
		}, nil
	}

	transformer := func(data json.RawMessage) ([]ExportedAsset, error) {
		var items []ExportedAsset
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// VulnsIterator returns an iterator over exported vulnerabilities.
func (e *ExportsAPI) VulnsIterator(ctx context.Context, req *ExportVulnsRequest) *base.Iterator[ExportedVuln] {
	var exportUUID string
	var status *ExportStatus
	var currentChunk int

	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		// Initialize export if needed
		if exportUUID == "" {
			uuid, err := e.VulnsExport(ctx, req)
			if err != nil {
				return nil, nil, err
			}
			exportUUID = uuid
		}

		// Wait for export and get status
		if status == nil {
			s, err := e.WaitForExport(ctx, "vulns", exportUUID, 5*time.Second)
			if err != nil {
				return nil, nil, err
			}
			status = s
		}

		// Check if we have more chunks
		if currentChunk >= len(status.ChunksAvailable) {
			return json.RawMessage("[]"), &base.PaginationInfo{Total: 0}, nil
		}

		// Download chunk
		chunkID := status.ChunksAvailable[currentChunk]
		reader, err := e.VulnsExportChunk(ctx, exportUUID, chunkID)
		if err != nil {
			return nil, nil, err
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}

		currentChunk++
		return data, &base.PaginationInfo{
			Total: len(status.ChunksAvailable),
		}, nil
	}

	transformer := func(data json.RawMessage) ([]ExportedVuln, error) {
		var items []ExportedVuln
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// ListExports lists all exports of a given type.
func (e *ExportsAPI) ListExports(ctx context.Context, exportType string) ([]ExportStatus, error) {
	params := map[string]string{
		"size": strconv.Itoa(1000),
	}

	var result struct {
		Exports []ExportStatus `json:"exports"`
	}

	_, err := e.client.GetWithParams(ctx, fmt.Sprintf("%s/export/status", exportType), params, &result)
	if err != nil {
		return nil, err
	}

	return result.Exports, nil
}

