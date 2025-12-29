package tio

import (
	"context"
	"fmt"
	"io"
)

// PoliciesAPI handles policy operations.
type PoliciesAPI struct {
	client *Client
}

// Policy represents a scan policy.
type Policy struct {
	ID                   int    `json:"id"`
	TemplateUUID         string `json:"template_uuid"`
	Name                 string `json:"name"`
	Description          string `json:"description,omitempty"`
	OwnerID              int    `json:"owner_id"`
	Owner                string `json:"owner"`
	Shared               int    `json:"shared"`
	UserPermissions      int    `json:"user_permissions"`
	CreationDate         int64  `json:"creation_date"`
	LastModificationDate int64  `json:"last_modification_date"`
	Visibility           string `json:"visibility"`
	NoTarget             bool   `json:"no_target"`
	IsScap               bool   `json:"is_scap"`
}

// PolicyDetails contains detailed policy information.
type PolicyDetails struct {
	UUID        string                 `json:"uuid"`
	Settings    map[string]interface{} `json:"settings"`
	Plugins     map[string]interface{} `json:"plugins,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Audits      map[string]interface{} `json:"audits,omitempty"`
	Scap        map[string]interface{} `json:"scap,omitempty"`
}

// PolicyTemplate represents a policy template.
type PolicyTemplate struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CloudOnly   bool   `json:"cloud_only"`
	SubCategory string `json:"subscription_only"`
	IsAgent     bool   `json:"is_agent,omitempty"`
	MoreInfo    string `json:"more_info,omitempty"`
}

// List retrieves all policies.
func (p *PoliciesAPI) List(ctx context.Context) ([]Policy, error) {
	var result struct {
		Policies []Policy `json:"policies"`
	}

	_, err := p.client.Get(ctx, "policies", &result)
	if err != nil {
		return nil, err
	}

	return result.Policies, nil
}

// Get retrieves a specific policy.
func (p *PoliciesAPI) Get(ctx context.Context, policyID int) (*PolicyDetails, error) {
	var result PolicyDetails
	_, err := p.client.Get(ctx, fmt.Sprintf("policies/%d", policyID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new policy.
func (p *PoliciesAPI) Create(ctx context.Context, policy *PolicyDetails) (*Policy, error) {
	var result struct {
		PolicyID int    `json:"policy_id"`
		Name     string `json:"policy_name"`
	}

	_, err := p.client.Post(ctx, "policies", policy, &result)
	if err != nil {
		return nil, err
	}

	return &Policy{ID: result.PolicyID, Name: result.Name}, nil
}

// Update updates a policy.
func (p *PoliciesAPI) Update(ctx context.Context, policyID int, policy *PolicyDetails) error {
	_, err := p.client.Put(ctx, fmt.Sprintf("policies/%d", policyID), policy, nil)
	return err
}

// Delete removes a policy.
func (p *PoliciesAPI) Delete(ctx context.Context, policyID int) error {
	_, err := p.client.Delete(ctx, fmt.Sprintf("policies/%d", policyID))
	return err
}

// Copy duplicates a policy.
func (p *PoliciesAPI) Copy(ctx context.Context, policyID int) (*Policy, error) {
	var result Policy
	_, err := p.client.Post(ctx, fmt.Sprintf("policies/%d/copy", policyID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Export exports a policy to XML format.
func (p *PoliciesAPI) Export(ctx context.Context, policyID int) (io.Reader, error) {
	data, err := p.client.Download(ctx, fmt.Sprintf("policies/%d/export", policyID))
	if err != nil {
		return nil, err
	}
	return &bytesReader{data: data}, nil
}

// Import imports a policy from XML.
func (p *PoliciesAPI) Import(ctx context.Context, filename string) (*Policy, error) {
	var result Policy
	payload := map[string]string{"file": filename}
	_, err := p.client.Post(ctx, "policies/import", payload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Templates retrieves available policy templates.
func (p *PoliciesAPI) Templates(ctx context.Context) (map[string]string, error) {
	var result struct {
		Templates []PolicyTemplate `json:"templates"`
	}

	_, err := p.client.Get(ctx, "editor/policy/templates", &result)
	if err != nil {
		return nil, err
	}

	templates := make(map[string]string)
	for _, t := range result.Templates {
		templates[t.Name] = t.UUID
	}
	return templates, nil
}

