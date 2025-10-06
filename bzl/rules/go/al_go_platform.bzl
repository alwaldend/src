def al_go_platform(rctx):
    """
    Returns a normalized name of the host os and CPU architecture.

    Source: https://github.com/abrisco/rules_helm/blob/fafa9dfa6eded68d313a385aabdae9cc9e242d57/helm/repositories.bzl#L113C1-L158C32

    Alias archictures names are normalized:

    x86_64 => amd64
    aarch64 => arm64

    The result can be used to generate repository names for host toolchain
    repositories for toolchains that use these normalized names.

    Common os & architecture pairs that are returned are,

    - darwin_amd64
    - darwin_arm64
    - linux_amd64
    - linux_arm64
    - linux_s390x
    - linux_ppc64le
    - windows_amd64

    Args:
        rctx: rctx

    Returns:
        The normalized "<os>_<arch>" string of the host os and CPU architecture.
    """
    name = rctx.os.name.lower()
    if name.startswith("linux"):
        os = "linux"
    elif name.startswith("mac os"):
        os = "darwin"
    elif name.startswith("freebsd"):
        os = "freebsd"
    elif name.find("windows") != -1:
        os = "windows"
    else:
        fail("unrecognized os")

    arch = rctx.os.arch
    if arch == "aarch64":
        arch = "arm64"
    elif arch == "x86_64":
        arch = "amd64"

    return "{}_{}".format(os, arch)
