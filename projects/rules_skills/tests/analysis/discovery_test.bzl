"""Execution tests for the skills_write updater script."""

def _skills_write_updater_test_impl(ctx):
    updater = ctx.executable.updater
    discovery = ctx.attr.discovery_dir
    environment = {}
    if RunEnvironmentInfo in ctx.attr.updater:
        environment = ctx.attr.updater[RunEnvironmentInfo].environment

    script = ctx.actions.declare_file(ctx.label.name + ".sh")
    lines = [
        "#!/usr/bin/env bash",
        "set -euo pipefail",
        "",
        'runfiles="${TEST_SRCDIR}/${TEST_WORKSPACE}"',
        'updater="${runfiles}/' + updater.short_path + '"',
        'workspace="${TEST_TMPDIR}/workspace"',
        "discovery=" + discovery,
        'mkdir -p "${workspace}/tests/analysis"',
        "printf '%s\\n' name: analysis >\"${workspace}/tests/analysis/SKILL.md\"",
        "printf '%s\\n' '---' >\"${workspace}/MODULE.bazel\"",
        "ln -s '../stale' \"${workspace}/${discovery}/stale\" 2>/dev/null || true",
        "",
    ]
    for name in sorted(environment.keys()):
        lines.append("export {}=\"{}\"".format(name, environment[name]))
    lines.extend([
        "run() {",
        '    BUILD_WORKSPACE_DIRECTORY="${workspace}" "${updater}"',
        "}",
        "",
        "run",
        '[[ ! -e "${workspace}/${discovery}/stale" && ! -L "${workspace}/${discovery}/stale" ]]',
        "",
        "run",
        "",
        'mkdir "${workspace}/${discovery}.lock"',
        'if run >"${TEST_TMPDIR}/locked.stdout" 2>"${TEST_TMPDIR}/locked.stderr"; then',
        "    echo 'updater ignored an active discovery lock' >&2",
        "    exit 1",
        "fi",
        'rmdir "${workspace}/${discovery}.lock"',
        "run",
        "",
    ])
    if ctx.attr.expected_kind == "symlink":
        lines.append('[[ -L "${workspace}/${discovery}/analysis" ]]')
    else:
        lines.extend([
            '[[ -d "${workspace}/${discovery}/analysis" && ! -L "${workspace}/${discovery}/analysis" ]]',
            '[[ -f "${workspace}/${discovery}/analysis/SKILL.md" ]]',
            "printf '%s\\n' tampered >\"${workspace}/${discovery}/analysis/SKILL.md\"",
        ])
    lines.append("")
    ctx.actions.write(
        output = script,
        content = "\n".join(lines),
        is_executable = True,
    )

    runfiles = ctx.runfiles(files = [updater])
    runfiles = runfiles.merge(ctx.attr.updater[DefaultInfo].default_runfiles)
    return [DefaultInfo(executable = script, runfiles = runfiles)]

skills_write_updater_test = rule(
    implementation = _skills_write_updater_test_impl,
    attrs = {
        "discovery_dir": attr.string(default = "tests/discovery"),
        "expected_kind": attr.string(default = "symlink"),
        "updater": attr.label(
            cfg = "target",
            executable = True,
            mandatory = True,
        ),
    },
    test = True,
)
