package evaluation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// configuration declares one evaluation namespace and its bounded workload matrix.
type configuration struct {
	// Owner distinguishes these test instances from ordinary rentals.
	Owner string `json:"owner"`
	// MaxParallel limits simultaneous independent rentals.
	MaxParallel int `json:"max_parallel"`
	// MaxLifetimeSeconds bounds one worker from allocation through retrieval.
	MaxLifetimeSeconds float64 `json:"max_lifetime_seconds"`
	// BootTimeoutSeconds bounds readiness independently for each instance.
	BootTimeoutSeconds float64 `json:"boot_timeout_seconds"`
	// WorkloadTimeoutSeconds bounds generation and retrieval for each seed.
	WorkloadTimeoutSeconds float64 `json:"workload_timeout_seconds"`
	// CleanupTimeoutSeconds gives cancellation-independent deletion its own deadline.
	CleanupTimeoutSeconds float64 `json:"cleanup_timeout_seconds"`
	// PollSeconds sets the observation interval without a cross-worker barrier.
	PollSeconds float64 `json:"poll_seconds"`
	// MinimumReliability optionally overrides only the inherited reliability floor; zero removes it.
	MinimumReliability *float64 `json:"minimum_reliability,omitempty"`
	// Matrix declares compatible axes without repeating complete workloads.
	Matrix []matrixGroup `json:"matrix"`
	// Jobs is the validated expansion retained as a report projection.
	Jobs []job `json:"-"`
}

// matrixGroup crosses compatible axes with shared fixed trial settings.
type matrixGroup struct {
	// ID names this comparison group and prefixes derived workload identities.
	ID string `json:"id"`
	// HardwareProfiles selects canonical rental resource configurations.
	HardwareProfiles []string `json:"hardware_profiles"`
	// ModelProfiles selects canonical installed-model configurations.
	ModelProfiles []string `json:"model_profiles"`
	// Workflows lists compatible workflow purposes without the model-name prefix.
	Workflows []string `json:"workflows"`
	// CaseIDs selects reusable prompt/reference definitions.
	CaseIDs []string `json:"case_ids"`
	// Seeds is the unique fixed trial list shared by every expanded cell.
	Seeds []int64 `json:"seeds"`
	// Sampling optionally overrides controls for every cell in this group.
	Sampling *samplingOverrides `json:"sampling,omitempty"`
}

// rental pools every workload using the same hardware and installed model profile.
type rental struct {
	// ID is the first workload identity, used for labels and lifecycle artifacts.
	ID string `json:"id"`
	// ModelProfile owns the installed model assets shared by this rental.
	ModelProfile string `json:"model_profile"`
	// HardwareProfile owns the resource limits shared by this rental.
	HardwareProfile string `json:"hardware_profile"`
	// Jobs lists native workflows, cases and trials to queue on this installation.
	Jobs []job `json:"workloads"`
}

// job is one reproducible workflow/case cell, with fixed trials.
type job struct {
	// ID is a unique safe job basename and label component.
	ID string `json:"id"`
	// ModelProfile selects the owning model registry entry.
	ModelProfile string `json:"model_profile"`
	// HardwareProfile selects declared offer and storage constraints.
	HardwareProfile string `json:"hardware_profile"`
	// CaseID selects the reusable prompt, references and rubric.
	CaseID string `json:"case_id"`
	// Workflow selects a generated native workflow for this model.
	Workflow string `json:"workflow"`
	// Seeds supplies fixed independent trials in a stable order.
	Seeds []int64 `json:"seeds"`
	// Sampling optionally varies sampling/reference controls without copying a workflow.
	Sampling *samplingOverrides `json:"sampling,omitempty"`
}

// samplingOverrides supplies explicit controls for reproducible ablations.
type samplingOverrides struct {
	// Steps replaces the selected sampler's iteration count.
	Steps *int `json:"steps,omitempty"`
	// CFG replaces classifier-free guidance strength.
	CFG *float64 `json:"cfg,omitempty"`
	// Denoise replaces native img2img strength.
	Denoise *float64 `json:"denoise,omitempty"`
	// AdapterWeight replaces IP-Adapter reference influence when present.
	AdapterWeight *float64 `json:"adapter_weight,omitempty"`
	// AdapterEnd replaces the fraction at which adapter guidance ends.
	AdapterEnd *float64 `json:"adapter_end,omitempty"`
	// StyleWeight replaces the RES4LYF style guide influence when present.
	StyleWeight *float64 `json:"style_weight,omitempty"`
	// FrameSource selects previous-frame feedback or the original reference latent.
	FrameSource *string `json:"frame_source,omitempty"`
	// FrameSeedMode selects fixed noise or incrementing seeds within one GIF.
	FrameSeedMode *string `json:"frame_seed_mode,omitempty"`
}

