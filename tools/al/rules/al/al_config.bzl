AlConfigInfo = provider(
    doc = "Config info",
    fields = {
        "config": "Merged config file",
    },
)

def _impl(ctx):
    out = ctx.actions.declare_file("{}_/al.json".format(ctx.label.name))
    toolchain = ctx.toolchains["//tools/al/rules/al:toolchain_type"]
    args = ctx.actions.args()
    args.add_all(["config", "dump", "--out", out.path])
    for dep in ctx.attr.deps:
        args.add_all(["--config", dep[AlConfigInfo].config.path])
    for src in ctx.files.srcs:
        args.add_all(["--config", src.path])
    ctx.actions.run(
        executable = toolchain.executable,
        arguments = [args],
        inputs = ctx.files.srcs + ctx.files.deps,
        outputs = [out],
    )
    return [
        DefaultInfo(
            files = depset([out]),
        ),
        AlConfigInfo(
            config = out,
        ),
    ]

al_config = rule(
    doc = "Dump al config",
    toolchains = ["//tools/al/rules/al:toolchain_type"],
    provides = [DefaultInfo, AlConfigInfo],
    implementation = _impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Al configs (source files)",
        ),
        "deps": attr.label_list(
            providers = [AlConfigInfo],
            doc = "Al configs (targets)",
        ),
    },
)
