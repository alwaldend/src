def _impl(ctx):
    env = {
        "GIT_PATH": ctx.attr.git_path,
        "GIT_DIR": ctx.attr.git_dir,
        "GIT_ROOT": ctx.attr.git_root,
    }
    default_info = DefaultInfo(
        files = depset(ctx.files.invalidation),
    )
    return [
        default_info,
        platform_common.TemplateVariableInfo(env),
        platform_common.ToolchainInfo(
            default_info = default_info,
            env = env,
            git_path = ctx.attr.git_path,
            git_dir = ctx.attr.git_dir,
            git_root = ctx.attr.git_root,
            git_invalidation = ctx.files.invalidation,
        ),
    ]

git_toolchain = rule(
    doc = "Local git toolchain",
    implementation = _impl,
    attrs = {
        "git_path": attr.string(
            mandatory = True,
            doc = "Git binary path",
        ),
        "git_dir": attr.string(
            mandatory = True,
            doc = "Git directory path",
        ),
        "git_root": attr.string(
            mandatory = True,
            doc = "Git workspace path",
        ),
        "invalidation": attr.label_list(
            doc = "Files that should invalidate git actions",
        ),
    },
)

def _resolved_impl(ctx):
    toolchain = ctx.toolchains[Label("//rules:git_toolchain_type")]
    return [toolchain.default_info, toolchain, platform_common.TemplateVariableInfo(toolchain.env)]

git_resolved_toolchain = rule(
    implementation = _resolved_impl,
    toolchains = [Label("//rules:git_toolchain_type")],
    doc = "Expose the configured Git executable and checkout template variables.",
)
