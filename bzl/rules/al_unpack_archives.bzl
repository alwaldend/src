def _unpack_archives_impl(ctx):
    directory = ctx.attr.out or (ctx.label.name + "-directory")
    output = ctx.actions.declare_directory(directory)
    cmd = "mkdir '{}'".format(directory)
    for src in ctx.files.srcs:
        cmd += " && tar -xf '{}' -C '{}'".format(src.path, directory)

    ctx.actions.run_shell(
        command = cmd,
        outputs = [output],
        inputs = ctx.files.srcs,
        progress_message = "Unpacking %{label} to %{output}",
    )

    return DefaultInfo(
        files = depset([output]),
    )

al_unpack_archives = rule(
    implementation = _unpack_archives_impl,
    doc = "Unpack several archives using tar into a directory",
    attrs = {
        "srcs": attr.label_list(allow_files = [".tar"], mandatory = True),
        "out": attr.string(default = ""),
    },
)
