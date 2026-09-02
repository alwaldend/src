package fsstore

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/goccy/go-yaml"
)

type Store struct {
	workspaceRoot string
	runtimeRoot   string
	now           func() time.Time
	random        io.Reader
	beforeRename  func(string) error
}

func NewStore(workspaceRoot string) (*Store, error) {
	return NewStoreWithRuntimeDir(workspaceRoot, os.Getenv("XDG_RUNTIME_DIR"))
}

func NewStoreWithRuntimeDir(
	workspaceRoot string,
	runtimeDirectory string,
) (*Store, error) {
	root, err := absoluteClean(workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve workspace root: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("workspace root is not a directory: %s", root)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, fmt.Errorf("resolve workspace root symlinks: %w", err)
	}
	runtimeRoot, err := resolveRuntimeRoot(runtimeDirectory)
	if err != nil {
		return nil, err
	}
	if isWithin(root, goalLockRootPath(runtimeRoot)) {
		return nil, fmt.Errorf(
			"XDG_RUNTIME_DIR/alwaldend/goal/locks must be outside the workspace",
		)
	}
	return &Store{
		workspaceRoot: root,
		runtimeRoot:   runtimeRoot,
		now:           time.Now,
		random:        rand.Reader,
	}, nil
}

func resolveRuntimeRoot(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("XDG_RUNTIME_DIR is required for per-goal locks")
	}
	if !filepath.IsAbs(value) {
		return "", fmt.Errorf("XDG_RUNTIME_DIR must be an absolute path")
	}
	info, err := os.Lstat(value)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("XDG_RUNTIME_DIR is not a directory: %s", value)
	}
	if info.Mode().Perm() != 0o700 {
		return "", fmt.Errorf("XDG_RUNTIME_DIR must have mode 0700")
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok && int(stat.Uid) != os.Getuid() {
		return "", fmt.Errorf("XDG_RUNTIME_DIR must be owned by the current user")
	}
	resolved, err := filepath.EvalSymlinks(value)
	if err != nil {
		return "", fmt.Errorf("resolve XDG_RUNTIME_DIR symlinks: %w", err)
	}
	return filepath.Clean(resolved), nil
}

type InitOptions struct {
	GoalsRoot string
	Title     string
	GoalID    string
	Scope     string
	OwnerRoot string
	Criteria  []string
	Retention string
}

type GoalReference struct {
	GoalID          string `json:"goalID"`
	GoalRef         string `json:"goalRef"`
	ResourceVersion string `json:"resourceVersion"`
}

func (s *Store) Init(options InitOptions) (GoalReference, error) {
	goalsRoot, err := s.resolveInsideWorkspace(options.GoalsRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("goals root: %w", err)
	}
	if err := os.MkdirAll(goalsRoot, 0o755); err != nil {
		return GoalReference{}, fmt.Errorf("create goals root: %w", err)
	}
	goalID := options.GoalID
	if goalID == "" {
		goalID, err = s.generateGoalID(options.Title)
		if err != nil {
			return GoalReference{}, err
		}
	}
	if err := validateRecordID("goal ID", goalID); err != nil {
		return GoalReference{}, err
	}
	if options.Scope == "" {
		options.Scope = "workspace"
	}
	if options.Retention == "" {
		if options.Scope == "project" {
			options.Retention = "durable"
		} else {
			options.Retention = "ephemeral"
		}
	}
	if len(options.Criteria) > maxCriteria {
		return GoalReference{}, fmt.Errorf("criteria cardinality exceeds %d", maxCriteria)
	}
	goalDir := filepath.Join(goalsRoot, goalID)
	lock, err := s.acquireGoalLock(goalDir)
	if err != nil {
		return GoalReference{}, err
	}
	defer lock.release()
	if pathExists(goalDir) {
		return GoalReference{}, fmt.Errorf("goal %q already exists", goalID)
	}
	ownerRoot := options.OwnerRoot
	if ownerRoot == "" {
		ownerRoot = filepath.Dir(goalsRoot)
	}
	ownerRoot, err = s.resolveInsideWorkspace(ownerRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("owner root: %w", err)
	}
	ownerRef, err := portableOwnerRoot(s.workspaceRoot, ownerRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("owner root reference: %w", err)
	}
	now := s.timestamp()
	goal := GoalManifest{
		APIVersion: goalAPIVersion,
		Kind:       "Goal",
		Metadata: ObjectMeta{
			Name:              goalID,
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
			Annotations: map[string]string{
				localOwnerRootAnnotation: ownerRef,
			},
		},
		Spec: GoalSpec{
			Title:     strings.TrimSpace(options.Title),
			Scope:     options.Scope,
			Retention: Retention{Policy: options.Retention},
			Relationships: Relationships{
				DependsOnGoalRefs:  []LocalGoalReference{},
				SupersedesGoalRefs: []LocalGoalReference{},
			},
		},
		Status: GoalStatus{
			LifecycleGeneration: 1,
			Outcome:             "open",
			Execution:           "active",
			CriteriaRevision:    1,
			ObservedAt:          now,
		},
	}
	criteria := CriteriaManifest{
		APIVersion: goalAPIVersion,
		Kind:       "GoalCriteria",
		Metadata: ObjectMeta{
			Name:              goalID,
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
		},
		Spec: CriteriaSpec{
			GoalRef:  LocalGoalReference{Name: goalID},
			Revision: 1,
			Items:    make([]Criterion, 0, len(options.Criteria)),
		},
	}
	for index, statement := range options.Criteria {
		criteria.Spec.Items = append(criteria.Spec.Items, Criterion{
			CriterionID:    fmt.Sprintf("criterion-%03d", index+1),
			Revision:       1,
			Required:       true,
			Statement:      strings.TrimSpace(statement),
			EvidenceMethod: "Inspect linked evidence against the criterion.",
		})
	}
	if err := goal.validate(); err != nil {
		return GoalReference{}, fmt.Errorf("invalid goal: %w", err)
	}
	if err := criteria.validate(goal); err != nil {
		return GoalReference{}, fmt.Errorf("invalid criteria: %w", err)
	}
	temporary, err := os.MkdirTemp(goalsRoot, ".goal-init-")
	if err != nil {
		return GoalReference{}, fmt.Errorf("create initialization directory: %w", err)
	}
	installed := false
	defer func() {
		if !installed {
			_ = os.RemoveAll(temporary)
		}
	}()
	if err := os.Mkdir(filepath.Join(temporary, "attempts"), 0o755); err != nil {
		return GoalReference{}, err
	}
	if err := os.Mkdir(filepath.Join(temporary, "criteria-revisions"), 0o755); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(temporary, "goal.yaml"), goal); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(temporary, "criteria.yaml"), criteria); err != nil {
		return GoalReference{}, err
	}
	criteriaBytes, err := marshalYAML(criteria)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(
		criteriaSnapshotPath(temporary, 1),
		criteriaBytes,
		0o644,
	); err != nil {
		return GoalReference{}, err
	}
	readme, err := renderREADME(goal, criteria, nil, defaultOutputLimit)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(filepath.Join(temporary, "README.md"), readme, 0o644); err != nil {
		return GoalReference{}, err
	}
	if err := s.callBeforeRename(goalDir); err != nil {
		return GoalReference{}, err
	}
	if err := os.Rename(temporary, goalDir); err != nil {
		if errors.Is(err, os.ErrExist) || pathExists(goalDir) {
			return GoalReference{}, fmt.Errorf("goal %q already exists", goalID)
		}
		return GoalReference{}, fmt.Errorf("install goal: %w", err)
	}
	installed = true
	return GoalReference{
		GoalID:          goalID,
		GoalRef:         goalID,
		ResourceVersion: goal.Metadata.ResourceVersion,
	}, nil
}

