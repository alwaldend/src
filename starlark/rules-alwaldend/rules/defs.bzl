"""
Bazel rules
"""

load("@rules_alwaldend//providers:defs.bzl", "al_transitive_sources")

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
    runfiles = ctx.runfiles(files = ctx.files.data + ctx.files.tools)
    transitive_runfiles = []

    for runfiles_attr in (
        ctx.attr.data,
        ctx.attr.tools,
    ):
        for target in runfiles_attr:
            transitive_runfiles.append(target[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge_all(transitive_runfiles)

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
        tools = ctx.files.tools,
        outputs = ctx.outputs.outs,
        arguments = ["run", "--flagfile={}".format(flagfile.path)],
        execution_requirements = {
            "supports-workers": "1",
            "requires-worker-protocol": "proto",
        },
    )
    executable = None
    if ctx.attr._output_executable:
        executable = ctx.outputs.outs[0]
    return [
        DefaultInfo(
            executable = executable,
            runfiles = runfiles,
        ),
    ]

_al_genrule_attrs = {
    "srcs": attr.label_list(default = [], doc = "Sources, will not be added to runfiles", allow_files = True),
    "data": attr.label_list(default = [], doc = "Data, will be added to runfiles", allow_files = True),
    "tools": attr.label_list(default = [], doc = "Tools, will be added to runfiles"),
    "outs": attr.output_list(mandatory = True, doc = "Outputs"),
    "cmd": attr.string(mandatory = True, doc = "Script to execute"),
    "set_flags": attr.string_list(doc = "set flags", default = ["-eux"]),
    "shell": attr.string(doc = "shell to use", default = "/bin/sh"),
    "worker": attr.label(
        default = Label("@com-alwaldend-git-src//go/bazel-shell-worker"),
        executable = True,
        allow_single_file = True,
        cfg = "exec",
        doc = "Worker binary",
    ),
    "_output_executable": attr.bool(default = False),
}

al_genrule_test = rule(
    implementation = _al_genrule_impl,
    doc = "Test using shell worker",
    test = True,
    attrs = _al_genrule_attrs | {"_output_executable": attr.bool(default = True)},
)

al_genrule_executable = rule(
    implementation = _al_genrule_impl,
    doc = "Build executable using shell worker",
    executable = True,
    attrs = _al_genrule_attrs | {"_output_executable": attr.bool(default = True)},
)

al_genrule_regular = rule(
    implementation = _al_genrule_impl,
    doc = "Build shell worker rule",
    attrs = _al_genrule_attrs,
)

def _al_write_script_impl(ctx):
    pass

al_write_script = rule(
    implementation = _al_write_script_impl,
    doc = "Write a script",
    attrs = {
        "shebang": attr.string(default = "#!/usr/bin/env sh", doc = "Sheband to use"),
        "set_flags": attr.string_list(default = ["-eu"], doc = "Flags to pass to set"),
        "content": attr.string(mandatory = True, doc = "Script content"),
    },
)
