package main

import (
	"context"
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha256"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	gitLib "git.alwaldend.com/src/tools/git/main/go"
	"git.alwaldend.com/src/tools/release/main/proto/contracts"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Generator struct {
	gitInfo *gitLib.GitInfo
}

type GeneratorParams struct{}

func NewGenerator(gitInfo *gitLib.GitInfo) *Generator {
	return &Generator{gitInfo: gitInfo}
}

type GenerateOpts struct {
	Ctx                context.Context
	AddFiles           []string
	DeploymentInfo     []string
	MergeManifests     []string
	OutputManifests    []string
	OutputReleasePages []string
	OutputFileMode     string
	GitRoot            string
	MarshalOptions     *protojson.MarshalOptions
}

// Generate a manifest file
func (self *Generator) Generate(opts *GenerateOpts) error {
	release, err := self.parseRelease(opts)
	if err != nil {
		return fmt.Errorf("could not create release %v: %w", opts, err)
	}
	for _, manifest := range opts.OutputManifests {
		err := self.writeMessage(opts, release, manifest)
		if err != nil {
			return fmt.Errorf("could not output release %v: %w", manifest, err)
		}
	}
	for _, page := range opts.OutputReleasePages {
		releasePage, err := self.createReleasePage(release)
		if err != nil {
			return fmt.Errorf("could not create release page %v: %w", opts, err)
		}
		err = self.writeMessage(opts, releasePage, page)
		if err != nil {
			return fmt.Errorf("could not write release page to %s %v: %w", page, opts, err)
		}
	}
	return nil
}

func (self *Generator) writeMessage(opts *GenerateOpts, message proto.Message, outputPath string) error {
	data, err := opts.MarshalOptions.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not marshal message %v: %w", opts, err)
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

// Create a release page config
func (self *Generator) createReleasePage(release *contracts.Release) (*contracts.ReleasePage, error) {
	page := &contracts.ReleasePage{}
	changelogSection, err := self.changelogSection(release)
	if err != nil {
		return nil, fmt.Errorf("could not generate changelog section: %w", err)
	}
	artifactSection, err := self.artifactSection(release)
	if err != nil {
		return nil, fmt.Errorf("could not generate artifact section: %w", err)
	}
	page.Sections = append(page.Sections, changelogSection)
	page.Sections = append(page.Sections, artifactSection)
	return page, nil
}

func (self *Generator) artifactSection(release *contracts.Release) (*contracts.ReleasePageSection, error) {
	artifactSection := &contracts.ReleasePageSection{Title: "Artifacts"}
	for _, item := range release.Items {
		sectionItem, err := self.artifactSectionItem(release, item)
		if err != nil {
			return nil, fmt.Errorf("could not generate section for item %v: %w", item, err)
		}
		artifactSection.Items = append(artifactSection.Items, sectionItem)
	}
	return artifactSection, nil
}

func (self *Generator) artifactSectionItem(release *contracts.Release, item *contracts.ReleaseItem) (*contracts.ReleasePageSectionItem, error) {
	if item.File == nil {
		return nil, fmt.Errorf("missing file data")
	}
	sectionItem := &contracts.ReleasePageSectionItem{}
	size := float64(item.File.Size) / 1024 / 1024
	size = roundFloat(size, 2)
	sectionItem.Attrs = append(
		sectionItem.Attrs,
		&contracts.ReleasePageSectionItemAttr{
			Name:    "Type",
			Content: "File",
		},
	)
	for _, deployment := range item.Deployments {
		if deployment.Oci != nil {
			for _, tag := range deployment.Oci.Tags {
				sectionItem.Attrs = append(
					sectionItem.Attrs,
					&contracts.ReleasePageSectionItemAttr{
						Name: "Pull",
						Content: strings.Join(
							[]string{"oras pull", "--output \"${PWD}\"", tag},
							" ",
						),
					},
				)
				if sectionItem.Content == "" && strings.HasPrefix(deployment.Oci.Repository, "docker.io") {
					path := strings.SplitN(deployment.Oci.Repository, "/", 2)
					tagShort := strings.SplitN(tag, ":", 2)
					sectionItem.Content = fmt.Sprintf(
						"[%s](https://hub.docker.com/repository/docker/%s/tags/%s/)",
						item.File.Name,
						path[1],
						tagShort[len(tagShort)-1],
					)
				}
			}
		}
	}
	if sectionItem.Content == "" {
		sectionItem.Content = item.File.Name
	}
	sectionItem.Attrs = append(
		sectionItem.Attrs,
		&contracts.ReleasePageSectionItemAttr{
			Name:    "Size",
			Content: fmt.Sprintf("%.2f MB", size),
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
	return sectionItem, nil
}

func (self *Generator) changelogSection(release *contracts.Release) (*contracts.ReleasePageSection, error) {
	changelogSection := &contracts.ReleasePageSection{Title: "Changelog"}
	for _, commit := range release.Git.Commits {
		sectionItem := &contracts.ReleasePageSectionItem{
			Content: commit.Message,
			Attrs: []*contracts.ReleasePageSectionItemAttr{
				{
					Name:    "Hash",
					Content: commit.Hash,
				},
				{
					Name:    "Author",
					Content: fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email),
				},
				{
					Name:    "Committer",
					Content: fmt.Sprintf("%s <%s>", commit.Committer.Name, commit.Committer.Email),
				},
			},
		}
		if len(commit.Tags) != 0 {
			sectionItem.Attrs = append(sectionItem.Attrs, &contracts.ReleasePageSectionItemAttr{
				Name:    "Tags",
				Content: strings.Join(commit.Tags, ", "),
			})
		}
		if commit.PgpSignature != "" {
			sectionItem.Attrs = append(sectionItem.Attrs, &contracts.ReleasePageSectionItemAttr{
				Name:    "PGP",
				Content: commit.PgpSignature,
			})
		}
		changelogSection.Items = append(changelogSection.Items, sectionItem)
	}
	return changelogSection, nil
}

// Create a release config
func (self *Generator) parseRelease(opts *GenerateOpts) (*contracts.Release, error) {
	release := &contracts.Release{Project: &contracts.Project{}}
	deployments, err := self.parseDeployments(opts.DeploymentInfo)
	if err != nil {
		return nil, fmt.Errorf("could not parse deployments: %w", err)
	}
	for _, manifest := range opts.MergeManifests {
		if err := self.addManifest(release, manifest); err != nil {
			return nil, fmt.Errorf("could not add manifest %s to release: %w", manifest, err)
		}
	}
	if release.Project.Subdir != "" {
		release.Project.SafeSubdir = safeString(release.Project.Subdir)
	}
	for _, file := range opts.AddFiles {
		if err := self.addFile(release, deployments, file); err != nil {
			return nil, fmt.Errorf("could not add file %s to release: %w", file, err)
		}
	}
	encItems := map[string]*contracts.ReleaseItem{}
	for _, item := range release.Items {
		if item.File == nil {
			continue
		}
		if encItem, ok := encItems[item.File.Name]; ok {
			return nil, fmt.Errorf("filename collision: %v and %v", item, encItem)
		} else {
			encItems[item.File.Name] = item
		}
		for _, deployment := range item.Deployments {
			if deployment.Oci != nil {
				tag := fmt.Sprintf(
					"%s:%s_%s_%s",
					deployment.Oci.Repository,
					release.Project.SafeSubdir,
					item.File.SafeName,
					release.Name,
				)
				deployment.Oci.Tags = []string{tag}
			}
		}
	}
	if release.Git != nil {
		err := self.addGit(opts.Ctx, release, opts.GitRoot)
		if err != nil {
			return nil, fmt.Errorf("could not add git info to release: %w", err)
		}
	}
	return release, nil
}

func (self *Generator) parseDeployments(deployments []string) ([]*contracts.ReleaseDeployment, error) {
	res := []*contracts.ReleaseDeployment{}
	for _, path := range deployments {
		deployment := &contracts.ReleaseDeployment{}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not open deployment file %s: %w", path, err)
		}
		if err := protojson.Unmarshal(content, deployment); err != nil {
			return nil, fmt.Errorf("could not parse deployment file %s: %w", path, err)
		}
		res = append(res, deployment)
	}
	return res, nil
}

// Add git information in-place
func (self *Generator) addGit(ctx context.Context, release *contracts.Release, gitRoot string) error {
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
	rev, err := repo.ResolveRevision(plumbing.Revision(release.Git.Revision.Hash))
	if err != nil {
		return fmt.Errorf("could not resolve revision %s: %w", release.Git.Revision, err)
	}
	mutex := &sync.Mutex{}
	_, commitTags, err := self.gitInfo.ParseTags(ctx, repo)
	if err != nil {
		return fmt.Errorf("could not parse tags: %w", err)
	}
	headTags, _ := commitTags[rev.String()]
	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("could not create head ref: %w", err)
	}
	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return fmt.Errorf("could not create head commit: %w", err)
	}
	release.Git.Revision, err = self.gitInfo.ParseCommit(ctx, mutex, headCommit, headTags, false)
	filter := func(path string) bool {
		res := strings.HasPrefix(path, release.Project.Subdir)
		return res
	}
	commitIter, err := repo.Log(&git.LogOptions{From: *rev, PathFilter: filter})
	if err != nil {
		return fmt.Errorf("could not create commit iterator: %w", err)
	}
	err = commitIter.ForEach(func(commitRef *object.Commit) error {
		tags, hasTags := commitTags[commitRef.Hash.String()]
		if commitRef.Hash != *rev && hasTags {
			commitIter.Close()
			return nil
		}
		commit, err := self.gitInfo.ParseCommit(ctx, mutex, commitRef, tags, false)
		if err != nil {
			return fmt.Errorf("could not parse commit %s: %w", commitRef.Hash, err)
		}
		release.Git.Commits = append(release.Git.Commits, commit)
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

// Add file information in-place
func (self *Generator) addFile(
	release *contracts.Release,
	deployments []*contracts.ReleaseDeployment,
	path string,
) error {
	item := &contracts.ReleaseItem{
		File:        &contracts.ReleaseFile{},
		Deployments: deployments,
	}
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("could not stat file %s: %w", path, err)
	}

	item.File.Name = filepath.Base(path)
	item.File.SafeName = safeString(item.File.Name)
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

func safeString(val string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(val, "_")
}

// Round a float to a precision
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