type GoalSummary struct {
	GoalID              string `json:"goalID"`
	GoalRef             string `json:"goalRef"`
	Title               string `json:"title"`
	Scope               string `json:"scope"`
	ResourceVersion     string `json:"resourceVersion"`
	Generation          uint64 `json:"generation"`
	LifecycleGeneration uint64 `json:"lifecycleGeneration"`
	Outcome             string `json:"outcome"`
	Execution           string `json:"execution"`
}

type ListResult struct {
	APIVersion string        `json:"apiVersion"`
	Kind       string        `json:"kind"`
	Goals      []GoalSummary `json:"goals"`
	Returned   int           `json:"returned"`
	Total      int           `json:"total"`
	Truncated  bool          `json:"truncated"`
}

func (s *Store) List(goalsRoot string, limit int) (ListResult, error) {
	root, err := s.resolveInsideWorkspace(goalsRoot)
	if err != nil {
		return ListResult{}, err
	}
	limit, err = validateLimit(limit)
	if err != nil {
		return ListResult{}, err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return ListResult{}, fmt.Errorf("read goals root: %w", err)
	}
	if len(entries) > maxGoals {
		return ListResult{}, fmt.Errorf("goal cardinality exceeds %d", maxGoals)
	}
	summaries := make([]GoalSummary, 0, min(limit, len(entries)))
	total := 0
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		lock, err := s.acquireGoalLock(dir)
		if err != nil {
			return ListResult{}, fmt.Errorf("goal %q: %w", entry.Name(), err)
		}
		goal, readErr := s.readGoalManifest(dir)
		_ = lock.release()
		err = readErr
		if err != nil {
			return ListResult{}, fmt.Errorf("goal %q: %w", entry.Name(), err)
		}
		if goal.Metadata.Name != entry.Name() {
			return ListResult{}, fmt.Errorf("goal directory %q does not match metadata.name", entry.Name())
		}
		total++
		if len(summaries) >= limit {
			continue
		}
		summaries = append(summaries, summaryFromGoal(entry.Name(), goal))
	}
	sort.Slice(summaries, func(i, j int) bool { return summaries[i].GoalID < summaries[j].GoalID })
	return ListResult{
		APIVersion: goalAPIVersion,
		Kind:       "GoalList",
		Goals:      summaries,
		Returned:   len(summaries),
		Total:      total,
		Truncated:  total > len(summaries),
	}, nil
}

type AttemptSummary struct {
	AttemptID           string `json:"attemptID"`
	ResourceVersion     string `json:"resourceVersion"`
	WorkType            string `json:"workType"`
	State               string `json:"state"`
	CriteriaRevision    uint64 `json:"criteriaRevision"`
	LifecycleGeneration uint64 `json:"lifecycleGeneration"`
	ObservedAt          string `json:"observedAt"`
}

type GoalView struct {
	APIVersion   string           `json:"apiVersion"`
	Kind         string           `json:"kind"`
	Goal         GoalManifest     `json:"goal"`
	Criteria     []Criterion      `json:"criteria"`
	Attempts     []AttemptSummary `json:"attempts"`
	Returned     int              `json:"returned"`
	Total        int              `json:"total"`
	Truncated    bool             `json:"truncated"`
	Session      *SessionBinding  `json:"session,omitempty"`
	SessionStale bool             `json:"sessionStale,omitempty"`
}

