package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/goal/internal/fsstore"
	"github.com/spf13/cobra"
)

const (
	defaultOutputLimit = 20
	toolVersion        = "0.0.1"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	root := newRootCommand(ctx, getenv, stdout)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

func newRootCommand(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) *cobra.Command {
	_ = ctx
	workspaceRoot := ""
	root := &cobra.Command{
		Use:           "goal",
		Short:         "Manage local versioned goal records",
		Version:       toolVersion,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.PersistentFlags().StringVar(
		&workspaceRoot,
		"workspace-root",
		"",
		"workspace root (default BUILD_WORKSPACE_DIRECTORY or current directory)",
	)
	storeForCommand := func() (*fsstore.Store, error) {
		value := workspaceRoot
		if value == "" {
			value = getenv("BUILD_WORKSPACE_DIRECTORY")
		}
		if value == "" {
			var err error
			value, err = os.Getwd()
			if err != nil {
				return nil, err
			}
		}
		return fsstore.NewStoreWithRuntimeDir(
			value,
			getenv("XDG_RUNTIME_DIR"),
		)
	}
	root.AddCommand(
		newInitCommand(storeForCommand, stdout),
		newListCommand(storeForCommand, stdout),
		newGraphCommand(storeForCommand, stdout),
		newShowCommand(storeForCommand, stdout),
		newAttachCommand(storeForCommand, stdout),
		newCheckpointCommand(storeForCommand, stdout),
		newSetRelationshipsCommand(storeForCommand, stdout),
		newValidateCommand(storeForCommand, stdout),
		newPromoteCommand(storeForCommand, stdout),
		newRenderCommand(storeForCommand, stdout),
		newMigrateCommand(storeForCommand, stdout),
	)
	return root
}

type storeFactory func() (*fsstore.Store, error)

func newInitCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	options := fsstore.InitOptions{}
	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize one goal record",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.GoalsRoot == "" || options.Title == "" {
				return fmt.Errorf("--goals-root and --title are required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.Init(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.GoalsRoot, "goals-root", "", "directory containing goal records")
	flags.StringVar(&options.Title, "title", "", "goal title")
	flags.StringVar(&options.GoalID, "goal-id", "", "portable goal ID (generated when omitted)")
	flags.StringVar(&options.Scope, "scope", "workspace", "goal scope: workspace or project")
	flags.StringVar(&options.OwnerRoot, "owner-root", "", "owning project/task root")
	flags.StringArrayVar(&options.Criteria, "criterion", nil, "acceptance criterion (repeatable)")
	flags.StringVar(&options.Retention, "retention", "", "retention policy: ephemeral or durable")
	return command
}

func newListCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	goalsRoot := ""
	limit := defaultOutputLimit
	command := &cobra.Command{
		Use:   "list",
		Short: "List a bounded goal catalog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if goalsRoot == "" {
				return fmt.Errorf("--goals-root is required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.List(goalsRoot, limit)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	command.Flags().StringVar(&goalsRoot, "goals-root", "", "directory containing goal records")
	command.Flags().IntVar(&limit, "limit", defaultOutputLimit, "maximum returned goals")
	return command
}

func newGraphCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	goalsRoot := ""
	command := &cobra.Command{
		Use:   "graph",
		Short: "Derive the relationship graph for a goal catalog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if goalsRoot == "" {
				return fmt.Errorf("--goals-root is required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			analysis, err := store.Graph(goalsRoot)
			if err != nil {
				return err
			}
			return writeJSON(stdout, analysis)
		},
	}
	command.Flags().StringVar(
		&goalsRoot,
		"goals-root",
		"",
		"directory containing goal records",
	)
	return command
}

func newShowCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	goalDir := ""
	sessionRoot := ""
	sessionID := ""
	limit := defaultOutputLimit
	command := &cobra.Command{
		Use:   "show",
		Short: "Show one bounded goal view",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (goalDir == "") == (sessionID == "") {
				return fmt.Errorf("set exactly one of --goal-dir or --session-id")
			}
			if sessionID != "" && sessionRoot == "" {
				return fmt.Errorf("--session-root is required with --session-id")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			var result fsstore.GoalView
			if goalDir != "" {
				result, err = store.ShowGoal(goalDir, limit)
			} else {
				result, err = store.ShowSession(sessionRoot, sessionID, limit)
			}
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&goalDir, "goal-dir", "", "goal record directory")
	flags.StringVar(&sessionRoot, "session-root", "", "task-specific session binding directory")
	flags.StringVar(&sessionID, "session-id", "", "session binding ID")
	flags.IntVar(&limit, "limit", defaultOutputLimit, "maximum returned criteria and attempts")
	return command
}

func newAttachCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	options := fsstore.AttachOptions{}
	command := &cobra.Command{
		Use:   "attach",
		Short: "Attach a session to one explicit goal",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.SessionRoot == "" || options.SessionID == "" || options.GoalDir == "" {
				return fmt.Errorf("--session-root, --session-id, and --goal-dir are required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.Attach(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.SessionRoot, "session-root", "", "task-specific session binding directory")
	flags.StringVar(&options.SessionID, "session-id", "", "session binding ID")
	flags.StringVar(&options.GoalDir, "goal-dir", "", "goal record directory")
	return command
}

func newCheckpointCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	options := fsstore.CheckpointOptions{}
	command := &cobra.Command{
		Use:   "checkpoint",
		Short: "Publish one revision-guarded goal checkpoint",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.GoalDir == "" || options.ExpectedResourceVersion == "" {
				return fmt.Errorf("--goal-dir and --expected-resource-version are required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.Checkpoint(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.GoalDir, "goal-dir", "", "goal record directory")
	flags.StringVar(&options.ExpectedResourceVersion, "expected-resource-version", "", "exact current Goal resourceVersion")
	flags.StringVar(&options.AttemptID, "attempt-id", "", "attempt ID (generated when publication starts without one)")
	flags.StringVar(&options.WorkType, "work-type", "", "attempt work type (default change for a new attempt)")
	flags.StringVar(&options.PlanFile, "plan-file", "", "immutable plan Markdown for a new attempt")
	flags.StringVar(&options.ResultFile, "result-file", "", "result Markdown to publish")
	flags.StringSliceVar(&options.EvidenceFiles, "evidence", nil, "evidence file to copy (repeatable)")
	flags.StringVar(&options.ReviewFile, "review-file", "", "structured close review YAML")
	flags.StringVar(&options.CriteriaFile, "criteria-file", "", "desired criteria items YAML; requires paused execution")
	flags.BoolVar(&options.CloseAttempt, "close-attempt", false, "close the selected attempt after publication")
	flags.StringVar(&options.Outcome, "outcome", "", "goal outcome transition")
	flags.StringVar(&options.Execution, "execution", "", "goal execution transition")
	return command
}

func newSetRelationshipsCommand(
	factory storeFactory,
	stdout io.Writer,
) *cobra.Command {
	options := fsstore.SetRelationshipsOptions{}
	command := &cobra.Command{
		Use:   "set-relationships",
		Short: "Atomically replace one goal's relationships",
		Long: "Replace the complete dependency and supersession lists and " +
			"optionally change the parent. The goal must have no active " +
			"attempt. Cycle prevention uses a per-goal-locked catalog " +
			"snapshot; after concurrent writes settle, run graph again.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.GoalDir == "" || options.ExpectedResourceVersion == "" {
				return fmt.Errorf(
					"--goal-dir and --expected-resource-version are required",
				)
			}
			parentSet := cmd.Flags().Changed("parent-goal")
			if parentSet && options.ClearParent {
				return fmt.Errorf(
					"--parent-goal and --clear-parent are mutually exclusive",
				)
			}
			if parentSet && options.ParentGoal == "" {
				return fmt.Errorf("--parent-goal requires a non-empty Goal name")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.SetRelationships(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.GoalDir, "goal-dir", "", "goal record directory")
	flags.StringVar(
		&options.ExpectedResourceVersion,
		"expected-resource-version",
		"",
		"exact current Goal resourceVersion",
	)
	flags.StringVar(
		&options.ParentGoal,
		"parent-goal",
		"",
		"parent Goal name; omitted to preserve the current parent",
	)
	flags.BoolVar(
		&options.ClearParent,
		"clear-parent",
		false,
		"remove the current parent Goal reference",
	)
	flags.StringSliceVar(
		&options.DependsOn,
		"depends-on",
		nil,
		"complete dependency Goal-name set (repeatable)",
	)
	flags.StringSliceVar(
		&options.Supersedes,
		"supersedes",
		nil,
		"complete superseded Goal-name set (repeatable)",
	)
	return command
}

func newValidateCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	goalDir := ""
	goalsRoot := ""
	command := &cobra.Command{
		Use:   "validate",
		Short: "Validate one goal or a goal catalog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (goalDir == "") == (goalsRoot == "") {
				return fmt.Errorf("set exactly one of --goal-dir or --goals-root")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			count := 1
			if goalDir != "" {
				err = store.ValidateGoal(goalDir)
			} else {
				count, err = store.ValidateRoot(goalsRoot)
			}
			if err != nil {
				return err
			}
			return writeJSON(stdout, map[string]any{
				"apiVersion": fsstore.APIVersion,
				"kind":       "GoalValidation",
				"valid":      true,
				"count":      count,
			})
		},
	}
	command.Flags().StringVar(&goalDir, "goal-dir", "", "goal record directory")
	command.Flags().StringVar(&goalsRoot, "goals-root", "", "directory containing goal records")
	return command
}

func newPromoteCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	options := fsstore.PromoteOptions{}
	command := &cobra.Command{
		Use:   "promote",
		Short: "Promote a paused workspace goal into project storage",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.GoalDir == "" || options.DestinationGoalsRoot == "" ||
				options.ExpectedResourceVersion == "" {
				return fmt.Errorf("--goal-dir, --destination-goals-root, and --expected-resource-version are required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.Promote(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(&options.GoalDir, "goal-dir", "", "source workspace goal directory")
	flags.StringVar(&options.DestinationGoalsRoot, "destination-goals-root", "", "project goals directory")
	flags.StringVar(&options.ExpectedResourceVersion, "expected-resource-version", "", "exact source Goal resourceVersion")
	flags.StringVar(&options.OwnerRoot, "owner-root", "", "destination owner root")
	return command
}

func newRenderCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	goalDir := ""
	expected := ""
	limit := defaultOutputLimit
	command := &cobra.Command{
		Use:   "render",
		Short: "Regenerate the bounded README projection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if goalDir == "" || expected == "" {
				return fmt.Errorf("--goal-dir and --expected-resource-version are required")
			}
			store, err := factory()
			if err != nil {
				return err
			}
			if err := store.Render(goalDir, expected, limit); err != nil {
				return err
			}
			return writeJSON(stdout, map[string]any{
				"apiVersion":      fsstore.APIVersion,
				"kind":            "GoalProjection",
				"resourceVersion": expected,
			})
		},
	}
	flags := command.Flags()
	flags.StringVar(&goalDir, "goal-dir", "", "goal record directory")
	flags.StringVar(&expected, "expected-resource-version", "", "exact current Goal resourceVersion")
	flags.IntVar(&limit, "limit", defaultOutputLimit, "maximum rendered criteria and attempts")
	return command
}

func newMigrateCommand(factory storeFactory, stdout io.Writer) *cobra.Command {
	options := fsstore.MigrateOptions{}
	command := &cobra.Command{
		Use:   "migrate",
		Short: "Import an unversioned goal into a fresh v1alpha1 record",
		Long: "Read an unversioned Markdown goal without changing it, build a " +
			"complete v1alpha1 record under the destination goals root, and " +
			"publish that record only when the target Goal name is absent.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if options.SourceGoalDir == "" || options.DestinationGoalsRoot == "" {
				return fmt.Errorf(
					"--source-goal-dir and --destination-goals-root are required",
				)
			}
			store, err := factory()
			if err != nil {
				return err
			}
			result, err := store.Migrate(options)
			if err != nil {
				return err
			}
			return writeJSON(stdout, result)
		},
	}
	flags := command.Flags()
	flags.StringVar(
		&options.SourceGoalDir,
		"source-goal-dir",
		"",
		"existing unversioned Markdown goal directory; never modified",
	)
	flags.StringVar(
		&options.DestinationGoalsRoot,
		"destination-goals-root",
		"",
		"destination goals directory for the new <goal-id> record",
	)
	flags.StringVar(&options.GoalID, "goal-id", "", "must match the existing directory name")
	flags.StringVar(&options.Title, "title", "", "title override for ambiguous prose")
	flags.StringVar(&options.Scope, "scope", "workspace", "goal scope: workspace or project")
	flags.StringVar(&options.OwnerRoot, "owner-root", "", "owning project/task root")
	flags.StringArrayVar(&options.Criteria, "criterion", nil, "acceptance criterion override (repeatable)")
	return command
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
