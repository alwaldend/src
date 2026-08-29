"""Implementation of Promptfoo test rules."""

load(
    ":workspace.bzl",
    "PromptfooWorkspaceInfo",
    "promptfoo_workspace",
)

_RUNFILES_INIT = """\
set -uo pipefail; set +e; f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \\
  source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null || \\
  source "$0.runfiles/$f" 2>/dev/null || \\
  source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \\
  source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \\
  { echo >&2 "ERROR: cannot find $f"; exit 1; }; f=; set -e
runfiles_export_envvars
root="$(runfiles_current_repository || true)"
if [[ -z "$root" ]]; then
  root="_main"
fi
"""

_SCRIPT = """\
#!/usr/bin/env bash
set -euo pipefail
umask 077
{runfiles_init}

promptfoo={promptfoo}
config={config}

scratch_root="${{TEST_TMPDIR:?TEST_TMPDIR is required}}"
if [[ "$scratch_root" != /* ]]; then
  echo >&2 "TEST_TMPDIR must be an absolute path"
  exit 2
fi
if [[ ! -d "$scratch_root" || ! -w "$scratch_root" ]]; then
  echo >&2 "TEST_TMPDIR must be an existing writable directory"
  exit 2
fi
scratch_root="$(cd "$scratch_root" && pwd -P)"
ancestor="$scratch_root"
while true; do
  if [[ -e "$ancestor/.agents" || -L "$ancestor/.agents" ]]; then
    echo >&2 "TEST_TMPDIR must not be below a tree containing .agents"
    exit 2
  fi
  if [[ "$ancestor" == / ]]; then
    break
  fi
  ancestor="$(dirname -- "$ancestor")"
done
state=""
cleanup() {{
  status=$?
  trap - EXIT HUP INT TERM
  if [[ -n "$state" ]]; then
    if ! rm -rf -- "$state"; then
      echo >&2 "ERROR: could not remove Promptfoo runner state: $state"
      if [[ "$status" == 0 ]]; then
        status=1
      fi
    fi
  fi
  exit "$status"
}}
trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM
state="$(mktemp -d "$scratch_root/rules_promptfoo.XXXXXXXX")"
readonly state
chmod 700 "$state"
workspace="$state/workspace"
judge_workspace="$state/judge_workspace"
subject_codex_home="$state/subject_codex"
judge_codex_home="$state/judge_codex"
mkdir -p \
  "$state/config" \
  "$state/cache" \
  "$state/tmp" \
  "$subject_codex_home" \
  "$judge_codex_home" \
  "$judge_workspace" \
  "$workspace"
if [[ {is_eval} == 1 ]]; then
  mkdir -p "${{TEST_UNDECLARED_OUTPUTS_DIR:?TEST_UNDECLARED_OUTPUTS_DIR is required}}"
fi
{copy_commands}
export PROMPTFOO_CONFIG_DIR="$state/config"
export PROMPTFOO_CACHE_PATH="$state/cache"
export PROMPTFOO_DISABLE_TELEMETRY=1
export PROMPTFOO_DISABLE_UPDATE=1
export PROMPTFOO_JUDGE_CODEX_HOME="$judge_codex_home"
export PROMPTFOO_JUDGE_WORKSPACE="$judge_workspace"
export PROMPTFOO_STATE_DIR="$state"
export PROMPTFOO_SUBJECT_CODEX_HOME="$subject_codex_home"
export PROMPTFOO_SKILL_WORKSPACE="$workspace"
export TMPDIR="$state/tmp"
export TMP="$TMPDIR"
export TEMP="$TMPDIR"
# rules_js binaries normally start in Bazel's output tree. This is a test
# process, so permit the launcher to run from the isolated skill workspace.
export BAZEL_BINDIR=.
{codex_auth_fixture}
{codex_login}
export CODEX_HOME="$subject_codex_home"

args=(
  {command}
  --config "$config"
{arguments}
{safety_flags}
)
cd "$workspace"
"$promptfoo" "${{args[@]}}"
"""

def _shell_quote(value):
    return "'{}'".format(value.replace("'", "'\"'\"'"))

def _rlocation_expression(file):
    if file.short_path.startswith("../"):
        return "$(rlocation {})".format(_shell_quote(file.short_path[3:]))
    return "$(rlocation \"$root\"/{})".format(
        _shell_quote(file.short_path),
    )

