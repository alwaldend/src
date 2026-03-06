load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("@rules_pkg//pkg:mappings.bzl", "pkg_files")
load("@rules_template//main/bzl:template_run_binary.bzl", "template_run_binary")

def al_proto_docs(name, src, prefix = None, visibility = None, renames = None):
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
            "{{ range .DataFiles }}",
            "---",
            "title: {{ .BasenameWithoutExt }}",
            "description: Proto docs for {{ .Basename }}",
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
            '        "@{}//{}:{{{{ .BasenameWithoutExt }}}}_java_library",'.format(native.module_name(), native.package_name()),
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
            '        "@{}//{}:{}",'.format(
                native.module_name(),
                native.package_name(),
                native.package_name().rsplit("/", 1)[-1],
            ),
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
    template_run_binary(
        name = "{}.{}".format(name, src),
        srcs = ["{}.template".format(name)],
        data = [src],
        args = ["--extension", ".txt"],
        outs = ["{}.{}.md".format(name, src)],
    )
    pkg_files(
        name = name,
        srcs = ["{}.{}".format(name, src)],
        renames = renames or {
            ":{}.{}".format(name, src): "content/docs/{}/{}/index.md".format(native.package_name(), src),
        },
        visibility = visibility,
    )
