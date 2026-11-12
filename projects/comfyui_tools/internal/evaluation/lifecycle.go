package evaluation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
)

// duration converts a validated configuration interval to a context deadline.
func duration(seconds float64) time.Duration { return time.Duration(seconds * float64(time.Second)) }

// prefix is the exact namespace shared by the final and manual cleanup commands.
func prefix(owner string) string { return "imagegen-eval:" + owner + ":" }

// pause waits for an observation interval or cancellation without busy polling.
func pause(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return fmt.Errorf("wait for lifecycle observation: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

// destroy confirms ownership before deletion and requires confirmed absence afterward.
func destroy(ctx context.Context, api *vastclient.Client, id int64, namespace string) error {
	for {
		instance, err := api.Get(ctx, id)
		if errors.Is(err, vastclient.ErrNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("inspect cleanup candidate %d: %w", id, err)
		}
		if !strings.HasPrefix(instance.Label, namespace) {
			return fmt.Errorf("instance %d no longer belongs to the evaluation namespace", id)
		}
		deleteErr := api.Delete(ctx, id)
		_, readErr := api.Get(ctx, id)
		if errors.Is(readErr, vastclient.ErrNotFound) {
			return nil
		}
		if err := pause(ctx, time.Second); err != nil {
			return errors.Join(fmt.Errorf("confirm destruction of %d: %w", id, err), deleteErr, readErr)
		}
	}
}

// cleanup deletes every matching test instance, regardless of run or age.
func cleanup(ctx context.Context, api *vastclient.Client, owner string) error {
	for {
		instances, err := api.Instances(ctx)
		if err != nil {
			return fmt.Errorf("enumerate cleanup namespace: %w", err)
		}
		var candidates []int64
		for _, instance := range instances {
			if strings.HasPrefix(instance.Label, prefix(owner)) {
				candidates = append(candidates, instance.ID)
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i] < candidates[j] })
		var failures []error
		var group sync.WaitGroup
		var guard sync.Mutex
		semaphore := make(chan struct{}, 8)
		for _, id := range candidates {
			semaphore <- struct{}{}
			group.Add(1)
			go func(id int64) {
				defer group.Done()
				defer func() { <-semaphore }()
				if err := destroy(ctx, api, id, prefix(owner)); err != nil {
					guard.Lock()
					failures = append(failures, err)
					guard.Unlock()
				}
			}(id)
		}
		group.Wait()
		if len(failures) > 0 {
			return fmt.Errorf("evaluation cleanup incomplete: %w", errors.Join(failures...))
		}
	}
}

// Cleanup executes the same namespace-wide sweep used at the end of Run.
func Cleanup(ctx context.Context, options Options) error {
	configuration, err := loadConfiguration(options.Configuration, nil)
	if err != nil {
		return fmt.Errorf("configure cleanup: %w", err)
	}
	api, err := vastclient.New(vastclient.Options{Endpoint: options.Endpoint, APIKey: os.Getenv("VAST_API_KEY")})
	if err != nil {
		return fmt.Errorf("configure cleanup client: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, duration(configuration.CleanupTimeoutSeconds))
	defer cancel()
	if err := cleanup(ctx, api, configuration.Owner); err != nil {
		return fmt.Errorf("clean evaluation instances: %w", err)
	}
	if options.Progress != nil {
		log.New(options.Progress, "", 0).Printf("cleanup: namespace %s is empty", configuration.Owner)
	}
	return nil
}

// preparedRental contains validated local inputs compiled before any paid allocation.
type preparedRental struct {
	// rental pools compatible work on one independent installation.
	rental rental
	// hardware supplies approved offer constraints.
	hardware hardware
	// startup is the bounded encoded startup command, never written to reports.
	startup string
	// startupFingerprint binds installed public source, requirements, runtime and template.
	startupFingerprint string
	// query retains the shared image-compatible filters for this allocation.
	query map[string]any
	// retained identifies a validated existing installation instead of an allocation.
	retained *journal
}

// journal records non-secret lifecycle ownership and terminal result status.
type journal struct {
	// Rental preserves the reproducible pooled workload selection.
	Rental rental `json:"rental"`
	// Label is persisted before the rental request for ambiguous-create recovery.
	Label string `json:"label"`
	// InstanceID identifies a confirmed allocation, when available.
	InstanceID int64 `json:"instance_id,omitempty"`
	// Phase records the last observed lifecycle transition.
	Phase string `json:"phase"`
	// StartedAt is the start of this rental attempt.
	StartedAt time.Time `json:"started_at"`
	// CreatedAt records the confirmed rental acknowledgment.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// BootStatus is the last classified provider boot state, without free-form diagnostics.
	BootStatus string `json:"boot_status,omitempty"`
	// ReadyAt separates cold startup from generation time.
	ReadyAt time.Time `json:"ready_at,omitempty"`
	// ReleasedAt records deletion before local media inspection.
	ReleasedAt time.Time `json:"released_at,omitempty"`
	// CompletedAt records the terminal controller outcome after destruction or retention.
	CompletedAt time.Time `json:"completed_at,omitempty"`
	// Failure reports a contextual non-secret error, including cleanup failures.
	Failure string `json:"failure,omitempty"`
	// ReusedFrom names the prior retained run while preserving the original label.
	ReusedFrom string `json:"reused_from,omitempty"`
	// InstallationFingerprint binds the compiled installation without credential values.
	InstallationFingerprint string `json:"installation_fingerprint"`
}

// Run owns one bounded concurrent suite with automatic cleanup enabled by default.
func Run(ctx context.Context, options Options) (resultErr error) {
	if options.ReuseRun != "" && !options.DisableCleanup {
		return errors.New("retained-run reuse requires --cleanup=false; use explicit cleanup to release rentals")
	}
	if options.Progress == nil {
		options.Progress = io.Discard
	}
	options.logger = log.New(options.Progress, "", 0)
	configuration, err := loadConfiguration(options.Configuration, options.Groups)
	if err != nil {
		return fmt.Errorf("configure evaluation: %w", err)
	}
	if len(configuration.Jobs) == 0 {
		return errors.New("evaluation requires at least one job")
	}
	// False reservations exclude uncertain or operator-rejected capacity without waiting for release.
	usedOffers := map[int64]bool{}
	for _, id := range options.ExcludedOffers {
		if _, duplicate := usedOffers[id]; id <= 0 || duplicate {
			return errors.New("excluded offers require unique positive identities")
		}
		usedOffers[id] = false
	}
	if err := os.MkdirAll(options.Output, 0o700); err != nil {
		return fmt.Errorf("create evaluation output: %w", err)
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("locate shared controller lock: %w", err)
	}
	lockDirectory := filepath.Join(cache, "comfyui-evaluations", "locks")
	if err := os.MkdirAll(lockDirectory, 0o700); err != nil {
		return fmt.Errorf("create shared controller lock directory: %w", err)
	}
	lock, err := os.OpenFile(filepath.Join(lockDirectory, configuration.Owner+".lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open controller lock: %w", err)
	}
	defer func() {
		if err := lock.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("release controller lock: %w", err))
		}
	}()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("evaluation namespace already has a local controller: %w", err)
	}
	identifier := make([]byte, 16)
	if _, err := rand.Read(identifier); err != nil {
		return fmt.Errorf("generate evaluation identity: %w", err)
	}
	runID := hex.EncodeToString(identifier)
	directory := filepath.Join(options.Output, runID)
	if err := os.Mkdir(directory, 0o700); err != nil {
		return fmt.Errorf("create run directory: %w", err)
	}
	api, err := vastclient.New(vastclient.Options{Endpoint: options.Endpoint, APIKey: os.Getenv("VAST_API_KEY")})
	if err != nil {
		return fmt.Errorf("configure evaluation client: %w", err)
	}
	defer func() {
		var cleanupErr error
		if options.DisableCleanup {
			options.logger.Printf("run %s: automatic cleanup disabled; rentals retained in namespace %s", runID, configuration.Owner)
		} else {
			options.logger.Printf("run %s: final cleanup of namespace %s", runID, configuration.Owner)
			cleanupContext, cancel := context.WithTimeout(context.Background(), duration(configuration.CleanupTimeoutSeconds))
			defer cancel()
			cleanupErr = cleanup(cleanupContext, api, configuration.Owner)
		}
		resultErr = errors.Join(resultErr, cleanupErr)
		cleanupConfirmed := !options.DisableCleanup && cleanupErr == nil
		report := map[string]any{"run_id": runID, "owner": configuration.Owner, "provider_endpoint": options.Endpoint, "finished_at": time.Now().UTC(), "cleanup_confirmed": cleanupConfirmed, "cleanup_skipped": options.DisableCleanup, "execution_success": resultErr == nil, "quality_review": "unscored", "jobs": configuration.Jobs, "rentals": groupRentals(configuration.Jobs)}
		if options.ReuseRun != "" {
			report["reused_from"] = options.ReuseRun
		}
		if len(options.ExcludedOffers) > 0 {
			report["excluded_offers"] = options.ExcludedOffers
		}
		if resultErr != nil {
			report["failure"] = resultErr.Error()
		}
		if err := writeJSON(filepath.Join(directory, "run.json"), report); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("save run report: %w", err))
		}
		options.logger.Printf("run %s: cleanup confirmed=%t; report %s", runID, cleanupConfirmed, filepath.Join(directory, "run.json"))
	}()
	var profiles registry
	if err := readJSON(options.Profiles, &profiles); err != nil {
		return fmt.Errorf("load approved profiles: %w", err)
	}
	if strings.TrimSpace(profiles.Image) == "" {
		return errors.New("model registry requires an explicit image")
	}
	rentals, err := prepareRentals(ctx, options, configuration, profiles, directory, credentials{username: os.Getenv("IMAGEGEN_WEB_USERNAME"), password: os.Getenv("IMAGEGEN_WEB_PASSWORD"), huggingface: os.Getenv("HF_TOKEN")})
	if err != nil {
		return fmt.Errorf("prepare evaluation inputs: %w", err)
	}
	if options.ReuseRun != "" {
		if err := reuseRentals(ctx, api, options, configuration, profiles, rentals); err != nil {
			return fmt.Errorf("validate retained installations: %w", err)
		}
	}
	if err := writeJSON(filepath.Join(directory, "inputs.json"), map[string]any{"configuration": configuration, "profiles": profiles}); err != nil {
		return fmt.Errorf("save effective public inputs: %w", err)
	}
	options.logger.Printf("run %s: %d rentals serving %d workloads, at most %d active", runID, len(rentals), len(configuration.Jobs), configuration.MaxParallel)
	semaphore := make(chan struct{}, configuration.MaxParallel)
	var group sync.WaitGroup
	var guard sync.Mutex
	var failures []error
	var allocationErr error
	var offersMu sync.Mutex
