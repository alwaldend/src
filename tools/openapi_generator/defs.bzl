"""Offline OpenAPI Generator source generation."""

def _generate_impl(ctx):
    output = ctx.actions.declare_directory(ctx.label.name + ".srcs")
    log = ctx.actions.declare_file(ctx.label.name + ".log")
    args = ctx.actions.args()
    args.add("generate")
    args.add("--input-spec", ctx.file.schema.path)
    args.add("--config", ctx.file.config.path)
    args.add("--generator-name", ctx.attr.generator)
    args.add("--output", output.path)
    args.add_all(ctx.attr.args)
    ctx.actions.run(
        executable = ctx.executable._generator,
        arguments = [args],
        inputs = [ctx.file.schema, ctx.file.config, ctx.file._logging],
        outputs = [output, log],
        tools = [ctx.attr._generator[DefaultInfo].files_to_run],
        env = {"JAVA_TOOL_OPTIONS": "-Duser.timezone=UTC -Dfile.encoding=UTF-8 -Dlogback.configurationFile=" + ctx.file._logging.path + " -Dopenapi.generator.log=" + log.path},
        mnemonic = "OpenApiGenerate",
        progress_message = "Generating OpenAPI client %{label}",
    )
    return [
        DefaultInfo(files = depset([output])),
        OutputGroupInfo(logs = depset([log])),
    ]

openapi_generate = rule(
    implementation = _generate_impl,
    attrs = {
        "schema": attr.label(allow_single_file = [".json", ".yaml", ".yml"], mandatory = True),
        "config": attr.label(allow_single_file = [".json", ".yaml", ".yml"], mandatory = True),
        "generator": attr.string(default = "go"),
        "args": attr.string_list(),
        "_generator": attr.label(default = "//tools/openapi_generator", executable = True, cfg = "exec"),
        "_logging": attr.label(default = "//tools/openapi_generator:logback.xml", allow_single_file = True),
    },
)
