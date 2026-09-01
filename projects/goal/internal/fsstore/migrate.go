package fsstore

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PromoteOptions struct {
	GoalDir                 string
	DestinationGoalsRoot    string
	ExpectedResourceVersion string
	OwnerRoot               string
}

func (s *Store) Promote(options PromoteOptions) (GoalReference, error) {
	source, err := s.resolveInsideWorkspace(options.GoalDir)
	if err != nil {
		return GoalReference{}, err
	}
	if _, err := parseResourceVersion(options.ExpectedResourceVersion); err != nil {
		return GoalReference{}, fmt.Errorf("expected resource version: %w", err)
	}
	destinationRoot, err := s.resolveInsideWorkspace(options.DestinationGoalsRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("destination goals root: %w", err)
	}
	target := filepath.Join(destinationRoot, filepath.Base(source))
	sourceIdentity, targetIdentity, err := s.canonicalNonoverlappingGoalIdentities(
		source,
		target,
	)
	if err != nil {
		return GoalReference{}, err
	}
	ownerRoot := options.OwnerRoot
	if ownerRoot == "" {
		ownerRoot = filepath.Dir(destinationRoot)
	}
	ownerRoot, err = s.resolveInsideWorkspace(ownerRoot)
	if err != nil {
		return GoalReference{}, err
	}
	ownerRef, err := portableOwnerRoot(s.workspaceRoot, ownerRoot)
	if err != nil {
		return GoalReference{}, err
	}
	sourceLock, targetLock, err := s.acquireOrderedGoalLocks(sourceIdentity, targetIdentity)
	if err != nil {
		return GoalReference{}, err
	}
	defer sourceLock.release()
	defer targetLock.release()
	if err := s.checkNoPendingPublication(source); err != nil {
		return GoalReference{}, err
	}
	goal, criteria, attempts, err := s.loadAndValidate(source)
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
	if goal.Spec.Scope != "workspace" {
		return GoalReference{}, fmt.Errorf("only workspace goals can be promoted")
	}
	if goal.Status.Execution != "paused" || goal.Status.ActiveAttemptID != "" {
		return GoalReference{}, fmt.Errorf("promotion requires paused execution and no active attempt")
	}
	recordFiles, err := inspectPromotionRecord(source)
	if err != nil {
		return GoalReference{}, err
	}
	digest, err := digestRecord(source, recordFiles)
	if err != nil {
		return GoalReference{}, err
	}
	sourceStateDigest, err := goalStateDigest(goal)
	if err != nil {
		return GoalReference{}, err
	}
	sourceCriteriaDigest, err := criteriaPortableDigest(criteria)
	if err != nil {
		return GoalReference{}, err
	}
	if pathExists(target) {
		existing, _, _, validateErr := s.loadAndValidate(target)
		if validateErr != nil {
			return GoalReference{}, fmt.Errorf(
				"existing destination goal is invalid: %w",
				validateErr,
			)
		}
		if _, err := inspectPromotionRecord(target); err != nil {
			return GoalReference{}, fmt.Errorf(
				"existing destination goal layout is invalid: %w",
				err,
			)
		}
		if existing.Metadata.Annotations[localOwnerRootAnnotation] != ownerRef {
			return GoalReference{}, fmt.Errorf(
				"existing destination goal belongs to a different owner root",
			)
		}
		if existing.Status.Promotion.SourceDigest == digest &&
			existing.Status.Promotion.SourceGeneration == goal.Metadata.Generation &&
			existing.Status.Promotion.SourceLifecycleGeneration == goal.Status.LifecycleGeneration &&
			existing.Status.Promotion.SourceCriteriaRevision == goal.Status.CriteriaRevision &&
			existing.Status.Promotion.SourceCriteriaDigest == sourceCriteriaDigest &&
			existing.Status.Promotion.SourceStateDigest == sourceStateDigest {
			return GoalReference{
				GoalID:          existing.Metadata.Name,
				GoalRef:         existing.Metadata.Name,
				ResourceVersion: existing.Metadata.ResourceVersion,
			}, nil
		}
		return GoalReference{}, fmt.Errorf("destination goal already exists with different provenance")
	}
	if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
		return GoalReference{}, err
	}
	temporary, err := os.MkdirTemp(destinationRoot, ".goal-promote-")
	if err != nil {
		return GoalReference{}, err
	}
	installed := false
	defer func() {
		if !installed {
			_ = os.RemoveAll(temporary)
		}
	}()
	if err := copyPromotedRecordFiles(source, temporary, recordFiles); err != nil {
		return GoalReference{}, err
	}
	nextResourceVersion, err := incrementResourceVersion(goal.Metadata.ResourceVersion)
	if err != nil {
		return GoalReference{}, err
	}
	goal.Metadata.ResourceVersion = nextResourceVersion
	goal.Metadata.Generation++
	goal.Spec.Scope = "project"
	if goal.Metadata.Annotations == nil {
		goal.Metadata.Annotations = map[string]string{}
	}
	goal.Metadata.Annotations[localOwnerRootAnnotation] = ownerRef
	goal.Spec.Retention.Policy = "durable"
	goal.Status.Promotion = PromotionStatus{
		SourceScope:               "workspace",
		SourceGeneration:          goal.Metadata.Generation - 1,
		SourceLifecycleGeneration: goal.Status.LifecycleGeneration,
		SourceCriteriaRevision:    goal.Status.CriteriaRevision,
		SourceCriteriaDigest:      sourceCriteriaDigest,
		SourceStateDigest:         sourceStateDigest,
		SourceDigest:              digest,
		PromotedAt:                s.timestamp(),
	}
	goal.Status.ObservedAt = s.timestamp()
	if err := goal.validate(); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(temporary, "goal.yaml"), goal); err != nil {
		return GoalReference{}, err
	}
	readme, err := renderREADME(goal, criteria, attempts, defaultOutputLimit)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(filepath.Join(temporary, "README.md"), readme, 0o644); err != nil {
		return GoalReference{}, err
	}
	if err := s.scanPromotedArtifacts(temporary); err != nil {
		return GoalReference{}, err
	}
	if err := s.validateStagedGoal(temporary, goal, criteria, attempts); err != nil {
		return GoalReference{}, fmt.Errorf("validate staged promoted goal: %w", err)
	}
	if err := s.callBeforeRename(target); err != nil {
		return GoalReference{}, err
	}
	if err := os.Rename(temporary, target); err != nil {
		return GoalReference{}, fmt.Errorf("install promoted goal: %w", err)
	}
	installed = true
	return GoalReference{
		GoalID:          goal.Metadata.Name,
		GoalRef:         goal.Metadata.Name,
		ResourceVersion: goal.Metadata.ResourceVersion,
	}, nil
}

