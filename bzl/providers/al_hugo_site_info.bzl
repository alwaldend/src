AlHugoSiteInfo = provider(
    fields = {
        "tree": "Hugo site tree archive (file, .tar))",
        "env": "Environment variable (dict[str, str])",
        "env_script": "Shell script to export env variables",
    },
    doc = "Information about a hugo site",
)
