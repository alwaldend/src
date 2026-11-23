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
}

// Handoff represents a call from a state into another machine.
type Handoff struct {
	// From is the state whose action makes the call.
	From State
	// To is the ID of the called machine's section.
	To string
	// Label describes the call.
	Label string
}

// ComposeGraph embeds each Stateless export as a cluster and adds call edges.
// Stateless exports transitions for one machine at a time.
func ComposeGraph(sections []GraphSection, handoffs []Handoff) (string, error) {
	var graph strings.Builder
	graph.WriteString("digraph {\n\tcompound=true;\n\tnode [shape=box, style=rounded];\n\trankdir=TB;\n")
	for _, section := range sections {
		// The pinned Stateless format separates its defaults from the body
		// with a blank line. The composed graph supplies shared defaults.
		_, body, ok := strings.Cut(section.DOT, "\n\n")
		if !ok || !strings.HasPrefix(section.DOT, "digraph {\n") {
			return "", fmt.Errorf("unexpected DOT header in %s graph", section.ID)
		}
		body, ok = strings.CutSuffix(body, "}\n")
		if !ok {
			return "", fmt.Errorf("unexpected DOT footer in %s graph", section.ID)
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
		fmt.Fprintf(&graph, "\tsubgraph cluster_%s {\n\t\tlabel=%q;\n%s\t}\n", section.ID, section.Title, body)
	}
	for _, handoff := range handoffs {
		fmt.Fprintf(&graph, "\t%q -> %s_init [lhead=cluster_%s, style=dashed, label=%q];\n", handoff.From, handoff.To, handoff.To, handoff.Label)
	}
	graph.WriteString("}\n")
	return graph.String(), nil
}
