package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// main dispatches staging and publication.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	directory := filepath.Join(workspace, "out", "artifact-mirror")
	var helm, oras, owner string
	var inventoryFile string
	var artifacts inventory
	root := &cobra.Command{
		Use:          "artifact_mirror",
		SilenceUsage: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var err error
			artifacts, err = loadInventory(inventoryFile)
			if err != nil {
				return fmt.Errorf("load upstream pins: %w", err)
			}
			return nil
		},
	}
	root.PersistentFlags().StringVar(&inventoryFile, "inventory", "", "Artifact pin file")
	root.PersistentFlags().StringVar(&directory, "directory", directory, "Local mirror directory")
	root.PersistentFlags().StringVar(&helm, "helm", "helm", "Path to Helm")
	root.PersistentFlags().StringVar(&oras, "oras", "oras", "Path to ORAS")
	publish := &cobra.Command{
		Use:  "push",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return push(cmd.Context(), directory, oras, owner, artifacts)
		},
	}
	publish.Flags().StringVar(&owner, "owner", "", "Forgejo package owner")
	root.AddCommand(
		&cobra.Command{
			Use:  "pull",
			Args: cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				return pull(cmd.Context(), directory, helm, oras, artifacts)
			},
		},
		publish,
	)
	if err := root.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

// command runs a packaged tool with console output and optional standard input.
func command(ctx context.Context, tool string, input io.Reader, args ...string) error {
	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = input, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %q %s: %w", filepath.Base(tool), args[0], err)
	}
	return nil
}

// verifyDigest confirms that a stream has the pinned content.
func verifyDigest(data io.Reader, expected, name string) error {
	hash := sha256.New()
	if _, err := io.Copy(hash, data); err != nil {
		return fmt.Errorf("hash %q: %w", name, err)
	}
	if actual := fmt.Sprintf("sha256:%x", hash.Sum(nil)); actual != expected {
		return fmt.Errorf("digest mismatch for %q: got %s, want %s", name, actual, expected)
	}
	return nil
}

// verifyChart checks an archive before it can be used for publication.
func verifyChart(filename, digest string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open chart %q: %w", filename, err)
	}
	defer file.Close()
	if err := verifyDigest(file, digest, filename); err != nil {
		return fmt.Errorf("verify chart archive: %w", err)
	}
	return nil
}

// layoutReference names an image tag in its local OCI layout.
func layoutReference(directory string, artifact image) string {
	return filepath.Join(directory, "images", path.Base(artifact.Repository)) + ":" + artifact.Tag
}

// verifyImage compares the local or remote OCI digest with the pinned index.
func verifyImage(ctx context.Context, oras, reference, digest string, args ...string) error {
	cmd := exec.CommandContext(ctx, oras, append(append([]string{"resolve"}, args...), reference)...)
	cmd.Stderr = os.Stderr
	result, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("resolve OCI reference %q: %w", reference, err)
	}
	if actual := strings.TrimSpace(string(result)); actual != digest {
		return fmt.Errorf("digest mismatch for %q: got %s, want %s", reference, actual, digest)
	}
	return nil
}

// pull downloads charts and stages complete pinned image indexes.
func pull(ctx context.Context, directory, helm, oras string, artifacts inventory) error {
	chartDirectory := filepath.Join(directory, "charts")
	if err := os.MkdirAll(chartDirectory, 0o755); err != nil {
		return fmt.Errorf("create chart directory %q: %w", chartDirectory, err)
	}
	staging, err := os.MkdirTemp(chartDirectory, ".download-")
	if err != nil {
		return fmt.Errorf("create chart download directory: %w", err)
	}
	defer os.RemoveAll(staging)
	for _, artifact := range artifacts.Charts {
		args := []string{"pull", artifact.Source, "--destination", staging}
		if strings.HasPrefix(artifact.Source, "oci://") {
			args = append(args, "--version", artifact.Version)
		}
		if err := command(ctx, helm, nil, args...); err != nil {
			return fmt.Errorf("pull chart %q: %w", artifact.filename(), err)
		}
		download := filepath.Join(staging, artifact.filename())
		if err := verifyChart(download, artifact.Digest); err != nil {
			return fmt.Errorf("verify pulled chart: %w", err)
		}
		if err := os.Rename(download, filepath.Join(chartDirectory, artifact.filename())); err != nil {
			return fmt.Errorf("store verified chart %q: %w", artifact.filename(), err)
		}
	}
	for _, artifact := range artifacts.Images {
		source := artifact.Repository + "@" + artifact.Digest
		if err := command(ctx, oras, nil, "cp", "--to-oci-layout", source, layoutReference(directory, artifact)); err != nil {
			return fmt.Errorf("pull image %q: %w", source, err)
		}
	}
	return nil
}