allocation:
	for _, prepared := range rentals {
		select {
		case semaphore <- struct{}{}:
		case <-ctx.Done():
			allocationErr = fmt.Errorf("evaluation cancelled before allocation: %w", ctx.Err())
			break allocation
		}
		group.Add(1)
		go func(prepared preparedRental) {
			defer group.Done()
			defer func() { <-semaphore }()
			if err := runRental(ctx, api, options, configuration, profiles.Image, prepared, runID, directory, &offersMu, usedOffers); err != nil {
				guard.Lock()
				failures = append(failures, fmt.Errorf("rental %q: %w", prepared.rental.ID, err))
				guard.Unlock()
			}
		}(prepared)
	}
	group.Wait()
	if allocationErr != nil {
		failures = append(failures, allocationErr)
	}
	if len(failures) > 0 {
		return fmt.Errorf("evaluation failures: %w", errors.Join(failures...))
	}
	return nil
}

// runRental queues pooled workloads and releases one rental without a suite barrier.
func runRental(ctx context.Context, api *vastclient.Client, options Options, configuration configuration, image string, prepared preparedRental, runID, directory string, offersMu *sync.Mutex, usedOffers map[int64]bool) (resultErr error) {
	ctx, cancel := context.WithTimeout(ctx, duration(configuration.MaxLifetimeSeconds))
	defer cancel()
	rental := prepared.rental
	rentalDirectory := filepath.Join(directory, rental.ID)
	label := prefix(configuration.Owner) + runID + ":" + rental.ID + ":"
	state := journal{Rental: rental, Label: label, Phase: "planned", StartedAt: time.Now().UTC(), InstallationFingerprint: prepared.startupFingerprint}
	if previous := prepared.retained; previous != nil {
		state.Label = previous.Label
		state.InstanceID = previous.InstanceID
		state.CreatedAt = previous.CreatedAt
		state.ReusedFrom = options.ReuseRun
		state.Phase = "reused"
		label = previous.Label
		options.logger.Printf("rental %s: reusing instance %d from retained run", rental.ID, state.InstanceID)
	}
	statePath := filepath.Join(rentalDirectory, "journal.json")
	if err := writeJSON(statePath, state); err != nil {
		return fmt.Errorf("record allocation intent: %w", err)
	}
	var reservationID int64
	createAttempted := false
	defer func() {
		released := !createAttempted
		cleanupContext, cleanupCancel := context.WithTimeout(context.Background(), duration(configuration.CleanupTimeoutSeconds))
		defer cleanupCancel()
		if state.InstanceID > 0 && !options.DisableCleanup {
			options.logger.Printf("rental %s: destroying instance %d", rental.ID, state.InstanceID)
			cleanupErr := destroy(cleanupContext, api, state.InstanceID, prefix(configuration.Owner))
			resultErr = errors.Join(resultErr, cleanupErr)
			if cleanupErr == nil {
				released = true
				state.ReleasedAt = time.Now().UTC()
				options.logger.Printf("rental %s: instance released", rental.ID)
			}
		}
		if reservationID > 0 {
			offersMu.Lock()
			if released {
				delete(usedOffers, reservationID)
			} else {
				usedOffers[reservationID] = false
			}
			offersMu.Unlock()
		}
		// Media inspection follows the configured destruction or retention policy.
		if resultErr == nil && options.Worker == "" {
			for _, job := range rental.Jobs {
				for _, seed := range job.Seeds {
					if err := finalizeMedia(options, job, filepath.Join(rentalDirectory, job.ID, strconv.FormatInt(seed, 10))); err != nil {
						resultErr = errors.Join(resultErr, err)
					}
				}
			}
		}
		state.CompletedAt = time.Now().UTC()
		state.Phase = "finished"
		if resultErr != nil {
			state.Failure = resultErr.Error()
			state.Phase = "failed"
		}
		if options.DisableCleanup && state.InstanceID > 0 {
			state.Phase = "retained"
			options.logger.Printf("rental %s: instance %d retained; automatic cleanup disabled", rental.ID, state.InstanceID)
		}
		if err := writeJSON(statePath, state); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("save rental outcome: %w", err))
		}
	}()
	if prepared.retained == nil {
		hardware := prepared.hardware
		query := prepared.query
		for key, value := range map[string]any{"gpu_ram": map[string]any{"gte": hardware.MinimumVRAM * 1024}, "disk_space": map[string]any{"gte": hardware.Disk}, "dph_total": map[string]any{"lte": hardware.MaximumPrice}, "allocated_storage": hardware.Disk, "direct_port_count": map[string]any{"gte": 2}, "limit": 16} {
			query[key] = value
		}
		if hardware.GPU != "" {
			query["gpu_name"] = map[string]any{"eq": hardware.GPU}
		}
		var selected vastclient.Offer
		for {
			offers, err := api.Search(ctx, query)
			if err != nil {
				return fmt.Errorf("search evaluation offers: %w", err)
			}
			sort.Slice(offers, func(i, j int) bool { return offers[i].HourlyPrice < offers[j].HourlyPrice })
			offersMu.Lock()
			for _, offer := range offers {
				_, reserved := usedOffers[offer.ID]
				if !reserved && offer.HostingType != nil && *offer.HostingType == 1 && offer.Verification != nil && *offer.Verification == "verified" && offer.GPUs == 1 && offer.GPURAM >= hardware.MinimumVRAM*1024 && offer.HourlyPrice > 0 && offer.HourlyPrice <= hardware.MaximumPrice && (hardware.GPU == "" || offer.GPUName == hardware.GPU) {
					selected = offer
					usedOffers[offer.ID] = true
					break
				}
			}
			releasable := false
			for _, active := range usedOffers {
				releasable = releasable || active
			}
			offersMu.Unlock()
			if selected.ID != 0 {
				reservationID = selected.ID
				break
			}
			if !releasable {
				return errors.New("no unused verified Secure Cloud offer satisfies the configured limits")
			}
			if err := pause(ctx, duration(configuration.PollSeconds)); err != nil {
				return fmt.Errorf("wait for capacity held by another evaluation worker: %w", err)
			}
		}
		if err := writeJSON(filepath.Join(rentalDirectory, "offer.json"), selected); err != nil {
			return fmt.Errorf("save selected offer: %w", err)
		}
		createAttempted = true
		id, err := api.Create(ctx, selected.ID, vastclient.CreateInput{Image: image, Disk: hardware.Disk, Label: label, Runtype: "ssh_direct", State: "running", DockerOptions: "-p 7860:7860", Onstart: prepared.startup})
		if err != nil {
			reconcileContext, reconcileCancel := context.WithTimeout(context.Background(), duration(configuration.CleanupTimeoutSeconds))
			defer reconcileCancel()
			instances, reconcileErr := api.Instances(reconcileContext)
			if reconcileErr == nil {
				for _, instance := range instances {
					if instance.Label == label {
						state.InstanceID = instance.ID
						break
					}
				}
			}
			return errors.Join(fmt.Errorf("create evaluation rental (no retry): %w", err), reconcileErr)
		}
		state.InstanceID = id
		state.Phase = "created"
		state.CreatedAt = time.Now().UTC()
		options.logger.Printf("rental %s: instance %d created; waiting for ComfyUI", rental.ID, id)
		if err := writeJSON(statePath, state); err != nil {
			return fmt.Errorf("record rental identity: %w", err)
		}
	}
	id := state.InstanceID
	bootContext, bootCancel := context.WithTimeout(ctx, duration(configuration.BootTimeoutSeconds))
	defer bootCancel()
	var endpoint string
	for {
		instance, err := api.Get(bootContext, id)
		if err != nil {
			return fmt.Errorf("observe evaluation rental: %w", err)
		}
		if instance.Label != label {
			return errors.New("rental label changed before workload readiness")
		}
		status := "unknown"
		if instance.Status != nil {
			switch *instance.Status {
			case "created", "loading", "running", "exited", "error", "offline", "stopped":
				status = *instance.Status
			}
		}
		if status != state.BootStatus {
			state.BootStatus = status
			if err := writeJSON(statePath, state); err != nil {
				return fmt.Errorf("record provider boot state: %w", err)
			}
			options.logger.Printf("rental %s: provider state %s; waiting for ComfyUI", rental.ID, status)
		}
		if status == "error" || status == "exited" || status == "stopped" {
			return errors.New("rental stopped before ComfyUI readiness")
		}
		port, err := vastclient.PublishedPort(instance.Ports["7860/tcp"])
		if err != nil {
			return fmt.Errorf("read rental HTTPS port: %w", err)
		}
		if instance.PublicIP != nil && port > 0 && instance.Status != nil && *instance.Status == "running" {
			host := *instance.PublicIP
			if net.ParseIP(host) == nil {
				return errors.New("rental returned an invalid public IP")
			}
			endpoint = (&url.URL{Scheme: "https", Host: net.JoinHostPort(host, strconv.FormatInt(port, 10))}).String()
			if options.Worker != "" {
				break
			}
			if _, err := callCLI(bootContext, options, endpoint, "--json", "system-stats"); err == nil {
				break
			}
		}
		if err := pause(bootContext, duration(configuration.PollSeconds)); err != nil {
			return fmt.Errorf("wait for ComfyUI readiness: %w", err)
		}
	}
	state.Phase = "ready"
	state.ReadyAt = time.Now().UTC()
	readinessStart := state.CreatedAt
	if state.ReusedFrom != "" {
		readinessStart = state.StartedAt
	}
	options.logger.Printf("rental %s: ready after %.1fs", rental.ID, state.ReadyAt.Sub(readinessStart).Seconds())
	if err := writeJSON(statePath, state); err != nil {
		return fmt.Errorf("record readiness: %w", err)
	}
	workContext, stopMonitoring := context.WithCancelCause(ctx)
	monitorDone := make(chan error, 1)
	go func() {
		monitorDone <- monitorRental(workContext, api, id, label, duration(configuration.PollSeconds), stopMonitoring)
	}()
	defer func() { stopMonitoring(nil); resultErr = errors.Join(resultErr, <-monitorDone) }()
	ctx = workContext
	if options.Worker != "" {
		for _, job := range rental.Jobs {
			for _, seed := range job.Seeds {
				options.logger.Printf("rental %s: running and retrieving %s seed %d", rental.ID, job.ID, seed)
				workContext, workCancel := context.WithTimeout(ctx, duration(configuration.WorkloadTimeoutSeconds))
				output := filepath.Join(rentalDirectory, job.ID, strconv.FormatInt(seed, 10))
				child := processCommand(workContext, options.Worker, options.WorkerArgs...)
				child.Env = append(workloadEnvironment(), "COMFY_HTTPS_URL="+endpoint, "EVALUATION_JOB_ID="+job.ID, "EVALUATION_OUTPUT="+output, "EVALUATION_SEED="+strconv.FormatInt(seed, 10))
				child.Stdout = io.Discard
				child.Stderr = io.Discard
				err := child.Run()
				workCancel()
				if err != nil {
					return fmt.Errorf("execute fixture workload %q seed %d: %w", job.ID, seed, err)
				}
				if err := confirmRetrieval(output); err != nil {
					return fmt.Errorf("confirm workload %q seed %d: %w", job.ID, seed, err)
				}
			}
		}
	} else if err := runQueuedWorkloads(ctx, options, rental.Jobs, endpoint, rentalDirectory, duration(configuration.WorkloadTimeoutSeconds), duration(configuration.PollSeconds)); err != nil {
		return fmt.Errorf("run pooled rental workloads: %w", err)
	}

	state.Phase = "downloaded"
	options.logger.Printf("rental %s: all results downloaded", rental.ID)
	if err := writeJSON(statePath, state); err != nil {
		return fmt.Errorf("record completed retrieval: %w", err)
	}
	return nil
}

