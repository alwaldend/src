"""Build-time Cobra registration from the pinned upstream Anki schema."""

load("@protobuf//bazel/common:proto_info.bzl", "ProtoInfo")

def _generate_impl(ctx):
    descriptor = ctx.attr.schema[ProtoInfo].direct_descriptor_set
    output = ctx.actions.declare_file(ctx.label.name + ".go")
    args = ctx.actions.args()
    args.add("--descriptor", descriptor)
    args.add("--output", output)
    ctx.actions.run(
        executable = ctx.executable._generator,
        arguments = [args],
        inputs = [descriptor],
        outputs = [output],
        mnemonic = "AnkiCobraCommands",
    )
    return [DefaultInfo(files = depset([output]))]

anki_cli_source = rule(
    implementation = _generate_impl,
    attrs = {
        "schema": attr.label(
            mandatory = True,
            providers = [ProtoInfo],
            doc = "Upstream proto library containing the complete Anki descriptor set.",
        ),
        "_generator": attr.label(
            default = "//projects/anki_grpc/cmd/generate_anki_cli",
            executable = True,
            cfg = "exec",
        ),
    },
    doc = "Generates Go service/method commands from an upstream proto library.",
)
