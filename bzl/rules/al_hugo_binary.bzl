load("@bazel_skylib//lib:shell.bzl", "shell")
load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    default_info = ctx.attr.site[DefaultInfo]
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    runfiles = ctx.runfiles(files = [hugo.hugo] + default_info.files.to_list())
    runfiles.merge(default_info.default_runfiles)
    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))

    script_content = """\
        #!/usr/bin/env sh
        set -eux
        ln -s '{content}' ./content
        ln -s '{themes}' ./themes
        ln -s '{config}' ./
        '{hugo}' {arguments} "$${{@}}"
    """.format(
        hugo = hugo.hugo.short_path,
        content = info.content.short_path,
        themes = info.themes.short_path,
        config = info.config.short_path,
        arguments = " ".join([
            shell.quote(argument)
            for argument in ctx.attr.arguments
        ]),
    )
    script_content = ctx.expand_make_variables(
        "script",
        ctx.expand_location(script_content),
        {},
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )

    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = script,
        ),
    ]

al_hugo_binary = rule(
    implementation = _impl,
    executable = True,
    doc = "Run a hugo command",
    toolchains = [
        "//bzl/toolchain-types:hugo",
    ],
    attrs = {
        "arguments": attr.string_list(
            mandatory = True,
            doc = "Hugo arguments",
        ),
        "site": attr.label(
            mandatory = True,
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
    },
)
