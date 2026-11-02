package release

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// artifactReport describes resolved metadata or observations without authorizing writes.
type artifactReport struct {
	// Schema identifies the plan or inspection JSON contract.
	Schema string `json:"schema"`
	// Mode distinguishes offline metadata from a read-only destination snapshot.
	Mode string `json:"mode"`
	// ObservedAt records when inspection finished; plans make no observation claim.
	ObservedAt string `json:"observed_at,omitempty"`
	// Operations preserves input bundle order and stable destination selection.
	Operations []artifactObservation `json:"operations"`
}

// artifactObservation binds the publication preview and optional observed state.
type artifactObservation struct {
	// Namespace selects first-party or external artifact storage.
	Namespace string `json:"namespace"`
	// Project identifies the independent publication owner.
	Project string `json:"project"`
	// Version is the selected immutable publication identity.
	Version string `json:"version"`
	// Commit identifies the first-party source revision.
	Commit string `json:"commit,omitempty"`
	// Channel identifies the first-party source version channel.
	Channel string `json:"channel,omitempty"`
	// Selection is explicit for a host/deployment override, otherwise environment.
	Selection string `json:"selection"`
	// PublicationDigest hashes the canonical manifest including its final newline.
	PublicationDigest string `json:"publication_digest"`
	// Destination reports operational addressing without authentication settings.
	Destination observationDestination `json:"destination"`
	// PublicPath is the resolved immutable payload directory.
	PublicPath string `json:"public_path"`
	// SitePath is the first-party retained website directory, when applicable.
	SitePath string `json:"site_path,omitempty"`
	// SiteArchive names the selected archive without claiming its contents are safe.
	SiteArchive string `json:"site_archive,omitempty"`
	// Files contains the exact selected payload metadata in filename order.
	Files []publicationFile `json:"files"`
	// Steps contains the requested mutation steps; observation executes none of them.
	Steps []string `json:"steps"`
	// Local independently verifies selected local payloads during inspection only.
	Local []localObservation `json:"local,omitempty"`
	// Remote reports the destination snapshot during inspection only.
	Remote *remoteObservation `json:"remote,omitempty"`
}

// observationDestination does not infer addressing settings owned by OpenSSH.
type observationDestination struct {
	// Environment names the explicitly selected deployment environment.
	Environment string `json:"environment"`
	// Host is the configured host or SSH alias.
	Host string `json:"host"`
	// Port is null when OpenSSH configuration owns its resolution.
	Port *int `json:"port"`
	// SSHUser is null when OpenSSH configuration selects the login account.
	SSHUser *string `json:"ssh_user"`
	// PublishUser is null when commands run directly as the SSH login account.
	PublishUser *string `json:"publish_user"`
	// Root is the resolved remote content root.
	Root string `json:"root"`
}

// localObservation records payload availability separately from published bytes.
type localObservation struct {
	// Name identifies the selected payload basename.
	Name string `json:"name"`
	// Status is verified, missing, corrupt, or unreadable.
	Status string `json:"status"`
}

// remoteObservation reports state relative to the selected canonical publication.
type remoteObservation struct {
	// Status is absent, matching, conflicting, corrupt, or unavailable.
	Status string `json:"status"`
	// Site independently reports retained binding and current selection.
	Site siteObservation `json:"site"`
}

// siteObservation distinguishes website validity from the observed current link.
type siteObservation struct {
	// Status is absent, retained, active, conflicting, corrupt, unavailable, or not_applicable.
	Status string `json:"status"`
	// Active is independent of validity and null when selection cannot be determined.
	Active *bool `json:"active"`
	// Current is absent, selected, other, dangling, unsafe, unavailable, or not_applicable.
	Current string `json:"current"`
}

