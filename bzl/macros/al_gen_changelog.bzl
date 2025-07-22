load("@bazel_skylib//rules:write_file.bzl", "write_file")
load(":al_template_files.bzl", "al_template_files")

def al_gen_changelog(name, data, visibility):
    """
    Create changelog target

    Args:
        name: target name prefix
        data: commit data
        visibility: visibility
    """
    write_file(
        name = "{}-template".format(name),
        out = "{}-template.md".format(name),
        content = [
            "---",
            "title: Changelog",
            "tags: [generated, changelog]",
            "---",
            "",
            "{{ range .Data -}}",
            "{{ range .Data -}}",
            "- {{ .subject | html_escape }} ([{{ sliceString .hash 0 6 | html_escape }}](https://github.com/alwaldend/src/commit/{{ .hash | html_escape }}))",
            "{{- if not (eq .subject .message) }}",
            "",
            '{{ .message | html_escape | indent "  " }}',
            "{{ end }}",
            "{{ end }}",
            "{{ end }}",
        ],
    )
    al_template_files(
        name = "{}-changelog".format(name),
        srcs = ["{}-template".format(name)],
        data = data,
        outs = ["{}-changelog.md".format(name)],
        visibility = visibility,
    )
