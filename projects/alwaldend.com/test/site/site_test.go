package site_test

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"golang.org/x/net/html"
)

var siteIndex = flag.String("site-index", "", "Generated site index runfile")

// A reference keeps one example rather than every occurrence in the large
// print editions. HTML is tokenized one file at a time, without retaining DOMs.
type reference struct {
	path     string
	fragment string
}

type page struct {
	ids    map[string]bool
	h1     int
	links  map[string]bool
	images map[string]bool
}

type findings struct {
	count    int
	examples []string
}

func (f *findings) add(format string, args ...any) {
	f.count++
	if len(f.examples) < 20 {
		f.examples = append(f.examples, fmt.Sprintf(format, args...))
	}
}

func (f *findings) report(t *testing.T, category string) {
	t.Helper()
	if f.count != 0 {
		t.Errorf("%s: %d findings (first %d shown)\n%s",
			category, f.count, len(f.examples), strings.Join(f.examples, "\n"))
	}
}

func TestGeneratedSite(t *testing.T) {
	runfilesSet, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	if *siteIndex == "" {
		t.Fatal("-site-index is required")
	}
	index, err := runfilesSet.Rlocation(*siteIndex)
	if err != nil {
		t.Fatalf("locate generated site: %v", err)
	}
	root, err := filepath.EvalSymlinks(filepath.Dir(index))
	if err != nil {
		t.Fatalf("resolve generated site directory: %v", err)
	}

	files := map[string]bool{}
	pages := map[string]*page{}
	references := map[reference]string{}
	var markupProblems, invalidURLs findings
	err = filepath.WalkDir(root, func(filename string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, filename)
		if err != nil {
			return err
		}
		filePath := "/" + filepath.ToSlash(relative)
		files[filePath] = true
		if path.Ext(filePath) != ".html" {
			return nil
		}
		parsed, err := inspectPage(filename, filePath, references, &markupProblems, &invalidURLs)
		if err != nil {
			return fmt.Errorf("parse %s: %w", filePath, err)
		}
		pages[filePath] = parsed
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) == 0 {
		t.Fatal("generated site contains no HTML pages")
	}
	markupProblems.report(t, "invalid HTML IDs or image alternatives")
	invalidURLs.report(t, "invalid URLs or obsolete source links")

	// Sorting makes the bounded diagnostics stable across runs.
	ordered := make([]reference, 0, len(references))
	for ref := range references {
		ordered = append(ordered, ref)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].path == ordered[j].path {
			return ordered[i].fragment < ordered[j].fragment
		}
		return ordered[i].path < ordered[j].path
	})
	var missingFiles, missingFragments findings
	for _, ref := range ordered {
		target := ref.path
		if !files[target] {
			// A static server redirects /directory to /directory/ and serves
			// index.html. A directory without an index is not a valid target.
			target = path.Join(target, "index.html")
		}
		if !files[target] {
			missingFiles.add("%s, linked from %s", ref.path, references[ref])
			continue
		}
		fragment := ref.fragment
		// Browsers interpret text-fragment directives separately from IDs.
		fragment, _, _ = strings.Cut(fragment, ":~:text=")
		if fragment == "" {
			continue
		}
		if targetPage := pages[target]; targetPage != nil && !targetPage.ids[fragment] &&
			!strings.EqualFold(fragment, "top") {
			missingFragments.add("%s#%s, linked from %s", target, fragment, references[ref])
		}
	}
	missingFiles.report(t, "missing internal files")
	missingFragments.report(t, "missing HTML fragments")
	checkKnownRegressions(t, pages)
	t.Logf("checked %d HTML pages, %d generated files and %d unique internal destinations",
		len(pages), len(files), len(references))
}