// workloadEnvironment prevents account and Vault credentials reaching child workloads.
func workloadEnvironment() []string {
	var result []string
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name == "VAST_API_KEY" || name == "HF_TOKEN" || name == "HUGGING_FACE_HUB_TOKEN" || strings.HasPrefix(name, "VAULT_") || strings.HasPrefix(name, "AL_") || name == "TF_CLI_CONFIG_FILE" || strings.HasPrefix(name, "TF_VAR_") {
			continue
		}
		result = append(result, entry)
	}
	return result
}

// confirmRetrieval requires explicit download completion before rental destruction.
func confirmRetrieval(directory string) error {
	var receipt struct {
		// DownloadComplete is the workload's explicit retrieval acknowledgment.
		DownloadComplete bool `json:"download_complete"`
	}
	if err := readJSON(filepath.Join(directory, "result.json"), &receipt); err != nil {
		return fmt.Errorf("read workload retrieval receipt: %w", err)
	}
	if !receipt.DownloadComplete {
		return errors.New("workload returned without confirming complete retrieval")
	}
	return nil
}

// monitorRental cancels queued work when provider ownership or liveness is definitively lost.
func monitorRental(ctx context.Context, api *vastclient.Client, id int64, label string, poll time.Duration, cancel context.CancelCauseFunc) error {
	for {
		if err := pause(ctx, poll); err != nil {
			return nil
		}
		remote, err := api.Get(ctx, id)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if errors.Is(err, vastclient.ErrNotFound) {
				failure := fmt.Errorf("rental %d disappeared during workload: %w", id, err)
				cancel(failure)
				return failure
			}
			// Observation failures are transient; workload and rental deadlines remain authoritative.
			continue
		}
		var failure error
		if remote.Label != label {
			failure = fmt.Errorf("rental %d ownership changed during workload", id)
		} else if remote.Status != nil {
			switch *remote.Status {
			case "exited", "error", "stopped":
				failure = fmt.Errorf("rental %d %s during workload", id, *remote.Status)
			}
		}
		if failure != nil {
			cancel(failure)
			return failure
		}
	}
}

