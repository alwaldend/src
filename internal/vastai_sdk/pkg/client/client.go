package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	generatedapi "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/internal/generated"
)

// ErrNotFound indicates a missing resource confirmed by its endpoint's response.
var ErrNotFound = errors.New("Vast.ai resource not found")

// ErrValueRedacted indicates an existing environment variable with a masked value.
var ErrValueRedacted = errors.New("Vast.ai environment value is redacted")

// Options configures a client's required endpoint and bearer credential.
type Options struct {
	// Endpoint is the REST API root, including its version path.
	Endpoint string
	// APIKey is the bearer credential, never included in diagnostics.
	APIKey string
}

// Client owns the API endpoint and a bounded HTTP client. It is concurrency safe.
type Client struct {
	// endpoint is the validated API root without a trailing slash.
	endpoint string
	// key is the bearer credential, never included in diagnostics.
	key string
	// http sends requests with timeout and redirect protection.
	http *http.Client
	// api supplies operations generated from the official pinned schema.
	api *generatedapi.APIClient
}

// New constructs a client; plaintext HTTP is allowed only for literal loopback.
func New(options Options) (*Client, error) {
	u, err := url.Parse(options.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse API endpoint: %w", err)
	}
	if u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("API endpoint must have a host and no credentials, query, or fragment")
	}
	ip := net.ParseIP(u.Hostname())
	if u.Scheme != "https" && !(u.Scheme == "http" && ip != nil && ip.IsLoopback()) {
		return nil, errors.New("API endpoint requires HTTPS, except literal loopback HTTP for testing")
	}
	if strings.TrimSpace(options.APIKey) == "" {
		return nil, errors.New("API key is required")
	}
	escapedPath := strings.TrimRight(u.EscapedPath(), "/")
	apiPath, err := url.PathUnescape(escapedPath)
	if err != nil {
		return nil, fmt.Errorf("decode API endpoint path: %w", err)
	}
	client := &Client{endpoint: strings.TrimRight(options.Endpoint, "/"), key: options.APIKey, http: &http.Client{
		Timeout:       60 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
	}}
	config := generatedapi.NewConfiguration()
	config.Servers = generatedapi.ServerConfigurations{{URL: u.Scheme + "://" + u.Host}}
	config.HTTPClient = &http.Client{
		Timeout:       client.http.Timeout,
		CheckRedirect: client.http.CheckRedirect,
		Transport:     &generatedTransport{path: apiPath, escapedPath: escapedPath, key: options.APIKey},
	}
	client.api = generatedapi.NewAPIClient(config)
	return client, nil
}

// request executes one request without retrying or exposing API response bodies.
func (c *Client) request(ctx context.Context, method, path string, body, result any) (resultErr error) {
	if err := c.requestPolicy(ctx, method, path, body, result, true); err != nil {
		return fmt.Errorf("Vast.ai request: %w", err)
	}
	return nil
}

// requestPolicy permits endpoint-specific role confirmations while retaining error checks.
func (c *Client) requestPolicy(ctx context.Context, method, path string, body, result any, requireSuccess bool) (resultErr error) {
	var input io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode %s %s request: %w", method, path, err)
		}
		input = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, input)
	if err != nil {
		return fmt.Errorf("construct %s %s request: %w", method, path, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send %s %s request (mutation outcome may be uncertain; inspect your account before retrying): %w", method, path, err)
	}
	if err := c.decodeResponse(method, path, resp, result, requireSuccess); err != nil {
		return fmt.Errorf("validate %s %s response: %w", method, path, err)
	}
	return nil
}

