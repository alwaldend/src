load("//bzl/rules:al_toml_data.bzl", _al_toml_data = "al_toml_data")

def al_toml_data(**kwargs):
    """
    Wrapper around al_toml_data rule

    Args:
        **kwargs: kwargs for al_toml_data
    """
    _al_toml_data(**kwargs)
