---
title: Repository video evaluations
---

`eval_config_test` validates the configuration, cases, references, and staged
skill without model calls. Cases cover compatible delivery, real encoder
capability, complete audio, precise cuts, and interpolation routing.

A live model target is omitted because representative execution requires media
fixtures, writable artifact storage, working FFmpeg codecs, and browser or GPU
capabilities not supplied by this offline harness. Configuration validation
does not establish correct model answers or working media output. Actual export,
browser playback, and audio verification evidence are retained separately in
the task's ignored output directory.
