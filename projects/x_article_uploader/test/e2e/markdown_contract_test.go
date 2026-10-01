package e2e

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draft"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/draftjs"
	"git.alwaldend.com/alwaldend/src/projects/x_article_uploader/internal/markdown"
)

// Markdown contract failure modes, specified before implementation:
//   - code and tables reach the draft endpoint with immutable entities, unlike
//     the mutable Markdown entities used by the X Articles composer;
//   - changing mutability loses the original source, breaks the atomic entity
//     reference, or changes the mutability of unrelated entity kinds.
func TestMarkdownEntitiesRemainEditableThroughDraftHTTP(t *testing.T) {
	code := "```go\nfmt.Println(\"hello\")\n```"
	table := "| Item | Value |\n| --- | --- |\n| example | 1 |"
	body := "[Reference](https://example.invalid)\n\n" + code + "\n\n" + table + "\n\n---\n"
	root, postDir, source := bannerPost(t, "", body)
	path, artifactJSON := bannerArtifact(t, root, postDir, source)
	artifact, err := draft.ReadArtifact(path)
	if err != nil {
		t.Fatalf("read artifact: %v", err)
	}
	service, client := newBannerService(t)
	if _, err := draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	exchanges := service.recorded()
	writeBannerEvidence(t, artifactJSON, exchanges)
	if len(exchanges) != 1 || exchanges[0].Path != "/articles/draft" {
		t.Fatalf("expected one draft request, got %#v", exchanges)
	}
	document := exchanges[0].Payload["content_state"].(map[string]any)
	if len(document["entities"].([]any)) != 4 || len(document["blocks"].([]any)) != 4 {
		t.Fatalf("code, table, link, or divider was lost: %#v", document)
	}
	markdownKeys := map[string]bool{}
	seenSource := map[string]bool{}
	for _, raw := range document["entities"].([]any) {
		entity := raw.(map[string]any)
		value := entity["value"].(map[string]any)
		wantMutability := draftjs.MutabilityMutable
		if value["type"] == draftjs.EntityDivider {
			wantMutability = draftjs.MutabilityImmutable
		}
		if value["mutability"] != wantMutability {
			t.Errorf("%s entity mutability = %v, want %s", value["type"], value["mutability"], wantMutability)
		}
		if value["type"] == draftjs.EntityMarkdown {
			markdownKeys[entity["key"].(string)] = false
			payload := value["data"].(map[string]any)["markdown"].(string)
			seenSource[payload] = true
		}
	}
	if len(seenSource) != 2 || !seenSource[code] || !seenSource[table] {
		t.Errorf("Markdown source changed: %#v", seenSource)
	}
	for _, raw := range document["blocks"].([]any) {
		block := raw.(map[string]any)
		if block["type"] != draftjs.BlockAtomic {
			continue
		}
		ranges := block["entity_ranges"].([]any)
		if len(ranges) != 1 || block["text"] != " " {
			t.Fatalf("malformed atomic block: %#v", block)
		}
		selected := ranges[0].(map[string]any)
		if selected["offset"] != float64(0) || selected["length"] != float64(1) {
			t.Errorf("atomic entity range = %#v", selected)
		}
		key := strconv.Itoa(int(selected["key"].(float64)))
		if _, exists := markdownKeys[key]; exists {
			markdownKeys[key] = true
		}
	}
	for key, referenced := range markdownKeys {
		if !referenced {
			t.Errorf("Markdown entity %s has no atomic block", key)
		}
	}
}

