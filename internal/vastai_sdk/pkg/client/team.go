package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	generatedapi "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/internal/generated"
)

// apiKeyMetadata identifies keys and team contexts without retrieving credentials.
type apiKeyMetadata struct {
	// ID identifies an API key or the owner context used by save-context.
	ID int64 `json:"id"`
	// TeamID identifies the account represented by the key.
	TeamID int64 `json:"team_id"`
	// TeamName identifies the team's current name.
	TeamName string `json:"team_name"`
	// Name identifies the key or owner role.
	Name string `json:"name"`
}

// keyMetadata obtains the authenticated user's key collection, including owner contexts.
func (c *Client) keyMetadata(ctx context.Context) ([]apiKeyMetadata, error) {
	var result struct {
		// Keys is the validated key collection from the API.
		Keys []apiKeyMetadata `json:"apikeys"`
	}
	request := c.api.AccountsAPI.ShowApiKeys(ctx)
	if err := generated(c, http.MethodGet, "/auth/apikeys/", &result, true, request.Execute); err != nil {
		return nil, fmt.Errorf("list API key metadata: %w", err)
	}
	if result.Keys == nil {
		return nil, errors.New("API key response omitted its collection")
	}
	for _, key := range result.Keys {
		if key.ID <= 0 {
			return nil, errors.New("API key collection contains an invalid key identifier")
		}
	}
	return result.Keys, nil
}

// CreateTeam creates a separate team without transferring credits.
func (c *Client) CreateTeam(ctx context.Context, name string) error {
	var result struct {
		// Success confirms team creation.
		Success bool `json:"success"`
	}
	zero := float64(0)
	request := c.api.TeamAPI.CreateTeam(ctx).CreateTeamRequest(generatedapi.CreateTeamRequest{TeamName: name, TransferCredit: &zero})
	if err := generated(c, http.MethodPost, "/team/", &result, true, request.Execute); err != nil {
		return fmt.Errorf("create team %q: %w", name, err)
	}
	if !result.Success {
		return errors.New("team creation omitted success confirmation")
	}
	return nil
}

// TeamContext returns an owned team's account and owner-context identifiers.
func (c *Client) TeamContext(ctx context.Context, name string) (int64, int64, error) {
	keys, err := c.keyMetadata(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("read team %q: %w", name, err)
	}
	for _, key := range keys {
		if key.TeamName == name && key.Name == "owner" && key.TeamID > 0 && key.ID > 0 {
			return key.TeamID, key.ID, nil
		}
	}
	return 0, 0, fmt.Errorf("team %q: %w", name, ErrNotFound)
}

// TeamContextByID locates only the account captured after confirmed creation.
func (c *Client) TeamContextByID(ctx context.Context, id int64) (string, int64, error) {
	keys, err := c.keyMetadata(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("read team account %d: %w", id, err)
	}
	for _, key := range keys {
		if key.TeamID == id && key.Name == "owner" && key.ID > 0 && key.TeamName != "" {
			return key.TeamName, key.ID, nil
		}
	}
	return "", 0, fmt.Errorf("team account %d: %w", id, ErrNotFound)
}

// WithKey constructs an independent client using generated credential material.
func (c *Client) WithKey(key string) (*Client, error) {
	client, err := New(Options{Endpoint: c.endpoint, APIKey: key})
	if err != nil {
		return nil, fmt.Errorf("construct generated-key client: %w", err)
	}
	return client, nil
}

// ActivateTeam changes this credential's context and verifies the selected account.
// Callers must use an independently generated key to preserve bootstrap context.
func (c *Client) ActivateTeam(ctx context.Context, contextID, teamID int64) error {
	selected, err := c.teamContextMatches(ctx, teamID)
	if err != nil {
		return err
	}
	if selected {
		return nil
	}
	var result struct{}
	if err := c.requestPolicy(ctx, http.MethodPut, "/users/save-context/", map[string]int64{"context_key_id": contextID}, &result, false); err != nil {
		return fmt.Errorf("activate generated team credential: %w", err)
	}
	selected, err = c.teamContextMatches(ctx, teamID)
	if err != nil {
		return err
	}
	if !selected {
		return errors.New("generated management key did not enter its requested team context")
	}
	return nil
}

