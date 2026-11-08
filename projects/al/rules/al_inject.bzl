"""Flat named resource declarations for the built-in injector."""

load(":al_config.bzl", "al_inject_config")

def al_inject(name, config_name = None, **kwargs):
    """Declare an injector call using named res_<type> maps.

    Args:
        name: Bazel target and default runtime call name.
        config_name: Optional override of the runtime call name.
        **kwargs: Named resource maps prefixed res_, plus injector configuration attributes
            such as labels, extends, substitutions and data.

    Each resource entry contains structured type settings, optional deps and
    generic JSON-compatible data for local templates or provider input.
    Provider metadata remains opaque in args. Resource names must be unique across
    groups; dependency order is explicit. The prefix keeps resource maps separate
    from call attributes and supports new resource types without new parameters.
    """
    resources = {}
    for attribute in list(kwargs.keys()):
        if not attribute.startswith("res_"):
            continue
        kind = attribute[len("res_"):]
        if not kind:
            fail("Injector resource types must be nonempty")
        entries = kwargs.pop(attribute)
        if type(entries) != "dict":
            fail("Injector resource groups must be dictionaries: " + kind)
        for resource_name, settings in entries.items():
            if type(resource_name) != "string" or not resource_name:
                fail("Injector resource names must be nonempty strings")
            if resource_name in resources:
                fail("Duplicate injector resource: " + resource_name)
            if type(settings) != "dict":
                fail("Injector resource settings must be dictionaries: " + resource_name)
            values = dict(settings)
            resource = {"type": kind}
            if "deps" in values:
                resource["deps"] = values.pop("deps")
            if "data" in values:
                resource["data"] = values.pop("data")
            resource["args"] = values
            resources[resource_name] = resource
    al_inject_config(
        name = name,
        config_name = config_name,
        resources = resources,
        **kwargs
    )
