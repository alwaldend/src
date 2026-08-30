package fsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

type attemptTree struct {
	Manifest AttemptManifest
	Plan     []byte
	Result   []byte
	Evidence map[string][]byte
}

type stagedAttempt struct {
	temporary string
	target    string
}

func (s *Store) Checkpoint(options CheckpointOptions) (GoalReference, error) {
	dir, err := s.resolveInsideWorkspace(options.GoalDir)
	if err != nil {
		return GoalReference{}, err
	}
	if _, err := parseResourceVersion(options.ExpectedResourceVersion); err != nil {
		return GoalReference{}, fmt.Errorf("expected resource version: %w", err)
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return GoalReference{}, err
	}
	defer lock.release()
	goal, criteria, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return GoalReference{}, err
	}
	if goal.Metadata.ResourceVersion != options.ExpectedResourceVersion {
		return GoalReference{}, fmt.Errorf(
			"stale resourceVersion: expected %s, current is %s",
			options.ExpectedResourceVersion,
			goal.Metadata.ResourceVersion,
		)
	}
	if goal.Status.Outcome != "open" {
		return GoalReference{}, fmt.Errorf("closed goal cannot be checkpointed")
	}
	if options.CriteriaFile != "" {
		if options.AttemptID != "" || options.WorkType != "" || options.PlanFile != "" ||
			options.ResultFile != "" || len(options.EvidenceFiles) > 0 ||
			options.ReviewFile != "" || options.CloseAttempt || options.Outcome != "" ||
			options.Execution != "" {
			return GoalReference{}, fmt.Errorf(
				"--criteria-file cannot be combined with another checkpoint mutation",
			)
		}
		if goal.Status.Execution != "paused" || goal.Status.ActiveAttemptID != "" {
			return GoalReference{}, fmt.Errorf(
				"criteria update requires paused execution and no active attempt",
			)
		}
		return s.updateCriteria(dir, goal, options.CriteriaFile)
	}
	return s.checkpointRecord(dir, goal, criteria, attempts, options)
}

