package fsstore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
	"github.com/goccy/go-yaml"
)

// Publication intent files live inside the goal record (workspace-relative,
// owner-local) and are canonical recovery metadata, not locks and not a
// background service.
const (
	publicationIntentName = ".goal-publication.yaml"
	publicationStageDir   = ".goal-publication-stage"
)

// PublicationIncompleteError is a stable structured error returned by every
// normal mutation when a publication intent is pending. It must not be
// confused with a stale resource version or a generic validation failure.
type PublicationIncompleteError struct {
	OperationID      string
	IntendedRevision string
	Phase            string
	Kind             string
	Message          string
	Cause            error
}

func (e *PublicationIncompleteError) Error() string {
	// Keep the legacy committed-version phrasing for backward compatibility
	// with callers that parse the prefix, while adding the stable
	// publication-incomplete diagnosis and recovery action.
	prefix := "checkpoint committed"
	if e.Kind == "criteria" {
		prefix = "criteria update committed"
	} else if e.Kind == "relationships" {
		prefix = "relationships committed"
	}
	return fmt.Sprintf(
		"%s at resourceVersion %s; goal publication is incomplete: operation %q (phase %s); run goal doctor then goal recover: %s",
		prefix,
		e.IntendedRevision,
		e.OperationID,
		e.Phase,
		e.Message,
	)
}

// Unwrap exposes the underlying publication failure so callers can match the
// concrete injected error with errors.Is.
func (e *PublicationIncompleteError) Unwrap() error {
	return e.Cause
}

// publicationFileEntry is one canonical after-image in deterministic
// publication order. Path is workspace-relative to the goal record.
type publicationFileEntry struct {
	Path           string
	BeforeDigest   string
	Content        []byte
	StagedRelative string // optional override; defaults to Path
}

// publication evaluation states used by doctor/recover.
type publicationState int

const (
	publicationNone publicationState = iota
	publicationDiscardable
	publicationStaged
	publicationPartial
	publicationConflict
)

// DoctorResult is the structured, machine-readable diagnosis.
type DoctorResult struct {
	GoalID           string `json:"goalID"`
	PublicationState string `json:"publicationState"`
	OperationID      string `json:"operationID,omitempty"`
	ResourceVersion  string `json:"resourceVersion,omitempty"`
	Message          string `json:"message,omitempty"`
}

// RecoverResult is the structured recovery receipt.
type RecoverResult struct {
	GoalID          string   `json:"goalID"`
	OperationID     string   `json:"operationID"`
	Discarded       bool     `json:"discarded"`
	ResourceVersion string   `json:"resourceVersion"`
	Installed       []string `json:"installed,omitempty"`
}

// publicationIntent is the stored form of v1alpha1.GoalPublication.
type publicationIntent struct {
	APIVersion string                         `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                         `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta                     `json:"metadata" yaml:"metadata"`
	Spec       v1alpha1.GoalPublicationSpec   `json:"spec" yaml:"spec"`
	Status     v1alpha1.GoalPublicationStatus `json:"status" yaml:"status"`
}

func (i publicationIntent) validate() error {
	api := v1alpha1.GoalPublication{
		APIVersion: i.APIVersion,
		Kind:       i.Kind,
		Metadata:   v1alpha1.ObjectMeta(i.Metadata),
		Spec:       i.Spec,
		Status:     i.Status,
	}
	if err := api.Validate(); err != nil {
		return err
	}
	return validatePersistedMetadata(i.Metadata)
}

func publicationIntentPath(dir string) string {
	return filepath.Join(dir, publicationIntentName)
}

func publicationStageRoot(dir string) string {
	return filepath.Join(dir, publicationStageDir)
}

// readPublicationIntent loads a stored intent without validating the whole
// goal record. A malformed intent is a conflict and is refused, never
// silently deleted.
func (s *Store) readPublicationIntent(dir string) (*publicationIntent, error) {
	path := publicationIntentPath(dir)
	content, err := readRegularFile(path, maxManifestBytes)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var intent publicationIntent
	if err := yaml.UnmarshalWithOptions(content, &intent, yaml.Strict()); err != nil {
		return nil, fmt.Errorf("decode publication intent strictly: %w", err)
	}
	if err := intent.validate(); err != nil {
		return nil, fmt.Errorf("invalid stored publication intent: %w", err)
	}
	return &intent, nil
}

