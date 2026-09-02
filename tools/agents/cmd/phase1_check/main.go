package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	registrySchema   = "agents.alwaldend.com/phase1-registry/v1alpha1"
	operationsSchema = "agents.alwaldend.com/operations/v1alpha1"
	baselineSchema   = "agents.alwaldend.com/resource-baseline/v1alpha1"
)

var identifierPattern = regexp.MustCompile(
	`^[a-z][a-z0-9]*(?:[._/-][a-z0-9]+)*$`,
)

type registry struct {
	Schema               string              `json:"schema"`
	CriteriaRevision     int                 `json:"criteriaRevision"`
	Authorities          []authority         `json:"authorities"`
	Skills               []skill             `json:"skills"`
	RuntimeTools         []runtimeTool       `json:"runtimeTools"`
	OperationFiles       []string            `json:"operationFiles"`
	DirectBinaries       []directBinary      `json:"directBinaries"`
	GeneratedArtifacts   []generatedArtifact `json:"generatedArtifacts"`
	TerraformDefinitions string              `json:"terraformDefinitions"`
	CordisMCPSource      string              `json:"cordisMcpSource"`
}

type authority struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Source string `json:"source"`
}

type skill struct {
	ID                   string   `json:"id"`
	Owner                string   `json:"owner"`
	Layer                string   `json:"layer"`
	Activation           string   `json:"activation"`
	Exclusions           []string `json:"exclusions"`
	CapabilityRefs       []string `json:"capabilityRefs"`
	Dependencies         []string `json:"dependencies"`
	Conflicts            []string `json:"conflicts"`
	ProviderRequirements []string `json:"providerRequirements"`
	ContextCost          string   `json:"contextCost"`
	EvaluationMaturity   string   `json:"evaluationMaturity"`
}

type runtimeTool struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Classification string `json:"classification"`
}

type directBinary struct {
	ID             string `json:"id"`
	Owner          string `json:"owner"`
	Path           string `json:"path"`
	Classification string `json:"classification"`
}

type generatedArtifact struct {
	Path          string `json:"path"`
	Owner         string `json:"owner"`
	Formatter     string `json:"formatter"`
	CheckedOutput string `json:"checkedOutput"`
	Updater       string `json:"updater"`
	Exclusion     string `json:"exclusion,omitempty"`
}

type operationsFile struct {
	Schema     string      `json:"schema"`
	Owner      string      `json:"owner"`
	Provider   string      `json:"provider"`
	Definition string      `json:"definition"`
	Operations []operation `json:"operations"`
	Removed    []removed   `json:"removedAliases"`
}

type operation struct {
	ID                  string   `json:"id"`
	Selector            string   `json:"selector"`
	Classification      string   `json:"classification"`
	Effects             []string `json:"effects"`
	Inputs              []string `json:"inputs"`
	Outputs             []string `json:"outputs"`
	Information         []string `json:"information"`
	CredentialUse       string   `json:"credentialUse"`
	NetworkUse          string   `json:"networkUse"`
	EnvironmentSelector string   `json:"environmentSelector"`
	AuthorityGate       string   `json:"authorityGate"`
	Preflight           string   `json:"preflight"`
	Verification        string   `json:"verification"`
	Cost                string   `json:"cost"`
	Cacheability        string   `json:"cacheability"`
	Cancellation        string   `json:"cancellation"`
}

type removed struct {
	Selector    string `json:"selector"`
	Replacement string `json:"replacement"`
	Reason      string `json:"reason"`
}

type baseline struct {
	Schema           string         `json:"schema"`
	CriteriaRevision int            `json:"criteriaRevision"`
	RegistryDigest   string         `json:"registryDigest"`
	Ceilings         resourceLimits `json:"ceilings"`
	Observations     []observation  `json:"observations"`
}

type resourceLimits struct {
	CorrectnessFailuresMax  int64 `json:"correctnessFailuresMax"`
	UnsafeActionsMax        int64 `json:"unsafeActionsMax"`
	DiscoveryCallsMax       int64 `json:"discoveryCallsMax"`
	ContextBytesMax         int64 `json:"contextBytesMax"`
	ColdDurationMSMax       int64 `json:"coldDurationMsMax"`
	WarmDurationMSMax       int64 `json:"warmDurationMsMax"`
	UniverseUnclassifiedMax int64 `json:"universeUnclassifiedMax"`
	RedundantChecksMax      int64 `json:"redundantChecksMax"`
	ReusedChecksMin         int64 `json:"reusedChecksMin"`
}

