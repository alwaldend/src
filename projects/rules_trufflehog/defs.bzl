"""Scan-only TruffleHog tests with exact historical finding allowlists."""

_DEFAULT_SCANNER = Label("//:trufflehog")

def _runfile(ctx, file):
    if file.short_path.startswith("../"):
        return file.short_path[3:]
    return ctx.workspace_name + "/" + file.short_path

def _impl(ctx):
    executable = ctx.actions.declare_file(ctx.label.name)
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable._wrapper,
        is_executable = True,
    )
    dependencies = [ctx.attr._wrapper, ctx.attr.scanner, ctx.attr.git]
    runfiles = ctx.runfiles(files = [
        ctx.file.allowlist,
        ctx.file.repository_marker,
        ctx.executable.scanner,
        ctx.executable.git,
    ]).merge_all([item[DefaultInfo].default_runfiles for item in dependencies])
    return [
        DefaultInfo(executable = executable, runfiles = runfiles),
        RunEnvironmentInfo(environment = {
            "TRUFFLEHOG_SCAN_SCANNER": _runfile(ctx, ctx.executable.scanner),
            "TRUFFLEHOG_SCAN_GIT": _runfile(ctx, ctx.executable.git),
            "TRUFFLEHOG_SCAN_MARKER": _runfile(ctx, ctx.file.repository_marker),
            "TRUFFLEHOG_SCAN_ALLOWLIST": _runfile(ctx, ctx.file.allowlist),
            "TRUFFLEHOG_SCAN_ALL_REFS": "true" if ctx.attr.all_refs else "false",
        }),
    ]

_scan_test = rule(
    implementation = _impl,
    test = True,
    attrs = {
        "scanner": attr.label(default = _DEFAULT_SCANNER, executable = True, cfg = "exec", doc = "Optional scanner override; defaults to the module-owned TruffleHog toolchain."),
        "git": attr.label(mandatory = True, executable = True, cfg = "exec", doc = "Git executable used by the scanner."),
        "repository_marker": attr.label(mandatory = True, allow_single_file = True, doc = "File at the worktree root; its resolved directory identifies the repository."),
        "allowlist": attr.label(mandatory = True, allow_single_file = [".json"], doc = "Strict detector/path/commit/line finding policy."),
        "all_refs": attr.bool(default = False, doc = "Include every local ref instead of only HEAD-reachable history."),
        "_wrapper": attr.label(default = Label("//cmd/scan:scan"), executable = True, cfg = "exec"),
    },
)

def trufflehog_test(name, git, repository_marker, allowlist, scanner = _DEFAULT_SCANNER, all_refs = False, tags = [], **kwargs):
    """Runs a fresh Git-history scan without credential verification.

    Args:
        name: Test target name.
        git: Consumer-supplied Git executable.
        repository_marker: File at the worktree root.
        allowlist: JSON policy with exact detector/path/commit/line tuples.
        scanner: Optional executable override for a custom scanner or fixture.
        all_refs: Include all local refs when true.
        tags: Additional tags, such as manual for an all-ref audit.
        **kwargs: Common Bazel test attributes.
    """
    _scan_test(
        name = name,
        scanner = scanner,
        git = git,
        repository_marker = repository_marker,
        allowlist = allowlist,
        all_refs = all_refs,
        tags = ["external", "local", "no-cache"] + tags,
        **kwargs
    )