// reportArtifacts resolves metadata before optionally observing local and remote bytes.
// A completed observation may report corruption and still succeeds; unavailable
// transport or protocol evidence produces JSON followed by a command error.
func reportArtifacts(opts *DeployerOpts, inspect bool) error {
	if len(opts.Releases) == 0 {
		return fmt.Errorf("at least one release_dir is required")
	}
	report := artifactReport{Schema: "release/artifact_plan/v1", Mode: "offline_metadata", Operations: []artifactObservation{}}
	if inspect {
		report.Schema, report.Mode = "release/artifact_inspection/v1", "read_only_snapshot"
	}
	var observationErr error
	for _, directory := range opts.Releases {
		if !filepath.IsAbs(directory) && os.Getenv("BUILD_WORKSPACE_DIRECTORY") != "" {
			directory = filepath.Join(os.Getenv("BUILD_WORKSPACE_DIRECTORY"), directory)
		}
		publishers, err := selectRelease(opts, directory)
		if err != nil {
			return fmt.Errorf("select release %q: %w", directory, err)
		}
		for _, publisher := range publishers {
			manifest, data, err := makePublication(publisher.release, publisher.target, publisher.items)
			if err != nil {
				return fmt.Errorf("prepare publication metadata: %w", err)
			}
			for _, step := range opts.SSHSteps {
				if step == "extract" {
					if _, err := publicationArchive(manifest, publisher.target.SiteArchive); err != nil {
						return fmt.Errorf("select planned extraction archive: %w", err)
					}
				}
			}
			operation := artifactObservation{
				Namespace: manifest.Namespace, Project: manifest.Project, Version: manifest.Version, Commit: manifest.Commit, Channel: manifest.Channel,
				Selection: "environment", PublicationDigest: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(data))),
				Destination: observationDestination{Environment: opts.Environment, Host: publisher.target.Host, Root: publisher.target.Root},
				PublicPath:  publisher.publicDir, SiteArchive: publisher.target.SiteArchive, Files: manifest.Files, Steps: opts.SSHSteps,
			}
			if opts.SSHHost != "" || opts.SSHDeployment != "" {
				operation.Selection = "explicit"
			}
			if opts.SSHPort != 0 {
				operation.Destination.Port = &opts.SSHPort
			}
			if publisher.target.User != "" {
				operation.Destination.SSHUser = &publisher.target.User
			}
			if publisher.target.PublishUser != "" {
				operation.Destination.PublishUser = &publisher.target.PublishUser
			}
			if manifest.Namespace == "projects" {
				operation.SitePath = publisher.siteDir
			}
			if inspect {
				operation.Local, err = observeLocal(opts.Ctx, directory, manifest.Files)
				if err != nil {
					return fmt.Errorf("verify local inspection payloads: %w", err)
				}
				operation.Remote, err = observeRemote(publisher, manifest, data)
				if err != nil {
					observationErr = errors.Join(observationErr, fmt.Errorf("inspect %s/%s at %q: %w", manifest.Project, manifest.Version, publisher.target.Host, err))
				}
			}
			report.Operations = append(report.Operations, operation)
		}
	}
	if inspect {
		report.ObservedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	encoder := json.NewEncoder(opts.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return errors.Join(observationErr, fmt.Errorf("write artifact report: %w", err))
	}
	if observationErr != nil {
		return fmt.Errorf("incomplete destination observation: %w", observationErr)
	}
	return nil
}

// observeLocal preserves missing or invalid local bytes as evidence, never as a
// prerequisite for inspecting a successfully published destination.
func observeLocal(ctx context.Context, directory string, files []publicationFile) ([]localObservation, error) {
	observations := make([]localObservation, 0, len(files))
	for _, file := range files {
		state := "verified"
		if err := verifyPublicationFileContext(ctx, filepath.Join(directory, "files", file.Name), file); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, fmt.Errorf("verify payload %q: %w", file.Name, err)
			}
			var pathErr *os.PathError
			switch {
			case errors.Is(err, os.ErrNotExist):
				state = "missing"
			case errors.As(err, &pathErr):
				state = "unreadable"
			default:
				state = "corrupt"
			}
		}
		observations = append(observations, localObservation{Name: file.Name, Status: state})
	}
	return observations, nil
}

// boundedObservation cancels the process as soon as either captured stream exceeds
// its bound, so a noisy endpoint cannot consume memory or the full SSH deadline.
type boundedObservation struct {
	// data holds at most 16 KiB of non-authoritative protocol or diagnostic output.
	data bytes.Buffer
	// cancel stops the process immediately on overflow.
	cancel context.CancelFunc
}

// Write bounds capture without disclosing remote diagnostics or SSH arguments.
func (b *boundedObservation) Write(data []byte) (int, error) {
	if len(data) > 16*1024-b.data.Len() {
		b.cancel()
		return 0, fmt.Errorf("remote inspection output exceeds 16 KiB")
	}
	count, err := b.data.Write(data)
	if err != nil {
		return count, fmt.Errorf("capture remote inspection: %w", err)
	}
	return count, nil
}

// observeRemote executes only read operations with bounded time and captured output.
func observeRemote(publisher *sshPublisher, manifest publicationManifest, data string) (*remoteObservation, error) {
	unavailable := &remoteObservation{Status: "unavailable", Site: siteObservation{Status: "unavailable", Current: "unavailable"}}
	ctx, cancel := context.WithTimeout(publisher.opts.Ctx, publisher.opts.Timeout)
	defer cancel()
	reader := *publisher
	reader.ctx = ctx
	stdout, stderr := &boundedObservation{cancel: cancel}, &boundedObservation{cancel: cancel}
	if err := runRemoteCommand(&reader, observationScript(&reader, manifest, data), stdout, stderr); err != nil {
		return unavailable, fmt.Errorf("read destination: %w", err)
	}
	fields := strings.Fields(stdout.data.String())
	if len(fields) != 5 || fields[0] != "release-inspection-v1" || !oneOf(fields[1], "absent", "matching", "conflicting", "corrupt") || !oneOf(fields[2], "absent", "retained", "active", "conflicting", "corrupt", "not_applicable") || !oneOf(fields[3], "true", "false", "unknown") || !oneOf(fields[4], "absent", "selected", "other", "dangling", "unsafe", "not_applicable") {
		return unavailable, fmt.Errorf("invalid destination inspection response")
	}
	result := &remoteObservation{Status: fields[1], Site: siteObservation{Status: fields[2], Current: fields[4]}}
	if fields[3] != "unknown" {
		active := fields[3] == "true"
		result.Site.Active = &active
	}
	return result, nil
}

