package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// NetworksAPI handles network operations.
type NetworksAPI struct {
	client *Client
}

// Network represents a network in Tenable.io.
type Network struct {
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	Description       string    `json:"description,omitempty"`
	IsDefault         bool      `json:"is_default"`
	Created           time.Time `json:"created"`
	CreatedBy         string    `json:"created_by"`
	Modified          time.Time `json:"modified"`
	ModifiedBy        string    `json:"modified_by"`
	OwnerUUID         string    `json:"owner_uuid"`
	CreatedInSeconds  int64     `json:"created_in_seconds"`
	ModifiedInSeconds int64     `json:"modified_in_seconds"`
	ScannerCount      int       `json:"scanner_count,omitempty"`
	AssetsTTLDays     int       `json:"assets_ttl_days,omitempty"`
}

// List retrieves all networks.
func (n *NetworksAPI) List(ctx context.Context) *base.Iterator[Network] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			Networks   []Network `json:"networks"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := n.client.GetWithParams(ctx, "networks", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Networks)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]Network, error) {
		var items []Network
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// Get retrieves a specific network.
func (n *NetworksAPI) Get(ctx context.Context, networkUUID string) (*Network, error) {
	var result Network
	_, err := n.client.Get(ctx, fmt.Sprintf("networks/%s", networkUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new network.
func (n *NetworksAPI) Create(ctx context.Context, name, description string, assetsTTLDays int) (*Network, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	if assetsTTLDays > 0 {
		payload["assets_ttl_days"] = assetsTTLDays
	}

	var result Network
	_, err := n.client.Post(ctx, "networks", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a network.
func (n *NetworksAPI) Update(ctx context.Context, networkUUID, name, description string, assetsTTLDays int) (*Network, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}
	if assetsTTLDays > 0 {
		payload["assets_ttl_days"] = assetsTTLDays
	}

	var result Network
	_, err := n.client.Put(ctx, fmt.Sprintf("networks/%s", networkUUID), payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a network.
func (n *NetworksAPI) Delete(ctx context.Context, networkUUID string) error {
	_, err := n.client.Delete(ctx, fmt.Sprintf("networks/%s", networkUUID))
	return err
}

// AssignScanners assigns scanners to a network.
func (n *NetworksAPI) AssignScanners(ctx context.Context, networkUUID string, scannerUUIDs []string) error {
	payload := map[string]interface{}{
		"scanners": scannerUUIDs,
	}
	_, err := n.client.Post(ctx, fmt.Sprintf("networks/%s/scanners", networkUUID), payload, nil)
	return err
}

// ListScanners retrieves scanners assigned to a network.
func (n *NetworksAPI) ListScanners(ctx context.Context, networkUUID string) ([]Scanner, error) {
	var result struct {
		Scanners []Scanner `json:"scanners"`
	}
	_, err := n.client.Get(ctx, fmt.Sprintf("networks/%s/scanners", networkUUID), &result)
	if err != nil {
		return nil, err
	}
	return result.Scanners, nil
}

// NetworkAssetCount represents asset counts in a network.
type NetworkAssetCount struct {
	NumAssets int `json:"numAssets"`
}

// AssetCount retrieves the asset count for a network.
func (n *NetworksAPI) AssetCount(ctx context.Context, networkUUID string) (int, error) {
	var result NetworkAssetCount
	_, err := n.client.Get(ctx, fmt.Sprintf("networks/%s/counts/assets", networkUUID), &result)
	if err != nil {
		return 0, err
	}
	return result.NumAssets, nil
}
