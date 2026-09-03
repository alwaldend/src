"""Bazel-aware DNSControl configuration packaging."""

load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")

def _requires_impl(ctx):
    requires = []
    require_files = []
    for src in ctx.files.srcs:
        link_name = "{}__{}".format(
            ctx.attr.name.removesuffix("_manifest_raw"),
            src.short_path.replace("/", "_").replace(".", "_"),
        )
        link = ctx.actions.declare_file(link_name)
        ctx.actions.symlink(
            output = link,
            target_file = src,
        )
        requires.append(link.basename)
        require_files.append(link)

    content = "module.exports = [\n{}];\n".format(
        "".join([
            "    require(\"./{}\"),\n".format(path)
            for path in requires
        ]),
    )
    out = ctx.outputs.out
    ctx.actions.write(output = out, content = content)
    return [DefaultInfo(files = depset([out] + require_files))]

_requires = rule(
    implementation = _requires_impl,
    attrs = {
        "out": attr.output(mandatory = True),
        "srcs": attr.label_list(allow_files = True, mandatory = True),
    },
)

def dnscontrol_site(name, config, srcs, data = [], visibility = None):
    """Creates a Bazel-aware DNSControl configuration bundle.

    Args:
        name: Target name.
        config: DNSControl JavaScript entrypoint.
        srcs: Record configurations bundled as `name_manifest`.
        data: Additional files needed by DNSControl.
        visibility: Optional target visibility.
    """
    _requires(
        name = "{}_manifest_raw".format(name),
        srcs = srcs,
        out = "{}.js".format(name),
    )

    pkg_files(
        name = "{}_manifest".format(name),
        srcs = [":{}_manifest_raw".format(name)],
    )

    pkg_files(
        name = "{}_data".format(name),
        srcs = data,
        strip_prefix = "/",
    )

    kwargs = {
        "name": name,
        "srcs": [
            ":{}_manifest".format(name),
            ":{}_data".format(name),
        ],
    }
    if visibility != None:
        kwargs["visibility"] = visibility
    pkg_filegroup(**kwargs)
