// Package dump renders the repository's declared DNS records as documentation
// pages, one per destination view.
//
// The declarations under dnsconfig.json are the owning source, so these pages
// project declared state rather than reproducing a provider snapshot.
package dump

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"git.alwaldend.com/alwaldend/src/infra/dns/internal/lint"
)

// View selects the destination a record is published to.
type View string

const (
	// CloudflareView is the public, globally resolvable zone.
	CloudflareView View = "cloudflare"
	// MikrotikView is the RouterOS-served view for the dc1 network.
	MikrotikView View = "mikrotik"
)

// Page is the documentation page for one view.
type Page struct {
	// Path is the repository-relative path of the generated page.
	Path string
	// Body is the complete rendered Markdown document.
	Body string
}

// selector returns the dsp member that marks a record for this view.
func (v View) selector() (string, error) {
	switch v {
	case CloudflareView:
		return "global", nil
	case MikrotikView:
		return "dc1", nil
	default:
		return "", fmt.Errorf("unsupported DNS view %q", v)
	}
}

// Pages returns every generated page in a stable order.
func Pages() []Page {
	return []Page{
		{Path: "infra/dns/cloudflare_dns.md"},
		{Path: "infra/dns/mikrotik_dns.md"},
	}
}

// viewForPath returns the view a generated page renders.
func viewForPath(path string) (View, bool) {
	switch path {
	case "infra/dns/cloudflare_dns.md":
		return CloudflareView, true
	case "infra/dns/mikrotik_dns.md":
		return MikrotikView, true
	default:
		return "", false
	}
}

// row is a single rendered record value.
type row struct {
	Name  string
	Type  string
	TTL   string
	Value string
	Owner string
}

type member struct {
	Name     string `json:"name"`
	TTL      int    `json:"ttl"`
	Address  string `json:"address"`
	Target   string `json:"target"`
	Content  string `json:"content"`
	Priority *int   `json:"priority"`
}

// value returns the destination value the member carries.
func (m member) value() string {
	switch {
	case m.Address != "":
		return m.Address
	case m.Target != "":
		if m.Priority != nil {
			return fmt.Sprintf("%d %s", *m.Priority, m.Target)
		}
		return m.Target
	default:
		return m.Content
	}
}

// RenderPages renders every generated page from the workspace declarations.
func RenderPages(workspace fs.FS) ([]Page, error) {
	sources, err := lint.DiscoverSources(workspace)
	if err != nil {
		return nil, fmt.Errorf("discover DNS sources: %w", err)
	}
	pages := Pages()
	for i, page := range pages {
		view, ok := viewForPath(page.Path)
		if !ok {
			return nil, fmt.Errorf("no DNS view for page %s", page.Path)
		}
		body, err := render(workspace, sources, view)
		if err != nil {
			return nil, err
		}
		pages[i].Body = body
	}
	return pages, nil
}

// render builds the Markdown document for one view.
func render(workspace fs.FS, sources []string, view View) (string, error) {
	selector, err := view.selector()
	if err != nil {
		return "", err
	}
	rows := map[string]*row{}
	for _, source := range sources {
		data, err := fs.ReadFile(workspace, source)
		if err != nil {
			return "", fmt.Errorf("%s: %w", source, err)
		}
		var document struct {
			Records map[string]map[string]json.RawMessage `json:"records"`
		}
		if err := json.Unmarshal(data, &document); err != nil {
			return "", fmt.Errorf("%s: %w", source, err)
		}
		keys := make([]string, 0, len(document.Records))
		for key := range document.Records {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			members := document.Records[key]
			var views []string
			if err := json.Unmarshal(members["dsp"], &views); err != nil {
				return "", fmt.Errorf("%s record %q: %w", source, key, err)
			}
			if !selected(views, selector) {
				continue
			}
			types := make([]string, 0, len(members))
			for kind := range members {
				if kind != "dsp" {
					types = append(types, kind)
				}
			}
			sort.Strings(types)
			for _, kind := range types {
				var parsed member
				if err := json.Unmarshal(members[kind], &parsed); err != nil {
					return "", fmt.Errorf("%s record %q type %s: %w", source, key, kind, err)
				}
				ttl := "default"
				if parsed.TTL != 0 {
					ttl = fmt.Sprintf("%d", parsed.TTL)
				}
				id := strings.Join([]string{parsed.Name, kind, ttl, parsed.value()}, "\x00")
				item := rows[id]
				if item == nil {
					item = &row{Name: parsed.Name, Type: kind, TTL: ttl, Value: parsed.value(), Owner: source}
					rows[id] = item
				}
			}
		}
	}
	ordered := make([]*row, 0, len(rows))
	for _, item := range rows {
		ordered = append(ordered, item)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Name != ordered[j].Name {
			return ordered[i].Name < ordered[j].Name
		}
		if ordered[i].Type != ordered[j].Type {
			return ordered[i].Type < ordered[j].Type
		}
		return ordered[i].Value < ordered[j].Value
	})

	options := pageOptions(view)
	var page strings.Builder
	fmt.Fprintf(&page, "---\ntitle: %s\nlinkTitle: %s\ndescription: %s\n---\n\n", options.title, options.title, options.description)
	page.WriteString(options.intro)
	page.WriteString("\n\nThis page is generated from the declarations; run `bazel run //infra/dns/cmd/dump -- --write` after changing them. The authoritative source is each owner's `dnsconfig.json`.\n\n")
	page.WriteString("| Domain | Type | TTL | Value | Declaration |\n| --- | --- | --- | --- | --- |\n")
	for _, item := range ordered {
		fmt.Fprintf(&page, "| %s | %s | %s | %s | %s |\n",
			cell(item.Name), cell(item.Type), cell(item.TTL), cell(item.Value), cell(item.Owner))
	}
	return page.String(), nil
}

type pageMeta struct {
	title       string
	description string
	intro       string
}

func pageOptions(view View) pageMeta {
	if view == MikrotikView {
		return pageMeta{
			title:       "Mikrotik DNS",
			description: "Records declared for the internal dc1 view served by RouterOS",
			intro:       "Every record this repository declares for the internal `dc1` view, which RouterOS serves for that network.",
		}
	}
	return pageMeta{
		title:       "Cloudflare DNS",
		description: "Records declared for the public zone served by Cloudflare",
		intro:       "Every record this repository declares for the public zone, which Cloudflare serves.",
	}
}

// selected reports whether a declaration targets this view. "all" expands to
// every view.
func selected(views []string, selector string) bool {
	for _, view := range views {
		if view == selector || view == "all" {
			return true
		}
	}
	return false
}

func cell(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "|", "&#124;", "\n", " ", "\r", " ").Replace(value)
}
