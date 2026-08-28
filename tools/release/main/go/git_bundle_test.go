package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func TestOpenGitBundlePreservesHistoryAndTags(t *testing.T) {
	storage := memory.NewStorage()
	firstTree, firstBlob := storeTestTree(t, storage, "first\n")
	firstCommit := storeTestCommit(
		t,
		storage,
		firstTree,
		nil,
		"first release change\n",
		1,
	)
	headTree, headBlob := storeTestTree(t, storage, "head\n")
	headCommit := storeTestCommit(
		t,
		storage,
		headTree,
		[]plumbing.Hash{firstCommit},
		"head release change\n",
		2,
	)

	var pack bytes.Buffer
	encoder := packfile.NewEncoder(&pack, storage, false)
	_, err := encoder.Encode(
		[]plumbing.Hash{
			firstTree,
			firstCommit,
			headTree,
			headCommit,
		},
		0,
	)
	if err != nil {
		t.Fatalf("could not encode test pack: %v", err)
	}

	var bundle bytes.Buffer
	fmt.Fprintf(
		&bundle,
		"%s%s %s\n%s refs/tags/v1\n\n",
		gitBundleSignature,
		headCommit,
		releaseGitBranch,
		firstCommit,
	)
	bundle.Write(pack.Bytes())
	bundlePath := filepath.Join(t.TempDir(), "release.bundle")
	if err := os.WriteFile(bundlePath, bundle.Bytes(), 0o600); err != nil {
		t.Fatalf("could not write test bundle: %v", err)
	}

	repo, err := openGitBundle(bundlePath)
	if err != nil {
		t.Fatalf("could not open test bundle: %v", err)
	}
	head, err := repo.Head()
	if err != nil {
		t.Fatalf("could not resolve test bundle head: %v", err)
	}
	if head.Hash() != headCommit {
		t.Fatalf("head hash is %s, want %s", head.Hash(), headCommit)
	}
	tag, err := repo.Tag("v1")
	if err != nil {
		t.Fatalf("could not resolve test bundle tag: %v", err)
	}
	if tag.Hash() != firstCommit {
		t.Fatalf("tag hash is %s, want %s", tag.Hash(), firstCommit)
	}
	for _, blob := range []plumbing.Hash{firstBlob, headBlob} {
		if _, err := repo.BlobObject(blob); !errors.Is(err, plumbing.ErrObjectNotFound) {
			t.Fatalf("bundle unexpectedly contains blob %s: %v", blob, err)
		}
	}

	commitIterator, err := repo.Log(&git.LogOptions{
		From: headCommit,
		PathFilter: func(path string) bool {
			return path == "release.txt"
		},
	})
	if err != nil {
		t.Fatalf("could not read test bundle log: %v", err)
	}
	defer commitIterator.Close()
	var messages []string
	if err := commitIterator.ForEach(func(commit *object.Commit) error {
		messages = append(messages, commit.Message)
		return nil
	}); err != nil {
		t.Fatalf("could not iterate test bundle log: %v", err)
	}
	wantMessages := []string{"head release change\n", "first release change\n"}
	if fmt.Sprint(messages) != fmt.Sprint(wantMessages) {
		t.Fatalf("messages are %q, want %q", messages, wantMessages)
	}
}

func TestOpenGitBundleRejectsThinBundle(t *testing.T) {
	bundlePath := filepath.Join(t.TempDir(), "thin.bundle")
	content := gitBundleSignature +
		"-0000000000000000000000000000000000000001 prerequisite\n\n"
	if err := os.WriteFile(bundlePath, []byte(content), 0o600); err != nil {
		t.Fatalf("could not write test bundle: %v", err)
	}
	if _, err := openGitBundle(bundlePath); err == nil {
		t.Fatal("openGitBundle accepted a thin bundle")
	}
}

func storeTestTree(
	t *testing.T,
	storage *memory.Storage,
	content string,
) (plumbing.Hash, plumbing.Hash) {
	t.Helper()
	blob := storage.NewEncodedObject()
	blob.SetType(plumbing.BlobObject)
	writer, err := blob.Writer()
	if err != nil {
		t.Fatalf("could not create test blob writer: %v", err)
	}
	if _, err := io.WriteString(writer, content); err != nil {
		t.Fatalf("could not write test blob: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("could not close test blob: %v", err)
	}
	blobHash, err := storage.SetEncodedObject(blob)
	if err != nil {
		t.Fatalf("could not store test blob: %v", err)
	}

	tree := storage.NewEncodedObject()
	if err := (&object.Tree{Entries: []object.TreeEntry{
		{Name: "release.txt", Mode: filemode.Regular, Hash: blobHash},
	}}).Encode(tree); err != nil {
		t.Fatalf("could not encode test tree: %v", err)
	}
	treeHash, err := storage.SetEncodedObject(tree)
	if err != nil {
		t.Fatalf("could not store test tree: %v", err)
	}
	return treeHash, blobHash
}

func storeTestCommit(
	t *testing.T,
	storage *memory.Storage,
	tree plumbing.Hash,
	parents []plumbing.Hash,
	message string,
	second int64,
) plumbing.Hash {
	t.Helper()
	signature := object.Signature{
		Name:  "Release Test",
		Email: "release-test@example.invalid",
		When:  time.Unix(second, 0).UTC(),
	}
	commit := storage.NewEncodedObject()
	if err := (&object.Commit{
		Author:       signature,
		Committer:    signature,
		Message:      message,
		TreeHash:     tree,
		ParentHashes: parents,
	}).Encode(commit); err != nil {
		t.Fatalf("could not encode test commit: %v", err)
	}
	commitHash, err := storage.SetEncodedObject(commit)
	if err != nil {
		t.Fatalf("could not store test commit: %v", err)
	}
	return commitHash
}
