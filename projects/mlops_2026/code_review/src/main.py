"""Train the taxi-fare model and save its estimator and evaluation metadata."""

import argparse
import json
import os
import pathlib
import pickle

import pandas
from python.runfiles import runfiles

from projects.mlops_2026.code_review.src import data, features, model


def main(
    path: pathlib.Path,
    output_dir: pathlib.Path,
    random_state: int,
    test_size: float,
    limit: int | None = None,
) -> None:
    """Train on the selected records and write model.pkl with metrics.json."""
    if not 0 < test_size < 1:
        raise ValueError("test_size must be between zero and one")
    if limit is not None and limit < 1:
        raise ValueError("limit must be positive")
    raw_data: pandas.DataFrame = data.load_data(path, limit)
    processed_data: pandas.DataFrame = features.add_time_features(raw_data)
    x_train: pandas.DataFrame
    x_test: pandas.DataFrame
    y_train: pandas.Series
    y_test: pandas.Series
    x_train, x_test, y_train, y_test = data.split_data(
        processed_data, test_size, random_state
    )
    if len(y_test) < 2:
        raise ValueError("R² requires at least two test records")
    model_: model.TaxiFareModel = model.TaxiFareModel(
        random_state=random_state
    )
    model_.fit(x_train, y_train)
    score: float = float(model_.model.score(x_test, y_test))
    metrics: dict[str, object] = {
        "r2": score,
        "dataset": str(path),
        "split": {"test_size": test_size, "random_state": random_state},
        "model_parameters": model_.model.get_params(),
        "limit": limit,
        "rows_loaded": len(raw_data),
        "rows_used": len(processed_data),
    }
    serialized_metrics: str = (
        json.dumps(metrics, indent=2, allow_nan=False) + "\n"
    )
    output_dir.mkdir(parents=True, exist_ok=True)
    (output_dir / "model.pkl").write_bytes(pickle.dumps(model_.model))
    (output_dir / "metrics.json").write_text(
        serialized_metrics, encoding="utf-8"
    )
    print(f"R²: {score:.2f}")


def _cli() -> None:
    """Resolve Bazel data and workspace-relative output paths for the CLI."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(
        description=__doc__
    )
    parser.add_argument("--dataset", required=True)
    parser.add_argument(
        "--output-dir",
        type=pathlib.Path,
        default=pathlib.Path("out/code_review"),
    )
    parser.add_argument("--random-state", type=int, default=42)
    parser.add_argument("--test-size", type=float, default=0.2)
    parser.add_argument(
        "--limit", type=int, help="Read only this many CSV records"
    )
    args: argparse.Namespace = parser.parse_args()
    dataset: pathlib.Path = pathlib.Path(args.dataset)
    if not dataset.is_file():
        resolver: runfiles.Runfiles = runfiles.Create()
        resolved: str | None = (
            resolver.Rlocation(args.dataset) if resolver else None
        )
        if not resolved or not pathlib.Path(resolved).is_file():
            raise FileNotFoundError(args.dataset)
        dataset = pathlib.Path(resolved)
    output: pathlib.Path = args.output_dir
    if not output.is_absolute() and "BUILD_WORKSPACE_DIRECTORY" in os.environ:
        output = pathlib.Path(os.environ["BUILD_WORKSPACE_DIRECTORY"]) / output
    main(dataset, output, args.random_state, args.test_size, args.limit)


if __name__ == "__main__":
    _cli()
