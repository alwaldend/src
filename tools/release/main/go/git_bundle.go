package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	gitBundleSignature = "# v2 git bundle\n"
	releaseGitBranch   = plumbing.ReferenceName("refs/heads/release")
)

var gitObjectIDPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

func openGitBundle(path string) (_ *git.Repository, returnErr error) {
	bundle, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open Git bundle %s: %w", path, err)
	}
	defer func() {
		if err := bundle.Close(); err != nil && returnErr == nil {
			returnErr = fmt.Errorf("could not close Git bundle %s: %w", path, err)
		}
	}()

	reader := bufio.NewReader(bundle)
	signature, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("could not read Git bundle signature: %w", err)
	}
	if signature != gitBundleSignature {
		return nil, fmt.Errorf("unsupported Git bundle signature %q", signature)
	}

	references, err := readGitBundleReferences(reader)
	if err != nil {
		return nil, err
	}
	head, ok := references[releaseGitBranch]
	if !ok {
		return nil, fmt.Errorf("Git bundle is missing %s", releaseGitBranch)
	}

	storage := memory.NewStorage()
	repo, err := git.InitWithOptions(
		storage,
		nil,
		git.InitOptions{DefaultBranch: releaseGitBranch},
	)
	if err != nil {
		return nil, fmt.Errorf("could not initialize Git bundle storage: %w", err)
	}
	if err := packfile.UpdateObjectStorage(storage, reader); err != nil {
		return nil, fmt.Errorf("could not unpack Git bundle: %w", err)
	}
	for name, hash := range references {
		if err := storage.SetReference(plumbing.NewHashReference(name, hash)); err != nil {
			return nil, fmt.Errorf("could not set Git bundle reference %s: %w", name, err)
		}
	}
	if _, err := repo.CommitObject(head); err != nil {
		return nil, fmt.Errorf("could not resolve Git bundle head %s: %w", head, err)
	}
	return repo, nil
}

func readGitBundleReferences(
	reader *bufio.Reader,
) (map[plumbing.ReferenceName]plumbing.Hash, error) {
	references := make(map[plumbing.ReferenceName]plumbing.Hash)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("could not read Git bundle header: %w", err)
		}
		if line == "\n" {
			return references, nil
		}
		if strings.HasPrefix(line, "-") {
			return nil, fmt.Errorf("thin Git bundles are not supported")
		}

		fields := strings.SplitN(strings.TrimSuffix(line, "\n"), " ", 2)
		if len(fields) != 2 || !gitObjectIDPattern.MatchString(fields[0]) {
			return nil, fmt.Errorf("invalid Git bundle reference %q", line)
		}
		name := plumbing.ReferenceName(fields[1])
		if err := name.Validate(); err != nil || name == plumbing.HEAD {
			return nil, fmt.Errorf("invalid Git bundle reference name %q", name)
		}
		if _, exists := references[name]; exists {
			return nil, fmt.Errorf("duplicate Git bundle reference %s", name)
		}
		references[name] = plumbing.NewHash(fields[0])
	}
}