// compatibleInstallation compares the selected declared installation, not mutable prompts.
func compatibleInstallation(before, after registry, selected rental) bool {
	if before.Image != after.Image || !reflect.DeepEqual(before.Models[selected.ModelProfile], after.Models[selected.ModelProfile]) || !reflect.DeepEqual(before.Hardware[selected.HardwareProfile], after.Hardware[selected.HardwareProfile]) {
		return false
	}
	definition, found := before.Models[selected.ModelProfile]
	if !found {
		return false
	}
	for _, name := range definition.Downloads {
		if !reflect.DeepEqual(before.Downloads[name], after.Downloads[name]) {
			return false
		}
	}
	for _, name := range definition.Nodes {
		if !reflect.DeepEqual(before.Nodes[name], after.Nodes[name]) {
			return false
		}
	}
	return true
}

// reuseRentals validates every selected prior installation before any worker executes.
func reuseRentals(ctx context.Context, api *vastclient.Client, options Options, configuration configuration, profiles registry, rentals []preparedRental) error {
	ctx, cancel := context.WithTimeout(ctx, duration(configuration.BootTimeoutSeconds))
	defer cancel()
	var report struct {
		// Owner binds retained resources to the configured namespace.
		Owner string `json:"owner"`
		// ProviderEndpoint prevents identities being reused against another provider.
		ProviderEndpoint string `json:"provider_endpoint"`
		// CleanupSkipped establishes explicit retention rather than released resources.
		CleanupSkipped bool `json:"cleanup_skipped"`
		// Rentals identifies the prior hardware/model installations.
		Rentals []rental `json:"rentals"`
	}
	if err := readJSON(filepath.Join(options.ReuseRun, "run.json"), &report); err != nil {
		return fmt.Errorf("read retained run report: %w", err)
	}
	if report.Owner != configuration.Owner || report.ProviderEndpoint != options.Endpoint || !report.CleanupSkipped {
		return errors.New("source run must retain the same namespace and provider endpoint")
	}
	var inputs struct {
		// Profiles preserves the declared installed image, assets and hardware.
		Profiles registry `json:"profiles"`
	}
	if err := readJSON(filepath.Join(options.ReuseRun, "inputs.json"), &inputs); err != nil {
		return fmt.Errorf("read retained installation inputs: %w", err)
	}
	identities := map[int64]bool{}
	safe := regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)
	for index := range rentals {
		prepared := &rentals[index]
		var matches []rental
		for _, previous := range report.Rentals {
			if previous.ModelProfile == prepared.rental.ModelProfile && previous.HardwareProfile == prepared.rental.HardwareProfile {
				matches = append(matches, previous)
			}
		}
		if len(matches) != 1 || !safe.MatchString(matches[0].ID) {
			return fmt.Errorf("rental %q requires exactly one retained hardware/model match", prepared.rental.ID)
		}
		previous := matches[0]
		if !compatibleInstallation(inputs.Profiles, profiles, prepared.rental) {
			return fmt.Errorf("rental %q differs from its retained installation declaration", prepared.rental.ID)
		}
		var state journal
		if err := readJSON(filepath.Join(options.ReuseRun, previous.ID, "journal.json"), &state); err != nil {
			return fmt.Errorf("read retained rental %q: %w", previous.ID, err)
		}
		if state.InstallationFingerprint == "" || state.InstallationFingerprint != prepared.startupFingerprint {
			return fmt.Errorf("rental %q differs from its retained provisioning inputs or lacks an installation fingerprint", previous.ID)
		}
		parts := strings.Split(state.Label, ":")
		validLabel := len(parts) == 5 && parts[0] == "imagegen-eval" && parts[1] == configuration.Owner && len(parts[2]) == 32 && safe.MatchString(parts[3]) && parts[4] == ""
		if validLabel {
			_, err := hex.DecodeString(parts[2])
			validLabel = err == nil
		}
		if state.Rental.ID != previous.ID || state.Rental.ModelProfile != previous.ModelProfile || state.Rental.HardwareProfile != previous.HardwareProfile || state.Phase != "retained" || !state.ReleasedAt.IsZero() || state.CreatedAt.IsZero() || state.InstanceID <= 0 || identities[state.InstanceID] || !validLabel {
			return fmt.Errorf("rental %q has an invalid, duplicated or released retained identity", previous.ID)
		}
		identities[state.InstanceID] = true
		remote, err := api.Get(ctx, state.InstanceID)
		if err != nil {
			return fmt.Errorf("inspect retained instance %d: %w", state.InstanceID, err)
		}
		if remote.Label != state.Label || remote.Status == nil || *remote.Status != "running" {
			return fmt.Errorf("retained instance %d ownership or running state does not match", state.InstanceID)
		}
		prepared.retained = &state
	}
	return nil
}
