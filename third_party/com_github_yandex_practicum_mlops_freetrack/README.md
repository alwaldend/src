---
title: Yandex Practicum MLOps freetrack
description: Pinned upstream course source and data
---

Source: [Yandex-Practicum/mlops-freetrack](https://github.com/Yandex-Practicum/mlops-freetrack).
The immutable revision and source archive integrity are owned by
`include.MODULE.bazel`. The CSVs and training Python files stay upstream-owned;
this package only provides Bazel integration. The inspected revision includes
no LICENSE file; no first-party license is asserted for these materials.

The upstream training executable expects data.csv and dataset.csv in its working
directory and writes model.pkl there. Use the project model build target for a
sandboxed action with a declared output rather than running it in source.

The maintained load_model_path.patch gives load_model an optional file path,
preserving its model.pkl default. This allows Bazel runfiles resolution without
changing the process working directory. Training and CSV bytes remain upstream.