// decodeResponse validates raw endpoint contracts independently of optional schema fields.
func (c *Client) decodeResponse(method, path string, resp *http.Response, result any, requireSuccess bool) (resultErr error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close %s %s response: %w", method, path, err))
		}
	}()
	if resp.StatusCode == http.StatusNotFound && (strings.HasPrefix(path, "/instances/") || strings.HasPrefix(path, "/ssh/") || strings.HasPrefix(path, "/team/roles/") || (strings.HasPrefix(path, "/auth/apikeys/") && path != "/auth/apikeys/")) {
		return fmt.Errorf("%s %s: %w", method, path, ErrNotFound)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusBadRequest && method == http.MethodPost && path == "/team/" {
			var rejection struct {
				// Message is compared only against a documented constant, never rendered.
				Message string `json:"msg"`
			}
			if err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&rejection); err == nil && rejection.Message == "Team name already exists" {
				return errors.New("Vast.ai team name already exists; choose a different name")
			}
		}
		return fmt.Errorf("Vast.ai %s %s returned HTTP %d", method, path, resp.StatusCode)
	}
	const maximumBody = 8 << 20
	data, err := io.ReadAll(io.LimitReader(resp.Body, maximumBody+1))
	if err != nil {
		return fmt.Errorf("read %s %s response: %w", method, path, err)
	}
	if len(data) > maximumBody {
		return fmt.Errorf("Vast.ai %s %s response exceeds 8 MiB", method, path)
	}
	var envelope struct {
		// Success is optional on read operations and required for mutations.
		Success *bool `json:"success"`
		// Error identifies an application failure; its value is not disclosed.
		Error json.RawMessage `json:"error"`
	}
	if !json.Valid(data) {
		return fmt.Errorf("Vast.ai %s %s returned invalid JSON", method, path)
	}
	if bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
		if err := json.Unmarshal(data, &envelope); err != nil {
			return fmt.Errorf("decode %s %s response: %w", method, path, err)
		}
	}
	if (envelope.Success != nil && !*envelope.Success) || (len(envelope.Error) > 0 && string(envelope.Error) != "null" && string(envelope.Error) != `""`) {
		return fmt.Errorf("Vast.ai rejected %s %s; inspect account permissions and request parameters", method, path)
	}
	if requireSuccess && (method == http.MethodPut || method == http.MethodDelete) {
		if envelope.Success == nil || !*envelope.Success {
			return fmt.Errorf("Vast.ai %s %s omitted success confirmation", method, path)
		}
	}
	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("decode %s %s result: %w", method, path, err)
	}
	return nil
}

// CreateInput contains the supported on-demand rental parameters.
type CreateInput struct {
	// Image selects the Docker image.
	Image string `json:"image"`
	// Disk allocates local storage in GB.
	Disk float64 `json:"disk"`
	// Label names the rental.
	Label string `json:"label"`
	// Runtype selects the container launch mode.
	Runtype string `json:"runtype"`
	// State selects the initial power state.
	State string `json:"target_state"`
	// Onstart is the optional startup command.
	Onstart string `json:"onstart,omitempty"`
}

// Create accepts an offer and returns its new contract ID without waiting for boot.
func (c *Client) Create(ctx context.Context, offerID int64, input CreateInput) (int64, error) {
	var result struct {
		// ID is the newly allocated contract.
		ID int64 `json:"new_contract"`
	}
	body := generatedapi.CreateInstanceRequest{Image: input.Image, Disk: &input.Disk, Label: &input.Label, Runtype: &input.Runtype, TargetState: &input.State}
	if input.Onstart != "" {
		body.Onstart = &input.Onstart
	}
	request := c.api.InstancesAPI.CreateInstance(ctx, offerID).CreateInstanceRequest(body)
	if err := generated(c, http.MethodPut, fmt.Sprintf("/asks/%d/", offerID), &result, true, request.Execute); err != nil {
		return 0, fmt.Errorf("rent offer %d: %w", offerID, err)
	}
	if result.ID <= 0 {
		return 0, errors.New("rental accepted without a valid contract ID; inspect your account before retrying")
	}
	return result.ID, nil
}

// Instance contains only supported non-secret fields from an instance response.
type Instance struct {
	// ID is the contract identifier.
	ID int64 `json:"id"`
	// Label is the operator-assigned rental name.
	Label string `json:"label"`
	// State is the desired power state reported by the API.
	State string `json:"intended_status"`
	// Status is the observed container state.
	Status *string `json:"actual_status"`
	// Image is the current container image identifier.
	Image string `json:"image_uuid"`
	// Disk is the allocated disk size in GB.
	Disk float64 `json:"disk_space"`
	// Runtype is the container launch mode.
	Runtype string `json:"image_runtype"`
	// Onstart is the startup command, when returned.
	Onstart *string `json:"onstart"`
	// SSHHost is the SSH endpoint, when assigned.
	SSHHost *string `json:"ssh_host"`
	// SSHPort is the SSH port, when assigned.
	SSHPort *int64 `json:"ssh_port"`
	// PublicIP is the public address, when assigned.
	PublicIP *string `json:"public_ipaddr"`
}

