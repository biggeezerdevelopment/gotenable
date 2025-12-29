package tio

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"
)

// WorkbenchesAPI handles workbench operations.
// Note: This API is deprecated in favor of the Exports API.
type WorkbenchesAPI struct {
	client *Client
}

// WorkbenchVuln represents a vulnerability in workbench.
type WorkbenchVuln struct {
	PluginID     int    `json:"plugin_id"`
	PluginName   string `json:"plugin_name"`
	PluginFamily string `json:"plugin_family"`
	Count        int    `json:"count"`
	VulnIndex    int    `json:"vuln_index"`
	Severity     int    `json:"severity"`
}

// WorkbenchAsset represents an asset in workbench.
type WorkbenchAsset struct {
	ID              string   `json:"id"`
	HasAgent        bool     `json:"has_agent"`
	LastSeen        string   `json:"last_seen"`
	Sources         []string `json:"sources,omitempty"`
	IPv4            []string `json:"ipv4,omitempty"`
	IPv6            []string `json:"ipv6,omitempty"`
	FQDN            []string `json:"fqdn,omitempty"`
	NetbiosName     []string `json:"netbios_name,omitempty"`
	OperatingSystem []string `json:"operating_system,omitempty"`
	AgentName       []string `json:"agent_name,omitempty"`
	Severities      struct {
		Total    int `json:"total"`
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
		Info     int `json:"info"`
	} `json:"severities"`
}

// WorkbenchOptions represents options for workbench queries.
type WorkbenchOptions struct {
	DateRange  int
	FilterType string
	Filters    []WorkbenchFilter
}

// WorkbenchFilter represents a filter for workbench queries.
type WorkbenchFilter struct {
	Name     string
	Operator string
	Value    string
}

// Assets retrieves assets from the workbench.
func (w *WorkbenchesAPI) Assets(ctx context.Context, opts *WorkbenchOptions) ([]WorkbenchAsset, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.DateRange > 0 {
			params["date_range"] = strconv.Itoa(opts.DateRange)
		}
		if opts.FilterType != "" {
			params["filter.search_type"] = opts.FilterType
		}
		for i, f := range opts.Filters {
			params[fmt.Sprintf("filter.%d.filter", i)] = f.Name
			params[fmt.Sprintf("filter.%d.quality", i)] = f.Operator
			params[fmt.Sprintf("filter.%d.value", i)] = f.Value
		}
	}

	var result struct {
		Assets []WorkbenchAsset `json:"assets"`
	}

	_, err := w.client.GetWithParams(ctx, "workbenches/assets", params, &result)
	if err != nil {
		return nil, err
	}

	return result.Assets, nil
}

// Vulnerabilities retrieves vulnerabilities from the workbench.
func (w *WorkbenchesAPI) Vulnerabilities(ctx context.Context, opts *WorkbenchOptions) ([]WorkbenchVuln, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.DateRange > 0 {
			params["date_range"] = strconv.Itoa(opts.DateRange)
		}
		if opts.FilterType != "" {
			params["filter.search_type"] = opts.FilterType
		}
		for i, f := range opts.Filters {
			params[fmt.Sprintf("filter.%d.filter", i)] = f.Name
			params[fmt.Sprintf("filter.%d.quality", i)] = f.Operator
			params[fmt.Sprintf("filter.%d.value", i)] = f.Value
		}
	}

	var result struct {
		Vulnerabilities []WorkbenchVuln `json:"vulnerabilities"`
	}

	_, err := w.client.GetWithParams(ctx, "workbenches/vulnerabilities", params, &result)
	if err != nil {
		return nil, err
	}

	return result.Vulnerabilities, nil
}

// AssetVulnerabilities retrieves vulnerabilities for a specific asset.
func (w *WorkbenchesAPI) AssetVulnerabilities(ctx context.Context, assetUUID string, dateRange int) ([]WorkbenchVuln, error) {
	params := map[string]string{}
	if dateRange > 0 {
		params["date_range"] = strconv.Itoa(dateRange)
	}

	var result struct {
		Vulnerabilities []WorkbenchVuln `json:"vulnerabilities"`
	}

	_, err := w.client.GetWithParams(ctx, fmt.Sprintf("workbenches/assets/%s/vulnerabilities", assetUUID), params, &result)
	if err != nil {
		return nil, err
	}

	return result.Vulnerabilities, nil
}

// VulnerabilityAssets retrieves assets affected by a vulnerability.
func (w *WorkbenchesAPI) VulnerabilityAssets(ctx context.Context, pluginID, dateRange int) ([]WorkbenchAsset, error) {
	params := map[string]string{}
	if dateRange > 0 {
		params["date_range"] = strconv.Itoa(dateRange)
	}

	var result struct {
		Assets []WorkbenchAsset `json:"assets"`
	}

	_, err := w.client.GetWithParams(ctx, fmt.Sprintf("workbenches/vulnerabilities/%d/outputs", pluginID), params, &result)
	if err != nil {
		return nil, err
	}

	return result.Assets, nil
}

// Export exports workbench data.
func (w *WorkbenchesAPI) Export(ctx context.Context, format string, opts *WorkbenchOptions, chapters []string) (io.Reader, error) {
	params := map[string]string{
		"format": format,
	}
	if opts != nil {
		if opts.DateRange > 0 {
			params["date_range"] = strconv.Itoa(opts.DateRange)
		}
		if opts.FilterType != "" {
			params["filter.search_type"] = opts.FilterType
		}
		for i, f := range opts.Filters {
			params[fmt.Sprintf("filter.%d.filter", i)] = f.Name
			params[fmt.Sprintf("filter.%d.quality", i)] = f.Operator
			params[fmt.Sprintf("filter.%d.value", i)] = f.Value
		}
	}
	if len(chapters) > 0 {
		for i, ch := range chapters {
			params[fmt.Sprintf("chapter.%d", i)] = ch
		}
	}

	var exportResp struct {
		File int `json:"file"`
	}

	_, err := w.client.GetWithParams(ctx, "workbenches/export", params, &exportResp)
	if err != nil {
		return nil, err
	}

	fileID := exportResp.File

	// Wait for export to be ready
	for {
		var statusResp struct {
			Status string `json:"status"`
		}
		_, err := w.client.Get(ctx, fmt.Sprintf("workbenches/export/%d/status", fileID), &statusResp)
		if err != nil {
			return nil, err
		}

		if statusResp.Status == "ready" {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2500 * time.Millisecond):
		}
	}

	// Download the file
	data, err := w.client.Download(ctx, fmt.Sprintf("workbenches/export/%d/download", fileID))
	if err != nil {
		return nil, err
	}

	return &bytesReader{data: data}, nil
}
