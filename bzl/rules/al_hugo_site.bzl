load("//bzl/providers:al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/toolchain-types:hugo"]
    themes = ctx.actions.declare_directory("{}-themes".format(ctx.label.name))
    content = ctx.actions.declare_directory("{}-content".format(ctx.label.name))
    data = ctx.actions.declare_directory("{}-data".format(ctx.label.name))
    config = ctx.actions.declare_directory("{}-config".format(ctx.label.name))
    files = [hugo.hugo, themes, config, content, data]
    runfiles = ctx.runfiles(
        files = files,
        transitive_files = depset(ctx.files.tools),
    )
    for tool in ctx.attr.tools:
        runfiles.merge(tool[DefaultInfo].default_runfiles)
    env_script = " && ".join([
        cmd
        for name, value in ctx.attr.env.items()
        for cmd in ['{}="{}"'.format(name, value), "export {}".format(name)]
    ])
    env_script = ctx.expand_make_variables(
        "env_script",
        ctx.expand_location(env_script),
        {},
    )

    ctx.actions.run_shell(
        outputs = [config],
        inputs = ctx.files.configs,
        command = """
            mkdir '{config_dir}/_default' && \
                cp {configs} '{config_dir}/_default/'
        """.format(
            config_dir = config.path,
            configs = " ".join([file.path for file in ctx.files.configs]),
        ),
    )

    for output, input in [
        [content, ctx.file.content],
        [themes, ctx.file.themes],
        [data, ctx.file.hugo_data],
    ]:
        ctx.actions.run_shell(
            outputs = [output],
            inputs = [input],
            command = "tar -xf '{}' -C '{}'".format(input.path, output.path),
        )

    return [
        DefaultInfo(
            files = depset(files),
            runfiles = runfiles,
        ),
        AlHugoSiteInfo(
            content = content,
            themes = themes,
            data = data,
            config = config,
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
        "themes": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo themes",
        ),
        "content": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo content",
        ),
        "hugo_data": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo data",
        ),
        "configs": attr.label_list(
            mandatory = True,
            allow_files = True,
            doc = "Hugo configs",
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