func (s *Store) ShowGoal(goalDir string, limit int) (GoalView, error) {
	dir, err := s.resolveInsideWorkspace(goalDir)
	if err != nil {
		return GoalView{}, err
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return GoalView{}, err
	}
	defer lock.release()
	if err := s.checkNoPendingPublication(dir); err != nil {
		return GoalView{}, err
	}
	goal, criteria, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return GoalView{}, err
	}
	return makeGoalView(goal, criteria, attempts, limit)
}

func makeGoalView(goal GoalManifest, criteria CriteriaManifest, attempts []AttemptManifest, limit int) (GoalView, error) {
	limit, err := validateLimit(limit)
	if err != nil {
		return GoalView{}, err
	}
	sort.Slice(attempts, func(i, j int) bool {
		if attempts[i].Status.ObservedAt == attempts[j].Status.ObservedAt {
			return attempts[i].Metadata.Name > attempts[j].Metadata.Name
		}
		return attempts[i].Status.ObservedAt > attempts[j].Status.ObservedAt
	})
	returned := min(limit, len(attempts))
	summaries := make([]AttemptSummary, 0, returned)
	for _, attempt := range attempts[:returned] {
		summaries = append(summaries, AttemptSummary{
			AttemptID:           attempt.Metadata.Name,
			ResourceVersion:     attempt.Metadata.ResourceVersion,
			WorkType:            attempt.Spec.WorkType,
			State:               attempt.Status.State,
			CriteriaRevision:    attempt.Spec.CriteriaRevision,
			LifecycleGeneration: attempt.Spec.LifecycleGeneration,
			ObservedAt:          attempt.Status.ObservedAt,
		})
	}
	criteriaItems := criteria.Spec.Items
	if len(criteriaItems) > limit {
		criteriaItems = criteriaItems[:limit]
	}
	return GoalView{
		APIVersion: goalAPIVersion,
		Kind:       "GoalView",
		Goal:       goal,
		Criteria:   criteriaItems,
		Attempts:   summaries,
		Returned:   returned,
		Total:      len(attempts),
		Truncated:  len(attempts) > returned || len(criteria.Spec.Items) > len(criteriaItems),
	}, nil
}

type AttachOptions struct {
	SessionRoot string
	SessionID   string
	GoalDir     string
}

func (s *Store) Attach(options AttachOptions) (SessionBinding, error) {
	if err := validateRecordID("session ID", options.SessionID); err != nil {
		return SessionBinding{}, err
	}
	sessionRoot, err := s.resolveSessionRoot(options.SessionRoot)
	if err != nil {
		return SessionBinding{}, fmt.Errorf("session root: %w", err)
	}
	goalDir, err := s.resolveInsideWorkspace(options.GoalDir)
	if err != nil {
		return SessionBinding{}, fmt.Errorf("goal directory: %w", err)
	}
	goalLock, err := s.acquireGoalLock(goalDir)
	if err != nil {
		return SessionBinding{}, err
	}
	goal, criteria, _, err := s.loadAndValidate(goalDir)
	goalLock.release()
	if err != nil {
		return SessionBinding{}, err
	}
	criteriaDigest, err := criteriaPortableDigest(criteria)
	if err != nil {
		return SessionBinding{}, err
	}
	stateDigest, err := goalStateDigest(goal)
	if err != nil {
		return SessionBinding{}, err
	}
	goalRef, err := portableRelative(s.workspaceRoot, goalDir, false)
	if err != nil {
		return SessionBinding{}, err
	}
	if err := os.MkdirAll(sessionRoot, 0o755); err != nil {
		return SessionBinding{}, fmt.Errorf("create session root: %w", err)
	}
	path := filepath.Join(sessionRoot, options.SessionID+".yaml")
	lock, err := s.acquirePathLock(path)
	if err != nil {
		return SessionBinding{}, err
	}
	defer lock.release()
	now := s.timestamp()
	binding := SessionBinding{
		APIVersion: goalAPIVersion,
		Kind:       "GoalSessionBinding",
		Metadata: ObjectMeta{
			Name:              options.SessionID,
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
			Annotations: map[string]string{
				localGoalReferenceAnnotation: goalRef,
			},
		},
		Spec: SessionSpec{GoalRef: LocalGoalReference{Name: goal.Metadata.Name}},
		Status: SessionStatus{
			GoalID:                      goal.Metadata.Name,
			AttachedGeneration:          goal.Metadata.Generation,
			AttachedLifecycleGeneration: goal.Status.LifecycleGeneration,
			AttachedCriteriaRevision:    criteria.Spec.Revision,
			AttachedCriteriaDigest:      criteriaDigest,
			AttachedGoalStateDigest:     stateDigest,
			ObservedAt:                  now,
		},
	}
	if pathExists(path) {
		var previous SessionBinding
		if err := s.readYAML(path, &previous); err != nil {
			return SessionBinding{}, err
		}
		if err := previous.validate(); err != nil {
			return SessionBinding{}, fmt.Errorf("existing session binding: %w", err)
		}
		next, err := incrementResourceVersion(previous.Metadata.ResourceVersion)
		if err != nil {
			return SessionBinding{}, err
		}
		binding.Metadata.ResourceVersion = next
		binding.Metadata.CreationTimestamp = previous.Metadata.CreationTimestamp
		if previous.Spec.GoalRef.Name == goal.Metadata.Name &&
			previous.Metadata.Annotations[localGoalReferenceAnnotation] == goalRef {
			binding.Metadata.Generation = previous.Metadata.Generation
		} else {
			binding.Metadata.Generation = previous.Metadata.Generation + 1
		}
	}
	if err := binding.validate(); err != nil {
		return SessionBinding{}, err
	}
	if err := s.writeYAML(path, binding); err != nil {
		return SessionBinding{}, err
	}
	return binding, nil
}

