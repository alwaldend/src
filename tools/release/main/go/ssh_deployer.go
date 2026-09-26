package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/release/main/proto/contracts"
)

func (self *Deployer) deploySSH(opts *DeployerOpts, release *contracts.Release, releaseDir string, target *contracts.ReleaseDeploymentSsh, items []*contracts.ReleaseItem) error {
	publisher, err := newSSHPublisher(opts, release, target)
	if err != nil {
		return fmt.Errorf("create SSH publisher: %w", err)
	}
	ctx, cancel := context.WithTimeout(opts.Ctx, opts.Timeout)
	defer cancel()
	publisher.ctx = ctx
	for _, step := range opts.SSHSteps {
		switch step {
		case "upload":
			err = publisher.uploadRelease(releaseDir, items)
		case "extract":
			err = publisher.extractSite()
		case "link":
			err = publisher.linkSite()
		}
		if err != nil {
			return fmt.Errorf("%s %s/%s in %s: %w", step, publisher.target.Project, release.Name, opts.Environment, err)
		}
		fmt.Fprintf(opts.Stdout, "Completed %s for %s/%s in %s\n", step, publisher.target.Project, release.Name, opts.Environment)
	}
	return nil
}

func (p *sshPublisher) uploadRelease(releaseDir string, items []*contracts.ReleaseItem) error {
	if len(items) == 0 {
		return fmt.Errorf("release contains no files")
	}
	files := make([]string, 0, len(items))
	seen := map[string]bool{}
	filesDir := filepath.Join(releaseDir, "files")
	for _, item := range items {
		if item.File == nil {
			return fmt.Errorf("release item has no file")
		}
		name := item.File.Name
		if err := pathComponent(name); err != nil {
			return fmt.Errorf("release filename: %w", err)
		}
		if seen[name] {
			return fmt.Errorf("duplicate release filename %q", name)
		}
		seen[name] = true
		info, err := os.Stat(filepath.Join(filesDir, name))
		if err != nil {
			return fmt.Errorf("inspect release file %q: %w", name, err)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("release payload is not a regular file: %s", name)
		}
		files = append(files, name)
	}
	staging := path.Join(p.target.Root, "staging")
	if _, err := p.remote("test -d " + shellQuote(staging) + "; " + shellJoin([]string{"mkdir", "-p", "--", p.publicDir})); err != nil {
		return fmt.Errorf("prepare release directory %q: %w", p.publicDir, err)
	}
	input := strings.NewReader(strings.Join(files, "\x00") + "\x00")
	extra := []string{"--copy-links", "--from0", "--files-from=-", "--temp-dir=" + staging}
	if err := p.rsync(filesDir+string(filepath.Separator), p.remotePath(p.publicDir)+"/", extra, input); err != nil {
		return fmt.Errorf("copy release files; earlier completed files may already be public: %w", err)
	}
	return nil
}

func (p *sshPublisher) extractSite() error {
	retained, err := p.remote("if test -d " + shellQuote(p.siteDir) + " && test ! -L " + shellQuote(p.siteDir) + "; then test -f " + shellQuote(p.siteDir+"/index.html") + "; printf retained; fi")
	if err != nil {
		return fmt.Errorf("inspect site directory %q: %w", p.siteDir, err)
	}
	if string(retained) == "retained" {
		return nil
	}
	local, err := publicationTemp()
	if err != nil {
		return fmt.Errorf("prepare extraction workspace: %w", err)
	}
	defer os.RemoveAll(local)
	archive := filepath.Join(local, p.target.SiteArchive)
	publishedArchive := path.Join(p.publicDir, p.target.SiteArchive)
	if err := p.rsync(p.remotePath(publishedArchive), archive, nil, nil); err != nil {
		return fmt.Errorf("read published archive %q: %w", publishedArchive, err)
	}
	extracted := filepath.Join(local, "site")
	if err := os.Mkdir(extracted, 0o755); err != nil {
		return fmt.Errorf("create extraction directory %q: %w", extracted, err)
	}
	if err := extractWebsite(archive, extracted); err != nil {
		return fmt.Errorf("extract published archive %q: %w", publishedArchive, err)
	}
	staging := path.Join(p.target.Root, "staging")
	remote, err := p.remote(shellJoin([]string{"mktemp", "-d", staging + "/publish.XXXXXXXX"}))
	if err != nil {
		return fmt.Errorf("prepare remote extraction staging: %w", err)
	}
	attempt := strings.TrimSpace(string(remote))
	if path.Dir(attempt) != staging || !strings.HasPrefix(path.Base(attempt), "publish.") || pathComponent(path.Base(attempt)) != nil {
		return fmt.Errorf("invalid remote staging directory")
	}
	defer p.cleanup(attempt)
	if err := p.rsync(extracted+string(filepath.Separator), p.remotePath(attempt+"/site")+"/", []string{"--temp-dir=" + attempt}, nil); err != nil {
		return fmt.Errorf("copy extracted site to %q: %w", attempt, err)
	}
	script := shellJoin([]string{"mkdir", "-p", "--", p.siteRoot + "/releases"}) + "; "
	script += "if test ! -e " + shellQuote(p.siteDir) + "; then " + shellJoin([]string{"mv", "-T", "--", attempt + "/site", p.siteDir}) + "; fi; " + p.checkSite()
	if err := p.locked(script); err != nil {
		return fmt.Errorf("retain extracted site %q: %w", p.siteDir, err)
	}
	return nil
}

