// Command goal_check compiles the bounded deterministic GoalCatalog over
// the registered repository goals root (projects/agents/goals).
//
// It performs no network or stateful operation, and never mutates goal
// records. Outputs are checked repository artifacts: portable JSON plus a
// human Markdown render that states the JSON digest.
package goal_check

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	goalv1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
	goalcatalog "git.alwaldend.com/alwaldend/src/tools/agents/catalog/v1alpha1"
	"github.com/goccy/go-yaml"
)

const (
	goalsRegistryAuthorityID   = "repository.goals"
	goalsRegistryAuthorityKind = "goals"
	goalsRegistryAuthorityRoot = "projects/agents/goals"
	currentCriteriaFileName    = "criteria.yaml"
	criteriaRevisionsDirectory = "criteria-revisions"
	attemptsDirectory          = "attempts"
	attemptManifestFileName    = "attempt.yaml"
	planFileName               = "plan.md"
	resultFileName             = "result.md"
	evidenceDirectory          = "evidence"
)

const goalsRegistrySchema = "agents.alwaldend.com/phase1-registry/v1alpha1"

type options struct {
	workspaceRoot  string
	registryPath   string
	outputPath     string
	markdownPath   string
	check          bool
	sourceRevision string
	producerRef    string
}

func parseFlags(args []string) (options, error) {
	var opts options
	flags := flag.NewFlagSet("goal_check", flag.ContinueOnError)
	flags.StringVar(&opts.workspaceRoot, "workspace-root", "",
		"repository workspace root (required)")
	flags.StringVar(&opts.registryPath, "registry",
		"tools/agents/declarations/registry.json",
		"registry JSON path relative to workspace root")
	flags.StringVar(&opts.outputPath, "output", "tools/agents/catalogs/goal.json",
		"JSON output path relative to workspace root")
	flags.StringVar(&opts.markdownPath, "markdown", "tools/agents/catalogs/goal.md",
		"Markdown output path relative to workspace root")
	flags.BoolVar(&opts.check, "check", false,
		"validate and emit, then exit nonzero on completeness failure")
	flags.StringVar(&opts.sourceRevision, "source-revision", "",
		"exact Git tree/commit identity (default: content-addressed inputs)")
	flags.StringVar(&opts.producerRef, "producer-ref", "repository.goal-compiler",
		"producer reference for the catalog")
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	if opts.workspaceRoot == "" {
		return options{}, fmt.Errorf("--workspace-root is required")
	}
	return opts, nil
}

type registry struct {
	Schema      string `json:"schema"`
	Authorities []struct {
		ID     string `json:"id"`
		Kind   string `json:"kind"`
		Source string `json:"source"`
	} `json:"authorities"`
}

type compiler struct {
	root        string
	opts        options
	inputs      []goalcatalog.CatalogInput
	goals       []goalcatalog.GoalRecord
	problems    []string
	eligible    int
	emitted     int
	unavailable int
}

func (c *compiler) invalidRecord(format string, args ...any) string {
	reason := fmt.Sprintf(format, args...)
	c.problems = append(c.problems, reason)
	return reason
}

func (c *compiler) problem(format string, args ...any) {
	c.problems = append(c.problems, fmt.Sprintf(format, args...))
}

func (c *compiler) input(relativePath, role string, content []byte) {
	value := sha256.Sum256(content)
	c.inputs = append(c.inputs, goalcatalog.CatalogInput{
		Path:   filepath.ToSlash(relativePath),
		Role:   role,
		Digest: "sha256:" + hex.EncodeToString(value[:]),
	})
}

// compile discovers the registered goals root and compiles every eligible
// non-hidden goal directory into an availability record.
func (c *compiler) compile() error {
	goalsRoot, err := c.goalsRoot()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(goalsRoot)
	if err != nil {
		return fmt.Errorf("read goals root: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		goalDir := filepath.ToSlash(filepath.Join(
			goalsRegistryAuthorityRoot,
			entry.Name(),
		))
		c.eligible++
		record, reason := c.compileGoal(goalDir)
		c.goals = append(c.goals, record)
		if record.Availability == "available" {
			c.emitted++
		} else {
			c.unavailable++
		}
		if reason != "" && record.Reason == reason {
			c.problem("goal %s unavailable: %s", goalDir, reason)
		}
	}
	return nil
}

