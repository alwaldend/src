load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    build_dir = ctx.actions.declare_directory("{}-build".format(ctx.label.name))
    script = ctx.actions.declare_file("{}-script.sh".format(ctx.label.name))
    files = [build_dir] + ctx.files.data
    cmd = ["#!/usr/bin/env sh", "set -eux"]

    for key, value in ctx.attr.env.items():
        cmd.append('{}="{}"'.format(key, value))
        cmd.append("export {}".format(key))

    cmd.extend(
        [
            "tar -xf '{}' -C '{}'".format(ctx.file.tree.path, build_dir.path),
            # "ls -la $(rootpath //js/docsy:node_modules/autoprefixer/dir) || true",
            "$(execpath :postcss) $$(realpath {build}/themes/docsy/assets/scss/main.scss) --config $$(realpath {build}/themes/docsy)".format(build = build_dir.path),
            "exit 1",
            "'{}' build --source '{}' --logLevel '{}' \"$${{@}}\"".format(
                hugo.hugo.path,
                build_dir.path,
                ctx.attr.log_level,
            ),
        ],
    )
    script_content = ctx.expand_make_variables(
        "script",
        ctx.expand_location("\n".join(cmd + [""])),
        {},
    )

    ctx.actions.write(
        content = script_content,
        output = script,
        is_executable = True,
    )

    ctx.actions.run(
        inputs = [ctx.file.tree, hugo.hugo] + ctx.files.data,
        executable = script,
        outputs = [build_dir],
        arguments = ctx.attr.arguments,
        tools = [tool.files for tool in ctx.attr.tools] + [tool[DefaultInfo].default_runfiles.files for tool in ctx.attr.tools],
        progress_message = "Building hugo site %{label} to %{output}",
    )

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = ctx.runfiles(files),
        ),
        AlHugoSiteInfo(
            build = build_dir,
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Build a hugo site",
    toolchains = [
        #"@com-alwaldend-src-nodejs-toolchains//:resolved_toolchain",
        "//bzl/toolchain-types:hugo",
    ],
    provides = [AlHugoSiteInfo],
    attrs = {
        "log_level": attr.string(
            doc = "Value for --logLevel",
            default = "info",
        ),
        "arguments": attr.string_list(
            doc = "Build arguments",
            default = [],
        ),
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
