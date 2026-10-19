//go:build !anki_cgo

package main

// nativeServerRunner records that this executable has no linked native engine.
var nativeServerRunner func([]string) error
