load(":al_secrets_backend_info.bzl", "AlSecretsBackendInfo")

def _impl(ctx):
    return [
        DefaultInfo(
            files = depset(ctx.files.srcs),
        ),
        AlSecretsBackendInfo(
            srcs = ctx.files.srcs,
        ),
    ]

al_secrets_backend = rule(
    implementation = _impl,
    provides = [AlSecretsBackendInfo, DefaultInfo],
    doc = "Systemd-creds backend",
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Backend files",
        ),
    },
)
