AlReleaseFilesInfo = provider(
    doc = "Release files",
    fields = {
        "files": "File dict, keys are filenames, values are Files",
        "manifest": "Release manifest for srcs",
    },
)
