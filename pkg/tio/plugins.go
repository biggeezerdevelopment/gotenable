package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/biggeezerdevelopment/gotenable/pkg/base"
)

// PluginsAPI handles plugin operations.
type PluginsAPI struct {
	client *Client
}

// Plugin represents a Nessus plugin.
type Plugin struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
}

// PluginDetails contains detailed plugin information.
type PluginDetails struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	FamilyName string            `json:"family_name"`
	Attributes []PluginAttribute `json:"attributes"`
}

// PluginFamily represents a plugin family.
type PluginFamily struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// PluginFamilyDetails contains detailed plugin family information.
type PluginFamilyDetails struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Plugins []Plugin `json:"plugins"`
}

// Families retrieves all plugin families.
func (p *PluginsAPI) Families(ctx context.Context) ([]PluginFamily, error) {
	var result struct {
		Families []PluginFamily `json:"families"`
	}

	_, err := p.client.Get(ctx, "plugins/families", &result)
	if err != nil {
		return nil, err
	}

	return result.Families, nil
}

// FamilyDetails retrieves details about a plugin family.
func (p *PluginsAPI) FamilyDetails(ctx context.Context, familyID int) (*PluginFamilyDetails, error) {
	var result PluginFamilyDetails
	_, err := p.client.Get(ctx, fmt.Sprintf("plugins/families/%d", familyID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves details about a specific plugin.
func (p *PluginsAPI) Get(ctx context.Context, pluginID int) (*PluginDetails, error) {
	var result PluginDetails
	_, err := p.client.Get(ctx, fmt.Sprintf("plugins/plugin/%d", pluginID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginListOptions contains options for listing plugins.
type PluginListOptions struct {
	Size        int
	Page        int
	LastUpdated int64
}

// List retrieves plugins with pagination.
func (p *PluginsAPI) List(ctx context.Context, opts *PluginListOptions) *base.Iterator[Plugin] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"size": strconv.Itoa(limit),
			"page": strconv.Itoa(offset / limit),
		}
		if opts != nil && opts.LastUpdated > 0 {
			params["last_updated"] = strconv.FormatInt(opts.LastUpdated, 10)
		}

		var result struct {
			PluginData struct {
				FamilyCount int `json:"family_count"`
				PluginCount int `json:"plugin_count"`
			} `json:"plugin_data"`
			PluginFamilyDetails []struct {
				ID      int      `json:"id"`
				Name    string   `json:"name"`
				Plugins []Plugin `json:"plugins"`
			} `json:"plugin_family_details"`
		}

		_, err := p.client.GetWithParams(ctx, "plugins", params, &result)
		if err != nil {
			return nil, nil, err
		}

		var plugins []Plugin
		for _, family := range result.PluginFamilyDetails {
			for _, plugin := range family.Plugins {
				plugin.FamilyName = family.Name
				plugins = append(plugins, plugin)
			}
		}

		data, _ := json.Marshal(plugins)
		return data, &base.PaginationInfo{
			Total: result.PluginData.PluginCount,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]Plugin, error) {
		var items []Plugin
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}
