load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    runfiles = ctx.runfiles(
        files = [hugo.hugo, info.build],
    )
    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))

    ctx.actions.write(
        output = script,
        is_executable = True,
        content = """\
            #!/usr/bin/env sh
            set -eux

            '{hugo}' serve --source '{build}' "${{@}}"
        """.format(
            hugo = hugo.hugo.short_path,
            build = info.build.short_path,
        ),
    )

    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = script,
        ),
    ]

al_hugo_serve = rule(
    implementation = _impl,
    executable = True,
    doc = "Serve a hugo site",
    toolchains = ["//bzl/toolchain-types:hugo"],
    attrs = {
        "site": attr.label(
            mandatory = True,
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
    },
)