// Options supplies deployment-owned tools and non-secret asset paths.
type Options struct {
	// DisableCleanup retains rentals and skips the final sweep when explicitly requested.
	DisableCleanup bool
	// Progress receives non-secret lifecycle events and the final report location.
	Progress io.Writer
	// logger serializes progress from parallel workers during one run.
	logger *log.Logger
	// Groups optionally selects named matrix groups; cleanup always ignores it.
	Groups []string
	// ExcludedOffers omits operator-observed failed capacity for this run only.
	ExcludedOffers []int64
	// Configuration is the checked-in run matrix path.
	Configuration string
	// Profiles is the public model and hardware registry path.
	Profiles string
	// Endpoint is the Vast API root; HTTP is accepted only for literal loopback.
	Endpoint string
	// Output is the private results and journal parent directory.
	Output string
	// Cases contains reusable case directories.
	Cases string
	// Workflows contains generated profile directories.
	Workflows string
	// Terraform renders the existing template without managing resources or state.
	Terraform string
	// StartupTemplate is the canonical image-generation startup template.
	StartupTemplate string
	// RuntimeDefaults is the canonical ComfyUI runtime declaration.
	RuntimeDefaults string
	// OfferDefaults supplies the shared image-compatible offer filters.
	OfferDefaults string
	// Requirements is the approved hash-locked node requirements path.
	Requirements string
	// SequenceSource supplies the first-party Prompt Sequence node.
	SequenceSource string
	// GIFSource supplies the first-party chained GIF node.
	GIFSource string
	// QwenGIFSource supplies model-specific image-edit feedback without enlarging other bundles.
	QwenGIFSource string
	// Bridge is the existing authenticated CLI bridge executable.
	Bridge string
	// CLI is the existing pinned upstream ComfyUI command.
	CLI string
	// CAFile supplies the official Vast root certificate.
	CAFile string
	// Worker optionally selects another explicitly configured workload executable.
	Worker string
	// WorkerArgs contains literal arguments for the configured workload executable.
	WorkerArgs []string
	// StartupFile optionally selects an explicitly supplied startup recipe.
	StartupFile string
}

// registry resolves approved profiles without independently maintaining their pins.
type registry struct {
	// Image is the existing approved image reference.
	Image string `json:"image"`
	// Models maps public model profile names to their native download definitions.
	Models map[string]model `json:"model_profiles"`
	// Hardware maps hardware names to declared limits.
	Hardware map[string]hardware `json:"hardware_profiles"`
	// Downloads owns each model file's publisher URL and destination kind.
	Downloads map[string]asset `json:"downloads"`
	// Nodes owns approved custom-node repository revisions.
	Nodes map[string]nodeSource `json:"custom_nodes"`
}

// model selects already-approved model assets and node packs.
type model struct {
	// Downloads lists registry filenames in the declared order.
	Downloads []string `json:"downloads"`
	// Nodes lists approved node packs.
	Nodes []string `json:"custom_nodes"`
	// CheckpointDirectory distinguishes complete SDXL bundles from denoisers.
	CheckpointDirectory string `json:"checkpoint_directory"`
	// Requirements selects the existing approved dependency lock.
	Requirements string `json:"python_requirements"`
}

// hardware declares the resource and price bounds for an offer.
type hardware struct {
	// MinimumVRAM is the minimum per-GPU VRAM in GiB.
	MinimumVRAM float64 `json:"min_gpu_ram_gb"`
	// Disk is the allocated storage in GB.
	Disk float64 `json:"disk_gb"`
	// MaximumPrice is the maximum API-reported hourly price.
	MaximumPrice float64 `json:"max_hourly_price"`
	// GPU optionally restricts the GPU model.
	GPU string `json:"gpu_name"`
}

