package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type deliverOptions struct {
	MessageFile string
	ReceiptFile string
	OwnerRoot   string
	TaskPaths   []string
}

type deliverStepReport struct {
	Step   string `json:"step"`
	Status string `json:"status"`
}

type deliverReport struct {
	Steps   []deliverStepReport `json:"steps"`
	Publish *publishReport      `json:"publish,omitempty"`
	Verify  *verifyReport       `json:"verify,omitempty"`
}

func newDeliverCommand(
	ctx context.Context,
	config *deliveryConfig,
	getenv func(string) string,
	stdout io.Writer,
	runner commandRunner,
) *cobra.Command {
	options := &deliverOptions{}
	command := &cobra.Command{
		Use:   "deliver",
		Short: "Prepare, publish, and verify in one idempotent command",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			delivery, err := deliveryFromConfig(ctx, config, getenv, runner)
			if err != nil {
				return err
			}
			report, err := delivery.deliver(ctx, *options)
			if report != nil {
				if outputErr := writeJSON(stdout, report); outputErr != nil {
					return outputErr
				}
			}
			return err
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
	_ = command.MarkFlagRequired("message-file")
	_ = command.MarkFlagRequired("receipt-file")
	_ = command.MarkFlagRequired("owner-root")
	return command
}

func (d *delivery) deliver(
	ctx context.Context,
	options deliverOptions,
) (*deliverReport, error) {
	if err := d.requireSameTaskOutputs(
		options.MessageFile,
		options.ReceiptFile,
	); err != nil {
		return nil, err
	}
	if len(options.TaskPaths) == 0 {
		return nil, fmt.Errorf("--task-path is required")
	}
	authorizationPath := buildAuthorizationPath(options.ReceiptFile)
	report := &deliverReport{}
	inspection, err := d.inspect(ctx)
	if err != nil {
		return report, fmt.Errorf("inspect: %w", err)
	}
	headOID := inspection.LocalHeadOID
	remoteHeadOID := ""
	if inspection.RemoteHeadOID != "" {
		remoteHeadOID = inspection.RemoteHeadOID
	}
	if inspection.RemoteHeadDiverged && remoteHeadOID != "" {
		authorization, err := d.buildRewriteAuthorization(ctx,
			rewriteAuthorizeOptions{
				OldRemoteOID:      remoteHeadOID,
				NewHeadOID:        headOID,
				OwnerRoot:         options.OwnerRoot,
				TaskPaths:         options.TaskPaths,
				AuthorizationFile: authorizationPath,
			})
		if err != nil {
			return report, fmt.Errorf("rewrite-authorize: %w", err)
		}
		if err := d.writeAtomicIgnoredJSON(ctx, authorizationPath,
			"rewrite authorization", authorization); err != nil {
			return report, fmt.Errorf("write rewrite authorization: %w", err)
		}
		report.Steps = append(report.Steps, deliverStepReport{
			Step:   "rewrite-authorize",
			Status: "ok",
		})
	}
	prepareOptions := prepareOptions{
		MessageFile: options.MessageFile,
		ReceiptFile: options.ReceiptFile,
		MessageOnly: true,
		RewriteOID:  headOID,
	}
	if inspection.RemoteHeadDiverged && remoteHeadOID != "" {
		prepareOptions.RewriteAuthorization = authorizationPath
	}
	_, err = d.prepare(ctx, prepareOptions)
	if err != nil {
		return report, fmt.Errorf("prepare: %w", err)
	}
	report.Steps = append(report.Steps, deliverStepReport{
		Step:   "prepare",
		Status: "ok",
	})
	receipt, err := d.readReceipt(ctx, options.ReceiptFile)
	if err != nil {
		return report, fmt.Errorf("read prepared receipt: %w", err)
	}
	publishReport, err := d.publish(ctx, publishOptions{
		ReceiptFile: options.ReceiptFile,
	})
	if err != nil {
		return report, fmt.Errorf("publish: %w", err)
	}
	report.Publish = publishReport
	report.Steps = append(report.Steps, deliverStepReport{
		Step:   "publish",
		Status: "ok",
	})
	verifyReport, err := d.verify(ctx, &receipt)
	if err != nil {
		return report, fmt.Errorf("verify: %w", err)
	}
	report.Verify = verifyReport
	report.Steps = append(report.Steps, deliverStepReport{
		Step:   "verify",
		Status: "ok",
	})
	return report, nil
}

func buildAuthorizationPath(receiptPath string) string {
	directory := filepath.Dir(receiptPath)
	base := filepath.Base(receiptPath)
	base = strings.TrimSuffix(base, ".json")
	return filepath.Join(directory, "."+base+".rewrite-authorization.json")
}
