package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"git.alwaldend.com/alwaldend/src/tools/git/main/proto/contracts"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
	repo, err := self.OpenRepo(opts.GitRoot, opts.InMemory)
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
	tagInfo, commitTags, err := self.ParseTags(ctx, repo)
	if err != nil {
		return fmt.Errorf("could not parse tags: %w", err)
	}
	proto.Merge(res, tagInfo)
	commitInfo, err := self.ParseCommits(ctx, repo, commitTags)
	if err != nil {
		return fmt.Errorf("could not parse commits: %w", err)
	}
	proto.Merge(res, commitInfo)
	data, err := protojson.MarshalOptions{}.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not marshal git info %v: %w", res, err)
	}
	if err := os.WriteFile(opts.Output, data, 0o444); err != nil {
		return fmt.Errorf("could not write data to file %s: %w", opts.Output, err)
	}
	return nil
}

func (self *GitInfo) ParseTags(
	ctx context.Context,
	repo *git.Repository,
) (*contracts.GitInfo, map[string][]string, error) {
	tagIter, err := repo.Tags()
	if err != nil {
		return nil, nil, fmt.Errorf("could not create tag iterator ref: %w", err)
	}
	res := &contracts.GitInfo{
		Tags: make(map[string]*contracts.GitTag),
	}
	commitTags := make(map[string][]string)
	tagsChannel := make(chan *contracts.GitTag)
	tagsErrorChannel := make(chan error)
	tagsCount := 0
	err = tagIter.ForEach(func(tagRef *plumbing.Reference) error {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("could not parse tag %s: %w", tagRef.Name(), ctx.Err())
		default:
			tagsCount += 1
			go func() {
				tag, err := self.ParseTag(repo, tagRef)
				if err == nil {
					tagsChannel <- tag
				} else {
					tagsErrorChannel <- err
				}
			}()
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not iterate over tags: %w", err)
	}
	for range tagsCount {
		select {
		case tag := <-tagsChannel:
			res.Tags[tag.Name] = tag
			commitTags[tag.Target] = append(commitTags[tag.Target], tag.Name)
		case curErr := <-tagsErrorChannel:
			err = errors.Join(err, curErr)
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse tags: %w", err)
	}
	return res, commitTags, nil
}

func (self *GitInfo) ParseTag(repo *git.Repository, tagRef *plumbing.Reference) (*contracts.GitTag, error) {
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

func (self *GitInfo) OpenRepo(gitRoot string, inMemory bool) (*git.Repository, error) {
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

func (self *GitInfo) ParseCommits(ctx context.Context, repo *git.Repository, commitTags map[string][]string) (*contracts.GitInfo, error) {
	res := &contracts.GitInfo{
		Commits: make(map[string]*contracts.GitCommit),
	}
	commitIter, err := repo.CommitObjects()
	if err != nil {
		return nil, fmt.Errorf("could not create commit iterator: %w", err)
	}
	mutex := &sync.Mutex{}
	commitsChannel := make(chan *contracts.GitCommit)
	commitsErrorChannel := make(chan error)
	commitCount := 0
	err = commitIter.ForEach(func(commitRef *object.Commit) error {
		select {
		case <-ctx.Done():
			return fmt.Errorf("could not parse commit %s: %w", commitRef.Hash.String(), ctx.Err())
		default:
			commitCount += 1
			tags, _ := commitTags[commitRef.Hash.String()]
			go func() {
				commit, err := self.ParseCommit(ctx, mutex, commitRef, tags, false)
				if err == nil {
					commitsChannel <- commit
				} else {
					commitsErrorChannel <- err
				}
			}()
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not iterate over commits: %w", err)
	}
	for range commitCount {
		select {
		case commit := <-commitsChannel:
			res.CommitsOrder = append(res.CommitsOrder, commit.Hash)
			res.Commits[commit.Hash] = commit
		case curErr := <-commitsErrorChannel:
			err = errors.Join(err, curErr)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not parse commits: %w", err)
	}
	return res, nil
}

func (self *GitInfo) ParseCommit(
	ctx context.Context,
	mutex *sync.Mutex,
	commitRef *object.Commit,
	tags []string,
	addChangedFiles bool,
) (*contracts.GitCommit, error) {
	commitHash := commitRef.Hash.String()
	commit := &contracts.GitCommit{
		Hash:         commitHash,
		ChangedFiles: make(map[string]bool),
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
	if addChangedFiles {
		changes, err := self.commitChanges(ctx, mutex, commitRef)
		if err != nil {
			return nil, fmt.Errorf("could not get commit changes: %w", err)
		}
		for _, change := range changes {
			commit.ChangedFiles[change.From.Name] = false
			commit.ChangedFiles[change.To.Name] = false
		}
	}
	return commit, nil
}

func (self *GitInfo) commitChanges(ctx context.Context, mutex *sync.Mutex, commitRef *object.Commit) (object.Changes, error) {
	mutex.Lock()
	defer mutex.Unlock()
	fromTree, err := commitRef.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get commit tree: %w", err)
	}
	toTree := &object.Tree{}
	if commitRef.NumParents() != 0 {
		firstParent, err := commitRef.Parents().Next()
		if err != nil {
			return nil, fmt.Errorf("could not get first parent: %w", err)
		}
		toTree, err = firstParent.Tree()
		if err != nil {
			return nil, fmt.Errorf("could not get first parent tree: %w", err)
		}
	}
	changes, err := object.DiffTreeWithOptions(ctx, toTree, fromTree, &object.DiffTreeOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get changes: %w", err)
	}
	return changes, nil
}
