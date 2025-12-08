load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_files")
load("//tools/template_files/main/bzl:al_template_files.bzl", "al_template_files")

def al_leetcode_submissions(name, srcs, visibility = None, **kwargs):
    """
    Generate leetcode submission targets

    Args:
        name (str): generated md archive name
        srcs (list[str]): leetcode submission configs
        visibility: visibility
        **kwargs: kwargs for al_template_files
    """
    write_file(
        name = "{}.template".format(name),
        out = "{}.template.md".format(name),
        content = [
            "{{ range .Data }}",
            "{{ range .Data.Submissions }}",
            "---",
            "title: {{ .Timestamp | timestamp_to_date }}",
            "description: |-",
            "  {{ indent .Title \"  \" }}",
            "tags: [generated, leetcode_submission]",
            "categories: []",
            "date: {{ .Timestamp | timestamp_to_date }}",
            "toc_hide: true",
            "---",
            "",
            "## Links",
            "",
            "- https://leetcode.com{{ .Url }}",
            "",
            "## Code",
            "",
            "```{{ .Lang }}",
            "{{ .Code }}",
            "```",
            "{{ end }}",
            "{{ end }}",
        ],
    )
    src_names = []
    for i, src in enumerate(srcs):
        src_name = "{}.template_{}".format(name, i)
        src_names.append(src_name)
        al_template_files(
            name = src_name,
            srcs = ["{}.template".format(name)],
            data = [src],
            outs = ["{}.md".format(src_name)],
        )
    pkg_files(
        name = name,
        srcs = src_names,
        visibility = visibility,
    )
