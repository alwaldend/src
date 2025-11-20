AlReleaseFilesInfo = provider(
    doc = "Release files",
    fields = {
        "srcs": "Sources (depset of File)",
        "deps": "Deps (Depset of AlReleaseFilesInfo)",
        "manifest": "Release manifest for srcs",
    },
)
