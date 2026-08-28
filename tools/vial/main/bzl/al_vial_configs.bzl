load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_docs//docs:defs.bzl", "docs_filegroup")
load("@rules_template//main/bzl:template_run_binary.bzl", "template_run_binary")

def al_vial_configs(name, srcs, visibility = None, **kwargs):
    """
    Generate vial config targets

    Args:
        name (str): generated docs archive name
        srcs (list[str]): vial config
        visibility: visibility
        **kwargs: kwargs for template_files
    """
    write_file(
        name = "{}.template".format(name),
        out = "{}.template.md".format(name),
        content = [
            "{{ range .DataFiles }}",
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
        template_run_binary(
            name = "{}.{}".format(name, src),
            srcs = ["{}.template".format(name)],
            data = [src],
            outs = ["{}.{}.md".format(name, src)],
        )
    docs_filegroup(
        name = name,
        srcs = ["{}.{}".format(name, src) for src in srcs],
        visibility = visibility,
    )
