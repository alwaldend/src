// Package mmdc asserts the contract of the interactive Mermaid CLI: a
// run applies the checked-in repository theme, and it measures label text
// through the repository's pinned fonts rather than through whatever
// Fontconfig the host provides.
//
// The font check is the load-bearing one. Mermaid asks Chrome for text
// metrics, so a CLI run that followed the host would render the same diagram
// with different geometry on every machine.
package mmdc
