load("//tools/release/main/bzl:al_release_deployment_info.bzl", "AlReleaseDeploymentInfo")
load("//tools/release/main/bzl:al_release_info.bzl", "AlReleaseInfo")

_SCRIPT = """\
#!/usr/bin/env sh
set -eu
for root in "" "${{RUNFILES_DIR:-}}/{workspace_name}/" "${{0}}.runfiles/{workspace_name}/"; do
    if [ -x "${{root}}{bin}" ]; then
        exec "${{root}}{bin}" {arguments} "${{@}}"
    fi
done
echo "Could not find the binary"
exit 1
"""

def _impl(ctx):
    runfiles = ctx.runfiles()
    if ctx.attr.oras:
        runfiles = runfiles.merge(ctx.attr.oras[DefaultInfo].default_runfiles)
    runfiles = runfiles.merge(ctx.attr.release_tool[DefaultInfo].default_runfiles)
    script = ctx.actions.declare_file("{}.script.sh".format(ctx.label.name))
    arguments = [ctx.attr.cmd]
    if ctx.attr.oras:
        arguments.extend(["--soras_path", "$${{root}}{}".format(ctx.executable.oras.short_path)])
    if ctx.attr.ssh_deployment:
        deployment = ctx.attr.ssh_deployment[AlReleaseDeploymentInfo].info_file
        runfiles = runfiles.merge(ctx.runfiles(files = [deployment]))
        arguments.extend(["--ssh_deployment", "$${{root}}{}".format(deployment.short_path)])

    symlinks = {}
    for release_target in ctx.attr.srcs:
        release = release_target[AlReleaseInfo]
        dir = "releases/{}/{}".format(release.project, release.release_name)
        symlinks["{}/release.json".format(dir)] = release.manifest
        for path, file in release.files.items():
            symlinks["{}/files/{}".format(dir, path)] = file
        arguments.extend(("--release_dir", "$${{root}}{}".format(dir)))
    runfiles = runfiles.merge(ctx.runfiles(symlinks = symlinks))

    ctx.actions.write(
        output = script,
        content = _SCRIPT.format(
            workspace_name = ctx.workspace_name,
            bin = ctx.executable.release_tool.short_path,
            arguments = " ".join(
                [
                    '"{}"'.format(ctx.expand_make_variables(str(ctx.label), arg, {}))
                    for arg in arguments + ctx.attr.arguments
                ],
            ),
        ),
    )

    return [
        DefaultInfo(
            executable = script,
            runfiles = runfiles,
        ),
    ]

_al_release_binary = rule(
    implementation = _impl,
    doc = "Release binary",
    executable = True,
    attrs = {
        "cmd": attr.string(
            default = "deploy",
            doc = "Cmd",
        ),
        "arguments": attr.string_list(
            doc = "Arguments",
        ),
        "srcs": attr.label_list(
            providers = [AlReleaseInfo],
            doc = "Releases",
        ),
        "ssh_deployment": attr.label(
            providers = [AlReleaseDeploymentInfo],
            doc = "SSH destination for all files in runtime release_dir inputs",
        ),
        "oras": attr.label(
            executable = True,
            cfg = "exec",
            doc = "Oras binary",
        ),
        "release_tool": attr.label(
            executable = True,
            cfg = "exec",
            doc = "Release tool",
            default = "//tools/release/main/go",
        ),
    },
)

def al_release_binary(name, oras = "//third_party/land_oras_oras:oras", **kwargs):
    """Create a release command; oras=None omits the OCI executable.

    Args:
        name: Target name.
        oras: OCI executable, or None for an SSH-only command.
        **kwargs: Remaining release rule attributes.
    """
    _al_release_binary(name = name, oras = oras, **kwargs)