def _promptfoo_runner_test_impl(ctx):
    workspace = ctx.attr.workspace[PromptfooWorkspaceInfo]
    executable = ctx.actions.declare_file("{}.sh".format(ctx.label.name))
    safety_flags = []
    if ctx.attr.command == "eval":
        safety_flags = [
            "--no-cache",
            "--no-write",
            "--no-share",
        ]
        if ctx.attr.reuse_codex_login:
            safety_flags.extend([
                "--max-concurrency",
                "1",
            ])

    copy_commands = []
    for destination in sorted(workspace.files_by_path.keys()):
        source = workspace.files_by_path[destination]
        copy_commands.extend([
            "source={}".format(_rlocation_expression(source)),
            "destination=\"$workspace\"/{}".format(
                _shell_quote(destination),
            ),
            "mkdir -p \"$(dirname \"$destination\")\"",
            "cp -- \"$source\" \"$destination\"",
        ])

    ctx.actions.write(
        output = executable,
        is_executable = True,
        content = _SCRIPT.format(
            arguments = "\n".join([
                "  {}".format(_shell_quote(argument))
                for argument in ctx.attr.promptfoo_args
            ]),
            codex_auth_fixture = """\
fixture_codex_home="$state/fixture_codex_home"
fixture_auth_dir="$fixture_codex_home/real"
mkdir -p "$fixture_auth_dir"
cp -- {auth_fixture} "$fixture_auth_dir/auth.json"
chmod 600 "$fixture_auth_dir/auth.json"
ln -s -- "real/auth.json" "$fixture_codex_home/auth.json"
export CODEX_HOME="$fixture_codex_home"
""".format(
                auth_fixture = _rlocation_expression(ctx.file.codex_auth_fixture),
            ) if ctx.file.codex_auth_fixture else "",
            codex_login = """\
host_codex_home="${CODEX_HOME:-}"
if [[ -z "$host_codex_home" ]]; then
  echo >&2 "CODEX_HOME is required when reuse_codex_login is enabled"
  exit 2
fi
host_codex_auth="$host_codex_home/auth.json"
if [[ "$host_codex_auth" != /* ]]; then
  host_codex_auth="$(rlocation "$root/$host_codex_auth" || true)"
fi
if [[ -z "$host_codex_auth" || ! -f "$host_codex_auth" ]]; then
  echo >&2 "CODEX_HOME/auth.json is required when reuse_codex_login is enabled"
  exit 2
fi
host_codex_auth="$(readlink -f -- "$host_codex_auth")"
if [[ ! -w "$host_codex_auth" ]]; then
  echo >&2 "CODEX_HOME/auth.json must be writable so Codex can persist refreshes"
  exit 2
fi
# Lock the stable login directory so serialization does not depend on the
# current auth file inode.
host_codex_home="$(dirname "$host_codex_auth")"
exec {auth_lock_fd}<"$host_codex_home"
flock -x "$auth_lock_fd"
export PROMPTFOO_ASSERTIONS_MAX_CONCURRENCY=1
export PROMPTFOO_SUGGESTIONS_MAX_CONCURRENCY=1
for isolated_codex_home in "$subject_codex_home" "$judge_codex_home"; do
  ln -s -- "$host_codex_auth" "$isolated_codex_home/auth.json"
done
""" if ctx.attr.reuse_codex_login else "",
            command = _shell_quote(ctx.attr.command),
            config = _rlocation_expression(ctx.file.config),
            copy_commands = "\n".join(copy_commands),
            is_eval = "1" if ctx.attr.command == "eval" else "0",
            promptfoo = _rlocation_expression(ctx.executable.promptfoo),
            runfiles_init = _RUNFILES_INIT,
            safety_flags = "\n".join([
                "  {}".format(_shell_quote(flag))
                for flag in safety_flags
            ] + ([
                "  --output \"$TEST_UNDECLARED_OUTPUTS_DIR/results.json\"",
            ] if ctx.attr.command == "eval" else [])),
        ),
    )

    direct_files = [ctx.file.config] + ctx.files.data
    if ctx.file.codex_auth_fixture:
        direct_files.append(ctx.file.codex_auth_fixture)
    runfiles = ctx.runfiles(
        files = direct_files,
    ).merge_all([
        ctx.attr.promptfoo[DefaultInfo].default_runfiles,
        ctx.attr.shell_runfiles[DefaultInfo].default_runfiles,
        ctx.attr.workspace[DefaultInfo].default_runfiles,
    ] + [
        data[DefaultInfo].default_runfiles
        for data in ctx.attr.data
    ])

    return [
        DefaultInfo(
            executable = executable,
            runfiles = runfiles,
        ),
        RunEnvironmentInfo(
            environment = ctx.attr.env,
            inherited_environment = ctx.attr.env_inherit,
        ),
    ]

