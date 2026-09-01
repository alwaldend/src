package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
	"git.alwaldend.com/alwaldend/src/tools/agents/evidence"
)

func execute(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	root := &cobra.Command{
		Use:          "evidence_tool",
		Short:        "emit immutable validation evidence",
		SilenceUsage: true,
	}
	root.AddCommand(newEmitCommand())
	root.AddCommand(newAssertCommand())
	root.SetOut(stdout)
	root.SetErr(stderr)
	return root.ExecuteContext(ctx)
}

type emitOptions struct {
	profile        string
	candidate      string
	baseOID        string
	treeOID        string
	commitOID      string
	workingScope   string
	result         []string
	provider       []string
	configDigest   []string
	toolchain      []string
	policyDigest   []string
	outputCalls    int64
	outputBytes    int64
	outputDuration int64
	rawLog         []string
	cleanPre       bool
	cleanPost      bool
}

func newEmitCommand() *cobra.Command {
	var opts emitOptions
	emit := &cobra.Command{
		Use:   "emit",
		Short: "emit a validation set as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			set, err := runEmit(cmd, opts)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), set)
		},
	}
	flags := emit.Flags()
	flags.StringVar(&opts.profile, "profile", "", "impact profile")
	flags.StringVar(&opts.candidate, "candidate", "", "candidate reference")
	flags.StringVar(&opts.baseOID, "base-oid", "", "base OID")
	flags.StringVar(&opts.treeOID, "tree-oid", "", "tree OID")
	flags.StringVar(&opts.commitOID, "commit-oid", "", "commit OID")
	flags.StringVar(&opts.workingScope, "working-scope", "", "working scope")
	flags.StringSliceVar(&opts.result, "result", nil, "checkId:status")
	flags.StringSliceVar(&opts.provider, "provider", nil, "provider ref")
	flags.StringSliceVar(&opts.configDigest, "config-digest", nil, "name=digest")
	flags.StringSliceVar(&opts.toolchain, "toolchain-digest", nil, "name=digest")
	flags.StringSliceVar(&opts.policyDigest, "policy-digest", nil, "name=digest")
	flags.Int64Var(&opts.outputCalls, "output-calls", 0, "output call limit")
	flags.Int64Var(&opts.outputBytes, "output-bytes", 0, "output byte limit")
	flags.Int64Var(
		&opts.outputDuration,
		"output-duration-ms",
		0,
		"output duration limit in ms",
	)
	flags.StringSliceVar(&opts.rawLog, "raw-log", nil, "raw log artifact ref")
	flags.BoolVar(&opts.cleanPre, "clean-pre", false, "pre-state was clean")
	flags.BoolVar(&opts.cleanPost, "clean-post", false, "post-state was clean")
	return emit
}

func runEmit(cmd *cobra.Command, opts emitOptions) (v1alpha1.ValidationSet, error) {
	profile, err := parseImpactProfile(opts.profile)
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	candidate, err := parseReference(opts.candidate, "candidate")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	if candidate.Kind != v1alpha1.ReferenceSubject &&
		candidate.Kind != v1alpha1.ReferenceRepository {
		return v1alpha1.ValidationSet{}, fmt.Errorf(
			"candidate must be a subject or repository reference",
		)
	}
	providers, err := parseReferences(opts.provider, "provider")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	rawLogs, err := parseReferences(opts.rawLog, "raw-log")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	configDigests, err := parseDigests(opts.configDigest, "config-digest")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	toolchainDigests, err := parseDigests(opts.toolchain, "toolchain-digest")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	policyDigests, err := parseDigests(opts.policyDigest, "policy-digest")
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	results, err := parseResults(opts.result)
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	if len(results) == 0 {
		return v1alpha1.ValidationSet{}, fmt.Errorf(
			"at least one --result is required",
		)
	}
	set, err := evidence.EmitValidationSet(evidence.EmitOptions{
		Profile:          profile,
		Candidate:        candidate,
		BaseOID:          opts.baseOID,
		TreeOID:          opts.treeOID,
		CommitOID:        opts.commitOID,
		SanitizedArgs:    cmd.Flags().Args(),
		WorkingScope:     opts.workingScope,
		ProviderRefs:     providers,
		ConfigDigests:    configDigests,
		ToolchainDigests: toolchainDigests,
		PolicyDigests:    policyDigests,
		OutputBounds: v1alpha1.Budget{
			Calls:      opts.outputCalls,
			Bytes:      opts.outputBytes,
			DurationMS: opts.outputDuration,
		},
		RawLogs:         rawLogs,
		TotalDurationMS: totalDuration(results),
		CleanPreState:   opts.cleanPre,
		CleanPostState:  opts.cleanPost,
	}, results)
	if err != nil {
		return v1alpha1.ValidationSet{}, err
	}
	return set, nil
}

func totalDuration(results []v1alpha1.ValidationResult) int64 {
	var total int64
	for _, result := range results {
		total += result.DurationMS
	}
	return total
}