func (s *Store) ShowSession(sessionRoot string, sessionID string, limit int) (GoalView, error) {
	if err := validateRecordID("session ID", sessionID); err != nil {
		return GoalView{}, err
	}
	root, err := s.resolveSessionRoot(sessionRoot)
	if err != nil {
		return GoalView{}, err
	}
	var binding SessionBinding
	if err := s.readYAML(filepath.Join(root, sessionID+".yaml"), &binding); err != nil {
		return GoalView{}, fmt.Errorf("read session binding: %w", err)
	}
	if err := binding.validate(); err != nil {
		return GoalView{}, fmt.Errorf("invalid session binding: %w", err)
	}
	if binding.Metadata.Name != sessionID {
		return GoalView{}, fmt.Errorf("session binding filename does not match metadata.name")
	}
	goalDir, err := s.resolveInsideWorkspace(
		binding.Metadata.Annotations[localGoalReferenceAnnotation],
	)
	if err != nil {
		return GoalView{}, err
	}
	lock, err := s.acquireGoalLock(goalDir)
	if err != nil {
		return GoalView{}, err
	}
	defer lock.release()
	goal, criteria, attempts, err := s.loadAndValidate(goalDir)
	if err != nil {
		return GoalView{}, err
	}
	view, err := makeGoalView(goal, criteria, attempts, limit)
	if err != nil {
		return GoalView{}, err
	}
	if view.Goal.Metadata.Name != binding.Status.GoalID {
		return GoalView{}, fmt.Errorf("session goal ID no longer matches referenced record")
	}
	view.Session = &binding
	criteriaDigest, err := criteriaPortableDigest(criteria)
	if err != nil {
		return GoalView{}, err
	}
	stateDigest, err := goalStateDigest(view.Goal)
	if err != nil {
		return GoalView{}, err
	}
	view.SessionStale = binding.Status.AttachedGeneration != view.Goal.Metadata.Generation ||
		binding.Status.AttachedLifecycleGeneration != view.Goal.Status.LifecycleGeneration ||
		binding.Status.AttachedCriteriaRevision != view.Goal.Status.CriteriaRevision ||
		binding.Status.AttachedCriteriaDigest != criteriaDigest ||
		binding.Status.AttachedGoalStateDigest != stateDigest
	return view, nil
}

type CheckpointOptions struct {
	GoalDir                 string
	ExpectedResourceVersion string
	AttemptID               string
	PlanID                  string
	PlanStrategy            string
	PlanState               string
	PlanRejectionReason     string
	PlanOnly                bool
	WorkType                string
	PlanFile                string
	ResultFile              string
	EvidenceFiles           []string
	ReviewFile              string
	CriteriaFile            string
	CloseAttempt            bool
	Outcome                 string
	Execution               string
	StableDefect            string
	Hypothesis              string
	Subject                 string
	AffectedCriteria        []string
	RegressionRefs          []string
	PriorAttemptID          string
	DominantFailure         string
	MeasurableDelta         string
	NextAction              string
	Blocker                 string
	ResumeCondition         string
}

type desiredCriterion struct {
	CriterionID    string `json:"criterionID" yaml:"criterionID"`
	Required       *bool  `json:"required,omitempty" yaml:"required,omitempty"`
	Statement      string `json:"statement" yaml:"statement"`
	EvidenceMethod string `json:"evidenceMethod" yaml:"evidenceMethod"`
}

type desiredCriteria struct {
	Items []desiredCriterion `json:"items" yaml:"items"`
}

func nextCriteriaRevision(current uint64) (uint64, error) {
	if current >= maxCriteriaRevisions {
		return 0, fmt.Errorf(
			"criteria revision cardinality cannot exceed %d",
			maxCriteriaRevisions,
		)
	}
	return current + 1, nil
}

