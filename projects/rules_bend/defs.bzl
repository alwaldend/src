"""Bazel rules for running and checking Bend 2 programs."""

_EXEC_TOOLCHAIN_TYPE = Label("//projects/rules_bend:exec_toolchain_type")
_RUNTIME_TOOLCHAIN_TYPE = Label("//projects/rules_bend:runtime_toolchain_type")
_CHECK_EXEC_GROUP = "bend_check"

def _bend_toolchain_impl(ctx):
    compiler = ctx.attr.compiler[DefaultInfo]
    return [platform_common.ToolchainInfo(compiler = compiler)]

def _toolchain_implementation(configuration):
    return rule(
        implementation = _bend_toolchain_impl,
        attrs = {
            "compiler": attr.label(cfg = configuration, executable = True, mandatory = True),
        },
    )

_bend_exec_toolchain = _toolchain_implementation("exec")
_bend_runtime_toolchain = _toolchain_implementation("target")

def bend_toolchain(name, compiler, compatible_with = [], visibility = None):
    """Defines registrations named <name>_exec and <name>_runtime.

    Args:
        name: Prefix for the toolchain registrations.
        compiler: Executable Bend 2 target with its runtime files.
        compatible_with: Platform constraints for the compiler.
        visibility: Visibility of the registration targets.
    """
    _bend_exec_toolchain(name = name + "_exec_impl", compiler = compiler)
    _bend_runtime_toolchain(name = name + "_runtime_impl", compiler = compiler)
    native.toolchain(
        name = name + "_exec",
        toolchain = ":" + name + "_exec_impl",
        toolchain_type = _EXEC_TOOLCHAIN_TYPE,
        exec_compatible_with = compatible_with,
        visibility = visibility,
    )
    native.toolchain(
        name = name + "_runtime",
        toolchain = ":" + name + "_runtime_impl",
        toolchain_type = _RUNTIME_TOOLCHAIN_TYPE,
        target_compatible_with = compatible_with,
        visibility = visibility,
    )

def _source_path(file):
    path = file.short_path
    if path.startswith("../"):
        return "external/" + path[3:]
    return "workspace/" + path

def _runfile_path(ctx, file):
    path = file.short_path
    if path.startswith("../"):
        return path[3:]
    return ctx.workspace_name + "/" + path

def _inputs(ctx):
    files = depset([ctx.file.main], transitive = [
        target[DefaultInfo].files
        for target in ctx.attr.srcs
    ])
    sources = {}
    for file in files.to_list():
        if file.is_directory or file.extension != "bend":
            fail("expected a .bend source file, got " + file.path, attr = "srcs")
        path = _source_path(file)
        if path in sources:
            fail("multiple sources stage to " + path, attr = "srcs")
        sources[path] = file
    return files, sources

def _bend_check_impl(ctx):
    files, sources = _inputs(ctx)
    compiler = ctx.exec_groups[_CHECK_EXEC_GROUP].toolchains[_EXEC_TOOLCHAIN_TYPE].compiler.files_to_run
    output = ctx.actions.declare_file(ctx.label.name + ".checked")
    manifest = ctx.actions.declare_file(ctx.label.name + ".bend-inputs.json")
    ctx.actions.write(manifest, json.encode({
        "compiler": compiler.executable.path,
        "main": _source_path(ctx.file.main),
        "output": output.path,
        "sources": {path: file.path for path, file in sources.items()},
    }))
    ctx.actions.run(
        executable = ctx.executable._runner,
        tools = [compiler],
        inputs = depset([manifest], transitive = [files]),
        outputs = [output],
        arguments = [manifest.path],
        execution_requirements = {"block-network": "1"},
        exec_group = _CHECK_EXEC_GROUP,
        mnemonic = "BendCheck",
        progress_message = "Checking %{label}",
    )
    return [DefaultInfo(files = depset([output]))]

def _bend_binary_impl(ctx):
    files, sources = _inputs(ctx)
    compiler = ctx.toolchains[_RUNTIME_TOOLCHAIN_TYPE].compiler
    binary = compiler.files_to_run.executable
    launcher_path = ctx.label.name + ".bend-launcher/" + ctx.label.name.split("/")[-1]
    executable = ctx.actions.declare_file(launcher_path)
    manifest = ctx.actions.declare_file(launcher_path + ".bend.json")
    ctx.actions.write(manifest, json.encode({
        "compiler": _runfile_path(ctx, binary),
        "main": _source_path(ctx.file.main),
        "sources": {path: _runfile_path(ctx, file) for path, file in sources.items()},
    }))
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable._runner,
        is_executable = True,
    )
    runfiles = ctx.runfiles(
        files = [manifest, binary],
        transitive_files = files,
    ).merge(compiler.default_runfiles)
    return [DefaultInfo(
        executable = executable,
        files = depset([executable, manifest]),
        runfiles = runfiles,
    )]

def _bend_attrs(configuration):
    return {
        "main": attr.label(
            doc = "Entry point passed to Bend.",
            allow_single_file = [".bend"],
            mandatory = True,
        ),
        "srcs": attr.label_list(
            doc = "Local .bend imports and other Bend inputs, including transitive imports.",
            allow_files = [".bend"],
        ),
        "_runner": attr.label(
            default = Label("//projects/rules_bend/cmd/bend_action"),
            cfg = configuration,
            executable = True,
        ),
    }

bend_binary = rule(
    doc = "Runs a Bend 2 program with its declared imports.",
    implementation = _bend_binary_impl,
    attrs = _bend_attrs("target"),
    toolchains = [_RUNTIME_TOOLCHAIN_TYPE],
    executable = True,
)

bend_check = rule(
    doc = "Checks a Bend 2 entry point and emits a .checked marker on success.",
    implementation = _bend_check_impl,
    attrs = _bend_attrs(config.exec(exec_group = _CHECK_EXEC_GROUP)),
    exec_groups = {_CHECK_EXEC_GROUP: exec_group(toolchains = [_EXEC_TOOLCHAIN_TYPE])},
)
