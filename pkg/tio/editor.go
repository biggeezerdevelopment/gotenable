package tio

import (
	"context"
	"fmt"
)

// EditorAPI handles editor operations for scan/policy configuration.
type EditorAPI struct {
	client *Client
}

// EditorTemplate represents a scan/policy template.
type EditorTemplate struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CloudOnly   bool   `json:"cloud_only"`
	MoreInfo    string `json:"more_info,omitempty"`
	IsAgent     bool   `json:"is_agent,omitempty"`
}

// EditorDetails represents detailed editor configuration.
type EditorDetails struct {
	UUID        string                 `json:"uuid"`
	Name        string                 `json:"name,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Plugins     map[string]interface{} `json:"plugins,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Compliance  map[string]interface{} `json:"compliance,omitempty"`
}

// Templates retrieves available templates.
func (e *EditorAPI) Templates(ctx context.Context, objectType string) ([]EditorTemplate, error) {
	var result struct {
		Templates []EditorTemplate `json:"templates"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/templates", objectType), &result)
	if err != nil {
		return nil, err
	}

	return result.Templates, nil
}

// ScanTemplates retrieves available scan templates.
func (e *EditorAPI) ScanTemplates(ctx context.Context) ([]EditorTemplate, error) {
	return e.Templates(ctx, "scan")
}

// PolicyTemplates retrieves available policy templates.
func (e *EditorAPI) PolicyTemplates(ctx context.Context) ([]EditorTemplate, error) {
	return e.Templates(ctx, "policy")
}

// Details retrieves the editor details for a scan or policy.
func (e *EditorAPI) Details(ctx context.Context, objectType string, objectID int) (*EditorDetails, error) {
	var result EditorDetails
	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/%d", objectType, objectID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanDetails retrieves editor details for a scan.
func (e *EditorAPI) ScanDetails(ctx context.Context, scanID int) (*EditorDetails, error) {
	return e.Details(ctx, "scan", scanID)
}

// PolicyDetails retrieves editor details for a policy.
func (e *EditorAPI) PolicyDetails(ctx context.Context, policyID int) (*EditorDetails, error) {
	return e.Details(ctx, "policy", policyID)
}

// TemplateDetails retrieves details about a specific template.
func (e *EditorAPI) TemplateDetails(ctx context.Context, objectType, templateUUID string) (*EditorDetails, error) {
	var result EditorDetails
	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/templates/%s", objectType, templateUUID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginDescription represents a plugin description.
type PluginDescription struct {
	PluginID       int               `json:"pluginid"`
	Name           string            `json:"pluginname"`
	Family         string            `json:"pluginfamily"`
	Severity       int               `json:"severity"`
	Attributes     []PluginAttribute `json:"pluginattributes,omitempty"`
}

// PluginFamilies retrieves plugin families for a template.
func (e *EditorAPI) PluginFamilies(ctx context.Context, objectType string, objectID int) ([]PluginFamily, error) {
	var result struct {
		Families []PluginFamily `json:"families"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/%d/families", objectType, objectID), &result)
	if err != nil {
		return nil, err
	}

	return result.Families, nil
}

// FamilyPlugins retrieves plugins in a family for a template.
func (e *EditorAPI) FamilyPlugins(ctx context.Context, objectType string, objectID, familyID int) ([]Plugin, error) {
	var result struct {
		Plugins []Plugin `json:"plugins"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/%d/families/%d", objectType, objectID, familyID), &result)
	if err != nil {
		return nil, err
	}

	return result.Plugins, nil
}

// PluginDetails retrieves details about a specific plugin.
func (e *EditorAPI) PluginDetails(ctx context.Context, objectType string, objectID, pluginID int) (*PluginDescription, error) {
	var result struct {
		PluginDescription PluginDescription `json:"plugindescription"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/%d/plugins/%d", objectType, objectID, pluginID), &result)
	if err != nil {
		return nil, err
	}

	return &result.PluginDescription, nil
}

// Audits retrieves audit files for a template.
func (e *EditorAPI) Audits(ctx context.Context, objectType string, objectID int) ([]interface{}, error) {
	var result struct {
		Audits []interface{} `json:"audits"`
	}

	_, err := e.client.Get(ctx, fmt.Sprintf("editor/%s/%d/audits", objectType, objectID), &result)
	if err != nil {
		return nil, err
	}

	return result.Audits, nil
}

