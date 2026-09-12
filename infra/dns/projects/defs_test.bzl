"""Check external source identity and paths in nested DNS runtime packaging."""

load("@bazel_skylib//lib:unittest.bzl", "analysistest", "asserts")

def _project_files_test_impl(ctx):
    env = analysistest.begin(ctx)
    target = analysistest.target_under_test(env)
    aliases = {
        entry.path: entry.target_file
        for entry in target[DefaultInfo].default_runfiles.symlinks.to_list()
    }
    records = ctx.attr.records[DefaultInfo].files.to_list()
    asserts.equals(env, 1, len(records))
    asserts.equals(env, records[0], aliases.get("projects/rules_docs/dnsconfig.json"))
    sources = ctx.attr.sources[DefaultInfo].files.to_list()
    asserts.equals(env, len(sources) + 1, len(aliases))
    for source in sources:
        asserts.equals(env, source, aliases.get("projects/rules_docs/tf/" + source.basename))
    return analysistest.end(env)

project_files_test = analysistest.make(
    _project_files_test_impl,
    attrs = {
        "records": attr.label(mandatory = True),
        "sources": attr.label(mandatory = True),
    },
)
