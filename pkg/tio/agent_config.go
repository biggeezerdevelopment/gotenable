package tio

import (
	"context"
	"fmt"
)

// AgentConfigAPI handles agent configuration operations.
type AgentConfigAPI struct {
	client *Client
}

// AgentConfig represents agent configuration settings.
type AgentConfig struct {
	AutoUnlink     AgentAutoUnlink `json:"auto_unlink"`
	SoftwareUpdate bool            `json:"software_update"`
}

// AgentAutoUnlink represents auto-unlink settings.
type AgentAutoUnlink struct {
	Enabled    bool `json:"enabled"`
	Expiration int  `json:"expiration"`
}

// Get retrieves agent configuration for a scanner.
func (a *AgentConfigAPI) Get(ctx context.Context, scannerID int) (*AgentConfig, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentConfig
	_, err := a.client.Get(ctx, fmt.Sprintf("scanners/%d/agents/config", scannerID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Edit updates agent configuration for a scanner.
func (a *AgentConfigAPI) Edit(ctx context.Context, scannerID int, config *AgentConfig) (*AgentConfig, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentConfig
	_, err := a.client.Put(ctx, fmt.Sprintf("scanners/%d/agents/config", scannerID), config, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// EnableAutoUnlink enables auto-unlink with the specified expiration.
func (a *AgentConfigAPI) EnableAutoUnlink(ctx context.Context, scannerID, expirationDays int) (*AgentConfig, error) {
	config := &AgentConfig{
		AutoUnlink: AgentAutoUnlink{
			Enabled:    true,
			Expiration: expirationDays,
		},
	}
	return a.Edit(ctx, scannerID, config)
}

// DisableAutoUnlink disables auto-unlink.
func (a *AgentConfigAPI) DisableAutoUnlink(ctx context.Context, scannerID int) (*AgentConfig, error) {
	config := &AgentConfig{
		AutoUnlink: AgentAutoUnlink{
			Enabled: false,
		},
	}
	return a.Edit(ctx, scannerID, config)
}

// EnableSoftwareUpdate enables software updates.
func (a *AgentConfigAPI) EnableSoftwareUpdate(ctx context.Context, scannerID int) (*AgentConfig, error) {
	config := &AgentConfig{
		SoftwareUpdate: true,
	}
	return a.Edit(ctx, scannerID, config)
}

// DisableSoftwareUpdate disables software updates.
func (a *AgentConfigAPI) DisableSoftwareUpdate(ctx context.Context, scannerID int) (*AgentConfig, error) {
	config := &AgentConfig{
		SoftwareUpdate: false,
	}
	return a.Edit(ctx, scannerID, config)
}

