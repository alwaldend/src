package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// options identifies the packaged tools and a new caller-owned receipt.
type options struct {
	// terraform is the declared Terraform launcher with the rental configuration.
	terraform string
	// vastai is the pinned CLI using the same injected rental credential.
	vastai string
	// receipt is a new file recording non-sensitive lifecycle evidence.
	receipt string
}

// receipt records observed results without API responses, credentials, or personal data.
type receipt struct {
	// Success means creation, release, and independent cleanup verification passed.
	Success bool `json:"success"`
	// Team identifies the requested funded team.
	Team string `json:"team"`
	// StartedAt records when preflight began.
	StartedAt time.Time `json:"started_at"`
	// ReleaseStartedAt records when cleanup began after the creation command returned.
	ReleaseStartedAt time.Time `json:"release_started_at,omitempty"`
	// CompletedAt records when final verification ended.
	CompletedAt time.Time `json:"completed_at"`
	// ContractID identifies the test-owned rental acknowledged by Terraform.
	ContractID string `json:"contract_id,omitempty"`
	// HourlyPrice records the selected offer's total hourly price in USD.
	HourlyPrice float64 `json:"hourly_price_usd,omitempty"`
	// StateEmpty confirms the dedicated test state contains no resources.
	StateEmpty bool `json:"state_empty"`
	// ContractAbsent confirms the API no longer lists the test contract.
	ContractAbsent bool `json:"contract_absent"`
	// CleanupAttempted distinguishes preflight failure from a failed lifecycle.
	CleanupAttempted bool `json:"cleanup_attempted"`
}

// main adapts process cancellation and Cobra input to the lifecycle operation.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	var input options
	cmd := &cobra.Command{
		Use: "rental-test", Short: "Interactively rent the cheapest GPU and immediately release it", SilenceUsage: true, SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := input.run(cmd.Context()); err != nil {
				return fmt.Errorf("run rental lifecycle: %w", err)
			}
			return nil
		}, Args: cobra.NoArgs,
	}
	cmd.Flags().StringVar(&input.terraform, "terraform", "", "Packaged Terraform launcher")
	cmd.Flags().StringVar(&input.vastai, "vastai", "", "Pinned Vast.ai CLI")
	cmd.Flags().StringVar(&input.receipt, "receipt", "", "New absolute receipt path under task scratch")
	for _, name := range []string{"terraform", "vastai", "receipt"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// capture executes a read-only tool operation without exposing its response.
func capture(ctx context.Context, executable string, args ...string) ([]byte, error) {
	cmd := toolCommand(ctx, executable, args...)
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run read-only tool operation: %w", err)
	}
	return data, nil
}

// toolCommand gives the Terraform launcher time to interrupt and reap its child.
func toolCommand(ctx context.Context, executable string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Cancel = func() error { return cmd.Process.Signal(os.Interrupt) }
	// The launcher allows ten seconds for Terraform to stop before forcing exit.
	cmd.WaitDelay = 20 * time.Second
	return cmd
}

// stateResources treats Terraform's explicit absent-state diagnostic as an empty backend.
func (o options) stateResources(ctx context.Context) ([]byte, error) {
	data, err := capture(ctx, o.terraform, "state", "list")
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 && strings.Contains(string(exitError.Stderr), "No state file was found!") {
			return nil, nil
		}
		return nil, fmt.Errorf("list rental state resources: %w", err)
	}
	return data, nil
}

// interactive runs a Terraform operation with its original confirmation prompt.
func interactive(ctx context.Context, executable string, args ...string) ([]byte, error) {
	cmd := toolCommand(ctx, executable, args...)
	var output bytes.Buffer
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, io.MultiWriter(os.Stdout, &output), os.Stderr
	if err := cmd.Run(); err != nil {
		return output.Bytes(), fmt.Errorf("run interactive Terraform operation: %w", err)
	}
	return output.Bytes(), nil
}

