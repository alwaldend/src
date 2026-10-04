"""Verify the build-produced model with a real recommendation."""

import json
import os
import pathlib
import subprocess
import tempfile
import typing

import joblib
from python.runfiles import runfiles
import recommender


def main() -> None:
    resolver: runfiles.Runfiles | None = runfiles.Create()
    assert resolver is not None
    location: str | None = resolver.Rlocation(
        "_main/projects/mlops_2026/model.pkl"
    )
    assert location is not None
    artifact: pathlib.Path = pathlib.Path(location)
    assert artifact.stat().st_size > 0
    model: recommender.RecommenderModel = typing.cast(
        recommender.RecommenderModel, joblib.load(artifact)
    )
    title: str = str(model.data_encoded.iloc[0]["track_name"])
    recommendations: list[dict[str, str]] = model.recommend(title)
    assert recommendations and len(recommendations) <= 5
    recommendation: dict[str, str]
    for recommendation in recommendations:
        assert set(recommendation) == {"track_name", "artists", "album_name"}
    runner_location: str | None = resolver.Rlocation(
        "_main/projects/mlops_2026/internal/recommender/training/build_model"
    )
    trainer_location: str | None = resolver.Rlocation(
        "com_github_yandex_practicum_mlops_freetrack/train_and_save_model"
    )
    data_location: str | None = resolver.Rlocation(
        "com_github_yandex_practicum_mlops_freetrack/data.csv"
    )
    assert runner_location and trainer_location and data_location
    temporary: str
    with tempfile.TemporaryDirectory(
        dir=os.environ.get("TEST_TMPDIR")
    ) as temporary:
        directory: pathlib.Path = pathlib.Path(temporary)
        malformed: pathlib.Path = directory / "malformed.csv"
        malformed.write_text("wrong_column\ninvalid\n")
        failed_output: pathlib.Path = directory / "model.pkl"
        completed: subprocess.CompletedProcess[bytes] = subprocess.run(
            [
                runner_location,
                "--trainer",
                trainer_location,
                "--dataset",
                str(malformed),
                "--data",
                data_location,
                "--output",
                str(failed_output),
            ],
            cwd=directory,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            timeout=30,
        )
        assert completed.returncode != 0
        assert b"KeyError" in completed.stdout
        assert not failed_output.exists()
    output: pathlib.Path = pathlib.Path(
        os.environ["TEST_UNDECLARED_OUTPUTS_DIR"]
    )
    output.joinpath("model-e2e.json").write_text(
        json.dumps(
            {
                "model_bytes": artifact.stat().st_size,
                "recommendations": len(recommendations),
                "training_failure_without_artifact": "passed",
            }
        )
        + "\n"
    )


if __name__ == "__main__":
    main()
