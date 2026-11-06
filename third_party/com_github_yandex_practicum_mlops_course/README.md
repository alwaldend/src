---
title: Yandex Practicum MLOps course
description: Pinned upstream course source archive
---

Source: [Yandex-Practicum/mlops-course](https://github.com/Yandex-Practicum/mlops-course).
The immutable commit archive and SRI integrity are owned by
`include.MODULE.bazel`. Upstream had no tags when this dependency was added.

`//third_party/com_github_yandex_practicum_mlops_course:source` exposes the
upstream README, MIT LICENSE, lesson code, and data as repository-internal
Bazel inputs. The package provides source integration; each lesson's environment
and dependencies remain upstream-owned.

License: [MIT](https://github.com/Yandex-Practicum/mlops-course/blob/74dc45062f43f9b4abfbca738ffa6b56db4d5daa/LICENSE).
The archive retains the upstream license and copyright notice.

`//third_party/com_github_yandex_practicum_mlops_course:taxi_data` exposes the
full `lesson-08-code-review/data/uber.csv` to the maintained code-review lesson.
The lesson uses the repository Python lock; the source archive pin is unchanged.
