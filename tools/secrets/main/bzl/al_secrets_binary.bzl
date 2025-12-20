load(":al_secrets_backend_info.bzl", "AlSecretsBackendInfo")
load(":al_secrets_secret_info.bzl", "AlSecretsSecretInfo")

_SCRIPT = """\
#!/usr/bin/env sh
set -eux
for root in "" "${{0}}.runfiles/{workspace_name}/" "${{RUNFILES_DIR:-}}/{workspace_name}/"; do
    if [ -x "${{root}}{bin}" ]; then
        exec "${{root}}{bin}" {arguments} "${{@}}"
    fi
done
echo "Could not find the binary"
exit 1
"""

def _impl(ctx):
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    runfiles = ctx.runfiles(files = ctx.files.secrets_tool)
    runfiles.merge(ctx.attr.secrets_tool[DefaultInfo].default_runfiles)
    arguments = [ctx.attr.cmd]
    for src in ctx.attr.srcs:
        secret = src[AlSecretsSecretInfo].secret
        runfiles = runfiles.merge(ctx.runfiles(files = [secret]))
        arguments.extend(["--secret", "$${{root}}{}".format(secret.short_path)])
    for src in ctx.attr.backends:
        for backend in src[AlSecretsBackendInfo].srcs:
            runfiles = runfiles.merge(ctx.runfiles(files = [backend]))
            arguments.extend(["--backend", "$${{root}}{}".format(backend.short_path)])
    arguments = [
        '"{}"'.format(ctx.expand_make_variables(str(ctx.label), arg, {}))
        for arg in arguments + ctx.attr.arguments
    ]
    ctx.actions.write(
        output = script,
        is_executable = True,
        content = _SCRIPT.format(
            workspace_name = ctx.workspace_name,
            bin = ctx.executable.secrets_tool.short_path,
            arguments = " ".join(arguments),
        ),
    )
    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
    ]

al_secrets_binary = rule(
    implementation = _impl,
    executable = True,
    doc = "Secrets binary",
    provides = [DefaultInfo],
    attrs = {
        "cmd": attr.string(
            default = "export",
            doc = "Cmd to use",
        ),
        "arguments": attr.string_list(
            doc = "Arguments",
        ),
        "srcs": attr.label_list(
            providers = [AlSecretsSecretInfo],
            doc = "Secrets",
        ),
        "backends": attr.label_list(
            providers = [AlSecretsBackendInfo],
            doc = "Backends",
        ),
        "secrets_tool": attr.label(
            executable = True,
            cfg = "exec",
            default = "//tools/secrets/main/go",
            doc = "Secrets tool",
        ),
    },
)