func (s *Store) canonicalNonoverlappingGoalIdentities(
	source string,
	target string,
) (string, string, error) {
	sourceIdentity, err := s.resolveInsideWorkspace(source)
	if err != nil {
		return "", "", fmt.Errorf("resolve source goal identity: %w", err)
	}
	targetIdentity, err := s.resolveInsideWorkspace(target)
	if err != nil {
		return "", "", fmt.Errorf("resolve destination goal identity: %w", err)
	}
	if sourceIdentity == targetIdentity {
		return "", "", fmt.Errorf("source and destination goal are the same")
	}
	if isWithin(sourceIdentity, targetIdentity) || isWithin(targetIdentity, sourceIdentity) {
		return "", "", fmt.Errorf("source and destination goals must not overlap")
	}
	return sourceIdentity, targetIdentity, nil
}

func (s *Store) acquireOrderedGoalLocks(
	sourceIdentity string,
	targetIdentity string,
) (*heldLock, *heldLock, error) {
	// Promotion reads one goal while creating another. Sort the canonical goal
	// paths so concurrent cross-catalog promotions cannot deadlock.
	requests := []string{
		sourceIdentity,
		targetIdentity,
	}
	sort.Slice(requests, func(i, j int) bool {
		return requests[i] < requests[j]
	})

	held := make(map[string]*heldLock, len(requests))
	for _, request := range requests {
		lock, err := s.acquireGoalLock(request)
		if err != nil {
			for _, acquired := range held {
				_ = acquired.release()
			}
			return nil, nil, err
		}
		held[request] = lock
	}
	return held[sourceIdentity], held[targetIdentity], nil
}

