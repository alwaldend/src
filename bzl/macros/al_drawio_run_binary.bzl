load("//bzl/rules:al_drawio_run_binary.bzl", _impl = "al_drawio_run_binary")

def al_drawio_run_binary(**kwargs):
    """
    Wrapper around al_drawio_run_binary rule

    Args:
        **kwargs: kwargs for al_drawio_run_binary
    """
    _impl(**kwargs)
