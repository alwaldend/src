load("@aspect_bazel_lib//lib:write_source_files.bzl", "write_source_file")
load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("//tools/template_files:al_template_files.bzl", "al_template_files")

def al_genquery_write_to_source_file(
        name,
        expression,
        scope,
        var_name,
        out_file):
    """
    Write genquery result to a bzl file

    Example:

    ```starlark
        al_genquery_write_to_source_file(
            name = "al_bzl_libs",
            expression = \"\"\"
                filter(
                    "^//",
                    attr(
                        "srcs",
                        ".{3,}",
                        kind(
                            "bzl_library",
                            deps("//bzl")
                        )
                    )
                )
            \"\"\",
            out_file = "al_bzl_libs.bzl",
            scope = ["//bzl"],
            var_name = "AL_BZL_LIBS",
        )
    ```

    Args:
        name (string): name prefix
        expression (string): genquery expression
        scope (string): genquery scope
        var_name (string): variable name in the generated .bzl file
        out_file (string): output bzl file
    """

    native.genquery(
        name = "{}-query.txt".format(name),
        expression = expression,
        scope = scope,
    )

    write_file(
        name = "{}-template".format(name),
        out = "{}-template.txt.tpl".format(name),
        content = [
            "{} = [".format(var_name),
            "{{ range .Data -}}",
            '{{- range .Lines }}    "{{ . }}",',
            "{{ end }}",
            "{{- end -}}",
            "]",
            "",
        ],
    )

    al_template_files(
        name = "{}-templated".format(name),
        srcs = [":{}-template".format(name)],
        outs = ["{}-gen.bzl".format(name)],
        data = [":{}-query.txt".format(name)],
    )

    write_source_file(
        name = "{}-write".format(name),
        in_file = ":{}-templated".format(name),
        out_file = out_file,
    )
