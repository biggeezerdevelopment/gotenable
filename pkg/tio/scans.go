package tio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/tenable/gotenable/pkg/base"
)

// ScansAPI handles scan-related operations.
type ScansAPI struct {
	client *Client
}

// Scan represents a scan configuration.
type Scan struct {
	ID                   int    `json:"id"`
	UUID                 string `json:"uuid"`
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Owner                string `json:"owner"`
	Enabled              bool   `json:"enabled"`
	FolderID             int    `json:"folder_id"`
	Read                 bool   `json:"read"`
	Status               string `json:"status"`
	Shared               bool   `json:"shared"`
	UserPermissions      int    `json:"user_permissions"`
	CreationDate         int64  `json:"creation_date"`
	LastModificationDate int64  `json:"last_modification_date"`
	Control              bool   `json:"control"`
	StartTime            string `json:"starttime,omitempty"`
	Timezone             string `json:"timezone,omitempty"`
	RRules               string `json:"rrules,omitempty"`
}

// ScanDetails represents detailed scan information.
type ScanDetails struct {
	Info            ScanInfo         `json:"info"`
	Hosts           []ScanHost       `json:"hosts,omitempty"`
	Comphosts       []ScanHost       `json:"comphosts,omitempty"`
	Notes           []ScanNote       `json:"notes,omitempty"`
	Remediations    ScanRemediations `json:"remediations,omitempty"`
	Vulnerabilities []ScanVuln       `json:"vulnerabilities,omitempty"`
	Compliance      []ScanCompliance `json:"compliance,omitempty"`
	History         []ScanHistory    `json:"history,omitempty"`
	Filters         []ScanFilter     `json:"filters,omitempty"`
}

// ScanInfo contains scan metadata.
type ScanInfo struct {
	EditAllowed   bool   `json:"edit_allowed"`
	Status        string `json:"status"`
	Policy        string `json:"policy"`
	PCICanUpload  bool   `json:"pci-can-upload"`
	HasAuditTrail bool   `json:"hasaudittrail"`
	ScannerName   string `json:"scanner_name"`
	Control       bool   `json:"control"`
	Name          string `json:"name"`
	ObjectID      int    `json:"object_id"`
	NoTarget      bool   `json:"no_target"`
	UUID          string `json:"uuid"`
	HostCount     int    `json:"hostcount"`
	ScannerStart  int64  `json:"scanner_start,omitempty"`
	ScannerEnd    int64  `json:"scanner_end,omitempty"`
	Timestamp     int64  `json:"timestamp"`
	Targets       string `json:"targets"`
	ACLs          []ACL  `json:"acls,omitempty"`
	ScheduleUUID  string `json:"schedule_uuid,omitempty"`
	IsArchived    bool   `json:"is_archived"`
	FolderID      int    `json:"folder_id"`
	Haskb         bool   `json:"haskb"`
	ScanType      string `json:"scan_type"`
}

// ACL represents an access control list entry.
type ACL struct {
	ID          int    `json:"id"`
	Owner       int    `json:"owner"`
	Type        string `json:"type"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// ScanHost represents a host in scan results.
type ScanHost struct {
	HostID              int    `json:"host_id"`
	HostIndex           int    `json:"host_index"`
	Hostname            string `json:"hostname"`
	Progress            string `json:"progress"`
	Critical            int    `json:"critical"`
	High                int    `json:"high"`
	Medium              int    `json:"medium"`
	Low                 int    `json:"low"`
	Info                int    `json:"info"`
	TotalChecksCount    int    `json:"totalchecksconsidered"`
	NumChecksCount      int    `json:"numchecksconsidered"`
	ScanProgressTotal   int    `json:"scanprogresstotal"`
	ScanProgressCurrent int    `json:"scanprogresscurrent"`
	Score               int    `json:"score"`
}

// ScanNote represents a note in scan results.
type ScanNote struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Severity int    `json:"severity"`
}

// ScanRemediations represents remediation information.
type ScanRemediations struct {
	Remediations []Remediation `json:"remediations,omitempty"`
	NumHosts     int           `json:"num_hosts"`
	NumCVEs      int           `json:"num_cves"`
	NumImpacted  int           `json:"num_impacted_hosts"`
	NumPlugins   int           `json:"num_remediated_cves"`
}

// Remediation represents a single remediation.
type Remediation struct {
	Value       string `json:"value"`
	Remediation string `json:"remediation"`
	Hosts       int    `json:"hosts"`
	Vulns       int    `json:"vulns"`
}

// ScanVuln represents a vulnerability in scan results.
type ScanVuln struct {
	PluginID      int    `json:"plugin_id"`
	PluginName    string `json:"plugin_name"`
	PluginFamily  string `json:"plugin_family"`
	Count         int    `json:"count"`
	VulnIndex     int    `json:"vuln_index"`
	SeverityIndex int    `json:"severity_index"`
}

// ScanCompliance represents compliance information.
type ScanCompliance struct {
	PluginID      int    `json:"plugin_id"`
	PluginName    string `json:"plugin_name"`
	PluginFamily  string `json:"plugin_family"`
	Count         int    `json:"count"`
	SeverityIndex int    `json:"severity_index"`
}

// ScanHistory represents a scan history entry.
type ScanHistory struct {
	HistoryID            int    `json:"history_id"`
	UUID                 string `json:"uuid"`
	OwnerID              int    `json:"owner_id"`
	Status               string `json:"status"`
	CreationDate         int64  `json:"creation_date"`
	LastModificationDate int64  `json:"last_modification_date"`
	Type                 string `json:"type"`
	IsArchived           bool   `json:"is_archived"`
}

// ScanFilter represents an available filter.
type ScanFilter struct {
	Name         string   `json:"name"`
	ReadableName string   `json:"readable_name"`
	Control      Control  `json:"control"`
	Operators    []string `json:"operators"`
}

// Control represents filter control options.
type Control struct {
	Type          string   `json:"type"`
	ReadableRegex string   `json:"readable_regex,omitempty"`
	Regex         string   `json:"regex,omitempty"`
	List          []string `json:"list,omitempty"`
}

// ScanCreateRequest represents a request to create a scan.
type ScanCreateRequest struct {
	UUID        string                 `json:"uuid,omitempty"`
	Settings    ScanSettings           `json:"settings"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	Plugins     map[string]interface{} `json:"plugins,omitempty"`
}