func (s *Store) checkpointRecord(
	dir string,
	goal GoalManifest,
	criteria CriteriaManifest,
	attempts []AttemptManifest,
	options CheckpointOptions,
) (GoalReference, error) {
	hasAttemptMutation := options.AttemptID != "" || options.WorkType != "" ||
		options.PlanFile != "" || options.ResultFile != "" ||
		len(options.EvidenceFiles) > 0 || options.ReviewFile != "" || options.CloseAttempt
	if !hasAttemptMutation && options.Outcome == "" && options.Execution == "" {
		return GoalReference{}, fmt.Errorf("checkpoint has no mutation")
	}
	if hasAttemptMutation && goal.Status.Execution != "active" {
		return GoalReference{}, fmt.Errorf("attempt publication requires active execution")
	}
	if options.Outcome != "" &&
		!oneOf(options.Outcome, "open", "achieved", "abandoned", "superseded") {
		return GoalReference{}, fmt.Errorf("invalid outcome %q", options.Outcome)
	}
	if options.Execution != "" &&
		!oneOf(options.Execution, "active", "paused", "waiting", "blocked") {
		return GoalReference{}, fmt.Errorf("invalid execution %q", options.Execution)
	}
	if len(options.EvidenceFiles) > 64 {
		return GoalReference{}, fmt.Errorf("one checkpoint accepts at most 64 evidence files")
	}
	if options.ReviewFile != "" && !options.CloseAttempt {
		return GoalReference{}, fmt.Errorf("--review-file requires --close-attempt")
	}
	if options.CloseAttempt && options.ReviewFile == "" {
		return GoalReference{}, fmt.Errorf("closing an attempt requires --review-file")
	}

	attemptByID := make(map[string]AttemptManifest, len(attempts))
	for _, attempt := range attempts {
		attemptByID[attempt.Metadata.Name] = attempt
	}
	attemptID := options.AttemptID
	if attemptID == "" && hasAttemptMutation {
		attemptID = goal.Status.ActiveAttemptID
	}
	if attemptID == "" && hasAttemptMutation {
		var err error
		attemptID, err = s.generateAttemptID()
		if err != nil {
			return GoalReference{}, err
		}
	}

	newGoal := goal
	prospectiveAttempts := append([]AttemptManifest(nil), attempts...)
	newAttempt := false
	var tree *attemptTree
	var err error
	if attemptID != "" {
		if err := validateRecordID("attempt ID", attemptID); err != nil {
			return GoalReference{}, err
		}
		if goal.Status.ActiveAttemptID != "" && goal.Status.ActiveAttemptID != attemptID {
			return GoalReference{}, fmt.Errorf(
				"another attempt %q is active",
				goal.Status.ActiveAttemptID,
			)
		}
		existing, exists := attemptByID[attemptID]
		if exists {
			if existing.Status.State != "open" {
				return GoalReference{}, fmt.Errorf("attempt %q is closed and immutable", attemptID)
			}
			if existing.Spec.LifecycleGeneration != goal.Status.LifecycleGeneration {
				return GoalReference{}, fmt.Errorf(
					"attempt %q belongs to an obsolete lifecycle generation",
					attemptID,
				)
			}
			if options.PlanFile != "" {
				return GoalReference{}, fmt.Errorf("an existing attempt plan is immutable")
			}
			tree, err = s.buildAttemptTree(dir, goal, criteria, &existing, attemptID, options)
			if err != nil {
				return GoalReference{}, err
			}
			for index := range prospectiveAttempts {
				if prospectiveAttempts[index].Metadata.Name == attemptID {
					prospectiveAttempts[index] = tree.Manifest
				}
			}
		} else {
			if len(attempts) >= maxAttempts {
				return GoalReference{}, fmt.Errorf("attempt cardinality exceeds %d", maxAttempts)
			}
			newAttempt = true
			tree, err = s.buildAttemptTree(dir, goal, criteria, nil, attemptID, options)
			if err != nil {
				return GoalReference{}, err
			}
			prospectiveAttempts = append(prospectiveAttempts, tree.Manifest)
		}
		if options.CloseAttempt {
			newGoal.Status.ActiveAttemptID = ""
		} else {
			newGoal.Status.ActiveAttemptID = attemptID
		}
	}

	newOutcome := goal.Status.Outcome
	if options.Outcome != "" {
		newOutcome = options.Outcome
	}
	newExecution := goal.Status.Execution
	if options.Execution != "" {
		newExecution = options.Execution
	}
	if newOutcome != "open" {
		if newGoal.Status.ActiveAttemptID != "" {
			return GoalReference{}, fmt.Errorf("close the active attempt before closing the goal")
		}
		newExecution = "paused"
	}
	if newExecution != "active" && newGoal.Status.ActiveAttemptID != "" {
		return GoalReference{}, fmt.Errorf("close the active attempt before leaving active execution")
	}
	if newOutcome == "achieved" {
		if !options.CloseAttempt || tree == nil {
			return GoalReference{}, fmt.Errorf(
				"achieved must close an accept attempt in the same checkpoint",
			)
		}
		criteriaDigest, err := criteriaPortableDigest(criteria)
		if err != nil {
			return GoalReference{}, err
		}
		if tree.Manifest.Spec.CriteriaRevision != criteria.Spec.Revision ||
			tree.Manifest.Spec.CriteriaDigest != criteriaDigest {
			return GoalReference{}, fmt.Errorf("achieved attempt is not bound to current criteria")
		}
		if err := reviewAcceptsRequired(tree.Manifest.Status.Review, criteria); err != nil {
			return GoalReference{}, err
		}
		newGoal.Status.AcceptedAttemptID = tree.Manifest.Metadata.Name
		newGoal.Status.AcceptedResultDigest = tree.Manifest.Status.Artifacts.ResultDigest
	}
	if newOutcome != goal.Status.Outcome || newExecution != goal.Status.Execution {
		newGoal.Status.LifecycleGeneration++
	}
	newGoal.Status.Outcome = newOutcome
	newGoal.Status.Execution = newExecution
	next, err := incrementResourceVersion(goal.Metadata.ResourceVersion)
	if err != nil {
		return GoalReference{}, err
	}
	newGoal.Metadata.ResourceVersion = next
	newGoal.Status.ObservedAt = s.timestamp()
	if err := validateProspectiveRecord(newGoal, criteria, prospectiveAttempts); err != nil {
		return GoalReference{}, err
	}
	readme, err := renderREADME(newGoal, criteria, prospectiveAttempts, defaultOutputLimit)
	if err != nil {
		return GoalReference{}, err
	}
	var staged *stagedAttempt
	if tree != nil && newAttempt {
		staged, err = s.stageNewAttempt(dir, attemptID, *tree)
		if err != nil {
			return GoalReference{}, err
		}
		defer func() {
			_ = os.RemoveAll(staged.temporary)
		}()
	}
	committedGoal := newGoal
	finalizeGoalAfterAttempt := tree != nil && newAttempt && options.CloseAttempt
	if finalizeGoalAfterAttempt {
		// A newly created closed attempt otherwise has no pointer from the goal.
		// Keep a deliberately inconsistent active pointer until the staged
		// directory is published. If either following rename is interrupted,
		// validation fails closed instead of silently accepting a goal that lost
		// the closed attempt.
		committedGoal.Status.ActiveAttemptID = attemptID
	}
	// goal.yaml is the optimistic-concurrency commit point. Stage a new
	// attempt first, but publish the advanced token before any canonical
	// attempt content. A later publication failure therefore always leaves the
	// record guarded by the committed resourceVersion rather than the caller's
	// stale token.
	if err := s.writeYAML(filepath.Join(dir, "goal.yaml"), committedGoal); err != nil {
		return GoalReference{}, err
	}
	reference := GoalReference{
		GoalID:          newGoal.Metadata.Name,
		GoalRef:         ".",
		ResourceVersion: newGoal.Metadata.ResourceVersion,
	}
	if tree != nil {
		if newAttempt {
			err = s.publishStagedAttempt(staged)
		} else {
			err = s.updateExistingAttempt(dir, attemptID, *tree, options)
		}
		if err != nil {
			return reference, fmt.Errorf(
				"checkpoint committed at resourceVersion %s; publish attempt %q: %w",
				reference.ResourceVersion,
				attemptID,
				err,
			)
		}
	}
	if finalizeGoalAfterAttempt {
		if err := s.writeYAML(filepath.Join(dir, "goal.yaml"), newGoal); err != nil {
			return reference, fmt.Errorf(
				"checkpoint committed at resourceVersion %s; finalize goal after attempt %q: %w",
				reference.ResourceVersion,
				attemptID,
				err,
			)
		}
	}
	if err := s.atomicWrite(filepath.Join(dir, "README.md"), readme, 0o644); err != nil {
		return reference, fmt.Errorf(
			"checkpoint committed at resourceVersion %s; refresh README projection: %w",
			reference.ResourceVersion,
			err,
		)
	}
	return reference, nil
}

