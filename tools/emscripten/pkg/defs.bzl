"""Public Emscripten build toolchain and explicitly declared JavaScript outputs."""

_EMSCRIPTEN = "//tools/emscripten:toolchain_type"
_NODE = "@rules_nodejs//nodejs:toolchain_type"
_PYTHON = "@rules_python//python:toolchain_type"

def _toolchain_impl(ctx):
    return [platform_common.ToolchainInfo(
        emcc = ctx.file.emcc,
        files = ctx.attr.sdk[DefaultInfo].files,
    )]

emscripten_toolchain = rule(
    implementation = _toolchain_impl,
    doc = "Declare an official SDK archive for the Emscripten build toolchain.",
    attrs = {
        "emcc": attr.label(allow_single_file = True, mandatory = True),
        "sdk": attr.label(mandatory = True),
    },
)

def _javascript_impl(ctx):
    sdk = ctx.toolchains[_EMSCRIPTEN]
    python = ctx.toolchains[_PYTHON].py3_runtime
    node = ctx.toolchains[_NODE].nodeinfo.node
    if not python.interpreter or not node:
        fail("Emscripten requires declared Python and Node interpreter files")

    output = ctx.outputs.js
    if not output.basename.endswith(".js"):
        fail("js must name a .js output")
    outputs = [output]
    if ctx.attr.wasm:
        outputs.append(ctx.actions.declare_file(output.basename.removesuffix(".js") + ".wasm", sibling = output))
    config = ctx.actions.declare_file(ctx.label.name + ".em_config.py")
    sdk_root = sdk.emcc.dirname.removesuffix("/emscripten")
    ctx.actions.write(config, "\n".join([
        "import os",
        "LLVM_ROOT = os.path.abspath(%r)" % (sdk_root + "/bin"),
        "BINARYEN_ROOT = os.path.abspath(%r)" % sdk_root,
        "NODE_JS = [os.path.abspath(%r)]" % node.path,
        "CACHE = os.path.abspath(%r)" % (sdk.emcc.dirname + "/cache"),
        "FROZEN_CACHE = True",
        "",
    ]))
    args = ctx.actions.args()
    args.add(sdk.emcc)
    args.add("--em-config", config)
    args.add_all(ctx.files.srcs)
    args.add_all(ctx.attr.arguments)
    args.add("-sWASM=" + ("1" if ctx.attr.wasm else "0"))
    args.add("-o", output)
    ctx.actions.run(
        executable = python.interpreter,
        arguments = [args],
        inputs = depset(ctx.files.srcs + [config], transitive = [sdk.files]),
        tools = depset([node, python.interpreter], transitive = [python.files]),
        outputs = outputs,
        env = {
            "EMCC_CORES": "1",
            "PYTHONDONTWRITEBYTECODE": "1",
            "TMPDIR": output.dirname,
        },
        mnemonic = "EmscriptenJavaScript",
        progress_message = "Generating Emscripten JavaScript %{label}",
    )
    files = depset(outputs)
    return [DefaultInfo(files = files, runfiles = ctx.runfiles(transitive_files = files))]

emscripten_javascript = rule(
    implementation = _javascript_impl,
    doc = "Compile C inputs to JavaScript and an optional declared Wasm sidecar.",
    attrs = {
        "arguments": attr.string_list(),
        "js": attr.output(mandatory = True),
        "srcs": attr.label_list(allow_files = [".c"], mandatory = True),
        "wasm": attr.bool(default = True),
    },
    toolchains = [_EMSCRIPTEN, _NODE, _PYTHON],
)
