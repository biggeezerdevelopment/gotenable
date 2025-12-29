package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// AgentsAPI handles agent-related operations.
type AgentsAPI struct {
	client *Client
}

// Agent represents a Nessus agent.
type Agent struct {
	ID                     int          `json:"id"`
	UUID                   string       `json:"uuid"`
	Name                   string       `json:"name"`
	Platform               string       `json:"platform"`
	Distro                 string       `json:"distro"`
	IP                     string       `json:"ip"`
	LastScanned            time.Time    `json:"last_scanned,omitempty"`
	PluginFeedID           string       `json:"plugin_feed_id"`
	CoreBuild              string       `json:"core_build"`
	CoreVersion            string       `json:"core_version"`
	LinkedOn               time.Time    `json:"linked_on"`
	LastConnect            time.Time    `json:"last_connect"`
	Status                 string       `json:"status"`
	Groups                 []AgentGroup `json:"groups,omitempty"`
	NetworkUUID            string       `json:"network_uuid,omitempty"`
	NetworkName            string       `json:"network_name,omitempty"`
	SupportsRemoteLogs     bool         `json:"supports_remote_logs"`
	SupportsRemoteSettings bool         `json:"supports_remote_settings"`
	Restart                bool         `json:"restart_pending"`
	NeedsRestart           bool         `json:"needs_restart"`
	NeedsPluginUpdate      bool         `json:"needs_plugin_update"`
}

// AgentListOptions contains options for listing agents.
type AgentListOptions struct {
	ScannerID      int
	Limit          int
	Offset         int
	Sort           string
	Filter         string
	Wildcard       string
	WildcardFields string
}

// List retrieves all agents.
func (a *AgentsAPI) List(ctx context.Context, opts *AgentListOptions) *base.Iterator[Agent] {
	scannerID := 1 // Default scanner
	if opts != nil && opts.ScannerID > 0 {
		scannerID = opts.ScannerID
	}

	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}
		if opts != nil {
			if opts.Sort != "" {
				params["sort"] = opts.Sort
			}
			if opts.Filter != "" {
				params["f"] = opts.Filter
			}
			if opts.Wildcard != "" {
				params["w"] = opts.Wildcard
			}
			if opts.WildcardFields != "" {
				params["wf"] = opts.WildcardFields
			}
		}

		var result struct {
			Agents     []Agent `json:"agents"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := a.client.GetWithParams(ctx, fmt.Sprintf("scanners/%d/agents", scannerID), params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Agents)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]Agent, error) {
		var items []Agent
		err := json.Unmarshal(data, &items)
		return items, err
	}

	limit := 100
	offset := 0
	if opts != nil {
		if opts.Limit > 0 {
			limit = opts.Limit
		}
		if opts.Offset > 0 {
			offset = opts.Offset
		}
	}

	return base.NewIterator(ctx, fetcher, transformer,
		base.WithLimit[Agent](limit),
		base.WithOffset[Agent](offset),
	)
}

// Get retrieves a specific agent.
func (a *AgentsAPI) Get(ctx context.Context, scannerID, agentID int) (*Agent, error) {
	var result Agent
	_, err := a.client.Get(ctx, fmt.Sprintf("scanners/%d/agents/%d", scannerID, agentID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an agent.
func (a *AgentsAPI) Delete(ctx context.Context, scannerID, agentID int) error {
	_, err := a.client.Delete(ctx, fmt.Sprintf("scanners/%d/agents/%d", scannerID, agentID))
	return err
}

// Unlink unlinks an agent from the scanner.
func (a *AgentsAPI) Unlink(ctx context.Context, scannerID, agentID int) error {
	return a.Delete(ctx, scannerID, agentID)
}

// BulkUnlink unlinks multiple agents.
func (a *AgentsAPI) BulkUnlink(ctx context.Context, scannerID int, agentIDs []int) error {
	payload := map[string]interface{}{
		"items": agentIDs,
	}
	_, err := a.client.Post(ctx, fmt.Sprintf("scanners/%d/agents/_bulk/unlink", scannerID), payload, nil)
	return err
}

// BulkAddToGroup adds multiple agents to a group.
func (a *AgentsAPI) BulkAddToGroup(ctx context.Context, scannerID, groupID int, agentIDs []int) error {
	payload := map[string]interface{}{
		"items": agentIDs,
	}
	_, err := a.client.Post(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d/agents/_bulk/add", scannerID, groupID), payload, nil)
	return err
}

// BulkRemoveFromGroup removes multiple agents from a group.
func (a *AgentsAPI) BulkRemoveFromGroup(ctx context.Context, scannerID, groupID int, agentIDs []int) error {
	payload := map[string]interface{}{
		"items": agentIDs,
	}
	_, err := a.client.Post(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d/agents/_bulk/remove", scannerID, groupID), payload, nil)
	return err
}

// TaskStatus retrieves the status of a bulk operation task.
func (a *AgentsAPI) TaskStatus(ctx context.Context, scannerID int, taskUUID string) (*TaskStatusResponse, error) {
	var result TaskStatusResponse
	_, err := a.client.Get(ctx, fmt.Sprintf("scanners/%d/agents/_bulk/%s", scannerID, taskUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TaskStatusResponse contains task status information.
type TaskStatusResponse struct {
	TaskID      string `json:"task_id"`
	ContainerID string `json:"container_id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	StartTime   int64  `json:"start_time"`
	EndTime     int64  `json:"end_time"`
	TotalWork   int    `json:"total_work_units"`
	WorkDone    int    `json:"total_work_units_completed"`
	LastUpdate  int64  `json:"last_update_time"`
}

