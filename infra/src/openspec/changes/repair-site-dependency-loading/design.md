## Context

The pinned rules_python 2.4.2 importer evaluates each package's resolution
markers but does not follow incoming dependency markers. The autoscroll lock
therefore exposes PyObjC source builds on Linux even though pynput requires
them only on macOS. Native genquery visits every select branch in its scope.

## Decisions

Proceed with uv's supported `environments` setting in Autoscroll. Resolve
`sys_platform == 'darwin'` and `sys_platform != 'darwin'` separately; their union
preserves all platforms. The regenerated lock gives PyObjC package records
Darwin-qualified resolution markers understood by the existing importer.
Remove the local rules_python patch and its isolated implementation fixtures.

Marker-qualified overrides correctly alter dependency edges but leave PyObjC
package records applicable on Linux, so they do not fix this import path.
Platform-specific requirements were also compiled successfully for Linux,
macOS, and Windows; they add export files and wheel-platform configuration.
The native environment partition is the smaller supported solution. The latest
inspected rules_python 2.4.5 still selects package markers without incoming-edge
reachability, so upgrading alone does not fix the original lock.

Preserve the 25 existing package identities and all distribution hashes.
Regeneration also selects tomli 2.5.0 for the existing PySide6 dependency on
Python 3.10; the user explicitly approved this addition. Verify Linux, macOS,
and Windows native import metadata at the repository Python version, stable
repeated locking, and the actual Autoscroll/site consumers. Implementation and
probe records are linked from the component-local change
`projects/autoscroll/openspec/changes/partition-platform-resolution` and ignored
`out/bazel_deps_analyzer/uv-platform-import/`.

Use Bazel 9.3.1 for Rust Git acquisition. A clean isolated test of the exact
locked repository attributes reproduces Bazel 8.7's broken cache reuse: the
second output base reports a successful fetch but its manifest and Rust source
still point into the retired first output base. Enabling external-file checks
does not repair that reuse. The same test under 9.3.1 retains readable manifest,
source, and license files with identical bytes and relative links.

Decision: proceed with the user-requested Bazel upgrade and remove every Anki
archive-specific projection, template, update rule, override, and formatting
exception. Keep Cargo.lock and rules_rust unchanged. The current rules_rust
0.75.0 still invokes the same built-in Git rule, so upgrading it alone would
not fix Bazel 8.7's behavior. Disabling shared repository-content caching avoids
the failure but discards caching for unrelated dependencies. The original
archive also omitted license targets outside the stripped crate subdirectory;
the normal Git recipe under Bazel 9 preserves the complete source tree.

Select bazel_jar_jar 0.1.15, an already-consumed transitive dependency, to replace
0.1.7's implicit JavaInfo usage with explicit rules_java imports. Reconcile
the Bazel 9 lockfile through the normal module workflow. Raw probe inputs,
cache paths, file comparisons, and logs are under
`out/bazel_deps_analyzer/anki-acquisition/`.

## Risks / Trade-offs

- Resolution partitions can change package selection across Python versions →
  retain existing package identities/hashes and review the generated markers.
- Later lock generation could lose the platform distinction → keep both
  complementary environments documented and verify actual native imports.
- Bazel major-version compatibility → validate the actual Rust crate, native
  site capture/build, both diagram pipelines, and affected nested workspaces.