// ScanSettings represents scan settings.
type ScanSettings struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	FolderID    int    `json:"folder_id,omitempty"`
	ScannerID   string `json:"scanner_id,omitempty"`
	PolicyID    int    `json:"policy_id,omitempty"`
	TextTargets string `json:"text_targets,omitempty"`
	FileTargets string `json:"file_targets,omitempty"`
	Emails      string `json:"emails,omitempty"`
	Enabled     bool   `json:"enabled"`
	Launch      string `json:"launch,omitempty"`
	RRules      string `json:"rrules,omitempty"`
	Starttime   string `json:"starttime,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	ACLs        []ACL  `json:"acls,omitempty"`
}

// ScanLaunchResponse represents the response from launching a scan.
type ScanLaunchResponse struct {
	ScanUUID string `json:"scan_uuid"`
}

// ScanExportRequest represents an export request.
type ScanExportRequest struct {
	Format     string   `json:"format"`
	Password   string   `json:"password,omitempty"`
	Chapters   []string `json:"chapters,omitempty"`
	FilterType string   `json:"filter.search_type,omitempty"`
}

// ScanListOptions contains options for listing scans.
type ScanListOptions struct {
	FolderID             *int
	LastModificationDate *time.Time
}

// List retrieves all configured scans.
func (s *ScansAPI) List(ctx context.Context, opts *ScanListOptions) ([]Scan, error) {
	params := make(map[string]string)
	if opts != nil {
		if opts.FolderID != nil {
			params["folder_id"] = strconv.Itoa(*opts.FolderID)
		}
		if opts.LastModificationDate != nil {
			params["last_modification_date"] = strconv.FormatInt(opts.LastModificationDate.Unix(), 10)
		}
	}

	var result struct {
		Scans []Scan `json:"scans"`
	}

	_, err := s.client.GetWithParams(ctx, "scans", params, &result)
	if err != nil {
		return nil, err
	}

	return result.Scans, nil
}

// Create creates a new scan.
func (s *ScansAPI) Create(ctx context.Context, req *ScanCreateRequest) (*Scan, error) {
	var result struct {
		Scan Scan `json:"scan"`
	}

	_, err := s.client.Post(ctx, "scans", req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Scan, nil
}

// Details retrieves detailed information about a scan.
func (s *ScansAPI) Details(ctx context.Context, scanID int) (*ScanDetails, error) {
	var result ScanDetails
	_, err := s.client.Get(ctx, fmt.Sprintf("scans/%d", scanID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Configure updates a scan configuration.
func (s *ScansAPI) Configure(ctx context.Context, scanID int, req *ScanCreateRequest) (*Scan, error) {
	var result struct {
		Scan Scan `json:"scan"`
	}

	_, err := s.client.Put(ctx, fmt.Sprintf("scans/%d", scanID), req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Scan, nil
}

// Delete removes a scan.
func (s *ScansAPI) Delete(ctx context.Context, scanID int) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("scans/%d", scanID))
	return err
}

// Copy duplicates a scan.
func (s *ScansAPI) Copy(ctx context.Context, scanID int, name string, folderID *int) (*Scan, error) {
	payload := map[string]interface{}{}
	if name != "" {
		payload["name"] = name
	}
	if folderID != nil {
		payload["folder_id"] = *folderID
	}

	var result Scan
	_, err := s.client.Post(ctx, fmt.Sprintf("scans/%d/copy", scanID), payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Launch starts a scan.
func (s *ScansAPI) Launch(ctx context.Context, scanID int, targets []string) (string, error) {
	payload := map[string]interface{}{}
	if len(targets) > 0 {
		payload["alt_targets"] = targets
	}

	var result ScanLaunchResponse
	_, err := s.client.Post(ctx, fmt.Sprintf("scans/%d/launch", scanID), payload, &result)
	if err != nil {
		return "", err
	}

	return result.ScanUUID, nil
}

// Pause pauses a running scan.
func (s *ScansAPI) Pause(ctx context.Context, scanID int) error {
	_, err := s.client.Post(ctx, fmt.Sprintf("scans/%d/pause", scanID), nil, nil)
	return err
}

// Resume resumes a paused scan.
func (s *ScansAPI) Resume(ctx context.Context, scanID int) error {
	_, err := s.client.Post(ctx, fmt.Sprintf("scans/%d/resume", scanID), nil, nil)
	return err
}

// Stop stops a running scan.
func (s *ScansAPI) Stop(ctx context.Context, scanID int) error {
	_, err := s.client.Post(ctx, fmt.Sprintf("scans/%d/stop", scanID), nil, nil)
	return err
}

// Status returns the status of the latest scan instance.
func (s *ScansAPI) Status(ctx context.Context, scanID int) (string, error) {
	var result struct {
		Status string `json:"status"`
	}
	_, err := s.client.Get(ctx, fmt.Sprintf("scans/%d/latest-status", scanID), &result)
	if err != nil {
		return "", err
	}
	return result.Status, nil
}

// History retrieves scan history.
func (s *ScansAPI) History(ctx context.Context, scanID int, limit, offset int) *base.Iterator[ScanHistory] {
	fetcher := func(ctx context.Context, offset, limit int) (json.RawMessage, *base.PaginationInfo, error) {
		params := map[string]string{
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		var result struct {
			History    []ScanHistory `json:"history"`
			Pagination struct {
				Total  int `json:"total"`
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"pagination"`
		}

		_, err := s.client.GetWithParams(ctx, fmt.Sprintf("scans/%d/history", scanID), params, &result)
		if err != nil {
			return nil, nil, err
		}

		data, _ := json.Marshal(result.History)
		return data, &base.PaginationInfo{
			Total:  result.Pagination.Total,
			Limit:  result.Pagination.Limit,
			Offset: result.Pagination.Offset,
		}, nil
	}

	transformer := func(data json.RawMessage) ([]ScanHistory, error) {
		var items []ScanHistory
		err := json.Unmarshal(data, &items)
		return items, err
	}

	return base.NewIterator(ctx, fetcher, transformer,
		base.WithLimit[ScanHistory](limit),
		base.WithOffset[ScanHistory](offset),
	)
}

