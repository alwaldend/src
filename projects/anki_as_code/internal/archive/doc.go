// Package archive manages offline collection files and Anki archive resources.
//
// It streams ZIP and Zstandard payloads, creates private database copies, and
// publishes completed files with integrity and input-change checks. It does not
// interpret collection entities or execute SQLite queries.
package archive
