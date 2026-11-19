package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

type hostServer struct {
	ID      string `json:"id"`
	Host    string `json:"host"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
	PoolID  string `json:"poolId"`
}

func findHost(r *rpc, host string) (*hostServer, error) {
	var servers []hostServer
	if err := r.call("server.getAll", map[string]any{}, &servers); err != nil {
		return nil, err
	}
	var found *hostServer
	for _, server := range servers {
		if server.Host != host {
			continue
		}
		if found != nil {
			return nil, errors.New("multiple XO server registrations match XCPNG_HOST")
		}
		copy := server
		found = &copy
	}
	return found, nil
}

// Connected entries are verified without changes; connecting entries are only
// polled. Credentials are used to repair disconnected entries. This bootstrap
// command deliberately does not rotate credentials of healthy connections.
func registerHost(r *rpc, host, username, password string, insecure bool, out io.Writer, pause func(time.Duration)) error {
	if host == "" || strings.ContainsAny(host, "/\\?#@ \t\r\n") || username == "" || password == "" {
		return errors.New("XCPNG_HOST must be a host address and XCPNG_USERNAME/XCPNG_PASSWORD are required")
	}
	server, err := findHost(r, host)
	if err != nil {
		return err
	}
	var id string
	if server == nil {
		if err := r.call("server.add", map[string]any{"host": host, "username": username, "password": password, "autoConnect": true, "allowUnauthorized": insecure}, &id); err != nil {
			return err
		}
		if id == "" {
			return errors.New("server.add returned an empty registration ID")
		}
	} else {
		id = server.ID
		if id == "" {
			return errors.New("matching XO registration has no ID")
		}
		switch server.Status {
		case "connected", "connecting", "disconnected":
		default:
			return errors.New("matching XO registration has an unknown connection status")
		}
		if server.Status == "disconnected" {
			// Stop retries of an enabled-but-disconnected entry before changing
			// their credentials. Connected or connecting entries are never disabled.
			if server.Enabled {
				if err := r.call("server.disable", map[string]string{"id": id}, nil); err != nil {
					return err
				}
			}
			if err := r.call("server.set", map[string]any{"id": id, "username": username, "password": password, "allowUnauthorized": insecure}, nil); err != nil {
				return err
			}
			if err := r.call("server.enable", map[string]string{"id": id}, nil); err != nil {
				return err
			}
		}
	}
	deadline := time.Now().Add(90 * time.Second)
	for attempt := 0; attempt < 30 && time.Now().Before(deadline); attempt++ {
		actual, err := findHost(r, host)
		if err != nil {
			return err
		}
		if actual == nil || actual.ID != id {
			return errors.New("XO host registration changed during connection verification")
		}
		if actual.Status == "connected" && actual.Enabled && actual.PoolID != "" {
			return json.NewEncoder(out).Encode(actual)
		}
		switch actual.Status {
		case "connected", "connecting", "disconnected":
		default:
			return errors.New("XO returned an unknown connection status")
		}
		if attempt < 29 {
			pause(2 * time.Second)
		}
	}
	return errors.New("XO host did not reach connected status with a pool ID within the bounded verification window")
}

func registerHostEnv(r *rpc, out io.Writer) error {
	insecure := os.Getenv("XCPNG_INSECURE")
	if insecure != "" && insecure != "false" && insecure != "true" {
		return errors.New("XCPNG_INSECURE must be true or false")
	}
	return registerHost(r, os.Getenv("XCPNG_HOST"), os.Getenv("XCPNG_USERNAME"), os.Getenv("XCPNG_PASSWORD"), insecure == "true", out, time.Sleep)
}
