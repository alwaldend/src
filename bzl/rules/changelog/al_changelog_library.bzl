load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")
load("//bzl/rules/template_files:al_template_files.bzl", "al_template_files")

def al_changelog_library(name, visibility, subpackages = []):
    """
    Create changelog target

    Args:
        name: target name prefix
        visibility: visibility
        subpackages: subpackages
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
    write_file(
        name = "{}-template".format(name),
        out = "{}-template.md".format(name),
        content = [
            "---",
            "title: Changelog",
            "tags: [generated, changelog]",
            "toc_hide: true",
            "hide_summary: true",
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
    native.genrule(
        name = "{}-changelog-data".format(name),
        srcs = ["//cfg/misc:git-log"],
        outs = [
            "{}-changelog-data.yaml".format(name),
        ],
        cmd = '''
            $(execpath //tools:git) log \
                "--pretty=format:$$(cat $(execpath //cfg/misc:git-log))" \
                -- \
                '{}' \
                >$(@)
        '''.format(package_name if package_name else "."),
        tools = ["//tools:git"],
    )
    al_template_files(
        name = "{}-changelog".format(name),
        srcs = ["{}-template".format(name)],
        data = ["{}-changelog-data".format(name)],
        outs = ["{}.md".format(name)],
        visibility = visibility,
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
        srcs = ["{}-changelog".format(name)],
        package_dir = package_dir,
        deps = ["{}-children".format(name)],
    )