// teamContextMatches reads the selected account before changing an owned key's context.
func (c *Client) teamContextMatches(ctx context.Context, teamID int64) (bool, error) {
	var current struct {
		// ID is the authenticated account after selection.
		ID int64 `json:"id"`
		// IsTeam confirms a team context rather than a personal account.
		IsTeam bool `json:"is_team"`
	}
	request := c.api.AccountsAPI.ShowUser(ctx)
	if err := generated(c, http.MethodGet, "/users/current/", &current, true, request.Execute); err != nil {
		return false, fmt.Errorf("read generated team context: %w", err)
	}
	if current.ID <= 0 {
		return false, errors.New("current account response omitted its identifier")
	}
	return current.ID == teamID && current.IsTeam, nil
}

// DeleteTeam destroys the team selected on this independent management client.
func (c *Client) DeleteTeam(ctx context.Context) error {
	var result struct{}
	request := c.api.TeamAPI.DestroyTeam(ctx)
	if err := generated(c, http.MethodDelete, "/team/", &result, true, request.Execute); err != nil {
		return fmt.Errorf("delete selected team: %w", err)
	}
	return nil
}

// CreateAPIKey generates a key in the selected team, returning its one-time value.
func (c *Client) CreateAPIKey(ctx context.Context, name, permissions string) (int64, string, error) {
	var result struct {
		// ID is the generated key ID.
		ID int64 `json:"id"`
		// Key is the generated credential; callers must mark it sensitive.
		Key string `json:"key"`
	}
	body := generatedapi.CreateApiKeyRequest{Name: name}
	if permissions != "" {
		if err := json.Unmarshal([]byte(permissions), &body.Permissions); err != nil {
			return 0, "", fmt.Errorf("decode API key permissions: %w", err)
		}
	}
	request := c.api.AccountsAPI.CreateApiKey(ctx).CreateApiKeyRequest(body)
	if err := generated(c, http.MethodPost, "/auth/apikeys/", &result, true, request.Execute); err != nil {
		return 0, "", fmt.Errorf("create API key: %w", err)
	}
	if result.ID <= 0 || result.Key == "" {
		return 0, "", errors.New("API key response omitted ID or generated credential; inspect account before retrying")
	}
	return result.ID, result.Key, nil
}

// ReadAPIKey checks identity without returning or replacing stored credential material.
func (c *Client) ReadAPIKey(ctx context.Context, id int64) error {
	keys, err := c.keyMetadata(ctx)
	if err != nil {
		return fmt.Errorf("read team API key: %w", err)
	}
	for _, key := range keys {
		if key.ID == id {
			return nil
		}
	}
	return fmt.Errorf("API key %d: %w", id, ErrNotFound)
}

// DeleteAPIKey revokes only the explicitly identified key in its owning team.
func (c *Client) DeleteAPIKey(ctx context.Context, id int64) error {
	if err := c.ReadAPIKey(ctx, id); err != nil {
		return err
	}
	var result json.RawMessage
	request := c.api.AccountsAPI.DeleteApiKey(ctx, id)
	if err := generated(c, http.MethodDelete, "/auth/apikeys/"+strconv.FormatInt(id, 10)+"/", &result, false, request.Execute); err != nil {
		return fmt.Errorf("revoke team API key: %w", err)
	}
	if err := c.ReadAPIKey(ctx, id); errors.Is(err, ErrNotFound) {
		return nil
	} else if err != nil {
		return fmt.Errorf("verify API key revocation: %w", err)
	}
	return errors.New("API key remains present after revocation")
}
