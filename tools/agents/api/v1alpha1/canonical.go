package v1alpha1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func DecodeOperationContract(content []byte) (OperationContract, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var contract OperationContract
	if err := decoder.Decode(&contract); err != nil {
		return OperationContract{}, fmt.Errorf("decode operation contract: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return OperationContract{}, err
	}
	if err := contract.Validate(); err != nil {
		return OperationContract{}, err
	}
	return contract, nil
}

func CanonicalOperationJSON(contract OperationContract) ([]byte, error) {
	if err := contract.Validate(); err != nil {
		return nil, err
	}
	content, err := json.Marshal(contract)
	if err != nil {
		return nil, fmt.Errorf("encode operation contract: %w", err)
	}
	return append(content, '\n'), nil
}

func DecodeArtifactEnvelope(content []byte) (ArtifactEnvelope, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var envelope ArtifactEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		return ArtifactEnvelope{}, fmt.Errorf("decode artifact envelope: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return ArtifactEnvelope{}, err
	}
	if err := envelope.Validate(); err != nil {
		return ArtifactEnvelope{}, err
	}
	return envelope, nil
}

func CanonicalArtifactJSON(envelope ArtifactEnvelope) ([]byte, error) {
	if err := envelope.Validate(); err != nil {
		return nil, err
	}
	content, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("encode artifact envelope: %w", err)
	}
	return append(content, '\n'), nil
}

func DecodeTaskRunManifest(content []byte) (TaskRunManifest, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var manifest TaskRunManifest
	if err := decoder.Decode(&manifest); err != nil {
		return TaskRunManifest{}, fmt.Errorf("decode task-run manifest: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return TaskRunManifest{}, err
	}
	if err := manifest.Validate(); err != nil {
		return TaskRunManifest{}, err
	}
	return manifest, nil
}

func CanonicalTaskRunJSON(manifest TaskRunManifest) ([]byte, error) {
	if err := manifest.Validate(); err != nil {
		return nil, err
	}
	content, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("encode task-run manifest: %w", err)
	}
	return append(content, '\n'), nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("decode trailing JSON: %w", err)
	}
	return fmt.Errorf("operation contract contains multiple JSON values")
}
