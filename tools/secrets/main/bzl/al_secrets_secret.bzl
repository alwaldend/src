load(":al_secrets_backend_info.bzl", "AlSecretsBackendInfo")
load(":al_secrets_secret_info.bzl", "AlSecretsSecretInfo")

def _impl(ctx):
    secret = ctx.actions.declare_file("{}.secrets.json".format(ctx.label.name))
    inputs = []
    args = ctx.actions.args()
    args.add_all(["build", "secret", "--output", secret.path])
    for key, value in ctx.attr.env.items():
        args.add_all(["--env", "=".join([key, value])])
    for backend in ctx.attr.backends:
        for backend_file in backend[AlSecretsBackendInfo].srcs:
            args.add_all(["--backend", backend_file.path])
        inputs.append(backend_file)
    for dep in ctx.attr.deps:
        dep_secret = dep[AlSecretsSecretInfo].secret
        args.add_all(["--secret", dep_secret.path])
        inputs.append(dep_secret)

    ctx.actions.run(
        executable = ctx.executable.secrets_tool,
        arguments = [args],
        outputs = [secret],
        inputs = inputs,
    )

    return [
        DefaultInfo(
            files = depset(direct = [secret]),
        ),
        AlSecretsSecretInfo(
            secret = secret,
        ),
    ]

al_secrets_secret = rule(
    implementation = _impl,
    doc = "Secret info",
    provides = [DefaultInfo, AlSecretsSecretInfo],
    attrs = {
        "deps": attr.label_list(
            providers = [AlSecretsSecretInfo],
            doc = "Secrets",
        ),
        "env": attr.string_dict(
            doc = "Env dict, keys are values, values are templates",
        ),
        "backends": attr.label_list(
            providers = [AlSecretsBackendInfo],
            doc = "Backends",
        ),
        "secrets_tool": attr.label(
            cfg = "exec",
            executable = True,
            default = "//tools/secrets",
            doc = "Secrets tool",
        ),
    },
)
