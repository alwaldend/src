load(":release_deployment_info.bzl", "ReleaseDeploymentInfo")

def _impl(ctx):
    info_file = ctx.actions.declare_file("{}.deployment.json".format(ctx.label.name))
    if bool(ctx.attr.oci_repository) == bool(ctx.attr.ssh_host):
        fail("select exactly one of oci_repository or ssh_host")
    if ctx.attr.ssh_host:
        if not ctx.attr.ssh_environment:
            fail("ssh_environment is required for SSH deployment")
        info = {"ssh": {
            "environment": ctx.attr.ssh_environment,
            "host": ctx.attr.ssh_host,
            "user": ctx.attr.ssh_user,
            "root": ctx.attr.ssh_root,
            "project": ctx.attr.ssh_project,
            "namespace": ctx.attr.ssh_namespace,
            "site_archive": ctx.attr.ssh_site_archive,
            "publish_user": ctx.attr.ssh_publish_user,
            "public_url": ctx.attr.ssh_public_url,
        }}
    else:
        info = {"oci": {"repository": ctx.attr.oci_repository}}
    ctx.actions.write(
        output = info_file,
        content = json.encode_indent(info, prefix = "", indent = "    "),
    )
    return [
        DefaultInfo(
            files = depset([info_file]),
        ),
        ReleaseDeploymentInfo(
            info_file = info_file,
            info = info,
        ),
    ]

release_deployment = rule(
    implementation = _impl,
    doc = "Deployment info",
    provides = [ReleaseDeploymentInfo],
    attrs = {
        "oci_repository": attr.string(
            doc = "OCI repository url",
        ),
        "ssh_environment": attr.string(doc = "Explicit environment selector"),
        "ssh_host": attr.string(doc = "SSH hostname or SSH configuration alias"),
        "ssh_user": attr.string(doc = "Administrator login; empty uses SSH configuration"),
        "ssh_root": attr.string(default = "/srv/download", doc = "Provisioned content root"),
        "ssh_project": attr.string(doc = "Public artifact name; empty removes the manifest namespace prefix"),
        "ssh_namespace": attr.string(doc = "Public artifact namespace: projects (default) or third_party"),
        "ssh_site_archive": attr.string(doc = "Optional website archive filename"),
        "ssh_publish_user": attr.string(doc = "Content account selected through sudo -n -u"),
        "ssh_public_url": attr.string(doc = "Public download origin for generated artifact links"),
    },
)
