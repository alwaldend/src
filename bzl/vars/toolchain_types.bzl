_TOOLCHAIN_TYPES = ["al-qt", "al-drawio"]

TOOLCHAIN_TYPES = {
    toolchain_type: "//bzl/toolchain-types:{}".format(toolchain_type)
    for toolchain_type in _TOOLCHAIN_TYPES
}
