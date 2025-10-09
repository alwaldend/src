load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load("//bzl/rules/git:al_git_library_info.bzl", "AlGitLibraryInfo")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    files = ctx.files.site
    transitive_files = []
    git_archive = None
    if ctx.attr.git:
        git_archive = ctx.attr.git[AlGitLibraryInfo].git_archive
        transitive_files.append(depset([git_archive]))
    for tool in ctx.attr.tools:
        transitive_files.append(tool[DefaultInfo].default_runfiles.files)
        transitive_files.append(tool[DefaultInfo].files)
    transitive_files.append(ctx.attr.postcss[DefaultInfo].default_runfiles.files)
    transitive_files.append(hugo.default_info.default_runfiles.files)
    runfiles = ctx.runfiles(
        files = files,
        transitive_files = depset(transitive = transitive_files),
    )

    path = ":".join((
        "$${{PWD}}/$$(dirname '{}')".format(ctx.executable.postcss.short_path),
        "$${{PWD}}/$$(dirname '{}')".format(ctx.executable.postcss.path),
        "$${{PWD}}/$${{0}}.runfiles/{}/$$(dirname '{}')".format(ctx.workspace_name, ctx.executable.postcss.short_path),
        ctx.attr.env.get("PATH", ""),
        "$${PATH}",
    ))
    env = {}
    env.update(ctx.attr.env)
    env["PATH"] = path
    env["BAZEL_BINDIR"] = "$${PWD}/$(BINDIR)"
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
            site_archive = ctx.file.site,
            env = ctx.attr.env,
            env_script = env_script,
            git_archive = git_archive,
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
            allow_single_file = [".tar"],
            doc = "Hugo site archive",
        ),
        "git": attr.label(
            providers = [AlGitLibraryInfo],
            doc = "Git info",
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
