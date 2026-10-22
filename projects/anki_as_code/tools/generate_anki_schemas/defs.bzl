"""Generate Anki schema registrations with execution-platform descriptors."""

def _anki_schema_imports_impl(ctx):
    descriptors = ctx.file.descriptors
    arguments = ctx.actions.args()
    arguments.add("--descriptors", descriptors.path + "/anki_descriptors.bin")
    arguments.add("--output", ctx.outputs.out.path)
    ctx.actions.run(
        executable = ctx.executable._generator,
        arguments = [arguments],
        inputs = [descriptors],
        outputs = [ctx.outputs.out],
        mnemonic = "AnkiSchemaImports",
    )
    return [DefaultInfo(files = depset([ctx.outputs.out]))]

anki_schema_imports = rule(
    implementation = _anki_schema_imports_impl,
    attrs = {
        "descriptors": attr.label(
            allow_single_file = True,
            cfg = "exec",
            mandatory = True,
            doc = "Upstream build-script output directory containing anki_descriptors.bin.",
        ),
        "out": attr.output(mandatory = True),
        "_generator": attr.label(
            default = Label("//projects/anki_as_code/tools/generate_anki_schemas"),
            executable = True,
            cfg = "exec",
        ),
    },
    doc = "Generate Go registration imports from host-built Anki descriptors.",
)
