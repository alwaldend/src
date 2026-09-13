_SCRIPT = """\
#!/usr/bin/env bash

# https://github.com/bazelbuild/rules_shell/blob/main/shell/private/sh_executable.bzl
# https://github.com/bazel-contrib/rules_shell/blob/main/shell/runfiles/runfiles.bash
# --- begin runfiles.bash initialization v3 ---
if [ -z "${{RUNFILES_DIR:-}}" ] && [ -d "${{RUNFILES:-}}" ]; then
  export RUNFILES_DIR="$RUNFILES"
fi
set -uo pipefail; set +e; f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck disable=SC1090
source "${{RUNFILES_DIR:-${{RUNFILES:-/dev/null}}}}/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "${{RUNFILES_MANIFEST_FILE:-/dev/null}}" | cut -f2- -d' ')" 2>/dev/null || \
  source "$0.runfiles/$f" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
  source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
  {{ echo>&2 "ERROR: cannot find $f"; exit 1; }}; f=; set -e
# --- end runfiles.bash initialization v3 ---
runfiles_export_envvars
bin="{bin}"
if [[ "$bin" == ../* ]]; then
  bin="${{bin#../}}"
else
  root="$(runfiles_current_repository || true)"
  if [ -z "${{root}}" ]; then
    root="_main"
  fi
  bin="${{root}}/$bin"
fi
exec "$(rlocation "$bin")" {arguments} "${{@}}"
"""

def _impl(ctx):
    toolchain = ctx.toolchains[ctx.attr.toolchain_type]
    script = ctx.actions.declare_file("{}_/{}".format(ctx.label.name, toolchain.binary.basename))
    runfiles = ctx.runfiles().merge_all(
        [
            toolchain.default_info.default_runfiles,
            ctx.attr.shell_runfiles[DefaultInfo].default_runfiles,
            ctx.runfiles(
                transitive_files = ctx.attr.shell_runfiles[DefaultInfo].files,
            ),
            ctx.runfiles(files = ctx.files.data),
        ] + [
            data[DefaultInfo].default_runfiles
            for data in ctx.attr.data
        ],
    )
    ctx.actions.write(
        output = script,
        content = _SCRIPT.format(
            bin = toolchain.binary.short_path,
            arguments = " ".join([
                '"{}"'.format(
                    ctx.expand_make_variables(
                        str(ctx.label),
                        ctx.expand_location(arg, ctx.attr.data),
                        {},
                    ),
                )
                for arg in ctx.attr.arguments
            ]),
        ),
    )
    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
        RunEnvironmentInfo(
            environment = ctx.attr.env,
        ),
    ]

def binary_toolchain_binary(toolchain_type, kwargs = {}, attrs = {}):
    """
    Generate a binary rule for a binary_toolchain

    Args:
        toolchain_type: toolchain type label
        kwargs: override for rule kwargs
        attrs: override for attrs for the rule

    Returns:
        Executable rule
    """
    default_kwargs = {
        "implementation": _impl,
        "executable": True,
        "doc": "Binary toolchain for {}".format(toolchain_type),
        "toolchains": [toolchain_type],
    }
    default_attrs = {
        "toolchain_type": attr.string(
            doc = "Toolchain type",
            default = toolchain_type,
        ),
        "data": attr.label_list(
            allow_files = True,
            doc = "Data that is added runfiles",
        ),
        "arguments": attr.string_list(
            doc = "Arguments (make variables and locations are expanded)",
        ),
        "env": attr.string_dict(
            doc = "Environment variables to set when running the binary.",
        ),
        "shell_runfiles": attr.label(
            doc = "Shell runfiles",
            default = "@rules_shell//shell/runfiles",
        ),
    }
    kwargs = default_kwargs | kwargs
    kwargs["attrs"] = default_attrs | attrs
    res = rule(**kwargs)
    return res
