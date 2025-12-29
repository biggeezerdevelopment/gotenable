package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/biggeezerdevelopment/gotenable/pkg/base"
)

// CredentialsAPI handles credential operations.
type CredentialsAPI struct {
	client *Client
}

// Credential represents a managed credential.
type Credential struct {
	UUID                 string            `json:"uuid"`
	Name                 string            `json:"name"`
	Description          string            `json:"description,omitempty"`
	Type                 CredentialType    `json:"type"`
	Category             CredentialCategory `json:"category"`
	Created              string            `json:"created_date"`
	CreatedBy            CredentialUser    `json:"created_by"`
	LastModified         string            `json:"last_used_by"`
	ModifiedBy           CredentialUser    `json:"modified_by,omitempty"`
	Permissions          []CredentialPermission `json:"permissions,omitempty"`
	UserPermissions      int               `json:"user_permissions"`
}

// CredentialType represents a credential type.
type CredentialType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CredentialCategory represents a credential category.
type CredentialCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CredentialUser represents a user associated with a credential.
type CredentialUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// CredentialPermission represents a permission on a credential.
type CredentialPermission struct {
	GranteeUUID string `json:"grantee_uuid"`
	Type        string `json:"type"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name,omitempty"`
}

// CredentialDetails contains detailed credential information.
type CredentialDetails struct {
	Credential
	Settings map[string]interface{} `json:"settings,omitempty"`
	AdHoc    bool                   `json:"ad_hoc"`
}

// CredentialCreateRequest represents a request to create a credential.
type CredentialCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"`
	Settings    map[string]interface{} `json:"settings"`
	Permissions []CredentialPermission `json:"permissions,omitempty"`
}

// List retrieves all credentials.
func (c *CredentialsAPI) List(ctx context.Context) *base.Iterator[Credential] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			Credentials []Credential `json:"credentials"`
			Pagination  struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := c.client.GetWithParams(ctx, "credentials", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Credentials)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]Credential, error) {
		var items []Credential
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// Get retrieves a specific credential.
func (c *CredentialsAPI) Get(ctx context.Context, credentialUUID string) (*CredentialDetails, error) {
	var result CredentialDetails
	_, err := c.client.Get(ctx, fmt.Sprintf("credentials/%s", credentialUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new credential.
func (c *CredentialsAPI) Create(ctx context.Context, req *CredentialCreateRequest) (*Credential, error) {
	var result Credential
	_, err := c.client.Post(ctx, "credentials", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a credential.
func (c *CredentialsAPI) Update(ctx context.Context, credentialUUID string, req *CredentialCreateRequest) (*Credential, error) {
	var result Credential
	_, err := c.client.Put(ctx, fmt.Sprintf("credentials/%s", credentialUUID), req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a credential.
func (c *CredentialsAPI) Delete(ctx context.Context, credentialUUID string) error {
	_, err := c.client.Delete(ctx, fmt.Sprintf("credentials/%s", credentialUUID))
	return err
}

// CredentialTypes retrieves available credential types.
func (c *CredentialsAPI) Types(ctx context.Context) ([]CredentialType, error) {
	var result struct {
		Credentials []struct {
			Category CredentialCategory `json:"category"`
			Types    []CredentialType   `json:"types"`
		} `json:"credentials"`
	}

	_, err := c.client.Get(ctx, "credentials/types", &result)
	if err != nil {
		return nil, err
	}

	var types []CredentialType
	for _, cat := range result.Credentials {
		types = append(types, cat.Types...)
	}
	return types, nil
}

