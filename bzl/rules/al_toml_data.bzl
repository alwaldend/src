load("//bzl/providers:al_toml_info.bzl", "AlTomlInfo")

def _impl(ctx):
    tomlv_script = ctx.actions.declare_file("{}-tomlv-script".format(ctx.label.name))
    tomlv_mark = ctx.actions.declare_file("{}-tomlv-mark".format(ctx.label.name))
    tomlv_script_content = """\
        #!/usr/bin/env sh
        set -eu
        '{tomlv}' "${{@}}"
        touch '{output}'
    """
    ctx.actions.write(
        output = tomlv_script,
        is_executable = True,
        content = tomlv_script_content.format(tomlv = ctx.executable.tomlv.path, output = tomlv_mark.path),
    )
    ctx.actions.run(
        executable = tomlv_script,
        inputs = [ctx.executable.tomlv] + ctx.files.srcs,
        outputs = [tomlv_mark],
        arguments = [file.path for file in ctx.files.srcs],
        progress_message = "Run tomlv on {}".format(ctx.label),
    )
    return [
        DefaultInfo(files = depset(direct = ctx.files.srcs, transitive = ctx.files.deps)),
        AlTomlInfo(srcs = ctx.files.srcs, deps = ctx.files.deps),
        OutputGroupInfo(validation = depset([tomlv_mark])),
    ]

al_toml_data = rule(
    implementation = _impl,
    provides = [AlTomlInfo, DefaultInfo],
    attrs = {
        "srcs": attr.label_list(doc = "Toml files", allow_files = [".toml"]),
        "deps": attr.label_list(doc = "Toml data targets", providers = [AlTomlInfo]),
        "tomlv": attr.label(
            executable = True,
            doc = "Tomlv target to use for validation",
            cfg = "exec",
            default = "//tools:tomlv",
        ),
    },
)