// Get recognizes explicit null as deletion and rejects omitted or mismatched instances.
func (c *Client) Get(ctx context.Context, id int64) (*Instance, error) {
	var result struct {
		// Instance preserves the distinction between an omitted field and explicit null.
		Instance json.RawMessage `json:"instances"`
	}
	request := c.api.InstancesAPI.ShowInstance(ctx, id)
	if err := generated(c, http.MethodGet, fmt.Sprintf("/instances/%d/", id), &result, true, request.Execute); err != nil {
		return nil, fmt.Errorf("read instance %d: %w", id, err)
	}
	if bytes.Equal(bytes.TrimSpace(result.Instance), []byte("null")) {
		return nil, fmt.Errorf("read instance %d: %w", id, ErrNotFound)
	}
	if len(result.Instance) == 0 {
		return nil, fmt.Errorf("instance %d response omitted the instance", id)
	}
	var instance Instance
	if err := json.Unmarshal(result.Instance, &instance); err != nil {
		return nil, fmt.Errorf("decode instance %d: %w", id, err)
	}
	if instance.ID != id {
		return nil, fmt.Errorf("instance %d response omitted or mismatched the instance", id)
	}
	return &instance, nil
}

// Update changes an existing instance's mutable name and desired power state.
func (c *Client) Update(ctx context.Context, id int64, label, state string) error {
	input := generatedapi.ManageInstanceRequest{Label: &label, State: &state}
	var result struct{}
	request := c.api.InstancesAPI.ManageInstance(ctx, id).ManageInstanceRequest(input)
	if err := generated(c, http.MethodPut, fmt.Sprintf("/instances/%d/", id), &result, true, request.Execute); err != nil {
		return fmt.Errorf("update instance %d: %w", id, err)
	}
	return nil
}

// Delete destroys an instance; callers decide whether a missing instance is success.
func (c *Client) Delete(ctx context.Context, id int64) error {
	var result struct{}
	request := c.api.InstancesAPI.DestroyInstance(ctx, id)
	if err := generated(c, http.MethodDelete, fmt.Sprintf("/instances/%d/", id), &result, true, request.Execute); err != nil {
		return fmt.Errorf("destroy instance %d: %w", id, err)
	}
	return nil
}

// Offer contains the typed public fields exposed by the offer data source.
type Offer struct {
	// ID identifies the offer accepted by Create.
	ID int64 `json:"id"`
	// GPUName is the GPU model.
	GPUName string `json:"gpu_name"`
	// GPUs is the GPU count.
	GPUs int64 `json:"num_gpus"`
	// GPURAM is memory per GPU in MiB.
	GPURAM float64 `json:"gpu_ram"`
	// HourlyPrice is the API's total hourly offer price in USD.
	HourlyPrice float64 `json:"dph_total"`
}

// Search submits a JSON object query and preserves the API's offer ordering.
func (c *Client) Search(ctx context.Context, query map[string]any) ([]Offer, error) {
	var result struct {
		// Offers contains the marketplace results.
		Offers json.RawMessage `json:"offers"`
	}
	request := c.api.SearchAPI.SearchOffers(ctx).SearchOffersRequest(query)
	if err := generated(c, http.MethodPost, "/bundles/", &result, true, request.Execute); err != nil {
		return nil, fmt.Errorf("search offers: %w", err)
	}
	var offers []Offer
	if len(result.Offers) == 0 || string(result.Offers) == "null" {
		return nil, errors.New("offer response omitted the offers collection")
	}
	if err := json.Unmarshal(result.Offers, &offers); err != nil {
		return nil, fmt.Errorf("decode offer collection: %w", err)
	}
	for _, offer := range offers {
		if offer.ID <= 0 {
			return nil, errors.New("offer response contains an invalid offer ID")
		}
	}
	return offers, nil
}