func (s *Store) scanPromotedArtifacts(root string) error {
	return filepath.WalkDir(filepath.Join(root, "attempts"), func(
		path string,
		entry os.DirEntry,
		err error,
	) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		name := entry.Name()
		if name == "attempt.yaml" {
			return nil
		}
		content, err := readMarkdownFile(path, maxEvidenceFileBytes)
		if err != nil {
			return err
		}
		text := string(content)
		if strings.Contains(text, s.workspaceRoot) ||
			strings.Contains(text, "file:///") ||
			strings.Contains(text, "file://localhost/") {
			return fmt.Errorf(
				"promotion artifact %s contains a non-portable absolute workspace reference",
				filepath.ToSlash(strings.TrimPrefix(path, root+string(filepath.Separator))),
			)
		}
		return nil
	})
}

func (s *Store) validateStagedGoal(
	dir string,
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
	history, err := s.loadCriteriaHistory(dir, goal, criteria)
	if err != nil {
		return err
	}
	for _, attempt := range attempts {
		attemptDir := filepath.Join(dir, "attempts", attempt.Metadata.Name)
		if err := attempt.validate(goal); err != nil {
			return err
		}
		if err := validateAttemptArtifacts(attemptDir, attempt.Status.Artifacts); err != nil {
			return err
		}
		snapshot, ok := history.Snapshots[attempt.Spec.CriteriaRevision]
		if !ok {
			return fmt.Errorf("attempt %q has no criteria snapshot", attempt.Metadata.Name)
		}
		digest, err := criteriaPortableDigest(snapshot)
		if err != nil || digest != attempt.Spec.CriteriaDigest {
			return fmt.Errorf("attempt %q criteria digest mismatch", attempt.Metadata.Name)
		}
	}
	return validateProspectiveRecord(goal, criteria, attempts)
}

type MigrateOptions struct {
	SourceGoalDir        string
	DestinationGoalsRoot string
	GoalID               string
	Title                string
	Scope                string
	OwnerRoot            string
	Criteria             []string
}

