// Package lint checks domain ownership directly from repository DNS files.
package lint

// Record identifies a domain declaration without transforming provider values.
type Record struct {
	Source string
	Key    string
	Name   string
	Type   string
	Views  []string
}

// Report contains every discovered file, including files with no records.
type Report struct {
	Sources []string
	Records []Record
}
