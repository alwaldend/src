def al_alias_map(aliases, visibility = ["//visibility:public"]):
    """
    Generate aliases from an alias map


    Args:
        aliases: alias map, keys are names, values are alias arguments
        visibility: default visibility
    """

    for name, config in aliases.items():
        native.alias(
            name = name,
            actual = config["actual"],
            visibility = config.get("visibility", visibility),
        )