func (s *Store) Migrate(options MigrateOptions) (GoalReference, error) {
	source, err := s.resolveInsideWorkspace(options.SourceGoalDir)
	if err != nil {
		return GoalReference{}, fmt.Errorf("source goal directory: %w", err)
	}
	destinationRoot, err := s.resolveInsideWorkspace(options.DestinationGoalsRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("destination goals root: %w", err)
	}
	goalID := filepath.Base(source)
	if options.GoalID != "" && options.GoalID != goalID {
		return GoalReference{}, fmt.Errorf("--goal-id must match the existing directory name %q", goalID)
	}
	if err := validateRecordID("goal ID", goalID); err != nil {
		return GoalReference{}, fmt.Errorf("rename the unversioned directory before migration: %w", err)
	}
	target := filepath.Join(destinationRoot, goalID)
	sourceIdentity, targetIdentity, err := s.canonicalNonoverlappingGoalIdentities(
		source,
		target,
	)
	if err != nil {
		return GoalReference{}, err
	}
	ownerRoot := options.OwnerRoot
	if ownerRoot == "" {
		ownerRoot = filepath.Dir(destinationRoot)
	}
	ownerRoot, err = s.resolveInsideWorkspace(ownerRoot)
	if err != nil {
		return GoalReference{}, fmt.Errorf("owner root: %w", err)
	}
	ownerRef, err := portableOwnerRoot(s.workspaceRoot, ownerRoot)
	if err != nil {
		return GoalReference{}, err
	}
	sourceLock, targetLock, err := s.acquireOrderedGoalLocks(
		sourceIdentity,
		targetIdentity,
	)
	if err != nil {
		return GoalReference{}, err
	}
	defer sourceLock.release()
	defer targetLock.release()

	if err := s.checkNoPendingPublication(target); err != nil {
		return GoalReference{}, err
	}
	legacyFiles, err := inspectUnversionedRecord(source)
	if err != nil {
		return GoalReference{}, err
	}
	readme := legacyFiles["README.md"]
	title := strings.TrimSpace(options.Title)
	if title == "" {
		title, err = extractUnambiguousTitle(readme)
		if err != nil {
			return GoalReference{}, err
		}
	}
	criteriaStatements := make([]string, len(options.Criteria))
	for index, statement := range options.Criteria {
		criteriaStatements[index] = strings.TrimSpace(statement)
	}
	if len(criteriaStatements) == 0 {
		criteriaStatements = extractCriteria(readme)
	}
	if len(criteriaStatements) > maxCriteria {
		return GoalReference{}, fmt.Errorf("criteria cardinality exceeds %d", maxCriteria)
	}
	scope := options.Scope
	if scope == "" {
		scope = "workspace"
	}
	digest := digestLegacyFiles(legacyFiles)
	retention := "ephemeral"
	if scope == "project" {
		retention = "durable"
	}
	if pathExists(target) {
		return s.matchExistingMigration(
			target,
			goalID,
			title,
			scope,
			retention,
			ownerRef,
			criteriaStatements,
			digest,
		)
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
			Title:     title,
			Scope:     scope,
			Retention: Retention{Policy: retention},
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
			Migration: MigrationStatus{
				SourceFormat: "unversioned",
				SourceDigest: digest,
				MigratedAt:   now,
			},
			ObservedAt: now,
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
			Items:    make([]Criterion, 0, len(criteriaStatements)),
		},
	}
	for index, statement := range criteriaStatements {
		criteria.Spec.Items = append(criteria.Spec.Items, Criterion{
			CriterionID:    fmt.Sprintf("criterion-%03d", index+1),
			Revision:       1,
			Required:       true,
			Statement:      statement,
			EvidenceMethod: "Inspect linked evidence against the criterion.",
		})
	}
	criteriaDigest, err := criteriaPortableDigest(criteria)
	if err != nil {
		return GoalReference{}, err
	}
	stateDigest, err := goalStateDigest(goal)
	if err != nil {
		return GoalReference{}, err
	}
	resultContent := []byte("# Result\n\nImported from an unversioned prose goal record.\n")
	artifacts := ArtifactManifest{
		PlanDigest:   digestBytes(readme),
		ResultDigest: digestBytes(resultContent),
		Evidence:     []ArtifactDigest{},
	}
	for name, content := range legacyFiles {
		if name != "README.md" {
			artifacts.Evidence = append(artifacts.Evidence, ArtifactDigest{
				Path: "evidence/" + name, Digest: digestBytes(content),
			})
		}
	}
	sort.Slice(artifacts.Evidence, func(i, j int) bool {
		return artifacts.Evidence[i].Path < artifacts.Evidence[j].Path
	})
	review := CloseReview{Decision: "refine", Criteria: []CriterionReview{}}
	for _, item := range criteria.Spec.Items {
		review.Criteria = append(review.Criteria, CriterionReview{
			CriterionID:       item.CriterionID,
			CriterionRevision: item.Revision,
			Verdict:           "unverified",
			EvidenceRefs:      []string{},
		})
	}
	attemptID := "imported-unversioned"
	attempt := AttemptManifest{
		APIVersion: goalAPIVersion,
		Kind:       "GoalAttempt",
		Metadata: ObjectMeta{
			Name:              attemptID,
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: now,
		},
		Spec: AttemptSpec{
			GoalRef:             LocalGoalReference{Name: goalID},
			GoalGeneration:      1,
			LifecycleGeneration: 1,
			CriteriaRevision:    1,
			CriteriaDigest:      criteriaDigest,
			GoalStateDigest:     stateDigest,
			WorkType:            "investigation",
		},
		Status: AttemptStatus{
			State: "closed", ClosedAt: now, Artifacts: artifacts,
			Review: review, ObservedAt: now,
		},
	}
	if err := goal.validate(); err != nil {
		return GoalReference{}, err
	}
	if err := criteria.validate(goal); err != nil {
		return GoalReference{}, err
	}
	if err := attempt.validate(goal); err != nil {
		return GoalReference{}, err
	}
	if err := os.MkdirAll(destinationRoot, 0o755); err != nil {
		return GoalReference{}, err
	}
	stagingRoot, err := os.MkdirTemp(
		destinationRoot,
		".goal-migrate-"+goalID+"-",
	)
	if err != nil {
		return GoalReference{}, err
	}
	defer func() {
		_ = os.RemoveAll(stagingRoot)
	}()
	staging := filepath.Join(stagingRoot, goalID)
	if err := os.Mkdir(staging, 0o755); err != nil {
		return GoalReference{}, err
	}
	attemptDir := filepath.Join(staging, "attempts", attemptID)
	if err := os.MkdirAll(filepath.Join(attemptDir, "evidence"), 0o755); err != nil {
		return GoalReference{}, err
	}
	if err := os.Mkdir(filepath.Join(staging, "criteria-revisions"), 0o755); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(staging, "goal.yaml"), goal); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(staging, "criteria.yaml"), criteria); err != nil {
		return GoalReference{}, err
	}
	criteriaBytes, err := marshalYAML(criteria)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(criteriaSnapshotPath(staging, 1), criteriaBytes, 0o644); err != nil {
		return GoalReference{}, err
	}
	if err := s.writeYAML(filepath.Join(attemptDir, "attempt.yaml"), attempt); err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(filepath.Join(attemptDir, "plan.md"), readme, 0o644); err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(
		filepath.Join(attemptDir, "result.md"),
		resultContent,
		0o644,
	); err != nil {
		return GoalReference{}, err
	}
	for name, content := range legacyFiles {
		if name == "README.md" {
			continue
		}
		if err := s.atomicWrite(filepath.Join(attemptDir, "evidence", name), content, 0o644); err != nil {
			return GoalReference{}, err
		}
	}
	projected, err := renderREADME(goal, criteria, []AttemptManifest{attempt}, defaultOutputLimit)
	if err != nil {
		return GoalReference{}, err
	}
	if err := s.atomicWrite(filepath.Join(staging, "README.md"), projected, 0o644); err != nil {
		return GoalReference{}, err
	}
	if _, _, _, err := s.loadAndValidate(staging); err != nil {
		return GoalReference{}, fmt.Errorf("validate staged migration: %w", err)
	}
	if err := s.callBeforeRename(target); err != nil {
		return GoalReference{}, err
	}
	if pathExists(target) {
		return GoalReference{}, fmt.Errorf("migration target appeared before publication")
	}
	latestFiles, err := inspectUnversionedRecord(source)
	if err != nil {
		return GoalReference{}, fmt.Errorf("re-read migration source: %w", err)
	}
	if latestDigest := digestLegacyFiles(latestFiles); latestDigest != digest {
		return GoalReference{}, fmt.Errorf("migration source changed before publication")
	}
	if err := os.Rename(staging, target); err != nil {
		return GoalReference{}, fmt.Errorf("publish migrated goal without overwrite: %w", err)
	}
	return GoalReference{GoalID: goalID, GoalRef: goalID, ResourceVersion: "1"}, nil
}

