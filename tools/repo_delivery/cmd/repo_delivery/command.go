package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func Execute(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	runner commandRunner,
) error {
	root := newRootCommand(ctx, getenv, stdout, runner)
	root.SetArgs(args)
	root.SetIn(stdin)
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
	runner commandRunner,
) *cobra.Command {
	config := &deliveryConfig{}
	root := &cobra.Command{
		Use:           "repo_delivery",
		Short:         "Safely prepare and publish one-commit feature branches",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	flags := root.PersistentFlags()
	flags.StringVar(
		&config.Repository,
		"repository",
		"",
		"repository worktree (default BUILD_WORKSPACE_DIRECTORY or cwd)",
	)
	flags.StringVar(&config.Remote, "remote", "origin", "Git remote name")
	flags.StringVar(
		&config.Base,
		"base",
		"",
		"base branch when no open pull request exists",
	)
	flags.StringVar(
		&config.Forge,
		"forge",
		"auto",
		"forge adapter (auto or github)",
	)
	flags.StringVar(
		&config.ForgeCLI,
		"forge-cli",
		"",
		"single forge CLI executable path or name",
	)
	flags.StringVar(
		&config.GitCLI,
		"git-cli",
		"git",
		"single Git executable path or name",
	)

	root.AddCommand(
		newProviderCommand(ctx, config, getenv, stdout, runner),
		newInspectCommand(ctx, config, getenv, stdout, runner),
		newRewriteAuthorizeCommand(ctx, config, getenv, stdout, runner),
		newPrepareCommand(ctx, config, getenv, stdout, runner),
		newPublishCommand(ctx, config, getenv, stdout, runner),
		newRebaseCommand(ctx, config, getenv, stdout, runner),
		newVerifyCommand(ctx, config, getenv, stdout, runner),
		newReviewCommand(ctx, config, getenv, stdout, runner),
	)
	return root
}

func newRewriteAuthorizeCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	options := &rewriteAuthorizeOptions{}
	command := &cobra.Command{
		Use:   "rewrite-authorize",
		Short: "Write a typed non-authorizing remote-rewrite authorization",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			authorization, err := delivery.buildRewriteAuthorization(
				ctx,
				*options,
			)
			if err != nil {
				return err
			}
			if err := delivery.writeAtomicIgnoredJSON(
				ctx,
				options.AuthorizationFile,
				"rewrite authorization",
				authorization,
			); err != nil {
				return err
			}
			return writeJSON(stdout, authorization)
		},
	}
	flags := command.Flags()
	flags.StringVar(
		&options.OldRemoteOID,
		"old-remote-oid",
		"",
		"exact fetched remote feature OID that will be replaced",
	)
	flags.StringVar(
		&options.NewHeadOID,
		"new-head-oid",
		"",
		"exact local candidate head OID that will replace the remote ref",
	)
	flags.StringVar(
		&options.OwnerRoot,
		"owner-root",
		"",
		"task-owned path root that establishes ownership",
	)
	flags.StringArrayVar(
		&options.TaskPaths,
		"task-path",
		nil,
		"exact task-owned path (repeatable)",
	)
	flags.StringVar(
		&options.SourceReceipt,
		"source-receipt",
		"",
		"ignored preparation or inspection receipt establishing provenance",
	)
	flags.StringVar(
		&options.AuthorizationFile,
		"authorization-file",
		"",
		"ignored output authorization receipt under out/<task>/",
	)
	_ = command.MarkFlagRequired("old-remote-oid")
	_ = command.MarkFlagRequired("new-head-oid")
	_ = command.MarkFlagRequired("owner-root")
	_ = command.MarkFlagRequired("authorization-file")
	return command
}

func newProviderCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	return &cobra.Command{
		Use:   "provider",
		Short: "Detect the delivery provider without exposing remote URLs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := inspectProvider(ctx, *config, getenv, runner)
			if err != nil {
				return err
			}
			return writeJSON(stdout, report)
		},
	}
}

func inspectProvider(
	ctx context.Context,
	config deliveryConfig,
	getenv func(string) string,
	runner commandRunner,
) (*providerReport, error) {
	remote := strings.TrimSpace(config.Remote)
	if remote == "" {
		remote = "origin"
	}
	if err := validateRemoteName(remote); err != nil {
		return nil, err
	}
	repository, err := openGitRepository(
		ctx,
		config.Repository,
		config.GitCLI,
		getenv,
		runner,
	)
	if err != nil {
		return nil, err
	}
	fetchURL, pushURL, err := repository.remoteURLs(ctx, remote)
	if err != nil {
		return nil, err
	}
	fetchRepository, err := parseRemoteRepository(fetchURL)
	if err != nil {
		return nil, fmt.Errorf("parse fetch remote: %w", err)
	}
	pushRepository, err := parseRemoteRepository(pushURL)
	if err != nil {
		return nil, fmt.Errorf("parse push remote: %w", err)
	}
	if !sameRemoteRepository(fetchRepository, pushRepository) {
		return nil, fmt.Errorf(
			"fetch and push URLs identify different repositories",
		)
	}
	provider, available := providerForRepository(fetchRepository)
	selected := strings.ToLower(strings.TrimSpace(config.Forge))
	if selected != "" && selected != "auto" {
		if selected != "github" {
			return nil, fmt.Errorf("forge adapter %q is not supported", selected)
		}
		provider = selected
		available = true
	}
	gitTransport, deliveryTransportAvailable := providerGitTransport(
		fetchURL,
		pushURL,
	)
	return &providerReport{
		RemoteRepository:           fetchRepository,
		ProviderHint:               provider,
		AdapterAvailable:           available,
		GitTransport:               gitTransport,
		DeliveryTransportAvailable: deliveryTransportAvailable,
	}, nil
}

