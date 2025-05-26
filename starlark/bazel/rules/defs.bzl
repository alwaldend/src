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

def _al_genrule_impl(ctx):
    flagfile = ctx.actions.declare_file(ctx.label.name + "-flagfile")
    cmd = ctx.expand_location(ctx.attr.cmd)
    ctx.actions.write(
        output = flagfile,
        content = json.encode_indent({
            "cmd": cmd,
            "set_args": ctx.attr.set_flags,
            "shell": ctx.attr.shell,
            "var": ctx.var,
        }),
    )
    ctx.actions.run(
        mnemonic = "Shell",
        executable = ctx.executable.worker,
        inputs = [flagfile] + ctx.files.srcs + ctx.files.tools,
        outputs = ctx.outputs.outs,
        arguments = ["run", "--flagfile={}".format(flagfile.path)],
        execution_requirements = {
            "supports-workers": "1",
            "requires-worker-protocol": "proto",
        },
    )

al_genrule = rule(
    implementation = _al_genrule_impl,
    doc = "Shell worker",
    attrs = {
        "srcs": attr.label_list(default = [], doc = "Sources", allow_files = True),
        "tools": attr.label_list(default = [], doc = "tools"),
        "outs": attr.output_list(mandatory = True, doc = "Outputs"),
        "cmd": attr.string(mandatory = True, doc = "Script to execute"),
        "set_flags": attr.string_list(doc = "set flags", default = ["-eux"]),
        "shell": attr.string(doc = "shell to use", default = "/bin/sh"),
        "worker": attr.label(
            default = Label("//golang/bazel-shell-worker"),
            executable = True,
            allow_single_file = True,
            cfg = "exec",
            doc = "Worker binary",
        ),
    },
)
