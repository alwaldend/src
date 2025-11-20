package main

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/src/tools/release/contracts"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Generator struct{}

type GeneratorParams struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

// Generate a manifest file
func (self *Generator) Generate(
	items []string,
	manifests []string,
	output string,
	gitRoot string,
	marshalOptions *protojson.MarshalOptions,
) error {
	release := &contracts.Release{}
	for _, item := range items {
		if err := self.addItem(release, item); err != nil {
			return fmt.Errorf("could not add item %s to release: %w", item, err)
		}
	}
	for _, manifest := range manifests {
		if err := self.addManifest(release, manifest); err != nil {
			return fmt.Errorf("could not add manifest %s to release: %w", manifest, err)
		}
	}
	if release.Git != nil {
		err := self.addGit(release, gitRoot)
		if err != nil {
			return fmt.Errorf("could not add git info to release: %w", err)
		}
	}
	data, err := marshalOptions.Marshal(release)
	if err != nil {
		return fmt.Errorf("could not marshal release %v: %w", release, err)
	}
	if err := os.WriteFile(output, data, 0o444); err != nil {
		return fmt.Errorf("could not write data to file %s: %w", output, err)
	}
	return nil
}

// Add git information in-place
func (self *Generator) addGit(release *contracts.Release, gitRoot string) error {
	if release.Git == nil {
		return fmt.Errorf("release is missing git info: %v", release)
	}
	if release.Project == nil {
		return fmt.Errorf("missing project info: %v", release)
	}
	repo, err := git.PlainOpenWithOptions(gitRoot, &git.PlainOpenOptions{})
	if err != nil {
		return fmt.Errorf("could not open repo %s: %w", gitRoot, err)
	}
	rev, err := repo.ResolveRevision(plumbing.Revision(release.Git.Revision))
	if err != nil {
		return fmt.Errorf("could not resolve revision %s: %w", release.Git.Revision, err)
	}
	release.Git.Revision = rev.String()
	tagIter, err := repo.Tags()
	if err != nil {
		return fmt.Errorf("could not create tag iterator ref: %w", err)
	}
	hashesToTags := make(map[string][]string)
	err = tagIter.ForEach(func(tag *plumbing.Reference) error {
		idx := tag.Hash().String()
		hashesToTags[idx] = append(hashesToTags[idx], tag.String())
		return nil
	})
	revTags, ok := hashesToTags[rev.String()]
	if ok {
		release.Git.Tags = append(release.Git.Tags, revTags...)
	}
	filter := func(path string) bool {
		res := strings.HasPrefix(path, release.Project.Subdir)
		return res
	}
	commitIter, err := repo.Log(&git.LogOptions{From: *rev, PathFilter: filter})
	if err != nil {
		return fmt.Errorf("could not create commit iterator: %w", err)
	}
	err = commitIter.ForEach(func(commit *object.Commit) error {
		_, hasTags := hashesToTags[commit.Hash.String()]
		if commit.Hash != *rev && hasTags {
			commitIter.Close()
			return nil
		}
		release.Git.Commits = append(release.Git.Commits, commit.Hash.String())
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not iterate over commits: %w", err)
	}
	return nil
}

// Add manifest information in-place
func (self *Generator) addManifest(release *contracts.Release, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}
	fileManifest := &contracts.Release{}
	if err := protojson.Unmarshal(content, fileManifest); err != nil {
		return fmt.Errorf("could not parse file %s: %w", path, err)
	}
	proto.Merge(release, fileManifest)
	return nil
}

// Add item information in-place
func (self *Generator) addItem(release *contracts.Release, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read item manifest %s: %w", path, err)
	}
	item := &contracts.ReleaseItem{}
	if err := protojson.Unmarshal(content, item); err != nil {
		return fmt.Errorf("could not unmarshal item manifest %s: %w", path, err)
	}
	if item.File == nil || item.File.LocalPath == "" {
		return fmt.Errorf("item manifest %s is missing file.local_path", path)
	}
	stat, err := os.Stat(item.File.LocalPath)
	if err != nil {
		return fmt.Errorf("could not stat local file %s: %w", path, err)
	}
	item.File.Name = filepath.Base(item.File.LocalPath)
	item.File.Size = stat.Size()
	for _, hashType := range []crypto.Hash{crypto.SHA256, crypto.MD5} {
		hashObj := hashType.New()
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("could not open file %s: %w", path, err)
		}
		defer file.Close()
		if _, err := io.Copy(hashObj, file); err != nil {
			return fmt.Errorf("could not copy file %s to hash: %w", path, err)
		}
		item.File.Hashes = append(item.File.Hashes, &contracts.ReleaseHash{
			Algo:    hashType.String(),
			Content: fmt.Sprintf("%x", hashObj.Sum(nil)),
		})
	}
	release.Items = append(release.Items, item)
	return nil
}
