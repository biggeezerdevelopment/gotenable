package tio

import (
	"context"
	"fmt"
	"time"
)

// ExclusionsAPI handles exclusion operations.
type ExclusionsAPI struct {
	client *Client
}

// Exclusion represents a scan exclusion.
type Exclusion struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description,omitempty"`
	CreationDate         int64     `json:"creation_date"`
	LastModificationDate int64     `json:"last_modification_date"`
	Schedule             *ExclusionSchedule `json:"schedule,omitempty"`
	Members              string    `json:"members"`
	NetworkID            string    `json:"network_id,omitempty"`
}

// ExclusionSchedule represents the schedule for an exclusion.
type ExclusionSchedule struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"starttime,omitempty"`
	EndTime   string `json:"endtime,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	RRules    string `json:"rrules,omitempty"`
}

// ExclusionCreateRequest represents a request to create an exclusion.
type ExclusionCreateRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Members     string             `json:"members"`
	Schedule    *ExclusionSchedule `json:"schedule,omitempty"`
	NetworkID   string             `json:"network_id,omitempty"`
}

// List retrieves all exclusions.
func (e *ExclusionsAPI) List(ctx context.Context) ([]Exclusion, error) {
	var result struct {
		Exclusions []Exclusion `json:"exclusions"`
	}

	_, err := e.client.Get(ctx, "exclusions", &result)
	if err != nil {
		return nil, err
	}

	return result.Exclusions, nil
}

// Get retrieves a specific exclusion.
func (e *ExclusionsAPI) Get(ctx context.Context, exclusionID int) (*Exclusion, error) {
	var result Exclusion
	_, err := e.client.Get(ctx, fmt.Sprintf("exclusions/%d", exclusionID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new exclusion.
func (e *ExclusionsAPI) Create(ctx context.Context, req *ExclusionCreateRequest) (*Exclusion, error) {
	var result Exclusion
	_, err := e.client.Post(ctx, "exclusions", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an exclusion.
func (e *ExclusionsAPI) Update(ctx context.Context, exclusionID int, req *ExclusionCreateRequest) (*Exclusion, error) {
	var result Exclusion
	_, err := e.client.Put(ctx, fmt.Sprintf("exclusions/%d", exclusionID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an exclusion.
func (e *ExclusionsAPI) Delete(ctx context.Context, exclusionID int) error {
	_, err := e.client.Delete(ctx, fmt.Sprintf("exclusions/%d", exclusionID))
	return err
}

// Import imports exclusions from a file.
func (e *ExclusionsAPI) Import(ctx context.Context, filename string) error {
	payload := map[string]string{"file": filename}
	_, err := e.client.Post(ctx, "exclusions/import", payload, nil)
	return err
}

// AgentExclusionsAPI handles agent exclusion operations.
type AgentExclusionsAPI struct {
	client *Client
}

// AgentExclusion represents an agent exclusion.
type AgentExclusion struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description,omitempty"`
	CreationDate         int64     `json:"creation_date"`
	LastModificationDate int64     `json:"last_modification_date"`
	Schedule             *AgentExclusionSchedule `json:"schedule,omitempty"`
}

// AgentExclusionSchedule represents the schedule for an agent exclusion.
type AgentExclusionSchedule struct {
	Enabled   bool   `json:"enabled"`
	StartTime string `json:"starttime,omitempty"`
	EndTime   string `json:"endtime,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	RRules    string `json:"rrules,omitempty"`
}

// AgentExclusionCreateRequest represents a request to create an agent exclusion.
type AgentExclusionCreateRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Schedule    *AgentExclusionSchedule `json:"schedule,omitempty"`
}

// List retrieves all agent exclusions.
func (e *AgentExclusionsAPI) List(ctx context.Context, scannerID int) ([]AgentExclusion, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result struct {
		Exclusions []AgentExclusion `json:"exclusions"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("scanners/%d/agents/exclusions", scannerID), &result)
	if err != nil {
		return nil, err
	}

	return result.Exclusions, nil
}

// Get retrieves a specific agent exclusion.
func (e *AgentExclusionsAPI) Get(ctx context.Context, scannerID, exclusionID int) (*AgentExclusion, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentExclusion
	_, err := e.client.Get(ctx, fmt.Sprintf("scanners/%d/agents/exclusions/%d", scannerID, exclusionID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new agent exclusion.
func (e *AgentExclusionsAPI) Create(ctx context.Context, scannerID int, req *AgentExclusionCreateRequest) (*AgentExclusion, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentExclusion
	_, err := e.client.Post(ctx, fmt.Sprintf("scanners/%d/agents/exclusions", scannerID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an agent exclusion.
func (e *AgentExclusionsAPI) Update(ctx context.Context, scannerID, exclusionID int, req *AgentExclusionCreateRequest) (*AgentExclusion, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentExclusion
	_, err := e.client.Put(ctx, fmt.Sprintf("scanners/%d/agents/exclusions/%d", scannerID, exclusionID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an agent exclusion.
func (e *AgentExclusionsAPI) Delete(ctx context.Context, scannerID, exclusionID int) error {
	if scannerID == 0 {
		scannerID = 1
	}
	_, err := e.client.Delete(ctx, fmt.Sprintf("scanners/%d/agents/exclusions/%d", scannerID, exclusionID))
	return err
}

// CreateSchedule creates a schedule for agent exclusions.
func CreateSchedule(enabled bool, startTime, endTime time.Time, timezone, rrules string) *AgentExclusionSchedule {
	return &AgentExclusionSchedule{
		Enabled:   enabled,
		StartTime: startTime.Format("2006-01-02T15:04:05"),
		EndTime:   endTime.Format("2006-01-02T15:04:05"),
		Timezone:  timezone,
		RRules:    rrules,
	}
}

