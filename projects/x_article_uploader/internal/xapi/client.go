package xapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
)

// DefaultBaseURL is the X API v2 endpoint.
const DefaultBaseURL = "https://api.x.com/2"

// Client performs signed requests against the X API.
type Client struct {
	BaseURL     string
	Credentials Credentials
	HTTPClient  *http.Client
	// Now supplies the signing timestamp; tests set it for a stable request.
	Now func() time.Time
}

// New returns a client for the given credentials.
func New(credentials Credentials) *Client {
	return &Client{
		BaseURL:     DefaultBaseURL,
		Credentials: credentials,
		HTTPClient:  &http.Client{Timeout: 60 * time.Second},
		Now:         time.Now,
	}
}

// APIError is a rejected request. It carries the failing operation, the status,
// and the API's own error body, and it is never retried automatically.
type APIError struct {
	Operation string
	Status    int
	Body      string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: X API rejected the request with status %d: %s", e.Operation, e.Status, e.Body)
}

// MediaUpload is one media item referenced by an image entity.
//
// The v2 media upload endpoint returns its result inside a `data` envelope and
// names the identifier `id`, not `media_id`; the payload carries the identifier
// and the lifetime after which the uploaded media expires.
type MediaUpload struct {
	MediaID          string `json:"id"`
	MediaCategory    string `json:"media_category"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
}

// mediaUploadResponse is the v2 envelope the media upload endpoint returns.
type mediaUploadResponse struct {
	Data MediaUpload `json:"data"`
}

// Draft is the response to a draft creation request.
type Draft struct {
	ID string `json:"id"`
}

// draftResponse is the v2 envelope the draft endpoint returns; the identifier
// is readable only inside `data`, not at the top level.
type draftResponse struct {
	Data Draft `json:"data"`
}

// CreateDraft sends the title and content_state to the draft endpoint. It does
// not publish anything.
func (c *Client) CreateDraft(title string, contentState any) (*Draft, error) {
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("create draft: the artifact carries no parsed title")
	}
	payload := map[string]any{"title": title, "content_state": contentState}
	var envelope draftResponse
	if err := c.doJSON(http.MethodPost, "/articles/draft", payload, &envelope); err != nil {
		return nil, err
	}
	if strings.TrimSpace(envelope.Data.ID) == "" {
		return nil, fmt.Errorf("create draft: the API returned no draft id")
	}
	return &envelope.Data, nil
}

// UploadImage uploads image bytes to the media endpoint and returns the media
// identifier the image entity must reference. The bytes are supplied by the
// caller rather than read here, so the content the digest was verified against
// is the content the request carries even if the file changes on disk between
// verification and upload.
func (c *Client) UploadImage(name string, content []byte) (*MediaUpload, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("media", filepath.Base(name))
	if err != nil {
		return nil, fmt.Errorf("build media upload body: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return nil, fmt.Errorf("write media upload body: %w", err)
	}
	// The upload endpoint requires media_category alongside the file, and
	// rejects the request when it is missing. Every image this project uploads
	// is referenced by a tweet-style image entity, so the category is constant.
	if err := writer.WriteField("media_category", draftjs.MediaCategoryImage); err != nil {
		return nil, fmt.Errorf("write media upload category: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close media upload body: %w", err)
	}

	endpoint, err := url.Parse(c.BaseURL + "/media/upload")
	if err != nil {
		return nil, fmt.Errorf("parse media endpoint: %w", err)
	}
	authorization, err := c.Credentials.Sign(http.MethodPost, endpoint.String(), nil, nil, c.Now())
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("build media request: %w", err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", authorization)

	var envelope mediaUploadResponse
	if err := c.roundTrip("upload image", request, &envelope); err != nil {
		return nil, err
	}
	if strings.TrimSpace(envelope.Data.MediaID) == "" {
		return nil, fmt.Errorf("upload image: the API returned no media id")
	}
	return &envelope.Data, nil
}

// doJSON sends a JSON request to a relative API path.
func (c *Client) doJSON(method, path string, payload, into any) error {
	endpoint, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("parse endpoint: %w", err)
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode request body: %w", err)
	}
	authorization, err := c.Credentials.Sign(method, endpoint.String(), nil, nil, c.Now())
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, endpoint.String(), bytes.NewReader(encoded))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authorization)
	return c.roundTrip(strings.TrimPrefix(path, "/"), request, into)
}

// roundTrip performs the request and decodes the response, reporting a rejected
// request as an APIError.
func (c *Client) roundTrip(operation string, request *http.Request, into any) error {
	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%s: read response: %w", operation, err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &APIError{Operation: operation, Status: response.StatusCode, Body: strings.TrimSpace(string(body))}
	}
	if into == nil {
		return nil
	}
	if err := json.Unmarshal(body, into); err != nil {
		return fmt.Errorf("%s: decode response: %w", operation, err)
	}
	return nil
}