func (s *Store) buildAttemptTree(
	dir string,
	goal GoalManifest,
	criteria CriteriaManifest,
	existing *AttemptManifest,
	attemptID string,
	options CheckpointOptions,
) (*attemptTree, error) {
	now := s.timestamp()
	tree := &attemptTree{Evidence: map[string][]byte{}}
	if existing == nil {
		workType := options.WorkType
		if workType == "" {
			workType = "change"
		}
		criteriaDigest, err := criteriaPortableDigest(criteria)
		if err != nil {
			return nil, err
		}
		stateDigest, err := goalStateDigest(goal)
		if err != nil {
			return nil, err
		}
		tree.Manifest = AttemptManifest{
			APIVersion: goalAPIVersion,
			Kind:       "GoalAttempt",
			Metadata: ObjectMeta{
				Name:              attemptID,
				ResourceVersion:   "1",
				Generation:        1,
				CreationTimestamp: now,
			},
			Spec: AttemptSpec{
				GoalRef:             LocalGoalReference{Name: goal.Metadata.Name},
				GoalGeneration:      goal.Metadata.Generation,
				LifecycleGeneration: goal.Status.LifecycleGeneration,
				CriteriaRevision:    criteria.Spec.Revision,
				CriteriaDigest:      criteriaDigest,
				GoalStateDigest:     stateDigest,
				WorkType:            workType,
			},
			Status: AttemptStatus{State: "open", ObservedAt: now},
		}
		tree.Plan = []byte("# Plan\n\nNo plan was supplied.\n")
		tree.Result = []byte("# Result\n\nNo result has been published.\n")
	} else {
		if options.WorkType != "" && options.WorkType != existing.Spec.WorkType {
			return nil, fmt.Errorf("an existing attempt workType is immutable")
		}
		tree.Manifest = *existing
		attemptDir := filepath.Join(dir, "attempts", attemptID)
		var err error
		tree.Plan, err = readMarkdownFile(filepath.Join(attemptDir, "plan.md"), maxPlanResultBytes)
		if err != nil {
			return nil, err
		}
		tree.Result, err = readMarkdownFile(filepath.Join(attemptDir, "result.md"), maxPlanResultBytes)
		if err != nil {
			return nil, err
		}
		entries, err := os.ReadDir(filepath.Join(attemptDir, "evidence"))
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".md") {
				return nil, fmt.Errorf(
					"stored evidence %q is not a Markdown file",
					entry.Name(),
				)
			}
			content, err := readMarkdownFile(
				filepath.Join(attemptDir, "evidence", entry.Name()),
				maxEvidenceFileBytes,
			)
			if err != nil {
				return nil, err
			}
			tree.Evidence[entry.Name()] = content
		}
	}
	var err error
	if options.PlanFile != "" {
		tree.Plan, err = readMarkdownFile(options.PlanFile, maxPlanResultBytes)
		if err != nil {
			return nil, fmt.Errorf("read plan: %w", err)
		}
	}
	if options.ResultFile != "" {
		tree.Result, err = readMarkdownFile(options.ResultFile, maxPlanResultBytes)
		if err != nil {
			return nil, fmt.Errorf("read result: %w", err)
		}
	}
	newEvidenceNames := map[string]bool{}
	for _, source := range options.EvidenceFiles {
		name := filepath.Base(source)
		if !safeEvidenceName(name) || !strings.HasSuffix(name, ".md") ||
			newEvidenceNames[name] {
			return nil, fmt.Errorf(
				"evidence basename %q must be a safe, unique Markdown filename",
				name,
			)
		}
		newEvidenceNames[name] = true
		if _, exists := tree.Evidence[name]; exists {
			return nil, fmt.Errorf("evidence %q already exists and is immutable", name)
		}
		content, err := readMarkdownFile(source, maxEvidenceFileBytes)
		if err != nil {
			return nil, fmt.Errorf("read evidence %q: %w", name, err)
		}
		tree.Evidence[name] = content
	}
	if len(tree.Evidence) > maxEvidenceFiles {
		return nil, fmt.Errorf("evidence cardinality exceeds %d", maxEvidenceFiles)
	}
	tree.Manifest.Status.Artifacts = artifactManifestForTree(*tree)
	if options.CloseAttempt {
		reviewContent, err := readRegularFile(options.ReviewFile, maxManifestBytes)
		if err != nil {
			return nil, fmt.Errorf("read close review: %w", err)
		}
		var review CloseReview
		if err := yaml.UnmarshalWithOptions(reviewContent, &review, yaml.Strict()); err != nil {
			return nil, fmt.Errorf("decode close review strictly: %w", err)
		}
		tree.Manifest.Status.State = "closed"
		tree.Manifest.Status.ClosedAt = now
		tree.Manifest.Status.Review = review
		if err := validateReviewAgainstCriteria(review, criteria); err != nil {
			return nil, err
		}
	}
	if existing != nil {
		next, err := incrementResourceVersion(existing.Metadata.ResourceVersion)
		if err != nil {
			return nil, err
		}
		tree.Manifest.Metadata.ResourceVersion = next
	}
	tree.Manifest.Status.ObservedAt = now
	if err := tree.Manifest.validate(goal); err != nil {
		return nil, err
	}
	return tree, nil
}