// Weighted-length failure modes, specified before implementation:
//   - UTF-8 byte length rejects valid CJK, emoji, or NFC-equivalent text;
//   - short URLs expand beyond the budget while their source bytes fit;
//   - punctuation is mistaken for a URL, or separate Markdown entities bypass
//     the per-article budget;
//   - counting normalizes the emitted source, or a conservative estimate is
//     reported as an exact backend length.
func TestMarkdownWeightedBudgetThroughDraftHTTP(t *testing.T) {
	for _, test := range []struct {
		name     string
		payloads []string
		failure  bool
	}{
		{"ascii_at_budget", []string{strings.Repeat("x", 9492)}, false},
		{"ascii_over_budget", []string{strings.Repeat("x", 9493)}, true},
		{"cjk_at_budget", []string{strings.Repeat("字", 4746)}, false},
		{"cjk_over_budget", []string{strings.Repeat("字", 4747)}, true},
		{"emoji_at_budget", []string{strings.Repeat("😀", 4746)}, false},
		{"emoji_over_budget", []string{strings.Repeat("😀", 4747)}, true},
		{"composed_at_budget", []string{strings.Repeat("é", 9492)}, false},
		{"decomposed_at_budget", []string{strings.Repeat("e\u0301", 9492)}, false},
		{"short_bare_domains", []string{strings.Repeat("a.co ", 400)}, true},
		{"short_urls", []string{strings.Repeat("http://a.co ", 400)}, true},
		{"ordinary_punctuation_at_budget", []string{strings.Repeat("x. ", 3164)}, false},
		{"separate_entities_share_budget", []string{strings.Repeat("x", 5000), strings.Repeat("y", 5000)}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var sources []string
			for _, payload := range test.payloads {
				sources = append(sources, "```\n"+payload+"\n```")
			}
			root, postDir, source := bannerPost(t, "", strings.Join(sources, "\n\n"))
			article, convertErr := markdown.New(postDir, "post").Convert(source)
			if article == nil {
				t.Fatalf("conversion produced no inspectable artifact: %v", convertErr)
			}
			encoded, err := article.JSON()
			if err != nil {
				t.Fatalf("serialize artifact: %v", err)
			}
			var artifactJSON map[string]any
			if err := json.Unmarshal(encoded, &artifactJSON); err != nil {
				t.Fatalf("decode artifact evidence: %v", err)
			}
			writeBannerEvidence(t, artifactJSON, nil)
			if test.failure {
				if convertErr == nil || !hasDiagnostic(diagnosticCodes(article), "markdown-payload-over-budget") {
					t.Fatalf("expected local Markdown budget rejection, got %v", convertErr)
				}
				assertDiagnosticHasPosition(t, article, "markdown-payload-over-budget")
				if !strings.Contains(convertErr.Error(), "weighted-length estimate") || strings.Contains(convertErr.Error(), "bytes") {
					t.Errorf("budget error must identify an estimate, got %v", convertErr)
				}
				return
			}
			if convertErr != nil {
				t.Fatalf("expected conversion within the local weighted budget: %v", convertErr)
			}
			path := filepath.Join(root, "artifact.json")
			writeFile(t, path, encoded)
			artifact, err := draft.ReadArtifact(path)
			if err != nil {
				t.Fatalf("read artifact: %v", err)
			}
			service, client := newBannerService(t)
			if _, err := draft.New(client, root, filepath.Join(root, "cache.json")).CreateDraft(artifact); err != nil {
				t.Fatalf("create draft: %v", err)
			}
			exchanges := service.recorded()
			writeBannerEvidence(t, artifactJSON, exchanges)
			if len(exchanges) != 1 || exchanges[0].Path != "/articles/draft" {
				t.Fatalf("expected one draft request, got %#v", exchanges)
			}
			document := exchanges[0].Payload["content_state"].(map[string]any)
			entities := document["entities"].([]any)
			if len(entities) != len(sources) {
				t.Fatalf("got %d Markdown entities, want %d", len(entities), len(sources))
			}
			for index, raw := range entities {
				value := raw.(map[string]any)["value"].(map[string]any)
				if value["type"] != draftjs.EntityMarkdown || value["data"].(map[string]any)["markdown"] != sources[index] {
					t.Errorf("counting changed original Markdown entity %d", index)
				}
			}
		})
	}
}
