package tio

import (
	"context"
)

// ServerAPI handles server operations.
type ServerAPI struct {
	client *Client
}

// ServerStatus represents the server status.
type ServerStatus struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

// ServerProperties represents server properties.
type ServerProperties struct {
	Analytics          ServerAnalytics      `json:"analytics"`
	Capabilities       ServerCapabilities   `json:"capabilities"`
	ContainerDBVersion string               `json:"containerdbversion,omitempty"`
	Enterprise         bool                 `json:"enterprise"`
	Evaluation         bool                 `json:"evaluation"`
	Expiration         int64                `json:"expiration"`
	ExpirationTime     int64                `json:"expiration_time"`
	IdleTimeout        int                  `json:"idle_timeout"`
	License            ServerLicense        `json:"license"`
	LoadedPluginSet    string               `json:"loaded_plugin_set"`
	LoginBanner        string               `json:"login_banner,omitempty"`
	NessusType         string               `json:"nessus_type"`
	NessusUIBuild      string               `json:"nessus_ui_build"`
	NessusUIVersion    string               `json:"nessus_ui_version"`
	Notifications      []ServerNotification `json:"notifications,omitempty"`
	PluginSet          string               `json:"plugin_set"`
	RestrictedFeatures bool                 `json:"restricted_features"`
	ScannerBootTime    int64                `json:"scanner_boottime"`
	ServerBuild        string               `json:"server_build"`
	ServerUUID         string               `json:"server_uuid"`
	ServerVersion      string               `json:"server_version"`
	ShowNPV            bool                 `json:"show_npv"`
	Update             ServerUpdate         `json:"update,omitempty"`
}

// ServerAnalytics represents analytics settings.
type ServerAnalytics struct {
	Enabled bool `json:"enabled"`
}

// ServerCapabilities represents server capabilities.
type ServerCapabilities struct {
	MultiScanner   bool          `json:"multi_scanner"`
	Reports        bool          `json:"reports"`
	MultiUser      bool          `json:"multi_user"`
	TwoFactor      TwoFactorCaps `json:"two_factor"`
	ScanVulnGroups bool          `json:"scan_vuln_groups"`
}

// TwoFactorCaps represents two-factor capabilities.
type TwoFactorCaps struct {
	SMS   bool `json:"sms"`
	Email bool `json:"email"`
}

// ServerLicense represents license information.
type ServerLicense struct {
	Type       string `json:"type"`
	Agents     int    `json:"agents"`
	IPS        int    `json:"ips"`
	Scanners   int    `json:"scanners"`
	Users      int    `json:"users"`
	Apps       int    `json:"apps,omitempty"`
	Expiration int64  `json:"expiration_date"`
}

// ServerNotification represents a server notification.
type ServerNotification struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// ServerUpdate represents update information.
type ServerUpdate struct {
	Href      string `json:"href,omitempty"`
	NewUpdate bool   `json:"new_update"`
	Restart   bool   `json:"restart"`
}

// Status retrieves the server status.
func (s *ServerAPI) Status(ctx context.Context) (*ServerStatus, error) {
	var result ServerStatus
	_, err := s.client.Get(ctx, "server/status", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Properties retrieves server properties.
func (s *ServerAPI) Properties(ctx context.Context) (*ServerProperties, error) {
	var result ServerProperties
	_, err := s.client.Get(ctx, "server/properties", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
