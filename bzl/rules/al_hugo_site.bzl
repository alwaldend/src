load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    files = [hugo.hugo, ctx.file.tree] + ctx.files.tools + ctx.files.data
    runfiles = ctx.runfiles(files)
    for tool in ctx.attr.tools:
        runfiles.merge(tool[DefaultInfo].default_runfiles)
    env_script = " && ".join([
        cmd
        for name, value in ctx.attr.env.items()
        for cmd in ['{}="{}"'.format(name, value), "export {}".format(name)]
    ])

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlHugoSiteInfo(
            tree = ctx.file.tree,
            env = ctx.attr.env,
            env_script = env_script,
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Define a hugo site",
    toolchains = [
        "//bzl/toolchain-types:hugo",
    ],
    provides = [AlHugoSiteInfo],
    attrs = {
        "tree": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo tree",
        ),
        "tools": attr.label_list(
            doc = "Tools that should be available for the build",
            default = [],
        ),
        "data": attr.label_list(
            doc = "Data that is made available for the site",
            default = [],
        ),
        "env": attr.string_dict(
            default = {},
            doc = """
                Hugo environment variables
                (support location statements, support make variables, support
                shell commands)
            """,
        ),
    },
)
