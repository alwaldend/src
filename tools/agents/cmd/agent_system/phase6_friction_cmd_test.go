package agent_system

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestRunFrictionAggregatesRecords(t *testing.T) {
	records := []v1alpha1.FrictionRecord{
		fixtureFrictionRecord(t, "friction/phase6-test/alpha", "alpha-defect", 2, 3, 120000),
		fixtureFrictionRecord(t, "friction/phase6-test/beta", "beta-defect", 1, 1, 30000),
		fixtureFrictionRecord(t, "friction/phase6-test/gamma", "alpha-defect", 1, 2, 60000),
	}
	directory := t.TempDir()
	inputs := make([]string, 0, len(records))
	for index, record := range records {
		content, err := v1alpha1.CanonicalFrictionRecordJSON(record)
		if err != nil {
			t.Fatalf("canonicalize record %d: %v", index, err)
		}
		path := filepath.Join(
			directory,
			fmt.Sprintf("record-%d.json", index),
		)
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
		inputs = append(inputs, path)
	}
	outputPath := filepath.Join(directory, "aggregate.json")
	args := append([]string{"--output", outputPath}, "--input")
	for _, input := range inputs {
		args = append(args, input, "--input")
	}
	args = args[:len(args)-1]
	var stdout syncBuffer
	if err := runFriction(args, &stdout); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var aggregate frictionAggregate
	if err := json.Unmarshal(content, &aggregate); err != nil {
		t.Fatalf("invalid aggregate JSON: %v", err)
	}
	if aggregate.TotalRecords != 3 {
		t.Fatalf("totalRecords = %d, want 3", aggregate.TotalRecords)
	}
	if aggregate.UniqueSignatures != 2 {
		t.Fatalf("uniqueSignatures = %d, want 2", aggregate.UniqueSignatures)
	}
	if aggregate.TotalAvoidableReads != 4 {
		t.Fatalf("totalAvoidableReads = %d, want 4", aggregate.TotalAvoidableReads)
	}
	if aggregate.TotalAvoidableCommands != 6 {
		t.Fatalf(
			"totalAvoidableCommands = %d, want 6",
			aggregate.TotalAvoidableCommands,
		)
	}
	if aggregate.TotalLatencyMS != 210000 {
		t.Fatalf("totalLatencyMS = %d, want 210000", aggregate.TotalLatencyMS)
	}
	if len(aggregate.Groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(aggregate.Groups))
	}
	if aggregate.Groups[0].DefectSignature != "alpha-defect" {
		t.Fatalf(
			"first group = %s, want alpha-defect (higher measured cost)",
			aggregate.Groups[0].DefectSignature,
		)
	}
}

func TestRunFrictionRejectsDuplicateIDs(t *testing.T) {
	record := fixtureFrictionRecord(
		t, "friction/phase6-test/duplicate", "duplicate-defect", 1, 1, 1000,
	)
	directory := t.TempDir()
	content, err := v1alpha1.CanonicalFrictionRecordJSON(record)
	if err != nil {
		t.Fatal(err)
	}
	firstPath := filepath.Join(directory, "first.json")
	secondPath := filepath.Join(directory, "second.json")
	for _, path := range []string{firstPath, secondPath} {
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var stdout syncBuffer
	err = runFriction([]string{
		"--input", firstPath,
		"--input", secondPath,
		"--output", filepath.Join(directory, "aggregate.json"),
	}, &stdout)
	if err == nil {
		t.Fatal("expected duplicate record ID to fail")
	}
}

func fixtureFrictionRecord(
	t *testing.T,
	id string,
	defectSignature string,
	avoidableReads int,
	avoidableCommands int,
	latencyMS int64,
) v1alpha1.FrictionRecord {
	t.Helper()
	return v1alpha1.FrictionRecord{
		APIVersion: "agents.alwaldend.com/v1alpha1",
		Kind:       "FrictionRecord",
		ID:         id,
		TaskRef: v1alpha1.Reference{
			Kind: v1alpha1.ReferenceTask,
			ID:   "task/phase-six-test",
		},
		AvoidableReads:    avoidableReads,
		AvoidableCommands: avoidableCommands,
		LatencyMS:         latencyMS,
		EvidenceRef: v1alpha1.Reference{
			Kind: v1alpha1.ReferenceArtifact,
			ID:   "artifact/test-" + id,
		},
		EvidenceTier:    v1alpha1.TierFixtureTested,
		DefectSignature: defectSignature,
	}
}

type syncBuffer struct {
	data []byte
}

func (buffer *syncBuffer) Write(value []byte) (int, error) {
	buffer.data = append(buffer.data, value...)
	return len(value), nil
}