type assertOptions struct {
	id            string
	criterion     string
	rev           string
	verdict       string
	validation    []string
	validationSet string
}

func newAssertCommand() *cobra.Command {
	var opts assertOptions
	assert := &cobra.Command{
		Use:   "assert",
		Short: "emit an evidence assertion as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			assertion, err := runAssert(cmd, opts)
			if err != nil {
				return err
			}
			return writeJSON(cmd.OutOrStdout(), assertion)
		},
	}
	flags := assert.Flags()
	flags.StringVar(&opts.id, "id", "", "assertion id")
	flags.StringVar(&opts.criterion, "criterion", "", "criterion reference")
	flags.StringVar(&opts.rev, "rev", "", "criterion revision")
	flags.StringVar(&opts.verdict, "verdict", "", "satisfied|unsatisfied|blocked|unknown")
	flags.StringSliceVar(
		&opts.validation,
		"validation",
		nil,
		"validation set JSON file",
	)
	flags.StringVar(
		&opts.validationSet,
		"validation-set",
		"",
		"single validation set JSON file",
	)
	return assert
}

func runAssert(
	cmd *cobra.Command,
	opts assertOptions,
) (v1alpha1.EvidenceAssertion, error) {
	criterion, err := parseReference(opts.criterion, "criterion")
	if err != nil {
		return v1alpha1.EvidenceAssertion{}, err
	}
	if criterion.Kind != v1alpha1.ReferenceGoal &&
		criterion.Kind != v1alpha1.ReferenceArtifact {
		return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
			"criterion must be a goal or artifact reference",
		)
	}
	if opts.id == "" {
		return v1alpha1.EvidenceAssertion{}, fmt.Errorf("--id is required")
	}
	files := append([]string(nil), opts.validation...)
	if opts.validationSet != "" {
		files = append(files, opts.validationSet)
	}
	if len(files) == 0 {
		return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
			"at least one --validation or --validation-set is required",
		)
	}
	sets := make([]v1alpha1.ValidationSet, 0, len(files))
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
				"could not read %s: %w",
				path,
				err,
			)
		}
		set, err := v1alpha1.DecodeValidationSet(content)
		if err != nil {
			return v1alpha1.EvidenceAssertion{}, fmt.Errorf(
				"invalid validation set %s: %w",
				path,
				err,
			)
		}
		sets = append(sets, set)
	}
	return evidence.NewAssertion(
		opts.id,
		criterion,
		opts.rev,
		opts.verdict,
		sets...,
	)
}

func parseImpactProfile(value string) (v1alpha1.ImpactProfile, error) {
	if value == "" {
		return "", fmt.Errorf("--profile is required")
	}
	switch v1alpha1.ImpactProfile(value) {
	case v1alpha1.ImpactProfileChangedFast,
		v1alpha1.ImpactProfileWorkspace,
		v1alpha1.ImpactProfileFreshEvidence,
		v1alpha1.ImpactProfileFullAudit,
		v1alpha1.ImpactProfileDiagnose:
		return v1alpha1.ImpactProfile(value), nil
	}
	return "", fmt.Errorf("unknown impact profile %q", value)
}

func parseReference(value, flagName string) (v1alpha1.Reference, error) {
	kind, id, ok := strings.Cut(value, ":")
	if !ok {
		return v1alpha1.Reference{}, fmt.Errorf(
			"--%s must be kind:id, got %q",
			flagName,
			value,
		)
	}
	reference := v1alpha1.Reference{Kind: v1alpha1.ReferenceKind(kind), ID: id}
	if err := reference.Validate(); err != nil {
		return v1alpha1.Reference{}, fmt.Errorf("--%s: %w", flagName, err)
	}
	return reference, nil
}

func parseReferences(
	values []string,
	flagName string,
) ([]v1alpha1.Reference, error) {
	references := make([]v1alpha1.Reference, 0, len(values))
	for _, value := range values {
		reference, err := parseReference(value, flagName)
		if err != nil {
			return nil, err
		}
		references = append(references, reference)
	}
	return references, nil
}

func parseDigests(
	values []string,
	flagName string,
) (map[string]string, error) {
	digests := make(map[string]string, len(values))
	for _, value := range values {
		name, digest, ok := strings.Cut(value, "=")
		if !ok {
			return nil, fmt.Errorf("--%s must be name=digest, got %q", flagName, value)
		}
		digests[name] = digest
	}
	return digests, nil
}

func parseResults(values []string) ([]v1alpha1.ValidationResult, error) {
	results := make([]v1alpha1.ValidationResult, 0, len(values))
	for _, value := range values {
		checkID, status, ok := strings.Cut(value, ":")
		if !ok {
			return nil, fmt.Errorf("--result must be checkId:status, got %q", value)
		}
		results = append(results, v1alpha1.ValidationResult{
			CheckID: checkID,
			Status:  status,
		})
	}
	return results, nil
}

func writeJSON(writer io.Writer, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "%s\n", content)
	return err
}