func artifactManifestForTree(tree attemptTree) ArtifactManifest {
	manifest := ArtifactManifest{
		PlanDigest:   digestBytes(tree.Plan),
		ResultDigest: digestBytes(tree.Result),
		Evidence:     make([]ArtifactDigest, 0, len(tree.Evidence)),
	}
	for name, content := range tree.Evidence {
		manifest.Evidence = append(manifest.Evidence, ArtifactDigest{
			Path: "evidence/" + name, Digest: digestBytes(content),
		})
	}
	sort.Slice(manifest.Evidence, func(i, j int) bool {
		return manifest.Evidence[i].Path < manifest.Evidence[j].Path
	})
	return manifest
}

func (s *Store) writeAttemptTree(root string, tree attemptTree) error {
	if err := os.Mkdir(filepath.Join(root, "evidence"), 0o755); err != nil {
		return err
	}
	if err := s.atomicWrite(filepath.Join(root, "plan.md"), tree.Plan, 0o644); err != nil {
		return err
	}
	if err := s.atomicWrite(filepath.Join(root, "result.md"), tree.Result, 0o644); err != nil {
		return err
	}
	names := make([]string, 0, len(tree.Evidence))
	for name := range tree.Evidence {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if !strings.HasSuffix(name, ".md") {
			return fmt.Errorf("evidence %q is not a Markdown file", name)
		}
		if err := s.atomicWrite(
			filepath.Join(root, "evidence", name),
			tree.Evidence[name],
			0o644,
		); err != nil {
			return err
		}
	}
	return s.writeYAML(filepath.Join(root, "attempt.yaml"), tree.Manifest)
}