func (s *Store) matchExistingMigration(
	target string,
	goalID string,
	title string,
	scope string,
	retention string,
	ownerRef string,
	criteriaStatements []string,
	sourceDigest string,
) (GoalReference, error) {
	goal, criteria, _, err := s.loadAndValidate(target)
	if err != nil {
		return GoalReference{}, fmt.Errorf("existing migration target is invalid: %w", err)
	}
	if goal.Metadata.Name != goalID ||
		goal.Status.Migration.SourceFormat != "unversioned" ||
		goal.Status.Migration.SourceDigest != sourceDigest ||
		goal.Spec.Title != title ||
		goal.Spec.Scope != scope ||
		goal.Spec.Retention.Policy != retention ||
		goal.Metadata.Annotations[localOwnerRootAnnotation] != ownerRef {
		return GoalReference{}, fmt.Errorf(
			"existing migration target has different provenance or options",
		)
	}
	if len(criteria.Spec.Items) != len(criteriaStatements) {
		return GoalReference{}, fmt.Errorf(
			"existing migration target has different criteria options",
		)
	}
	for index, statement := range criteriaStatements {
		if criteria.Spec.Items[index].Statement != statement {
			return GoalReference{}, fmt.Errorf(
				"existing migration target has different criteria options",
			)
		}
	}
	return GoalReference{
		GoalID:          goal.Metadata.Name,
		GoalRef:         goal.Metadata.Name,
		ResourceVersion: goal.Metadata.ResourceVersion,
	}, nil
}

func inspectUnversionedRecord(dir string) (map[string][]byte, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	hasREADME := false
	evidenceCount := 0
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 || entry.IsDir() ||
			!strings.HasSuffix(strings.ToLower(entry.Name()), ".md") ||
			!safeEvidenceName(entry.Name()) {
			return nil, fmt.Errorf(
				"ambiguous unversioned record entry %q; migration accepts only root Markdown files",
				entry.Name(),
			)
		}
		if entry.Name() == "README.md" {
			hasREADME = true
		} else {
			evidenceCount++
		}
	}
	if !hasREADME {
		return nil, fmt.Errorf("unversioned migration requires README.md")
	}
	if evidenceCount > maxEvidenceFiles {
		return nil, fmt.Errorf("legacy evidence cardinality exceeds %d", maxEvidenceFiles)
	}
	files := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		maximum := int64(maxEvidenceFileBytes)
		if entry.Name() == "README.md" {
			maximum = maxPlanResultBytes
		}
		content, err := readMarkdownFile(
			filepath.Join(dir, entry.Name()),
			maximum,
		)
		if err != nil {
			return nil, err
		}
		files[entry.Name()] = content
	}
	return files, nil
}

