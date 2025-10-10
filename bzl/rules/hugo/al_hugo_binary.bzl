load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    default_info = ctx.attr.site[DefaultInfo]
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))

    runfiles = ctx.runfiles(
        transitive_files = depset(transitive = [
            default_info.default_runfiles.files,
            default_info.files,
        ]),
    )
    script_content = """\
        #!/usr/bin/env sh
        set -eux
        {env_script}
        if [ ! -f '{site_archive}' ]; then
            ls -la ./
            cd "${{0}}.runfiles/{workspace_name}"
            ls -la ./node_modules
        fi
        tar -xf '{site_archive}'
        if [ -f '{git_archive}' ]; then
            tar -xf '{git_archive}'
        fi
        which postcss
        exec '{hugo}' \
            {arguments} \
            "${{@}}"
    """.format(
        hugo = hugo.hugo.short_path,
        site_archive = info.site_archive.short_path,
        git_archive = info.git_archive.short_path if info.git_archive else None,
        workspace_name = ctx.workspace_name,
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
        "//bzl/rules/hugo:toolchain_type",
    ],
    attrs = {
        "arguments": attr.string_list(
            default = [],
            doc = "Hugo arguments",
        ),
        "site": attr.label(
            mandatory = True,
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
    },
)
