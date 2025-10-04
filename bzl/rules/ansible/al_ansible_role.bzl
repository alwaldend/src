load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def al_ansible_role(name, srcs, visibility):
    """
    Create targets for an ansible role

    Args:
        name (str): name prefix
        srcs (list[str]): role sources
        visibility (list[str]): resulting archive visibility
    """
    role_name = native.package_name().rsplit("/", 1)[-1]
    if "defaults/main.yaml" in srcs:
        defaults = "defaults/main.yaml"
    elif "defaults/main.yml" in srcs:
        defaults = "defaults/main.yml"
    else:
        fail("role {} is missing defaults: {}".format(native.package_name(), srcs))
    native.genrule(
        name = "{}-defaults".format(name),
        outs = ["{}-defaults.md".format(name)],
        srcs = [defaults],
        cmd = """
            {{
                echo "---"
                echo "title: Defaults"
                echo "description: Defaults for {role_name}"
                echo "tags: [generated]"
                echo "---"
                echo ""
                echo '```yaml'
                cat $(<)
                echo '```'
            }} >$(@)
        """.format(role_name = role_name),
    )
    pkg_tar(
        name = "{}-docs".format(name),
        srcs = [":{}-defaults".format(name)],
        package_dir = role_name,
        remap_paths = {"/{}-defaults.md".format(name): "defaults/_index.md"},
        visibility = visibility,
    )
    pkg_tar(
        name = name,
        srcs = srcs,
        package_dir = role_name,
        strip_prefix = ".",
        visibility = visibility,
    )