// forgejoUsername resolves the identity associated with the AL login token.
func forgejoUsername(ctx context.Context, client *http.Client, baseURL, token string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/v1/user", nil)
	if err != nil {
		return "", fmt.Errorf("construct Forgejo identity request: %w", err)
	}
	request.Header.Set("Authorization", "token "+token)
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("read Forgejo identity: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("read Forgejo identity: HTTP %d", response.StatusCode)
	}
	var account struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(response.Body).Decode(&account); err != nil {
		return "", fmt.Errorf("decode Forgejo identity: %w", err)
	}
	if account.Login == "" {
		return "", fmt.Errorf("Forgejo identity has no username")
	}
	return account.Login, nil
}

// publishChart uploads a Helm archive and verifies the stored bytes on success or retry.
func publishChart(ctx context.Context, client *http.Client, registryURL, username, token, filename string, artifact chart) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open chart %q for upload: %w", filename, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat chart %q: %w", filename, err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, registryURL+"/api/charts", file)
	if err != nil {
		return fmt.Errorf("construct chart upload: %w", err)
	}
	request.ContentLength = info.Size()
	request.SetBasicAuth(username, token)
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("upload chart %q: %w", filename, err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusConflict {
		return fmt.Errorf("upload chart %q: HTTP %d", filename, response.StatusCode)
	}
	request, err = http.NewRequestWithContext(ctx, http.MethodGet, registryURL+"/"+artifact.filename(), nil)
	if err != nil {
		return fmt.Errorf("construct chart verification request: %w", err)
	}
	request.SetBasicAuth(username, token)
	response, err = client.Do(request)
	if err != nil {
		return fmt.Errorf("read published chart %q: %w", filename, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("read published chart %q: HTTP %d", filename, response.StatusCode)
	}
	if err := verifyDigest(response.Body, artifact.Digest, artifact.filename()); err != nil {
		return fmt.Errorf("verify published chart: %w", err)
	}
	fmt.Println("Verified published chart", artifact.filename())
	return nil
}

// push verifies local inputs before writing packages through the Forgejo registry.
func push(ctx context.Context, directory, oras, owner string, artifacts inventory) error {
	if owner == "" || owner == "." || owner == ".." || strings.ContainsAny(owner, "/\\?#") {
		return fmt.Errorf("provide a valid Forgejo package owner with --owner")
	}
	token := os.Getenv("FORGEJO_API_TOKEN")
	if token == "" {
		return fmt.Errorf("Forgejo login supplied no API token")
	}
	baseURL := strings.TrimRight(os.Getenv("FORGEJO_HOST"), "/")
	registry, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("parse Forgejo registry URL: %w", err)
	}
	if registry.Scheme != "https" || registry.Host == "" || registry.User != nil ||
		registry.Path != "" || registry.RawQuery != "" || registry.Fragment != "" {
		return fmt.Errorf("Forgejo login must supply an HTTPS registry origin")
	}
	for _, artifact := range artifacts.Charts {
		if err := verifyChart(filepath.Join(directory, "charts", artifact.filename()), artifact.Digest); err != nil {
			return fmt.Errorf("check staged chart: %w", err)
		}
	}
	for _, artifact := range artifacts.Images {
		if err := verifyImage(ctx, oras, layoutReference(directory, artifact), artifact.Digest, "--oci-layout"); err != nil {
			return fmt.Errorf("check staged image: %w", err)
		}
	}
	client := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	username, err := forgejoUsername(ctx, client, baseURL, token)
	if err != nil {
		return fmt.Errorf("authenticate artifact publication: %w", err)
	}
	credentials, err := os.MkdirTemp(directory, ".registry-")
	if err != nil {
		return fmt.Errorf("create private registry configuration: %w", err)
	}
	defer os.RemoveAll(credentials)
	registryConfig := filepath.Join(credentials, "config.json")
	if err := command(ctx, oras, strings.NewReader(token+"\n"), "login", "--registry-config", registryConfig, "--username", username, "--password-stdin", registry.Host); err != nil {
		return fmt.Errorf("authenticate OCI publication: %w", err)
	}
	registryURL := baseURL + "/api/packages/" + url.PathEscape(owner) + "/helm"
	for _, artifact := range artifacts.Charts {
		if err := publishChart(ctx, client, registryURL, username, token, filepath.Join(directory, "charts", artifact.filename()), artifact); err != nil {
			return fmt.Errorf("publish charts: %w", err)
		}
	}
	for _, artifact := range artifacts.Images {
		destination := registry.Host + "/" + owner + "/" + path.Base(artifact.Repository) + ":" + artifact.Tag
		if err := command(ctx, oras, nil, "cp", "--from-oci-layout", "--to-registry-config", registryConfig, layoutReference(directory, artifact), destination); err != nil {
			return fmt.Errorf("publish image %q: %w", destination, err)
		}
		if err := verifyImage(ctx, oras, destination, artifact.Digest, "--registry-config", registryConfig); err != nil {
			return fmt.Errorf("verify published image: %w", err)
		}
	}
	return nil
}
