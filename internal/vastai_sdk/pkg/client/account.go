package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	generatedapi "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/internal/generated"
)

// CreateSSHKey registers a public key and returns the account-local numeric ID.
func (c *Client) CreateSSHKey(ctx context.Context, key string) (int64, error) {
	var result struct {
		// Success confirms key creation.
		Success bool `json:"success"`
		// Key identifies the registered key.
		Key struct {
			// ID is the server-assigned key ID.
			ID int64 `json:"id"`
		} `json:"key"`
	}
	request := c.api.AccountsAPI.CreateSshKey(ctx).CreateSshKeyRequest(generatedapi.CreateSshKeyRequest{SshKey: key})
	if err := generated(c, http.MethodPost, "/ssh/", &result, true, request.Execute); err != nil {
		return 0, fmt.Errorf("create SSH key: %w", err)
	}
	if !result.Success || result.Key.ID <= 0 {
		return 0, errors.New("SSH key response omitted creation confirmation or ID")
	}
	return result.Key.ID, nil
}

// SSHKey retrieves a key by ID from the account's key collection.
func (c *Client) SSHKey(ctx context.Context, id int64) (string, error) {
	var keys []struct {
		// ID identifies the registered key.
		ID int64 `json:"id"`
		// Key is the public key in list responses.
		Key string `json:"key"`
		// PublicKey is the public-key alias used by some API versions.
		PublicKey string `json:"public_key"`
		// DeletedAt marks soft-deleted keys.
		DeletedAt *string `json:"deleted_at"`
	}
	if err := c.request(ctx, http.MethodGet, "/ssh/", nil, &keys); err != nil {
		return "", fmt.Errorf("list SSH keys: %w", err)
	}
	if keys == nil {
		return "", errors.New("SSH key response is not a key collection")
	}
	for _, key := range keys {
		if key.ID <= 0 {
			return "", errors.New("SSH key collection contains an invalid key ID")
		}
		if key.ID == id && key.DeletedAt == nil {
			value := key.Key
			if value == "" {
				value = key.PublicKey
			}
			if value == "" {
				return "", errors.New("SSH key response omitted public key")
			}
			return value, nil
		}
	}
	return "", fmt.Errorf("SSH key %d: %w", id, ErrNotFound)
}

// UpdateSSHKey replaces the public key associated with one key ID.
func (c *Client) UpdateSSHKey(ctx context.Context, id int64, key string) error {
	var result struct{}
	request := c.api.AccountsAPI.UpdateSshKey(ctx, id).UpdateSshKeyRequest(generatedapi.UpdateSshKeyRequest{SshKey: key})
	if err := generated(c, http.MethodPut, "/ssh/"+strconv.FormatInt(id, 10)+"/", &result, true, request.Execute); err != nil {
		return fmt.Errorf("update SSH key %d: %w", id, err)
	}
	return nil
}

// DeleteSSHKey deletes one registered public key.
func (c *Client) DeleteSSHKey(ctx context.Context, id int64) error {
	var result struct{}
	request := c.api.AccountsAPI.DeleteSshKey(ctx, id)
	if err := generated(c, http.MethodDelete, "/ssh/"+strconv.FormatInt(id, 10)+"/", &result, true, request.Execute); err != nil {
		return fmt.Errorf("delete SSH key %d: %w", id, err)
	}
	return nil
}

// EnvironmentVariable reads one variable, returning ErrValueRedacted for a mask.
func (c *Client) EnvironmentVariable(ctx context.Context, name string) (string, error) {
	var result struct {
		// Secrets holds the account's variable mapping.
		Secrets map[string]string `json:"secrets"`
	}
	request := c.api.AccountsAPI.ShowEnvVars(ctx)
	if err := generated(c, http.MethodGet, "/secrets/", &result, true, request.Execute); err != nil {
		return "", fmt.Errorf("read environment variable: %w", err)
	}
	if result.Secrets == nil {
		return "", errors.New("environment response omitted variable collection")
	}
	value, ok := result.Secrets[name]
	if !ok {
		return "", fmt.Errorf("environment variable absent: %w", ErrNotFound)
	}
	if value == "********" {
		return "", fmt.Errorf("read environment variable: %w", ErrValueRedacted)
	}
	return value, nil
}

// WriteEnvironmentVariable creates or updates a named account variable.
func (c *Client) WriteEnvironmentVariable(ctx context.Context, name, value string, create bool) error {
	method := http.MethodPut
	if create {
		method = http.MethodPost
	}
	var result struct {
		// Success confirms the write.
		Success bool `json:"success"`
	}
	var err error
	if create {
		request := c.api.AccountsAPI.CreateEnvVar(ctx).CreateEnvVarRequest(generatedapi.CreateEnvVarRequest{Key: name, Value: value})
		err = generated(c, method, "/secrets/", &result, true, request.Execute)
	} else {
		request := c.api.AccountsAPI.UpdateEnvVar(ctx).UpdateEnvVarRequest(generatedapi.UpdateEnvVarRequest{Key: name, Value: value})
		err = generated(c, method, "/secrets/", &result, true, request.Execute)
	}
	if err != nil {
		return fmt.Errorf("write environment variable: %w", err)
	}
	if !result.Success {
		return errors.New("environment variable response omitted success confirmation")
	}
	return nil
}

