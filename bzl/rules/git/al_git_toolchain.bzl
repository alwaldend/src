def _impl(ctx):
    if not ctx.attr.local_bin and not ctx.attr.label_bin:
        fail("missing both local_bin and label_bin")
    bin = ctx.attr.local_bin or ctx.executable.label_bin.path
    env = {
        "GIT_BIN": bin,
    }
    return [
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            env = env,
            git_bin = bin,
        ),
    ]

al_git_toolchain = rule(
    doc = "Local git toolchain",
    implementation = _impl,
    attrs = {
        "local_bin": attr.string(doc = "Local binary to use"),
        "label_bin": attr.label(
            executable = True,
            cfg = "exec",
            doc = "Label binary to use",
        ),
    },
)
