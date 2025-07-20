load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load(":al_template_files.bzl", "al_template_files")

def al_gen_leetcode_submissions(name, srcs, visibility, **kwargs):
    """
    Generate leetcode submission targets

    Args:
        name (str): generated md archive name
        srcs (list[str]): leetcode submission json files
        visibility: visibility
        **kwargs: kwargs for al_template_files
    """
    write_file(
        name = "{}-template".format(name),
        out = "{}-template.md".format(name),
        content = [
            "{{ range .Data }}",
            "---",
            "title: {{ .Data.timestamp | timestamp_to_date }}",
            "description: {{ .Data.title }}",
            "tags: [generated, leetcode]",
            "date: {{ .Data.timestamp | timestamp_to_date }}",
            "toc_hide: true",
            "---",
            "",
            "## Links",
            "",
            "- https://leetcode.com{{ .Data.url }}",
            "",
            "## Code",
            "```{{ .Data.lang }}",
            "{{ .Data.code }}",
            "```",
            "{{ end }}",
        ],
    )
    src_names = []
    for src in srcs:
        src_name = "{}-{}".format(name, src.replace(".json", ""))
        src_names.append(src_name)
        al_template_files(
            name = src_name,
            srcs = ["{}-template".format(name)],
            data = [src],
            outs = ["{}.md".format(src_name)],
        )
    pkg_tar(
        name = name,
        package_dir = native.package_name(),
        srcs = src_names,
        visibility = visibility,
    )
