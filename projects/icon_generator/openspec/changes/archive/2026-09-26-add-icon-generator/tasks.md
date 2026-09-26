## 1. Define acceptance before implementation

- [x] 1.1 Enumerate command failure cases from the specification and design
      before writing implementation code; verify the list covers invalid
      options, allocation limits, seed failures, and output failures.
- [x] 1.2 Write the end-to-end harness and hand-derived geometry cases
      before the CLI implementation; verify it invokes the real command,
      decodes PNGs, and defines a repeatable artifact manifest without
      duplicating the rendering algorithm.

## 2. Implement the command and renderer

- [x] 2.1 Add the Go entry point, internal packages, and Bazel targets using
      only the existing toolchain and standard library; verify the command
      target builds and the change introduces no external dependencies.
- [x] 2.2 Implement flags, defaults, color parsing, input bounds, and help;
      verify the prewritten invalid-input and help scenarios pass without
      creating image files.
- [x] 2.3 Implement seeded grid occupancy, palette selection, background
      fill, square and circle rasterization, and edge clipping; verify decoded
      output against the prewritten geometry, color, alpha, density-endpoint,
      and explicit-seed replay cases.
- [x] 2.4 Implement fresh-seed generation and reporting, exclusive output
      creation, PNG encoding, and error cleanup; verify automatic-seed replay,
      existing-file preservation, and output failure behavior.

## 3. Document and integrate the project

- [x] 3.1 Add the project README with runnable examples, option defaults,
      density semantics, clipping, and seed replay; verify every example
      against the built command and check documented values against help.
- [x] 3.2 Add the normal project catalog, landing-page, and documentation
      wiring; verify the affected documentation targets and inspect the
      generated project page without deploying it.
- [x] 3.3 Expose the project's OpenSpec source target and register it with
      the shared validation workflow; verify strict validation covers this
      workspace and its new capability delta.

## 4. Verify and deliver the implementation

- [x] 4.1 Run the complete end-to-end suite and retain square, circle, and
      transparent PNG examples with commands, seeds, and decoded checksums;
      inspect representative images and replay the artifact manifest.
- [x] 4.2 Run the project checks, affected semantic lint, repository quality
      gate, and strict OpenSpec validation against the final candidate; record
      exact commands, candidate identity, results, and any remaining failures.

## Delivery

After implementation acceptance, archive the change and verify the resulting
baseline specification and archive checks. Publication state is recorded by
the `repo-delivery` receipt, rather than an implementation checkbox.