promptfoo_runner_test = rule(
    implementation = _promptfoo_runner_test_impl,
    attrs = {
        "command": attr.string(
            mandatory = True,
            values = ["eval", "validate"],
        ),
        "codex_auth_fixture": attr.label(
            allow_single_file = True,
            doc = "Test-only fake auth file copied into writable runner state.",
        ),
        "config": attr.label(
            allow_single_file = [".json", ".yaml", ".yml"],
            mandatory = True,
        ),
        "data": attr.label_list(allow_files = True),
        "env": attr.string_dict(
            doc = "Non-secret environment values for the test.",
        ),
        "env_inherit": attr.string_list(
            doc = "Host environment names to inherit at runtime.",
        ),
        "reuse_codex_login": attr.bool(
            doc = (
                "Share only CODEX_HOME/auth.json with the otherwise " +
                "isolated runner homes for this serialized local test; the " +
                "Codex executable must preserve auth-file symlinks on save."
            ),
        ),
        "promptfoo": attr.label(
            cfg = "exec",
            default = Label("//:promptfoo"),
            executable = True,
        ),
        "promptfoo_args": attr.string_list(
            doc = "Additional literal Promptfoo CLI arguments.",
        ),
        "shell_runfiles": attr.label(
            default = Label("@bazel_tools//tools/bash/runfiles"),
        ),
        "workspace": attr.label(
            mandatory = True,
            providers = [[PromptfooWorkspaceInfo]],
        ),
    },
    doc = "Runs Promptfoo against a staged set of Codex skills.",
    test = True,
)

def _workspace_name(name):
    return "{}_skill_workspace".format(name)

_RUNNER_OWNED_ENVIRONMENT = [
    "TEMP",
    "TEST_TMPDIR",
    "TEST_UNDECLARED_OUTPUTS_DIR",
    "TMP",
    "TMPDIR",
]

def _create_test(
        name,
        command,
        config,
        skills,
        args,
        data,
        env,
        env_inherit,
        reuse_codex_login,
        tags,
        promptfoo = None,
        **kwargs):
    for environment_name in _RUNNER_OWNED_ENVIRONMENT:
        if environment_name in env or environment_name in env_inherit:
            fail(
                "{} is owned by the Promptfoo runner".format(
                    environment_name,
                ),
            )
    workspace_name = _workspace_name(name)
    promptfoo_workspace(
        name = workspace_name,
        skills = skills,
        testonly = True,
        visibility = ["//visibility:private"],
    )
    runner_kwargs = dict(kwargs)
    if promptfoo != None:
        runner_kwargs["promptfoo"] = promptfoo
    promptfoo_runner_test(
        name = name,
        command = command,
        config = config,
        data = data,
        env = env,
        env_inherit = env_inherit,
        reuse_codex_login = reuse_codex_login,
        promptfoo_args = args,
        tags = tags,
        workspace = ":{}".format(workspace_name),
        **runner_kwargs
    )

def promptfoo_test(
        name,
        config,
        skills = [],
        args = [],
        data = [],
        env = {},
        env_inherit = [],
        reuse_codex_login = False,
        promptfoo = None,
        tags = [],
        **kwargs):
    """Declares a manual Promptfoo evaluation test.

    Args:
        name: Target name.
        config: Promptfoo configuration file.
        skills: rules_skill targets staged for Codex discovery.
        args: Additional literal Promptfoo CLI arguments.
        data: Runtime files needed by the Promptfoo configuration.
        env: Non-secret environment variables.
        env_inherit: Secret-bearing host environment variable names to inherit.
        reuse_codex_login: Whether to share only the inherited local Codex login
            with the otherwise isolated runner homes for this serialized test;
            the Codex executable must preserve auth-file symlinks on save.
        promptfoo: Optional Promptfoo-compatible executable, useful in tests.
        tags: Additional Bazel tags.
        **kwargs: Additional test attributes.
    """
    if reuse_codex_login and not (
        "CODEX_HOME" in env or
        "CODEX_HOME" in env_inherit
    ):
        fail(
            "promptfoo_test with reuse_codex_login requires CODEX_HOME in " +
            "env or env_inherit",
        )
    _create_test(
        name = name,
        command = "eval",
        config = config,
        skills = skills,
        args = args,
        data = data,
        env = env,
        env_inherit = env_inherit,
        reuse_codex_login = reuse_codex_login,
        promptfoo = promptfoo,
        tags = tags + [
            "local",
            "manual",
            "no-cache",
            "no-remote",
            "requires-network",
        ],
        **kwargs
    )

def promptfoo_validate_test(
        name,
        config,
        skills = [],
        args = [],
        data = [],
        env = {},
        env_inherit = [],
        reuse_codex_login = False,
        promptfoo = None,
        tags = [],
        **kwargs):
    """Declares an offline test that validates a Promptfoo configuration."""
    if env_inherit:
        fail(
            "promptfoo_validate_test does not accept env_inherit; " +
            "validation must not inherit host credentials",
        )
    if reuse_codex_login:
        fail(
            "promptfoo_validate_test does not accept reuse_codex_login; " +
            "validation must not receive host credentials",
        )
    _create_test(
        name = name,
        command = "validate",
        config = config,
        skills = skills,
        args = args,
        data = data,
        env = env,
        env_inherit = env_inherit,
        reuse_codex_login = reuse_codex_login,
        promptfoo = promptfoo,
        tags = tags,
        **kwargs
    )
