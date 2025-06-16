load("//bzl/providers:al_toml_info.bzl", "AlTomlInfo")

def _impl(target, ctx):
    if ctx.label.repo_name or not hasattr(ctx.rule.attr, "srcs"):
        return []

    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = """\
            #!/usr/bin/env sh
            set -eu
            mark="${1}"
            shift
            "${@}"
            touch "${mark}"
        """,
    )

    tomlv = ctx.rule.executable.tomlv
    outputs = []
    for i, file in enumerate(ctx.rule.files.srcs):
        output = ctx.actions.declare_file("{}-tomlv-mark-{}".format(ctx.label.name, i))
        outputs.append(output)
        ctx.actions.run(
            executable = script,
            inputs = [file, tomlv],
            outputs = [output],
            arguments = [output.path, tomlv.path, file.path],
            progress_message = "Run tomlv on {}".format(file),
        )

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