// goalsRoot reads the registry, verifies the repository.goals authority, and
// returns the absolute eligible goals root.
func (c *compiler) goalsRoot() (string, error) {
	registryPath := filepath.Join(c.root, filepath.FromSlash(c.opts.registryPath))
	registryContent, err := os.ReadFile(registryPath)
	if err != nil {
		return "", fmt.Errorf("read registry: %w", err)
	}
	var registry registry
	decoder := json.NewDecoder(bytes.NewReader(registryContent))
	if err := decoder.Decode(&registry); err != nil {
		return "", fmt.Errorf("decode registry: %w", err)
	}
	if registry.Schema != goalsRegistrySchema {
		return "", fmt.Errorf("registry schema mismatch: %s", registry.Schema)
	}
	c.input(c.opts.registryPath, "registry", registryContent)
	for _, authority := range registry.Authorities {
		if authority.ID == goalsRegistryAuthorityID &&
			authority.Kind == goalsRegistryAuthorityKind {
			if authority.Source != goalsRegistryAuthorityRoot {
				return "", fmt.Errorf(
					"goals authority root mismatch: %s",
					authority.Source,
				)
			}
			return filepath.Join(c.root, filepath.FromSlash(authority.Source)), nil
		}
	}
	return "", fmt.Errorf("registry lacks goals authority %q",
		goalsRegistryAuthorityID)
}

// compileGoal builds exactly one availability record for an eligible goal
// directory. It never mutates the record: availability requires the complete
// record to parse and validate purely from disk.
func (c *compiler) compileGoal(goalDir string) (goalcatalog.GoalRecord, string) {
	candidatePath := filepath.ToSlash(goalDir)
	unavailable := goalcatalog.GoalRecord{
		CandidatePath: candidatePath,
		Availability:  "unavailable",
	}
	goal, criteria, criteriaRevision, attempts, reason := c.loadRecord(goalDir)
	if reason != "" {
		unavailable.Reason = reason
		return unavailable, reason
	}
	if err := validateCompleteRecord(
		goal, criteria, criteriaRevision, attempts,
	); err != nil {
		unavailable.Reason = c.invalidRecord(
			"goal %s invalid record: %s",
			candidatePath,
			boundedReason(err),
		)
		return unavailable, unavailable.Reason
	}
	ownerRoot := goal.Metadata.Annotations[goalv1alpha1.LocalOwnerRootAnnotation]
	if ownerRoot == "" {
		ownerRoot = "projects/agents"
	}
	record := goalcatalog.GoalRecord{
		CandidatePath: candidatePath,
		Availability:  "available",
		Identity: &goalcatalog.GoalCoarseIdentity{
			Name:      goal.Metadata.Name,
			OwnerRoot: ownerRoot,
			Scope:     goal.Spec.Scope,
		},
		CoarseStatus: &goalcatalog.GoalCoarseStatus{
			Outcome:   goal.Status.Outcome,
			Execution: goal.Status.Execution,
		},
	}
	if goal.Status.Outcome == "open" {
		record.Continuation = continuationFor(goal, attempts)
	}
	return record, ""
}

// continuationFor builds the bounded resume projection for an open goal. It
// prefers the active attempt and falls back to the most recent open attempt;
// closed goals carry no continuation.
func continuationFor(
	goal goalv1alpha1.Goal,
	attempts []goalv1alpha1.GoalAttempt,
) *goalcatalog.GoalContinuation {
	byID := make(map[string]goalv1alpha1.GoalAttempt, len(attempts))
	for _, attempt := range attempts {
		byID[attempt.Metadata.Name] = attempt
	}
	var selected *goalv1alpha1.GoalAttempt
	if attempt, ok := byID[goal.Status.ActiveAttemptID]; ok &&
		attempt.Status.State == "open" {
		selected = &attempt
	} else {
		for index := range attempts {
			attempt := attempts[index]
			if attempt.Status.State != "open" {
				continue
			}
			if selected == nil || attempt.Status.ObservedAt >
				selected.Status.ObservedAt {
				selected = &attempts[index]
			}
		}
	}
	if selected == nil {
		// An open goal without an open attempt has no resumable work yet.
		return nil
	}
	return &goalcatalog.GoalContinuation{
		ActiveAttempt:    selected.Metadata.Name,
		StableDefect:     selected.Spec.StableDefect,
		Hypothesis:       selected.Spec.Hypothesis,
		Subject:          selected.Spec.Subject,
		AffectedCriteria: selected.Spec.AffectedCriteria,
		RegressionRefs:   selected.Spec.RegressionRefs,
		PriorAttemptID:   selected.Spec.PriorAttemptID,
		DominantFailure:  selected.Spec.DominantFailure,
		MeasurableDelta:  selected.Spec.MeasurableDelta,
		NextAction:       selected.Spec.NextAction,
		Blocker:          selected.Spec.Blocker,
		ResumeCondition:  selected.Spec.ResumeCondition,
	}
}

