package plan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// MaxDocumentBytes bounds encoded JSON independently of the protobuf payload.
// It allows field names, default values, escaping and file indentation; writers
// check the final representation so every successfully saved document is loadable.
const MaxDocumentBytes = 128 << 20

// Document binds a protobuf plan to trusted metadata outside script-controlled data.
// Saved files and HTTP responses use the same envelope and strict protobuf JSON.
type Document struct {
	// AccountID is the operator-selected account validated against the worker response.
	AccountID string
	// EvaluatedAt fixes the time used during this plan's evaluation.
	EvaluatedAt time.Time
	// Plan is the canonical payload; no parallel action or explanation model is kept.
	Plan *workerpb.Plan
}

// Validate checks the metadata and complete typed payload without authorizing X access.
func (d Document) Validate() error {
	if !model.ValidID(d.AccountID) || d.EvaluatedAt.IsZero() || d.EvaluatedAt.Year() < 1 || d.EvaluatedAt.Year() > 9999 {
		return fmt.Errorf("invalid plan account or evaluation time")
	}
	if err := Validate(d.Plan); err != nil {
		return fmt.Errorf("validate document plan: %w", err)
	}
	return nil
}

// Clone isolates mutable protobuf state while preserving trusted metadata.
func (d Document) Clone() Document {
	if d.Plan != nil {
		d.Plan = proto.Clone(d.Plan).(*workerpb.Plan)
	}
	return d
}

// MarshalJSON emits the sole saved-plan contract using canonical protobuf field names.
func (d Document) MarshalJSON() ([]byte, error) {
	if err := d.Validate(); err != nil {
		return nil, fmt.Errorf("encode plan document: %w", err)
	}
	payload, err := (protojson.MarshalOptions{UseProtoNames: true, EmitDefaultValues: true}).Marshal(d.Plan)
	if err != nil {
		return nil, fmt.Errorf("encode protobuf plan: %w", err)
	}
	data, err := json.Marshal(documentJSON{AccountID: d.AccountID, EvaluatedAt: d.EvaluatedAt, Plan: payload})
	if err != nil {
		return nil, fmt.Errorf("encode plan envelope: %w", err)
	}
	if len(data) > MaxDocumentBytes {
		return nil, fmt.Errorf("plan document exceeds JSON byte bound")
	}
	return data, nil
}

// UnmarshalJSON rejects obsolete schemas, unknown fields and malformed protobuf data.
// Failed decoding leaves the receiver unchanged.
func (d *Document) UnmarshalJSON(data []byte) error {
	if len(data) > MaxDocumentBytes {
		return fmt.Errorf("plan document exceeds JSON byte bound")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var envelope documentJSON
	if err := decoder.Decode(&envelope); err != nil {
		return fmt.Errorf("decode plan envelope: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return fmt.Errorf("decode trailing plan content: %w", err)
		}
		return fmt.Errorf("trailing plan content")
	}
	value := &workerpb.Plan{}
	if err := protojson.Unmarshal(envelope.Plan, value); err != nil {
		return fmt.Errorf("decode protobuf plan: %w", err)
	}
	result := Document{AccountID: envelope.AccountID, EvaluatedAt: envelope.EvaluatedAt, Plan: value}
	if err := result.Validate(); err != nil {
		return fmt.Errorf("validate decoded document: %w", err)
	}
	*d = result
	return nil
}

// documentJSON is the serialization projection of trusted metadata and protobuf bytes.
type documentJSON struct {
	// AccountID binds the saved payload to its owner.
	AccountID string `json:"account_id"`
	// EvaluatedAt is the fixed evaluation timestamp.
	EvaluatedAt time.Time `json:"evaluated_at"`
	// Plan uses protobuf JSON rather than encoding/json's Go oneof representation.
	Plan json.RawMessage `json:"plan"`
}
