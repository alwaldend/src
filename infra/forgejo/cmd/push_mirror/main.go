// Command push_mirror bootstraps a native Forgejo HTTPS push mirror.
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
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type mirrorConfig struct {
	host, repository, target, forge, username, vaultPath, apiToken string
}

var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

func validRepository(repository string) bool {
	owner, name, found := strings.Cut(repository, "/")
	return found && repositoryPattern.MatchString(repository) && owner != "." && owner != ".." && name != "." && name != ".."
}

type pushMirror struct {
	Name     string `json:"remote_name"`
	Address  string `json:"remote_address"`
	Key      string `json:"public_key"`
	Interval string `json:"interval"`
	Sync     bool   `json:"sync_on_commit"`
	Filter   string `json:"branch_filter"`
}

func main() {
	if len(os.Args) != 2 || os.Args[1] != "ensure" {
		fmt.Fprintln(os.Stderr, "usage: push_mirror ensure")
		os.Exit(2)
	}
	config := mirrorConfig{
		host:       strings.TrimSuffix(os.Getenv("FORGEJO_HOST"), "/"),
		repository: os.Getenv("FORGEJO_MIRROR_REPOSITORY"),
		target:     os.Getenv("FORGEJO_MIRROR_TARGET"),
		forge:      os.Getenv("FORGEJO_MIRROR_DESTINATION"),
		username:   os.Getenv("FORGEJO_MIRROR_USERNAME"),
		vaultPath:  os.Getenv("FORGEJO_MIRROR_VAULT_PATH"),
		apiToken:   os.Getenv("FORGEJO_API_TOKEN"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	client := &http.Client{Timeout: 15 * time.Second, CheckRedirect: rejectRedirect}
	if err := config.ensure(ctx, client); err != nil {
		fmt.Fprintln(os.Stderr, "push mirror:", err)
		os.Exit(1)
	}
}

func rejectRedirect(_ *http.Request, _ []*http.Request) error {
	return errors.New("unexpected redirect")
}

func (c mirrorConfig) ensure(ctx context.Context, client *http.Client) error {
	host, err := url.Parse(c.host)
	if err != nil || host.Scheme != "https" || host.Host == "" || host.User != nil ||
		host.Path != "" || host.RawQuery != "" || host.Fragment != "" || !validRepository(c.repository) || !validRepository(c.target) ||
		(c.forge != "github" && c.forge != "gitlab") || c.apiToken == "" {
		return errors.New("invalid Forgejo mirror configuration")
	}
	destination := "https://" + c.forge + ".com/" + c.target + ".git"
	endpoint := c.host + "/api/v1/repos/" + c.repository + "/push_mirrors"
	var inventory []pushMirror
	header, err := c.request(ctx, client, http.MethodGet, endpoint+"?limit=50&page=1", nil, &inventory, http.StatusOK)
	if err != nil {
		return err
	}
	total := len(inventory)
	if value := header.Get("X-Total-Count"); value != "" {
		total, err = strconv.Atoi(value)
		if err != nil {
			return errors.New("invalid mirror inventory count")
		}
	}
	if len(inventory) >= 50 || total != len(inventory) || strings.Contains(header.Get("Link"), `rel="next"`) {
		return errors.New("incomplete mirror inventory")
	}
	var current bool
	for _, entry := range inventory {
		if entry.Address == destination && entry.Key == "" && entry.Sync &&
			entry.Interval == "8h0m0s" && entry.Filter == "" {
			if current {
				return errors.New("multiple matching HTTPS mirrors")
			}
			current = true
			continue
		}
		remote, err := url.Parse(entry.Address)
		sameDestination := err == nil && remote.Hostname() == c.forge+".com" &&
			strings.TrimSuffix(remote.Path, ".git") == "/"+c.target
		if sameDestination {
			return errors.New("unexpected mirror for this destination")
		}
	}

	if !current {
		sourceRefs, err := refs(ctx, c.host+"/"+c.repository+".git")
		if err != nil {
			return err
		}
		destinationRefs, err := refs(ctx, destination)
		if err != nil {
			return err
		}
		if err := verifyRefs(sourceRefs, destinationRefs); err != nil {
			return err
		}
		token, err := readToken(ctx, c.vaultPath)
		if err != nil {
			return err
		}
		if c.username == "" {
			return errors.New("mirror username is missing")
		}
		var created pushMirror
		_, err = c.request(ctx, client, http.MethodPost, endpoint, map[string]any{
			"remote_address": destination, "remote_username": c.username, "remote_password": token,
			"interval": "8h0m0s", "sync_on_commit": true, "use_ssh": false,
		}, &created, http.StatusOK)
		if err != nil {
			return err
		}
		if created.Name == "" || created.Address != destination || created.Key != "" ||
			created.Interval != "8h0m0s" || !created.Sync || created.Filter != "" {
			return errors.New("unexpected created mirror configuration")
		}
	}
	fmt.Printf("Verified native HTTPS mirror from %s to %s/%s\n", c.repository, c.forge, c.target)
	return nil
}

func (c mirrorConfig) request(ctx context.Context, client *http.Client, method, endpoint string, payload map[string]any, result any, want int) (http.Header, error) {
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.New("encode mirror request")
		}
		body = bytes.NewReader(raw)
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, errors.New("invalid Forgejo request")
	}
	request.Header.Set("Authorization", "token "+c.apiToken)
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%s Forgejo mirror request failed", method)
	}
	defer response.Body.Close()
	if response.StatusCode != want {
		return nil, fmt.Errorf("%s Forgejo mirror returned HTTP %d", method, response.StatusCode)
	}
	if result != nil && json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(result) != nil {
		return nil, errors.New("invalid Forgejo mirror response")
	}
	return response.Header, nil
}

