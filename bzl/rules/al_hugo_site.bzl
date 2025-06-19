def _impl(ctx):
    hugo_args = ctx.actions.args()
    destination_dir = ctx.actions.declare_directory("{}-destination".format(ctx.label.name))
    src_dir = ctx.actions.declare_directory("{}-src".format(ctx.label.name))

    hugo_args.add_all(["build"])
    hugo_args.add_all(["--destination", destination_dir.path])
    hugo_args.add_all(["--themesDir", "{}/{}".format(src_dir.path, "themes")])
    hugo_args.add_all(["--contentDir", "{}/{}".format(src_dir.path, "content")])
    hugo_args.add_all(["--config", ",".join([config.path for config in ctx.files.configs])])
    hugo_args.add_all(["--destination", destination_dir.path])
    hugo_args.add_all(ctx.attr.arguments)

    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [src_dir],
        command = "tar -xf '{}' -C '{}'".format(ctx.file.src.path, src_dir.path),
    )

    ctx.actions.run(
        executable = ctx.executable.hugo,
        inputs = [src_dir] + ctx.files.configs,
        outputs = [destination_dir],
        arguments = [hugo_args],
    )

    return [
        DefaultInfo(
            files = depset([destination_dir]),
        ),
    ]

al_hugo_site = rule(
    implementation = _impl,
    doc = "Build a hugo site",
    toolchains = [
        # "@com-alwaldend-src-nodejs-toolchains//:resolved_toolchain",
    ],
    attrs = {
        "configs": attr.label_list(
            allow_files = True,
            mandatory = True,
            doc = "Hugo configs",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Additional CLI arguments for hugo",
        ),
        "src": attr.label(
            allow_single_file = [".tar"],
            mandatory = True,
            doc = "File tree for hugo",
        ),
        "hugo": attr.label(
            default = "//tools:hugo",
            executable = True,
            cfg = "exec",
            doc = "Hugo binary",
        ),
    },
)
