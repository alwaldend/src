load("@aspect_bazel_lib//lib:write_source_files.bzl", "write_source_file")
load("@bazel_skylib//rules:write_file.bzl", "write_file")
load("//tools/template_files/main/bzl:al_template_files.bzl", "al_template_files")

def al_helm_chart_lock(name, lock, lock_out):
    """
    Generate targets for a chart lock

    Args:
        name: name
        lock: lock label
        lock_out: parsed lock filename
    """
    write_file(
        name = "{}.lock_template".format(name),
        out = "{}.lock_template.json".format(name),
        content = [
            "{{- with (index .Data 0) -}}",
            '{{- to_json_indent .Data "" "  " -}}',
            "{{- end -}}",
        ],
    )
    al_template_files(
        name = name,
        srcs = ["{}.lock_template".format(name)],
        data = [lock],
        args = ["--extension", ".yaml"],
        outs = ["{}.json".format(name)],
    )
    write_source_file(
        name = "{}.update".format(name),
        in_file = name,
        out_file = lock_out,
    )
