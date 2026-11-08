"""Typed AL configuration targets with explicit extension relationships."""

AlConfigInfo = provider(
    doc = "A generated AL configuration and its runtime dependencies.",
    fields = {
        "config": "Merged protobuf JSON configuration file.",
        "definitions": "Named atomic definitions and their extension ancestry.",
        "kind": "Atomic definition kind, or empty for an aggregate.",
        "owner": "This configuration target's label.",
        "ancestors": "Atomic targets extended by this target.",
        "depth": "Number of explicit extensions in this target's chain.",
    },
)

def _runfile(file):
    if file.short_path.startswith("../"):
        return file.short_path[3:]
    return "_main/" + file.short_path

def _pb_json(value):
    # Starlark forbids recursion. Process nested JSON breadth-first, retaining
    # references to each child output until the next level has populated it.
    result = {}
    pending = [(value, result)]
    for _ in range(64):
        children = []
        for item, output in pending:
            kind = type(item)
            if item == None:
                continue
            if kind == "dict":
                values = {}
                output["val_map"] = {"val": values}
                for key, child in item.items():
                    values[key] = {}
                    children.append((child, values[key]))
            elif kind == "list" or kind == "tuple":
                values = [{} for _ in item]
                output["val_list"] = {"val": values}
                children.extend(zip(item, values))
            elif kind == "string":
                output["val_string"] = item
            elif kind == "bool":
                output["val_bool"] = item
            elif kind == "int":
                output["val_int64"] = str(item)
            elif kind == "float":
                output["val_float"] = item
            else:
                fail("Unsupported configuration argument type: " + kind)
        pending = children
        if not pending:
            return result
    fail("Configuration arguments exceed 64 nesting levels")

def _merge_definitions(definitions, incoming):
    for key, candidate in incoming.items():
        previous = definitions.get(key)
        if previous == None or previous.owner == candidate.owner:
            definitions[key] = candidate
        elif previous.owner in candidate.ancestors.to_list():
            definitions[key] = candidate
        elif candidate.owner not in previous.ancestors.to_list():
            fail("Conflicting %s configuration %r from %s and %s; declare an explicit extends relationship" % (
                key[0],
                key[1],
                previous.owner,
                candidate.owner,
            ))

def _definition_order(definition):
    return (definition.kind, definition.name)

def _emit(ctx, content, runtime = [], kind = "", runtime_parent = ""):
    source = ctx.actions.declare_file(ctx.label.name + "_/declarations.json")
    out = ctx.actions.declare_file(ctx.label.name + "_/al.json")
    ctx.actions.write(source, json.encode(content))
    runfiles = ctx.runfiles(files = ctx.files.data + runtime)
    for target in ctx.attr.data:
        runfiles = runfiles.merge(target[DefaultInfo].default_runfiles)

    definitions = {}
    dependencies = list(getattr(ctx.attr, "configs", []))
    ancestors = depset()
    depth = 0
    parent = getattr(ctx.attr, "extends", None)
    if parent:
        base = parent[AlConfigInfo]
        if not kind or base.kind != kind:
            fail("extends must select an atomic %s target" % kind)
        ancestors = depset([base.owner], transitive = [base.ancestors])
        depth = base.depth + 1
        dependencies.append(parent)
    for dependency in dependencies:
        _merge_definitions(definitions, dependency[AlConfigInfo].definitions)
        runfiles = runfiles.merge(dependency[DefaultInfo].default_runfiles)

    toolchain = ctx.toolchains[Label("//rules:al_toolchain_type")]
    declaration = source
    substitutions = getattr(ctx.attr, "substitutions", {})
    if substitutions:
        declaration = ctx.actions.declare_file(ctx.label.name + "_/local.json")
        inputs = [source]
        substitution_args = ctx.actions.args()
        substitution_args.add_all(["config", "dump", "--out", declaration.path, "--config", source.path])
        for target, marker in substitutions.items():
            files = target[DefaultInfo].files.to_list()
            if len(files) != 1:
                fail("Each substitution must select exactly one file")
            inputs.extend(files)
            substitution_args.add_all(["--substitution", marker + "=" + files[0].path])
        ctx.actions.run(
            executable = toolchain.executable,
            arguments = [substitution_args],
            inputs = depset(inputs),
            outputs = [declaration],
            mnemonic = "AlConfigSubstitution",
        )
    if kind:
        definition = struct(
            owner = ctx.label,
            ancestors = ancestors,
            depth = depth,
            kind = kind,
            name = ctx.attr.config_name,
            file = declaration,
            runtime_parent = runtime_parent,
        )
        _merge_definitions(definitions, {(kind, ctx.attr.config_name): definition})
    inputs = [definition.file for definition in sorted(definitions.values(), key = _definition_order)]
    if not inputs:
        inputs = [source]
    args = ctx.actions.args()
    args.add_all(["config", "dump", "--out", out.path])
    for file in inputs:
        args.add_all(["--config", file.path])
    ctx.actions.run(
        executable = toolchain.executable,
        arguments = [args],
        inputs = depset(inputs),
        outputs = [out],
        mnemonic = "AlConfig",
    )
    return [
        DefaultInfo(files = depset([out]), runfiles = runfiles),
        AlConfigInfo(
            config = out,
            definitions = definitions,
            kind = kind,
            owner = ctx.label,
            ancestors = ancestors,
            depth = depth,
        ),
    ]

