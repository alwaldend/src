package release

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// artifactReceipt records observations of one invocation, without granting authority.
type artifactReceipt struct {
	// Schema identifies the local evidence format.
	Schema string `json:"schema"`
	// Status is running, completed, or incomplete.
	Status string `json:"status"`
	// UpdatedAt records the latest persisted observation in UTC.
	UpdatedAt string `json:"updated_at"`
	// Operations lists destinations in their execution order.
	Operations []*artifactOperation `json:"operations"`
}

// artifactOperation binds ordered step observations to a resolved publication.
type artifactOperation struct {
	// Namespace selects first-party or external artifacts.
	Namespace string `json:"namespace"`
	// Project identifies the independently published project.
	Project string `json:"project"`
	// Version names the selected release identity.
	Version string `json:"version"`
	// Commit binds first-party artifacts to source.
	Commit string `json:"commit,omitempty"`
	// PublicationDigest hashes the canonical public manifest, including its newline.
	PublicationDigest string `json:"publication_digest"`
	// Destination identifies the resolved endpoint without authentication material.
	Destination artifactDestination `json:"destination"`
	// Steps preserves pending and uncertain steps alongside completed observations.
	Steps []artifactStep `json:"steps"`
}

// artifactDestination identifies the selected content store without credentials.
type artifactDestination struct {
	// Environment is the explicitly selected environment.
	Environment string `json:"environment"`
	// Host is the configured SSH host or alias.
	Host string `json:"host"`
	// Port is the explicit override, or zero for OpenSSH configuration.
	Port int `json:"port,omitempty"`
	// Root is the resolved content root on that host.
	Root string `json:"root"`
}

// artifactStep records confirmation, rather than assuming an error means no effects.
type artifactStep struct {
	// Name selects upload, extract, link, or delete.
	Name string `json:"name"`
	// Status is pending, unknown once started, or completed after confirmation.
	Status string `json:"status"`
}

// newArtifactOperation registers every requested step before a destination is contacted.
func newArtifactOperation(p *sshPublisher) (*artifactOperation, error) {
	if p.opts.receipt == nil {
		return nil, nil
	}
	_, manifest, err := makePublication(p.release, p.target, p.items)
	if err != nil {
		return nil, fmt.Errorf("bind artifact receipt: %w", err)
	}
	digest := sha256.Sum256([]byte(manifest))
	operation := &artifactOperation{
		Namespace: p.target.Namespace, Project: p.target.Project, Version: p.version,
		Commit: p.release.SourceCommit, PublicationDigest: fmt.Sprintf("sha256:%x", digest),
		Destination: artifactDestination{Environment: p.opts.Environment, Host: p.target.Host, Port: p.opts.SSHPort, Root: p.target.Root},
	}
	for _, name := range p.opts.SSHSteps {
		operation.Steps = append(operation.Steps, artifactStep{Name: name, Status: "pending"})
	}
	p.opts.receipt.Operations = append(p.opts.receipt.Operations, operation)
	if err := checkpointReceipt(p.opts); err != nil {
		return nil, fmt.Errorf("record artifact destination: %w", err)
	}
	return operation, nil
}

// checkpointReceipt atomically replaces private local evidence before or after a step.
func checkpointReceipt(opts *DeployerOpts) error {
	if opts.receipt == nil {
		return nil
	}
	opts.receipt.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	data, err := json.MarshalIndent(opts.receipt, "", "  ")
	if err != nil {
		return fmt.Errorf("encode artifact receipt: %w", err)
	}
	destination, err := filepath.Abs(opts.Receipt)
	if err != nil {
		return fmt.Errorf("resolve artifact receipt: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return fmt.Errorf("create artifact receipt directory: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(destination), ".artifact-receipt-*")
	if err != nil {
		return fmt.Errorf("create artifact receipt: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(append(data, '\n')); err != nil {
		_ = file.Close()
		return fmt.Errorf("write artifact receipt: %w", err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync artifact receipt: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close artifact receipt: %w", err)
	}
	if err := os.Rename(file.Name(), destination); err != nil {
		return fmt.Errorf("install artifact receipt: %w", err)
	}
	directory, err := os.Open(filepath.Dir(destination))
	if err != nil {
		return fmt.Errorf("open artifact receipt directory: %w", err)
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	if syncErr != nil {
		return fmt.Errorf("sync artifact receipt directory: %w", syncErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close artifact receipt directory: %w", closeErr)
	}
	return nil
}