func refs(ctx context.Context, remote string) (map[string]string, error) {
	git := exec.CommandContext(ctx, "git", "-c", "credential.helper=", "ls-remote", "--refs", remote, "refs/heads/*", "refs/tags/*")
	git.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := git.Output()
	if err != nil {
		return nil, errors.New("could not read remote refs")
	}
	refs := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, errors.New("invalid remote ref response")
		}
		refs[fields[1]] = fields[0]
	}
	return refs, nil
}

func verifyRefs(source, destination map[string]string) error {
	if len(source) == 0 {
		return errors.New("source has no advertised branches or tags")
	}
	var mismatches []string
	for ref, oid := range destination {
		if source[ref] != oid {
			mismatches = append(mismatches, ref)
		}
	}
	if len(mismatches) > 0 {
		sort.Strings(mismatches)
		return fmt.Errorf("native activation would change destination history; reconcile refs: %s", strings.Join(mismatches, ", "))
	}
	return nil
}

func readToken(ctx context.Context, path string) (string, error) {
	address, err := url.Parse(os.Getenv("VAULT_ADDR"))
	if err != nil || address.Scheme != "https" || address.Host == "" || address.User != nil ||
		path == "" || os.Getenv("VAULT_TOKEN") == "" {
		return "", errors.New("invalid Vault configuration")
	}
	if value := os.Getenv("VAULT_SKIP_VERIFY"); value != "" && value != "false" && value != "0" {
		return "", errors.New("Vault certificate verification must be enabled")
	}
	config := vault.DefaultConfig()
	config.DisableRedirects = true
	config.HttpClient.CheckRedirect = rejectRedirect
	client, err := vault.NewClient(config)
	if err != nil {
		return "", errors.New("configure Vault client")
	}
	client.SetMaxRetries(0)
	record, err := client.KVv2("secrets").Get(ctx, path)
	if err != nil {
		return "", errors.New("read mirror token from Vault")
	}
	token, _ := record.Data["token"].(string)
	if token == "" || strings.ContainsAny(token, "\r\n") {
		return "", errors.New("missing or invalid mirror token")
	}
	return token, nil
}