func providerGitTransport(fetchURL string, pushURL string) (string, bool) {
	fetchSSH := requireCanonicalSSHEndpoint(fetchURL) == nil
	pushSSH := requireCanonicalSSHEndpoint(pushURL) == nil
	if fetchSSH && pushSSH {
		return "ssh", true
	}
	fetchHTTPS := strings.HasPrefix(strings.ToLower(fetchURL), "https://")
	pushHTTPS := strings.HasPrefix(strings.ToLower(pushURL), "https://")
	if fetchHTTPS && pushHTTPS {
		return "https", false
	}
	return "mixed_or_unsupported", false
}

func deliveryFromConfig(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	runner commandRunner,
) (*delivery, error) {
	return newDelivery(ctx, *config, getenv, runner)
}

func newInspectCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Inspect exact Git and pull-request delivery state",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.inspect(ctx)
			if err != nil {
				return err
			}
			return writeJSON(stdout, report)
		},
	}
}

func newPrepareCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	options := &prepareOptions{}
	command := &cobra.Command{
		Use:   "prepare",
		Short: "Create, amend, or consolidate the feature commit and rebase it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.prepare(ctx, *options)
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(
		&options.MessageFile,
		"message-file",
		"",
		"aggregate commit message under out/<task>/",
	)
	flags.StringVar(
		&options.ReceiptFile,
		"receipt-file",
		"",
		"ignored out/<task>/ preparation receipt",
	)
	flags.StringArrayVar(
		&options.Paths,
		"path",
		nil,
		"explicit fully task-owned path (repeatable)",
	)
	flags.BoolVar(
		&options.UseIndex,
		"use-index",
		false,
		"commit the already prepared index without staging paths",
	)
	flags.BoolVar(
		&options.MessageOnly,
		"message-only",
		false,
		"amend only the aggregate message, preserving the tree",
	)
	flags.StringVar(
		&options.RewriteOID,
		"rewrite",
		"",
		"exact local commit OID authorized for amendment",
	)
	flags.StringVar(
		&options.ConsolidateOID,
		"consolidate",
		"",
		"exact local multi-commit head authorized for ownership consolidation",
	)
	flags.StringVar(
		&options.ReplaceRemoteOID,
		"replace-remote",
		"",
		"exact task-owned remote OID authorized for replacement",
	)
	flags.StringVar(
		&options.RewriteAuthorization,
		"rewrite-authorization",
		"",
		"typed rewrite-authorization receipt under out/<task>/ authorizing a remote replacement",
	)
	_ = command.MarkFlagRequired("message-file")
	_ = command.MarkFlagRequired("receipt-file")
	return command
}

func newPublishCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	options := &publishOptions{}
	command := &cobra.Command{
		Use:   "publish",
		Short: "Refetch, exact-lease push, synchronize the PR, and verify",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.publish(ctx, *options)
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(
		&options.ValidatedHead,
		"validated-head",
		"",
		"exact HEAD commit OID covered by the reported validations",
	)
	flags.StringVar(
		&options.ReceiptFile,
		"receipt-file",
		"",
		"ignored preparation receipt emitted by prepare",
	)
	flags.BoolVar(
		&options.NoPullRequest,
		"no-pull-request",
		false,
		"push the feature branch without creating or synchronizing a pull request",
	)
	_ = command.MarkFlagRequired("validated-head")
	_ = command.MarkFlagRequired("receipt-file")
	return command
}

func newVerifyCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	options := &struct {
		ReceiptFile string
	}{}
	command := &cobra.Command{
		Use:   "verify",
		Short: "Fetch and verify final Git and pull-request invariants",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			receipt, err := delivery.readReceipt(ctx, options.ReceiptFile)
			if err != nil {
				return err
			}
			report, err := delivery.verify(ctx, &receipt)
			if err != nil {
				return err
			}
			return writeJSON(stdout, report)
		},
	}
	command.Flags().StringVar(
		&options.ReceiptFile,
		"receipt-file",
		"",
		"ignored preparation receipt emitted by prepare",
	)
	_ = command.MarkFlagRequired("receipt-file")
	return command
}

func newRebaseCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	return &cobra.Command{
		Use:   "rebase",
		Short: "Fetch the advanced base, replay the task-owned commit, and lease-push it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.rebase(ctx, rebaseOptions{})
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			return err
		},
	}
}
