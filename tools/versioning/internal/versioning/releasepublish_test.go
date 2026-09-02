package versioning

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestLocalGitPublisherPreflightRequiresRemote(t *testing.T) {
	publisher := &LocalGitReleasePublisher{Git: nil, Remote: "origin"}
	if err := publisher.Preflight(context.Background()); err == nil {
		t.Fatal("Preflight() accepted a nil runner")
	}
	publisher = &LocalGitReleasePublisher{
		Git:    &fakeRunner{responses: map[string]string{}},
		Remote: "",
	}
	if err := publisher.Preflight(context.Background()); err == nil {
		t.Fatal("Preflight() accepted an empty remote")
	}
}

func TestLocalGitPublisherRefusesAtomicWithoutSupport(t *testing.T) {
	pushKey := "push origin " + abc123() + ":refs/heads/releases/2026.35 " +
		abc123() + ":refs/tags/v2026.35.0 --atomic"
	runner := &fakeRunner{responses: map[string]string{
		"remote get-url origin": "file:///tmp/remote",
		"ls-remote origin":      abc123() + "\trefs/heads/master\n",
	}, errors: map[string]error{
		pushKey: errors.New("atomic push not supported"),
	}}
	publisher := &LocalGitReleasePublisher{Git: runner, Remote: "origin"}
	plan, err := BuildReleaseRefPlan(State{
		Version: "2026.35.0", Channel: "release", Branch: "releases/2026.35",
		Commit: abc123(), TreeState: "clean",
	}, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = "lease"
	_, err = PublishReleaseRefs(context.Background(), publisher, plan)
	if err == nil || !strings.Contains(err.Error(), "atomic multi-ref") {
		t.Fatalf("PublishReleaseRefs() error = %v, want atomicity refusal", err)
	}
	for _, call := range runner.calls {
		if strings.HasPrefix(call, "push ") && !strings.Contains(call, "--atomic") {
			t.Fatalf("publisher issued a non-atomic push: %v", runner.calls)
		}
	}
}

func TestLocalGitPublisherVerifiesRemoteAfterPublish(t *testing.T) {
	plan, err := BuildReleaseRefPlan(State{
		Version: "2026.35.0", Channel: "release", Branch: "releases/2026.35",
		Commit: abc123(), TreeState: "clean",
	}, "master")
	if err != nil {
		t.Fatal(err)
	}
	plan.Lease = "lease"
	pushKey := "push origin " + abc123() + ":refs/heads/releases/2026.35 " +
		abc123() + ":refs/tags/v2026.35.0 --atomic"
	runner := &fakeRunner{
		responses: map[string]string{
			"remote get-url origin": "file:///tmp/remote",
			pushKey:                 "",
		},
		responseLists: map[string][]string{
			"ls-remote origin": {
				abc123() + "\trefs/heads/master\n",
				abc123() + "\trefs/tags/v2026.35.0\n" + abc123() + "\trefs/heads/releases/2026.35\n",
			},
		},
	}
	publisher := &LocalGitReleasePublisher{Git: runner, Remote: "origin"}
	if err := publisher.Preflight(context.Background()); err != nil {
		t.Fatal(err)
	}
	receipt, err := PublishReleaseRefs(context.Background(), publisher, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !receipt.Verified || len(receipt.Refs) != 2 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
}
