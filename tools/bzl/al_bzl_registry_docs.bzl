load("@bazel_skylib//lib:paths.bzl", "paths")
load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")
load("//tools/template_files:al_template_files.bzl", "al_template_files")

def al_bzl_registry_docs(name, srcs, visibility = None):
    """
    Create bazel registry docs

    Args:
        name: name
        srcs: registry sources
        visibility: visibility
    """
    docs = []

    write_file(
        name = "{}.source_template".format(name),
        out = "{}.source_template.md".format(name),
        content = [
            "{{ range .Data }}",
            "---",
            "title: {{ .Dirname | basename }}",
            "description: Version {{ .Dirname | basename }} of module {{ .Dirname | dirname | basename }}",
            "tags: [generated, bzl_registry_module_source]",
            "---",
            "",
            "## Usage",
            "",
            "```bzl",
            'bazel_dep(name = "{{ .Dirname | dirname | basename }}", version = "{{ .Dirname | basename }}")',
            "```",
            "",
            "## Source",
            "",
            "{{ to_html_table .Data }}",
            "",
            "{{ end }}",
        ],
    )
    write_file(
        name = "{}.metadata_template".format(name),
        out = "{}.metadata_template.md".format(name),
        content = [
            "{{ range .Data }}",
            "---",
            "title: {{ .Dirname | basename }}",
            "description: Module {{ .Dirname | basename }}",
            "tags: [generated, bzl_registry_module]",
            "---",
            "",
            "## Usage",
            "",
            "```bzl",
            'bazel_dep(name = "{{ .Dirname | basename }}", version = "{{ .Data.versions | first }}")',
            "```",
            "",
            "## Metadata",
            "",
            "{{ to_html_table .Data }}",
            "",
            "{{ end }}",
        ],
    )

    for i, src in enumerate(srcs):
        basename = paths.basename(src)
        if basename == "source.json":
            template = "{}.source_template".format(name)
        elif basename == "metadata.json":
            template = "{}.metadata_template".format(name)
        else:
            continue
        cur_name = "{}.doc.{}".format(name, i)
        al_template_files(
            name = "{}.gen_md".format(cur_name),
            srcs = [template],
            data = [src],
            outs = ["{}.md".format(cur_name)],
        )
        pkg_files(
            name = cur_name,
            srcs = ["{}.gen_md".format(cur_name)],
            renames = {
                ":{}.gen_md".format(cur_name): "{}/{}/_index.md".format(
                    native.package_name(),
                    paths.dirname(src),
                ),
            },
        )
        docs.append(cur_name)
    pkg_filegroup(
        name = name,
        srcs = docs,
        visibility = visibility,
    )