// asset is one model file in the owning public registry.
type asset struct {
	// URL is the immutable publisher URL.
	URL string `json:"url"`
	// Kind selects the native model directory.
	Kind string `json:"kind"`
}

// nodeSource is one previously approved custom-node revision.
type nodeSource struct {
	// URL is the public Git repository URL.
	URL string `json:"url"`
	// Ref is the full immutable commit.
	Ref string `json:"ref"`
}

// readJSON decodes a bounded local configuration with contextual errors.
func readJSON(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("read JSON %q: %w", path, err)
	}
	data, readErr := io.ReadAll(io.LimitReader(file, (16<<20)+1))
	if err := errors.Join(readErr, file.Close()); err != nil {
		return fmt.Errorf("read JSON %q: %w", path, err)
	}
	if len(data) > 16<<20 {
		return fmt.Errorf("JSON %q exceeds 16 MiB", path)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode JSON %q: %w", path, err)
	}
	return nil
}

// loadConfiguration rejects unsafe namespaces and unbounded operations before allocation.
func loadConfiguration(path string, selected []string) (configuration, error) {
	var configuration configuration
	if err := readJSON(path, &configuration); err != nil {
		return configuration, fmt.Errorf("load evaluation configuration: %w", err)
	}
	if len(selected) > 0 {
		requested := map[string]bool{}
		for _, id := range selected {
			if requested[id] {
				return configuration, fmt.Errorf("matrix group %q was selected twice", id)
			}
			requested[id] = true
		}
		var groups []matrixGroup
		seen := map[string]bool{}
		for _, group := range configuration.Matrix {
			if seen[group.ID] {
				return configuration, fmt.Errorf("matrix requires unique group names")
			}
			seen[group.ID] = true
			if requested[group.ID] {
				groups = append(groups, group)
				delete(requested, group.ID)
			}
		}
		if len(requested) > 0 {
			return configuration, fmt.Errorf("selection includes an unknown matrix group")
		}
		configuration.Matrix = groups
	}
	safe := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)
	if !safe.MatchString(configuration.Owner) || configuration.MaxParallel < 1 || configuration.MaxParallel > 16 || configuration.MaxLifetimeSeconds <= 0 || configuration.MaxLifetimeSeconds > 86400 || configuration.BootTimeoutSeconds <= 0 || configuration.WorkloadTimeoutSeconds <= 0 || configuration.CleanupTimeoutSeconds <= 0 || configuration.CleanupTimeoutSeconds > 600 || configuration.PollSeconds <= 0 {
		return configuration, fmt.Errorf("evaluation requires a safe namespace and positive bounded concurrency, lifetimes and timeouts")
	}
	if configuration.MinimumReliability != nil && (*configuration.MinimumReliability < 0 || *configuration.MinimumReliability > 1) {
		return configuration, fmt.Errorf("minimum reliability must be between zero and one")
	}
	jobs, err := expandMatrix(configuration.Matrix, safe)
	if err != nil {
		return configuration, fmt.Errorf("expand evaluation matrix: %w", err)
	}
	configuration.Jobs = jobs
	seen := map[string]bool{}
	if configuration.BootTimeoutSeconds > configuration.MaxLifetimeSeconds || configuration.WorkloadTimeoutSeconds > configuration.MaxLifetimeSeconds || configuration.PollSeconds > 60 || len(configuration.Jobs) > 128 {
		return configuration, fmt.Errorf("evaluation timeouts or job count exceed their configured bounds")
	}
	for _, job := range configuration.Jobs {
		if !safe.MatchString(job.ID) || seen[job.ID] || !safe.MatchString(job.ModelProfile) || !safe.MatchString(job.HardwareProfile) || !safe.MatchString(job.CaseID) || !safe.MatchString(job.Workflow) || len(job.Seeds) == 0 {
			return configuration, fmt.Errorf("job %q requires unique safe names and fixed seeds", job.ID)
		}
		seen[job.ID] = true
		if controls := job.Sampling; controls != nil {
			if controls.FrameSource != nil && *controls.FrameSource != "previous" && *controls.FrameSource != "original" {
				return configuration, fmt.Errorf("job %q has an invalid frame source", job.ID)
			}
			if controls.FrameSeedMode != nil && *controls.FrameSeedMode != "increment" && *controls.FrameSeedMode != "fixed" {
				return configuration, fmt.Errorf("job %q has an invalid frame seed mode", job.ID)
			}
			if controls.Steps != nil && (*controls.Steps < 1 || *controls.Steps > 150) {
				return configuration, fmt.Errorf("job %q requires 1–150 sampling steps", job.ID)
			}
			for name, entry := range map[string]struct {
				// value optionally selects one bounded numeric control.
				value *float64
				// maximum is the supported upper bound for the control.
				maximum float64
			}{"cfg": {controls.CFG, 30}, "denoise": {controls.Denoise, 1}, "adapter_weight": {controls.AdapterWeight, 2}, "adapter_end": {controls.AdapterEnd, 1}, "style_weight": {controls.StyleWeight, 2}} {
				if entry.value != nil && (*entry.value < 0 || *entry.value > entry.maximum) {
					return configuration, fmt.Errorf("job %q has an out-of-range %s", job.ID, name)
				}
			}
		}
		if len(job.Seeds) > 100 {
			return configuration, fmt.Errorf("job %q has more than 100 trials", job.ID)
		}
		seenSeeds := map[int64]bool{}
		for _, seed := range job.Seeds {
			if seenSeeds[seed] {
				return configuration, fmt.Errorf("job %q requires unique seeds", job.ID)
			}
			seenSeeds[seed] = true
			if seed < 0 || seed > 1<<53-1 {
				return configuration, fmt.Errorf("job %q requires nonnegative exactly representable seeds", job.ID)
			}
		}
	}
	return configuration, nil
}

