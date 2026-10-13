# Predicting taxi fares

Train a gradient boosting regressor on the pinned upstream Uber CSV and save
`model.pkl` with `metrics.json`. This maintained adaptation applies the lesson's
review suggestions. The full dataset is fetched through Bazel; no CSV is checked in.

## Run

From the repository root, with the repository Bazel setup available:

```sh
bazel run //projects/mlops_2026/code_review:train
```

Outputs default to `out/code_review/` in this checkout. Bazel supplies the Python
runtime and dependencies from `tools/py/uv.lock`; a separate virtual environment
and the upstream's old requirements file are not needed.

For a quick run using the first 99 data rows (equivalent to the former 100-line CSV):

```sh
bazel run //projects/mlops_2026/code_review:train -- --limit 99 --output-dir out/code_review_sample
```

The CLI accepts `--dataset`, `--output-dir`, `--random-state` (default 42),
`--test-size` (default 0.2), and optional `--limit`. A supplied dataset is a local
CSV path with the upstream schema: unnamed index, `key`, `fare_amount`,
`pickup_datetime`, pickup/dropoff longitude and latitude, and `passenger_count`.
Use an absolute path for a custom dataset. Relative output directories are
resolved against the checkout when invoked with `bazel run`.

## Build artifacts

```sh
bazel build //projects/mlops_2026/code_review:model
```

This sandboxed, network-free action trains on the full upstream CSV and declares:

- `bazel-bin/projects/mlops_2026/code_review/model/model.pkl`
- `bazel-bin/projects/mlops_2026/code_review/model/metrics.json`

The target is marked manual because full-data training takes longer than the
sample test. Data fetching occurs through the existing integrity-pinned archive
before the action. The metric artifact records R², dataset path, split fraction
and seed, estimator parameters, optional row limit, and loaded/retained row counts.
The pickle contains the fitted scikit-learn estimator and expects the numeric
feature columns produced by `src/features.py` and `src/data.py`.

## Verify

```sh
bazel test //projects/mlops_2026/code_review/tests:training_test
```

The standard-library test invokes the real CLI twice from an unrelated working
directory on the first 99 upstream records. It reloads each model, checks finite
predictions and R², validates metadata, and checks reproducibility. The test
retains `verification.json` in Bazel's undeclared test outputs.

## Data handling

Training removes rows missing required values and rows whose pickup or dropoff
coordinate pair is `(0, 0)`. It extracts UTC hour and day of week, removes the
original timestamp, and excludes `key`, the CSV index, and `fare_amount` from
features. The fare remains the target. Both the split and estimator receive the
same seed. This exercise does not implement comprehensive anomaly validation
or hyperparameter search; its R² is a holdout score, not a quality guarantee.

## Source

Adapted from `lesson-08-code-review` in the existing pinned
[Yandex Practicum MLOps course](https://github.com/Yandex-Practicum/mlops-course).
The archive revision and integrity are owned by
`third_party/com_github_yandex_practicum_mlops_course/include.MODULE.bazel`.
The upstream MIT notice is retained in
[the change record](../openspec/changes/copy-code-review-lesson/upstream.LICENSE.txt).
