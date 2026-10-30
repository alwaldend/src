package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/web"
	"github.com/spf13/cobra"
)

// main establishes cancellation and reports application errors at the CLI boundary.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// command constructs the plan, apply, and serve commands without global state.
func command() *cobra.Command {
	var base string
	var policyTimeout time.Duration
	root := &cobra.Command{Use: "x_profile_cleaner", Short: "Plan and apply Lua-configured X profile cleanup", SilenceUsage: true, SilenceErrors: true}
	root.PersistentFlags().StringVar(&base, "api-base-url", "https://api.x.com/2", "X API v2 base URL (HTTP allowed only on loopback)")
	root.PersistentFlags().DurationVar(&policyTimeout, "policy-timeout", 10*time.Second, "deadline for Lua configuration, callbacks, and lazy reads")
	application := &app.App{}
	logs := &web.Log{}
	logger := slog.New(slog.NewTextHandler(io.MultiWriter(logs, root.ErrOrStderr()), nil))
	root.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		credentials := oauth1.Credentials{
			APIKey:            os.Getenv("X_API_KEY"),
			APISecret:         os.Getenv("X_API_SECRET"),
			AccessToken:       os.Getenv("X_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("X_ACCESS_TOKEN_SECRET"),
		}
		if err := application.Configure(base, credentials, policyTimeout, logger); err != nil {
			return fmt.Errorf("configure cleanup application: %w", err)
		}
		return nil
	}
	root.AddCommand(newPlanCommand(application), newApplyCommand(application), newServeCommand(application, logger, logs))
	return root
}

// newServeCommand supplies an optional initial policy to the loopback HTTP adapter.
func newServeCommand(application *app.App, logger *slog.Logger, logs *web.Log) *cobra.Command {
	var address, path string
	command := &cobra.Command{Use: "serve", Short: "Review and apply cleanup in a local browser", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		source, err := web.DefaultPolicy()
		if path != "" {
			source, err = app.ReadPolicy(path)
		}
		if err != nil {
			return fmt.Errorf("load browser policy: %w", err)
		}
		if err := web.Serve(cmd.Context(), application, web.Config{Address: address, Policy: source, Output: cmd.OutOrStdout(), Logger: logger, Logs: logs}); err != nil {
			return fmt.Errorf("serve cleanup interface: %w", err)
		}
		return nil
	}}
	command.Flags().StringVar(&address, "listen", "127.0.0.1:0", "loopback listen address; port 0 selects an available port")
	command.Flags().StringVar(&path, "config", "", "initial Lua policy file; defaults to the maintained example")
	return command
}

// newPlanCommand owns plan input flags and output reporting.
func newPlanCommand(application *app.App) *cobra.Command {
	var config, output string
	var explain bool
	planCommand := &cobra.Command{Use: "plan", Short: "Prepare a deletion list without mutations", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		plan, err := application.PreparePlan(cmd.Context(), config, "")
		if err != nil {
			return fmt.Errorf("plan cleanup: %w", err)
		}
		if output == "-" {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(plan); err != nil {
				return fmt.Errorf("write plan to stdout: %w", err)
			}
		} else {
			if err := app.WritePlan(output, plan); err != nil {
				return fmt.Errorf("save cleanup plan: %w", err)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Planned %d actions for account %s in %s\n", len(plan.Actions), plan.AccountID, output); err != nil {
				return fmt.Errorf("report saved plan: %w", err)
			}
		}
		if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "Coverage: available timeline inventory only; complete historical discovery is outside scope."); err != nil {
			return fmt.Errorf("report discovery coverage: %w", err)
		}
		if explain {
			if err := displayDecisions(cmd.ErrOrStderr(), plan); err != nil {
				return fmt.Errorf("explain plan: %w", err)
			}
		}
		return nil
	}}
	planCommand.Flags().StringVar(&config, "config", "", "Lua configuration path")
	planCommand.Flags().StringVar(&output, "out", "-", "plan JSON output path, or - for stdout")
	planCommand.Flags().BoolVar(&explain, "explain", false, "print informational policy decisions on stderr")
	if err := planCommand.MarkFlagRequired("config"); err != nil {
		panic(fmt.Sprintf("declare config flag: %v", err))
	}
	return planCommand
}

