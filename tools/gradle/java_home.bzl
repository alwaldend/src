load("@rules_java//java/common:java_common.bzl", "java_common")

def _java_home_impl(ctx):
    java_runtime_info = ctx.attr.runtime[java_common.JavaRuntimeInfo] if java_common.JavaRuntimeInfo in ctx.attr.runtime else None
    java_runtime = java_runtime_info
    java_home = java_runtime.java_home
    java = java_runtime.java_executable_runfiles_path
    runfiles = ctx.runfiles().merge(ctx.attr.runtime[DefaultInfo].default_runfiles)
    return [
        DefaultInfo(runfiles = runfiles),
        platform_common.ToolchainInfo(
            java_home = java_home,
            java = java,
        ),
    ]

java_home = rule(
    doc = "Exposes the Java runtime as a single-file binary.",
    implementation = _java_home_impl,
    attrs = {
        "runtime": attr.label(mandatory = True),
    },
)
