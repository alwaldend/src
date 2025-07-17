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
            "title: {{ .Data.title_slug }}",
            "description: {{ .Data.title }}",
            "tags: [generated, leetcode]",
            "date: {{ .Data.timestamp | timestamp_to_date }}",
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
    for src in srcs:
        al_template_files(
            name = "{}-{}".format(name, src),
            srcs = ["{}-template".format(name)],
            data = [src],
            outs = ["{}-{}.md".format(name, src)],
        )
    pkg_tar(
        name = name,
        package_dir = native.package_name(),
        srcs = ["{}-{}".format(name, src) for src in srcs],
        visibility = visibility,
    )
