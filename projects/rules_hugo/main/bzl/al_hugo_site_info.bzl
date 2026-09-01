AlHugoSiteInfo = provider(
    fields = {
        "site_archive": "Site archive File (.tar)",
        "postcss": "PostCSS executable File",
        "postcss_files_to_run": "PostCSS FilesToRunProvider",
        "env": "Environment variables",
        "env_script": "Shell script to export env variables",
    },
    doc = "Information about a hugo site",
)
