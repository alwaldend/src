load(":al_md_data.bzl", "al_md_data")

def al_readme(name, **kwargs):
    """
    Create readme targets

    Args:
        name: target name
        **kwargs: al_md_data kwargs
    """
    al_md_data(
        name = name,
        srcs = [":README.md"],
        tags = ["al-readme"],
        **kwargs
    )
