package dump

import (
	"strings"
	"testing"
	"testing/fstest"
)

const declaration = `{
  "records": {
    "public": {
      "A": {"name": "host", "address": "192.0.2.1"},
      "dsp": ["global"]
    },
    "internal": {
      "A": {"name": "host", "address": "198.51.100.1"},
      "dsp": ["dc1"]
    },
    "both": {
      "TXT": {"name": "host", "content": "v=spf1 -all"},
      "dsp": ["all"]
    }
  }
}`

func workspace() fstest.MapFS {
	return fstest.MapFS{
		"projects/example/dnsconfig.json": &fstest.MapFile{Data: []byte(declaration)},
	}
}

func TestRenderPagesSeparatesViews(t *testing.T) {
	pages, err := RenderPages(workspace())
	if err != nil {
		t.Fatalf("RenderPages: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(pages))
	}
	cloudflare, mikrotik := pages[0].Body, pages[1].Body
	if !strings.Contains(cloudflare, "192.0.2.1") {
		t.Error("cloudflare page is missing the global record")
	}
	if strings.Contains(cloudflare, "198.51.100.1") {
		t.Error("cloudflare page contains a dc1-only record")
	}
	if !strings.Contains(mikrotik, "198.51.100.1") {
		t.Error("mikrotik page is missing the dc1 record")
	}
	if strings.Contains(mikrotik, "192.0.2.1") {
		t.Error("mikrotik page contains a global-only record")
	}
	// A record declared for every view appears in both pages.
	for name, page := range map[string]string{"cloudflare": cloudflare, "mikrotik": mikrotik} {
		if !strings.Contains(page, "v=spf1 -all") {
			t.Errorf("%s page is missing the all-views record", name)
		}
	}
}

func TestRenderPagesAreStable(t *testing.T) {
	first, err := RenderPages(workspace())
	if err != nil {
		t.Fatalf("RenderPages: %v", err)
	}
	second, err := RenderPages(workspace())
	if err != nil {
		t.Fatalf("RenderPages: %v", err)
	}
	for i := range first {
		if first[i].Body != second[i].Body {
			t.Errorf("page %s is not rendered deterministically", first[i].Path)
		}
	}
}
