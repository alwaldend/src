package release

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/tools/rules_al/api/contracts"
)

func (self *Deployer) deploySSH(opts *DeployerOpts, release *contracts.Release, releaseDir string, target *contracts.ReleaseDeploymentSsh, items []*contracts.ReleaseItem) error {
	publisher, err := newSSHPublisher(opts, release, target)
	if err != nil {
		return fmt.Errorf("create SSH publisher: %w", err)
	}
	ctx, cancel := context.WithTimeout(opts.Ctx, opts.Timeout)
	defer cancel()
	publisher.ctx = ctx
	publisher.items = items
	operation, err := newArtifactOperation(publisher)
	if err != nil {
		return fmt.Errorf("prepare artifact evidence: %w", err)
	}
	for index, step := range opts.SSHSteps {
		if operation != nil {
			operation.Steps[index].Status = "unknown"
			if err := checkpointReceipt(opts); err != nil {
				return fmt.Errorf("checkpoint before %s: %w", step, err)
			}
		}
		switch step {
		case "delete":
			err = publisher.deleteRelease()
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
		if operation != nil {
			operation.Steps[index].Status = "completed"
			if err := checkpointReceipt(opts); err != nil {
				return fmt.Errorf("checkpoint after %s: %w", step, err)
			}
		}
		fmt.Fprintf(opts.Stdout, "Completed %s for %s/%s in %s\n", step, publisher.target.Project, release.Name, opts.Environment)
	}
	return nil
}

func (p *sshPublisher) uploadRelease(releaseDir string, items []*contracts.ReleaseItem) error {
	manifest, data, err := makePublication(p.release, p.target, items)
	if err != nil {
		return fmt.Errorf("prepare publication: %w", err)
	}
	// Verify a completed identity under the destination lock before requiring
	// local payloads, staging space, or another network transfer.
	script := p.checkAncestry(p.publicDir) + "; if test -e " + shellQuote(p.publicDir) + " || test -L " + shellQuote(p.publicDir) + "; then " + p.checkPublication(p.publicDir, manifest, data) + "; printf present; else printf absent; fi"
	state, err := p.lockedOutput(script)
	if err != nil {
		return fmt.Errorf("inspect immutable release: %w", err)
	}
	if string(state) == "present" {
		return nil
	}
	if string(state) != "absent" {
		return fmt.Errorf("unrecognized destination publication state")
	}
	local, err := publicationTemp()
	if err != nil {
		return fmt.Errorf("prepare local staging: %w", err)
	}
	defer os.RemoveAll(local)
	for _, file := range manifest.Files {
		source := filepath.Join(releaseDir, "files", file.Name)
		if err := verifyPublicationFile(source, file); err != nil {
			return fmt.Errorf("validate release payload: %w", err)
		}
		if err := copyBundleFile(source, filepath.Join(local, file.Name)); err != nil {
			return fmt.Errorf("stage release payload: %w", err)
		}
		if err := verifyPublicationFile(filepath.Join(local, file.Name), file); err != nil {
			return fmt.Errorf("validate staged payload: %w", err)
		}
	}
	if err := os.WriteFile(filepath.Join(local, "release.json"), []byte(data), 0o644); err != nil {
		return fmt.Errorf("stage publication manifest: %w", err)
	}
	attempt, err := p.newStaging()
	if err != nil {
		return fmt.Errorf("prepare remote staging: %w", err)
	}
	defer p.cleanup(attempt)
	staged := attempt + "/release"
	if err := p.rsync(local+string(filepath.Separator), p.remotePath(staged)+"/", nil, nil); err != nil {
		return fmt.Errorf("transfer staged release: %w", err)
	}
	script = p.checkAncestry(p.publicDir) + "; " + p.checkPublication(staged, manifest, data) + "; "
	script += shellJoin([]string{"mkdir", "-p", "--", path.Dir(p.publicDir)}) + "; "
	script += "if test -e " + shellQuote(p.publicDir) + " || test -L " + shellQuote(p.publicDir) + "; then " + p.checkPublication(p.publicDir, manifest, data) + "; else " + shellJoin([]string{"mv", "-T", "--", staged, p.publicDir}) + "; fi"
	if err := p.locked(script); err != nil {
		return fmt.Errorf("publish immutable release: %w", err)
	}
	return nil
}

// newStaging allocates a private directory on the destination content filesystem.
func (p *sshPublisher) newStaging() (string, error) {
	staging := path.Join(p.target.Root, "staging")
	output, err := p.remote(p.checkAncestry(staging) + "; " + shellJoin([]string{"mktemp", "-d", staging + "/publish.XXXXXXXX"}))
	if err != nil {
		return "", fmt.Errorf("allocate destination staging: %w", err)
	}
	attempt := strings.TrimSpace(string(output))
	if path.Dir(attempt) != staging || !strings.HasPrefix(path.Base(attempt), "publish.") || pathComponent(path.Base(attempt)) != nil {
		return "", fmt.Errorf("invalid remote staging directory")
	}
	return attempt, nil
}

// checkPublication generates destination checks for identity, file set, and bytes.
func (p *sshPublisher) checkPublication(directory string, manifest publicationManifest, data string) string {
	script := "test -d " + shellQuote(directory) + "; test ! -L " + shellQuote(directory) + "; test -f " + shellQuote(directory+"/release.json") + "; test ! -L " + shellQuote(directory+"/release.json") + "; printf %s " + shellQuote(data) + " | cmp - " + shellQuote(directory+"/release.json")
	script += "; test \"$(find " + shellQuote(directory) + " -mindepth 1 -maxdepth 1 | wc -l)\" -eq " + strconv.Itoa(len(manifest.Files)+1)
	for _, file := range manifest.Files {
		filename := shellQuote(path.Join(directory, file.Name))
		script += "; test -f " + filename + "; test ! -L " + filename + "; test \"$(wc -c < " + filename + ")\" -eq " + strconv.FormatInt(file.Size, 10)
		script += "; test \"$(sha256sum < " + filename + ")\" = " + shellQuote(file.SHA256+"  -")
	}
	return script
}

// deleteRelease removes an inactive identity without retaining restoration state.
func (p *sshPublisher) deleteRelease() error {
	current := shellQuote(p.siteRoot + "/current")
	script := p.checkRemoval(p.publicDir) + "; "
	if p.target.Namespace == "projects" {
		script += p.checkRemoval(p.siteDir) + "; if test -L " + current + "; then test \"$(readlink -f -- " + current + ")\" != \"$(readlink -f -- " + shellQuote(p.siteDir) + ")\"; fi; "
		script += shellJoin([]string{"rm", "-rf", "--one-file-system", "--", p.siteDir}) + "; "
	}
	script += shellJoin([]string{"rm", "-rf", "--one-file-system", "--", p.publicDir})
	if err := p.locked(script); err != nil {
		return fmt.Errorf("delete inactive release: %w", err)
	}
	return nil
}

func (p *sshPublisher) extractSite() error {
	manifest, data, err := makePublication(p.release, p.target, p.items)
	if err != nil {
		return fmt.Errorf("prepare extraction identity: %w", err)
	}
	archiveFile := publicationFile{}
	for _, file := range manifest.Files {
		if file.Name == p.target.SiteArchive {
			archiveFile = file
		}
	}
	if archiveFile.Name == "" {
		return fmt.Errorf("site archive is not a declared payload")
	}
	if err := p.locked(p.checkPublication(p.publicDir, manifest, data)); err != nil {
		return fmt.Errorf("verify published archive: %w", err)
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
	if err := verifyPublicationFile(archive, archiveFile); err != nil {
		return fmt.Errorf("verify retrieved archive: %w", err)
	}
	extracted := filepath.Join(local, "site")
	if err := os.Mkdir(extracted, 0o755); err != nil {
		return fmt.Errorf("create extraction directory %q: %w", extracted, err)
	}
	if err := extractWebsite(archive, extracted); err != nil {
		return fmt.Errorf("extract published archive %q: %w", publishedArchive, err)
	}
	if err := os.WriteFile(filepath.Join(extracted, ".release.json"), []byte(data), 0o644); err != nil {
		return fmt.Errorf("write site binding: %w", err)
	}
	attempt, err := p.newStaging()
	if err != nil {
		return fmt.Errorf("stage extracted site: %w", err)
	}
	defer p.cleanup(attempt)
	if err := p.rsync(extracted+string(filepath.Separator), p.remotePath(attempt+"/site")+"/", []string{"--temp-dir=" + attempt}, nil); err != nil {
		return fmt.Errorf("copy extracted site to %q: %w", attempt, err)
	}
	script := p.checkAncestry(p.siteDir) + "; " + p.checkAncestry(p.publicDir) + "; " + p.checkPublication(p.publicDir, manifest, data) + "; " + shellJoin([]string{"mkdir", "-p", "--", p.siteRoot + "/releases"}) + "; "
	script += "if test ! -e " + shellQuote(p.siteDir) + "; then " + shellJoin([]string{"mv", "-T", "--", attempt + "/site", p.siteDir}) + "; fi; " + p.checkSite() + "; printf %s " + shellQuote(data) + " | cmp - " + shellQuote(p.siteDir+"/.release.json")
	if err := p.locked(script); err != nil {
		return fmt.Errorf("retain extracted site %q: %w", p.siteDir, err)
	}
	return nil
}

func (p *sshPublisher) linkSite() error {
	manifest, data, err := makePublication(p.release, p.target, p.items)
	if err != nil {
		return fmt.Errorf("prepare activation identity: %w", err)
	}
	current := shellQuote(p.siteRoot + "/current")
	script := p.checkAncestry(p.publicDir) + "; " + p.checkRemoval(p.siteDir) + "; " + p.checkPublication(p.publicDir, manifest, data) + "; " + p.checkSite() + "; printf %s " + shellQuote(data) + " | cmp - " + shellQuote(p.siteDir+"/.release.json") + "; if test -e " + current + "; then test -L " + current + "; fi; "
	script += "selection=$(" + shellJoin([]string{"mktemp", "-d", path.Join(p.target.Root, "staging", "select.XXXXXXXX")}) + "); "
	script += "trap 'rm -rf -- \"$selection\"' EXIT; ln -s -- " + shellQuote("releases/"+p.version) + " \"$selection/current\"; mv -Tf -- \"$selection/current\" " + current
	if err := p.locked(script); err != nil {
		return fmt.Errorf("select site directory %q: %w", p.siteDir, err)
	}
	return nil
}

// checkAncestry rejects symlinks, non-directories and mounts below the content root.
// Missing components are permitted for creation and idempotent deletion.
func (p *sshPublisher) checkAncestry(directory string) string {
	script := "command -v mountpoint >/dev/null; test -d " + shellQuote(p.target.Root) + "; test ! -L " + shellQuote(p.target.Root)
	current := p.target.Root
	for _, component := range strings.Split(strings.TrimPrefix(directory, p.target.Root+"/"), "/") {
		current = path.Join(current, component)
		quoted := shellQuote(current)
		script += "; if test -e " + quoted + " || test -L " + quoted + "; then test -d " + quoted + "; test ! -L " + quoted + "; "
		script += "if mountpoint -q -- " + quoted + "; then printf '%s\\n' 'Mounted release ancestry prevents operation' >&2; exit 1; fi; fi"
	}
	return script
}

// checkRemoval checks the operand and all descendant mounts before recursive work.
func (p *sshPublisher) checkRemoval(directory string) string {
	script := p.checkAncestry(directory) + "; if test -d " + shellQuote(directory) + "; then "
	script += "release_mount=$(" + shellJoin([]string{"find", "-P", directory, "-xdev", "-mindepth", "1", "-exec", "mountpoint", "-q", "--", "{}", ";", "-print", "-quit"}) + "); "
	script += "if test -n \"$release_mount\"; then printf '%s\\n' 'Mounted release prevents operation' >&2; exit 1; fi; fi"
	return script
}

func (p *sshPublisher) checkSite() string {
	return "test -d " + shellQuote(p.siteDir) + "; test ! -L " + shellQuote(p.siteDir) + "; test -f " + shellQuote(p.siteDir+"/index.html")
}

// locked executes a serialized destination operation without returning output.
func (p *sshPublisher) locked(script string) error {
	if _, err := p.lockedOutput(script); err != nil {
		return fmt.Errorf("run locked publication operation: %w", err)
	}
	return nil
}

// lockedOutput returns output from an operation serialized with every project mutation.
func (p *sshPublisher) lockedOutput(script string) ([]byte, error) {
	lock := path.Join(p.target.Root, "staging", p.target.Project+".lock")
	output, err := p.remote(p.checkAncestry(path.Join(p.target.Root, "staging")) + "; " + shellJoin([]string{"flock", "--", lock, "sh", "-c", "set -eu; umask 022; " + script}))
	if err != nil {
		return nil, fmt.Errorf("update release under lock %q: %w", lock, err)
	}
	return output, nil
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
	// The content account may not be able to traverse the administrator's home.
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