func inspectPage(filename, filePath string, references map[reference]string,
	markupProblems, invalidURLs *findings,
) (*page, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	parsed := &page{ids: map[string]bool{}}
	seenIDs := map[string]bool{}
	// Only a few pages need their external links and image URLs retained for
	// focused regressions; the complete site still receives structural checks.
	if filePath == "/index.html" || filePath == "/projects/index.html" ||
		strings.HasPrefix(filePath, "/docs/misc/") ||
		strings.Contains(filePath, "/docs/projects/goal/") {
		parsed.links = map[string]bool{}
		parsed.images = map[string]bool{}
	}
	pagePath := filePath
	if strings.HasSuffix(pagePath, "/index.html") {
		pagePath = strings.TrimSuffix(pagePath, "index.html")
	}
	base := &url.URL{Scheme: "https", Host: "alwaldend.com", Path: pagePath}
	z := html.NewTokenizer(file)
	for {
		switch z.Next() {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return parsed, nil
			}
			return nil, z.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			token := z.Token()
			if id, present := attribute(token, "id"); present {
				if id == "" || strings.ContainsAny(id, "\t\n\f\r ") {
					markupProblems.add("%s: invalid id %q", filePath, id)
				}
				if seenIDs[id] {
					markupProblems.add("%s: duplicate id %q", filePath, id)
				}
				seenIDs[id] = true
				parsed.ids[id] = true
			}
			if token.Data == "h1" {
				parsed.h1++
			}
			// Named anchors remain valid fragment destinations.
			if token.Data == "a" {
				if name, present := attribute(token, "name"); present && name != "" {
					parsed.ids[name] = true
				}
			}
			if token.Data == "img" {
				// Empty alt is valid for a decorative image; absence is not.
				if _, present := attribute(token, "alt"); !present {
					src, _ := attribute(token, "src")
					markupProblems.add("%s: image %q has no alt attribute", filePath, src)
				}
			}
			attributeName := ""
			switch token.Data {
			case "a", "link":
				attributeName = "href"
			case "img", "script":
				attributeName = "src"
			}
			if attributeName == "" {
				continue
			}
			raw, present := attribute(token, attributeName)
			if !present {
				continue
			}
			if attributeName == "src" && strings.TrimSpace(raw) == "" {
				invalidURLs.add("%s: %s has an empty src", filePath, token.Data)
			}
			if token.Data == "a" && parsed.links != nil {
				parsed.links[raw] = true
			}
			if token.Data == "img" && parsed.images != nil {
				parsed.images[raw] = true
			}
			if strings.HasPrefix(raw, "https://github.com/alwaldend/src/") &&
				strings.Contains(raw, "/docs/content/misc") {
				invalidURLs.add("%s: obsolete source link %q", filePath, raw)
			}
			ref, internal, err := internalReference(base, raw)
			if err != nil {
				invalidURLs.add("%s: %s=%q: %v", filePath, attributeName, raw, err)
			} else if internal {
				if _, found := references[ref]; !found {
					references[ref] = filePath + " (" + raw + ")"
				}
			}
		}
	}
}

func attribute(token html.Token, name string) (string, bool) {
	for _, attr := range token.Attr {
		if attr.Key == name {
			return attr.Val, true
		}
	}
	return "", false
}

func internalReference(base *url.URL, raw string) (reference, bool, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return reference{}, false, err
	}
	resolved := base.ResolveReference(parsed)
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return reference{}, false, nil
	}
	switch strings.ToLower(resolved.Hostname()) {
	case "alwaldend.com", "www.alwaldend.com", "localhost", "127.0.0.1", "::1":
	default:
		return reference{}, false, nil
	}
	// URL parsing unescapes paths/fragments; queries do not change static files.
	target := path.Clean("/" + strings.TrimPrefix(resolved.Path, "/"))
	if target != "/" && strings.HasSuffix(resolved.Path, "/") {
		target += "/"
	}
	return reference{
		path:     target,
		fragment: resolved.Fragment,
	}, true, nil
}

func checkKnownRegressions(t *testing.T, pages map[string]*page) {
	t.Helper()
	for _, filename := range []string{"/index.html", "/projects/index.html"} {
		parsed := pages[filename]
		if parsed == nil {
			t.Errorf("missing landing page %s", filename)
		} else if parsed.h1 != 1 {
			t.Errorf("%s has %d H1 headings, want 1", filename, parsed.h1)
		}
	}
	for filename, source := range map[string]string{
		"/docs/misc/index.html":         "README.md",
		"/docs/misc/books/index.html":   "books/README.md",
		"/docs/misc/android/index.html": "android.md",
	} {
		parsed := pages[filename]
		if parsed == nil {
			t.Errorf("missing Misc page %s", filename)
			continue
		}
		for _, action := range []string{"tree", "edit"} {
			want := "https://github.com/alwaldend/src/" + action +
				"/master/projects/alwaldend.com/content/docs/misc/" + source
			if !parsed.links[want] {
				t.Errorf("%s is missing corrected source link %s", filename, want)
			}
		}
	}
	for _, filename := range []string{
		"/docs/projects/goal/docs/architecture/index.html",
		"/_print/docs/projects/goal/index.html",
	} {
		parsed := pages[filename]
		if parsed == nil {
			t.Errorf("missing Goal documentation page %s", filename)
			continue
		}
		for _, image := range []string{"goal-loop.svg", "goal-tool.svg"} {
			want := "/docs/projects/goal/docs/" + image
			if !parsed.images[want] {
				t.Errorf("%s is missing image with normal resource permalink %s", filename, want)
			}
		}
	}
}
