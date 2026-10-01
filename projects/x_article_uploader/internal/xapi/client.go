package xapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
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
	// MediaProcessingTimeout bounds asynchronous media processing; zero uses
	// the default two-minute deadline.
	MediaProcessingTimeout time.Duration
	// Now supplies the signing timestamp; tests set it for a stable request.
	Now func() time.Time
}

// New returns a client for the given credentials.
func New(credentials Credentials) *Client {
	return &Client{
		BaseURL:     DefaultBaseURL,
		Credentials: credentials,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
			// A redirect must not spend another API request. Preserve the
			// original response so its status and rate headers reach the caller.
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		Now: time.Now,
	}
}

// MediaCacheScope binds cached identifiers to the media endpoint, category, and
// exact OAuth credential set without persisting credential material.
func (c *Client) MediaCacheScope() string {
	if err := c.Credentials.Validate(); err != nil {
		return ""
	}
	encoded, err := json.Marshal([]string{
		"x-article-media-cache-v1", c.BaseURL + "/media/upload", draftjs.MediaCategoryImage,
		c.Credentials.APIKey, c.Credentials.APISecret,
		c.Credentials.AccessToken, c.Credentials.AccessTokenSecret,
	})
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(encoded))
}

// APIError is a rejected request. It carries the failing operation, the status,
// the API's own error body, and the allowed rate-limit response headers. It is
// never retried automatically.
type APIError struct {
	Operation string
	Status    int
	Body      string
	Headers   http.Header
	// ReceivedAt anchors relative Retry-After values and reset freshness.
	ReceivedAt time.Time
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: X API rejected the request with status %d: %s", e.Operation, e.Status, e.Body) + e.rateLimitDetails()
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
	// UploadedAt anchors the returned lifetime conservatively to the beginning
	// of the upload request, including time spent awaiting media processing.
	UploadedAt     time.Time            `json:"-"`
	ProcessingInfo *MediaProcessingInfo `json:"processing_info,omitempty"`
}

// MediaProcessingInfo describes whether an uploaded image can be attached yet.
type MediaProcessingInfo struct {
	State          string `json:"state"`
	CheckAfterSecs int64  `json:"check_after_secs"`
}

// mediaUploadResponse is the v2 envelope the media upload endpoint returns.
type mediaUploadResponse struct {
	Data   MediaUpload       `json:"data"`
	Errors []json.RawMessage `json:"errors,omitempty"`
}