// loadRecord parses every manifest of one goal record from disk without
// writing anything. It returns a bounded reason for any unavailable record.
func (c *compiler) loadRecord(goalDir string) (
	goalv1alpha1.Goal,
	goalv1alpha1.GoalCriteria,
	goalv1alpha1.GoalCriteria,
	[]goalv1alpha1.GoalAttempt,
	string,
) {
	var goal goalv1alpha1.Goal
	if !c.loadYAML(goalDir, "goal.yaml", "goal-manifest", &goal) {
		return goal, goalv1alpha1.GoalCriteria{},
			goalv1alpha1.GoalCriteria{}, nil, "invalid record"
	}
	if goal.Metadata.Name != filepath.Base(goalDir) {
		return goal, goalv1alpha1.GoalCriteria{},
			goalv1alpha1.GoalCriteria{}, nil,
			"invalid record: goal name does not match directory"
	}
	var criteria goalv1alpha1.GoalCriteria
	if !c.loadYAML(goalDir, currentCriteriaFileName, "criteria-manifest", &criteria) {
		return goal, criteria, goalv1alpha1.GoalCriteria{}, nil, "invalid record"
	}
	if criteria.Metadata.Name != filepath.Base(goalDir) {
		return goal, criteria, goalv1alpha1.GoalCriteria{}, nil,
			"invalid record: criteria name does not match directory"
	}
	criteriaRevision := goalv1alpha1.GoalCriteria{}
	revisionPath := filepath.ToSlash(filepath.Join(
		criteriaRevisionsDirectory,
		fmt.Sprintf("%d.yaml", goal.Status.CriteriaRevision),
	))
	if !c.readYAMLOnly(goalDir, revisionPath, &criteriaRevision) {
		return goal, criteria, criteriaRevision, nil, "invalid record"
	}
	attemptsDirectoryPath := filepath.ToSlash(filepath.Join(goalDir, attemptsDirectory))
	attemptEntries, err := os.ReadDir(filepath.Join(
		c.root,
		filepath.FromSlash(attemptsDirectoryPath),
	))
	if err != nil {
		return goal, criteria, criteriaRevision, nil, "invalid record"
	}
	attempts := make([]goalv1alpha1.GoalAttempt, 0, len(attemptEntries))
	for _, attemptEntry := range attemptEntries {
		if !attemptEntry.IsDir() || strings.HasPrefix(attemptEntry.Name(), ".") {
			continue
		}
		attemptDir := filepath.ToSlash(filepath.Join(
			attemptsDirectoryPath,
			attemptEntry.Name(),
		))
		content := attemptManifestFile(c.root, attemptDir)
		if len(content) == 0 {
			return goal, criteria, criteriaRevision, attempts, "invalid record"
		}
		var attempt goalv1alpha1.GoalAttempt
		if err := yaml.UnmarshalWithOptions(content, &attempt, yaml.Strict()); err != nil {
			return goal, criteria, criteriaRevision, attempts, "invalid record"
		}
		if attempt.Metadata.Name != attemptEntry.Name() {
			return goal, criteria, criteriaRevision, attempts,
				"invalid record: attempt name does not match directory"
		}
		attempts = append(attempts, attempt)
	}
	return goal, criteria, criteriaRevision, attempts, ""
}

// loadYAML reads one YAML manifest, records it as an input, parses it, and
// checks that its metadata name matches the goal directory name.
func (c *compiler) loadYAML(goalDir, relativePath, role string,
	destination any,
) bool {
	path := filepath.ToSlash(filepath.Join(goalDir, relativePath))
	full := filepath.Join(c.root, filepath.FromSlash(path))
	content, err := os.ReadFile(full)
	if err != nil {
		c.problem("manifest missing: %s", path)
		return false
	}
	c.input(path, role, content)
	if err := yaml.UnmarshalWithOptions(content, destination, yaml.Strict()); err != nil {
		c.problem("manifest unparseable: %s", path)
		return false
	}
	return true
}

// readYAMLOnly reads and strictly parses one manifest without recording it
// as a catalog input. It is used for required record components that are not
// part of the declared input universe (criteria revisions and attempt files).
func (c *compiler) readYAMLOnly(goalDir, relativePath string,
	destination any,
) bool {
	path := filepath.ToSlash(filepath.Join(goalDir, relativePath))
	full := filepath.Join(c.root, filepath.FromSlash(path))
	content, err := os.ReadFile(full)
	if err != nil {
		c.problem("manifest missing: %s", path)
		return false
	}
	if err := yaml.UnmarshalWithOptions(content, destination, yaml.Strict()); err != nil {
		c.problem("manifest unparseable: %s", path)
		return false
	}
	return true
}

