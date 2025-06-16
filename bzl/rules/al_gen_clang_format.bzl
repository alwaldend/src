def _impl(ctx):
    pass

al_gen_clang_format = rule(
    impl = _impl,
    toolchains = ["//bzl/toolchain-types:clang-format"],
)
