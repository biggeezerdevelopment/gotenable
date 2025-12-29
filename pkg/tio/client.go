// Package tio provides a Go client for the Tenable.io (Vulnerability Management) API.
package tio

import (
	"time"

	"github.com/biggeezerdevelopment/gotenable/pkg/base"
)

const (
	// DefaultURL is the default Tenable.io API URL.
	DefaultURL = "https://cloud.tenable.com"
	// EnvBase is the environment variable prefix for TIO configuration.
	EnvBase = "TIO"
)

// Client is the Tenable.io API client.
type Client struct {
	*base.Client

	// API endpoint interfaces
	AccessControl    *AccessControlAPI
	AgentConfig      *AgentConfigAPI
	AgentExclusions  *AgentExclusionsAPI
	AgentGroups      *AgentGroupsAPI
	Agents           *AgentsAPI
	Assets           *AssetsAPI
	AuditLog         *AuditLogAPI
	Credentials      *CredentialsAPI
	Editor           *EditorAPI
	Exclusions       *ExclusionsAPI
	Exports          *ExportsAPI
	Files            *FilesAPI
	Filters          *FiltersAPI
	Folders          *FoldersAPI
	Groups           *GroupsAPI
	Networks         *NetworksAPI
	Permissions      *PermissionsAPI
	Plugins          *PluginsAPI
	Policies         *PoliciesAPI
	RemediationScans *RemediationScansAPI
	ScannerGroups    *ScannerGroupsAPI
	Scanners         *ScannersAPI
	Scans            *ScansAPI
	Server           *ServerAPI
	Session          *SessionAPI
	Tags             *TagsAPI
	Users            *UsersAPI
	Workbenches      *WorkbenchesAPI

	// Cached timezone list
	timezones []string
}

// Option is a function that configures the TIO Client.
type Option func(*options)

type options struct {
	baseOpts []base.ClientOption
}

// New creates a new Tenable.io client.
func New(opts ...Option) (*Client, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	baseClient, err := base.NewClient(EnvBase, DefaultURL, o.baseOpts...)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client: baseClient,
	}

	// Initialize all API endpoints
	c.AccessControl = &AccessControlAPI{client: c}
	c.AgentConfig = &AgentConfigAPI{client: c}
	c.AgentExclusions = &AgentExclusionsAPI{client: c}
	c.AgentGroups = &AgentGroupsAPI{client: c}
	c.Agents = &AgentsAPI{client: c}
	c.Assets = &AssetsAPI{client: c}
	c.AuditLog = &AuditLogAPI{client: c}
	c.Credentials = &CredentialsAPI{client: c}
	c.Editor = &EditorAPI{client: c}
	c.Exclusions = &ExclusionsAPI{client: c}
	c.Exports = &ExportsAPI{client: c}
	c.Files = &FilesAPI{client: c}
	c.Filters = &FiltersAPI{client: c}
	c.Folders = &FoldersAPI{client: c}
	c.Groups = &GroupsAPI{client: c}
	c.Networks = &NetworksAPI{client: c}
	c.Permissions = &PermissionsAPI{client: c}
	c.Plugins = &PluginsAPI{client: c}
	c.Policies = &PoliciesAPI{client: c}
	c.RemediationScans = &RemediationScansAPI{client: c}
	c.ScannerGroups = &ScannerGroupsAPI{client: c}
	c.Scanners = &ScannersAPI{client: c}
	c.Scans = &ScansAPI{client: c}
	c.Server = &ServerAPI{client: c}
	c.Session = &SessionAPI{client: c}
	c.Tags = &TagsAPI{client: c}
	c.Users = &UsersAPI{client: c}
	c.Workbenches = &WorkbenchesAPI{client: c}

	return c, nil
}

// WithAPIKeys sets the API access and secret keys.
func WithAPIKeys(accessKey, secretKey string) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithAPIKeys(accessKey, secretKey))
	}
}

// WithURL sets the API base URL.
func WithURL(url string) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithURL(url))
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithTimeout(timeout))
	}
}

// WithRetries sets the number of retries.
func WithRetries(retries int) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithRetries(retries))
	}
}

// WithBackoff sets the backoff duration for rate limiting.
func WithBackoff(backoff time.Duration) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithBackoff(backoff))
	}
}

// WithVendor sets the vendor name for the User-Agent.
func WithVendor(vendor string) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithVendor(vendor))
	}
}

// WithProduct sets the product name for the User-Agent.
func WithProduct(product string) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithProduct(product))
	}
}

// WithBuild sets the build version for the User-Agent.
func WithBuild(build string) Option {
	return func(o *options) {
		o.baseOpts = append(o.baseOpts, base.WithBuild(build))
	}
}
