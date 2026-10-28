_SCRIPT = """\
#!/usr/bin/env bash
set -euo pipefail
if [ -z "${{RUNFILES_DIR:-}}" ]; then
    RUNFILES_DIR="$0.runfiles"
fi

f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
if ! source "$0.runfiles/$f" 2>/dev/null; then
    echo >&2 "ERROR: cannot find $f"
    exit 1
fi
runfiles_export_envvars

gradle="$(rlocation "{gradle}")"
java="$(rlocation "{java}")"
java_home="$(dirname "$(dirname "$java")")"

export JAVA_HOME="$java_home"
export PATH="$java_home/bin:$PATH"
cache_root="${{GRADLE_CACHE_ROOT:-${{BUILD_WORKSPACE_DIRECTORY:-$PWD}}/out/gradle}}"
mkdir -p "$cache_root"

exec "$gradle" \
    --no-daemon \
    -Dorg.gradle.java.home="$java_home" \
    --project-cache-dir "$cache_root/project-cache" \
    --gradle-user-home "$cache_root/gradle-home" \
    "$@"
"""

def _runfile_path(path, workspace_name):
    if path.startswith("../"):
        return path[3:]
    return workspace_name + "/" + path

def _gradle_wrapper_binary_impl(ctx):
    gradle = ctx.file.gradle
    java_home = ctx.attr.java_home[platform_common.ToolchainInfo]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles().merge_all([
        ctx.attr.java_home[DefaultInfo].default_runfiles,
        ctx.attr._bash_runfiles[DefaultInfo].default_runfiles,
    ] + [
        data[DefaultInfo].default_runfiles
        for data in ctx.attr.data
    ]).merge(ctx.runfiles(
        files = [gradle],
        transitive_files = ctx.attr._bash_runfiles[DefaultInfo].files,
    ))
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = _SCRIPT.format(
            gradle = _runfile_path(gradle.short_path, ctx.workspace_name),
            java = _runfile_path(java_home.java, ctx.workspace_name),
        ),
    )
    return [DefaultInfo(executable = script, runfiles = runfiles)]

gradle_wrapper_binary = rule(
    doc = "Run Gradle from a pinned Bazel distribution and Java runtime.",
    implementation = _gradle_wrapper_binary_impl,
    executable = True,
    attrs = {
        "gradle": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "java_home": attr.label(
            mandatory = True,
            providers = [platform_common.ToolchainInfo],
        ),
        "data": attr.label_list(
            allow_files = True,
        ),
        "_bash_runfiles": attr.label(
            default = "@rules_shell//shell/runfiles",
        ),
    },
)