func (s *Store) updateCriteria(
	dir string,
	goal GoalManifest,
	desiredPath string,
) (GoalReference, error) {
	var current CriteriaManifest
	if err := s.readYAML(filepath.Join(dir, "criteria.yaml"), &current); err != nil {
		return GoalReference{}, err
	}
	_, _, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return GoalReference{}, err
	}
	nextRevision, err := nextCriteriaRevision(current.Spec.Revision)
	if err != nil {
		return GoalReference{}, err
	}
	history, err := s.loadCriteriaHistory(dir, goal, current)
	if err != nil {
		return GoalReference{}, err
	}
	desiredContent, err := s.readWorkspaceRegularFile(desiredPath, maxManifestBytes)
	if err != nil {
		return GoalReference{}, fmt.Errorf("read criteria file: %w", err)
	}
	var desired desiredCriteria
	if err := yaml.UnmarshalWithOptions(desiredContent, &desired, yaml.Strict()); err != nil {
		return GoalReference{}, fmt.Errorf("decode criteria file strictly: %w", err)
	}
	if len(desired.Items) > maxCriteria {
		return GoalReference{}, fmt.Errorf("criteria cardinality exceeds %d", maxCriteria)
	}
	existing := map[string]Criterion{}
	for _, item := range current.Spec.Items {
		existing[item.CriterionID] = item
	}
	historicalLatest := map[string]Criterion{}
	historicalMaximum := map[string]uint64{}
	for revision := uint64(1); revision <= current.Spec.Revision; revision++ {
		for _, item := range history.Snapshots[revision].Spec.Items {
			if item.Revision >= historicalMaximum[item.CriterionID] {
				historicalMaximum[item.CriterionID] = item.Revision
				historicalLatest[item.CriterionID] = item
			}
		}
	}
	items := make([]Criterion, 0, len(desired.Items))
	seen := map[string]bool{}
	for _, input := range desired.Items {
		if seen[input.CriterionID] {
			return GoalReference{}, fmt.Errorf("duplicate criterionID %q", input.CriterionID)
		}
		seen[input.CriterionID] = true
		required := true
		if input.Required != nil {
			required = *input.Required
		}
		item := Criterion{
			CriterionID:    input.CriterionID,
			Revision:       1,
			Required:       required,
			Statement:      strings.TrimSpace(input.Statement),
			EvidenceMethod: strings.TrimSpace(input.EvidenceMethod),
		}
		if old, ok := existing[item.CriterionID]; ok {
			item.Revision = old.Revision
			if old.Required != item.Required || old.Statement != item.Statement ||
				old.EvidenceMethod != item.EvidenceMethod {
				item.Revision++
			}
		} else if old, ok := historicalLatest[item.CriterionID]; ok {
			item.Revision = historicalMaximum[item.CriterionID]
			if old.Required != item.Required || old.Statement != item.Statement ||
				old.EvidenceMethod != item.EvidenceMethod {
				item.Revision++
			}
		}
		items = append(items, item)
	}
	newCriteria := current
	criteriaResourceVersion, err := incrementResourceVersion(current.Metadata.ResourceVersion)
	if err != nil {
		return GoalReference{}, err
	}
	newCriteria.Metadata.ResourceVersion = criteriaResourceVersion
	newCriteria.Metadata.Generation++
	newCriteria.Spec.Revision = nextRevision
	newCriteria.Spec.Items = items
	newGoal := goal
	goalResourceVersion, err := incrementResourceVersion(goal.Metadata.ResourceVersion)
	if err != nil {
		return GoalReference{}, err
	}
	newGoal.Metadata.ResourceVersion = goalResourceVersion
	newGoal.Status.CriteriaRevision = newCriteria.Spec.Revision
	newGoal.Status.LifecycleGeneration++
	newGoal.Status.ObservedAt = s.timestamp()
	if err := newCriteria.validate(newGoal); err != nil {
		return GoalReference{}, fmt.Errorf("invalid desired criteria: %w", err)
	}
	if err := newGoal.validate(); err != nil {
		return GoalReference{}, err
	}
	newCriteriaBytes, err := marshalYAML(newCriteria)
	if err != nil {
		return GoalReference{}, err
	}
	if err := validateProspectiveRecord(newGoal, newCriteria, attempts); err != nil {
		return GoalReference{}, err
	}
	goalBytes, err := marshalYAML(newGoal)
	if err != nil {
		return GoalReference{}, err
	}
	// goal.yaml is the optimistic-concurrency commit point and must publish
	// first; criteria.yaml and the immutable criteria-revisions/<n>.yaml
	// snapshot follow so a partial publication fails closed and recovery can
	// replay every after-image in order.
	entries := []publicationFileEntry{
		{
			Path:         "goal.yaml",
			BeforeDigest: "",
			Content:      goalBytes,
		},
		{
			Path:         "criteria.yaml",
			BeforeDigest: "",
			Content:      newCriteriaBytes,
		},
		{
			Path:         "criteria-revisions/" + strconv.FormatUint(newCriteria.Spec.Revision, 10) + ".yaml",
			BeforeDigest: "",
			Content:      newCriteriaBytes,
		},
	}
	// criteria-revisions directory must exist.
	if err := os.MkdirAll(filepath.Join(dir, "criteria-revisions"), 0o755); err != nil {
		return GoalReference{}, err
	}
	intent, err := s.beginPublication(
		dir,
		newGoal.Metadata.Name,
		goal.Metadata.ResourceVersion,
		newGoal.Metadata.ResourceVersion,
		entries,
		nil,
	)
	if err != nil {
		return GoalReference{}, err
	}
	reference := GoalReference{
		GoalID:          newGoal.Metadata.Name,
		GoalRef:         ".",
		ResourceVersion: newGoal.Metadata.ResourceVersion,
	}
	if err := s.publishIntentFiles(dir, intent, ""); err != nil {
		if !s.commitPointPublished(dir, intent) {
			_ = s.finishPublication(dir)
			return GoalReference{}, err
		}
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "publish",
			Kind:             "criteria",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	if err := s.refreshREADMEProjection(dir); err != nil {
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "projection",
			Kind:             "criteria",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	if err := s.finishPublication(dir); err != nil {
		return reference, &PublicationIncompleteError{
			OperationID:      intent.Spec.OperationID,
			IntendedRevision: reference.ResourceVersion,
			Phase:            "finish",
			Kind:             "criteria",
			Message:          err.Error(),
			Cause:            err,
		}
	}
	return reference, nil
}

func (s *Store) ValidateGoal(goalDir string) error {
	dir, err := s.resolveInsideWorkspace(goalDir)
	if err != nil {
		return err
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return err
	}
	defer lock.release()
	if err := s.checkNoPendingPublication(dir); err != nil {
		return err
	}
	_, _, _, err = s.loadAndValidate(dir)
	return err
}

func (s *Store) ValidateRoot(goalsRoot string) (int, error) {
	root, err := s.resolveInsideWorkspace(goalsRoot)
	if err != nil {
		return 0, err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, err
	}
	validated := 0
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		dir := filepath.Join(root, entry.Name())
		lock, err := s.acquireGoalLock(dir)
		if err != nil {
			return validated, err
		}
		_, _, _, validateErr := s.loadAndValidate(dir)
		lock.release()
		if validateErr != nil {
			return validated, fmt.Errorf("goal %q: %w", entry.Name(), validateErr)
		}
		validated++
	}
	return validated, nil
}

