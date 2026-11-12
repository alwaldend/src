package evaluation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// credentials stays in memory and is never part of a journal or report.
type credentials struct {
	// username is the component's Vault-owned HTTPS login name.
	username string
	// password is the component's Vault-owned HTTPS password.
	password string
	// huggingface optionally authorizes already-approved gated model downloads.
	huggingface string
}

// prepareRentals compiles every local input before any allocation is possible.
func prepareRentals(ctx context.Context, options Options, configuration configuration, profiles registry, directory string, login credentials) ([]preparedRental, error) {
	var rentals []preparedRental
	for _, selection := range groupRentals(configuration.Jobs) {
		model, found := profiles.Models[selection.ModelProfile]
		if !found {
			return nil, fmt.Errorf("unknown model profile %q", selection.ModelProfile)
		}
		hardware, found := profiles.Hardware[selection.HardwareProfile]
		if !found || hardware.MinimumVRAM < 32 || hardware.Disk <= 0 || hardware.MaximumPrice <= 0 {
			return nil, fmt.Errorf("rental %q has missing or invalid hardware limits", selection.ID)
		}
		if options.Worker == "" {
			for _, workload := range selection.Jobs {
				if err := validateWorkload(options, workload, model); err != nil {
					return nil, fmt.Errorf("validate workload %q: %w", workload.ID, err)
				}
			}
		}
		jobLogin := login
		if selection.ModelProfile != "flux" {
			jobLogin.huggingface = ""
		}
		startup, fingerprint, err := prepareStartup(ctx, options, profiles, model, selection.ModelProfile == "qwen", filepath.Join(directory, "render-"+selection.ID), jobLogin)
		if err != nil {
			return nil, fmt.Errorf("compile rental %q provisioning: %w", selection.ID, err)
		}
		query := map[string]any{}
		if err := readJSON(options.OfferDefaults, &query); err != nil {
			return nil, fmt.Errorf("read shared offer filters: %w", err)
		}
		for _, name := range []string{"datacenter", "verified"} {
			predicate, ok := query[name].(map[string]any)
			if !ok || predicate["eq"] != true {
				return nil, fmt.Errorf("shared offer declaration must require %s", name)
			}
		}
		if configuration.MinimumReliability != nil {
			if *configuration.MinimumReliability == 0 {
				delete(query, "reliability")
			} else {
				query["reliability"] = map[string]any{"gte": *configuration.MinimumReliability}
			}
		}
		rentals = append(rentals, preparedRental{rental: selection, hardware: hardware, startup: startup, startupFingerprint: fingerprint, query: query})
	}
	return rentals, nil
}

// Check validates real packaged inputs with public dummy login values and no API client.
func Check(ctx context.Context, options Options) error {
	configuration, err := loadConfiguration(options.Configuration, options.Groups)
	if err != nil {
		return fmt.Errorf("configure offline check: %w", err)
	}
	if len(configuration.Jobs) == 0 {
		return fmt.Errorf("evaluation matrix requires at least one workload")
	}
	var profiles registry
	if err := readJSON(options.Profiles, &profiles); err != nil {
		return fmt.Errorf("read offline profiles: %w", err)
	}
	if err := os.MkdirAll(options.Output, 0o700); err != nil {
		return fmt.Errorf("create offline output: %w", err)
	}
	rentals, err := prepareRentals(ctx, options, configuration, profiles, options.Output, credentials{username: "imagegen", password: strings.Repeat("x", 48)})
	if err != nil {
		return fmt.Errorf("validate packaged evals: %w", err)
	}
	var summaries []map[string]any
	for _, prepared := range rentals {
		summaries = append(summaries, map[string]any{"rental": prepared.rental, "startup_bytes": len(prepared.startup), "limit_bytes": 16384})
	}
	if err := writeJSON(filepath.Join(options.Output, "configuration-check.json"), map[string]any{"rentals": summaries, "workloads": configuration.Jobs, "live_operations": 0, "gpu_acceptance": false}); err != nil {
		return fmt.Errorf("save configuration check: %w", err)
	}
	return nil
}
