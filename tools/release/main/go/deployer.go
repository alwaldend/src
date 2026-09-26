package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/release/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
)

type Deployer struct{}

type DeployerOpts struct {
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
	SiteArchive   string
	SSHSteps      []string
	SSHPath       string
	SSHConfig     string
	SSHIdentity   string
	SSHKnownHosts string
	RsyncPath     string
	SSHPort       int
	Timeout       time.Duration
	Stdout        io.Writer
	Stderr        io.Writer
}

func NewDeployer() *Deployer {
	return &Deployer{}
}

func (self *Deployer) Deploy(opts *DeployerOpts) error {
	if len(opts.Releases) == 0 {
		return fmt.Errorf("at least one release_dir is required")
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
			data, err := os.ReadFile(opts.SSHDeployment)
			if err != nil {
				return fmt.Errorf("read SSH deployment %q: %w", opts.SSHDeployment, err)
			}
			wrapper := &contracts.ReleaseDeployment{}
			if err := protojson.Unmarshal(data, wrapper); err != nil {
				return fmt.Errorf("parse SSH deployment %q: %w", opts.SSHDeployment, err)
			}
			if wrapper.Ssh == nil {
				return fmt.Errorf("ssh_deployment does not contain SSH configuration")
			}
			deployment = wrapper.Ssh
		}
		if err := self.deploySSH(opts, manifest, releasePath, deployment, manifest.Items); err != nil {
			return fmt.Errorf("deploy manifest over SSH: %w", err)
		}
		return nil
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
				if opts.SorasPath == "" {
					return fmt.Errorf("OCI deployment requires soras_path")
				}
				selected = true
				filePath := filepath.Join(releasePath, "files", item.File.Name)
				for _, tag := range deployment.Oci.Tags {
					cmd := exec.CommandContext(opts.Ctx, opts.SorasPath, "push", tag, filePath)
					cmd.Stdout, cmd.Stderr = opts.Stdout, opts.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("could not run %s: %w", cmd, err)
					}
				}
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
