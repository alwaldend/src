load(":al_release_deployment_info.bzl", "AlReleaseDeploymentInfo")

def _impl(ctx):
    info_file = ctx.actions.declare_file("{}.deployment.json".format(ctx.label.name))
    info = {
        "oci": {
            "repository": ctx.attr.oci_repository,
        },
    }
    ctx.actions.write(
        output = info_file,
        content = json.encode_indent(info, prefix = "", indent = "    "),
    )
    return [
        DefaultInfo(
            files = depset([info_file]),
        ),
        AlReleaseDeploymentInfo(
            info_file = info_file,
            info = info,
        ),
    ]

al_release_deployment = rule(
    implementation = _impl,
    doc = "Deployment info",
    provides = [AlReleaseDeploymentInfo],
    attrs = {
        "oci_repository": attr.string(
            doc = "OCI repository url",
        ),
    },
)