func (p *sshPublisher) linkSite() error {
	current := shellQuote(p.siteRoot + "/current")
	script := p.checkSite() + "; if test -e " + current + "; then test -L " + current + "; fi; "
	script += "selection=$(" + shellJoin([]string{"mktemp", "-d", path.Join(p.target.Root, "staging", "select.XXXXXXXX")}) + "); "
	script += "trap 'rm -rf -- \"$selection\"' EXIT; ln -s -- " + shellQuote("releases/"+p.version) + " \"$selection/current\"; mv -Tf -- \"$selection/current\" " + current
	if err := p.locked(script); err != nil {
		return fmt.Errorf("select site directory %q: %w", p.siteDir, err)
	}
	return nil
}

func (p *sshPublisher) checkSite() string {
	return "test -d " + shellQuote(p.siteDir) + "; test ! -L " + shellQuote(p.siteDir) + "; test -f " + shellQuote(p.siteDir+"/index.html")
}

func (p *sshPublisher) locked(script string) error {
	lock := path.Join(p.target.Root, "staging", p.target.Project+".lock")
	if _, err := p.remote(shellJoin([]string{"flock", "--", lock, "sh", "-c", "set -eu; umask 022; " + script})); err != nil {
		return fmt.Errorf("update site under lock %q: %w", lock, err)
	}
	return nil
}

func (p *sshPublisher) cleanup(directory string) {
	// Cleanup remains bounded even after a transfer's context is cancelled.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cleanup := *p
	cleanup.ctx = ctx
	if _, err := cleanup.remote(shellJoin([]string{"rm", "-rf", "--", directory})); err != nil {
		fmt.Fprintf(p.opts.Stderr, "Could not clean staging %s: %v\n", directory, err)
	}
}

func publicationTemp() (string, error) {
	base := ""
	if workspace := os.Getenv("BUILD_WORKSPACE_DIRECTORY"); workspace != "" {
		base = filepath.Join(workspace, "out", "release")
		if err := os.MkdirAll(base, 0o700); err != nil {
			return "", fmt.Errorf("create publication workspace %q: %w", base, err)
		}
	}
	directory, err := os.MkdirTemp(base, "publish-")
	if err != nil {
		return "", fmt.Errorf("create extraction temporary directory: %w", err)
	}
	return directory, nil
}

func (p *sshPublisher) sudo(command string) string {
	if p.target.PublishUser != "" {
		return shellJoin([]string{"sudo", "-n", "-u", p.target.PublishUser, "--"}) + " " + command
	}
	return command
}

func (p *sshPublisher) remote(script string) ([]byte, error) {
	args := append(append([]string{}, p.args...), "--", p.target.Host, p.sudo(shellJoin([]string{"sh", "-c", "set -eu; umask 022; " + script})))
	cmd := exec.CommandContext(p.ctx, p.opts.SSHPath, args...)
	var stdout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, p.opts.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("remote publication command: %w", err)
	}
	return stdout.Bytes(), nil
}

func (p *sshPublisher) remotePath(filename string) string {
	host := p.target.Host
	if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return host + ":" + filename
}

func (p *sshPublisher) rsync(source, destination string, extra []string, input io.Reader) error {
	rsh := shellJoin(append([]string{p.opts.SSHPath}, p.args...))
	args := []string{"--recursive", "--times", "--perms", "--checksum", "--chmod=D755,F644", "--protect-args", "--itemize-changes", "--rsh=" + rsh, "--rsync-path=" + p.sudo("rsync")}
	args = append(args, extra...)
	args = append(args, "--", source, destination)
	cmd := exec.CommandContext(p.ctx, p.opts.RsyncPath, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = input, p.opts.Stdout, p.opts.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsync transfer: %w", err)
	}
	return nil
}

// shellQuote preserves one literal POSIX shell argument, including apostrophes.
func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func shellJoin(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = shellQuote(arg)
	}
	return strings.Join(quoted, " ")
}
