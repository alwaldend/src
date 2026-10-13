"""Exercise training, artifact reload, and reproducibility through the CLI."""

import json
import math
import os
import pathlib
import pickle
import subprocess
import tempfile
import unittest

import numpy
import numpy.typing
import pandas
import sklearn.ensemble
from python.runfiles import runfiles


class TrainingArtifactsTest(unittest.TestCase):
    """Verify observable training behavior using the upstream CSV prefix."""

    def test_training_artifacts(self) -> None:
        """Train twice outside the source tree and retain verification evidence."""
        resolver: runfiles.Runfiles = runfiles.Create()
        trainer: str = resolver.Rlocation(os.environ["TRAINER_RLOCATION"])
        dataset: str = resolver.Rlocation(os.environ["DATASET_RLOCATION"])
        temporary: str
        with tempfile.TemporaryDirectory(
            dir=os.environ.get("TEST_TMPDIR")
        ) as temporary:
            directory: pathlib.Path = pathlib.Path(temporary)
            runs: list[dict[str, object]] = []
            predictions: list[numpy.typing.NDArray[numpy.float64]] = []
            index: int
            for index in range(2):
                output: pathlib.Path = directory / str(index)
                subprocess.run(
                    [
                        trainer,
                        "--dataset",
                        dataset,
                        "--output-dir",
                        str(output),
                        "--random-state",
                        "42",
                        "--test-size",
                        "0.2",
                        "--limit",
                        "99",
                    ],
                    cwd=directory,
                    check=True,
                    capture_output=True,
                    text=True,
                )
                model_path: pathlib.Path = output / "model.pkl"
                self.assertGreater(model_path.stat().st_size, 0)
                estimator: sklearn.ensemble.GradientBoostingRegressor = (
                    pickle.loads(model_path.read_bytes())
                )
                inputs: pandas.DataFrame = pandas.DataFrame(
                    {
                        "pickup_longitude": [-73.99981689453125],
                        "pickup_latitude": [40.73835372924805],
                        "dropoff_longitude": [-73.99951171875],
                        "dropoff_latitude": [40.72321701049805],
                        "passenger_count": [1],
                        "hour": [19],
                        "day_of_week": [3],
                    }
                )
                predictions.append(estimator.predict(inputs))
                metrics: dict[str, object] = json.loads(
                    (output / "metrics.json").read_text()
                )
                score: object = metrics["r2"]
                self.assertIsInstance(score, (int, float))
                self.assertTrue(math.isfinite(float(score)))
                self.assertEqual(
                    metrics["split"], {"test_size": 0.2, "random_state": 42}
                )
                self.assertEqual(metrics["dataset"], dataset)
                self.assertEqual(
                    metrics["model_parameters"], estimator.get_params()
                )
                self.assertEqual(metrics["rows_loaded"], 99)
                self.assertEqual(metrics["rows_used"], 95)
                runs.append(metrics)
            self.assertEqual(runs[0], runs[1])
            numpy.testing.assert_array_equal(predictions[0], predictions[1])
            self.assertTrue(numpy.isfinite(predictions[0]).all())
            report: pathlib.Path = (
                pathlib.Path(
                    os.environ.get("TEST_UNDECLARED_OUTPUTS_DIR", temporary)
                )
                / "verification.json"
            )
            report.write_text(
                json.dumps(
                    {
                        "metrics": runs[0],
                        "prediction": predictions[0].tolist(),
                        "reproducible": True,
                    },
                    indent=2,
                )
                + "\n"
            )


if __name__ == "__main__":
    unittest.main()
