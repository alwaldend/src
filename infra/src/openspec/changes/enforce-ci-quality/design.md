## Decisions

Use Gazelle's existing diff test and the repository lint profile. Share the
Gazelle configuration between generation and verification. Keep the root
Bazel server alive through all phases and report failures in the job summary.

Check the root workspace. Resolve imports owned by other Bazel packages
with Gazelle directives. Retain explicit targets whose inputs, platform
selection, or runtime dependencies are determined by Bazel.
