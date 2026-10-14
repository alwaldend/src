// Package layout reads and writes the fixed Anki configuration hierarchy.
//
// It owns discovery, path-derived deck titles, note deck defaults, supported
// versions, and file I/O. Codecs own syntax/schema conversion. Loading leaves
// sources unchanged; exporting writes supplied objects directly without semantic
// validation, cleanup, staging, or publication. No Anki collection is opened.
package layout
