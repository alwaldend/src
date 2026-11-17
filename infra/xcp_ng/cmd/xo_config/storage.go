package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// storageObject contains only fields needed to identify storage and attachments.
type storageObject struct {
	// ID is the XO object's stable identifier.
	ID string `json:"id"`
	// Type distinguishes hosts, repositories, VMs and disk objects.
	Type string `json:"type"`
	// Name is the human-readable inventory selector.
	Name string `json:"name_label"`
	// Description carries the bootstrap ownership marker.
	Description string `json:"name_description"`
	// Pool identifies the physical XCP-ng pool.
	Pool string `json:"$pool"`
	// SRType identifies the storage driver.
	SRType string `json:"SR_type"`
	// Shared reports whether the repository is shared by pool hosts.
	Shared bool `json:"shared"`
}

// storageInventory reads inventory without exposing credential-bearing objects.
func storageInventory(r *rpc) (map[string]storageObject, error) {
	var objects map[string]storageObject
	if err := r.call("xo.getAllObjects", map[string]any{}, &objects); err != nil {
		return nil, fmt.Errorf("read storage inventory: %w", err)
	}
	return objects, nil
}

// namedStorageObject requires an unambiguous named object, allowing absence.
func namedStorageObject(objects map[string]storageObject, kind, name string) (*storageObject, error) {
	var found *storageObject
	for _, object := range objects {
		if object.Type != kind || object.Name != name {
			continue
		}
		if found != nil {
			return nil, errors.New("multiple storage inventory objects match the configured selector")
		}
		copy := object
		found = &copy
	}
	if found != nil && found.ID == "" {
		return nil, errors.New("selected storage inventory object has no ID")
	}
	return found, nil
}

// ensureCacheStorage creates an absent SMB SR and preserves matching owned SRs.
// Terraform's backend lock serializes callers; this command does not lock XO.
func ensureCacheStorage(r *rpc, out io.Writer) error {
	hostName, name := os.Getenv("XO_CACHE_HOST"), os.Getenv("XO_CACHE_SR")
	share, user, password := os.Getenv("XO_CACHE_SMB_SHARE"), os.Getenv("XO_CACHE_SMB_USER"), os.Getenv("XO_CACHE_SMB_PASSWORD")
	if hostName == "" || name == "" || !strings.HasPrefix(share, "//") || user == "" || password == "" {
		return errors.New("XO_CACHE_HOST, XO_CACHE_SR and SMB share, user and password are required")
	}
	marker := "Managed by infra/xcp_ng/tf: " + share + " (" + user + ")"
	objects, err := storageInventory(r)
	if err != nil {
		return fmt.Errorf("inspect cache repository: %w", err)
	}
	host, err := namedStorageObject(objects, "host", hostName)
	if err != nil {
		return fmt.Errorf("select cache host: %w", err)
	}
	if host == nil || host.Pool == "" {
		return errors.New("cache host is absent or has no pool")
	}
	sr, err := namedStorageObject(objects, "SR", name)
	if err != nil {
		return fmt.Errorf("select cache repository: %w", err)
	}
	createdID := ""
	if sr == nil {
		if err := r.call("sr.createSmb", map[string]any{
			"host": host.ID, "nameLabel": name, "nameDescription": marker,
			"server": share, "user": user, "password": password,
		}, &createdID); err != nil {
			// A failed response can follow a successful create. Never create twice
			// in one invocation; the next authorized run rediscovers owned storage.
			return fmt.Errorf("create cache repository; inspect before retrying: %w", err)
		}
		if createdID == "" {
			return errors.New("cache repository creation returned no ID")
		}
	}
	for attempt := 0; attempt < 20; attempt++ {
		if sr != nil {
			if sr.SRType != "smb" || !sr.Shared || sr.Pool != host.Pool || sr.Description != marker || (createdID != "" && sr.ID != createdID) {
				return errors.New("cache repository conflicts with its expected type, pool, ownership or creation identity")
			}
			if err := json.NewEncoder(out).Encode(map[string]string{"id": sr.ID, "name": sr.Name}); err != nil {
				return fmt.Errorf("write cache repository result: %w", err)
			}
			return nil
		}
		objects, err = storageInventory(r)
		if err != nil {
			return fmt.Errorf("verify cache repository: %w", err)
		}
		sr, err = namedStorageObject(objects, "SR", name)
		if err != nil {
			return fmt.Errorf("verify unique cache repository: %w", err)
		}
		if sr == nil {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return errors.New("cache repository did not appear within the bounded verification window")
}