func extractUnambiguousTitle(readme []byte) (string, error) {
	var titles []string
	for _, line := range strings.Split(string(readme), "\n") {
		if strings.HasPrefix(line, "# ") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
			if title != "" {
				titles = append(titles, title)
			}
		}
	}
	if len(titles) != 1 {
		return "", fmt.Errorf("unversioned README must have exactly one H1 or --title must be supplied")
	}
	return titles[0], nil
}

func extractCriteria(readme []byte) []string {
	lines := strings.Split(string(readme), "\n")
	inSection := false
	var criteria []string
	for _, line := range lines {
		lower := strings.ToLower(strings.TrimSpace(line))
		if strings.HasPrefix(lower, "## ") {
			if inSection {
				break
			}
			inSection = strings.Contains(lower, "acceptance") && strings.Contains(lower, "criteria")
			continue
		}
		if inSection && (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ")) {
			value := strings.TrimSpace(line[2:])
			if value != "" {
				criteria = append(criteria, value)
			}
		}
	}
	return criteria
}

func digestLegacyFiles(files map[string][]byte) string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	hash := sha256.New()
	writeDigestFrame(hash, []byte("goal-legacy-record-v1"))
	for _, name := range names {
		writeDigestFrame(hash, []byte(name))
		writeDigestFrame(hash, files[name])
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil))
}

type promotionRecordFile struct {
	relative string
	maximum  int64
}

func inspectPromotionRecord(root string) ([]promotionRecordFile, error) {
	rootEntries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	requiredRootEntries := map[string]bool{
		"attempts":           false,
		"criteria-revisions": false,
		"criteria.yaml":      false,
		"goal.yaml":          false,
	}
	allowedRootEntries := map[string]bool{
		"README.md":          true,
		"attempts":           true,
		"criteria-revisions": true,
		"criteria.yaml":      true,
		"goal.yaml":          true,
	}
	files := make([]promotionRecordFile, 0, len(rootEntries))
	for _, entry := range rootEntries {
		if !allowedRootEntries[entry.Name()] {
			return nil, fmt.Errorf("unexpected goal record entry %q", entry.Name())
		}
		if _, required := requiredRootEntries[entry.Name()]; required {
			requiredRootEntries[entry.Name()] = true
		}
		switch entry.Name() {
		case "attempts", "criteria-revisions":
			if err := validatePromotionDirectory(filepath.Join(root, entry.Name())); err != nil {
				return nil, err
			}
		case "goal.yaml", "criteria.yaml":
			files, err = appendPromotionFile(files, root, entry.Name(), maxManifestBytes)
			if err != nil {
				return nil, err
			}
		case "README.md":
			files, err = appendPromotionFile(files, root, entry.Name(), maxPlanResultBytes)
			if err != nil {
				return nil, err
			}
		}
	}
	for name, present := range requiredRootEntries {
		if !present {
			return nil, fmt.Errorf("missing goal record entry %q", name)
		}
	}

	criteriaRoot := filepath.Join(root, "criteria-revisions")
	criteriaEntries, err := os.ReadDir(criteriaRoot)
	if err != nil {
		return nil, err
	}
	if len(criteriaEntries) == 0 || len(criteriaEntries) > maxCriteriaRevisions {
		return nil, fmt.Errorf(
			"criteria revision snapshot cardinality must be between 1 and %d",
			maxCriteriaRevisions,
		)
	}
	for _, entry := range criteriaEntries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") ||
			!strings.HasSuffix(entry.Name(), ".yaml") {
			return nil, fmt.Errorf("invalid criteria snapshot entry %q", entry.Name())
		}
		relative := filepath.Join("criteria-revisions", entry.Name())
		files, err = appendPromotionFile(files, root, relative, maxManifestBytes)
		if err != nil {
			return nil, err
		}
	}

	attemptsRoot := filepath.Join(root, "attempts")
	attemptEntries, err := os.ReadDir(attemptsRoot)
	if err != nil {
		return nil, err
	}
	if len(attemptEntries) > maxAttempts {
		return nil, fmt.Errorf("attempt cardinality exceeds %d", maxAttempts)
	}
	for _, attemptEntry := range attemptEntries {
		if !attemptEntry.IsDir() || strings.HasPrefix(attemptEntry.Name(), ".") {
			return nil, fmt.Errorf("unexpected entry in attempts: %s", attemptEntry.Name())
		}
		attemptRelative := filepath.Join("attempts", attemptEntry.Name())
		attemptRoot := filepath.Join(root, attemptRelative)
		if err := validatePromotionDirectory(attemptRoot); err != nil {
			return nil, err
		}
		if err := validateAttemptFiles(attemptRoot); err != nil {
			return nil, err
		}
		for _, file := range []struct {
			name    string
			maximum int64
		}{
			{name: "attempt.yaml", maximum: maxManifestBytes},
			{name: "plan.md", maximum: maxPlanResultBytes},
			{name: "result.md", maximum: maxPlanResultBytes},
		} {
			relative := filepath.Join(attemptRelative, file.name)
			files, err = appendPromotionFile(files, root, relative, file.maximum)
			if err != nil {
				return nil, err
			}
		}
		evidenceRoot := filepath.Join(attemptRoot, "evidence")
		evidenceEntries, err := os.ReadDir(evidenceRoot)
		if err != nil {
			return nil, err
		}
		if len(evidenceEntries) > maxEvidenceFiles {
			return nil, fmt.Errorf("evidence cardinality exceeds %d", maxEvidenceFiles)
		}
		for _, entry := range evidenceEntries {
			if entry.IsDir() || !safeEvidenceName(entry.Name()) {
				return nil, fmt.Errorf("invalid evidence entry %q", entry.Name())
			}
			relative := filepath.Join(attemptRelative, "evidence", entry.Name())
			files, err = appendPromotionFile(files, root, relative, maxEvidenceFileBytes)
			if err != nil {
				return nil, err
			}
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].relative < files[j].relative
	})
	return files, nil
}

