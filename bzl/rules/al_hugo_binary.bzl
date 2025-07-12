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
        {env_script}
        ln -s '{data}' ./data
        '{hugo}' \
            --ignoreCache \
            --noChmod \
            --configDir '{config}' \
            --contentDir '{content}' \
            --layoutDir '{layouts}' \
            --themesDir '{themes}' \
            "${{@}}"
    """.format(
        hugo = hugo.hugo.short_path,
        content = info.content.short_path,
        themes = info.themes.short_path,
        config = info.config.short_path,
        data = info.data.short_path,
        layouts = info.layouts.short_path,
        env_script = ctx.expand_make_variables(
            "env_script",
            ctx.expand_location(info.env_script),
            {},
        ),
        arguments = " ".join([
            shell.quote(argument)
            for argument in ctx.attr.arguments
        ]),
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