// oneOf validates the small, closed destination observation protocol.
func oneOf(value string, choices ...string) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

// observationCheck runs legacy semicolon-separated assertions in an independent
// shell; conditional context must not disable their required errexit behavior.
func observationCheck(script string) string {
	return shellJoin([]string{"sh", "-ec", script}) + " >/dev/null 2>&1"
}

// readableAncestry requires traversal from the filesystem root and read access
// within the content root. A missing content root is unavailable, not an absent
// release. Publication guards separately classify unsafe symlinks and mounts.
func readableAncestry(root, directory string) string {
	script := ""
	current := "/"
	components := append([]string{""}, strings.Split(strings.TrimPrefix(directory, "/"), "/")...)
	for _, component := range components {
		current = path.Join(current, component)
		quoted := shellQuote(current)
		script += "if test -d " + quoted + "; then test -x " + quoted + " || exit 1; "
		if current == root || strings.HasPrefix(current, root+"/") {
			script += "if test ! -L " + quoted + "; then test -r " + quoted + " || exit 1; fi; "
		}
		script += "fi; "
	}
	return script + "test -e " + shellQuote(root) + " || test -L " + shellQuote(root) + " || exit 1; "
}

// observationScript reads canonical publication bytes, retained binding and current
// selection without taking a lock, creating staging, transferring or repairing data.
func observationScript(p *sshPublisher, manifest publicationManifest, data string) string {
	public, canonical := shellQuote(p.publicDir), shellQuote(p.publicDir+"/release.json")
	script := "command -v mountpoint >/dev/null; command -v sha256sum >/dev/null; command -v cmp >/dev/null; command -v find >/dev/null; command -v wc >/dev/null; command -v readlink >/dev/null; "
	script += readableAncestry(p.target.Root, p.publicDir)
	script += "publication=absent; site=not_applicable; active=unknown; current=not_applicable; "
	script += "if ! " + observationCheck(p.checkAncestry(p.publicDir)) + "; then publication=corrupt; elif test -e " + public + " || test -L " + public + "; then "
	script += "if test ! -f " + canonical + " || test -L " + canonical + "; then publication=corrupt; else test -r " + canonical + " || exit 1; "
	script += "if printf %s " + shellQuote(data) + " | cmp -s - " + canonical + "; then "
	for _, file := range manifest.Files {
		filename := shellQuote(path.Join(p.publicDir, file.Name))
		script += "if test -f " + filename + " && test ! -L " + filename + "; then test -r " + filename + " || exit 1; fi; "
	}
	script += "if " + observationCheck(p.checkPublication(p.publicDir, manifest, data)) + "; then publication=matching; else result=$?; test \"$result\" -eq 1 || exit 1; publication=corrupt; fi; else result=$?; test \"$result\" -eq 1 || exit 1; publication=conflicting; fi; fi; fi; "
	if p.target.Namespace == "projects" {
		website, current := shellQuote(p.siteDir), shellQuote(p.siteRoot+"/current")
		binding := shellQuote(p.siteDir + "/.release.json")
		script += readableAncestry(p.target.Root, p.siteDir)
		script += "site=absent; active=false; current=absent; "
		script += "if ! " + observationCheck(p.checkAncestry(p.siteRoot)) + "; then site=corrupt; active=unknown; current=unsafe; else "
		script += "if test -L " + current + "; then selected=$(readlink -n -- " + current + " && printf .); selected=${selected%.}; case \"$selected\" in *[[:cntrl:]]*) current=unsafe;; " + shellQuote("releases/"+p.version) + ") active=true; current=selected;; releases/*) selected_name=${selected#releases/}; case \"$selected_name\" in ''|.|..|*/*) current=unsafe;; *) current=other;; esac;; *) current=unsafe;; esac; "
		script += "if test \"$current\" != unsafe && test ! -e " + current + "; then current=dangling; fi; elif test -e " + current + "; then current=unsafe; fi; if test \"$current\" = unsafe; then active=unknown; fi; "
		script += "if ! " + observationCheck(p.checkAncestry(p.siteDir)) + "; then site=corrupt; elif test -e " + website + " || test -L " + website + "; then "
		checks := p.checkSite() + "; test ! -L " + shellQuote(p.siteDir+"/index.html") + "; test -f " + binding + "; test ! -L " + binding
		script += "if ! " + observationCheck(checks) + "; then site=corrupt; else test -r " + binding + " || exit 1; if printf %s " + shellQuote(data) + " | cmp -s - " + binding + "; then if test \"$active\" = true; then site=active; else site=retained; fi; else result=$?; test \"$result\" -eq 1 || exit 1; site=conflicting; fi; fi; fi; fi; "
	}
	script += "printf 'release-inspection-v1 %s %s %s %s\\n' \"$publication\" \"$site\" \"$active\" \"$current\""
	return script
}
