load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")
load("//bzl/rules/template_files:al_template_files.bzl", "al_template_files")

def al_git_changelog(name, visibility, subpackages = []):
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
        name = "{}.template".format(name),
        out = "{}.template.md".format(name),
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
        name = "{}.changelog_data".format(name),
        srcs = ["//bzl/rules/git:git_log_config"],
        outs = [
            "{}.changelog_data.yaml".format(name),
        ],
        cmd = '''
            $(execpath //tools:git) log \
                "--pretty=format:$$(cat $(execpath //bzl/rules/git:git_log_config))" \
                -- \
                '{}' \
                >$(@)
        '''.format(package_name if package_name else "."),
        tools = ["//tools:git"],
    )
    al_template_files(
        name = "{}.changelog".format(name),
        srcs = ["{}.template".format(name)],
        data = ["{}.changelog_data".format(name)],
        outs = ["{}.doc.md".format(name)],
    )
    pkg_files(
        name = "{}.changelog_files".format(name),
        srcs = ["{}.changelog".format(name)],
    )
    pkg_filegroup(
        name = name,
        visibility = visibility,
        srcs = ["{}.changelog_files".format(name)] + [
            "{}{}:{}".format(package_name_prefix, dep, name)
            for dep in subpackages
        ],
        prefix = package_dir,
    )
