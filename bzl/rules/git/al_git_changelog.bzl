load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_filegroup", "pkg_files")
load("//bzl/rules/git:al_git_run_binary.bzl", "al_git_run_binary")
load("//bzl/rules/template_files:al_template_files.bzl", "al_template_files")

_GIT_LOG_FORMAT = """
- hash: >2-
    %H
  author_date: >2-
    %aI
  committer_date: >2-
    %cI
  subject: >2-
    %s
  author: |2-
    %w(0,0,4)%an <%ae>%w(0,0,0)
  committer: |2-
    %w(0,0,4)%cn <%ae>%w(0,0,0)
  message: |2-
    %w(0,0,4)%B%w(0,0,0)
  changed_files: []
"""

def al_git_changelog(name, visibility, git_binary = "@git//:git", subpackages = []):
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
    al_git_run_binary(
        name = "{}.changelog_data".format(name),
        srcs = ["//bzl/rules/git:git_log_config"],
        git = git_binary,
        outs = [
            "{}.changelog_data.yaml".format(name),
        ],
        arguments = [
            "log",
            "--pretty=format:{}".format(_GIT_LOG_FORMAT),
            "--output",
            "$(execpath :{}.changelog_data.yaml)".format(name),
            "--",
            package_name if package_name else ".",
        ],
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
