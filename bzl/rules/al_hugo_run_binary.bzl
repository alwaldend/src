load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    info = ctx.attr.site[AlHugoSiteInfo]
    default_info = ctx.attr.site[DefaultInfo]
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    script = ctx.actions.declare_file("{}-script".format(ctx.label.name))
    destination = ctx.actions.declare_directory("{}-destination".format(ctx.label.name))

    script_content = """\
        #!/usr/bin/env sh
        set -eux
        {env_script}
        top="${{PWD}}"
        ln -s "${{top}}/{content}" ./content
        ln -s "${{top}}/{themes}" ./themes
        ln -s "${{top}}/{data}" ./data
        ln -s "${{top}}/{config}" ./config
        ln -s "${{top}}/{layouts}" ./layouts
        chmod -R 700 ./content ./ ./themes
        find content/ -name "README.md" -exec sh -c 'mv "{{}}" "$(dirname "{{}}")/_index.md"' ";"
        "${{top}}/{hugo}" --destination '{destination}' "${{@}}"
    """.format(
        hugo = hugo.hugo.path,
        content = info.content.path,
        themes = info.themes.path,
        config = info.config.path,
        data = info.data.path,
        layouts = info.layouts.path,
        destination = destination.path,
        env_script = info.env_script,
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
        "//bzl/toolchain-types:hugo",
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
