package fsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	v1alpha1 "git.alwaldend.com/alwaldend/src/projects/goal/api/v1alpha1"
	"github.com/goccy/go-yaml"
)

type criteriaHistory struct {
	Snapshots map[uint64]CriteriaManifest
	Digests   map[uint64]string
}

func cleanupTemporaryResidue(root string) error {
	if err := cleanupTemporaryFiles(root, ".goal-write-"); err != nil {
		return err
	}
	if err := cleanupTemporaryFiles(
		filepath.Join(root, "criteria-revisions"),
		".goal-write-",
		".goal-immutable-",
	); err != nil {
		return err
	}

	attemptsRoot := filepath.Join(root, "attempts")
	attemptsInfo, err := os.Lstat(attemptsRoot)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect attempts directory: %w", err)
	}
	if !attemptsInfo.IsDir() || attemptsInfo.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	entries, err := os.ReadDir(attemptsRoot)
	if err != nil {
		return fmt.Errorf("inspect temporary attempt residue: %w", err)
	}
	for _, entry := range entries {
		path := filepath.Join(attemptsRoot, entry.Name())
		if internalTemporaryName(entry.Name(), ".goal-attempt-") {
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("inspect temporary attempt %q: %w", entry.Name(), err)
			}
			if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf(
					"temporary attempt residue %q has an unexpected type",
					entry.Name(),
				)
			}
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("remove temporary attempt %q: %w", entry.Name(), err)
			}
			continue
		}
		if strings.HasPrefix(entry.Name(), ".") || !entry.IsDir() {
			continue
		}
		if err := cleanupTemporaryFiles(path, ".goal-write-"); err != nil {
			return err
		}
		if err := cleanupTemporaryFiles(
			filepath.Join(path, "evidence"),
			".goal-write-",
		); err != nil {
			return err
		}
	}
	return nil
}

func cleanupTemporaryFiles(directory string, prefixes ...string) error {
	directoryInfo, err := os.Lstat(directory)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect temporary-file directory: %w", err)
	}
	if !directoryInfo.IsDir() || directoryInfo.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("inspect temporary files in %s: %w", filepath.Base(directory), err)
	}
	for _, entry := range entries {
		matched := false
		for _, prefix := range prefixes {
			if internalTemporaryName(entry.Name(), prefix) {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect temporary file %q: %w", entry.Name(), err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf(
				"temporary file residue %q has an unexpected type",
				entry.Name(),
			)
		}
		if err := os.Remove(filepath.Join(directory, entry.Name())); err != nil {
			return fmt.Errorf("remove temporary file %q: %w", entry.Name(), err)
		}
	}
	return nil
}

func internalTemporaryName(name string, prefix string) bool {
	if !strings.HasPrefix(name, prefix) || len(name) == len(prefix) {
		return false
	}
	for _, char := range name[len(prefix):] {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') {
			return false
		}
	}
	return true
}

func marshalYAML(value any) ([]byte, error) {
	content, err := yaml.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(content) > maxManifestBytes {
		return nil, fmt.Errorf("encoded manifest exceeds %d bytes", maxManifestBytes)
	}
	return content, nil
}

func criteriaSnapshotPath(dir string, revision uint64) string {
	return filepath.Join(dir, "criteria-revisions", strconv.FormatUint(revision, 10)+".yaml")
}

func (s *Store) installImmutableFile(path string, content []byte, mode os.FileMode) error {
	if existing, err := readRegularFile(path, maxManifestBytes); err == nil {
		if digestBytes(existing) == digestBytes(content) {
			return nil
		}
		return fmt.Errorf("immutable file %s already exists with a different digest", filepath.Base(path))
	} else if !os.IsNotExist(err) {
		return err
	}
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".goal-immutable-")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryPath)
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
		return err
	}
	return nil
}