// evaluatePublication classifies a stored intent against the current record.
// It is read-only and returns the exact per-target state.
func (s *Store) evaluatePublication(
	dir string,
	intent *publicationIntent,
) (publicationState, string, error) {
	before := true
	for _, file := range intent.Spec.Files {
		path := filepath.Join(dir, filepath.FromSlash(file.Path))
		current, err := readRegularFile(path, maxManifestBytes)
		switch {
		case err == nil:
			digest := digestBytes(current)
			switch digest {
			case file.BeforeDigest:
				// Target still has its prior content.
			case file.AfterDigest:
				before = false
			default:
				return publicationConflict, fmt.Sprintf(
					"publication target %s matches neither its before nor after digest",
					file.Path,
				), nil
			}
		case os.IsNotExist(err):
			if file.BeforeDigest != "" {
				return publicationConflict, fmt.Sprintf(
					"publication target %s is absent but a before digest was recorded",
					file.Path,
				), nil
			}
		default:
			return publicationNone, "", err
		}
	}
	if before {
		// Every target is still the prior content. If staging exists with the
		// matching operation ID, the intent is fully staged and safe to replay;
		// otherwise it must be discarded to preserve the prior record.
		root := publicationStageRoot(dir)
		entries, err := os.ReadDir(root)
		if err == nil && len(entries) == 1 &&
			entries[0].IsDir() && entries[0].Name() == intent.Spec.OperationID {
			return publicationStaged, "", nil
		}
		return publicationDiscardable, "", nil
	}
	// At least one target is already in its after state. Recovery must replay
	// every remaining after-image and refuse any conflict.
	return publicationPartial, "", nil
}

// Doctor classifies one goal record. It returns a structured diagnosis without
// mutating canonical state.
func (s *Store) Doctor(goalDir string) (DoctorResult, error) {
	dir, err := s.resolveInsideWorkspace(goalDir)
	if err != nil {
		return DoctorResult{}, err
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return DoctorResult{}, err
	}
	defer lock.release()
	intent, err := s.readPublicationIntent(dir)
	if err != nil {
		return DoctorResult{}, err
	}
	if intent != nil {
		state, message, err := s.evaluatePublication(dir, intent)
		if err != nil {
			return DoctorResult{}, err
		}
		return DoctorResult{
			GoalID:           filepath.Base(dir),
			PublicationState: publicationStateName(state),
			OperationID:      intent.Spec.OperationID,
			ResourceVersion:  intent.Spec.IntendedResourceVersion,
			Message:          message,
		}, nil
	}
	// No intent: a valid record is stable; a README-only issue is a stale
	// projection; anything else is a genuine validation failure.
	_, _, _, err = s.loadAndValidate(dir)
	if err == nil {
		return DoctorResult{
			GoalID:           filepath.Base(dir),
			PublicationState: "stable",
		}, nil
	}
	// Distinguish a replaceable README projection issue from canonical record
	// corruption. loadCanonical validates goal/criteria/attempts without
	// requiring README.
	canonicalErr := s.loadCanonical(dir)
	if canonicalErr == nil {
		stored, readErr := readRegularFile(filepath.Join(dir, "README.md"), maxManifestBytes)
		if readErr != nil {
			stored = nil
		}
		want, renderErr := s.renderREADMEFromCanonical(dir)
		if renderErr != nil {
			return DoctorResult{}, fmt.Errorf("goal validation: %w", err)
		}
		if !equalBytes(stored, want) {
			return DoctorResult{
				GoalID:           filepath.Base(dir),
				PublicationState: "committed-projection-stale",
				Message:          "README projection does not match the canonical record; run goal render",
			}, nil
		}
	}
	return DoctorResult{}, fmt.Errorf("goal validation: %w", err)
}