type observation struct {
	Metric string `json:"metric"`
	Value  *int64 `json:"value"`
	Reason string `json:"reason,omitempty"`
}

type report struct {
	Schema            string         `json:"schema"`
	CriteriaRevision  int            `json:"criteriaRevision"`
	RegistryDigest    string         `json:"registryDigest"`
	Counts            map[string]int `json:"counts"`
	Missing           []string       `json:"missing"`
	Unclassified      []string       `json:"unclassified"`
	RequiresMigration []string       `json:"requiresMigration"`
	Valid             bool           `json:"valid"`
}

type checker struct {
	root               string
	registryPath       string
	baselinePath       string
	registry           registry
	operations         map[string]operation
	operationSelectors map[string]map[string]bool
	report             report
}

func strictJSON(path string, destination any) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values")
		}
		return nil, err
	}
	return content, nil
}

func digest(content []byte) string {
	value := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(value[:])
}

func inside(root, path string) (string, error) {
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute path is forbidden: %s", path)
	}
	joined := filepath.Clean(filepath.Join(root, path))
	relative, err := filepath.Rel(root, joined)
	if err != nil || relative == ".." ||
		strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes workspace: %s", path)
	}
	return joined, nil
}

func (c *checker) issue(kind string, message string) {
	switch kind {
	case "missing":
		c.report.Missing = append(c.report.Missing, message)
	case "unclassified":
		c.report.Unclassified = append(c.report.Unclassified, message)
	case "requires_migration":
		c.report.RequiresMigration = append(
			c.report.RequiresMigration,
			message,
		)
	}
}

func validateID(value string) bool {
	return len(value) <= 253 && identifierPattern.MatchString(value) &&
		!strings.Contains(value, "..")
}

func (c *checker) load() error {
	registryPath, err := inside(c.root, c.registryPath)
	if err != nil {
		return err
	}
	content, err := strictJSON(registryPath, &c.registry)
	if err != nil {
		return fmt.Errorf("registry: %w", err)
	}
	c.report.RegistryDigest = digest(content)
	if c.registry.Schema != registrySchema || c.registry.CriteriaRevision != 4 {
		return fmt.Errorf("registry schema or criteria revision mismatch")
	}
	c.operations = map[string]operation{}
	c.operationSelectors = map[string]map[string]bool{}
	providerOwners := map[string]string{}
	for _, relative := range c.registry.OperationFiles {
		path, err := inside(c.root, relative)
		if err != nil {
			return err
		}
		var manifest operationsFile
		if _, err := strictJSON(path, &manifest); err != nil {
			return fmt.Errorf("operation file %s: %w", relative, err)
		}
		if manifest.Schema != operationsSchema || manifest.Owner == "" ||
			manifest.Provider == "" || manifest.Definition == "" {
			return fmt.Errorf("operation file %s has incomplete identity", relative)
		}
		if prior, found := providerOwners[manifest.Provider]; found &&
			prior != manifest.Owner {
			return fmt.Errorf(
				"provider %s has owners %s and %s",
				manifest.Provider,
				prior,
				manifest.Owner,
			)
		}
		providerOwners[manifest.Provider] = manifest.Owner
		selectors := map[string]bool{}
		for _, candidate := range manifest.Operations {
			if err := validateOperation(candidate); err != nil {
				return fmt.Errorf("operation file %s: %w", relative, err)
			}
			if _, found := c.operations[candidate.ID]; found {
				return fmt.Errorf("duplicate operation %s", candidate.ID)
			}
			if selectors[candidate.Selector] {
				return fmt.Errorf("duplicate selector %s", candidate.Selector)
			}
			c.operations[candidate.ID] = candidate
			selectors[candidate.Selector] = true
			if candidate.Classification == "requires_migration" {
				c.issue("requires_migration", "operation:"+candidate.ID)
			}
		}
		for _, alias := range manifest.Removed {
			if alias.Replacement == "" || alias.Reason == "" {
				return fmt.Errorf("operation file %s has invalid removed alias", relative)
			}
			if _, found := c.operations[alias.Replacement]; !found {
				return fmt.Errorf(
					"removed alias %q lacks replacement %q",
					alias.Selector,
					alias.Replacement,
				)
			}
			if selectors[alias.Selector] {
				return fmt.Errorf("removed alias %q is still registered", alias.Selector)
			}
		}
		c.operationSelectors[manifest.Definition] = selectors
	}
	return c.validateBaseline()
}

