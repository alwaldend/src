load("@bazel_skylib//lib:shell.bzl", "shell")

def _impl(ctx):
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = ctx.files.binary).merge(ctx.attr.binary[DefaultInfo].default_runfiles)
    args = []
    args.extend(ctx.attr.arguments)
    script_content = """\
        #!/usr/bin/env sh
        exec '{binary}' \
            {arguments} \
            "${{@}}"
    """.format(
        binary = ctx.executable.binary.short_path,
        workspace_name = ctx.workspace_name,
        arguments = " ".join([shell.quote(arg) for arg in args]),
    )
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = script_content,
    )
    default_info = DefaultInfo(
        runfiles = runfiles,
        executable = script,
    )
    run_environment_info = RunEnvironmentInfo(
        inherited_environment = ctx.attr.env_inherit,
    )
    return [default_info, run_environment_info]

al_dnscontrol_test = rule(
    implementation = _impl,
    test = True,
    doc = "Dnscontrol test",
    attrs = {
        "binary": attr.label(
            mandatory = True,
            executable = True,
            cfg = "exec",
            doc = "Dnscontrol binary",
        ),
        "env_inherit": attr.string_list(
            doc = "Environment variables to inherit",
        ),
        "arguments": attr.string_list(
            default = [],
            doc = "Dnscontrol arguments",
        ),
    },
)
