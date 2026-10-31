package gitinfo

import (
	"context"
	"fmt"
	"strings"
	"sync"

	contracts "git.alwaldend.com/alwaldend/src/projects/al/api/git"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitInfo projects Git tags and commits for optional release changelogs.
type GitInfo struct{}

// ParseTags resolves tags and indexes them by target commit. Storage reads remain
// sequential because the filesystem store lazily updates shared pack indexes.
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
	err = tagIter.ForEach(func(tagRef *plumbing.Reference) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("could not parse tag %s: %w", tagRef.Name(), err)
		}
		tag, err := self.ParseTag(repo, tagRef)
		if err != nil {
			return fmt.Errorf("could not parse tag %s: %w", tagRef.Name(), err)
		}
		res.Tags[tag.Name] = tag
		commitTags[tag.Target] = append(commitTags[tag.Target], tag.Name)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not iterate over tags: %w", err)
	}
	return res, commitTags, nil
}

// ParseTag resolves one lightweight or annotated Git tag.
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

// ParseCommit projects one Git commit into the shared metadata schema.
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

// commitChanges reads a commit tree difference while serializing storage access.
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