func validateOperation(candidate operation) error {
	if !validateID(candidate.ID) || candidate.Selector == "" ||
		candidate.AuthorityGate == "" || candidate.Preflight == "" ||
		candidate.Verification == "" || candidate.Cancellation == "" ||
		candidate.Cost == "" || candidate.EnvironmentSelector == "" ||
		len(candidate.Effects) == 0 || len(candidate.Information) == 0 {
		return fmt.Errorf("operation %q is incomplete", candidate.ID)
	}
	if candidate.Classification != "classified" &&
		candidate.Classification != "requires_migration" {
		return fmt.Errorf(
			"operation %s has unknown classification %q",
			candidate.ID,
			candidate.Classification,
		)
	}
	return nil
}

func (c *checker) validateBaseline() error {
	path, err := inside(c.root, c.baselinePath)
	if err != nil {
		return err
	}
	var value baseline
	if _, err := strictJSON(path, &value); err != nil {
		return fmt.Errorf("resource baseline: %w", err)
	}
	if value.Schema != baselineSchema || value.CriteriaRevision != 4 ||
		value.RegistryDigest != c.report.RegistryDigest {
		return fmt.Errorf("resource baseline is not bound to this registry")
	}
	limits := value.Ceilings
	if limits.CorrectnessFailuresMax != 0 || limits.UnsafeActionsMax != 0 ||
		limits.DiscoveryCallsMax <= 0 || limits.ContextBytesMax <= 0 ||
		limits.ColdDurationMSMax <= 0 || limits.WarmDurationMSMax <= 0 ||
		limits.UniverseUnclassifiedMax != 0 || limits.RedundantChecksMax < 0 ||
		limits.ReusedChecksMin < 0 {
		return fmt.Errorf("resource baseline has invalid numeric ceilings")
	}
	metrics := map[string]bool{}
	for _, item := range value.Observations {
		if item.Metric == "" || metrics[item.Metric] {
			return fmt.Errorf("resource baseline observation identity is invalid")
		}
		metrics[item.Metric] = true
		if item.Value == nil && item.Reason == "" {
			return fmt.Errorf("unavailable metric %s lacks a reason", item.Metric)
		}
		if item.Value != nil && item.Reason != "" {
			return fmt.Errorf("measured metric %s also declares a reason", item.Metric)
		}
	}
	return nil
}

func (c *checker) checkAuthorities() {
	seen := map[string]bool{}
	allowedKinds := map[string]bool{
		"bazel_operations": true,
		"direct_binaries":  true,
		"goals":            true,
		"projects":         true,
		"runtime_tools":    true,
		"skills":           true,
		"workspaces":       true,
	}
	for _, candidate := range c.registry.Authorities {
		if !validateID(candidate.ID) || !allowedKinds[candidate.Kind] ||
			candidate.Source == "" || seen[candidate.ID] {
			c.issue("unclassified", "authority:"+candidate.ID)
			continue
		}
		seen[candidate.ID] = true
		path, err := inside(c.root, candidate.Source)
		if err != nil {
			c.issue("missing", "authority:"+candidate.ID)
			continue
		}
		if _, err := os.Lstat(path); err != nil {
			c.issue("missing", "authority:"+candidate.ID)
		}
	}
	c.report.Counts["authorities"] = len(seen)
}

func (c *checker) checkSkills() error {
	root := filepath.Join(c.root, ".agents", "skills")
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	declared := map[string]skill{}
	for _, candidate := range c.registry.Skills {
		if !validateID(candidate.ID) || candidate.Owner == "" ||
			candidate.Layer == "" || candidate.Activation == "" ||
			candidate.ContextCost == "" || candidate.EvaluationMaturity == "" {
			c.issue("unclassified", "skill:"+candidate.ID)
		}
		declared[candidate.ID] = candidate
	}
	observed := map[string]bool{}
	for _, entry := range entries {
		observed[entry.Name()] = true
		if _, found := declared[entry.Name()]; !found {
			c.issue("missing", "skill:"+entry.Name())
		}
	}
	for id := range declared {
		if !observed[id] {
			c.issue("missing", "skill-discovery:"+id)
		}
	}
	c.report.Counts["skills"] = len(observed)
	return nil
}

func frontmatterStatuses(content []byte) []string {
	lines := strings.Split(string(content), "\n")
	insideStatuses := false
	var values []string
	for _, line := range lines {
		if line == "statuses:" {
			insideStatuses = true
			continue
		}
		if insideStatuses {
			if strings.HasPrefix(line, "  - ") {
				values = append(values, strings.TrimSpace(line[4:]))
				continue
			}
			break
		}
	}
	return values
}

