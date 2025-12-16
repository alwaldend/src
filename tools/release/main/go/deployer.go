package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"git.alwaldend.com/src/tools/release/main/proto/contracts"
	"google.golang.org/protobuf/encoding/protojson"
)

type Deployer struct{}

type DeployerOpts struct {
	Ctx       context.Context
	Releases  []string
	SorasPath string
}

func NewDeployer() *Deployer {
	return &Deployer{}
}

func (self *Deployer) Deploy(opts *DeployerOpts) error {
	for _, release := range opts.Releases {
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
	for _, item := range manifest.Items {
		for _, deployment := range item.Deployments {
			if item.File != nil && deployment.Oci != nil {
				filePath := filepath.Join(releasePath, "files", item.File.Name)
				for _, tag := range deployment.Oci.Tags {
					cmd := exec.Command(opts.SorasPath, "push", tag, filePath)
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("could not run %s: %w", cmd, err)
					}
				}
			}
		}
	}
	return nil
}
