load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("@stardoc//stardoc:stardoc.bzl", "stardoc")

def al_bzl_library_map(name, visibility, libs = {}, deps = [], **kwargs):
    """
    Create bzl_library targets from a map

    Args:
        name (str): combined bzl_library target name
        libs (dict[str, list[str]]): bzl_library names
        deps (list[str]): other al_bzl_library_map targets
        visibiliy: visibiliy
        **kwargs: bzl_library kwargs
    """
    bzl_library(
        name = name,
        deps = libs.keys() + deps,
        visibility = visibility,
        **kwargs
    )
    pkg_tar(
        name = "{}-stardoc".format(name),
        package_dir = native.package_name().split("/")[-1],
        deps = ["{}-stardoc-src".format(name), "{}-stardoc-deps".format(name)],
        visibility = visibility,
    )
    stardoc_readme_srcs = []
    if libs:
        stardoc_readme_srcs = ["{}-stardoc-readme".format(name)]
        write_file(
            name = "{}-stardoc-readme".format(name),
            out = "{}-stardoc-readme.md".format(name),
            content = [
                "---",
                "title: Stardoc",
                "description: Stardoc documentation",
                "tags: [generated, stardoc]",
                "---",
            ],
        )
    pkg_tar(
        name = "{}-stardoc-src".format(name),
        srcs = ["{}-stardoc".format(lib) for lib in libs.keys()] + stardoc_readme_srcs,
        package_dir = "stardoc",
        remap_paths = {"/{}-stardoc-readme.md".format(name): "_index.md"},
    )
    pkg_tar(
        name = "{}-stardoc-deps".format(name),
        deps = ["{}-stardoc".format(dep) for dep in deps],
    )
    for lib_name, lib_deps in libs.items():
        bzl_library(
            name = lib_name,
            srcs = ["{}.bzl".format(lib_name)],
            deps = lib_deps,
            visibility = visibility,
            **kwargs
        )
        stardoc(
            name = "{}-stardoc-raw".format(lib_name),
            out = "{}-stardoc-raw.md".format(lib_name),
            input = "{}.bzl".format(lib_name),
            deps = [lib_name],
        )
        native.genrule(
            name = "{}-stardoc".format(lib_name),
            outs = ["{}-stardoc.md".format(lib_name)],
            srcs = ["{}-stardoc-raw".format(lib_name)],
            cmd = """
                {{
                    echo "---"
                    echo "title: {title}"
                    echo "tags: [{tags}]"
                    echo "---"
                    cat $(<)
                }} >$(@)
            """.format(title = lib_name, tags = "generated"),
        )
