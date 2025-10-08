load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    themes = ctx.actions.declare_directory("{}.themes".format(ctx.label.name))
    files = [themes] + ctx.files.site
    runfiles = ctx.runfiles()
    for tool in ctx.attr.tools:
        runfiles.merge(tool[DefaultInfo].default_runfiles)
    runfiles.merge(ctx.attr.postcss[DefaultInfo].default_runfiles)

    ctx.actions.run_shell(
        inputs = [ctx.file.themes],
        outputs = [themes],
        command = "tar -xf '{}' -C '{}'".format(ctx.file.themes.path, themes.path),
    )

    path = ":".join((
        "$$(dirname '{}')".format(ctx.executable.postcss.path),
        "$$(dirname '{}')".format(ctx.executable.postcss.short_path),
        ctx.attr.env.get("PATH", ""),
        "$${PATH}",
    ))
    env = {}
    env.update(ctx.attr.env)
    env["PATH"] = path
    env["BAZEL_BINDIR"] = "$(BINDIR)"
    env_script = " && ".join([
        cmd
        for name, value in env.items()
        for cmd in ['{}="{}"'.format(name, value), "export {}".format(name)]
    ])

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlHugoSiteInfo(
            themes = themes,
            package_filegroup_info = ctx.attr.site[PackageFilegroupInfo],
            env = ctx.attr.env,
            env_script = env_script,
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Define a hugo site",
    toolchains = [
        "//bzl/rules/hugo:toolchain_type",
    ],
    provides = [AlHugoSiteInfo],
    attrs = {
        "site": attr.label(
            mandatory = True,
            providers = [PackageFilegroupInfo],
            doc = "Hugo site layout",
        ),
        "themes": attr.label(
            mandatory = True,
            allow_single_file = [".tar"],
            doc = "Hugo themes archive (hugo doesn't support symlinks in themes/)",
        ),
        "postcss": attr.label(
            mandatory = True,
            doc = "Postcss target",
            executable = True,
            cfg = "exec",
        ),
        "tools": attr.label_list(
            doc = "Tools that should be available for the build",
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
