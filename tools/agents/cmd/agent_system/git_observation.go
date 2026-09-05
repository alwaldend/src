package agent_system

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

const gitOutputLimit = 65536

type gitCommand func(context.Context, string, ...string) ([]byte, error)

func observeGit(root string, now func() time.Time, run gitCommand) *v1alpha1.CapsuleGitObservation {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	observation := &v1alpha1.CapsuleGitObservation{}
	inside, err := run(ctx, root, "rev-parse", "--is-inside-work-tree")
	if err != nil || strings.TrimSpace(string(inside)) != "true" {
		observation.Unavailable = append(observation.Unavailable,
			"workspace Git metadata unavailable (missing Git, non-worktree, command failure, or time/output limit)")
	} else {
		revision, err := run(ctx, root, "rev-parse", "--verify", "HEAD")
		if err == nil && validGitRevision(strings.TrimSpace(string(revision))) {
			observation.Revision = strings.TrimSpace(string(revision))
		} else {
			observation.Unavailable = append(observation.Unavailable,
				"Git HEAD unavailable (unborn branch, command failure, or time/output limit)")
		}
		status, err := run(ctx, root, "status", "--porcelain=v1", "-z", "--untracked-files=normal", "--ignore-submodules=none", "--", ".")
		if err == nil {
			dirty := len(status) != 0
			observation.Dirty = &dirty
		} else {
			observation.Unavailable = append(observation.Unavailable,
				"Git dirty state unavailable (command failure or time/output limit)")
		}
	}
	observation.ObservedAt = now().UTC().Format(time.RFC3339Nano)
	return observation
}

func validGitRevision(revision string) bool {
	if len(revision) != 40 && len(revision) != 64 {
		return false
	}
	for _, char := range revision {
		if !(char >= '0' && char <= '9' || char >= 'a' && char <= 'f') {
			return false
		}
	}
	return true
}

func applyGitObservation(capsule *v1alpha1.ContextCapsule, observation *v1alpha1.CapsuleGitObservation) {
	capsule.Identity.Git = observation
	if capsule.Identity.RevisionSource != "caller-declared" && observation.Revision != "" {
		capsule.Identity.Revision = observation.Revision
		capsule.Identity.RevisionSource = "observed-git-head"
	}
	if capsule.Identity.DirtyInputsSource != "caller-declared" && observation.Dirty != nil {
		capsule.Identity.DirtyInputs = *observation.Dirty
		capsule.Identity.DirtyInputsSource = "observed-git-status"
	}
}

// Prevent ambient Git variables from redirecting the observation to another
// checkout. Disable the optional index refresh and configured fsmonitor hooks.
// Never include Git stderr or changed path names in the capsule.
func runGitCommand(ctx context.Context, root string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", append([]string{
		"--no-optional-locks", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-C", root,
	}, args...)...)
	command.Env = []string{}
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "GIT_") {
			command.Env = append(command.Env, entry)
		}
	}
	command.WaitDelay = 100 * time.Millisecond
	output := &boundedGitOutput{}
	command.Stdout = output
	command.Stderr = &boundedGitOutput{}
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("Git observation command failed")
	}
	return output.buffer.Bytes(), nil
}

type boundedGitOutput struct {
	buffer bytes.Buffer
}

func (output *boundedGitOutput) Write(content []byte) (int, error) {
	if len(content) > gitOutputLimit-output.buffer.Len() {
		return 0, fmt.Errorf("Git observation exceeds output limit")
	}
	return output.buffer.Write(content)
}
