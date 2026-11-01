def _impl(ctx):
    files = {"BUILD.bazel": ""} | ctx.attr.files
    for path, content in files.items():
        ctx.file(path, content)
    return ctx.repo_metadata(
        reproducible = True,
    )

al_bzl_generate_repository = repository_rule(
    implementation = _impl,
    doc = "Generate a repository",
    attrs = {
        "files": attr.string_dict(
            doc = "Files to generate, keys are paths, values are file contents",
        ),
    },
)
