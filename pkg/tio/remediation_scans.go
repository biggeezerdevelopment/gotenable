package tio

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/biggeezerdevelopment/gotenable/pkg/base"
)

// RemediationScansAPI handles remediation scan operations.
type RemediationScansAPI struct {
	client *Client
}

// RemediationScan represents a remediation scan.
type RemediationScan struct {
	ID                   int       `json:"id"`
	UUID                 string    `json:"uuid"`
	Name                 string    `json:"name"`
	Description          string    `json:"description,omitempty"`
	Status               string    `json:"status"`
	Schedule             string    `json:"schedule,omitempty"`
	ScannerID            int       `json:"scanner_id"`
	ScannerName          string    `json:"scanner_name,omitempty"`
	PolicyID             int       `json:"policy_id,omitempty"`
	Remediation          RemediationInfo `json:"remediation,omitempty"`
	CreationDate         time.Time `json:"creation_date"`
	LastModificationDate time.Time `json:"last_modification_date"`
	Owner                string    `json:"owner"`
	OwnerID              int       `json:"owner_id"`
}

// RemediationInfo contains remediation-specific information.
type RemediationInfo struct {
	ID             int      `json:"id"`
	PluginID       int      `json:"plugin_id"`
	AssetCount     int      `json:"asset_count"`
	VulnCount      int      `json:"vuln_count"`
	Hosts          []string `json:"hosts,omitempty"`
}

// RemediationScanCreateRequest represents a request to create a remediation scan.
type RemediationScanCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ScannerID   int    `json:"scanner_id"`
	PluginID    int    `json:"plugin_id"`
	Assets      []string `json:"assets,omitempty"`
}

// List retrieves all remediation scans.
func (r *RemediationScansAPI) List(ctx context.Context) *base.Iterator[RemediationScan] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			Scans      []RemediationScan `json:"scans"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := r.client.GetWithParams(ctx, "remediation-scans", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Scans)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]RemediationScan, error) {
		var items []RemediationScan
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// Create creates a new remediation scan.
func (r *RemediationScansAPI) Create(ctx context.Context, req *RemediationScanCreateRequest) (*RemediationScan, error) {
	var result RemediationScan
	_, err := r.client.Post(ctx, "remediation-scans", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

