AlHelmChartInfo = provider(
    fields = {
        "source": "Chart sources (PackageFilegroupInfo, optional)",
        "package": "Chart package (tgz file)",
        "deps": "Chart deps (depset of AlHelmChartInfo)",
        "files_info": "Chart file structure (PackageFilesInfo)",
    },
    doc = "Information about a helm chart",
)
