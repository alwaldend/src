load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    default_info = ctx.attr.site[DefaultInfo]
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    destination = ctx.actions.declare_directory("{}.destination".format(ctx.label.name))

    script_content = """\
        #!/usr/bin/env sh
        set -eux
        {env_script}
        exec '{hugo}' -d '{destination}' "${{@}}"
    """.format(
        hugo = hugo.hugo.path,
        destination = destination.path,
        env_script = ctx.expand_make_variables(
            "env_script",
            ctx.expand_location(info.env_script),
            {},
        ),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    ctx.actions.run(
        executable = script,
        inputs = [hugo.hugo] + ctx.files.site,
        outputs = [destination],
        tools = [default_info.default_runfiles.files, default_info.files],
        arguments = ctx.attr.arguments,
    )

    return [
        DefaultInfo(
            files = depset([destination]),
        ),
    ]

al_hugo_run_binary = rule(
    implementation = _impl,
    doc = "Run hugo binary as a build action",
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