// writeJSON atomically publishes one private journal or result document.
func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode JSON %q: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create JSON directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".evaluation-*")
	if err != nil {
		return fmt.Errorf("create JSON temporary: %w", err)
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("write JSON temporary: %w", errors.Join(err, temporary.Close()))
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync JSON temporary: %w", errors.Join(err, temporary.Close()))
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close JSON temporary: %w", err)
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("publish JSON %q: %w", path, err)
	}
	return nil
}

// expandMatrix bounds and validates every axis before forming reproducible coordinates.
func expandMatrix(groups []matrixGroup, safe *regexp.Regexp) ([]job, error) {
	var jobs []job
	seenGroups := map[string]bool{}
	for _, group := range groups {
		if !safe.MatchString(group.ID) || len(group.ID) > 56 || seenGroups[group.ID] {
			return nil, fmt.Errorf("matrix requires unique safe group names of at most 56 characters")
		}
		seenGroups[group.ID] = true
		for name, axis := range map[string][]string{"hardware_profiles": group.HardwareProfiles, "model_profiles": group.ModelProfiles, "workflows": group.Workflows, "case_ids": group.CaseIDs} {
			if len(axis) == 0 {
				return nil, fmt.Errorf("matrix group %q has an empty %s axis", group.ID, name)
			}
			seen := map[string]bool{}
			for _, value := range axis {
				if !safe.MatchString(value) || seen[value] {
					return nil, fmt.Errorf("matrix group %q requires unique safe %s choices", group.ID, name)
				}
				seen[value] = true
			}
		}
		ordinal := 0
		for _, hardware := range group.HardwareProfiles {
			for _, model := range group.ModelProfiles {
				for _, workflow := range group.Workflows {
					for _, caseID := range group.CaseIDs {
						if len(jobs) >= 128 {
							return nil, fmt.Errorf("matrix expansion exceeds 128 workloads")
						}
						ordinal++
						jobs = append(jobs, job{ID: fmt.Sprintf("%s_%03d", group.ID, ordinal), HardwareProfile: hardware, ModelProfile: model, Workflow: model + "-" + workflow, CaseID: caseID, Seeds: group.Seeds, Sampling: group.Sampling})
					}
				}
			}
		}
	}
	return jobs, nil
}

// groupRentals preserves matrix order while sharing each hardware/model installation.
func groupRentals(jobs []job) []rental {
	var rentals []rental
	indexes := map[string]int{}
	for _, job := range jobs {
		key := job.HardwareProfile + ":" + job.ModelProfile
		index, found := indexes[key]
		if !found {
			index = len(rentals)
			indexes[key] = index
			rentals = append(rentals, rental{ID: job.ID, HardwareProfile: job.HardwareProfile, ModelProfile: job.ModelProfile})
		}
		rentals[index].Jobs = append(rentals[index].Jobs, job)
	}
	return rentals
}