// AgentGroup represents an agent group.
type AgentGroup struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	UUID                 string `json:"uuid,omitempty"`
	OwnerID              int    `json:"owner_id"`
	Owner                string `json:"owner"`
	OwnerName            string `json:"owner_name"`
	OwnerUUID            string `json:"owner_uuid"`
	Shared               int    `json:"shared"`
	UserPermissions      int    `json:"user_permissions"`
	CreationDate         int64  `json:"creation_date"`
	LastModificationDate int64  `json:"last_modification_date"`
	Timestamp            int64  `json:"timestamp"`
	AgentsCount          int    `json:"agents_count"`
}

// AgentGroupsAPI handles agent group operations.
type AgentGroupsAPI struct {
	client *Client
}

// List retrieves all agent groups.
func (g *AgentGroupsAPI) List(ctx context.Context, scannerID int) ([]AgentGroup, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result struct {
		Groups []AgentGroup `json:"groups"`
	}

	_, err := g.client.Get(ctx, fmt.Sprintf("scanners/%d/agent-groups", scannerID), &result)
	if err != nil {
		return nil, err
	}

	return result.Groups, nil
}

// Create creates a new agent group.
func (g *AgentGroupsAPI) Create(ctx context.Context, scannerID int, name string) (*AgentGroup, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	payload := map[string]string{"name": name}
	var result AgentGroup

	_, err := g.client.Post(ctx, fmt.Sprintf("scanners/%d/agent-groups", scannerID), payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes an agent group.
func (g *AgentGroupsAPI) Delete(ctx context.Context, scannerID, groupID int) error {
	if scannerID == 0 {
		scannerID = 1
	}
	_, err := g.client.Delete(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d", scannerID, groupID))
	return err
}

// Details retrieves details about an agent group.
func (g *AgentGroupsAPI) Details(ctx context.Context, scannerID, groupID int) (*AgentGroup, error) {
	if scannerID == 0 {
		scannerID = 1
	}

	var result AgentGroup
	_, err := g.client.Get(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d", scannerID, groupID), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Configure updates an agent group.
func (g *AgentGroupsAPI) Configure(ctx context.Context, scannerID, groupID int, name string) error {
	if scannerID == 0 {
		scannerID = 1
	}

	payload := map[string]string{"name": name}
	_, err := g.client.Put(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d", scannerID, groupID), payload, nil)
	return err
}

// AddAgent adds an agent to a group.
func (g *AgentGroupsAPI) AddAgent(ctx context.Context, scannerID, groupID, agentID int) error {
	if scannerID == 0 {
		scannerID = 1
	}
	_, err := g.client.Post(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d/agents/%d", scannerID, groupID, agentID), nil, nil)
	return err
}

// RemoveAgent removes an agent from a group.
func (g *AgentGroupsAPI) RemoveAgent(ctx context.Context, scannerID, groupID, agentID int) error {
	if scannerID == 0 {
		scannerID = 1
	}
	_, err := g.client.Delete(ctx, fmt.Sprintf("scanners/%d/agent-groups/%d/agents/%d", scannerID, groupID, agentID))
	return err
}