func (s *Store) loadCriteriaHistory(
	dir string,
	goal GoalManifest,
	current CriteriaManifest,
) (criteriaHistory, error) {
	root := filepath.Join(dir, "criteria-revisions")
	entries, err := os.ReadDir(root)
	if err != nil {
		return criteriaHistory{}, fmt.Errorf("read criteria revision snapshots: %w", err)
	}
	if len(entries) == 0 || len(entries) > maxCriteriaRevisions {
		return criteriaHistory{}, fmt.Errorf(
			"criteria revision snapshot cardinality must be between 1 and %d",
			maxCriteriaRevisions,
		)
	}
	history := criteriaHistory{
		Snapshots: make(map[uint64]CriteriaManifest, len(entries)),
		Digests:   make(map[uint64]string, len(entries)),
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") ||
			!strings.HasSuffix(entry.Name(), ".yaml") {
			return criteriaHistory{}, fmt.Errorf("invalid criteria snapshot entry %q", entry.Name())
		}
		revision, err := strconv.ParseUint(strings.TrimSuffix(entry.Name(), ".yaml"), 10, 64)
		if err != nil || revision == 0 || entry.Name() != strconv.FormatUint(revision, 10)+".yaml" {
			return criteriaHistory{}, fmt.Errorf("invalid criteria snapshot name %q", entry.Name())
		}
		var snapshot CriteriaManifest
		path := filepath.Join(root, entry.Name())
		if err := s.readYAML(path, &snapshot); err != nil {
			return criteriaHistory{}, err
		}
		if err := snapshot.validateSnapshot(goal.Metadata.Name, revision); err != nil {
			return criteriaHistory{}, fmt.Errorf("criteria snapshot %d: %w", revision, err)
		}
		content, err := readRegularFile(path, maxManifestBytes)
		if err != nil {
			return criteriaHistory{}, err
		}
		history.Snapshots[revision] = snapshot
		history.Digests[revision] = digestBytes(content)
	}
	for revision := uint64(1); revision <= goal.Status.CriteriaRevision; revision++ {
		if _, ok := history.Snapshots[revision]; !ok {
			return criteriaHistory{}, fmt.Errorf("missing immutable criteria snapshot %d", revision)
		}
	}
	if len(history.Snapshots) != int(goal.Status.CriteriaRevision) {
		return criteriaHistory{}, fmt.Errorf("criteria snapshots extend beyond current revision")
	}
	latest := history.Snapshots[goal.Status.CriteriaRevision]
	if !reflect.DeepEqual(latest, current) {
		return criteriaHistory{}, fmt.Errorf("criteria.yaml differs from its immutable current snapshot")
	}
	currentBytes, err := readRegularFile(filepath.Join(dir, "criteria.yaml"), maxManifestBytes)
	if err != nil {
		return criteriaHistory{}, err
	}
	if digestBytes(currentBytes) != history.Digests[goal.Status.CriteriaRevision] {
		return criteriaHistory{}, fmt.Errorf("criteria.yaml byte digest differs from its immutable snapshot")
	}
	return history, nil
}

func computeArtifactManifest(attemptDir string) (ArtifactManifest, error) {
	plan, err := readMarkdownFile(filepath.Join(attemptDir, "plan.md"), maxPlanResultBytes)
	if err != nil {
		return ArtifactManifest{}, fmt.Errorf("read plan artifact: %w", err)
	}
	result, err := readMarkdownFile(filepath.Join(attemptDir, "result.md"), maxPlanResultBytes)
	if err != nil {
		return ArtifactManifest{}, fmt.Errorf("read result artifact: %w", err)
	}
	entries, err := os.ReadDir(filepath.Join(attemptDir, "evidence"))
	if err != nil {
		return ArtifactManifest{}, fmt.Errorf("read evidence directory: %w", err)
	}
	if len(entries) > maxEvidenceFiles {
		return ArtifactManifest{}, fmt.Errorf("evidence cardinality exceeds %d", maxEvidenceFiles)
	}
	manifest := ArtifactManifest{
		PlanDigest:   digestBytes(plan),
		ResultDigest: digestBytes(result),
		Evidence:     make([]ArtifactDigest, 0, len(entries)),
	}
	for _, entry := range entries {
		if entry.IsDir() || !safeEvidenceName(entry.Name()) {
			return ArtifactManifest{}, fmt.Errorf("invalid evidence entry %q", entry.Name())
		}
		content, err := readMarkdownFile(
			filepath.Join(attemptDir, "evidence", entry.Name()),
			maxEvidenceFileBytes,
		)
		if err != nil {
			return ArtifactManifest{}, err
		}
		manifest.Evidence = append(manifest.Evidence, ArtifactDigest{
			Path:   "evidence/" + entry.Name(),
			Digest: digestBytes(content),
		})
	}
	sort.Slice(manifest.Evidence, func(i, j int) bool {
		return manifest.Evidence[i].Path < manifest.Evidence[j].Path
	})
	return manifest, manifest.Validate()
}

func validateAttemptArtifacts(attemptDir string, expected ArtifactManifest) error {
	actual, err := computeArtifactManifest(attemptDir)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(actual, expected) {
		return fmt.Errorf("attempt artifact digest manifest does not match frozen files")
	}
	return nil
}

func goalStateDigest(goal GoalManifest) (string, error) {
	return v1alpha1.GoalStateDigest(v1alpha1.Goal(goal))
}

func criteriaPortableDigest(criteria CriteriaManifest) (string, error) {
	return v1alpha1.CriteriaDigest(v1alpha1.GoalCriteria(criteria))
}

func validateReviewAgainstCriteria(
	review CloseReview,
	criteria CriteriaManifest,
) error {
	return v1alpha1.ValidateReviewAgainstCriteria(
		review,
		v1alpha1.GoalCriteria(criteria),
	)
}

func reviewAcceptsRequired(review CloseReview, criteria CriteriaManifest) error {
	return v1alpha1.ReviewAcceptsRequired(
		review,
		v1alpha1.GoalCriteria(criteria),
	)
}
