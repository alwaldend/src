package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
)

// configuration identifies one runner and its persistent credential cache.
type configuration struct {
	// instance is the canonical HTTPS Forgejo origin.
	instance string
	// repository identifies the owner/name scope.
	repository string
	// name is the runner's stable name within that repository.
	name string
	// mount is the Vault KV v2 mount.
	mount string
	// path is the credential path below the mount.
	path string
}

// credentials contains the service identity and confidential one-time response.
type credentials struct {
	// Instance binds the cached credential to its Forgejo origin.
	Instance string `json:"instance"`
	// Repository binds the credential to its repository scope.
	Repository string `json:"repository"`
	// Name records the configured runner name.
	Name string `json:"name"`
	// ID is the Forgejo runner identifier used for reconciliation.
	ID int64 `json:"id"`
	// UUID identifies the runner connection.
	UUID string `json:"uuid"`
	// Token authenticates the runner and never leaves Vault or process memory.
	Token string `json:"token"`
}

// redactedError retains an underlying cause without disclosing upstream text.
type redactedError struct {
	// cause may contain response bodies, URLs or credentials.
	cause error
}

// Error supplies the only upstream error text exposed by the command.
func (e redactedError) Error() string { return "request failed (upstream details withheld)" }

// Unwrap preserves cancellation and transport error inspection.
func (e redactedError) Unwrap() error { return e.cause }

// pathPattern allows only explicit non-secret Vault path segments.
var pathPattern = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9_.-]*(/[A-Za-z0-9_-][A-Za-z0-9_.-]*)*$`)

// validate rejects ambiguous endpoints and missing invocation credentials.
func (c *configuration) validate() error {
	origin, err := url.Parse(c.instance)
	if err != nil || origin.Scheme != "https" || origin.Host == "" || origin.User != nil || (origin.Path != "" && origin.Path != "/") || origin.RawQuery != "" || origin.Fragment != "" {
		return errors.New("FORGEJO_HOST must be an HTTPS origin")
	}
	c.instance = strings.TrimSuffix(c.instance, "/")
	if !pathPattern.MatchString(c.repository) || strings.Count(c.repository, "/") != 1 || !pathPattern.MatchString(c.mount) || strings.Contains(c.mount, "/") || !pathPattern.MatchString(c.path) || strings.TrimSpace(c.name) == "" {
		return errors.New("set a repository owner/name, runner name, KV mount and credential path")
	}
	for _, variable := range []string{"FORGEJO_API_TOKEN", "VAULT_TOKEN"} {
		value := os.Getenv(variable)
		if value == "" || strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("missing or invalid %s", variable)
		}
	}
	address, err := url.Parse(os.Getenv("VAULT_ADDR"))
	if err != nil || address.Scheme != "https" || address.Host == "" || address.User != nil || address.RawQuery != "" || address.Fragment != "" {
		return errors.New("VAULT_ADDR must be an HTTPS URL")
	}
	if value := os.Getenv("VAULT_SKIP_VERIFY"); value != "" && value != "false" && value != "0" {
		return errors.New("Vault certificate verification must be enabled")
	}
	return nil
}

// ensure reuses a matching Vault record or performs one guarded registration.
// The owning Terraform backend lock serializes invocations; no concurrent
// registration outside that workflow is supported.
func (c configuration) ensure(ctx context.Context) (bool, error) {
	if err := c.validate(); err != nil {
		return false, fmt.Errorf("validate bootstrap configuration: %w", err)
	}
	vaultConfig := vault.DefaultConfig()
	vaultConfig.DisableRedirects = true
	// Vault's logical reader can treat a returned 3xx response as an absent
	// secret. Reject redirects as transport errors before that conversion.
	vaultConfig.HttpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("Vault redirect rejected")
	}
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return false, fmt.Errorf("configure Vault client: %w", redactedError{err})
	}
	client.SetMaxRetries(0)
	kv := client.KVv2(c.mount)
	secret, err := kv.Get(ctx, c.path)
	if err == nil {
		raw, err := json.Marshal(secret.Data)
		if err != nil {
			return false, fmt.Errorf("read cached credentials: %w", redactedError{err})
		}
		var stored credentials
		if err := json.Unmarshal(raw, &stored); err != nil {
			return false, fmt.Errorf("decode cached credentials: %w", redactedError{err})
		}
		if stored.Instance != c.instance || stored.Repository != c.repository || stored.Name != c.name || stored.ID <= 0 || stored.UUID == "" || stored.Token == "" {
			return false, errors.New("cached credentials are incomplete or belong to another runner; reconcile the Vault record")
		}
		return true, nil
	}
	if !errors.Is(err, vault.ErrSecretNotFound) {
		return false, fmt.Errorf("read Vault credential cache: %w", redactedError{err})
	}
	// A soft-deleted secret still has history; do not replace its cached token.
	if secret != nil {
		return false, errors.New("credential history exists without readable data; reconcile the Vault record")
	}
	httpClient := &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	endpoint := c.instance + "/api/v1/repos/" + c.repository + "/actions/runners"
	for page := 1; page <= 1000; page++ {
		var runners []credentials
		header, err := request(ctx, httpClient, http.MethodGet, endpoint+"?limit=50&page="+strconv.Itoa(page), nil, &runners, http.StatusOK)
		if err != nil {
			return false, fmt.Errorf("check existing runners: %w", err)
		}
		for _, runner := range runners {
			if runner.Name == c.name {
				return false, errors.New("a same-name runner exists without cached credentials; reconcile it before retrying")
			}
		}
		if !strings.Contains(header.Get("Link"), `rel="next"`) {
			break
		}
		if page == 1000 {
			return false, errors.New("runner pagination exceeded the safety bound")
		}
	}
	stored := credentials{Instance: c.instance, Repository: c.repository, Name: c.name}
	if _, err := request(ctx, httpClient, http.MethodPost, endpoint, map[string]string{"name": c.name}, &stored, http.StatusCreated); err != nil {
		return false, fmt.Errorf("create runner; inspect same-name registrations before retrying: %w", err)
	}
	if stored.ID <= 0 || stored.UUID == "" || stored.Token == "" {
		return false, errors.New("creation returned incomplete credentials; reconcile the remote runner before retrying")
	}
	values := map[string]any{"instance": c.instance, "repository": c.repository, "name": c.name, "id": stored.ID, "uuid": stored.UUID, "token": stored.Token}
	if _, err := kv.Put(ctx, c.path, values, vault.WithCheckAndSet(0)); err != nil {
		return false, fmt.Errorf("runner created but Vault save failed; reconcile before retrying: %w", redactedError{err})
	}
	return false, nil
}

// request performs one bounded Forgejo operation without logging secret material.
func request(ctx context.Context, client *http.Client, method, endpoint string, body any, output any, expected int) (http.Header, error) {
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("encode API request: %w", redactedError{err})
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("prepare API request: %w", redactedError{err})
	}
	req.Header.Set("Authorization", "token "+os.Getenv("FORGEJO_API_TOKEN"))
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send API request: %w", redactedError{err})
	}
	raw, readErr := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	closeErr := response.Body.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return nil, fmt.Errorf("read API response: %w", redactedError{err})
	}
	if response.StatusCode != expected {
		return nil, fmt.Errorf("unexpected API status %d", response.StatusCode)
	}
	if len(raw) > 1<<20 {
		return nil, errors.New("API response exceeds one MiB")
	}
	if err := json.Unmarshal(raw, output); err != nil {
		return nil, fmt.Errorf("decode API response: %w", redactedError{err})
	}
	return response.Header, nil
}
