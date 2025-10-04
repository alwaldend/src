load(":al_hugo_site_info.bzl", "AlHugoSiteInfo")

def _impl(ctx):
    hugo = ctx.toolchains["//bzl/rules/hugo:toolchain_type"]
    themes = ctx.actions.declare_directory("{}-themes".format(ctx.label.name))
    content = ctx.actions.declare_directory("{}-content".format(ctx.label.name))
    data = ctx.actions.declare_directory("{}-data".format(ctx.label.name))
    config = ctx.actions.declare_directory("{}-config".format(ctx.label.name))
    layouts = ctx.actions.declare_directory("{}-layouts".format(ctx.label.name))
    i18n = ctx.actions.declare_directory("{}-i18n".format(ctx.label.name))

    files = [
        hugo.hugo,
        themes,
        config,
        content,
        data,
        layouts,
        i18n,
    ] + ctx.files.postcss + ctx.files.tools
    runfiles = ctx.runfiles(
        files = files,
        transitive_files = depset(ctx.files.tools),
    )
    for tool in ctx.attr.tools:
        runfiles.merge(tool[DefaultInfo].default_runfiles)
    runfiles.merge(ctx.attr.postcss[DefaultInfo].default_runfiles)

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
        [themes, ctx.file.themes],
        [data, ctx.file.hugo_data],
        [layouts, ctx.file.hugo_layouts],
        [i18n, ctx.file.hugo_i18n],
    ]:
        ctx.actions.run_shell(
            outputs = [output],
            inputs = [input],
            command = "tar -xf '{}' -C '{}'".format(input.path, output.path),
        )

    ctx.actions.run_shell(
        outputs = [content],
        inputs = [ctx.file.content],
        command = """
            tar -xf '{input}' -C '{output}'
            find '{output}/' -name "README.md" -exec sh -c 'chmod 700 "$(dirname "{{}}")" && mv "{{}}" "$(dirname "{{}}")/_index.md"' ";"
        """.format(input = ctx.file.content.path, output = content.path),
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
            layouts = layouts,
            i18n = i18n,
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
        "hugo_layouts": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo layouts",
        ),
        "hugo_i18n": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "Hugo i18n",
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