// newApplyCommand adapts execution flags, terminal confirmation, and progress.
func newApplyCommand(application *app.App) *cobra.Command {
	var applyConfig, savedPlan string
	var yes, explain bool
	applyCommand := &cobra.Command{Use: "apply", Short: "Execute a saved plan or prepare and confirm one", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		plan, err := application.PreparePlan(cmd.Context(), applyConfig, savedPlan)
		if err != nil {
			return fmt.Errorf("prepare cleanup plan: %w", err)
		}
		if err := displayPlan(cmd.OutOrStdout(), plan); err != nil {
			return fmt.Errorf("display cleanup plan: %w", err)
		}
		if explain {
			if err := displayDecisions(cmd.ErrOrStderr(), plan); err != nil {
				return fmt.Errorf("explain plan: %w", err)
			}
		}
		if savedPlan == "" && !yes {
			confirmed, err := confirm(cmd.Context(), cmd.InOrStdin(), cmd.ErrOrStderr())
			if err != nil {
				return fmt.Errorf("confirm cleanup: %w", err)
			}
			if !confirmed {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Not applied."); err != nil {
					return fmt.Errorf("report declined plan: %w", err)
				}
				return nil
			}
		}
		completed, applyErr := application.Apply(cmd.Context(), plan, nil)
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Completed %d of %d actions.\n", completed, len(plan.Actions)); err != nil {
			if applyErr != nil {
				return fmt.Errorf("apply cleanup: %w; report progress: %w", applyErr, err)
			}
			return fmt.Errorf("report apply progress: %w", err)
		}
		if applyErr != nil {
			return fmt.Errorf("apply cleanup: %w", applyErr)
		}
		return nil
	}}
	applyCommand.Flags().StringVar(&applyConfig, "config", "", "Lua configuration for an inline plan")
	applyCommand.Flags().StringVar(&savedPlan, "plan", "", "saved plan JSON; executes without reevaluating Lua")
	applyCommand.Flags().BoolVarP(&yes, "yes", "y", false, "apply an inline plan without confirmation")
	applyCommand.Flags().BoolVar(&explain, "explain", false, "print informational policy decisions on stderr")
	return applyCommand
}

// confirm waits for an affirmative line or cancellation without closing borrowed input.
// The buffered reader result can finish after cancellation; the CLI then exits.
func confirm(ctx context.Context, input io.Reader, output io.Writer) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("before confirmation: %w", err)
	}
	if _, err := fmt.Fprint(output, "Apply this plan? [y/N] "); err != nil {
		return false, fmt.Errorf("write confirmation prompt: %w", err)
	}
	// response carries a single line and any terminal input error.
	type response struct {
		// line is the user's answer, including its newline when present.
		line string
		// err distinguishes EOF from a complete answer or a failed read.
		err error
	}
	answers := make(chan response, 1)
	go func() {
		line, err := bufio.NewReader(input).ReadString('\n')
		answers <- response{line: line, err: err}
	}()
	select {
	case <-ctx.Done():
		return false, fmt.Errorf("wait for confirmation: %w", ctx.Err())
	case answer := <-answers:
		if err := ctx.Err(); err != nil {
			return false, fmt.Errorf("finish confirmation: %w", err)
		}
		if answer.err != nil && answer.err != io.EOF {
			return false, fmt.Errorf("read confirmation: %w", answer.err)
		}
		line := strings.ToLower(strings.TrimSpace(answer.line))
		return answer.err != io.EOF && (line == "y" || line == "yes"), nil
	}
}

// displayPlan prints the exact actions which confirmation will authorize.
func displayPlan(output io.Writer, plan model.Plan) error {
	if _, err := fmt.Fprintf(output, "Account %s: %d actions\n", plan.AccountID, len(plan.Actions)); err != nil {
		return fmt.Errorf("write plan summary: %w", err)
	}
	for _, action := range plan.Actions {
		level := "ordinary"
		if action.Forced {
			level = "forced"
		}
		if _, err := fmt.Fprintf(output, "%s %s (%s)\n", action.Kind, action.TargetID, level); err != nil {
			return fmt.Errorf("write plan action: %w", err)
		}
	}
	return nil
}

// displayDecisions prints quoted snapshot evidence; only displayed actions authorize execution.
func displayDecisions(output io.Writer, plan model.Plan) error {
	for _, decision := range plan.Decisions {
		if _, err := fmt.Fprintf(output, "%s %s %s %s %q\n", decision.Kind, decision.ItemID, decision.CreatedAt.UTC().Format(time.RFC3339), decision.Outcome, decision.Text); err != nil {
			return fmt.Errorf("write decision: %w", err)
		}
		for _, reason := range decision.Reasons {
			if _, err := fmt.Fprintf(output, "  rule %d %s %q seed=%s\n", reason.RuleID, reason.Operation, reason.Name, reason.SeedID); err != nil {
				return fmt.Errorf("write decision reason: %w", err)
			}
		}
	}
	return nil
}
