load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("//bzl/rules/template_files:al_template_files.bzl", "al_template_files")

def al_vial_configs(name, srcs, visibility, **kwargs):
    """
    Generate vial config targets

    Args:
        name (str): generated docs archive name
        srcs (list[str]): vial config
        visibility: visibility
        **kwargs: kwargs for al_template_files
    """
    write_file(
        name = "{}-template".format(name),
        out = "{}-template.md".format(name),
        content = [
            "{{ range .Data }}",
            "---",
            "title: {{ .Basename }}",
            "description: Vial config {{ .Data.Basename }}",
            "tags: [generated, vial]",
            "---",
            "",
            "## Config",
            "```json",
            "{{ range .Lines -}}",
            "{{ . }}",
            "{{ end -}}",
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
