"""Select the reactor startup object from the active Rust WASI toolchain."""

def _wasi_reactor_crt_impl(ctx):
    toolchain = ctx.toolchains["@rules_rust//rust:toolchain_type"]
    objects = [
        file
        for file in toolchain.all_files.to_list()
        if file.is_source and file.path.endswith("/lib/rustlib/wasm32-wasip1/lib/self-contained/crt1-reactor.o")
    ]
    if len(objects) != 1:
        fail("The Rust WASI toolchain must supply exactly one crt1-reactor.o")
    return [DefaultInfo(files = depset(objects))]

wasi_reactor_crt = rule(
    implementation = _wasi_reactor_crt_impl,
    toolchains = ["@rules_rust//rust:toolchain_type"],
)
