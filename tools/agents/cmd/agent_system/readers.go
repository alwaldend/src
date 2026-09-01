package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	catalogv1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
)

func (c *capsuleBuilder) readRaw(path string) ([]byte, error) {
	return os.ReadFile(filepath.Join(c.root, filepath.FromSlash(path)))
}

func (c *capsuleBuilder) readTopology() (*catalogv1alpha1.TopologyCatalog, error) {
	path := "tools/agents/catalogs/topology.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeTopologyStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readPolicy() (*catalogv1alpha1.PolicyCatalog, error) {
	path := "tools/agents/catalogs/policy.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodePolicyStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readAction() (*catalogv1alpha1.ActionCatalog, error) {
	path := "tools/agents/catalogs/action.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeActionStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readCapability() (*catalogv1alpha1.CapabilityCatalog, error) {
	path := "tools/agents/catalogs/capability.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeCapabilityStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readWorkspaceCheck() (*catalogv1alpha1.WorkspaceCheckCatalog, error) {
	path := "tools/agents/catalogs/workspace-check.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeWorkspaceCheckStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readGoal() (*catalogv1alpha1.GoalCatalog, error) {
	path := "tools/agents/catalogs/goal.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeGoalStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readIndex() (*catalogv1alpha1.AgentSystemIndex, error) {
	path := "tools/agents/catalogs/index.json"
	content, err := c.readRaw(path)
	if err != nil {
		return nil, err
	}
	catalog, err := catalogv1alpha1.DecodeIndexStrict(content)
	if err != nil {
		return nil, err
	}
	c.input(path, "catalog", digestBytes(content))
	return &catalog, nil
}

func (c *capsuleBuilder) readStrictDecode(path string, destination any) error {
	content, err := c.readRaw(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("trailing JSON")
	}
	return nil
}

func digestBytes(content []byte) string {
	value := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(value[:])
}
