// Package collectionpb defines editable resource configs, configuration sets, plans, and persisted state.
//
// Config messages describe editable TOML resources, ConfigSet messages group them,
// change-plan messages describe reconciliation actions, and state messages bind
// source paths to native identities. All are generated from
// collection.proto. Deck and note-type configuration reuse upstream Anki schemas.
package collectionpb