// Recover finishes or discards a stored publication intent. It is
// idempotent: an already-after file is skipped and no input is regenerated.
func (s *Store) Recover(goalDir string) (RecoverResult, error) {
	dir, err := s.resolveInsideWorkspace(goalDir)
	if err != nil {
		return RecoverResult{}, err
	}
	lock, err := s.acquireGoalLock(dir)
	if err != nil {
		return RecoverResult{}, err
	}
	defer lock.release()
	intent, err := s.readPublicationIntent(dir)
	if err != nil {
		return RecoverResult{}, err
	}
	if intent == nil {
		return RecoverResult{}, fmt.Errorf("no publication intent is pending")
	}
	state, message, err := s.evaluatePublication(dir, intent)
	if err != nil {
		return RecoverResult{}, err
	}
	switch state {
	case publicationConflict:
		return RecoverResult{}, fmt.Errorf("refusing recovery: %s", message)
	case publicationDiscardable:
		if err := os.Remove(publicationIntentPath(dir)); err != nil {
			return RecoverResult{}, err
		}
		if err := os.RemoveAll(publicationStageRoot(dir)); err != nil {
			return RecoverResult{}, err
		}
		return RecoverResult{
			GoalID:          filepath.Base(dir),
			OperationID:     intent.Spec.OperationID,
			Discarded:       true,
			ResourceVersion: intent.Spec.PriorResourceVersion,
		}, nil
	case publicationStaged, publicationPartial:
		return s.replayPublication(dir, intent)
	default:
		return RecoverResult{}, fmt.Errorf("unexpected publication state %q", publicationStateName(state))
	}
}

func (s *Store) replayPublication(
	dir string,
	intent *publicationIntent,
) (RecoverResult, error) {
	stage := filepath.Join(publicationStageRoot(dir), intent.Spec.OperationID)
	installed := make([]string, 0, len(intent.Spec.Files))
	// Identify the new-attempt directory from the file paths.
	newAttemptDir := ""
	for _, file := range intent.Spec.Files {
		if strings.HasPrefix(file.Path, "attempts/") {
			parts := strings.SplitN(file.Path, "/", 3)
			if len(parts) >= 2 {
				candidate := "attempts/" + parts[1]
				if newAttemptDir == "" || candidate == newAttemptDir {
					newAttemptDir = candidate
				} else {
					return RecoverResult{}, fmt.Errorf(
						"recovery cannot handle multiple attempt directories in one intent",
					)
				}
			}
		}
	}
	if newAttemptDir != "" {
		targetAttempt := filepath.Join(dir, newAttemptDir)
		if !pathExists(targetAttempt) {
			stagedAttempt := filepath.Join(stage, newAttemptDir)
			if err := s.callBeforeRename(targetAttempt); err != nil {
				return RecoverResult{}, err
			}
			if err := os.MkdirAll(filepath.Dir(targetAttempt), 0o755); err != nil {
				return RecoverResult{}, err
			}
			if err := os.Rename(stagedAttempt, targetAttempt); err != nil {
				return RecoverResult{}, err
			}
			installed = append(installed, newAttemptDir)
		}
	}
	for _, file := range intent.Spec.Files {
		target := filepath.Join(dir, filepath.FromSlash(file.Path))
		if newAttemptDir != "" && strings.HasPrefix(file.Path, newAttemptDir+"/") {
			// The attempt directory rename already installed these files.
			current, err := readRegularFile(target, maxManifestBytes)
			if err != nil {
				return RecoverResult{}, err
			}
			if digestBytes(current) != file.AfterDigest {
				return RecoverResult{}, fmt.Errorf(
					"recovered attempt file %s does not match its staged after digest",
					file.Path,
				)
			}
			installed = append(installed, file.Path)
			continue
		}
		current, err := readRegularFile(target, maxManifestBytes)
		if err == nil {
			switch digestBytes(current) {
			case file.AfterDigest:
				installed = append(installed, file.Path)
				continue
			case file.BeforeDigest:
				// Replay below.
			default:
				return RecoverResult{}, fmt.Errorf(
					"refusing recovery: target %s matches neither its before nor after digest",
					file.Path,
				)
			}
		} else if !os.IsNotExist(err) {
			return RecoverResult{}, err
		} else if file.BeforeDigest != "" {
			return RecoverResult{}, fmt.Errorf(
				"refusing recovery: target %s is absent but a before digest was recorded",
				file.Path,
			)
		}
		if err := s.callBeforeRename(target); err != nil {
			return RecoverResult{}, err
		}
		if err := os.Rename(
			filepath.Join(stage, filepath.FromSlash(file.StagedRelative)),
			target,
		); err != nil {
			return RecoverResult{}, err
		}
		installed = append(installed, file.Path)
	}
	// The recovered record must validate, then the replaceable README
	// projection is refreshed.
	if err := s.refreshREADMEProjection(dir); err != nil {
		return RecoverResult{}, err
	}
	if err := os.Remove(publicationIntentPath(dir)); err != nil {
		return RecoverResult{}, err
	}
	if err := os.RemoveAll(publicationStageRoot(dir)); err != nil {
		return RecoverResult{}, err
	}
	return RecoverResult{
		GoalID:          filepath.Base(dir),
		OperationID:     intent.Spec.OperationID,
		Discarded:       false,
		ResourceVersion: intent.Spec.IntendedResourceVersion,
		Installed:       installed,
	}, nil
}

