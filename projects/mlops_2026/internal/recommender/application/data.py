"""Resolve packaged course inputs independently of the working directory."""

import pathlib

from python.runfiles import runfiles


def _runfile(path: str, missing_message: str) -> pathlib.Path:
    resolver: runfiles.Runfiles | None = runfiles.Create()
    if resolver is None:
        raise RuntimeError("Application runfiles are unavailable")
    location: str | None = resolver.Rlocation(path, source_repo="")
    if location is None or not pathlib.Path(location).is_file():
        raise FileNotFoundError(missing_message)
    return pathlib.Path(location)


def course_data_paths() -> dict[str, pathlib.Path]:
    paths: dict[str, pathlib.Path] = {}
    name: str
    for name in ("data.csv", "dataset.csv"):
        paths[name] = _runfile(
            "com_github_yandex_practicum_mlops_freetrack/" + name,
            f"Course runfile is unavailable: {name}",
        )
    return paths


def model_path() -> pathlib.Path:
    return _runfile(
        "_main/projects/mlops_2026/model.pkl",
        "Recommendation model runfile is unavailable",
    )
