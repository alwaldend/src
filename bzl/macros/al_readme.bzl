load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("//bzl/rules/md:al_md_data.bzl", "al_md_data")

def al_readme(name, srcs = [":README.md"], subpackages = [], **kwargs):
    """
    Create readme targets

    Targets:
    - ${name}: al_md_data target
    - ${name}-children: pkg_tar containing ${name}-children from subpackages

    Args:
        name (str): target name
        subpackages (list[str]): override list of subpackages
        **kwargs: al_md_data kwargs
    """
    if not subpackages:
        subpackages = native.subpackages(include = ["*"], allow_empty = True)
    package_name = native.package_name()
    if package_name:
        package_dir = package_name.split("/")[-1]
        package_name_prefix = "//{}/".format(package_name)
    else:
        package_dir = "."
        package_name_prefix = "//"
    al_md_data(
        name = name,
        srcs = srcs,
        **kwargs
    )
    pkg_tar(
        name = "{}-children".format(name),
        srcs = [name],
        out = "{}-children.tar".format(name),
        package_dir = package_dir,
        deps = [
            "{}{}:{}-children".format(package_name_prefix, dep, name)
            for dep in subpackages
        ],
        **kwargs
    )
