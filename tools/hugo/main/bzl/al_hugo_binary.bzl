load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//tools/hugo/main/bzl:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))

    runfiles = ctx.runfiles(files = [hugo.env_file])
    if ctx.attr.site:
        site_info = ctx.attr.site[AlHugoSiteInfo]
        site_archive = site_info.site_archive.short_path
        env_script = site_info.env_script
        runfiles = runfiles.merge(ctx.attr.site[DefaultInfo].default_runfiles)
    else:
        site_archive = ""
        env_script = ""

    script_content = """\
        #!/usr/bin/env sh
        set -eu
        {env_script}
        if [ ! -f '{site_archive}' ]; then
            cd "${{0}}.runfiles/{workspace_name}"
        fi
        if [ -f '{site_archive}' ]; then
            tar -xf '{site_archive}'
        fi
        mkdir -p static
        mv '{env_file}' static/hugo_env.txt
        exec '{hugo}' \
            {arguments} \
            "${{@}}"
    """.format(
        hugo = hugo.hugo.short_path,
        env_file = hugo.env_file.short_path,
        site_archive = site_archive,
        workspace_name = ctx.workspace_name,
        env_script = ctx.expand_make_variables(
            "env_script",
            ctx.expand_location(env_script),
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
        "//tools/hugo/main/bzl:toolchain_type",
    ],
    attrs = {
        "arguments": attr.string_list(
            default = [],
            doc = "Hugo arguments",
        ),
        "site": attr.label(
            providers = [AlHugoSiteInfo],
            doc = "Hugo site",
        ),
    },
)
