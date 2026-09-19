// Package renderer_test asserts the observable contract of the Mermaid theme:
// which colors and shapes the renderer applies, how it isolates fonts, and the
// geometry it gives boxed labels and nested containers.
//
// The checks read the rendered SVG rather than the theme JSON wherever the
// renderer projects the setting into the document, so a theme edit that stops
// taking effect fails here even when the config file still looks correct.
package renderer_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/tools/mermaid/test/svgdoc"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	smoke  = flag.String("smoke", "", "runfile path of the smoke diagram SVG")
	edges  = flag.String("edges", "", "runfile path of the edge-label diagram SVG")
	nested = flag.String("nested", "", "runfile path of the nested-cluster diagram SVG")
	theme  = flag.String("theme", "", "runfile path of the theme JSON")
)

// readRunfile resolves a declared runfile argument and returns its contents.
func readRunfile(t *testing.T, label, path string) string {
	t.Helper()
	if path == "" {
		t.Fatalf("%s runfile was not passed", label)
	}
	resolved, err := runfiles.Rlocation(path)
	if err != nil {
		t.Fatalf("resolve %s: %v", label, err)
	}
	content, err := os.ReadFile(resolved)
	if err != nil {
		t.Fatalf("read %s: %v", label, err)
	}
	if len(content) == 0 {
		t.Fatalf("%s is empty", label)
	}
	return string(content)
}

// readSVG reads a rendered diagram and asserts it is a non-empty SVG document.
func readSVG(t *testing.T, label, path string) string {
	t.Helper()
	svg := readRunfile(t, label, path)
	if !strings.Contains(svg, "<svg") {
		t.Fatalf("%s is not an SVG document", label)
	}
	return svg
}

// themeRule returns every occurrence of a CSS rule the theme emits for the
// given selector. Mermaid emits its own rules for some of the same selectors,
// so a check must select the repository rule rather than the first match.
func themeRule(svg, selector string) []string {
	var rules []string
	for _, candidate := range strings.Split(svg, selector+"{")[1:] {
		end := strings.Index(candidate, "}")
		if end < 0 {
			continue
		}
		rules = append(rules, candidate[:end])
	}
	return rules
}

// hasRuleContaining reports whether any rule for the selector carries all of
// the named CSS declarations.
func hasRuleContaining(svg, selector string, fragments ...string) (bool, string) {
	for _, rule := range themeRule(svg, selector) {
		matches := true
		for _, fragment := range fragments {
			if !strings.Contains(rule, fragment) {
				matches = false
				break
			}
		}
		if matches {
			return true, rule
		}
	}
	return false, strings.Join(themeRule(svg, selector), " | ")
}

func TestThemeAppliesRepositoryPalette(t *testing.T) {
	svg := readSVG(t, "smoke", *smoke)

	for _, expected := range []string{"fill:#ffffff", "stroke:#1e1e1e", "rough-node"} {
		if !strings.Contains(svg, expected) {
			t.Errorf("rendered SVG does not contain %q", expected)
		}
	}

	// The built-in Mermaid theme would render lavender nodes on the
	// documentation canvas, so its distinctive colors must not survive.
	for _, forbidden := range []string{"#ECECFF", "#9370DB"} {
		if strings.Contains(svg, forbidden) {
			t.Errorf("rendered SVG still contains the Mermaid default theme color %s", forbidden)
		}
	}

	if !strings.Contains(svg, "font-family:Architects Daughter") {
		t.Error("rendered SVG does not request the pinned handwriting family")
	}
}

// TestFontIsolationIsStable measures the width the renderer computed for a
// known label. Chrome measured that width through the pinned Fontconfig, so a
// width produced by the host's fonts means font isolation stopped working.
//
// The check reads one label rather than the canvas: the canvas also grows with
// the configured node and rank spacing, so it would report a spacing change as
// a font-isolation failure.
func TestFontIsolationIsStable(t *testing.T) {
	svg := readSVG(t, "smoke", *smoke)

	const label = "Bazel-managed renderer"
	// The label sits in a foreignObject whose measured width reaches 181 with
	// the pinned fonts and 177 with the host's default proportional font.
	const wantWidth = `width="181"`

	index := strings.Index(svg, "<p>"+label+"</p>")
	if index < 0 {
		t.Fatalf("fixture no longer renders the %q label", label)
	}
	start := strings.LastIndex(svg[:index], "<foreignObject ")
	if start < 0 {
		t.Fatalf("label %q is not inside a foreignObject", label)
	}
	element := svg[start:index]
	if !strings.Contains(element, wantWidth) {
		t.Errorf("label %q was not measured with the pinned fonts; want %s in %q",
			label, wantWidth, element)
	}
}

