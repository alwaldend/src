AlReleaseInfo = provider(
    doc = "Release information",
    fields = {
        "release_name": "Release name (string)",
        "project": "Project subdir (string)",
        "files": "File dict, keys are filenames, values are Files",
        "manifest": "Release manifest (File)",
    },
)
