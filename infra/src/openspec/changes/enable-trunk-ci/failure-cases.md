# Failure cases

- Feature branches, tags and PR events must not schedule repository CI.
- A local action cannot load itself from an empty workspace; checkout must precede it.
- The removed AL wrapper must not intercept build/test invocations.
- A failed build or test must not prevent subsequent phases or workspaces.
- A missing Bazel executable must fail without reporting success.
- Invalid workspace boundaries must fail before starting commands.
- Excluded cache directories and repeated boundaries must not add workspaces.
