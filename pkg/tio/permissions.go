package tio

import (
	"context"
	"fmt"
)

// PermissionsAPI handles permission operations.
type PermissionsAPI struct {
	client *Client
}

// Permission represents a permission entry.
type Permission struct {
	ID          int    `json:"id"`
	Owner       int    `json:"owner,omitempty"`
	Type        string `json:"type"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

// PermissionConstants defines permission values.
const (
	PermissionNone    = 0
	PermissionView    = 16
	PermissionScan    = 32
	PermissionControl = 64
	PermissionDefault = 16
)

// List retrieves permissions for an object.
func (p *PermissionsAPI) List(ctx context.Context, objectType string, objectID int) ([]Permission, error) {
	var result []Permission
	_, err := p.client.Get(ctx, fmt.Sprintf("%s/%d/permissions", objectType, objectID), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Change updates permissions for an object.
func (p *PermissionsAPI) Change(ctx context.Context, objectType string, objectID int, acls []Permission) error {
	payload := map[string]interface{}{
		"acls": acls,
	}
	_, err := p.client.Put(ctx, fmt.Sprintf("%s/%d/permissions", objectType, objectID), payload, nil)
	return err
}

// AccessControlAPI handles access control operations.
type AccessControlAPI struct {
	client *Client
}

// AccessGroup represents an access group.
type AccessGroup struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	ContainerUUID     string   `json:"container_uuid,omitempty"`
	AllAssets         bool     `json:"all_assets"`
	AllUsers          bool     `json:"all_users"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	CreatedByUUID     string   `json:"created_by_uuid,omitempty"`
	UpdatedByUUID     string   `json:"updated_by_uuid,omitempty"`
	CreatedByName     string   `json:"created_by_name,omitempty"`
	UpdatedByName     string   `json:"updated_by_name,omitempty"`
	ProcessingPercent int      `json:"processing_percent_complete"`
	Status            string   `json:"status"`
	Principals        []Principal `json:"principals,omitempty"`
	Rules             []AccessRule `json:"rules,omitempty"`
}

// Principal represents a principal (user or group) in an access group.
type Principal struct {
	Type        string   `json:"type"`
	PrincipalID string   `json:"principal_id"`
	PrincipalName string `json:"principal_name,omitempty"`
	Permissions []string `json:"permissions"`
}

// AccessRule represents an access rule.
type AccessRule struct {
	Type     string   `json:"type"`
	Operator string   `json:"operator"`
	Terms    []string `json:"terms"`
}

// List retrieves all access groups.
func (a *AccessControlAPI) List(ctx context.Context) ([]AccessGroup, error) {
	var result struct {
		AccessGroups []AccessGroup `json:"access_groups"`
	}

	_, err := a.client.Get(ctx, "access-groups", &result)
	if err != nil {
		return nil, err
	}

	return result.AccessGroups, nil
}

// Get retrieves a specific access group.
func (a *AccessControlAPI) Get(ctx context.Context, groupUUID string) (*AccessGroup, error) {
	var result AccessGroup
	_, err := a.client.Get(ctx, fmt.Sprintf("access-groups/%s", groupUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new access group.
func (a *AccessControlAPI) Create(ctx context.Context, group *AccessGroup) (*AccessGroup, error) {
	var result AccessGroup
	_, err := a.client.Post(ctx, "access-groups", group, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an access group.
func (a *AccessControlAPI) Update(ctx context.Context, groupUUID string, group *AccessGroup) (*AccessGroup, error) {
	var result AccessGroup
	_, err := a.client.Put(ctx, fmt.Sprintf("access-groups/%s", groupUUID), group, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an access group.
func (a *AccessControlAPI) Delete(ctx context.Context, groupUUID string) error {
	_, err := a.client.Delete(ctx, fmt.Sprintf("access-groups/%s", groupUUID))
	return err
}

