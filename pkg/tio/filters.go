package tio

import (
	"context"
)

// FiltersAPI handles filter operations.
type FiltersAPI struct {
	client *Client
}

// Filter represents an available filter option.
type Filter struct {
	Name         string   `json:"name"`
	ReadableName string   `json:"readable_name"`
	Control      Control  `json:"control"`
	Operators    []string `json:"operators"`
	GroupName    string   `json:"group_name,omitempty"`
}

// FilterSet represents a collection of filters for a specific endpoint.
type FilterSet struct {
	Filters []Filter `json:"filters"`
}

// ScanFilters retrieves filters available for scan results.
func (f *FiltersAPI) ScanFilters(ctx context.Context) (map[string]Filter, error) {
	var result struct {
		Filters []Filter `json:"filters"`
	}

	_, err := f.client.Get(ctx, "filters/scans/reports", &result)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]Filter)
	for _, filter := range result.Filters {
		filters[filter.Name] = filter
	}
	return filters, nil
}

// VulnFilters retrieves filters available for vulnerability workbench.
func (f *FiltersAPI) VulnFilters(ctx context.Context) (map[string]Filter, error) {
	var result struct {
		Filters []Filter `json:"filters"`
	}

	_, err := f.client.Get(ctx, "filters/workbenches/vulnerabilities", &result)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]Filter)
	for _, filter := range result.Filters {
		filters[filter.Name] = filter
	}
	return filters, nil
}

// AssetFilters retrieves filters available for assets.
func (f *FiltersAPI) AssetFilters(ctx context.Context) (map[string]Filter, error) {
	var result struct {
		Filters []Filter `json:"filters"`
	}

	_, err := f.client.Get(ctx, "filters/workbenches/assets", &result)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]Filter)
	for _, filter := range result.Filters {
		filters[filter.Name] = filter
	}
	return filters, nil
}

// AgentFilters retrieves filters available for agents.
func (f *FiltersAPI) AgentFilters(ctx context.Context, scannerID int) (map[string]Filter, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result struct {
		Filters []Filter `json:"filters"`
	}

	_, err := f.client.Get(ctx, "filters/scans/agents", &result)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]Filter)
	for _, filter := range result.Filters {
		filters[filter.Name] = filter
	}
	return filters, nil
}

// CredentialFilters retrieves filters available for credentials.
func (f *FiltersAPI) CredentialFilters(ctx context.Context) (map[string]Filter, error) {
	var result struct {
		Filters []Filter `json:"filters"`
	}

	_, err := f.client.Get(ctx, "filters/credentials", &result)
	if err != nil {
		return nil, err
	}

	filters := make(map[string]Filter)
	for _, filter := range result.Filters {
		filters[filter.Name] = filter
	}
	return filters, nil
}

