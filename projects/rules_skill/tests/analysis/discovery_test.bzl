"""Execution test for a generated skill discovery link updater."""

def _skill_discovery_updater_test_impl(ctx):
    executable = ctx.actions.declare_file(ctx.label.name + ".sh")
    updater = ctx.executable.updater
    lines = [
        "#!/usr/bin/env bash",
        "set -euo pipefail",
        "",
        "updater=\"${TEST_SRCDIR}/${TEST_WORKSPACE}/" + updater.short_path + "\"",
        "workspace=\"${TEST_TMPDIR}/workspace\"",
        "canonical=\"${workspace}/tests/analysis\"",
        "discovery=\"${workspace}/tests/discovery\"",
        "mkdir -p \"${canonical}\" \"${discovery}\"",
        "printf '%s\\n' fixture >\"${canonical}/SKILL.md\"",
        "ln -s '../stale' \"${discovery}/stale\"",
        "",
        "BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\"",
        "[[ -d \"${discovery}\" && ! -L \"${discovery}\" ]]",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../analysis' ]]",
        "[[ ! -e \"${discovery}/stale\" && ! -L \"${discovery}/stale\" ]]",
        "",
        "BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\"",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../analysis' ]]",
        "",
        "mkdir \"${discovery}.lock\"",
        "unlink \"${discovery}/analysis\"",
        "ln -s '../locked' \"${discovery}/analysis\"",
        "if BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\" \\",
        "    >\"${TEST_TMPDIR}/locked.stdout\" \\",
        "    2>\"${TEST_TMPDIR}/locked.stderr\"; then",
        "    echo 'updater ignored an active discovery lock' >&2",
        "    exit 1",
        "fi",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../locked' ]]",
        "rmdir \"${discovery}.lock\"",
        "BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\"",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../analysis' ]]",
        "",
        "unlink \"${discovery}/analysis\"",
        "ln -s '../wrong' \"${discovery}/analysis\"",
        "printf '%s\\n' collision >\"${discovery}/unexpected\"",
        "if BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\" \\",
        "    >\"${TEST_TMPDIR}/unexpected.stdout\" \\",
        "    2>\"${TEST_TMPDIR}/unexpected.stderr\"; then",
        "    echo 'updater accepted an unexpected regular file' >&2",
        "    exit 1",
        "fi",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../wrong' ]]",
        "unlink \"${discovery}/unexpected\"",
        "BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\"",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../analysis' ]]",
        "",
        "unlink \"${discovery}/analysis\"",
        "rmdir \"${discovery}\"",
        "mkdir \"${workspace}/tests/legacy-discovery\"",
        "ln -s 'legacy-discovery' \"${discovery}\"",
        "BUILD_WORKSPACE_DIRECTORY=\"${workspace}\" \"${updater}\"",
        "[[ -d \"${discovery}\" && ! -L \"${discovery}\" ]]",
        "[[ \"$(readlink \"${discovery}/analysis\")\" == '../analysis' ]]",
        "",
    ]
    ctx.actions.write(
        output = executable,
        content = "\n".join(lines),
        is_executable = True,
    )

    runfiles = ctx.runfiles(files = [updater])
    runfiles = runfiles.merge(ctx.attr.updater[DefaultInfo].default_runfiles)
    return [DefaultInfo(executable = executable, runfiles = runfiles)]

skill_discovery_updater_test = rule(
    implementation = _skill_discovery_updater_test_impl,
    attrs = {
        "updater": attr.label(
            cfg = "target",
            executable = True,
            mandatory = True,
        ),
    },
    test = True,
)
