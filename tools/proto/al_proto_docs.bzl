load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_files")
load("//tools/template_files:al_template_files.bzl", "al_template_files")

def al_proto_docs(name, src, visibility = None):
    """
    Generate protobuf documentation

    Args:
        name: name
        src: protobuf source file
        visibility: visibility
    """
    write_file(
        name = "{}.template".format(name),
        out = "{}.md".format(name),
        content = [
            "{{ range .Data }}",
            "---",
            "title: {{ .BasenameWithoutExt }}",
            "description: Proto docs for {{ .BasenameWithoutExt }}",
            "tags: [generated, proto_doc]",
            "---",
            "",
            "## Java",
            "",
            "```bzl",
            'load("@rules_java//java:defs.bzl", "java_library")',
            "",
            "java_library(",
            '    name = "name",',
            "    deps = [",
            '        "//{}:{{{{ .BasenameWithoutExt }}}}_java_library",'.format(native.package_name()),
            "    ],",
            ")",
            "```",
            "",
            "## Golang",
            "",
            "```bzl",
            'load("@rules_go//go:def.bzl", "go_library")',
            "",
            "go_library(",
            '    name = "name",',
            "    deps = [",
            '        "//{}:{}",'.format(native.package_name(), native.package_name().rsplit("/", 1)[-1]),
            "    ],",
            ")",
            "```",
            "",
            "## Proto",
            "```proto",
            "{{ range .Lines -}}",
            "{{ . }}",
            "{{ end -}}",
            "```",
            "{{ end }}",
        ],
    )
    al_template_files(
        name = "{}.{}".format(name, src),
        srcs = ["{}.template".format(name)],
        data = [src],
        args = ["--extension", ".txt"],
        outs = ["{}.{}.md".format(name, src)],
    )
    pkg_files(
        name = name,
        prefix = native.package_name(),
        srcs = ["{}.{}".format(name, src)],
        visibility = visibility,
    )
