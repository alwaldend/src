package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
)

// main adapts command errors to a nonzero process exit without printing secrets.
func main() {
	if err := command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// command configures the explicit schema update operation.
func command() *cobra.Command {
	var endpoint, output string
	var maxBytes int64
	cmd := &cobra.Command{
		Use: "fetch_schema", Short: "Fetch the authenticated X OpenAPI schema for mirror publication",
		SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !filepath.IsAbs(output) {
				if workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY"); workspace != "" {
					output = filepath.Join(workspace, output)
				}
			}
			if err := fetch(cmd.Context(), endpoint, output, maxBytes); err != nil {
				return fmt.Errorf("update X schema: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Updated X OpenAPI schema")
			return nil
		},
		Args: cobra.NoArgs,
	}
	cmd.Flags().StringVar(&endpoint, "url", "https://api.x.com/2/openapi.json", "Schema URL; HTTP is permitted only for local fixtures")
	cmd.Flags().StringVar(&output, "out", "out/x_sdk/openapi.json", "Destination relative to the Bazel workspace")
	cmd.Flags().Int64Var(&maxBytes, "max-bytes", 32<<20, "Maximum response size")
	return cmd
}

// fetch downloads and validates one response before replacing the source file.
func fetch(ctx context.Context, endpoint, output string, maxBytes int64) error {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid schema URL")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Host == "" ||
		(parsed.Scheme != "https" && !(parsed.Scheme == "http" && (parsed.Hostname() == "127.0.0.1" || parsed.Hostname() == "::1"))) {
		return fmt.Errorf("schema URL must use HTTPS without credentials, query or fragment; local HTTP fixtures are permitted")
	}
	if maxBytes <= 0 || maxBytes > 256<<20 {
		return fmt.Errorf("max-bytes must be between 1 and 268435456")
	}
	credentials := oauth1.Credentials{
		APIKey: os.Getenv("X_API_KEY"), APISecret: os.Getenv("X_API_SECRET"),
		AccessToken: os.Getenv("X_ACCESS_TOKEN"), AccessTokenSecret: os.Getenv("X_ACCESS_TOKEN_SECRET"),
	}
	if err := credentials.Validate(); err != nil {
		return fmt.Errorf("load injected X credentials: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("create schema request: %w", err)
	}
	authorization, err := credentials.Sign(http.MethodGet, endpoint, nil, nil, time.Now())
	if err != nil {
		return fmt.Errorf("sign schema request: %w", err)
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download schema: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("schema endpoint returned HTTP %d; response body omitted", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		return fmt.Errorf("read schema response: %w", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("close schema response: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return fmt.Errorf("schema exceeds %d bytes", maxBytes)
	}
	// schema retains only the fields needed to distinguish OpenAPI from API errors.
	var schema struct {
		// OpenAPI identifies the document version.
		OpenAPI string `json:"openapi"`
		// Info contains required API metadata.
		Info map[string]json.RawMessage `json:"info"`
		// Paths contains the API operations.
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(body, &schema); err != nil {
		return fmt.Errorf("schema response is not an OpenAPI JSON document")
	}
	if !strings.HasPrefix(schema.OpenAPI, "3.") || len(schema.Info) == 0 || len(schema.Paths) == 0 {
		return fmt.Errorf("schema response lacks OpenAPI 3 metadata or API paths")
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return fmt.Errorf("create schema directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(output), ".openapi-*.json")
	if err != nil {
		return fmt.Errorf("create temporary schema: %w", err)
	}
	defer os.Remove(temporary.Name())
	defer temporary.Close()
	if _, err := temporary.Write(append(bytes.TrimSpace(body), '\n')); err != nil {
		return fmt.Errorf("write temporary schema: %w", err)
	}
	if err := temporary.Chmod(0o644); err != nil {
		return fmt.Errorf("set schema permissions: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary schema: %w", err)
	}
	if err := os.Rename(temporary.Name(), output); err != nil {
		return fmt.Errorf("replace schema source: %w", err)
	}
	return nil
}