func appendPromotionFile(
	files []promotionRecordFile,
	root string,
	relative string,
	maximum int64,
) ([]promotionRecordFile, error) {
	info, err := os.Lstat(filepath.Join(root, relative))
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Size() > maximum {
		return nil, fmt.Errorf(
			"goal record file %q must be regular and at most %d bytes",
			filepath.ToSlash(relative),
			maximum,
		)
	}
	return append(files, promotionRecordFile{
		relative: filepath.Clean(relative),
		maximum:  maximum,
	}), nil
}

func validatePromotionDirectory(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("goal record entry %q must be a directory", filepath.Base(path))
	}
	return nil
}

func digestRecord(root string, files []promotionRecordFile) (string, error) {
	hash := sha256.New()
	writeDigestFrame(hash, []byte("goal-promoted-record-v1"))
	for _, file := range files {
		if file.relative == "README.md" {
			continue
		}
		content, err := readRegularFile(filepath.Join(root, file.relative), file.maximum)
		if err != nil {
			return "", err
		}
		writeDigestFrame(hash, []byte(filepath.ToSlash(file.relative)))
		writeDigestFrame(hash, content)
	}
	return "sha256:" + hex.EncodeToString(hash.Sum(nil)), nil
}

func writeDigestFrame(destination io.Writer, value []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = destination.Write(length[:])
	_, _ = destination.Write(value)
}

func copyPromotedRecordFiles(
	source string,
	destination string,
	files []promotionRecordFile,
) error {
	for _, relative := range []string{"attempts", "criteria-revisions"} {
		if err := os.MkdirAll(filepath.Join(destination, relative), 0o755); err != nil {
			return err
		}
	}
	for _, file := range files {
		if file.relative == "goal.yaml" || file.relative == "README.md" {
			continue
		}
		target := filepath.Join(destination, file.relative)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if filepath.Base(file.relative) == "attempt.yaml" {
			if err := os.MkdirAll(filepath.Join(filepath.Dir(target), "evidence"), 0o755); err != nil {
				return err
			}
		}
		content, err := readRegularFile(filepath.Join(source, file.relative), file.maximum)
		if err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return err
		}
	}
	if _, err := os.Stat(filepath.Join(destination, "criteria.yaml")); err != nil {
		return err
	}
	return nil
}
