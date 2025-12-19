load(":al_secrets_backend_info.bzl", "AlSecretsBackendInfo")

def _impl(ctx):
    backend = ctx.actions.declare_file("{}.backend.json".format(ctx.label.name))
    args = ctx.actions.args()
    args.add_all(["build", "backend", "--output", backend.path])
    if not ctx.attr.system:
        args.add("--systemd_creds_user")
    for src in ctx.files.srcs:
        args.add_all(["--systemd_creds_file", src.path])
    args.add_all(["build"])
    ctx.actions.run(
        executable = ctx.executable.secrets_tool,
        outputs = [backend],
        args = [args],
        inputs = ctx.files.srcs,
    )
    return [
        DefaultInfo(
            files = depset([depset]),
        ),
        AlSecretsBackendInfo(
            backend = backend,
        ),
    ]

al_secrets_systemd_creds = rule(
    implementation = _impl,
    provides = [AlSecretsBackendInfo, DefaultInfo],
    doc = "Systemd-creds backend",
    attrs = {
        "srcs": attr.label_list(
            doc = "Sytemd-creds files",
        ),
        "system": attr.bool(
            doc = "Whether it's a system credential or a user one",
        ),
        "secrets_tool": attr.label(
            exec = "cfg",
            executable = True,
            default = "//tools/secrets",
            doc = "Secrets tool",
        ),
    },
)
