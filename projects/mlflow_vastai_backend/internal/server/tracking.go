package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
)

// tracking reads and updates lifecycle status at the configured MLflow endpoint.
type tracking struct {
	// endpoint is the sole destination permitted for tracking credentials.
	endpoint string
	// token authenticates tracking requests when configured.
	token string
	// client bounds HTTP calls and refuses redirects.
	client *http.Client
}

// newTracking scopes MLflow requests to the coordinator's configured endpoint.
func newTracking(endpoint *url.URL, token string) tracking {
	return tracking{
		endpoint: strings.TrimRight(endpoint.String(), "/"), token: token,
		client: &http.Client{Timeout: 8 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}},
	}
}

// publish updates a confirmed terminal status when the MLflow run remains active.
func (tracking tracking) publish(ctx context.Context, job *jobpb.Job) error {
	var result struct {
		// Run contains the tracking run's status.
		Run struct {
			// Info contains the run lifecycle.
			Info struct {
				// Status records whether the run is still active.
				Status string `json:"status"`
			} `json:"info"`
		} `json:"run"`
	}
	query := url.Values{"run_id": {job.GetJobId()}}
	if err := tracking.request(ctx, http.MethodGet, "get?"+query.Encode(), nil, &result); err != nil {
		return fmt.Errorf("read tracking run: %w", err)
	}
	switch result.Run.Info.Status {
	case "RUNNING", "SCHEDULED":
	case "FINISHED", "FAILED", "KILLED":
		return nil
	default:
		return errors.New("tracking run response has no recognized status")
	}
	payload, err := json.Marshal(struct {
		// ID names the existing MLflow run.
		ID string `json:"run_id"`
		// Status is the confirmed terminal worker state.
		Status string `json:"status"`
		// EndTime gives MLflow the observation time in Unix milliseconds.
		EndTime int64 `json:"end_time"`
	}{job.GetJobId(), strings.TrimPrefix(job.GetState().String(), "JOB_STATE_"), time.Now().UnixMilli()})
	if err != nil {
		return fmt.Errorf("encode tracking update: %w", err)
	}
	if err := tracking.request(ctx, http.MethodPost, "update", payload, nil); err != nil {
		return fmt.Errorf("update tracking run: %w", err)
	}
	return nil
}

// request calls a bounded run endpoint without logging bodies or credentials.
func (tracking tracking) request(ctx context.Context, method, route string, payload []byte, result any) (err error) {
	request, err := http.NewRequestWithContext(ctx, method, tracking.endpoint+"/api/2.0/mlflow/runs/"+route, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create tracking request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if tracking.token != "" {
		request.Header.Set("Authorization", "Bearer "+tracking.token)
	}
	response, err := tracking.client.Do(request)
	if err != nil {
		return fmt.Errorf("send tracking request: %w", err)
	}
	defer func() { err = errors.Join(err, response.Body.Close()) }()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("tracking endpoint returned HTTP %d", response.StatusCode)
	}
	if result != nil {
		if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(result); err != nil {
			return fmt.Errorf("decode tracking response: %w", err)
		}
	}
	return nil
}
