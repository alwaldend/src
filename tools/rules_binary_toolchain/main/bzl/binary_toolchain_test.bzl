load(":binary_toolchain_binary.bzl", "binary_toolchain_binary")

def binary_toolchain_test(toolchain_type, kwargs = {}, attrs = {}):
    """
    Args:
        toolchain_type: toolchain type
        kwargs: rule kwargs override
        attrs: attrs override

    Returns:
        Binary toolchain test rule
    """
    default_kwargs = {
        "executable": False,
        "test": True,
    }
    kwargs = default_kwargs | kwargs
    return binary_toolchain_binary(toolchain_type, kwargs, attrs)
