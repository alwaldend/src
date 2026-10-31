package release

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/api/contracts"
	"google.golang.org/protobuf/encoding/protojson"
)

type Deployer struct{}

type DeployerOpts struct {
	// Receipt selects optional checkpointed local operation evidence.
	Receipt string
	// receipt holds only this invocation's observed progress.
	receipt       *artifactReceipt
	Ctx           context.Context
	Releases      []string
	SorasPath     string
	Environment   string
	SSHDeployment string
	SSHHost       string
	SSHUser       string
	SSHRoot       string
	PublishUser   string
	Project       string
	// Namespace selects the public artifact tree; empty retains manifest defaults.
	Namespace   string
	SiteArchive string
	SSHSteps    []string
	SSHPath     string
	SSHConfig   string
	SSHIdentity string
	// SSHDisableAgent restricts authentication to configured identity files.
	SSHDisableAgent bool
	// SSHCertificate supplies the certificate paired with the SSH identity.
	SSHCertificate string
	SSHKnownHosts  string
	RsyncPath      string
	SSHPort        int
	Timeout        time.Duration
	Stdout         io.Writer
	Stderr         io.Writer
}

func NewDeployer() *Deployer {
	return &Deployer{}
}

func (self *Deployer) Deploy(opts *DeployerOpts) (resultErr error) {
	if len(opts.Releases) == 0 {
		return fmt.Errorf("at least one release_dir is required")
	}
	if opts.Receipt != "" {
		opts.receipt = &artifactReceipt{Schema: "release/artifact_receipt/v1", Status: "running", Operations: []*artifactOperation{}}
		if err := checkpointReceipt(opts); err != nil {
			return fmt.Errorf("initialize artifact receipt: %w", err)
		}
		defer func() {
			opts.receipt.Status = "completed"
			if resultErr != nil {
				opts.receipt.Status = "incomplete"
			}
			if err := checkpointReceipt(opts); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("finalize artifact receipt: %w", err))
			}
		}()
	}
	for _, release := range opts.Releases {
		if !filepath.IsAbs(release) {
			if workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY"); workspace != "" {
				release = filepath.Join(workspace, release)
			}
		}
		err := self.deployRelease(opts, release)
		if err != nil {
			return fmt.Errorf("could not deploy release %s: %w", release, err)
		}
	}
	return nil
}

func (self *Deployer) deployRelease(opts *DeployerOpts, releasePath string) error {
	manifestPath := filepath.Join(releasePath, "release.json")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("could not read manifest file %s: %w", manifestPath, err)
	}
	manifest := &contracts.Release{}
	if err := protojson.Unmarshal(content, manifest); err != nil {
		return fmt.Errorf("could not parse manifest %s: %w", manifestPath, err)
	}
	if opts.SSHDeployment != "" || opts.SSHHost != "" {
		deployment := &contracts.ReleaseDeploymentSsh{}
		if opts.SSHDeployment != "" {
			var err error
			deployment, err = readSSHDeployment(opts.SSHDeployment)
			if err != nil {
				return fmt.Errorf("configure release destination: %w", err)
			}
		}
		if err := self.deploySSH(opts, manifest, releasePath, deployment, manifest.Items); err != nil {
			return fmt.Errorf("deploy manifest over SSH: %w", err)
		}
		return nil
	}
	if len(opts.SSHSteps) == 1 && opts.SSHSteps[0] == "delete" && opts.Environment == "" {
		return fmt.Errorf("deletion requires an explicit SSH destination and environment")
	}
	selected := false
	sshGroups := map[string][]*contracts.ReleaseItem{}
	sshDeployments := map[string]*contracts.ReleaseDeploymentSsh{}
	for _, item := range manifest.Items {
		for _, deployment := range item.Deployments {
			if deployment.Ssh != nil && deployment.Ssh.Environment == opts.Environment && opts.Environment != "" {
				key := deployment.Ssh.String()
				sshGroups[key] = append(sshGroups[key], item)
				sshDeployments[key] = deployment.Ssh
			}
			if item.File != nil && deployment.Oci != nil {
				if opts.Environment != "" {
					continue
				}
				return fmt.Errorf("OCI publication is disabled: mutable registry tags cannot enforce immutable project releases; select an explicit SSH environment and destination")
			}
		}
	}
	for key, items := range sshGroups {
		selected = true
		if err := self.deploySSH(opts, manifest, releasePath, sshDeployments[key], items); err != nil {
			return fmt.Errorf("deploy SSH destination %q: %w", sshDeployments[key].Host, err)
		}
	}
	if !selected {
		return fmt.Errorf("no deployment selected; SSH requires an explicit environment")
	}
	return nil
}

// readSSHDeployment loads the existing deployment metadata contract.
func readSSHDeployment(filename string) (*contracts.ReleaseDeploymentSsh, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read SSH deployment %q: %w", filename, err)
	}
	wrapper := &contracts.ReleaseDeployment{}
	if err := protojson.Unmarshal(data, wrapper); err != nil {
		return nil, fmt.Errorf("parse SSH deployment %q: %w", filename, err)
	}
	if wrapper.Ssh == nil {
		return nil, fmt.Errorf("deployment %q does not contain SSH configuration", filename)
	}
	return wrapper.Ssh, nil
}