func (s *Store) loadAndValidate(dir string) (GoalManifest, CriteriaManifest, []AttemptManifest, error) {
	if err := cleanupTemporaryResidue(dir); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	if err := rejectRecordSymlinks(dir); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	goal, err := s.readGoalManifest(dir)
	if err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	if filepath.Base(dir) != goal.Metadata.Name {
		return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("directory name does not match Goal metadata.name")
	}
	var criteria CriteriaManifest
	if err := s.readYAML(filepath.Join(dir, "criteria.yaml"), &criteria); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	if err := criteria.validate(goal); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	history, err := s.loadCriteriaHistory(dir, goal, criteria)
	if err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	if _, err := os.Stat(filepath.Join(dir, "criteria-revisions")); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("missing criteria-revisions")
	}
	entries, err := os.ReadDir(filepath.Join(dir, "attempts"))
	if os.IsNotExist(err) {
		entries = nil
	} else if err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	if len(entries) > maxAttempts {
		return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("attempt cardinality exceeds %d", maxAttempts)
	}
	attempts := make([]AttemptManifest, 0, len(entries))
	openAttempt := ""
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("unexpected entry in attempts: %s", entry.Name())
		}
		attemptDir := filepath.Join(dir, "attempts", entry.Name())
		var attempt AttemptManifest
		if err := s.readYAML(filepath.Join(attemptDir, "attempt.yaml"), &attempt); err != nil {
			return GoalManifest{}, CriteriaManifest{}, nil, err
		}
		if attempt.Metadata.Name != entry.Name() {
			return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("attempt directory name mismatch")
		}
		if err := attempt.validate(goal); err != nil {
			return GoalManifest{}, CriteriaManifest{}, nil, err
		}
		if err := validateAttemptFiles(attemptDir); err != nil {
			return GoalManifest{}, CriteriaManifest{}, nil, err
		}
		if err := validateAttemptArtifacts(attemptDir, attempt.Status.Artifacts); err != nil {
			return GoalManifest{}, CriteriaManifest{}, nil, err
		}
		snapshot, ok := history.Snapshots[attempt.Spec.CriteriaRevision]
		if !ok {
			return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf(
				"attempt %q refers to a missing criteria snapshot",
				attempt.Metadata.Name,
			)
		}
		criteriaDigest, err := criteriaPortableDigest(snapshot)
		if err != nil || attempt.Spec.CriteriaDigest != criteriaDigest {
			return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf(
				"attempt %q criteria digest does not match immutable snapshot",
				attempt.Metadata.Name,
			)
		}
		if attempt.Status.State == "closed" {
			if err := validateReviewAgainstCriteria(attempt.Status.Review, snapshot); err != nil {
				return GoalManifest{}, CriteriaManifest{}, nil, err
			}
		}
		if attempt.Status.State == "open" {
			if openAttempt != "" {
				return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("multiple open attempts")
			}
			openAttempt = attempt.Metadata.Name
		}
		attempts = append(attempts, attempt)
	}
	if openAttempt != goal.Status.ActiveAttemptID {
		return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("active attempt does not match open attempt set")
	}
	if openAttempt != "" {
		if goal.Status.Outcome != "open" || goal.Status.Execution != "active" {
			return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("open attempt requires open active goal")
		}
		for _, attempt := range attempts {
			if attempt.Metadata.Name == openAttempt && attempt.Spec.LifecycleGeneration != goal.Status.LifecycleGeneration {
				return GoalManifest{}, CriteriaManifest{}, nil, fmt.Errorf("active attempt lifecycle generation is stale")
			}
		}
	}
	if err := validateProspectiveRecord(goal, criteria, attempts); err != nil {
		return GoalManifest{}, CriteriaManifest{}, nil, err
	}
	return goal, criteria, attempts, nil
}

func (s *Store) readGoalManifest(dir string) (GoalManifest, error) {
	var goal GoalManifest
	if err := s.readYAML(filepath.Join(dir, "goal.yaml"), &goal); err != nil {
		return GoalManifest{}, err
	}
	if err := goal.validate(); err != nil {
		return GoalManifest{}, fmt.Errorf("invalid goal.yaml: %w", err)
	}
	return goal, nil
}

func (s *Store) readYAML(path string, destination any) error {
	content, err := readRegularFile(path, maxManifestBytes)
	if err != nil {
		return fmt.Errorf("read %s: %w", filepath.Base(path), err)
	}
	if err := yaml.UnmarshalWithOptions(content, destination, yaml.Strict()); err != nil {
		return fmt.Errorf("decode %s strictly: %w", filepath.Base(path), err)
	}
	return nil
}

func (s *Store) writeYAML(path string, value any) error {
	content, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode %s: %w", filepath.Base(path), err)
	}
	if len(content) > maxManifestBytes {
		return fmt.Errorf("encoded manifest exceeds %d bytes", maxManifestBytes)
	}
	return s.atomicWrite(path, content, 0o644)
}

