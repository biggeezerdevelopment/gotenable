package base

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 120 * time.Second
	// DefaultRetries is the default number of retries.
	DefaultRetries = 5
	// DefaultBackoff is the default backoff duration for rate limiting.
	DefaultBackoff = 1 * time.Second
	// Version is the SDK version.
	Version = "0.1.0"
)

// Client is the base HTTP client for Tenable APIs.
type Client struct {
	resty       *resty.Client
	baseURL     string
	basePath    string
	accessKey   string
	secretKey   string
	timeout     time.Duration
	retries     int
	backoff     time.Duration
	userAgent   string
	vendor      string
	product     string
	build       string
	authMech    string
	envBase     string
	lastReqUUID string
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// NewClient creates a new base Client with the given options.
func NewClient(envBase string, defaultURL string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		baseURL: defaultURL,
		timeout: DefaultTimeout,
		retries: DefaultRetries,
		backoff: DefaultBackoff,
		vendor:  "unknown",
		product: "unknown",
		build:   "unknown",
		envBase: envBase,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Check for environment variables if not set
	if c.baseURL == "" {
		if envURL := os.Getenv(fmt.Sprintf("%s_URL", envBase)); envURL != "" {
			c.baseURL = envURL
		} else {
			c.baseURL = defaultURL
		}
	}

	if c.accessKey == "" {
		c.accessKey = os.Getenv(fmt.Sprintf("%s_ACCESS_KEY", envBase))
	}
	if c.secretKey == "" {
		c.secretKey = os.Getenv(fmt.Sprintf("%s_SECRET_KEY", envBase))
	}

	// Validate URL
	if c.baseURL == "" {
		return nil, &ConnectionError{URL: "", Message: "no URL specified"}
	}

	// Build user agent
	c.userAgent = c.buildUserAgent()

	// Initialize resty client
	c.resty = resty.New().
		SetBaseURL(c.baseURL).
		SetTimeout(c.timeout).
		SetRetryCount(c.retries).
		SetRetryWaitTime(c.backoff).
		SetRetryMaxWaitTime(30*time.Second).
		SetHeader("User-Agent", c.userAgent).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on 429 (rate limit) and 5xx errors
			if err != nil {
				return true
			}
			return r.StatusCode() == http.StatusTooManyRequests ||
				(r.StatusCode() >= 500 && r.StatusCode() < 600)
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			// Store the request UUID for retry tracking
			if reqUUID := r.Header().Get("X-Tio-Last-Request-Uuid"); reqUUID != "" {
				c.SetHeader("X-Tio-Last-Request-Uuid", reqUUID)
			}
			return nil
		})

	// Set up authentication if keys are provided
	if c.accessKey != "" && c.secretKey != "" {
		c.setAPIKeyAuth()
	}

	return c, nil
}

// buildUserAgent constructs the User-Agent string.
func (c *Client) buildUserAgent() string {
	return fmt.Sprintf(
		"Integration/1.0 (%s; %s; Build/%s) goTenable/%s (Resty; Go/%s; %s/%s)",
		c.vendor,
		c.product,
		c.build,
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// setAPIKeyAuth sets up API key authentication headers.
func (c *Client) setAPIKeyAuth() {
	c.resty.SetHeader("X-APIKeys", fmt.Sprintf("accessKey=%s; secretKey=%s", c.accessKey, c.secretKey))
	c.authMech = "keys"
}

// IsAuthenticated returns true if the client has authentication configured.
func (c *Client) IsAuthenticated() bool {
	return c.authMech != ""
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Request creates a new request with the given method and path.
func (c *Client) Request(ctx context.Context) *resty.Request {
	return c.resty.R().SetContext(ctx)
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, result interface{}) (*resty.Response, error) {
	resp, err := c.Request(ctx).SetResult(result).Get(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Request(ctx).SetBody(body).SetResult(result).Post(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Request(ctx).SetBody(body).SetResult(result).Put(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*resty.Response, error) {
	resp, err := c.Request(ctx).Delete(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// GetWithParams performs a GET request with query parameters.
func (c *Client) GetWithParams(ctx context.Context, path string, params map[string]string, result interface{}) (*resty.Response, error) {
	resp, err := c.Request(ctx).SetQueryParams(params).SetResult(result).Get(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// PostWithParams performs a POST request with query parameters.
func (c *Client) PostWithParams(ctx context.Context, path string, params map[string]string, body interface{}, result interface{}) (*resty.Response, error) {
	resp, err := c.Request(ctx).SetQueryParams(params).SetBody(body).SetResult(result).Post(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	return resp, c.checkResponse(resp)
}

// Download performs a GET request and returns the response body as bytes.
func (c *Client) Download(ctx context.Context, path string) ([]byte, error) {
	resp, err := c.Request(ctx).Get(c.buildPath(path))
	if err != nil {
		return nil, &ConnectionError{URL: c.baseURL, Message: "request failed", Err: err}
	}
	if err := c.checkResponse(resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// buildPath constructs the full API path.
func (c *Client) buildPath(path string) string {
	if c.basePath != "" {
		return fmt.Sprintf("%s/%s", c.basePath, path)
	}
	return path
}

// checkResponse checks the response for errors.
func (c *Client) checkResponse(resp *resty.Response) error {
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode(),
		RequestID:  resp.Header().Get("X-Request-Uuid"),
		Response:   resp.Body(),
	}

	// Try to parse error message from response
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp.Body(), &errResp); err == nil {
		if errResp.Error != "" {
			apiErr.Message = errResp.Error
		} else if errResp.Message != "" {
			apiErr.Message = errResp.Message
		}
	}

	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode())
	}

	return apiErr
}

// SetBasePath sets the base path for API requests.
func (c *Client) SetBasePath(path string) {
	c.basePath = path
}

// Resty returns the underlying resty client for advanced usage.
func (c *Client) Resty() *resty.Client {
	return c.resty
}

// WithURL sets the base URL.
func WithURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithAPIKeys sets the API access and secret keys.
func WithAPIKeys(accessKey, secretKey string) ClientOption {
	return func(c *Client) {
		c.accessKey = accessKey
		c.secretKey = secretKey
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithRetries sets the number of retries.
func WithRetries(retries int) ClientOption {
	return func(c *Client) {
		c.retries = retries
	}
}

// WithBackoff sets the backoff duration for rate limiting.
func WithBackoff(backoff time.Duration) ClientOption {
	return func(c *Client) {
		c.backoff = backoff
	}
}

// WithVendor sets the vendor name for the User-Agent.
func WithVendor(vendor string) ClientOption {
	return func(c *Client) {
		c.vendor = vendor
	}
}

// WithProduct sets the product name for the User-Agent.
func WithProduct(product string) ClientOption {
	return func(c *Client) {
		c.product = product
	}
}

// WithBuild sets the build version for the User-Agent.
func WithBuild(build string) ClientOption {
	return func(c *Client) {
		c.build = build
	}
}

// WithBasePath sets the base path for API requests.
func WithBasePath(path string) ClientOption {
	return func(c *Client) {
		c.basePath = path
	}
}
