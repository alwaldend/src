"""Terraform executables with declared configuration and provider runfiles."""

load("@rules_terraform_cli//:defs.bzl", "terraform_toolchain_type")
load(":providers.bzl", "TerraformProviderInfo")

TerraformRunInfo = provider(
    doc = "Declared Terraform invocation and provider mirror.",
    fields = {
        "config": "Generated invocation configuration file.",
        "providers": "Selected TerraformProviderInfo values.",
    },
)

def _runfile(ctx, file):
    if file.short_path.startswith("../"):
        return file.short_path[3:]
    return ctx.workspace_name + "/" + file.short_path

def _terraform_impl(ctx):
    executable = ctx.actions.declare_file(ctx.label.name)
    config = ctx.actions.declare_file(ctx.label.name + ".json")
    marker = ctx.actions.declare_file(ctx.label.name + ".providers.marker")
    toolchain = ctx.toolchains[terraform_toolchain_type]
    terraform = ctx.executable.terraform or toolchain.binary
    package_prefix = ctx.label.package + "/" if ctx.label.package else ""
    mirror_path = package_prefix + ctx.label.name + ".providers"
    symlinks = {mirror_path + "/.bazel_provider_mirror": marker}
    selected = []
    provider_files = []
    versions = {}
    platforms = {}
    for target in ctx.attr.providers:
        info = target[TerraformProviderInfo]
        if info.source in versions and versions[info.source] != info.version:
            fail("conflicting versions for provider {}: {} and {}".format(info.source, versions[info.source], info.version))
        versions[info.source] = info.version
        key = info.source + "/" + info.platform
        if key in platforms:
            fail("duplicate provider selection: " + key)
        platforms[key] = True
        selected.append(info)
        provider_files.append(info.archive)
        symlinks[mirror_path + "/" + info.mirror_path] = info.archive

    data_runfiles = [target[DefaultInfo].default_runfiles for target in ctx.attr.data]
    inputs = depset(ctx.files.srcs + ctx.files.data, transitive = [runfiles.files for runfiles in data_runfiles]).to_list()
    owner = _runfile(ctx, config).removesuffix("/" + package_prefix + ctx.label.name + ".json")
    chdir = ctx.attr.chdir if ctx.attr.chdir else (ctx.label.package or ".")
    ctx.actions.write(marker, "Bazel-managed Terraform provider mirror\n")
    ctx.actions.write(config, json.encode_indent({
        "terraform": _runfile(ctx, terraform),
        "workspace": owner,
        "chdir": chdir,
        "working_directory": owner + ("/" + chdir if chdir != "." else ""),
        "args": ctx.attr.arguments,
        "mirror": owner + "/" + mirror_path,
        "providers": [{
            "source": info.source,
            "version": info.version,
            "platform": info.platform,
            "archive": _runfile(ctx, info.archive),
            "mirror_path": info.mirror_path,
        } for info in selected],
        "files": [{"runfile": _runfile(ctx, file), "path": _runfile(ctx, file).removeprefix(owner + "/")} for file in inputs if _runfile(ctx, file).startswith(owner + "/")],
    }) + "\n")
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable._launcher,
        is_executable = True,
    )
    runfiles = ctx.runfiles(
        files = inputs + provider_files + [config, marker, terraform],
        root_symlinks = {owner + "/" + path: file for path, file in symlinks.items()},
    ).merge_all(data_runfiles + [ctx.attr._launcher[DefaultInfo].default_runfiles])
    if ctx.attr.terraform:
        runfiles = runfiles.merge(ctx.attr.terraform[DefaultInfo].default_runfiles)
    else:
        runfiles = runfiles.merge(toolchain.default_info.default_runfiles)
    return [
        DefaultInfo(files = depset([executable]), executable = executable, runfiles = runfiles),
        TerraformRunInfo(config = config, providers = selected),
    ]

_ATTRS = {
    "arguments": attr.string_list(doc = "Fixed runner options and Terraform arguments; invocation arguments are appended."),
    "srcs": attr.label_list(allow_files = True, doc = "Terraform configuration and local module files."),
    "data": attr.label_list(allow_files = True, doc = "Declared files used by the configuration, including local modules."),
    "providers": attr.label_list(providers = [TerraformProviderInfo], doc = "Pinned provider archives available to this invocation."),
    "chdir": attr.string(doc = "Working directory relative to the consuming workspace; defaults to this package."),
    "terraform": attr.label(executable = True, cfg = "exec", doc = "Optional executable override; otherwise use the pinned Terraform toolchain."),
    "_launcher": attr.label(default = Label("//cmd/terraform:terraform"), executable = True, cfg = "target"),
}

