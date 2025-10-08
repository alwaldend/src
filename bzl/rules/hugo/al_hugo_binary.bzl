load("@bazel_skylib//lib:shell.bzl", "shell")
load("@rules_pkg//pkg:providers.bzl", "PackageFilegroupInfo", "PackageFilesInfo")
load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    default_info = ctx.attr.site[DefaultInfo]
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles_symlinks = {"themes": info.themes}

    for files, _ in info.package_filegroup_info.pkg_files:
        runfiles_symlinks.update(files.dest_src_map)

    runfiles = ctx.runfiles(
        files = [],
        transitive_files = hugo.default_info.default_runfiles.files,
        symlinks = runfiles_symlinks,
    )
    runfiles.merge(default_info.default_runfiles)

    script_content = """\
        #!/usr/bin/env sh
        set -eux
        {env_script}
        exec '{hugo}' \
            {arguments} \
            "${{@}}"
    """.format(
        hugo = hugo.hugo.short_path,
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