// DeleteEnvironmentVariable deletes one named account variable.
func (c *Client) DeleteEnvironmentVariable(ctx context.Context, name string) error {
	var result struct{}
	request := c.api.AccountsAPI.DeleteUserSecret(ctx).DeleteUserSecretRequest(generatedapi.DeleteUserSecretRequest{Key: name})
	if err := generated(c, http.MethodDelete, "/secrets/", &result, true, request.Execute); err != nil {
		return fmt.Errorf("delete environment variable: %w", err)
	}
	return nil
}

// Role contains the stable numeric identifier and JSON permission document.
type Role struct {
	// ID is the server-assigned role ID used by update operations.
	ID int64 `json:"id"`
	// Name identifies the role for reads and deletion.
	Name string `json:"name"`
	// Permissions is the role's permission document.
	Permissions json.RawMessage `json:"permissions"`
}

// TeamRole reads a role by name, treating only HTTP 404 as absence.
func (c *Client) TeamRole(ctx context.Context, name string) (*Role, error) {
	var role Role
	request := c.api.TeamAPI.ShowTeamRole(ctx, name)
	if err := generated(c, http.MethodGet, "/team/roles/"+url.PathEscape(name)+"/", &role, true, request.Execute); err != nil {
		return nil, fmt.Errorf("read team role %q: %w", name, err)
	}
	if role.ID <= 0 || role.Name != name || len(role.Permissions) == 0 || string(role.Permissions) == "null" {
		return nil, errors.New("team role response omitted or mismatched role fields")
	}
	return &role, nil
}

// WriteTeamRole creates a role by name or updates an existing numeric ID.
func (c *Client) WriteTeamRole(ctx context.Context, id int64, name, permissions string) error {
	method, path := http.MethodPost, "/team/roles/"
	if id > 0 {
		method, path = http.MethodPut, "/team/roles/"+strconv.FormatInt(id, 10)+"/"
	}
	var result json.RawMessage
	var err error
	if id > 0 {
		body := generatedapi.UpdateTeamRoleRequest{Name: name}
		if err := json.Unmarshal([]byte(permissions), &body.Permissions); err != nil {
			return fmt.Errorf("decode role permissions: %w", err)
		}
		request := c.api.TeamAPI.UpdateTeamRole(ctx, id).UpdateTeamRoleRequest(body)
		err = generated(c, method, path, &result, false, request.Execute)
	} else {
		body := generatedapi.CreateTeamRoleRequest{Name: name}
		if err := json.Unmarshal([]byte(permissions), &body.Permissions); err != nil {
			return fmt.Errorf("decode role permissions: %w", err)
		}
		request := c.api.TeamAPI.CreateTeamRole(ctx).CreateTeamRoleRequest(body)
		err = generated(c, method, path, &result, false, request.Execute)
	}
	if err != nil {
		return fmt.Errorf("write team role %q: %w", name, err)
	}
	if err := roleConfirmation(result, "Successfully Updated Team Role For "+name); err != nil {
		return fmt.Errorf("confirm team role %q: %w", name, err)
	}
	return nil
}

// DeleteTeamRole removes a role by name.
func (c *Client) DeleteTeamRole(ctx context.Context, name string) error {
	var result json.RawMessage
	request := c.api.TeamAPI.RemoveTeamRole(ctx, name)
	if err := generated(c, http.MethodDelete, "/team/roles/"+url.PathEscape(name)+"/", &result, false, request.Execute); err != nil {
		return fmt.Errorf("delete team role %q: %w", name, err)
	}
	if err := roleConfirmation(result, ""); err != nil {
		return fmt.Errorf("confirm role deletion: %w", err)
	}
	return nil
}

// roleConfirmation accepts the documented success object and update string forms.
func roleConfirmation(data json.RawMessage, expected string) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil && expected != "" && text == expected {
		return nil
	}
	var result struct {
		// Success is the boolean success variant.
		Success bool `json:"success"`
		// Status is the status-string confirmation used by deletion.
		Status string `json:"status"`
		// Message is the create confirmation.
		Message string `json:"msg"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("decode role confirmation: %w", err)
	}
	if result.Success || result.Status == "success" || result.Message == "success" {
		return nil
	}
	return errors.New("role response omitted success confirmation")
}