terraform_binary = rule(
    implementation = _terraform_impl,
    attrs = _ATTRS,
    executable = True,
    toolchains = [terraform_toolchain_type],
    doc = "Run Terraform with a packaged filesystem provider mirror and declared inputs.",
)

terraform_test = rule(
    implementation = _terraform_impl,
    attrs = _ATTRS,
    test = True,
    toolchains = [terraform_toolchain_type],
    doc = "Run an offline Terraform check using the same executable and provider contract.",
)

def _terraform_invocation(
        name,
        arguments,
        srcs = [],
        data = [],
        providers = [],
        chdir = None,
        terraform = None,
        is_test = False,
        wrapper = None,
        wrapper_kwargs = {},
        **kwargs):
    if wrapper_kwargs and wrapper == None:
        fail("wrapper_kwargs requires a wrapper")
    kwargs["testonly"] = kwargs.get("testonly", False) or is_test
    for key in wrapper_kwargs:
        if key in ["name", "args", "data"] or key in kwargs:
            fail("wrapper_kwargs cannot override " + key)
    invocation_kwargs = {
        "arguments": arguments,
        "srcs": srcs,
        "data": data,
        "providers": providers,
        "chdir": chdir,
        "terraform": terraform,
    }
    if wrapper == None:
        invocation = terraform_test if is_test else terraform_binary
        invocation(name = name, **(invocation_kwargs | kwargs))
        return

    invocation_name = name + ".terraform"
    inner_kwargs = {key: kwargs[key] for key in [
        "testonly",
        "tags",
        "target_compatible_with",
        "exec_compatible_with",
        "features",
    ] if key in kwargs}
    terraform_binary(
        name = invocation_name,
        visibility = ["//visibility:private"],
        **(invocation_kwargs | inner_kwargs)
    )
    wrapper(
        name = name,
        args = ["$(rootpath :{})".format(invocation_name)],
        data = data + [":" + invocation_name],
        **(dict(wrapper_kwargs) | kwargs)
    )

DEFAULT_TERRAFORM_BINARIES = {
    "direct": ["--direct"],
    "init": ["--direct", "init"],
    "migrate": ["--direct", "init", "--migrate-state"],
    "fmt": ["--direct", "fmt", "--write", "--recursive"],
    "fmt_check": ["--direct", "fmt", "--check", "--recursive"],
    "plan": ["plan"],
    "output": ["output"],
    "import": ["import"],
    "destroy": ["destroy"],
    "apply": ["apply"],
    "show": ["show"],
    "deploy": ["apply"],
    "deploy_y": ["apply", "-y"],
    "state": ["state"],
    "force_unlock": ["force-unlock"],
}

def terraform_binary_map(
        name,
        args = DEFAULT_TERRAFORM_BINARIES,
        wrapper = None,
        wrapper_kwargs = {},
        **kwargs):
    """
    Create several explicitly named Terraform binaries.

    The unnamed target used to alias `apply`. Mutating operations must remain
    explicit, so callers use `<name>.apply` instead.
    """
    for args_name, args_value in args.items():
        if args_name:
            cur_name = "{}.{}".format(name, args_name)
        else:
            cur_name = name
        _terraform_invocation(
            name = cur_name,
            arguments = args_value,
            wrapper = wrapper,
            wrapper_kwargs = wrapper_kwargs,
            **kwargs
        )

DEFAULT_TERRAFORM_TESTS = {
    "fmt": ["--direct", "fmt", "--check", "--recursive"],
}

def terraform_target_binary_map(name, target, **kwargs):
    """Plan one target in an owning root and apply only a saved plan."""
    terraform_binary_map(
        name = name,
        args = {
            "plan": ["plan", "-target=" + target],
            "show": ["show"],
            "apply": ["--require-saved-plan", "apply"],
        },
        **kwargs
    )

def terraform_test_map(
        name,
        args = DEFAULT_TERRAFORM_TESTS,
        size = "small",
        wrapper = None,
        wrapper_kwargs = {},
        **kwargs):
    """
    Create several terraform tests
    """
    for args_name, args_value in args.items():
        if not args_name:
            fail("empty test name")
        _terraform_invocation(
            name = "{}.{}_test".format(name, args_name),
            size = size,
            arguments = args_value,
            wrapper = wrapper,
            wrapper_kwargs = wrapper_kwargs,
            is_test = True,
            **kwargs
        )
