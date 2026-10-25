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

	"git.alwaldend.com/alwaldend/src/tools/release/api/contracts"
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
	if p.target.Namespace == "third_party" {
		if _, err := p.remote("test -d " + shellQuote(staging) + "; " + shellJoin([]string{"mkdir", "-p", "--", p.publicDir})); err != nil {
			return fmt.Errorf("prepare third-party release directory %q: %w", p.publicDir, err)
		}
		input := strings.NewReader(strings.Join(files, "\x00") + "\x00")
		extra := []string{"--copy-links", "--from0", "--files-from=-", "--temp-dir=" + staging}
		if err := p.rsync(filesDir+string(filepath.Separator), p.remotePath(p.publicDir)+"/", extra, input); err != nil {
			return fmt.Errorf("copy third-party files; earlier completed files may already be public: %w", err)
		}
		return nil
	}
	if _, err := p.remote(p.requireAbsent(p.publicDir) + "; " + p.requireAbsent(p.siteDir)); err != nil {
		return fmt.Errorf("check new release %q: %w", p.publicDir, err)
	}
	remote, err := p.remote(shellJoin([]string{"mktemp", "-d", staging + "/publish.XXXXXXXX"}))
	if err != nil {
		return fmt.Errorf("prepare upload staging: %w", err)
	}
	attempt := strings.TrimSpace(string(remote))
	if path.Dir(attempt) != staging || !strings.HasPrefix(path.Base(attempt), "publish.") || pathComponent(path.Base(attempt)) != nil {
		return fmt.Errorf("invalid upload staging directory")
	}
	defer p.cleanup(attempt)
	input := strings.NewReader(strings.Join(files, "\x00") + "\x00")
	extra := []string{"--copy-links", "--from0", "--files-from=-", "--temp-dir=" + attempt}
	if err := p.rsync(filesDir+string(filepath.Separator), p.remotePath(attempt+"/release")+"/", extra, input); err != nil {
		return fmt.Errorf("stage release files: %w", err)
	}
	script := p.requireAbsent(p.publicDir) + "; " + p.requireAbsent(p.siteDir) + "; "
	// --files-from omits the transfer root's attributes; sudo's umask can
	// otherwise leave the promoted archive directory unreadable by the server.
	script += shellJoin([]string{"chmod", "0755", "--", attempt + "/release"}) + "; "
	script += shellJoin([]string{"mkdir", "-p", "--", path.Dir(p.publicDir)}) + "; "
	script += shellJoin([]string{"mv", "-T", "--", attempt + "/release", p.publicDir})
	if err := p.locked(script); err != nil {
		return fmt.Errorf("publish complete release %q: %w", p.publicDir, err)
	}
	return nil
}

// requireAbsent rejects both existing content and dangling symlinks.
func (p *sshPublisher) requireAbsent(directory string) string {
	quoted := shellQuote(directory)
	return "if test -e " + quoted + " || test -L " + quoted + "; then printf '%s\\n' " + shellQuote("release version already exists: "+directory) + " >&2; exit 1; fi"
}

func (p *sshPublisher) extractSite() error {
	if _, err := p.remote(p.requireAbsent(p.siteDir)); err != nil {
		return fmt.Errorf("inspect site directory %q: %w", p.siteDir, err)
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
	script := p.requireAbsent(p.siteDir) + "; " + shellJoin([]string{"mkdir", "-p", "--", p.siteRoot + "/releases"}) + "; "
	script += shellJoin([]string{"mv", "-T", "--", attempt + "/site", p.siteDir}) + "; " + p.checkSite()
	if err := p.locked(script); err != nil {
		return fmt.Errorf("retain extracted site %q: %w", p.siteDir, err)
	}
	return nil
}

func (p *sshPublisher) linkSite() error {
	current := shellQuote(p.siteRoot + "/current")
	// Check each managed ancestor before allowing recursive site cleanup.
	script := ""
	for _, directory := range []string{p.target.Root, path.Join(p.target.Root, "sites"), p.siteRoot, p.siteRoot + "/releases"} {
		quoted := shellQuote(directory)
		script += "test -d " + quoted + "; test ! -L " + quoted + "; "
		if directory != p.target.Root {
			script += "if mountpoint -q -- " + quoted + "; then printf '%s\\n' 'Mounted site ancestry prevents cleanup' >&2; exit 1; fi; "
		}
	}
	// Check mount points before selection: rm's filesystem boundary is relative
	// to each operand and does not protect a version that is itself mounted.
	script = "command -v mountpoint >/dev/null; " + script
	script += "cleanup_mount=$(" + shellJoin([]string{"find", "-P", p.siteRoot + "/releases", "-xdev", "-mindepth", "1", "-type", "d", "-exec", "mountpoint", "-q", "--", "{}", ";", "-print", "-quit"}) + "); "
	script += "if test -n \"$cleanup_mount\"; then printf '%s\\n' 'Mounted release prevents cleanup' >&2; exit 1; fi; "
	// A repeat selection also repairs a release created with a restrictive
	// transfer umask. Keep archive bytes and absent artifact trees untouched.
	publicDirectory := shellQuote(p.publicDir)
	script += "if test -e " + publicDirectory + " || test -L " + publicDirectory + "; then "
	for _, directory := range []string{path.Join(p.target.Root, "projects"), path.Dir(path.Dir(p.publicDir)), path.Dir(p.publicDir), p.publicDir} {
		quoted := shellQuote(directory)
		script += "test -d " + quoted + "; test ! -L " + quoted + "; "
	}
	script += shellJoin([]string{"chmod", "0755", "--", p.publicDir}) + "; fi; "
	script += p.checkSite() + "; if test -e " + current + "; then test -L " + current + "; fi; "
	script += "selection=$(" + shellJoin([]string{"mktemp", "-d", path.Join(p.target.Root, "staging", "select.XXXXXXXX")}) + "); "
	script += "trap 'rm -rf -- \"$selection\"' EXIT; ln -s -- " + shellQuote("releases/"+p.version) + " \"$selection/current\"; mv -Tf -- \"$selection/current\" " + current
	// Enumerate only immediate version entries, comparing literal paths rather
	// than patterns. rm does not follow links and cannot cross nested mounts.
	cleanup := `set -eu; selected=$1; shift; for candidate do
if test "$candidate" != "$selected"; then
  rm -rf --one-file-system -- "$candidate"
fi
done`
	script += "; if ! " + shellJoin([]string{"find", "-P", p.siteRoot + "/releases", "-mindepth", "1", "-maxdepth", "1", "-exec", "sh", "-c", cleanup, "site-retention", p.siteDir, "{}", "+"}) + "; then "
	script += "printf '%s\\n' 'New site selected, but obsolete site cleanup failed; retry link' >&2; exit 1; fi"
	if err := p.locked(script); err != nil {
		return fmt.Errorf("select site directory and clean obsolete versions %q: %w", p.siteDir, err)
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
	// sudo may retain a login home that the content account cannot traverse.
	args := append(append([]string{}, p.args...), "--", p.target.Host, p.sudo(shellJoin([]string{"sh", "-c", "set -eu; cd /; umask 022; " + script})))
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
