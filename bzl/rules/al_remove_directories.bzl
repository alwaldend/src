def _impl(ctx):
    files = []
    for file in ctx.files.srcs:
        if not file.is_directory:
            files.append(file)
    return [
        DefaultInfo(
            files = depset(files),
        ),
    ]

al_remove_directories = rule(
    implementation = _impl,
    doc = "Remove directory files from DefaultInfo",
    attrs = {
        "srcs": attr.label_list(doc = "Sources to remove directories from"),
    },
)