// DeleteHistory removes a scan history instance.
func (s *ScansAPI) DeleteHistory(ctx context.Context, scanID, historyID int) error {
	_, err := s.client.Delete(ctx, fmt.Sprintf("scans/%d/history/%d", scanID, historyID))
	return err
}

// Export initiates a scan export.
func (s *ScansAPI) Export(ctx context.Context, scanID int, format string, historyID *int, chapters []string) (io.Reader, error) {
	params := make(map[string]string)
	if historyID != nil {
		params["history_id"] = strconv.Itoa(*historyID)
	}

	payload := map[string]interface{}{
		"format": format,
	}
	if len(chapters) > 0 {
		payload["chapters"] = chapters
	}

	// Initiate export
	var exportResp struct {
		File int `json:"file"`
	}
	_, err := s.client.PostWithParams(ctx, fmt.Sprintf("scans/%d/export", scanID), params, payload, &exportResp)
	if err != nil {
		return nil, err
	}

	fileID := exportResp.File

	// Wait for export to be ready
	for {
		var statusResp struct {
			Status string `json:"status"`
		}
		_, err := s.client.Get(ctx, fmt.Sprintf("scans/%d/export/%d/status", scanID, fileID), &statusResp)
		if err != nil {
			return nil, err
		}

		if statusResp.Status == "ready" {
			break
		}
		if statusResp.Status == "error" {
			return nil, &base.FileDownloadError{
				Resource:   "scans",
				ResourceID: strconv.Itoa(scanID),
				Filename:   strconv.Itoa(fileID),
			}
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2500 * time.Millisecond):
		}
	}

	// Download the file
	data, err := s.client.Download(ctx, fmt.Sprintf("scans/%d/export/%d/download", scanID, fileID))
	if err != nil {
		return nil, err
	}

	return &bytesReader{data: data}, nil
}