func TestEdgeLabelPlates(t *testing.T) {
	svg := readSVG(t, "edges", *edges)

	// Guard against a vacuous check: the fixture must exercise both a labelled
	// and an unlabelled edge, because the two render differently. The graph
	// must therefore carry two edges regardless of how the renderer emits them.
	if count := strings.Count(svg, "flowchart-link\" style="); count != 2 {
		t.Errorf("fixture renders %d edges, want 2", count)
	}
	if !strings.Contains(svg, `<span class="edgeLabel"><p>HTTPS</p></span>`) {
		t.Error("fixture no longer renders a labelled edge")
	}

	// A labelled edge carries a readable plate rather than text sitting on the
	// edge line. The 3px margin keeps the 1px border clear of the frame that
	// clips it: the dagre renderer sizes that frame to the text only, so
	// without the margin the outer half of the border is cut off.
	ok, rule := hasRuleContaining(svg, ".edgeLabel p", "border:1px solid", "margin:3px")
	if !ok {
		t.Errorf("edge label plates are not boxed with border slack: %s", rule)
	}

	// Styling .labelBkg directly would paint a background for unlabelled edges
	// too, so the plate must come from the <p> inside a real label. Mermaid
	// emits an empty plate for those edges under the dagre renderer but omits
	// it under ELK, so the rule is checked rather than any empty element.
	if ok, rule := hasRuleContaining(svg, ".edgeLabel .labelBkg", "background-color:transparent"); !ok {
		t.Errorf("empty edge labels would draw a stray plate: %s", rule)
	}
}

func TestNestedClusterTitles(t *testing.T) {
	svg := readSVG(t, "nested", *nested)

	// Guard against a vacuous check: the fixture must actually nest two
	// containers, which is the condition this test exists to cover.
	if count := strings.Count(svg, `class="cluster"`); count < 2 {
		t.Fatalf("fixture renders %d clusters, want at least 2", count)
	}

	// A nested cluster title gets a boxed plate like an edge label, and it must
	// stay short enough not to reach the container nested inside it. Under the
	// dagre renderer the inset between them is a fixed 20px, so the plate pins
	// line-height:1 to fit; the ELK renderer leaves more room but shares the
	// rule, so the plate stays consistent across renderers.
	ok, rule := hasRuleContaining(svg, ".cluster-label span p", "border:1px solid", "line-height:1", "margin:3px")
	if !ok {
		t.Errorf("nested cluster titles can outgrow the fixed inset: %s", rule)
	}
}

// TestNestedClusterTitlesAreLifted asserts the setting that keeps a parent
// title clear of the container nested inside it. Mermaid does not project this
// key into the SVG, so the check reads the theme file directly.
func TestNestedClusterTitlesAreLifted(t *testing.T) {
	config := readRunfile(t, "theme", *theme)

	// The dagre renderer insets a nested cluster by a fixed 20px that no
	// spacing option widens, so a negative top margin is the only way to lift
	// the title clear of the child border. The ELK renderer has more room and
	// does not need the lift, but both renderers share this one config.
	if !strings.Contains(config, `"top": -16`) {
		t.Error("nested cluster titles are not lifted clear of the child border")
	}
}

// TestFlowchartsUseElkRenderer asserts the layout engine the theme selects. ELK
// gives a nested cluster far more room around its parent title than dagre, whose
// fixed 20px inset is the reason the negative title margin exists at all.
// Mermaid records the chosen renderer in the root element, so a config edit that
// silently reverts the layout fails here.
func TestFlowchartsUseElkRenderer(t *testing.T) {
	for _, fixture := range []struct{ label, path string }{
		{"smoke", *smoke},
		{"edges", *edges},
		{"nested", *nested},
	} {
		svg := readSVG(t, fixture.label, fixture.path)
		if !strings.Contains(svg, `aria-roledescription="flowchart-elk"`) {
			t.Errorf("%s was not laid out by the ELK renderer", fixture.label)
		}
	}
}

// TestClusterTitlesPaintLast asserts the paint order the renderer establishes
// after Mermaid writes the document. Mermaid emits every container title inside
// the clusters group, ahead of the edges, so an edge crossing a container border
// would paint over the title plate. cmd/render/paint_order.mjs lifts each title
// to the end of the root group, and SVG paints in document order, so a title
// that appears before the nodes again means the lift stopped running.
func TestClusterTitlesPaintLast(t *testing.T) {
	svg := readSVG(t, "nested", *nested)

	// A diagram with no container title cannot show whether the lift ran, so
	// the fixture must render at least one before the order is meaningful.
	if count := strings.Count(svg, svgdoc.ClusterTitleMarker); count == 0 {
		t.Fatal("fixture renders no container title to order")
	}

	if !svgdoc.ClusterTitlesPaintLast(svg) {
		t.Error("container titles do not paint after the node layer; " +
			"an edge crossing a container border would cover the title")
	}
}
