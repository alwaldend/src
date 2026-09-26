## Context

See [proposal.md](proposal.md) for motivation and
[the capability specification](specs/random-icon-generation/spec.md) for
the CLI contract and defaults. The project introduces a command and renderer
without changing existing application behavior.

The user confirmed configurable canvas dimensions, background, palette,
logical pixel size, density, seed, and square or circle shapes. Technical
decisions are needed for rasterization, partial cells, and reproducibility,
so this change includes a design artifact.

Example intended invocation:

```sh
icon_generator \
  --width 256 --height 256 \
  --background transparent \
  --colors '#ff0000,#00ff00,#0000ff' \
  --pixel-size 16 --pixel-shape circle \
  --density 0.4 --seed 42 \
  --output icon.png
```

## Goals / Non-Goals

**Goals:**

- Produce PNGs directly from an in-memory raster using Go's standard
  library and the repository's existing build toolchain.
- Separate CLI validation, deterministic rendering, and output handling so
  each failure has a clear command-level diagnostic.
- Keep palette colors exact and verify results by decoding generated PNGs.

**Non-Goals:**

- Antialiasing, gradients, arbitrary vector paths, or image import.
- Animation, additional output formats, a GUI, or an HTTP service.
- Automatic symmetry, random cell positions outside the grid, or a
  guarantee that different seeds always produce different images.

## Decisions

### Encode PNG directly with the standard library

Allocate an `image.NRGBA`, initialize it with the background, assign shape
pixels, and write it with `image/png.Encode`. Use opaque palette colors and
either opaque RGB or fully transparent background values. This avoids
premultiplied-alpha conversions and preserves the specified colors.

Research verified the relevant APIs in the local Go standard-library source
and the official documentation:

- [`image.NRGBA`](https://pkg.go.dev/image#NRGBA) provides a mutable pixel
  buffer and direct color assignment.
- [`png.Encode`](https://pkg.go.dev/image/png#Encode) writes an image to an
  `io.Writer`; its documentation includes a standard-library-only image
  generation example.
- [`math/rand/v2`](https://pkg.go.dev/math/rand/v2) provides seeded random
  sources.

FFmpeg was the user's fallback if direct generation required dependencies.
That condition does not apply. A third-party graphics library adds no
required capability for square fills, circle masks, and PNG encoding.

### Treat logical pixels as grid cells

Traverse cells from top to bottom and left to right, anchored at `(0, 0)`.
Use the same cell occupancy and color decisions for either shape. This
makes shape selection independent of random placement and color selection.
Clip partial edge cells against the requested canvas without resizing them.
Requiring divisible dimensions would unnecessarily restrict canvas sizes.

Squares fill their cell bounds. Circles test output-pixel centers against
the circle centered in the full cell. For local integer coordinates
`x, y` and cell size `s`, use the integer condition
`(2*x + 1 - s)^2 + (2*y + 1 - s)^2 <= s^2`, with sufficiently wide integer
arithmetic. This handles odd sizes and size one consistently without a
drawing dependency. Hard boundaries preserve palette colors; antialiasing
would introduce additional colors.

### Keep randomness explicit and reproducible

Use a local `math/rand/v2` generator, with `NewPCG(seed, 0)` as the fixed
mapping from the unsigned CLI seed. For every grid cell, consume an
occupancy draw and then a uniform palette-index draw, even when that cell
is empty. Render according to the occupancy result and selected shape.
Density zero and one must retain their exact endpoint behavior.

When `--seed` is omitted, obtain a seed from `crypto/rand`, fail clearly if
that fails, and report the effective seed on stderr. Track whether the seed
flag was supplied so explicit zero does not mean "randomize." A package
global generator would obscure state and make repeatability harder to
control.

Reproducibility is scoped to the same generator version and settings.
Compare decoded pixels in checks: Go's PNG encoder does not promise stable
encoded bytes across toolchain versions. Fixed seeds do not imply that
every distinct seed produces a distinct image.

### Validate before allocating or opening the output

Use the standard `flag` package and strict color parsing. The specification
owns defaults, accepted values, and limits. Its 4096-per-dimension limit
bounds the main four-byte-per-pixel buffer at 64 MiB; encoding and process
overhead require additional memory. Reject non-finite density values as
well as values outside the allowed range.

Create the output exclusively after validation and rendering. Refusing an
existing path prevents accidental replacement without adding an overwrite
mode. Check both encoding and close errors and remove an incomplete file
created by this invocation on failure. Do not create missing parent
directories or delete pre-existing destinations.

### Keep the project small and owned locally

Place the entry point at `cmd/icon_generator`, rendering and configuration
support at `internal`, and command-level acceptance checks at `test` under
`projects/icon_generator`. Use the existing root Go/Bazel setup; a separate
module and new dependencies are unnecessary. Add only packages needed by
the implementation.

Implementation must add the project README, build targets, project catalog
entry, and the normal landing-page/documentation integration. Register the
project's OpenSpec source target with the shared validation workflow when
the project's BUILD targets are introduced.

## Risks / Trade-offs

- Probability does not guarantee an exact filled-cell count: document
  density as a probability and test its endpoints deterministically.
- Small rasterized circles are coarse and may resemble squares: document
  hard pixel boundaries and inspect representative odd and even sizes.
- Partial circles at canvas edges may be visibly truncated: retain fixed
  grid geometry and include a non-divisible canvas in acceptance output.
- PNG byte comparisons can fail after toolchain upgrades without visible
  changes: compare decoded colors and dimensions for reproducibility.
- Large dimensions consume memory: enforce the documented allocation
  bounds before allocating the raster.

## Validation Approach

Before implementation, enumerate failing command scenarios and write the
end-to-end harness. Run the actual CLI and decode its PNG output with the
standard library. Cover both shapes, odd and even sizes, size one, partial
edge cells, opaque and transparent backgrounds, single and multiple palette
entries, density endpoints, seeded replay, fresh-seed replay, invalid
options, existing output preservation, and unavailable output directories.

Check palette membership, alpha values, geometry, dimensions, and replayed
pixels rather than a snapshot of encoded bytes. For geometry, use small
hand-derived examples so the check does not duplicate the renderer's
algorithm. Avoid probabilistic assertions that can fail on valid output.

Retain a representative square PNG, circle PNG, and transparent PNG with a
manifest of commands, seeds, and decoded-image checksums in the test's
declared artifact directory. A reviewer must be able to inspect the images
and rerun the commands.

## Migration Plan

This is an additive CLI with no existing users or persistent state to
migrate. Build and validate it through repository targets before normal
source delivery. Publishing release binaries or deploying website changes
is outside this implementation plan's execution authority.

## Implementation Evidence

On 2026-09-26, the implementation at staged tree
`885ebebfa85eb5b568615f5230464e0d41b5da2c` passed:

- `bazel_agent bazel test //:repo_quality_test
//projects/icon_generator:test
//infra/src/openspec/validation:projects_icon_generator_validate_test
//infra/dns:config_test`: all 31 tests passed.
- Semantic lint of the new project packages and affected project, OpenSpec,
  and documentation packages passed.
- `bazel_agent bazel build //projects/alwaldend.com:site` passed after fixing
  the existing `download-host-defaults-doc-reference` defect: its README
  referenced an unpackaged YAML file instead of the generated defaults page.
- The project index, landing page, README page, images taxonomy, and active
  status page were inspected in the generated HTML.
- Square, circle, and transparent PNG examples were visually inspected.
  Replaying the artifact manifest reproduced the captured outputs using the
  same executable. The README command examples also completed successfully.
- The Linux executable is statically linked. Its E2E commands passed with
  an empty `PATH`, and its imports require no external graphics dependencies.

The E2E suite retains its subprocess-only write-failure check, as confirmed
by the user. No host tool installation was performed. Final source and
publication state are bound by the repository delivery receipt after archive.
