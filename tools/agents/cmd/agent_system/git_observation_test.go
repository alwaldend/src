package agent_system

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	v1alpha1 "git.alwaldend.com/alwaldend/src/tools/agents/api/v1alpha1"
)

func TestGitObservationUsesWorkspaceAndDistinguishesCleanDirtyUnknown(t *testing.T) {
	head := strings.Repeat("a", 40)
	instant := time.Date(2026, 9, 5, 10, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		name        string
		status      string
		statusError bool
		wantDirty   *bool
	}{
		{name: "clean", wantDirty: boolPointer(false)},
		{name: "tracked", status: " M file.go\x00", wantDirty: boolPointer(true)},
		{name: "untracked", status: "?? file.go\x00", wantDirty: boolPointer(true)},
		{name: "failed-status", statusError: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var calls [][]string
			observation := observeGit("/selected/workspace", func() time.Time { return instant },
				func(ctx context.Context, root string, args ...string) ([]byte, error) {
					if root != "/selected/workspace" {
						t.Fatalf("observed process directory instead of requested workspace: %s", root)
					}
					if deadline, exists := ctx.Deadline(); !exists || time.Until(deadline) > 3*time.Second {
						t.Fatal("Git observation has no bounded deadline")
					}
					calls = append(calls, args)
					switch len(calls) {
					case 1:
						return []byte("true\n"), nil
					case 2:
						return []byte(head + "\n"), nil
					default:
						if test.statusError {
							return nil, fmt.Errorf("failure")
						}
						return []byte(test.status), nil
					}
				})
			wantCalls := [][]string{
				{"rev-parse", "--is-inside-work-tree"},
				{"rev-parse", "--verify", "HEAD"},
				{"status", "--porcelain=v1", "-z", "--untracked-files=normal", "--ignore-submodules=none", "--", "."},
			}
			if !reflect.DeepEqual(calls, wantCalls) || observation.Revision != head ||
				!reflect.DeepEqual(observation.Dirty, test.wantDirty) || observation.ObservedAt != instant.Format(time.RFC3339Nano) {
				t.Fatalf("incorrect observation: %+v; calls %v", observation, calls)
			}
			if test.statusError && len(observation.Unavailable) == 0 {
				t.Fatal("unknown state needs reason")
			}
		})
	}
}

func TestGitObservationHandlesUnbornOrUnavailableWorktrees(t *testing.T) {
	for _, unborn := range []bool{false, true} {
		calls := 0
		observation := observeGit("/workspace", time.Now,
			func(_ context.Context, _ string, args ...string) ([]byte, error) {
				calls++
				if unborn && calls == 1 {
					return []byte("true"), nil
				}
				if unborn && args[0] == "status" {
					return []byte("?? new.go\x00"), nil
				}
				return nil, fmt.Errorf("do not expose stderr")
			})
		if observation.Revision != "" || len(observation.Unavailable) == 0 {
			t.Fatalf("invented Git state: %+v", observation)
		}
		if unborn && (calls != 3 || observation.Dirty == nil || !*observation.Dirty) {
			t.Fatalf("unborn HEAD should not suppress dirty observation: %+v", observation)
		}
		if !unborn && (calls != 1 || observation.Dirty != nil) {
			t.Fatalf("unavailable worktree should stop Git queries: %+v", observation)
		}
	}
}

func TestGitObservationPreservesDeclaredOverrides(t *testing.T) {
	for _, declared := range []bool{false, true} {
		capsule := v1alpha1.ContextCapsule{Identity: v1alpha1.CapsuleIdentity{
			Revision: "declared-or-digest", DirtyInputs: true,
			RevisionSource: "input-digest", DirtyInputsSource: "conservative-default",
		}}
		if declared {
			capsule.Identity.RevisionSource = "caller-declared"
			capsule.Identity.DirtyInputsSource = "caller-declared"
		}
		observation := &v1alpha1.CapsuleGitObservation{Revision: strings.Repeat("b", 40), Dirty: boolPointer(false)}
		applyGitObservation(&capsule, observation)
		if capsule.Identity.Git != observation {
			t.Fatal("missing independent Git evidence")
		}
		if declared {
			if capsule.Identity.Revision != "declared-or-digest" || !capsule.Identity.DirtyInputs {
				t.Fatal("overwrote caller declaration")
			}
		} else if capsule.Identity.Revision != observation.Revision || capsule.Identity.DirtyInputs ||
			capsule.Identity.RevisionSource != "observed-git-head" || capsule.Identity.DirtyInputsSource != "observed-git-status" {
			t.Fatalf("defaults did not use observed evidence: %+v", capsule.Identity)
		}
	}
}

func TestGitOutputIsBoundedThroughIOCopy(t *testing.T) {
	var output boundedGitOutput
	_, err := io.Copy(&output, bytes.NewReader(bytes.Repeat([]byte("x"), gitOutputLimit+1)))
	if err == nil || output.buffer.Len() > gitOutputLimit {
		t.Fatalf("unbounded Git output: %d bytes, %v", output.buffer.Len(), err)
	}
}

func boolPointer(value bool) *bool { return &value }
