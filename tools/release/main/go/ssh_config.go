package main

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode"

	"git.alwaldend.com/alwaldend/src/tools/release/main/proto/contracts"
	"google.golang.org/protobuf/proto"
)

type sshPublisher struct {
	ctx       context.Context
	opts      *DeployerOpts
	target    *contracts.ReleaseDeploymentSsh
	args      []string
	version   string
	publicDir string
	siteRoot  string
	siteDir   string
}

// pathComponent rejects traversal and separators in one publication path element.
func pathComponent(value string) error {
	if value == "" || value == "." || value == ".." || strings.ContainsAny(value, "/\\") || strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return fmt.Errorf("unsafe path component %q", value)
	}
	return nil
}

// newSSHPublisher resolves validated manifest identity, overrides and SSH options.
func newSSHPublisher(opts *DeployerOpts, release *contracts.Release, original *contracts.ReleaseDeploymentSsh) (*sshPublisher, error) {
	target := proto.Clone(original).(*contracts.ReleaseDeploymentSsh)
	for _, pair := range []struct {
		value string
		dest  *string
	}{
		{opts.SSHHost, &target.Host},
		{opts.SSHUser, &target.User},
		{opts.SSHRoot, &target.Root},
		{opts.PublishUser, &target.PublishUser},
		{opts.Project, &target.Project},
		{opts.Namespace, &target.Namespace},
		{opts.SiteArchive, &target.SiteArchive},
	} {
		if pair.value != "" {
			*pair.dest = pair.value
		}
	}
	if target.Root == "" {
		target.Root = "/srv/download"
	}
	namespace, project, err := publicationIdentity(release.GetProject().GetSubdir(), target.Namespace, target.Project)
	if err != nil {
		return nil, fmt.Errorf("configure public artifact identity: %w", err)
	}
	target.Namespace, target.Project = namespace, project
	if err := validateSSHTarget(opts, target, release.Name); err != nil {
		return nil, fmt.Errorf("configure SSH publication: %w", err)
	}
	args := []string{"-o", "BatchMode=yes", "-o", "StrictHostKeyChecking=yes", "-o", "ConnectTimeout=15", "-o", "ServerAliveInterval=15", "-o", "ServerAliveCountMax=3", "-o", "ForwardAgent=no", "-o", "ClearAllForwardings=yes", "-T"}
	for _, option := range []struct{ flag, value string }{{"-F", opts.SSHConfig}, {"-i", opts.SSHIdentity}, {"-l", target.User}} {
		if option.value != "" {
			args = append(args, option.flag, option.value)
		}
	}
	if opts.SSHPort != 0 {
		args = append(args, "-p", strconv.Itoa(opts.SSHPort))
	}
	if opts.SSHKnownHosts != "" {
		args = append(args, "-o", "UserKnownHostsFile="+opts.SSHKnownHosts, "-o", "GlobalKnownHostsFile=/dev/null")
	}
	if opts.SSHCertificate != "" {
		args = append(args, "-o", "CertificateFile="+opts.SSHCertificate)
	}
	if opts.SSHDisableAgent {
		args = append(args, "-o", "IdentityAgent=none", "-o", "IdentitiesOnly=yes")
	}
	siteRoot := path.Join(target.Root, "sites", target.Project)
	return &sshPublisher{
		ctx: opts.Ctx, opts: opts, target: target, args: args, version: release.Name,
		publicDir: path.Join(target.Root, target.Namespace, target.Project, "releases", release.Name),
		siteRoot:  siteRoot, siteDir: path.Join(siteRoot, "releases", release.Name),
	}, nil
}

// validateSSHTarget checks the selected destination and namespace-specific steps.
func validateSSHTarget(opts *DeployerOpts, target *contracts.ReleaseDeploymentSsh, version string) error {
	if opts.Environment == "" || (target.Environment != "" && target.Environment != opts.Environment) {
		return fmt.Errorf("select the configured SSH environment explicitly")
	}
	if err := pathComponent(target.Project); err != nil {
		return fmt.Errorf("project: %w", err)
	}
	if err := pathComponent(version); err != nil {
		return fmt.Errorf("release version: %w", err)
	}
	if !path.IsAbs(target.Root) || path.Clean(target.Root) != target.Root || target.Root == "/" || strings.IndexFunc(target.Root, unicode.IsControl) >= 0 {
		return fmt.Errorf("SSH root must be a clean absolute directory")
	}
	if target.Host == "" || strings.HasPrefix(target.Host, "-") || strings.ContainsAny(target.Host, " \t\r\n@/") {
		return fmt.Errorf("invalid SSH host")
	}
	if opts.SSHPort < 0 || opts.SSHPort > 65535 || opts.Timeout <= 0 {
		return fmt.Errorf("invalid SSH port or timeout")
	}
	if len(opts.SSHSteps) == 0 {
		return fmt.Errorf("steps must select upload, extract, or link")
	}
	previous := 0
	for _, step := range opts.SSHSteps {
		if target.Namespace == "third_party" && step != "upload" {
			return fmt.Errorf("third_party publication supports upload only")
		}
		order := map[string]int{"upload": 1, "extract": 2, "link": 3}[step]
		if order <= previous {
			return fmt.Errorf("steps must be a nonempty ordered subset of upload,extract,link: %q", step)
		}
		previous = order
		if step == "extract" {
			if err := pathComponent(target.SiteArchive); err != nil {
				return fmt.Errorf("extract site_archive: %w", err)
			}
			if !strings.HasSuffix(target.SiteArchive, ".tar.gz") && !strings.HasSuffix(target.SiteArchive, ".zip") {
				return fmt.Errorf("extract site_archive must name a published .tar.gz or .zip file")
			}
		}
	}
	return nil
}

// publicationIdentity resolves and validates the public namespace and artifact name.
func publicationIdentity(subdir, namespace, name string) (string, string, error) {
	if namespace == "" {
		namespace = "projects"
		if strings.HasPrefix(subdir, "third_party/") {
			namespace = "third_party"
		}
	}
	if namespace != "projects" && namespace != "third_party" {
		return "", "", fmt.Errorf("namespace must be projects or third_party")
	}
	if name == "" {
		name = strings.TrimPrefix(subdir, namespace+"/")
	}
	if err := pathComponent(name); err != nil {
		return "", "", fmt.Errorf("artifact name: %w", err)
	}
	return namespace, name, nil
}
