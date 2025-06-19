load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    runfiles = ctx.runfiles(
        files = [ctx.executable.hugo, info.themes, info.content] + ctx.files.configs,
    )

    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = """\
            #!/usr/bin/env sh
            set -eux

            '{hugo}' serve \
                --config '{configs}' \
                --themesDir '{themes}' \
                --contentDir '{content}'
        """.format(
            hugo = ctx.executable.hugo.short_path,
            configs = ",".join([config.short_path for config in ctx.files.configs]),
            themes = info.themes.short_path,
            content = info.content.short_path,
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
    attrs = {
        "configs": attr.label_list(
            allow_files = True,
            mandatory = True,
            doc = "Hugo configs",
        ),
        "site": attr.label(
            mandatory = True,
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
        "hugo": attr.label(
            default = "@com-github-gohugoio-hugo//:bin",
            executable = True,
            cfg = "exec",
            doc = "Hugo binary",
        ),
    },
)
