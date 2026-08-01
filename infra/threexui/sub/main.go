package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

var hosts = flag.String("hosts", "", "Hostnames for sub servers")
var subId = flag.String("sub_id", "", "Subscription id")
var subFile = flag.String("sub_file", "", "File with subs to fix")

func run() error {
	flag.Parse()
	base_path := os.Getenv("NODE_XUI_BASE_PATH")
	toFix := []string{}
	newline := regexp.MustCompile("\r?\n")
	for i, host := range strings.Split(*hosts, ",") {
		if host == "" {
			continue
		}
		sub := fmt.Sprintf("https://%s%s/sub/%s", host, base_path, *subId)
		log.Printf("Fetching subscription for host %s, %s", host, sub)
		sub = strings.TrimSpace(sub)
		resp, err := http.Get(sub)
		if err != nil {
			return fmt.Errorf("could not get subscription %d: %w", i, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read response body: %w", err)
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("invalid status code %s: %s", resp.Status, string(body))
		}
		bodyParsed, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			return fmt.Errorf("could not parse body as base64: %w", err)
		}
		for _, line := range newline.Split(string(bodyParsed), -1) {
			toFix = append(toFix, line)
		}
	}
	for path := range strings.SplitSeq(*subFile, ",") {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read sub file %s: %w", path, err)
		}
		for _, line := range newline.Split(string(content), -1) {
			toFix = append(toFix, line)
		}
	}
	urls := []*url.URL{}
	for i, line := range toFix {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		curUrl, err := url.Parse(line)
		if err != nil {
			return fmt.Errorf("could not parse url %d: %w", i, err)
		}
		query := curUrl.Query()
		query.Set("security", "tls")
		query.Set("sni", curUrl.Hostname())
		curUrl.RawQuery = query.Encode()
		urls = append(urls, curUrl)
	}
	sort.Slice(urls, func(i, j int) bool { return urls[i].Fragment < urls[j].Fragment })
	for _, url := range urls {
		fmt.Printf("%s\n", url.String())
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}