func (s *Store) atomicWrite(path string, content []byte, mode os.FileMode) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".goal-write-")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	installed := false
	defer func() {
		_ = temporary.Close()
		if !installed {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(mode); err != nil {
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := s.callBeforeRename(path); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("atomic rename %s: %w", filepath.Base(path), err)
	}
	installed = true
	return nil
}

func (s *Store) callBeforeRename(path string) error {
	if s.beforeRename != nil {
		if err := s.beforeRename(path); err != nil {
			return fmt.Errorf("injected failure before rename of %s: %w", filepath.Base(path), err)
		}
	}
	return nil
}

func (s *Store) timestamp() string {
	return s.now().UTC().Format(time.RFC3339Nano)
}

func (s *Store) generateGoalID(title string) (string, error) {
	base := slugify(title)
	if base == "" {
		base = "goal"
	}
	if len(base) > 54 {
		base = strings.Trim(base[:54], "-._")
	}
	suffix := make([]byte, 4)
	if _, err := io.ReadFull(s.random, suffix); err != nil {
		return "", fmt.Errorf("generate goal ID: %w", err)
	}
	return base + "-" + hex.EncodeToString(suffix), nil
}

func (s *Store) generateAttemptID() (string, error) {
	suffix := make([]byte, 6)
	if _, err := io.ReadFull(s.random, suffix); err != nil {
		return "", fmt.Errorf("generate attempt ID: %w", err)
	}
	return "attempt-" + hex.EncodeToString(suffix), nil
}

func slugify(value string) string {
	var builder strings.Builder
	previousDash := false
	for _, char := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			if char <= unicode.MaxASCII {
				builder.WriteRune(char)
				previousDash = false
			}
		} else if !previousDash && builder.Len() > 0 {
			builder.WriteByte('-')
			previousDash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func (s *Store) resolveInsideWorkspace(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("path is required")
	}
	if !filepath.IsAbs(value) {
		value = filepath.Join(s.workspaceRoot, value)
	}
	resolved, err := resolveExistingSymlinks(value)
	if err != nil {
		return "", err
	}
	if !isWithin(s.workspaceRoot, resolved) {
		return "", fmt.Errorf("path escapes workspace root")
	}
	return resolved, nil
}

func resolveExistingSymlinks(value string) (string, error) {
	absolute, err := absoluteClean(value)
	if err != nil {
		return "", err
	}
	probe := absolute
	var missing []string
	for {
		_, err := os.Lstat(probe)
		if err == nil {
			break
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(probe)
		if parent == probe {
			return "", fmt.Errorf("cannot resolve an existing path ancestor")
		}
		missing = append(missing, filepath.Base(probe))
		probe = parent
	}
	resolved, err := filepath.EvalSymlinks(probe)
	if err != nil {
		return "", err
	}
	for index := len(missing) - 1; index >= 0; index-- {
		resolved = filepath.Join(resolved, missing[index])
	}
	return filepath.Clean(resolved), nil
}

func (s *Store) resolveSessionRoot(value string) (string, error) {
	resolved, err := s.resolveInsideWorkspace(value)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(s.workspaceRoot, resolved)
	if err != nil {
		return "", err
	}
	parts := strings.Split(filepath.ToSlash(relative), "/")
	if len(parts) < 3 || parts[0] != "out" || parts[1] == "" || parts[1] == "." {
		return "", fmt.Errorf("session root must be task-specific under out/<task>/")
	}
	return resolved, nil
}

func portableRelative(base string, target string, allowParent bool) (string, error) {
	relative, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	portable := filepath.ToSlash(relative)
	if err := validatePortablePath(
		"portable reference",
		portable,
		allowParent,
		false,
	); err != nil {
		return "", err
	}
	return portable, nil
}

func portableOwnerRoot(workspaceRoot string, ownerRoot string) (string, error) {
	relative, err := filepath.Rel(workspaceRoot, ownerRoot)
	if err != nil {
		return "", err
	}
	portable := filepath.ToSlash(relative)
	if err := validatePortablePath(
		"local owner root annotation",
		portable,
		false,
		true,
	); err != nil {
		return "", err
	}
	return portable, nil
}

func absoluteClean(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("path is empty")
	}
	return filepath.Abs(filepath.Clean(value))
}

func isWithin(root string, candidate string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func summaryFromGoal(ref string, goal GoalManifest) GoalSummary {
	return GoalSummary{
		GoalID:              goal.Metadata.Name,
		GoalRef:             filepath.ToSlash(ref),
		Title:               goal.Spec.Title,
		Scope:               goal.Spec.Scope,
		ResourceVersion:     goal.Metadata.ResourceVersion,
		Generation:          goal.Metadata.Generation,
		LifecycleGeneration: goal.Status.LifecycleGeneration,
		Outcome:             goal.Status.Outcome,
		Execution:           goal.Status.Execution,
	}
}

func readRegularFile(path string, maximum int64) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maximum {
		return nil, fmt.Errorf("must be a regular file of at most %d bytes", maximum)
	}
	return os.ReadFile(path)
}

func (s *Store) readWorkspaceRegularFile(path string, maximum int64) ([]byte, error) {
	resolved, err := s.resolveInsideWorkspace(path)
	if err != nil {
		return nil, err
	}
	return readRegularFile(resolved, maximum)
}

func (s *Store) readWorkspaceMarkdownFile(path string, maximum int64) ([]byte, error) {
	resolved, err := s.resolveInsideWorkspace(path)
	if err != nil {
		return nil, err
	}
	return readMarkdownFile(resolved, maximum)
}

func rejectRecordSymlinks(root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("record contains symlink at %s", filepath.ToSlash(strings.TrimPrefix(path, root+string(filepath.Separator))))
		}
		return nil
	})
}

func validateAttemptFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read attempt directory: %w", err)
	}
	allowed := map[string]bool{
		"attempt.yaml": false,
		"evidence":     true,
		"plan.md":      false,
		"result.md":    false,
	}
	if len(entries) != len(allowed) {
		return fmt.Errorf(
			"attempt %s must contain only attempt.yaml, plan.md, result.md, and evidence/",
			filepath.Base(dir),
		)
	}
	for _, entry := range entries {
		wantDirectory, ok := allowed[entry.Name()]
		if !ok || entry.IsDir() != wantDirectory {
			return fmt.Errorf("unexpected attempt entry %q", entry.Name())
		}
	}
	for _, name := range []string{"plan.md", "result.md"} {
		if _, err := readMarkdownFile(filepath.Join(dir, name), maxPlanResultBytes); err != nil {
			return fmt.Errorf("invalid %s for attempt %s: %w", name, filepath.Base(dir), err)
		}
	}
	evidenceEntries, err := os.ReadDir(filepath.Join(dir, "evidence"))
	if err != nil {
		return fmt.Errorf("read evidence directory: %w", err)
	}
	if len(evidenceEntries) > maxEvidenceFiles {
		return fmt.Errorf("evidence cardinality exceeds %d", maxEvidenceFiles)
	}
	for _, entry := range evidenceEntries {
		if entry.IsDir() || !safeEvidenceName(entry.Name()) {
			return fmt.Errorf("invalid evidence entry %q", entry.Name())
		}
		if _, err := readMarkdownFile(filepath.Join(dir, "evidence", entry.Name()), maxEvidenceFileBytes); err != nil {
			return err
		}
	}
	return nil
}

func safeEvidenceName(value string) bool {
	if value == "" || len(value) > 128 || value == "." || value == ".." ||
		!strings.HasSuffix(value, ".md") {
		return false
	}
	for _, char := range value {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') && char != '.' && char != '_' && char != '-' {
			return false
		}
	}
	return true
}

type heldLock struct {
	file *os.File
}

func (s *Store) acquireGoalLock(path string) (*heldLock, error) {
	return s.acquirePathLock(path)
}

func (s *Store) acquirePathLock(path string) (*heldLock, error) {
	lockPath, err := s.pathLockPath(path)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(
		lockPath,
		os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW|syscall.O_NONBLOCK,
		0o600,
	)
	if err != nil {
		return nil, fmt.Errorf("open goal lock: %w", err)
	}
	if err := validateGoalLockFile(file); err != nil {
		_ = file.Close()
		return nil, err
	}
	if err := file.Chmod(0o600); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("set goal lock permissions: %w", err)
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("acquire goal lock: %w", err)
	}
	if err := writeLockHolder(file); err != nil {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
		return nil, err
	}
	return &heldLock{file: file}, nil
}

func (s *Store) pathLockPath(path string) (string, error) {
	canonical, err := resolveExistingSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("resolve canonical path for lock: %w", err)
	}
	if !filepath.IsAbs(canonical) {
		return "", fmt.Errorf("canonical path for lock is not absolute")
	}
	root, err := s.ensureGoalLockRoot()
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(canonical))
	return filepath.Join(root, hex.EncodeToString(digest[:])+".lock"), nil
}

func (s *Store) ensureGoalLockRoot() (string, error) {
	if s.runtimeRoot == "" || !filepath.IsAbs(s.runtimeRoot) {
		return "", fmt.Errorf("valid XDG runtime root is required for goal locks")
	}
	lockRoot := goalLockRootPath(s.runtimeRoot)
	if isWithin(s.workspaceRoot, lockRoot) {
		return "", fmt.Errorf(
			"XDG_RUNTIME_DIR/alwaldend/goal/locks must be outside the workspace",
		)
	}
	current := s.runtimeRoot
	for _, component := range []string{"alwaldend", "goal", "locks"} {
		current = filepath.Join(current, component)
		if err := os.Mkdir(current, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
			return "", fmt.Errorf("create goal lock directory: %w", err)
		}
		info, err := os.Lstat(current)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("goal lock path is not a directory: %s", current)
		}
		if err := os.Chmod(current, 0o700); err != nil {
			return "", fmt.Errorf("set goal lock directory permissions: %w", err)
		}
	}
	resolvedLockRoot, err := filepath.EvalSymlinks(current)
	if err != nil {
		return "", fmt.Errorf("resolve goal lock directory: %w", err)
	}
	if isWithin(s.workspaceRoot, resolvedLockRoot) {
		return "", fmt.Errorf("goal lock directory must be outside the workspace")
	}
	return filepath.Clean(resolvedLockRoot), nil
}

func goalLockRootPath(runtimeRoot string) string {
	return filepath.Join(runtimeRoot, "alwaldend", "goal", "locks")
}

func validateGoalLockFile(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("inspect goal lock: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("goal lock must be a regular file")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("goal lock ownership is unavailable")
	}
	if int(stat.Uid) != os.Getuid() {
		return fmt.Errorf("goal lock must be owned by the current user")
	}
	if stat.Nlink != 1 {
		return fmt.Errorf("goal lock must have exactly one link")
	}
	return nil
}

func writeLockHolder(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("truncate goal lock metadata: %w", err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek goal lock metadata: %w", err)
	}
	if _, err := fmt.Fprintf(file, "%d\n", os.Getpid()); err != nil {
		return fmt.Errorf("write goal lock metadata: %w", err)
	}
	return nil
}

func (lock *heldLock) release() error {
	clearErr := lock.file.Truncate(0)
	unlockErr := syscall.Flock(int(lock.file.Fd()), syscall.LOCK_UN)
	closeErr := lock.file.Close()
	return errors.Join(clearErr, unlockErr, closeErr)
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func validateLimit(limit int) (int, error) {
	if limit == 0 {
		limit = defaultOutputLimit
	}
	if limit < 1 || limit > maximumOutputLimit {
		return 0, fmt.Errorf("limit must be between 1 and %d", maximumOutputLimit)
	}
	return limit, nil
}

func digestBytes(value []byte) string {
	digest := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(digest[:])
}