func (s *Store) stageNewAttempt(
	dir string,
	attemptID string,
	tree attemptTree,
) (*stagedAttempt, error) {
	attemptsDir := filepath.Join(dir, "attempts")
	temporary, err := os.MkdirTemp(attemptsDir, ".goal-attempt-")
	if err != nil {
		return nil, fmt.Errorf("create temporary attempt directory: %w", err)
	}
	staged := false
	defer func() {
		if !staged {
			_ = os.RemoveAll(temporary)
		}
	}()
	if err := s.writeAttemptTree(temporary, tree); err != nil {
		return nil, err
	}
	target := filepath.Join(attemptsDir, attemptID)
	if pathExists(target) {
		return nil, fmt.Errorf("attempt %q already exists", attemptID)
	}
	staged = true
	return &stagedAttempt{temporary: temporary, target: target}, nil
}

func (s *Store) publishStagedAttempt(staged *stagedAttempt) error {
	if pathExists(staged.target) {
		return fmt.Errorf("attempt %q already exists", filepath.Base(staged.target))
	}
	if err := s.callBeforeRename(staged.target); err != nil {
		return err
	}
	if err := os.Rename(staged.temporary, staged.target); err != nil {
		return fmt.Errorf("install attempt %q: %w", filepath.Base(staged.target), err)
	}
	return nil
}

func (s *Store) updateExistingAttempt(
	dir string,
	attemptID string,
	tree attemptTree,
	options CheckpointOptions,
) error {
	attemptDir := filepath.Join(dir, "attempts", attemptID)
	if options.ResultFile != "" {
		if err := s.atomicWrite(
			filepath.Join(attemptDir, "result.md"),
			tree.Result,
			0o644,
		); err != nil {
			return err
		}
	}
	names := make([]string, 0, len(options.EvidenceFiles))
	for _, source := range options.EvidenceFiles {
		names = append(names, filepath.Base(source))
	}
	sort.Strings(names)
	for _, name := range names {
		target := filepath.Join(attemptDir, "evidence", name)
		if pathExists(target) {
			return fmt.Errorf("evidence %q already exists and is immutable", name)
		}
		if err := s.atomicWrite(target, tree.Evidence[name], 0o644); err != nil {
			return err
		}
	}
	return s.writeYAML(filepath.Join(attemptDir, "attempt.yaml"), tree.Manifest)
}

func validateProspectiveRecord(
	goal GoalManifest,
	criteria CriteriaManifest,
	attempts []AttemptManifest,
) error {
	if err := goal.validate(); err != nil {
		return err
	}
	if err := criteria.validate(goal); err != nil {
		return err
	}
	open := ""
	for _, attempt := range attempts {
		if err := attempt.validate(goal); err != nil {
			return err
		}
		if attempt.Status.State == "open" {
			if open != "" {
				return fmt.Errorf("multiple open attempts")
			}
			open = attempt.Metadata.Name
		}
	}
	if open != goal.Status.ActiveAttemptID {
		return fmt.Errorf("active attempt does not match open attempt set")
	}
	if open != "" && (goal.Status.Outcome != "open" || goal.Status.Execution != "active") {
		return fmt.Errorf("open attempt requires open active goal")
	}
	if goal.Status.Outcome == "achieved" {
		var accepted *AttemptManifest
		for index := range attempts {
			if attempts[index].Metadata.Name == goal.Status.AcceptedAttemptID {
				accepted = &attempts[index]
				break
			}
		}
		if accepted == nil || accepted.Status.State != "closed" ||
			accepted.Status.Review.Decision != "accept" ||
			accepted.Status.Artifacts.ResultDigest != goal.Status.AcceptedResultDigest ||
			accepted.Spec.CriteriaRevision != criteria.Spec.Revision {
			return fmt.Errorf("achieved acceptance pointer does not resolve to the exact accepted result")
		}
		criteriaDigest, err := criteriaPortableDigest(criteria)
		if err != nil {
			return err
		}
		if accepted.Spec.CriteriaDigest != criteriaDigest {
			return fmt.Errorf("accepted attempt is not bound to current portable criteria digest")
		}
		if err := reviewAcceptsRequired(accepted.Status.Review, criteria); err != nil {
			return err
		}
	}
	return nil
}
