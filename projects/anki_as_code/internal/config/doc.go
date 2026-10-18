// Package config provides a root-bound Store for Anki configuration files.
// It composes layout discovery/file I/O and codec serialization to Load, Dump
// and Export protobuf configuration sets. Reads return fresh caller-owned values.
// Layout owns fixed filenames, supported versions and path-derived defaults;
// codec owns syntax and schema types. Export overwrites declared files directly
// and retains other entries; failures can leave partial output.
// Callers own semantic validation, state-dependent checks, normalization,
// staging, publication and native lifetime. This package does not open Anki,
// invoke native services or retain parsed messages between operations.
package config