// validateCompleteRecord checks every parsed manifest against the goal API
// and the compiler's read-only directory/record invariants.
func validateCompleteRecord(
	goal goalv1alpha1.Goal,
	criteria goalv1alpha1.GoalCriteria,
	criteriaRevision goalv1alpha1.GoalCriteria,
	attempts []goalv1alpha1.GoalAttempt,
) error {
	if err := goal.Validate(); err != nil {
		return fmt.Errorf("goal: %w", err)
	}
	ownerRoot, ok := goal.Metadata.Annotations[goalv1alpha1.LocalOwnerRootAnnotation]
	if !ok {
		return fmt.Errorf("goal.yaml lacks local owner root annotation")
	}
	if ownerRoot == "" || filepath.IsAbs(ownerRoot) ||
		strings.Contains(ownerRoot, `\`) {
		return fmt.Errorf("goal.yaml has malformed local owner root annotation")
	}
	cleanOwner := filepath.ToSlash(filepath.Clean(filepath.FromSlash(ownerRoot)))
	if cleanOwner != ownerRoot || cleanOwner == "." ||
		strings.HasPrefix(cleanOwner, "../") {
		return fmt.Errorf("goal.yaml has non-normalized local owner root annotation")
	}
	if err := criteria.ValidateForGoal(goal); err != nil {
		return fmt.Errorf("criteria: %w", err)
	}
	if err := criteriaRevision.ValidateSnapshot(
		goal.Metadata.Name,
		goal.Status.CriteriaRevision,
	); err != nil {
		return fmt.Errorf("criteria revision: %w", err)
	}
	if goal.Status.CriteriaRevision != criteria.Spec.Revision {
		return fmt.Errorf("criteria revision does not match goal status")
	}
	for _, attempt := range attempts {
		if err := attempt.ValidateForGoal(goal); err != nil {
			return fmt.Errorf("attempt %s: %w", attempt.Metadata.Name, err)
		}
	}
	return nil
}

// attemptManifestFile returns the attempt.yaml content after checking the
// required per-attempt artifacts (plan.md, result.md, evidence/*.md)
// read-only.
func attemptManifestFile(root, attemptDir string) []byte {
	content, err := os.ReadFile(filepath.Join(
		root,
		filepath.FromSlash(filepath.Join(attemptDir, attemptManifestFileName)),
	))
	if err != nil {
		return nil
	}
	for _, required := range []string{planFileName, resultFileName} {
		if _, err := os.Stat(filepath.Join(
			root,
			filepath.FromSlash(filepath.Join(attemptDir, required)),
		)); err != nil {
			return nil
		}
	}
	evidenceDir := filepath.Join(root,
		filepath.FromSlash(filepath.Join(attemptDir, evidenceDirectory)))
	entries, err := os.ReadDir(evidenceDir)
	if err != nil {
		return nil
	}
	found := false
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" {
			continue
		}
		found = true
		break
	}
	if !found {
		return nil
	}
	return content
}

// boundedReason converts a validation error into a bounded stable reason,
// never embedding unvalidated YAML claims.
func boundedReason(reason error) string {
	if reason == nil {
		return "invalid record"
	}
	message := strings.TrimSpace(reason.Error())
	if message == "" {
		return "invalid record"
	}
	message = strings.Map(func(character rune) rune {
		if character < 32 || character > 126 {
			return ' '
		}
		return character
	}, message)
	fields := strings.Fields(message)
	if len(fields) == 0 {
		return "invalid record"
	}
	if len(fields) > 24 {
		fields = fields[:24]
	}
	return strings.Join(fields, " ")
}

func (c *compiler) catalog() (goalcatalog.GoalCatalog, error) {
	sourceRevision := c.opts.sourceRevision
	if sourceRevision == "" {
		var aggregate bytes.Buffer
		for _, input := range c.inputs {
			aggregate.WriteString(input.Path)
			aggregate.WriteString("\x00")
			aggregate.WriteString(input.Digest)
			aggregate.WriteString("\n")
		}
		value := sha256.Sum256(aggregate.Bytes())
		sourceRevision = hex.EncodeToString(value[:16])
	}
	completeness := goalcatalog.CompletenessComplete
	limitations := []string{}
	if len(c.problems) > 0 {
		completeness = goalcatalog.CompletenessPartial
		limitations = c.problems
	}
	inputs := c.inputs
	if inputs == nil {
		inputs = []goalcatalog.CatalogInput{}
	}
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Path != inputs[j].Path {
			return inputs[i].Path < inputs[j].Path
		}
		return inputs[i].Role < inputs[j].Role
	})
	goals := c.goals
	if goals == nil {
		goals = []goalcatalog.GoalRecord{}
	}
	sort.Slice(goals, func(i, j int) bool {
		return goals[i].CandidatePath < goals[j].CandidatePath
	})
	return goalcatalog.GoalCatalog{
		CatalogEnvelope: goalcatalog.CatalogEnvelope{
			Schema:            goalcatalog.APIVersion + "/" + goalcatalog.KindGoalCatalog,
			Kind:              goalcatalog.KindGoalCatalog,
			ID:                "agent-system.goal",
			DerivationVersion: "1.0.0",
			ProducerRef:       c.opts.producerRef,
			SourceRevision:    sourceRevision,
			Inputs:            inputs,
			Bounds: goalcatalog.CatalogBounds{
				Eligible:       c.eligible,
				Emitted:        c.emitted,
				Unavailable:    c.unavailable,
				MaxItems:       1000,
				MaxInputBytes:  32 << 20,
				MaxOutputBytes: 1 << 20,
			},
			Completeness: completeness,
			Limitations:  limitations,
			Conflicts:    []goalcatalog.CatalogConflict{},
		},
		Goals: goals,
	}, nil
}

func (c *compiler) buildOutputs(catalog goalcatalog.GoalCatalog) ([]byte, string, error) {
	jsonContent, err := goalcatalog.CanonicalJSONGoal(catalog)
	if err != nil {
		return nil, "", fmt.Errorf("canonical JSON: %w", err)
	}
	digestValue, err := goalcatalog.DecodeGoalStrict(jsonContent)
	if err != nil {
		return nil, "", err
	}
	catalog.Digest = digestValue.Digest
	return jsonContent, goalcatalog.RenderGoalMarkdown(catalog), nil
}

func Run(args []string, stdout io.Writer) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}
	compiler := &compiler{root: opts.workspaceRoot, opts: opts}
	if err := compiler.compile(); err != nil {
		return err
	}
	catalog, err := compiler.catalog()
	if err != nil {
		return err
	}
	jsonContent, markdown, err := compiler.buildOutputs(catalog)
	if err != nil {
		return err
	}
	if opts.check {
		jsonPath := filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))
		tracked, err := os.ReadFile(jsonPath)
		if err != nil {
			return fmt.Errorf("checked goal JSON unavailable: %w", err)
		}
		if !bytes.Equal(tracked, jsonContent) {
			return fmt.Errorf("checked goal JSON is stale; rerun the generator")
		}
		markdownPath := filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))
		trackedMarkdown, err := os.ReadFile(markdownPath)
		if err != nil {
			return fmt.Errorf("checked goal Markdown unavailable: %w", err)
		}
		if !bytes.Equal(trackedMarkdown, []byte(markdown)) {
			return fmt.Errorf("checked goal Markdown is stale; rerun the generator")
		}
	}
	if err := os.MkdirAll(
		filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath))),
		0o755,
	); err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(compiler.root, filepath.FromSlash(opts.outputPath)),
		jsonContent,
		0o644,
	); err != nil {
		return err
	}
	if err := os.MkdirAll(
		filepath.Dir(filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath))),
		0o755,
	); err != nil {
		return err
	}
	if err := os.WriteFile(
		filepath.Join(compiler.root, filepath.FromSlash(opts.markdownPath)),
		[]byte(markdown),
		0o644,
	); err != nil {
		return err
	}
	report := struct {
		CatalogID    string   `json:"catalogID"`
		Output       string   `json:"output"`
		Markdown     string   `json:"markdown"`
		Completeness string   `json:"completeness"`
		Status       string   `json:"status"`
		Problems     []string `json:"problems,omitempty"`
	}{
		CatalogID:    catalog.ID,
		Output:       opts.outputPath,
		Markdown:     opts.markdownPath,
		Completeness: string(catalog.Completeness),
		Status:       "ok",
		Problems:     compiler.problems,
	}
	if len(compiler.problems) > 0 {
		report.Status = "incomplete"
	}
	content, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')
	if _, err := stdout.Write(content); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := Run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "goal_check:", err)
		os.Exit(1)
	}
}
