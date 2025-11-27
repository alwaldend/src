def _impl(ctx):
    minisign = ctx.toolchains["//tools/minisign/main/bzl:toolchain_type"]
    exec = ctx.actions.declare_file("{}.minisign".format(ctx.label.name))
    ctx.actions.symlink(
        target_file = minisign.minisign,
        output = exec,
        is_executable = True,
    )
    return [
        DefaultInfo(
            files = minisign.default_info.files,
            runfiles = minisign.default_info.default_runfiles,
            executable = exec,
        ),
    ]

al_minisign_binary = rule(
    implementation = _impl,
    doc = "Minisign binary",
    executable = True,
    toolchains = ["//tools/minisign/main/bzl:toolchain_type"],
)
