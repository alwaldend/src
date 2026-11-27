package main

import "fmt"

// outputFormat selects how to serialize the parsed tree.
type outputFormat string

const (
	// jsonFormat emits the machine-readable document index.
	jsonFormat outputFormat = "json"
	// markdownFormat renders the embedded tree template.
	markdownFormat outputFormat = "markdown"
)

// String supplies the current output format to the flag parser.
func (format *outputFormat) String() string { return string(*format) }

// Set validates and records an output format.
func (format *outputFormat) Set(value string) error {
	switch outputFormat(value) {
	case jsonFormat, markdownFormat:
		*format = outputFormat(value)
		return nil
	default:
		return fmt.Errorf("unsupported output format %q (want json or markdown)", value)
	}
}

// Type supplies the format flag's name in help output.
func (format *outputFormat) Type() string { return "ParseOutputType" }

// readme is the metadata record consumed by JSON clients and the Markdown template.
type readme struct {
	// Path is the absolute source path.
	Path string `json:"path"`
	// PathRelative locates the file relative to the requested root.
	PathRelative string `json:"path_relative"`
	// Dir is the absolute directory containing the source.
	Dir string `json:"dir"`
	// DirRelative locates the directory relative to the root.
	DirRelative string `json:"dir_relative"`
	// DirRelativeSplit contains path components used to indent the Markdown tree.
	DirRelativeSplit []string `json:"dir_relative_split"`
	// DirName names the containing directory.
	DirName string `json:"dir_name"`
	// Data contains the decoded YAML metadata.
	Data map[string]any `json:"data"`
}

// parseConfig describes one traversal in normalized workspace coordinates.
type parseConfig struct {
	// format selects the output encoding.
	format outputFormat
	// useGit limits discovery to files tracked in the Git index.
	useGit bool
	// name selects the filename to match in directories.
	name string
	// root anchors relative output paths.
	root string
	// exclude lists directory subtrees omitted from discovery.
	exclude []string
}