func (r mediaUploadResponse) checkErrors() error {
	if len(r.Errors) == 0 {
		return nil
	}
	details := make([]string, 0, len(r.Errors))
	for _, problem := range r.Errors {
		details = append(details, string(problem))
	}
	return fmt.Errorf("media response contains errors: %s", strings.Join(details, "; "))
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

// CoverMedia references an image already uploaded through the media endpoint.
type CoverMedia struct {
	MediaCategory string `json:"media_category"`
	MediaID       string `json:"media_id"`
}

// DraftRequest is the draft endpoint's payload. A banner is optional and stays
// separate from the body document's image entities.
type DraftRequest struct {
	Title        string      `json:"title"`
	ContentState any         `json:"content_state"`
	CoverMedia   *CoverMedia `json:"cover_media,omitempty"`
}

// CreateDraft sends the title, content_state, and optional cover_media to the
// draft endpoint. It does not publish anything.
func (c *Client) CreateDraft(request DraftRequest) (*Draft, error) {
	if strings.TrimSpace(request.Title) == "" {
		return nil, fmt.Errorf("create draft: the artifact carries no parsed title")
	}
	var envelope draftResponse
	if err := c.doJSON(http.MethodPost, "/articles/draft", request, &envelope); err != nil {
		return nil, fmt.Errorf("send draft request: %w", err)
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
	// A buffered body otherwise permits automatic transport replay. Media
	// uploads are mutations, so a consumed body must not be sent again.
	request.GetBody = nil
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", authorization)

	var envelope mediaUploadResponse
	uploadedAt := time.Now()
	if err := c.roundTrip("upload image", request, &envelope); err != nil {
		return nil, err
	}
	if err := envelope.checkErrors(); err != nil {
		return nil, fmt.Errorf("upload image: %w", err)
	}
	if strings.TrimSpace(envelope.Data.MediaID) == "" {
		return nil, fmt.Errorf("upload image: the API returned no media id")
	}
	envelope.Data.UploadedAt = uploadedAt
	if err := c.awaitMediaProcessing(&envelope.Data); err != nil {
		return nil, fmt.Errorf("upload image: %w", err)
	}
	return &envelope.Data, nil
}

// awaitMediaProcessing polls only asynchronous uploads. Rejected status calls
// are returned immediately; the original upload POST is never repeated.
func (c *Client) awaitMediaProcessing(upload *MediaUpload) error {
	if upload.ProcessingInfo == nil {
		return nil
	}
	timeout := c.MediaProcessingTimeout
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		info := upload.ProcessingInfo
		switch info.State {
		case "succeeded":
			return nil
		case "failed":
			return fmt.Errorf("media processing failed for %q", upload.MediaID)
		case "pending", "in_progress":
			if info.CheckAfterSecs < 0 || info.CheckAfterSecs > (1<<63-1)/int64(time.Second) {
				return fmt.Errorf("media processing returned an invalid check_after_secs: %d", info.CheckAfterSecs)
			}
			wait := time.Duration(info.CheckAfterSecs) * time.Second
			if wait < time.Second {
				wait = time.Second
			}
			timer := time.NewTimer(wait)
			select {
			case <-ctx.Done():
				timer.Stop()
				return fmt.Errorf("await media processing: %w", ctx.Err())
			case <-timer.C:
			}
			status, err := c.mediaStatus(ctx, upload.MediaID)
			if err != nil {
				return fmt.Errorf("await media processing: %w", err)
			}
			if status.MediaID != "" && status.MediaID != upload.MediaID {
				return fmt.Errorf("media processing returned a different identifier: %q", status.MediaID)
			}
			if status.ProcessingInfo == nil {
				return fmt.Errorf("media processing status omitted processing_info")
			}
			upload.ProcessingInfo = status.ProcessingInfo
		default:
			return fmt.Errorf("media processing returned an unknown state: %q", info.State)
		}
	}
}

// mediaStatus signs the complete status query. JSON and multipart POST bodies
// are excluded from OAuth parameters, but GET query parameters must be included.
func (c *Client) mediaStatus(ctx context.Context, mediaID string) (*MediaUpload, error) {
	endpoint, err := url.Parse(c.BaseURL + "/media/upload")
	if err != nil {
		return nil, fmt.Errorf("parse media status endpoint: %w", err)
	}
	query := url.Values{"command": {"STATUS"}, "media_id": {mediaID}}
	endpoint.RawQuery = query.Encode()
	authorization, err := c.Credentials.Sign(http.MethodGet, endpoint.String(), query, nil, c.Now())
	if err != nil {
		return nil, fmt.Errorf("sign media status request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build media status request: %w", err)
	}
	request.Header.Set("Authorization", authorization)
	var envelope mediaUploadResponse
	if err := c.roundTrip("media status", request, &envelope); err != nil {
		return nil, fmt.Errorf("read media status: %w", err)
	}
	if err := envelope.checkErrors(); err != nil {
		return nil, fmt.Errorf("read media status: %w", err)
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
	// Draft creation is a mutation. Refuse automatic replay of its body after
	// transport errors, including HTTP/2 stream failures.
	request.GetBody = nil
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
	receivedAt := time.Now()
	if c.Now != nil {
		receivedAt = c.Now()
	}
	headers := copyRateLimitHeaders(response.Header)
	body, readErr := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	closeErr := response.Body.Close()
	if readErr != nil {
		readErr = fmt.Errorf("read response: %w", readErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("close response: %w", closeErr)
	}
	bodyErr := errors.Join(readErr, closeErr)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		apiErr := &APIError{
			Operation: operation, Status: response.StatusCode, Body: strings.TrimSpace(string(body)),
			Headers: headers, ReceivedAt: receivedAt,
		}
		if bodyErr != nil {
			return fmt.Errorf("rejected response: %w", errors.Join(apiErr, bodyErr))
		}
		return apiErr
	}
	if bodyErr != nil {
		return fmt.Errorf("%s: %w", operation, bodyErr)
	}
	if into == nil {
		return nil
	}
	if err := json.Unmarshal(body, into); err != nil {
		return fmt.Errorf("%s: decode response: %w", operation, err)
	}
	return nil
}
