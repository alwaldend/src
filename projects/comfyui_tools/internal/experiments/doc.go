// Package experiments coordinates producer-supplied workflows on an existing
// ComfyUI endpoint through a caller-owned CLI. It owns private run journals,
// duplicate detection, buffered submission and original artifact collection;
// authentication, workflow semantics and instance lifecycles remain external.
package experiments
