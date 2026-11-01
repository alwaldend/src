# Code screenshot

Capture a UTF-8 code file as a PNG using the alwaldend.com shared light/dark palette and typography. Source is inserted as literal text, never interpreted as HTML or executed. The command runs the pinned headless browser with page JavaScript disabled and network requests blocked.

```sh
al bazel run //tools/code_screenshot -- \
  --input tools/code_screenshot/cmd/code_screenshot/main.mjs \
  --output out/code_screenshot/example.png --theme dark
```

Paths are relative to the workspace when run through Bazel. Parent output directories are created. `--theme` accepts `light` (default) or `dark`; `--width` sets the canvas width (default 960); `--scale` sets pixel density (default 2, maximum 4). Line numbers identify original source lines. Long lines wrap within the canvas; continuation rows retain their original line number. Whitespace and tabs are preserved; known filename extensions receive monochrome syntax highlighting: muted comments and emphasized keywords. `--language NAME` overrides detection; unknown extensions stay plain. Multiline tokens retain their styling across numbered, wrapped lines. Images exceeding 30,000 pixels on either side or 100 million pixels are rejected. Reduce scale or use a smaller file in that case.

The command compiles the existing site application stylesheet and embeds pinned Liberation Mono. It adds no runtime network dependency and does not change the site theme. Chrome uses the repository's existing `--no-sandbox` launch convention; capture only files you intend to read.

Run `al bazel test //tools/code_screenshot/test:capture_test` for repeatable command-level checks. Bazel test outputs retain light and dark PNGs plus a JSON report.