// beginPublication stages exact after-images and installs a durable intent.
// It must be called while the goal lock is held and before any canonical file
// changes. Entries must be in deterministic publication order (goal.yaml first,
// then criteria, then attempt files, then README). It returns the operation ID
// so the caller can finish the intent after publishing.
// stageDirs lists goal-relative directories that must exist inside the
// staging tree before publication (for example a new attempt's evidence
// directory when it has no evidence files).
func (s *Store) beginPublication(
	dir string,
	goalID string,
	priorResourceVersion string,
	intendedResourceVersion string,
	entries []publicationFileEntry,
	stageDirs []string,
) (*publicationIntent, error) {
	operationID, err := s.generateOperationID()
	if err != nil {
		return nil, err
	}
	stageRoot := publicationStageRoot(dir)
	stage := filepath.Join(stageRoot, operationID)
	if err := os.MkdirAll(stage, 0o755); err != nil {
		return nil, err
	}
	for _, dirPath := range stageDirs {
		if err := os.MkdirAll(filepath.Join(stage, filepath.FromSlash(dirPath)), 0o755); err != nil {
			return nil, err
		}
	}
	installed := false
	defer func() {
		if !installed {
			_ = os.RemoveAll(stage)
		}
	}()
	files := make([]v1alpha1.GoalPublicationFile, 0, len(entries))
	for _, entry := range entries {
		staged := entry.StagedRelative
		if staged == "" {
			staged = portablePublicationName(entry.Path)
		} else {
			staged = portablePublicationName(staged)
		}
		if staged == "" || staged == "." {
			return nil, fmt.Errorf("unsafe publication path %q", entry.Path)
		}
		if entry.BeforeDigest != "" && !validDigest(entry.BeforeDigest) {
			return nil, fmt.Errorf("invalid before digest for %q", entry.Path)
		}
		target := filepath.Join(stage, filepath.FromSlash(staged))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, err
		}
		if err := writeStagedFile(target, entry.Content); err != nil {
			return nil, err
		}
		files = append(files, v1alpha1.GoalPublicationFile{
			Path:           filepath.ToSlash(entry.Path),
			BeforeDigest:   entry.BeforeDigest,
			AfterDigest:    digestBytes(entry.Content),
			StagedRelative: staged,
		})
	}
	intent := publicationIntent{
		APIVersion: goalAPIVersion,
		Kind:       v1alpha1.KindGoalPublication,
		Metadata: ObjectMeta{
			Name:              operationID,
			ResourceVersion:   "1",
			Generation:        1,
			CreationTimestamp: s.timestamp(),
		},
		Spec: v1alpha1.GoalPublicationSpec{
			GoalRef:                 LocalGoalReference{Name: goalID},
			OperationID:             operationID,
			PriorResourceVersion:    priorResourceVersion,
			IntendedResourceVersion: intendedResourceVersion,
			Files:                   files,
			SnapshotDigests:         map[uint64]string{},
		},
		Status: v1alpha1.GoalPublicationStatus{
			State:      v1alpha1.PublicationIncomplete,
			ObservedAt: s.timestamp(),
		},
	}
	if err := s.writePublicationIntent(dir, intent); err != nil {
		return nil, err
	}
	installed = true
	stored, err := s.readPublicationIntent(dir)
	if err != nil || stored == nil {
		return nil, fmt.Errorf("stored publication intent is unreadable: %w", err)
	}
	return stored, nil
}

