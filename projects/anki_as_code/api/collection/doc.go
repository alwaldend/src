// Package collectionpb defines protobuf models for editable configs, executable
// plans and persisted path-to-native bindings. collection.proto owns the schema;
// generated models reuse upstream Anki deck and note-type messages.
// This package supplies data contracts only. Callers own serialization, file
// layout, semantic validation, planning, native operations and resource lifetime.
package collectionpb
