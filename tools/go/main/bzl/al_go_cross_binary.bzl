load("@rules_go//go:def.bzl", "go_cross_binary")
load("@rules_go//go/platform:list.bzl", "GOOS_GOARCH")

_INVALID_GOOS = ("qnx", "osx", "darwin", "ios", "android", "nacl")

def _goos_goarch():
    """
    Generate all valid GOOS/GOARCH targets

    Returns:
        GOOS, GOARCH tuple
    """
    res = []
    for goos, goarch in GOOS_GOARCH:
        if goos not in _INVALID_GOOS:
            res.append((goos, goarch))
    return res

def al_go_cross_binary(name, target, visibility):
    """
    Generate go_cross_binary targets for all valid GOOS_GOARCH and wrap them \
    in a filegroup

    Args:
        name: name
        target: go_cross_binary target
        visibility: visibility
    """
    names = []
    for goos, goarch in _goos_goarch():
        target_name = "{}.{}_{}".format(name, goos, goarch)
        names.append(target_name)
        go_cross_binary(
            name = target_name,
            platform = "@rules_go//go/toolchain:{}_{}".format(goos, goarch),
            target = target,
        )
    native.filegroup(
        name = name,
        srcs = names,
        visibility = visibility,
    )