def _config_impl(ctx):
    return _emit(ctx, {})

def _auth_impl(ctx):
    auth = {}
    if ctx.attr.approle:
        auth["approle"] = {"name": ctx.attr.approle}
    if ctx.attr.token_helper:
        auth["token_helper"] = {}
    if ctx.attr.no_auth:
        auth["no_auth"] = True
    if len([value for value in [ctx.attr.approle, ctx.attr.token_helper, ctx.attr.no_auth] if value]) > 1:
        fail("Select at most one Vault authentication method")
    return _emit(ctx, {"vault_auth": {ctx.attr.config_name: auth}}, kind = "vault_auth")

def _conn_impl(ctx):
    tls = {}
    runtime = []
    for field in ["ca_cert", "client_cert", "client_key"]:
        file = getattr(ctx.file, field)
        if file:
            tls[field] = _runfile(file)
            runtime.append(file)
    return _emit(ctx, {"vault_conn": {ctx.attr.config_name: {
        "config": {"address": ctx.attr.address},
        "tls": tls,
    }}}, runtime, kind = "vault_conn")

def _plugin_config_impl(ctx):
    content = {
        "bin": _runfile(ctx.executable.executable),
        "args": ctx.attr.argv,
        "env": ctx.attr.env,
        "labels": ctx.attr.labels,
        "resource_provider": ctx.attr.resource_provider,
    }
    if ctx.attr.args_json:
        content["data"] = _pb_json(json.decode(ctx.attr.args_json))
    result = _emit(ctx, {"plugin_configs": {ctx.attr.config_name: content}}, [ctx.executable.executable], kind = "plugin_configs")
    default = result[0]
    return [
        DefaultInfo(
            files = default.files,
            runfiles = default.default_runfiles.merge(ctx.attr.executable[DefaultInfo].default_runfiles),
        ),
        result[1],
    ]

def _plugin_impl(ctx):
    content = {"plugin": ctx.attr.plugin, "labels": ctx.attr.labels}
    if ctx.attr.args_json:
        content["data"] = _pb_json(json.decode(ctx.attr.args_json))
    return _emit(ctx, {"plugins": {ctx.attr.config_name: content}}, kind = "plugins")

def _attrs(extra):
    return dict({
        "data": attr.label_list(allow_files = True, doc = "Runtime inputs."),
        "extends": attr.label(providers = [AlConfigInfo], doc = "Atomic target extended by this definition."),
    }, **extra)

def _injector_impl(ctx):
    resources = json.decode(ctx.attr.resources_json)
    if type(resources) != "dict":
        fail("Injector resources must be a dictionary keyed by name")
    for resource in resources.values():
        if type(resource) != "dict":
            fail("Injector resource definitions must be dictionaries")
        resource["args"] = _pb_json(resource.get("args", {}))
        if "data" in resource:
            resource["data"] = _pb_json(resource["data"])
    runtime_parent = ""
    if ctx.attr.extends:
        base = ctx.attr.extends[AlConfigInfo]
        own = [definition for definition in base.definitions.values() if definition.owner == base.owner]
        if len(own) != 1:
            fail("Injector extends must select one named group")
        parent = own[0]
        runtime_parent = parent.name if parent.name != ctx.attr.config_name else parent.runtime_parent
    content = {"labels": ctx.attr.labels, "res": resources}
    if runtime_parent:
        content["extends"] = runtime_parent
    return _emit(ctx, {"inject": {ctx.attr.config_name: content}}, kind = "inject", runtime_parent = runtime_parent)

_TOOLCHAINS = [Label("//rules:al_toolchain_type")]

_al_config = rule(
    doc = "Combine named configuration targets independently of input order.",
    implementation = _config_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({"configs": attr.label_list(providers = [AlConfigInfo])}),
)

