## Decisions

Use the existing rules_binary_toolchain extension and the upstream go-jsonnet
v0.22.0 Linux amd64 archive. The same archive provides jsonnet and jsonnetfmt.
Publisher SHA-256: e87a93ea44e34da92c15636205c7f2240bf3ac92d00ccb855dbb4e7e03ea6941.
The user approved this exact release after a separate dependency approval question.
Other external software remains outside the authorized scope.

Add a small Jsonnet rule in tools/jsonnet. The caller declares the entrypoint,
all imported files and each named JSON output. Use a sandboxed action, no shell
script or runtime fetch, and existing execution-tool dependency conventions.
Outputs retain consumer-owned package-relative paths. Do not install the tools
or pass them to a rented instance. Keep workflow composition with ComfyUI tools.

## Acceptance and failure cases

Check archive integrity and both version commands; import files must be declared;
output paths must be declared and remain within the output directory. Exercise
actual compilation and formatting with an isolated fixture and preserve a result
artifact. A generation action must fail on invalid Jsonnet or an undefined symbol.
No network is allowed in build actions. Regenerate MODULE.bazel with its owning
workflow and verify stable generation; update the module lock through Bazel only.

The compiler and formatter version commands both report v0.22.0. Upstream
archive bytes match the publisher checksum. Actual generation preserves all 14
parsed workflow values from the recorded baseline. Isolated compiler fixtures
reject syntax/undefined-symbol errors; public configuration fixtures preserve
empty manifests, additional files and an omitted optional checkpoint kind.
No new runtime package, shared-host installation or live operation was used.
