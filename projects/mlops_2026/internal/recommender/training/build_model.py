"""Stage upstream training inputs in an isolated action directory."""

import argparse
import os
import pathlib
import shutil
import subprocess
import tempfile


def main() -> None:
    parser: argparse.ArgumentParser = argparse.ArgumentParser()
    parser.add_argument("--trainer", required=True)
    parser.add_argument("--dataset", required=True)
    parser.add_argument("--data", required=True)
    parser.add_argument("--output", required=True)
    args: argparse.Namespace = parser.parse_args()
    trainer: pathlib.Path = pathlib.Path(args.trainer).absolute()
    dataset: pathlib.Path = pathlib.Path(args.dataset).absolute()
    data: pathlib.Path = pathlib.Path(args.data).absolute()
    output: pathlib.Path = pathlib.Path(args.output).absolute()
    temporary: str
    with tempfile.TemporaryDirectory(dir=output.parent) as temporary:
        directory: pathlib.Path = pathlib.Path(temporary)
        shutil.copyfile(dataset, directory / "dataset.csv")
        shutil.copyfile(data, directory / "data.csv")
        environment: dict[str, str] = dict(os.environ)
        key: str
        for key in ("RUNFILES_DIR", "RUNFILES_MANIFEST_FILE", "JAVA_RUNFILES"):
            environment.pop(key, None)
        environment.update(
            {
                "OPENBLAS_NUM_THREADS": "1",
                "OMP_NUM_THREADS": "1",
                "MKL_NUM_THREADS": "1",
            }
        )
        subprocess.run(
            [str(trainer)], cwd=directory, env=environment, check=True
        )
        trained: pathlib.Path = directory / "model.pkl"
        if not trained.is_file() or trained.stat().st_size == 0:
            raise RuntimeError("Training produced no model.pkl")
        shutil.move(str(trained), output)


if __name__ == "__main__":
    main()
