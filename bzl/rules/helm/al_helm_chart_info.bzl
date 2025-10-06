AlHelmChartInfo = provider(
    fields = {
        "source": "Chart sources (directory, optional)",
        "package": "Chart package (tgz file)",
        "deps": "Chart deps (depset of AlHelmChartInfo)",
    },
    doc = "Information about a helm chart",
)
