package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/client"
)

// TestGeneratedClient exercises a schema-defined operation and its error response.
func TestGeneratedClient(t *testing.T) {
	for _, status := range []int{200, 401} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/2/users/123" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				if r.Header.Get("Authorization") != "Bearer fixture-token" {
					t.Error("missing bearer authentication")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				if status == http.StatusOK {
					_, _ = w.Write([]byte(`{"data":{"id":"123","name":"Fixture User","username":"fixture"}}`))
				} else {
					_, _ = w.Write([]byte(`{"title":"Unauthorized","detail":"Fixture rejection","status":401}`))
				}
			}))
			defer server.Close()
			configuration := client.NewConfiguration()
			configuration.Servers = client.ServerConfigurations{{URL: server.URL}}
			configuration.HTTPClient = server.Client()
			api := client.NewAPIClient(configuration)
			ctx := context.WithValue(context.Background(), client.ContextAccessToken, "fixture-token")
			result, response, err := api.UsersAPI.GetUsersById(ctx, "123").Execute()
			if response == nil || response.StatusCode != status {
				t.Fatalf("unexpected HTTP response: %v, error=%v", response, err)
			}
			if status != http.StatusOK {
				if err == nil {
					t.Fatal("rejected request succeeded")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.Data == nil || result.Data.GetId() != "123" || result.Data.GetUsername() != "fixture" {
				t.Fatalf("unexpected decoded user: %+v", result)
			}
			artifact, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "user.json"), artifact, 0o644); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestGeneratedPartialResponse retains successful users alongside typed problems.
func TestGeneratedPartialResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/2/users" || r.URL.Query().Get("ids") != "123,456" {
			t.Errorf("unexpected user lookup: %s %s", r.Method, r.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"id":"123","name":"Fixture User","username":"fixture"}],"errors":[{"type":"https://api.x.com/2/problems/resource-not-found","title":"Not Found Error","detail":"User 456 does not exist","resource_type":"user","resource_id":"456"}]}`))
	}))
	defer server.Close()
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{{URL: server.URL}}
	configuration.HTTPClient = server.Client()
	api := client.NewAPIClient(configuration)
	result, response, err := api.UsersAPI.GetUsersByIds(context.Background()).Ids([]string{"123", "456"}).Execute()
	if err != nil {
		t.Fatalf("decode partial-success response: %v", err)
	}
	if response.StatusCode != http.StatusOK || len(result.Data) != 1 || result.Data[0].GetId() != "123" {
		t.Fatalf("successful user was lost: %+v", result)
	}
	if len(result.Errors) != 1 || result.Errors[0].ResourceNotFoundProblem == nil {
		t.Fatalf("problem did not decode to its declared type: %+v", result.Errors)
	}
	if problem := result.Errors[0].ResourceNotFoundProblem; problem.GetResourceId() != "456" || problem.GetResourceType() != "user" {
		t.Fatalf("problem lost its resource identity: %+v", problem)
	}
	artifact, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "partial-users.json"), artifact, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestGeneratedMediaUpload verifies base64 media JSON through a generated API call.
func TestGeneratedMediaUpload(t *testing.T) {
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/2/media/upload/456/append" {
			t.Errorf("unexpected media request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()
	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{{URL: server.URL}}
	configuration.HTTPClient = server.Client()
	api := client.NewAPIClient(configuration)
	request := client.NewAppendMediaUploadRequest(json.RawMessage(`"aGVsbG8="`), 0)
	_, response, err := api.MediaAPI.AppendMediaUpload(context.Background(), "456").AppendMediaUploadRequest(*request).Execute()
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || received["media"] != "aGVsbG8=" || received["segment_index"] != float64(0) {
		t.Fatalf("unexpected media response or request: status=%d, request=%v", response.StatusCode, received)
	}
	artifact, err := json.MarshalIndent(received, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"), "media.json"), artifact, 0o644); err != nil {
		t.Fatal(err)
	}
}