// instanceIDs reads only identifiers from an otherwise private account response.
func (o options) instanceIDs(ctx context.Context) (map[int64]bool, error) {
	data, err := capture(ctx, o.vastai, "--raw", "show", "instances")
	if err != nil {
		return nil, fmt.Errorf("inspect rental contracts: %w", err)
	}
	var instances []struct {
		// ID identifies a currently listed contract.
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(data, &instances); err != nil {
		return nil, fmt.Errorf("decode contract identifiers: %w", err)
	}
	if instances == nil {
		return nil, errors.New("contract inventory omitted its collection")
	}
	ids := make(map[int64]bool, len(instances))
	for _, instance := range instances {
		if instance.ID <= 0 {
			return nil, errors.New("contract inventory contains an invalid identifier")
		}
		ids[instance.ID] = true
	}
	return ids, nil
}

// run performs preflight, one creation attempt, immediate release, and final verification.
func (o options) run(ctx context.Context) (resultErr error) {
	f, err := os.OpenFile(o.receipt, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create rental receipt: %w", err)
	}
	report := receipt{Team: "alwaldend", StartedAt: time.Now().UTC()}
	defer func() {
		report.CompletedAt = time.Now().UTC()
		if err := json.NewEncoder(f).Encode(report); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("write rental receipt: %w", err))
		}
		if err := f.Close(); err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("close rental receipt: %w", err))
		}
	}()
	data, err := capture(ctx, o.vastai, "--raw", "show", "user")
	if err != nil {
		return fmt.Errorf("verify rental account: %w", err)
	}
	var account struct {
		// ID identifies the account actually selected by the credential.
		ID int64 `json:"id"`
		// Name confirms the funded team's expected name.
		Name string `json:"username"`
		// IsTeam distinguishes the team from a personal account with a similar name.
		IsTeam bool `json:"is_team"`
		// Credit is the prepaid balance available to this team.
		Credit float64 `json:"credit"`
	}
	if err := json.Unmarshal(data, &account); err != nil {
		return fmt.Errorf("decode rental account metadata: %w", err)
	}
	if account.ID <= 0 || !account.IsTeam || account.Name != "alwaldend" || account.Credit <= 0 {
		return errors.New("rental credential must select the funded alwaldend team")
	}
	state, err := o.stateResources(ctx)
	if err != nil {
		return fmt.Errorf("inspect dedicated test state: %w", err)
	}
	if len(bytes.TrimSpace(state)) != 0 {
		return errors.New("test state already owns resources; inspect and release them explicitly before another test")
	}
	before, err := o.instanceIDs(ctx)
	if err != nil {
		return fmt.Errorf("inspect pre-test contract inventory: %w", err)
	}
	created, createErr := interactive(ctx, o.terraform, "apply", "-no-color", "-var=label=test_infra_vastai_rental")
	report.ReleaseStartedAt = time.Now().UTC()
	report.CleanupAttempted = true
	cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, releaseErr := interactive(cleanupCtx, o.terraform, "apply", "-destroy", "-no-color", "-var=release=true", "-var=label=test_infra_vastai_rental")
	resultErr = errors.Join(createErr, releaseErr)
	if match := regexp.MustCompile(`(?m)^rental_contract_id = "([0-9]+)"\s*$`).FindSubmatch(created); len(match) == 2 {
		report.ContractID = string(match[1])
	}
	if match := regexp.MustCompile(`(?m)^rental_hourly_price = ([0-9.eE+-]+)\s*$`).FindSubmatch(created); len(match) == 2 {
		report.HourlyPrice, err = strconv.ParseFloat(string(match[1]), 64)
		if err != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("decode selected hourly price: %w", err))
		}
	}
	state, err = o.stateResources(cleanupCtx)
	if err != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("verify empty test state: %w", err))
	} else {
		report.StateEmpty = len(bytes.TrimSpace(state)) == 0
	}
	after, err := o.instanceIDs(cleanupCtx)
	if err != nil {
		resultErr = errors.Join(resultErr, fmt.Errorf("verify released contract inventory: %w", err))
	} else {
		id, _ := strconv.ParseInt(report.ContractID, 10, 64)
		report.ContractAbsent = id > 0 && !after[id]
		for id := range before {
			if !after[id] {
				resultErr = errors.Join(resultErr, errors.New("an unrelated contract disappeared during the test; inspect concurrent activity"))
			}
		}
		for id := range after {
			if !before[id] {
				resultErr = errors.Join(resultErr, errors.New("a new contract remains after cleanup; inspect the team before retrying"))
			}
		}
	}
	if !report.StateEmpty || !report.ContractAbsent || strings.TrimSpace(report.ContractID) == "" {
		resultErr = errors.Join(resultErr, errors.New("rental lifecycle was not fully confirmed; inspect the receipt and owning tf.release target"))
	}
	report.Success = resultErr == nil
	return resultErr
}
