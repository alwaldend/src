load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")
load("//tools/template_files:al_template_files.bzl", "al_template_files")

def al_bzl_target_doc(name, visibility, subpackages = []):
    """
    Document bazel targets

    Args:
        name (str): target name
        subpackages (list[str]): list of subpackages
        **kwargs: al_md_data kwargs
    """
    package_name = native.package_name()
    if not subpackages:
        exclude = []
        if not package_name:
            exclude.append("bazel-src")
        subpackages = native.subpackages(
            include = ["**"],
            exclude = exclude,
            allow_empty = True,
        )
    srcs = native.existing_rules()
    if package_name:
        package_dir = package_name.split("/")[-1]
        package_name_prefix = "//{}/".format(package_name)
    else:
        package_dir = "."
        package_name_prefix = "//"

    docs = []
    if srcs:
        native.genquery(
            name = "{}.query.ndjson".format(name),
            expression = " union ".join([
                "//{}:{}".format(native.package_name(), src)
                for src in srcs.keys()
                if "node_modules" not in src
            ]),
            scope = srcs.keys(),
            opts = ["--output", "streamed_jsonproto"],
        )
        write_file(
            name = "{}.template".format(name),
            out = "{}.template.md".format(name),
            content = [
                "{{ range .Data }}",
                "---",
                "title: Bazel targets",
                "tags: [generated, bzl_targets]",
                "toc_hide: true",
                "hide_summary: true",
                "---",
                "",
                # print({key: value.items() for key, value in srcs.items()}),
                "{{` {{< al_bzl_target_doc.inline >}} `}}",
                "<table>",
                "<thead><tr><th>Name</th><th>Info</th></thead>",
                "<tbody>",
                "{{ range .Data -}}",
                "{{ if .rule -}}",
                "{{ $data := .rule -}}",
                "{{ range .rule.attribute -}}",
                '{{ if eq .name "visibility" -}}',
                '{{ $_ := set_map_key $data "visibility" .stringListValue -}}',
                "{{ end -}}",
                "{{ end -}}",
                '{{ $_ := unset_map_key .rule "attribute" "ruleInput" -}}',
                "{{ end -}}",
                "<tr>",
                '<td><b>{{ last (split .rule.name ":") }}</b></td>',
                "<td>{{ to_html_table .rule }}</td>",
                "</tr>",
                "{{ end }}",
                "</tbody>",
                "</table>",
                "{{` {{< /al_bzl_target_doc.inline >}} `}}",
                "",
                "{{ end }}",
            ],
        )
        docs.append("{}.doc".format(name))
        al_template_files(
            name = "{}.doc".format(name),
            srcs = ["{}.template".format(name)],
            data = ["{}.query.ndjson".format(name)],
            outs = ["{}.doc.md".format(name)],
        )

    pkg_files(
        name = "{}.docs".format(name),
        srcs = docs,
        prefix = native.package_name(),
    )
    pkg_filegroup(
        name = name,
        srcs = [":{}.docs".format(name)] + [
            "{}{}:{}".format(package_name_prefix, dep, name)
            for dep in subpackages
        ],
        visibility = visibility,
    )
