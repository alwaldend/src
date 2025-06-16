load("//bzl/providers:al_toml_info.bzl", "AlTomlInfo")

def _run(ctx, script, tomlv, file):
    output = ctx.actions.declare_file("{}-tomlv-output".format(ctx.label.name))
    ctx.actions.run(
        executable = script,
        inputs = [file, tomlv],
        outputs = [output],
        arguments = [output.path, tomlv.path, file.path],
        progress_message = "Run tomlv on {}".format(file),
    )
    return output

def _impl(target, ctx):
    if not hasattr(ctx.rule.attr, "srcs"):
        return []

    tomlv = ctx.rule.executable.tomlv
    outputs = []
    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = """\
            #!/usr/bin/env sh
            set -eu
            output="${1}"
            shift
            "${@}" 2>&1 | tee "${output}"
        """,
    )
    for file in ctx.rule.files.srcs:
        outputs.append(_run(ctx, script, tomlv, file))
    return [
        OutputGroupInfo(
            al_toml_validate = depset(direct = outputs),
        ),
    ]

al_toml_validate = aspect(
    implementation = _impl,
    attr_aspects = ["deps"],
    required_providers = [AlTomlInfo],
    doc = "Aspect adding linters for toml files",
)
