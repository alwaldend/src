load("//bzl/providers:al_hugo_theme_info.bzl", "AlHugoThemeInfo")

def _impl(ctx):
    hugo_args = ctx.actions.args()
    cmd = ["set -eux", 'echo "${@}"', 'head -n 100 "${0}"', "exit 1"]

    cmd.append("mkdir themes")
    hugo_args.add_all(["--themesDir", "themes"])
    for theme in ctx.attr.themes:
        info = theme[AlHugoThemeInfo]
        dir = "themes/{}".format(info.name)
        cmd.append("mkdir '{}'".format(dir))
        for src in info.srcs:
            file, relative_path = src.file, src.relative_path
            cmd.append('mkdir -p "{}/$(dirname \'{}\')"'.format(dir, relative_path))
            cmd.append("ln -sf '{}' '{}/{}'".format(file.path, dir, relative_path))

    cmd.append("mkdir content")
    hugo_args.add_all(["--contentDir", "content"])
    for src in ctx.files.srcs:
        cmd.append("tar -xf '{}' -C 'content'".format(src.path))

    for config in ctx.files.configs:
        hugo_args.add_all(["--config", config.path])

    cmd.append("mkdir out")
    hugo_args.add_all(["--destination", "out"])
    cmd.append('set -x && \'{}\' "${{@}}"'.format(ctx.executable.hugo.path))
    cmd.append("tar -cf '{}' -C out".format(ctx.outputs.out))

    ctx.actions.run_shell(
        inputs = ctx.files.srcs + ctx.files.themes + ctx.files.configs,
        outputs = [ctx.outputs.out],
        tools = [ctx.executable.hugo],
        arguments = [hugo_args],
        command = "\n".join(cmd),
    )

al_hugo_build = rule(
    implementation = _impl,
    doc = "Run a hugo build",
    toolchains = [
        # "@com-alwaldend-src-nodejs-toolchains//:resolved_toolchain",
    ],
    attrs = {
        "configs": attr.label_list(
            allow_files = True,
            doc = "Hugo config files for --config",
        ),
        "out": attr.output(
            doc = "Build output (tar)",
            mandatory = True,
        ),
        "srcs": attr.label_list(
            allow_files = [".tar"],
            mandatory = True,
            doc = "Hugo content for --contentDir",
        ),
        "themes": attr.label_list(
            default = [],
            doc = "Themes",
            providers = [AlHugoThemeInfo],
        ),
        "hugo": attr.label(
            default = "//tools:hugo",
            executable = True,
            cfg = "exec",
        ),
    },
)