func (c *checker) checkProjects() error {
	entries, err := os.ReadDir(filepath.Join(c.root, "projects"))
	if err != nil {
		return err
	}
	allowed := map[string]bool{
		"abandoned":    true,
		"active":       true,
		"experimental": true,
		"finished":     true,
		"in_progress":  true,
		"maintenance":  true,
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		count++
		content, err := os.ReadFile(
			filepath.Join(c.root, "projects", entry.Name(), "README.md"),
		)
		if err != nil {
			c.issue("missing", "project-readme:"+entry.Name())
			continue
		}
		statuses := frontmatterStatuses(content)
		if len(statuses) != 1 || !allowed[statuses[0]] {
			c.issue("unclassified", "project-lifecycle:"+entry.Name())
		}
	}
	c.report.Counts["projects"] = count
	return nil
}

func walkFiles(root string, wanted func(string) bool) ([]string, error) {
	var result []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			name := entry.Name()
			if path != root && (name == ".git" || name == "out" ||
				strings.HasPrefix(name, "bazel-")) {
				return filepath.SkipDir
			}
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if wanted(filepath.ToSlash(relative)) {
			result = append(result, filepath.ToSlash(relative))
		}
		return nil
	})
	sort.Strings(result)
	return result, err
}

func (c *checker) checkWorkspacesAndGoals() error {
	modules, err := walkFiles(c.root, func(path string) bool {
		return filepath.Base(path) == "MODULE.bazel"
	})
	if err != nil {
		return err
	}
	goals, err := walkFiles(c.root, func(path string) bool {
		return strings.Contains(path, "/goals/") &&
			filepath.Base(path) == "goal.yaml"
	})
	if err != nil {
		return err
	}
	c.report.Counts["workspaces"] = len(modules)
	c.report.Counts["goals"] = len(goals)
	return nil
}

func setDifference(observed, declared map[string]bool) ([]string, []string) {
	var missing []string
	var stale []string
	for value := range observed {
		if !declared[value] {
			missing = append(missing, value)
		}
	}
	for value := range declared {
		if !observed[value] {
			stale = append(stale, value)
		}
	}
	sort.Strings(missing)
	sort.Strings(stale)
	return missing, stale
}

func (c *checker) checkRuntimeTools() error {
	path, err := inside(c.root, c.registry.CordisMCPSource)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	pattern := regexp.MustCompile(
		`server\.registerTool\(\s*"(cordis_[a-z_]+)"`,
	)
	observed := map[string]bool{}
	for _, match := range pattern.FindAllSubmatch(content, -1) {
		observed[string(match[1])] = true
	}
	declared := map[string]bool{}
	for _, tool := range c.registry.RuntimeTools {
		if !validateID(strings.ReplaceAll(tool.ID, "_", ".")) ||
			tool.Owner == "" || (tool.Classification != "classified" &&
			tool.Classification != "requires_migration") {
			c.issue("unclassified", "runtime-tool:"+tool.ID)
		}
		declared[tool.ID] = true
		if tool.Classification == "requires_migration" {
			c.issue("requires_migration", "runtime-tool:"+tool.ID)
		}
	}
	missing, stale := setDifference(observed, declared)
	for _, id := range missing {
		c.issue("missing", "runtime-tool:"+id)
	}
	for _, id := range stale {
		c.issue("missing", "runtime-source:"+id)
	}
	c.report.Counts["runtimeTools"] = len(observed)
	return nil
}

func terraformSelectors(content []byte) map[string]bool {
	result := map[string]bool{}
	insideMap := false
	pattern := regexp.MustCompile(`^\s*"([^"]*)":\s*\[`)
	for _, line := range strings.Split(string(content), "\n") {
		if line == "DEFAULT_TERRAFORM_BINARIES = {" {
			insideMap = true
			continue
		}
		if insideMap && line == "}" {
			break
		}
		if insideMap {
			if match := pattern.FindStringSubmatch(line); match != nil {
				result[match[1]] = true
			}
		}
	}
	return result
}

func (c *checker) checkOperations() error {
	path, err := inside(c.root, c.registry.TerraformDefinitions)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	observed := terraformSelectors(content)
	if observed[""] {
		c.issue("unclassified", "terraform-selector:<unnamed>")
	}
	declared := c.operationSelectors[c.registry.TerraformDefinitions]
	missing, stale := setDifference(observed, declared)
	for _, selector := range missing {
		c.issue("missing", "terraform-operation:"+selector)
	}
	for _, selector := range stale {
		c.issue("missing", "terraform-selector:"+selector)
	}
	c.report.Counts["operations"] = len(c.operations)
	return nil
}