al_vault_auth = rule(
    doc = "Declare a named Vault authentication method.",
    implementation = _auth_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({
        "config_name": attr.string(default = "default"),
        "approle": attr.string(),
        "token_helper": attr.bool(),
        "no_auth": attr.bool(),
    }),
)

al_vault_conn = rule(
    doc = "Declare a Vault connection with label-based TLS files.",
    implementation = _conn_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({
        "config_name": attr.string(default = "default"),
        "address": attr.string(mandatory = True),
        "ca_cert": attr.label(allow_single_file = True),
        "client_cert": attr.label(allow_single_file = True),
        "client_key": attr.label(allow_single_file = True),
    }),
)

_al_plugin_config = rule(
    implementation = _plugin_config_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({
        "config_name": attr.string(mandatory = True),
        "executable": attr.label(executable = True, cfg = "target", mandatory = True),
        "argv": attr.string_list(),
        "env": attr.string_dict(),
        "labels": attr.string_dict(),
        "args_json": attr.string(),
        "resource_provider": attr.bool(),
    }),
)

_al_plugin = rule(
    implementation = _plugin_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({
        "config_name": attr.string(mandatory = True),
        "plugin": attr.string(mandatory = True),
        "labels": attr.string_dict(),
        "args_json": attr.string(),
        "substitutions": attr.label_keyed_string_dict(allow_files = True),
    }),
)

def al_plugin_config(name, executable, config_name = None, args = None, **kwargs):
    """Declare a plugin executable and its ordinary JSON-compatible parameters."""
    _al_plugin_config(
        name = name,
        executable = executable,
        config_name = config_name or name,
        args_json = json.encode(args) if args != None else "",
        **kwargs
    )

_al_inject_config = rule(
    implementation = _injector_impl,
    provides = [AlConfigInfo],
    toolchains = _TOOLCHAINS,
    attrs = _attrs({
        "config_name": attr.string(mandatory = True),
        "labels": attr.string_dict(),
        "resources_json": attr.string(mandatory = True),
        "substitutions": attr.label_keyed_string_dict(allow_files = True),
    }),
)

def al_inject_config(name, resources = {}, config_name = None, **kwargs):
    """Declare built-in resources independently of external plugin calls."""
    _al_inject_config(
        name = name,
        config_name = config_name or name,
        resources_json = json.encode(resources),
        **kwargs
    )

def al_plugin(name, plugin = None, config_name = None, args = None, **kwargs):
    """Declare a plugin call; runtime call and plugin names default to name."""
    _al_plugin(
        name = name,
        plugin = plugin or name,
        config_name = config_name or name,
        args_json = json.encode(args) if args != None else "",
        **kwargs
    )

def al_config(name, config = {}, configs = [], **kwargs):
    """Build named configuration dictionaries and compose inherited targets.

    Args:
        name: Aggregate target name.
        config: Dictionaries of vault_auth, vault_conn, plugin_configs, plugins and inject.
        configs: Configuration targets combined by identity and explicit extensions.
        **kwargs: Aggregate target attributes; test/platform attributes propagate.
    """
    if type(config) != "dict":
        fail("AL config must be a dictionary")
    constructors = {
        "vault_auth": al_vault_auth,
        "vault_conn": al_vault_conn,
        "plugin_configs": al_plugin_config,
        "plugins": al_plugin,
        "inject": al_inject_config,
    }
    suffixes = {"vault_auth": "vault_auth", "vault_conn": "vault_conn", "plugin_configs": "plugin_config", "plugins": "plugin", "inject": "inject"}
    fragments = []
    common = {key: kwargs[key] for key in ["testonly", "tags", "target_compatible_with", "exec_compatible_with", "features"] if key in kwargs}
    for section, entries in config.items():
        if section not in constructors:
            fail("Unknown AL configuration section: " + section)
        if type(entries) != "dict":
            fail("AL configuration sections must be dictionaries keyed by name")
        for config_name, attributes in entries.items():
            if type(attributes) != "dict":
                fail("AL configuration entries must be attribute dictionaries")
            if "name" in attributes or "config_name" in attributes:
                fail("The dictionary key owns each AL configuration entry name")
            if section == "plugins" and "plugin" not in attributes:
                attributes = dict(attributes, plugin = config_name)
            target = name + "." + suffixes[section] + "." + config_name
            constructors[section](
                name = target,
                config_name = config_name,
                visibility = ["//visibility:private"],
                **dict(common, **attributes)
            )
            fragments.append(":" + target)
    _al_config(name = name, configs = configs + fragments, **kwargs)
