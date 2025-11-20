package main

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.alwaldend.com/src/tools/release/main/proto/contracts"
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

type GenerateOpts struct {
	MergeItems        []string
	MergeManifests    []string
	OutputManifest    string
	OutputReleasePage string
	OutputFileMode    string
	GitRoot           string
	MarshalOptions    *protojson.MarshalOptions
}

// Generate a manifest file
func (self *Generator) Generate(opts *GenerateOpts) error {
	release, err := self.createRelease(opts)
	if err != nil {
		return fmt.Errorf("could not create release %v: %w", opts, err)
	}
	if opts.OutputManifest != "" {
		err := self.writeMessage(opts, release, opts.OutputManifest)
		if err != nil {
			return fmt.Errorf("could not output release %v: %w", opts.OutputManifest, err)
		}
	}
	if opts.OutputReleasePage != "" {
		releasePage, err := self.createReleasePage(opts, release)
		if err != nil {
			return fmt.Errorf("could not create release page %v: %w", opts, err)
		}
		err = self.writeMessage(opts, releasePage, opts.OutputReleasePage)
		if err != nil {
			return fmt.Errorf("could not output release page %v: %w", opts, err)
		}
	}
	return nil
}

func (self *Generator) writeMessage(opts *GenerateOpts, message proto.Message, outputPath string) error {
	data, err := opts.MarshalOptions.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message: %w", opts, err)
	}
	fileMode, err := strconv.ParseUint(opts.OutputFileMode, 8, 32)
	if err != nil {
		return fmt.Errorf("could not create file mode for %s: %w", opts.OutputFileMode, err)
	}
	if err := os.WriteFile(outputPath, data, os.FileMode(fileMode)); err != nil {
		return fmt.Errorf("could not write data to file %s: %w", outputPath, err)
	}
	return nil
}

func (self *Generator) createReleasePage(opts *GenerateOpts, release *contracts.Release) (*contracts.ReleasePage, error) {
	page := &contracts.ReleasePage{}
	artifactSection := &contracts.ReleasePageSection{Title: "Artifacts"}
	changelogSection := &contracts.ReleasePageSection{Title: "Changelog"}
	page.Sections = append(page.Sections, changelogSection)
	page.Sections = append(page.Sections, artifactSection)
	for _, commit := range release.Git.Commits {
		sectionItem := &contracts.ReleasePageSectionItem{Content: commit}
		changelogSection.Items = append(changelogSection.Items, sectionItem)
		sectionItem.Attrs = append(
			sectionItem.Attrs,
			&contracts.ReleasePageSectionItemAttr{
				Name:    "Hash",
				Content: commit,
			},
		)
	}
	for _, item := range release.Items {
		sectionItem := &contracts.ReleasePageSectionItem{}
		if item.File != nil {
			size := float64(item.File.Size) / 1024 / 1024
			size = roundFloat(size, 2)
			sectionItem.Content = item.File.Name
			sectionItem.Attrs = append(
				sectionItem.Attrs,
				&contracts.ReleasePageSectionItemAttr{
					Name:    "Type",
					Content: "File",
				},
			)
			sectionItem.Attrs = append(
				sectionItem.Attrs,
				&contracts.ReleasePageSectionItemAttr{
					Name:    "Size",
					Content: fmt.Sprintf("%f MB", size),
				},
			)
			for _, hash := range item.File.Hashes {
				sectionItem.Attrs = append(
					sectionItem.Attrs,
					&contracts.ReleasePageSectionItemAttr{
						Name:    hash.Algo,
						Content: hash.Content,
					},
				)
			}
		}
		artifactSection.Items = append(artifactSection.Items, sectionItem)
	}
	return page, nil
}

func (self *Generator) createRelease(opts *GenerateOpts) (*contracts.Release, error) {
	release := &contracts.Release{}
	for _, item := range opts.MergeItems {
		if err := self.addItem(release, item); err != nil {
			return nil, fmt.Errorf("could not add item %s to release: %w", item, err)
		}
	}
	for _, manifest := range opts.MergeManifests {
		if err := self.addManifest(release, manifest); err != nil {
			return nil, fmt.Errorf("could not add manifest %s to release: %w", manifest, err)
		}
	}
	if release.Git != nil {
		err := self.addGit(release, opts.GitRoot)
		if err != nil {
			return nil, fmt.Errorf("could not add git info to release: %w", err)
		}
	}
	return release, nil
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

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