// writePublicationIntent writes the intent into the goal directory.
func (s *Store) writePublicationIntent(dir string, intent publicationIntent) error {
	if err := intent.validate(); err != nil {
		return err
	}
	return s.writeYAML(publicationIntentPath(dir), intent)
}

// publishIntentFiles installs the staged after-images over the canonical
// targets. The order follows the historical optimistic-concurrency contract:
// the commit-point goal.yaml publishes first, then a new attempt directory is
// published as one atomic rename (firing the directory before-rename hook),
// then any remaining files (final goal.yaml for immediate close) follow.
// Files inside the attempt directory are captured by the directory rename and
// skipped individually.
func (s *Store) publishIntentFiles(
	dir string,
	intent *publicationIntent,
	newAttemptDir string,
) error {
	stage := filepath.Join(publicationStageRoot(dir), intent.Spec.OperationID)
	install := func(file v1alpha1.GoalPublicationFile) error {
		target := filepath.Join(dir, filepath.FromSlash(file.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := s.callBeforeRename(target); err != nil {
			return err
		}
		return os.Rename(
			filepath.Join(stage, filepath.FromSlash(file.StagedRelative)),
			target,
		)
	}
	if newAttemptDir == "" {
		for _, file := range intent.Spec.Files {
			if err := install(file); err != nil {
				return err
			}
		}
		return nil
	}
	// New attempt: publish the commit-point goal.yaml (StagedRelative
	// "goal.yaml") first, then the attempt directory, then the remaining
	// files in order (final goal.yaml, plus any non-attempt projection files
	// that were withheld from the intent).
	commitPoint := ""
	for _, file := range intent.Spec.Files {
		if file.Path == "goal.yaml" && file.StagedRelative == "goal.yaml" {
			commitPoint = file.StagedRelative
			continue
		}
		if file.Path == "goal.yaml" && file.StagedRelative == "goal-final.yaml" {
			continue
		}
	}
	if commitPoint != "" {
		for _, file := range intent.Spec.Files {
			if file.StagedRelative != commitPoint {
				continue
			}
			if err := install(file); err != nil {
				return err
			}
			break
		}
	}
	targetAttempt := filepath.Join(dir, filepath.FromSlash(newAttemptDir))
	if err := s.callBeforeRename(targetAttempt); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(targetAttempt), 0o755); err != nil {
		return err
	}
	if err := os.Rename(
		filepath.Join(stage, filepath.FromSlash(newAttemptDir)),
		targetAttempt,
	); err != nil {
		return err
	}
	// Publish remaining non-attempt files in intent order, skipping the
	// already-installed commit point.
	for _, file := range intent.Spec.Files {
		if strings.HasPrefix(file.Path, newAttemptDir+"/") {
			continue
		}
		if commitPoint != "" && file.StagedRelative == commitPoint {
			continue
		}
		if err := install(file); err != nil {
			return err
		}
	}
	return nil
}

// goalCommitted reports whether the optimistic-concurrency commit point
// (goal.yaml) currently matches the intent's final intended after digest.
func (s *Store) goalCommitted(dir string, intent *publicationIntent) bool {
	finalAfter := ""
	for _, file := range intent.Spec.Files {
		if file.Path == "goal.yaml" {
			finalAfter = file.AfterDigest
		}
	}
	current, err := readRegularFile(filepath.Join(dir, "goal.yaml"), maxManifestBytes)
	if err != nil {
		return false
	}
	return finalAfter != "" && digestBytes(current) == finalAfter
}

// commitPointPublished reports whether the first goal.yaml after-image (the
// optimistic-concurrency commit point) has been installed. Unlike
// goalCommitted, it compares against the FIRST goal.yaml entry and ignores a
// finalization after-image (goal-final.yaml), so the checkpoint discard
// decision matches the actual first canonical rename.
func (s *Store) commitPointPublished(dir string, intent *publicationIntent) bool {
	firstAfter := ""
	for _, file := range intent.Spec.Files {
		if file.Path == "goal.yaml" {
			firstAfter = file.AfterDigest
			break
		}
	}
	current, err := readRegularFile(filepath.Join(dir, "goal.yaml"), maxManifestBytes)
	if err != nil {
		return false
	}
	return firstAfter != "" && digestBytes(current) == firstAfter
}

// finishPublication removes the intent and staging directory once a
// publication has fully validated.
func (s *Store) finishPublication(dir string) error {
	if err := os.Remove(publicationIntentPath(dir)); err != nil {
		return err
	}
	return os.RemoveAll(publicationStageRoot(dir))
}

// refreshREADMEProjection loads the now-complete canonical record and
// atomically rewrites the replaceable README projection, firing the
// before-rename hook so projection failures remain injectable.
func (s *Store) refreshREADMEProjection(dir string) error {
	goal, criteria, attempts, err := s.loadAndValidate(dir)
	if err != nil {
		return err
	}
	content, err := renderREADME(goal, criteria, attempts, defaultOutputLimit)
	if err != nil {
		return err
	}
	return s.atomicWrite(filepath.Join(dir, "README.md"), content, 0o644)
}

// checkNoPendingPublication fails closed when a publication intent is
// pending. It must be called while the goal lock is held before any normal
// mutation or strict read.
func (s *Store) checkNoPendingPublication(dir string) error {
	return s.pendingPublication(dir)
}

// pendingPublication returns a stable error when a publication intent exists
// before a normal mutation. It is used to fail closed instead of reporting a
// stale resource version or a generic validation failure.
func (s *Store) pendingPublication(dir string) error {
	intent, err := s.readPublicationIntent(dir)
	if err != nil {
		return err
	}
	if intent == nil {
		return nil
	}
	return &PublicationIncompleteError{
		OperationID:      intent.Spec.OperationID,
		IntendedRevision: intent.Spec.IntendedResourceVersion,
		Phase:            "preflight",
		Kind:             "checkpoint",
		Message:          "doctor/recover must resolve the pending intent before another mutation",
	}
}

func portablePublicationName(path string) string {
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	if clean == "." || strings.HasPrefix(clean, "../") || filepath.IsAbs(path) {
		return ""
	}
	return clean
}

func (s *Store) generateOperationID() (string, error) {
	suffix := make([]byte, 6)
	if _, err := io.ReadFull(s.random, suffix); err != nil {
		return "", fmt.Errorf("generate publication operation ID: %w", err)
	}
	return "pub-" + hexEncode(suffix), nil
}

func hexEncode(value []byte) string {
	const digits = "0123456789abcdef"
	result := make([]byte, 2*len(value))
	for index, b := range value {
		result[2*index] = digits[b>>4]
		result[2*index+1] = digits[b&0x0f]
	}
	return string(result)
}

func equalBytes(left []byte, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

// loadCanonical validates the canonical goal record without requiring or
// inspecting the replaceable README projection. It is used by Doctor to
// distinguish projection staleness from canonical corruption.
func (s *Store) loadCanonical(dir string) error {
	if err := rejectRecordSymlinks(dir); err != nil {
		return err
	}
	goal, err := s.readGoalManifest(dir)
	if err != nil {
		return err
	}
	if filepath.Base(dir) != goal.Metadata.Name {
		return fmt.Errorf("directory name does not match Goal metadata.name")
	}
	var criteria CriteriaManifest
	if err := s.readYAML(filepath.Join(dir, "criteria.yaml"), &criteria); err != nil {
		return err
	}
	if err := criteria.validate(goal); err != nil {
		return err
	}
	history, err := s.loadCriteriaHistory(dir, goal, criteria)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(filepath.Join(dir, "attempts"))
	if os.IsNotExist(err) {
		entries = nil
	} else if err != nil {
		return err
	}
	if len(entries) > maxAttempts {
		return fmt.Errorf("attempt cardinality exceeds %d", maxAttempts)
	}
	attempts := make([]AttemptManifest, 0, len(entries))
	openAttempt := ""
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return fmt.Errorf("unexpected entry in attempts: %s", entry.Name())
		}
		attemptDir := filepath.Join(dir, "attempts", entry.Name())
		var attempt AttemptManifest
		if err := s.readYAML(filepath.Join(attemptDir, "attempt.yaml"), &attempt); err != nil {
			return err
		}
		if attempt.Metadata.Name != entry.Name() {
			return fmt.Errorf("attempt directory name mismatch")
		}
		if err := attempt.validate(goal); err != nil {
			return err
		}
		if err := validateAttemptFiles(attemptDir); err != nil {
			return err
		}
		if err := validateAttemptArtifacts(attemptDir, attempt.Status.Artifacts); err != nil {
			return err
		}
		snapshot, ok := history.Snapshots[attempt.Spec.CriteriaRevision]
		if !ok {
			return fmt.Errorf("attempt %q refers to a missing criteria snapshot", attempt.Metadata.Name)
		}
		criteriaDigest, err := criteriaPortableDigest(snapshot)
		if err != nil || attempt.Spec.CriteriaDigest != criteriaDigest {
			return fmt.Errorf("attempt %q criteria digest does not match immutable snapshot", attempt.Metadata.Name)
		}
		if attempt.Status.State == "closed" {
			if err := validateReviewAgainstCriteria(attempt.Status.Review, snapshot); err != nil {
				return err
			}
		}
		if attempt.Status.State == "open" {
			if openAttempt != "" {
				return fmt.Errorf("multiple open attempts")
			}
			openAttempt = attempt.Metadata.Name
		}
		attempts = append(attempts, attempt)
	}
	if openAttempt != goal.Status.ActiveAttemptID {
		return fmt.Errorf("active attempt does not match open attempt set")
	}
	if openAttempt != "" {
		if goal.Status.Outcome != "open" || goal.Status.Execution != "active" {
			return fmt.Errorf("open attempt requires open active goal")
		}
	}
	return validateProspectiveRecord(goal, criteria, attempts)
}

