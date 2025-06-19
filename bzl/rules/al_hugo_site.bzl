load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    content_dir = ctx.actions.declare_directory("{}-content".format(ctx.label.name))
    ctx.actions.run_shell(
        inputs = [ctx.file.content],
        outputs = [content_dir],
        command = "tar -xf '{}' -C '{}'".format(ctx.file.content.path, content_dir.path),
    )

    themes_dir = ctx.actions.declare_directory("{}-themes".format(ctx.label.name))
    ctx.actions.run_shell(
        inputs = [ctx.file.themes],
        outputs = [themes_dir],
        command = "tar -xf '{}' -C '{}'".format(ctx.file.themes.path, themes_dir.path),
    )

    files = [content_dir, themes_dir]
    return [
        DefaultInfo(
            files = depset(files),
            runfiles = ctx.runfiles(files),
        ),
        AlHugoSiteInfo(
            content = content_dir,
            themes = themes_dir,
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Define a hugo site",
    toolchains = [
        # "@com-alwaldend-src-nodejs-toolchains//:resolved_toolchain",
    ],
    provides = [AlHugoSiteInfo],
    attrs = {
        "content": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo content tree",
        ),
        "themes": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo themes tree",
        ),
    },
)
