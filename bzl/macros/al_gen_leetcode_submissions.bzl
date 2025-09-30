load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load(":al_template_files.bzl", "al_template_files")

def al_gen_leetcode_submissions(name, srcs, package_dir, visibility = None, **kwargs):
    """
    Generate leetcode submission targets

    Args:
        name (str): generated md archive name
        srcs (list[str]): leetcode submission configs
        package_dir (str): package_dir to use
        visibility: visibility
        **kwargs: kwargs for al_template_files
    """
    write_file(
        name = "{}-template".format(name),
        out = "{}-template.md".format(name),
        content = [
            "{{ range .Data }}",
            "{{ range .Data.Submissions }}",
            "---",
            "title: {{ .Timestamp | timestamp_to_date }}",
            "description: |-",
            "  {{ indent .Title \"  \" }}",
            "tags: [generated, leetcode]",
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
        src_name = "{}-template-{}".format(name, i)
        src_names.append(src_name)
        al_template_files(
            name = src_name,
            srcs = ["{}-template".format(name)],
            data = [src],
            outs = ["{}.md".format(src_name)],
        )
    pkg_tar(
        name = name,
        package_dir = package_dir,
        srcs = src_names,
        visibility = visibility,
    )