// renderREADMEFromCanonical renders the README projection directly from the
// canonical record without requiring an already-validated load.
func (s *Store) renderREADMEFromCanonical(dir string) ([]byte, error) {
	goal, err := s.readGoalManifest(dir)
	if err != nil {
		return nil, err
	}
	var criteria CriteriaManifest
	if err := s.readYAML(filepath.Join(dir, "criteria.yaml"), &criteria); err != nil {
		return nil, err
	}
	if err := criteria.validate(goal); err != nil {
		return nil, err
	}
	history, err := s.loadCriteriaHistory(dir, goal, criteria)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(filepath.Join(dir, "attempts"))
	if os.IsNotExist(err) {
		entries = nil
	} else if err != nil {
		return nil, err
	}
	attempts := make([]AttemptManifest, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			return nil, fmt.Errorf("unexpected entry in attempts: %s", entry.Name())
		}
		var attempt AttemptManifest
		if err := s.readYAML(filepath.Join(entryDir(dir, entry.Name()), "attempt.yaml"), &attempt); err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	_ = history
	return renderREADME(goal, criteria, attempts, defaultOutputLimit)
}

func entryDir(dir string, name string) string {
	return filepath.Join(dir, "attempts", name)
}

// writeStagedFile writes a staged after-image atomically WITHOUT firing the
// before-rename hook. The hook is reserved for canonical target renames so
// fault injection observes only the real publication boundaries.
func writeStagedFile(path string, content []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".goal-write-")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	installed := false
	defer func() {
		_ = temporary.Close()
		if !installed {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		return err
	}
	if _, err := temporary.Write(content); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return err
	}
	installed = true
	return nil
}

// publicationStateName maps internal states to stable public names.
func publicationStateName(state publicationState) string {
	switch state {
	case publicationStaged:
		return v1alpha1.PublicationStaged
	case publicationPartial:
		return v1alpha1.PublicationPartial
	case publicationConflict:
		return v1alpha1.PublicationConflict
	case publicationDiscardable:
		return v1alpha1.PublicationDiscardable
	default:
		return v1alpha1.PublicationStable
	}
}
