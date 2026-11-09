package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

// main reports protocol failures on stderr and preserves a nonzero exit code.
func main() {
	if err := lookup(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// lookup selects the newest matching tag from Docker Hub's update-ordered page.
func lookup() (resultErr error) {
	var query map[string]string
	if err := json.NewDecoder(io.LimitReader(os.Stdin, 4096)).Decode(&query); err != nil {
		return fmt.Errorf("decode image lookup input: %w", err)
	}
	endpoint := "https://hub.docker.com/v2/repositories/vastai/sd-forge/tags?page_size=100&ordering=last_updated&name=neo-"
	if override := query["endpoint"]; override != "" {
		address, err := url.Parse(override)
		if err != nil {
			return fmt.Errorf("parse fixture endpoint: %w", err)
		}
		ip := net.ParseIP(address.Hostname())
		if address.Scheme != "http" || ip == nil || !ip.IsLoopback() || address.User != nil {
			return fmt.Errorf("fixture endpoint must be literal loopback HTTP")
		}
		endpoint = override
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("construct image lookup: %w", err)
	}
	client := http.Client{Timeout: 30 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("read public Forge tags: %w", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close tag response: %w", err))
		}
	}()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("public Forge tags returned HTTP %d", response.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, (2<<20)+1))
	if err != nil {
		return fmt.Errorf("read tag metadata: %w", err)
	}
	if len(data) > 2<<20 {
		return fmt.Errorf("tag metadata exceeds 2 MiB")
	}
	var page struct {
		// Results are ordered by latest update according to the registry request.
		Results []struct {
			// Name is the published image tag.
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(data, &page); err != nil {
		return fmt.Errorf("decode tag metadata: %w", err)
	}
	pattern := regexp.MustCompile(`^neo-[0-9a-f]+-[0-9]{4}-[0-9]{2}-[0-9]{2}-cuda-12[.]9$`)
	for _, tag := range page.Results {
		if pattern.MatchString(tag.Name) {
			if err := json.NewEncoder(os.Stdout).Encode(map[string]string{"image": "vastai/sd-forge:" + tag.Name}); err != nil {
				return fmt.Errorf("encode selected image: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("no current Forge Neo CUDA 12.9 tag was returned; refusing a different runtime")
}
