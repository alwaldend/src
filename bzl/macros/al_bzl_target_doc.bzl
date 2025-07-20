load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load(":al_template_files.bzl", "al_template_files")

def al_bzl_target_doc(name, visibility, srcs = [], subpackages = []):
    """
    Document bazel targets

    Args:
        name (str): target name
        srcs (dict[str, str]): targets to document
        subpackages (list[str]): list of subpackages
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

    docs = []
    if srcs:
        native.genquery(
            name = "{}-query.ndjson".format(name),
            expression = " union ".join(["{}:{}".format(native.package_name(), src) for src in srcs.keys()]),
            scope = srcs,
            opts = ["--output", "streamed_jsonproto"],
        )
        write_file(
            name = "{}-template".format(name),
            out = "{}-template.md".format(name),
            content = [
                "{{ range .Data }}",
                "---",
                "title: Bazel targets",
                "tags: [generated, bzl]",
                "---",
                "",
                "{{ range .Data }}",
                "## {{ .sourceFile.name }}",
                "",
                "{{ end }}",
                "{{ end }}",
            ],
        )
        docs.append("{}-doc".format(name))
        al_template_files(
            name = "{}-doc".format(name),
            srcs = ["{}-template".format(name)],
            data = ["{}-query.ndjson".format(name)],
            outs = ["{}-doc.md".format(name)],
        )

    pkg_tar(
        name = "{}-children".format(name),
        deps = [
            "{}{}:{}".format(package_name_prefix, dep, name)
            for dep in subpackages
        ],
    )
    pkg_tar(
        name = name,
        visibility = visibility,
        srcs = docs,
        package_dir = package_dir,
        deps = ["{}-children".format(name)],
    )
