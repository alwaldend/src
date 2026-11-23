package machine

import (
	"fmt"
	"strings"
)

// GraphSection contains the generated graph of one machine.
type GraphSection struct {
	// ID identifies the section and names its cluster.
	ID string
	// Title labels the cluster in the combined graph.
	Title string
	// DOT is the machine definition exported by Stateless.
	DOT string
	// Calls are the child graphs invoked from this section.
	Calls []GraphCall
}

// GraphCall represents a call from a state into a child component.
type GraphCall struct {
	// From is the state whose action makes the call.
	From State
	// Target contains the called component's graph.
	Target GraphSection
	// Label describes the call.
	Label string
}

// ComposeGraph embeds each Stateless export as a cluster and adds call edges.
// Stateless exports transitions for one machine at a time.
func ComposeGraph(root GraphSection) (string, error) {
	var graph, edges strings.Builder
	graph.WriteString("digraph {\n\tcompound=true;\n\tnode [shape=box, style=rounded];\n\trankdir=TB;\n")
	if err := appendSection(&graph, &edges, root); err != nil {
		return "", fmt.Errorf("compose graph: %w", err)
	}
	graph.WriteString(edges.String())
	graph.WriteString("}\n")
	return graph.String(), nil
}

// appendSection renders the component before visiting its child graphs.
func appendSection(graph, edges *strings.Builder, section GraphSection) error {
	// The pinned Stateless format separates its defaults from the body
	// with a blank line. The composed graph supplies shared defaults.
	_, body, ok := strings.Cut(section.DOT, "\n\n")
	if !ok || !strings.HasPrefix(section.DOT, "digraph {\n") {
		return fmt.Errorf("unexpected DOT header in %s graph", section.ID)
	}
	body, ok = strings.CutSuffix(body, "}\n")
	if !ok {
		return fmt.Errorf("unexpected DOT footer in %s graph", section.ID)
	}
	// Give Stateless's synthetic nodes a unique name in each section
	// and remove Go's method-value suffix from action labels.
	body = strings.NewReplacer(
		"\tinit ", "\t"+section.ID+"_init ",
		"\tlegend ", "\t"+section.ID+"_legend ",
		`-fm\l`, `\l`,
		`-fm"`, `"`,
		"🔒 Internal", "Internal transitions",
		"🚫 Ignored", "Ignored events",
		"🔄 Reentry", "Reentry transitions",
		"🔒 ", "[internal] ",
		"🚫 ", "[ignored] ",
		"🔄 ", "[reentry] ",
	).Replace(body)
	fmt.Fprintf(graph, "\tsubgraph cluster_%s {\n\t\tlabel=%q;\n%s\t}\n", section.ID, section.Title, body)
	for _, call := range section.Calls {
		fmt.Fprintf(edges, "\t%q -> %s_init [lhead=cluster_%s, style=dashed, label=%q];\n", call.From, call.Target.ID, call.Target.ID, call.Label)
		if err := appendSection(graph, edges, call.Target); err != nil {
			return fmt.Errorf("append child graph of %s: %w", section.ID, err)
		}
	}
	return nil
}