// bytesReader wraps bytes to implement io.Reader.
type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// HostDetails retrieves host details from a scan.
func (s *ScansAPI) HostDetails(ctx context.Context, scanID, hostID int, historyID *int) (*ScanHostDetails, error) {
	params := make(map[string]string)
	if historyID != nil {
		params["history_id"] = strconv.Itoa(*historyID)
	}

	var result ScanHostDetails
	_, err := s.client.GetWithParams(ctx, fmt.Sprintf("scans/%d/hosts/%d", scanID, hostID), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanHostDetails contains detailed host information.
type ScanHostDetails struct {
	Info            HostInfo         `json:"info"`
	Vulnerabilities []HostVuln       `json:"vulnerabilities,omitempty"`
	Compliance      []HostCompliance `json:"compliance,omitempty"`
}

// HostInfo contains host metadata.
type HostInfo struct {
	HostStart       string `json:"host_start"`
	HostEnd         string `json:"host_end"`
	HostFQDN        string `json:"host-fqdn,omitempty"`
	HostIP          string `json:"host-ip"`
	MacAddress      string `json:"mac-address,omitempty"`
	NetbiosName     string `json:"netbios-name,omitempty"`
	OperatingSystem string `json:"operating-system,omitempty"`
}

// HostVuln represents a vulnerability on a host.
type HostVuln struct {
	PluginID     int    `json:"plugin_id"`
	PluginName   string `json:"plugin_name"`
	PluginFamily string `json:"plugin_family"`
	Count        int    `json:"count"`
	Severity     int    `json:"severity"`
}

// HostCompliance represents compliance on a host.
type HostCompliance struct {
	PluginID     int    `json:"plugin_id"`
	PluginName   string `json:"plugin_name"`
	PluginFamily string `json:"plugin_family"`
	Count        int    `json:"count"`
	Severity     int    `json:"severity"`
}

// PluginOutput retrieves plugin output for a specific plugin on a host.
func (s *ScansAPI) PluginOutput(ctx context.Context, scanID, hostID, pluginID int, historyID *int) (*PluginOutputResponse, error) {
	params := make(map[string]string)
	if historyID != nil {
		params["history_id"] = strconv.Itoa(*historyID)
	}

	var result PluginOutputResponse
	_, err := s.client.GetWithParams(ctx, fmt.Sprintf("scans/%d/hosts/%d/plugins/%d", scanID, hostID, pluginID), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// PluginOutputResponse contains plugin output data.
type PluginOutputResponse struct {
	Info    PluginInfo     `json:"info"`
	Outputs []PluginOutput `json:"outputs"`
}

// PluginInfo contains plugin metadata.
type PluginInfo struct {
	PluginDetails ScanPluginDetails `json:"plugindescription"`
}

// ScanPluginDetails contains detailed plugin information from scan results.
type ScanPluginDetails struct {
	PluginID     int               `json:"pluginid"`
	PluginName   string            `json:"pluginname"`
	PluginFamily string            `json:"pluginfamily"`
	Severity     int               `json:"severity"`
	PluginAttrs  []PluginAttribute `json:"pluginattributes"`
}

// PluginAttribute represents a plugin attribute.
type PluginAttribute struct {
	Name  string `json:"attribute_name"`
	Value string `json:"attribute_value"`
}

// PluginOutput represents plugin output.
type PluginOutput struct {
	PluginOutput string            `json:"plugin_output"`
	Hosts        string            `json:"hosts,omitempty"`
	Ports        map[string][]Port `json:"ports,omitempty"`
	Severity     int               `json:"severity"`
}

// Port represents a port.
type Port struct {
	Hostname string `json:"hostname"`
}

// Timezones retrieves the list of valid timezones.
func (s *ScansAPI) Timezones(ctx context.Context) ([]string, error) {
	var result struct {
		Timezones []struct {
			Value string `json:"value"`
		} `json:"timezones"`
	}

	_, err := s.client.Get(ctx, "scans/timezones", &result)
	if err != nil {
		return nil, err
	}

	tzs := make([]string, len(result.Timezones))
	for i, tz := range result.Timezones {
		tzs[i] = tz.Value
	}
	return tzs, nil
}

// Schedule enables or disables scan scheduling.
func (s *ScansAPI) Schedule(ctx context.Context, scanID int, enabled bool) error {
	payload := map[string]bool{"enabled": enabled}
	_, err := s.client.Put(ctx, fmt.Sprintf("scans/%d/schedule", scanID), payload, nil)
	return err
}
