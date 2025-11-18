package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"git.alwaldend.com/src/tools/git/main/proto/contracts"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"google.golang.org/protobuf/encoding/protojson"
)

type GitInfo struct{}

func NewGitInfo() *GitInfo {
	return &GitInfo{}
}

type GitInfoGenerateOptions struct {
	Output   string
	GitRoot  string
	InMemory bool
}

func (self *GitInfo) Generate(ctx context.Context, opts *GitInfoGenerateOptions) error {
	res := &contracts.GitInfo{}
	repo, err := self.openRepo(opts.GitRoot, opts.InMemory)
	if err != nil {
		return fmt.Errorf("could not open repo '%s': %w", opts.GitRoot, err)
	}
	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("could not get remotes: %w", err)
	}
	for _, remote := range remotes {
		config := remote.Config()
		remote := &contracts.GitRemote{Name: config.Name, Urls: config.URLs}
		res.Remotes = append(res.Remotes, remote)
	}
	tagIter, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("could not create tag iterator ref: %w", err)
	}
	res.Commits = make(map[string]*contracts.GitCommit)
	res.Tags = make(map[string]*contracts.GitTag)
	commitTags := make(map[string][]string)

	err = tagIter.ForEach(func(tagRef *plumbing.Reference) error {
		var err error
		var tag *contracts.GitTag
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			tag, err = self.parseTag(repo, tagRef)
		}
		if err != nil {
			return fmt.Errorf("could not parse tag %s: %w", tagRef.Name(), err)
		}
		res.Tags[tag.Name] = tag
		commitTags[tag.Target] = append(commitTags[tag.Target], tag.Name)
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not parse tags: %w", err)
	}
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return fmt.Errorf("could not create commit iterator: %w", err)
	}
	err = commitIter.ForEach(func(commitRef *object.Commit) error {
		tags, _ := commitTags[commitRef.Hash.String()]
		var err error
		var commit *contracts.GitCommit
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			commit, err = self.parseCommit(ctx, commitRef, tags)
		}
		if err != nil {
			return fmt.Errorf("could not parse commit %s: %w", commitRef.Hash.String(), err)
		}
		res.CommitsOrder = append(res.CommitsOrder, commit.Hash)
		res.Commits[commit.Hash] = commit
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not parse commits: %w", err)
	}
	data, err := protojson.MarshalOptions{}.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not marshal git info %v: %w", res, err)
	}
	if err := os.WriteFile(opts.Output, data, 0o444); err != nil {
		return fmt.Errorf("could not write data to file %s: %w", opts.Output, err)
	}
	return nil
}

func (self *GitInfo) parseTag(repo *git.Repository, tagRef *plumbing.Reference) (*contracts.GitTag, error) {
	tag := &contracts.GitTag{
		Name:   strings.TrimPrefix(tagRef.Name().String(), "refs/tags/"),
		Target: tagRef.Hash().String(),
	}

	if tagObj, err := repo.TagObject(tagRef.Hash()); err == nil {
		annotated := &contracts.GitTagAnnotated{
			Message:      tagObj.Message,
			PgpSignature: tagObj.PGPSignature,
			Hash:         tagRef.Hash().String(),
			Tagger: &contracts.GitSignature{
				Name:  tagObj.Tagger.Name,
				Email: tagObj.Tagger.Email,
				When:  tagObj.Tagger.When.Unix(),
			},
		}
		tag.Target = tagObj.Target.String()
		tag.Name = tagObj.Name
		tag.Annotated = annotated
	}
	return tag, nil
}

func (self *GitInfo) openRepo(gitRoot string, inMemory bool) (*git.Repository, error) {
	var repo *git.Repository
	var err error
	if inMemory {
		repo, err = git.Clone(
			memory.NewStorage(),
			memfs.New(),
			&git.CloneOptions{URL: gitRoot},
		)
		if err != nil {
			return nil, fmt.Errorf("could not create in-memory repo clone: %w", err)
		}
	} else {
		repo, err = git.PlainOpen(gitRoot)
		if err != nil {
			return nil, fmt.Errorf("could not plain open: %w", err)
		}
	}
	return repo, nil
}

func (self *GitInfo) parseCommit(ctx context.Context, commitRef *object.Commit, tags []string) (*contracts.GitCommit, error) {
	commitHash := commitRef.Hash.String()
	commit := &contracts.GitCommit{
		Hash: commitHash,
		Author: &contracts.GitSignature{
			Name:  commitRef.Author.Name,
			Email: commitRef.Author.Email,
			When:  commitRef.Author.When.Unix(),
		},
		Committer: &contracts.GitSignature{
			Name:  commitRef.Committer.Name,
			Email: commitRef.Committer.Email,
			When:  commitRef.Committer.When.Unix(),
		},
		PgpSignature: commitRef.PGPSignature,
		MergeTag:     commitRef.MergeTag,
		Message:      commitRef.Message,
	}
	commit.Tags = tags
	return commit, nil
}
