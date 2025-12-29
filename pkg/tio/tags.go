package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// TagsAPI handles tag operations.
type TagsAPI struct {
	client *Client
}

// TagCategory represents a tag category.
type TagCategory struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Reserved    bool      `json:"reserved"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
	ProductUUID string    `json:"product_uuid,omitempty"`
}

// TagValue represents a tag value.
type TagValue struct {
	UUID                string            `json:"uuid"`
	CategoryUUID        string            `json:"category_uuid"`
	CategoryName        string            `json:"category_name,omitempty"`
	Value               string            `json:"value"`
	Description         string            `json:"description,omitempty"`
	Type                string            `json:"type,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
	CreatedBy           string            `json:"created_by"`
	UpdatedAt           time.Time         `json:"updated_at"`
	UpdatedBy           string            `json:"updated_by"`
	CategoryDescription string            `json:"category_description,omitempty"`
	AccessControl       *TagAccessControl `json:"access_control,omitempty"`
}

// TagAccessControl represents access control for a tag.
type TagAccessControl struct {
	AllUsersPermissions      int `json:"all_users_permissions"`
	CurrentUserPermissions   int `json:"current_user_permissions"`
	CurrentDomainPermissions int `json:"current_domain_permissions,omitempty"`
	DefinedDomainPermissions int `json:"defined_domain_permissions,omitempty"`
	Version                  int `json:"version"`
}

// TagFilter represents a filter for listing tags.
type TagFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// ListCategories retrieves all tag categories.
func (t *TagsAPI) ListCategories(ctx context.Context) *base.Iterator[TagCategory] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			Categories []TagCategory `json:"categories"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := t.client.GetWithParams(ctx, "tags/categories", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Categories)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]TagCategory, error) {
		var items []TagCategory
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// CreateCategory creates a new tag category.
func (t *TagsAPI) CreateCategory(ctx context.Context, name, description string) (*TagCategory, error) {
	payload := map[string]string{
		"name":        name,
		"description": description,
	}

	var result TagCategory
	_, err := t.client.Post(ctx, "tags/categories", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCategory retrieves a specific tag category.
func (t *TagsAPI) GetCategory(ctx context.Context, categoryUUID string) (*TagCategory, error) {
	var result TagCategory
	_, err := t.client.Get(ctx, fmt.Sprintf("tags/categories/%s", categoryUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCategory updates a tag category.
func (t *TagsAPI) UpdateCategory(ctx context.Context, categoryUUID, name, description string) (*TagCategory, error) {
	payload := map[string]string{
		"name":        name,
		"description": description,
	}

	var result TagCategory
	_, err := t.client.Put(ctx, fmt.Sprintf("tags/categories/%s", categoryUUID), payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCategory deletes a tag category.
func (t *TagsAPI) DeleteCategory(ctx context.Context, categoryUUID string) error {
	_, err := t.client.Delete(ctx, fmt.Sprintf("tags/categories/%s", categoryUUID))
	return err
}

// ListValues retrieves all tag values.
func (t *TagsAPI) ListValues(ctx context.Context, filters []TagFilter) *base.Iterator[TagValue] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		// Add filters
		for i, f := range filters {
			params[fmt.Sprintf("f.%d.field", i)] = f.Field
			params[fmt.Sprintf("f.%d.operator", i)] = f.Operator
			params[fmt.Sprintf("f.%d.value", i)] = f.Value
		}

		var result struct {
			Values     []TagValue `json:"values"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := t.client.GetWithParams(ctx, "tags/values", params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.Values)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]TagValue, error) {
		var items []TagValue
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer)
}

// CreateValue creates a new tag value.
func (t *TagsAPI) CreateValue(ctx context.Context, categoryUUID, value, description string) (*TagValue, error) {
	payload := map[string]string{
		"category_uuid": categoryUUID,
		"value":         value,
		"description":   description,
	}

	var result TagValue
	_, err := t.client.Post(ctx, "tags/values", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetValue retrieves a specific tag value.
func (t *TagsAPI) GetValue(ctx context.Context, valueUUID string) (*TagValue, error) {
	var result TagValue
	_, err := t.client.Get(ctx, fmt.Sprintf("tags/values/%s", valueUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateValue updates a tag value.
func (t *TagsAPI) UpdateValue(ctx context.Context, valueUUID, value, description string) (*TagValue, error) {
	payload := map[string]string{
		"value":       value,
		"description": description,
	}

	var result TagValue
	_, err := t.client.Put(ctx, fmt.Sprintf("tags/values/%s", valueUUID), payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteValue deletes a tag value.
func (t *TagsAPI) DeleteValue(ctx context.Context, valueUUID string) error {
	_, err := t.client.Delete(ctx, fmt.Sprintf("tags/values/%s", valueUUID))
	return err
}

// TagAssignment represents a tag assignment to an asset.
type TagAssignment struct {
	AssetUUID    string `json:"asset_uuid"`
	TagValueUUID string `json:"tag_value_uuid,omitempty"`
}

// AssignTags assigns tags to assets.
func (t *TagsAPI) AssignTags(ctx context.Context, assetUUIDs, tagValueUUIDs []string) error {
	payload := map[string]interface{}{
		"action": "add",
		"assets": assetUUIDs,
		"tags":   tagValueUUIDs,
	}
	_, err := t.client.Post(ctx, "tags/assets/assignments", payload, nil)
	return err
}

// UnassignTags removes tags from assets.
func (t *TagsAPI) UnassignTags(ctx context.Context, assetUUIDs, tagValueUUIDs []string) error {
	payload := map[string]interface{}{
		"action": "remove",
		"assets": assetUUIDs,
		"tags":   tagValueUUIDs,
	}
	_, err := t.client.Post(ctx, "tags/assets/assignments", payload, nil)
	return err
}

// GetAssetTags retrieves tags for an asset.
func (t *TagsAPI) GetAssetTags(ctx context.Context, assetUUID string) ([]TagValue, error) {
	var result struct {
		Tags []TagValue `json:"tags"`
	}
	_, err := t.client.Get(ctx, fmt.Sprintf("tags/assets/%s/assignments", assetUUID), &result)
	if err != nil {
		return nil, err
	}
	return result.Tags, nil
}