func (c *checker) checkOwnedFiles() {
	for _, binary := range c.registry.DirectBinaries {
		if !validateID(binary.ID) || binary.Owner == "" ||
			(binary.Classification != "classified" &&
				binary.Classification != "requires_migration") {
			c.issue("unclassified", "direct-binary:"+binary.ID)
		}
		path, err := inside(c.root, binary.Path)
		if err != nil {
			c.issue("missing", "direct-binary:"+binary.ID)
			continue
		}
		if _, err := os.Stat(path); err != nil {
			c.issue("missing", "direct-binary:"+binary.ID)
		}
		if binary.Classification == "requires_migration" {
			c.issue("requires_migration", "direct-binary:"+binary.ID)
		}
	}
	for _, artifact := range c.registry.GeneratedArtifacts {
		if artifact.Path == "" || artifact.Owner == "" ||
			artifact.Formatter == "" || artifact.CheckedOutput == "" ||
			(artifact.Updater == "" && artifact.Exclusion == "") {
			c.issue("unclassified", "generated-artifact:"+artifact.Path)
		}
		path, err := inside(c.root, artifact.Path)
		if err != nil {
			c.issue("missing", "generated-artifact:"+artifact.Path)
			continue
		}
		if _, err := os.Stat(path); err != nil {
			c.issue("missing", "generated-artifact:"+artifact.Path)
		}
	}
	c.report.Counts["directBinaries"] = len(c.registry.DirectBinaries)
	c.report.Counts["generatedArtifacts"] = len(c.registry.GeneratedArtifacts)
}

func (c *checker) run() error {
	c.report = report{
		Schema:            "agents.alwaldend.com/phase1-report/v1alpha1",
		CriteriaRevision:  4,
		Counts:            map[string]int{},
		Missing:           []string{},
		Unclassified:      []string{},
		RequiresMigration: []string{},
	}
	if err := c.load(); err != nil {
		return err
	}
	c.checkAuthorities()
	if err := c.checkSkills(); err != nil {
		return err
	}
	if err := c.checkProjects(); err != nil {
		return err
	}
	if err := c.checkWorkspacesAndGoals(); err != nil {
		return err
	}
	if err := c.checkRuntimeTools(); err != nil {
		return err
	}
	if err := c.checkOperations(); err != nil {
		return err
	}
	c.checkOwnedFiles()
	sort.Strings(c.report.Missing)
	sort.Strings(c.report.Unclassified)
	sort.Strings(c.report.RequiresMigration)
	c.report.Valid = len(c.report.Missing) == 0 &&
		len(c.report.Unclassified) == 0
	return nil
}

func encodeReport(value report) ([]byte, error) {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(content, '\n'), nil
}

func main() {
	workspaceRoot := flag.String("workspace-root", "", "repository workspace root")
	registryPath := flag.String(
		"registry",
		"tools/agents/declarations/registry.json",
		"workspace-relative registry path",
	)
	baselinePath := flag.String(
		"baseline",
		"tools/agents/declarations/resource_baseline.json",
		"workspace-relative resource baseline path",
	)
	reportPath := flag.String("report", "", "optional output report path")
	flag.Parse()
	if *workspaceRoot == "" {
		*workspaceRoot = os.Getenv("BUILD_WORKSPACE_DIRECTORY")
	}
	if *workspaceRoot == "" {
		fmt.Fprintln(os.Stderr, "phase1_check: --workspace-root is required")
		os.Exit(2)
	}
	root, err := filepath.Abs(*workspaceRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phase1_check: %v\n", err)
		os.Exit(2)
	}
	value := checker{
		root:         root,
		registryPath: *registryPath,
		baselinePath: *baselinePath,
	}
	if err := value.run(); err != nil {
		fmt.Fprintf(os.Stderr, "phase1_check: %v\n", err)
		os.Exit(1)
	}
	content, err := encodeReport(value.report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "phase1_check: %v\n", err)
		os.Exit(1)
	}
	if *reportPath != "" {
		path, pathErr := inside(root, *reportPath)
		if pathErr != nil {
			fmt.Fprintf(os.Stderr, "phase1_check: %v\n", pathErr)
			os.Exit(1)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "phase1_check: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "phase1_check: %v\n", err)
			os.Exit(1)
		}
	}
	_, _ = os.Stdout.Write(content)
	if !value.report.Valid {
		os.Exit(1)
	}
}
