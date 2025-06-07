def _al_write_script_impl(ctx):
    ctx.actions.write(
        out = ctx.outputs.out,
        content = "\n".join([
            "{shebang}",
            "set {set_flags}",
            "",
            "{content}",
        ]).format(
            set_flags = " ".join(ctx.attr.set_flags),
            shebang = ctx.attr.shebang,
            content = ctx.expand_make_variables(
                "script",
                ctx.expand_location(ctx.attr.content),
                ctx.attr.make_vars,
            ),
        ),
        is_executable = True,
    )

al_write_script = rule(
    implementation = _al_write_script_impl,
    doc = "Write a script and make it executable",
    attrs = {
        "content": attr.string(mandatory = True, doc = "Script content"),
        "out": attr.output(mandatory = True, doc = "Output file"),
        "shebang": attr.string(default = "#!/usr/bin/env sh", doc = "Sheband to use"),
        "set_flags": attr.string_list(default = ["-eu"], doc = "Flags to pass to set"),
        "make_vars": attr.string_dict(default = {}, doc = "Additional make vars"),
    },
)
