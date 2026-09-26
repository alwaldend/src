## Purpose

Generate PNG icons from random arrangements of configurable colored shapes,
with control over the canvas, background, palette, density, and seed.

## ADDED Requirements

### Requirement: Generate a PNG with configurable canvas dimensions

The `icon_generator` CLI SHALL accept `--width`, `--height`, and `--output`
and write one PNG with exactly the requested dimensions. Width and height
SHALL be integers from 1 through 4096, measured in output image pixels.
Generation SHALL work without installed image-processing programs or
external runtime libraries.

#### Scenario: Generate a rectangular image

- **WHEN** the user requests width 120, height 80, and an available output
  path with otherwise valid options
- **THEN** the command exits successfully and writes a decodable 120 by 80
  PNG to that path
- **AND** no image-processing executable is required on `PATH`

### Requirement: Configure background and allowed shape colors

The CLI SHALL accept `--background` as `#RRGGBB` or `transparent`, and
`--colors` as a comma-separated nonempty list of `#RRGGBB` colors. Hex digits
SHALL be case-insensitive. An occupied cell SHALL choose uniformly from
the supplied palette entries and use that opaque color for its shape.
Unoccupied cells and pixels outside shapes SHALL retain the background.
Rendering SHALL preserve palette colors without blended edge colors.

#### Scenario: Use a single foreground color

- **WHEN** the user selects a white background and a palette containing
  only `#123abc`
- **THEN** every rendered shape pixel has opaque RGB color `#123abc`
- **AND** every remaining pixel is opaque white

#### Scenario: Preserve transparent background

- **WHEN** the user selects `--background transparent`
- **THEN** pixels not covered by a shape have alpha zero
- **AND** shape pixels have alpha 255

### Requirement: Configure grid size and shape

The CLI SHALL accept `--pixel-size` as an integer from 1 through 4096 and
`--pixel-shape` as `square` or `circle`. The grid SHALL start at the top-left
canvas corner with cells of the requested size. A square SHALL fill its
cell; a circle SHALL be centered in its cell with a diameter equal to that
size. A circle SHALL include an output pixel when that pixel's center is
on or inside the circle. Shapes SHALL have hard rasterized boundaries.
The rightmost and bottommost cells SHALL be clipped at the canvas bounds
without changing their original shape size or center.

#### Scenario: Render square cells

- **WHEN** the user requests a 32 by 16 canvas, size 16, square shapes,
  and density 1
- **THEN** two adjacent 16 by 16 squares cover the canvas

#### Scenario: Render a circle within a cell

- **WHEN** the user requests a 16 by 16 canvas, size 16, circle shapes,
  and density 1
- **THEN** the cell contains a centered circle of diameter 16
- **AND** the canvas corners retain the background

#### Scenario: Clip incomplete edge cells

- **WHEN** the user requests a 25 by 18 canvas with size 16
- **THEN** the grid contains two columns and two rows
- **AND** shapes in the last column and row are clipped to the 25 by 18
  output bounds without padding or shrinking the shapes

#### Scenario: Render the smallest logical pixel

- **WHEN** size is 1, density is 1, and the shape is either supported value
- **THEN** every output pixel receives a palette color

### Requirement: Configure occupancy density

The CLI SHALL accept a finite `--density` from 0 through 1 inclusive.
Each cell SHALL independently be occupied with that probability. Density
SHALL represent an occupancy probability, not an exact occupied-cell count
or the percentage of output pixels covered by circles.

#### Scenario: Generate only background

- **WHEN** density is 0
- **THEN** every output pixel retains the configured background

#### Scenario: Occupy every cell

- **WHEN** density is 1
- **THEN** every cell contains the selected shape
- **AND** circle boundaries still leave background visible outside circles

### Requirement: Reproduce a pattern from its seed

The CLI SHALL accept an optional unsigned 64-bit integer `--seed`, including
zero. For the same tool version, seed, and rendering settings, the decoded
output dimensions and pixel values SHALL be identical. If omitted, the
command SHALL obtain a fresh seed and report the effective seed on stderr
so the user can reproduce the result. Reproducibility SHALL concern image
content, without requiring byte-identical PNG encodings across toolchains.

#### Scenario: Replay an explicitly seeded image

- **WHEN** two invocations use the same explicit seed and rendering options
  but different available output paths
- **THEN** decoding both PNGs produces identical dimensions and pixels

#### Scenario: Replay an automatically seeded image

- **WHEN** an invocation omits the seed and the user repeats its rendering
  options with the reported seed
- **THEN** the repeated image has identical decoded dimensions and pixels

### Requirement: Provide usable defaults and help

The CLI SHALL default to width 256, height 256, background `#ffffff`,
palette `#000000`, pixel size 16, square shapes, density 0.4, a fresh seed,
and output `icon.png`. The CLI SHALL provide `--help` describing options,
defaults, accepted values, seed replay, and the distinction between output
pixels and logical grid cells. Help SHALL exit successfully without writing
an image.

#### Scenario: Generate an image with defaults

- **WHEN** the command is invoked without options and `icon.png` does not
  already exist
- **THEN** it writes a 256 by 256 PNG using black 16 by 16 square cells,
  a white background, and occupancy probability 0.4

#### Scenario: Inspect usage without generating an image

- **WHEN** the user invokes `--help`
- **THEN** the command describes all supported options and exits
  successfully without creating or changing an image file

### Requirement: Reject invalid input and report output failures

The CLI SHALL reject invalid dimensions, size, density, palette, background,
shape, seed, unknown flags, and unexpected positional arguments with a
nonzero exit status and a diagnostic on stderr. Input validation SHALL
complete before creating the output file. The command SHALL refuse to
overwrite an existing output file. File creation, encoding, and close
failures SHALL produce a nonzero exit status and a diagnostic; a failed
write SHALL remove only the incomplete output created by that invocation.

#### Scenario: Reject invalid rendering options without output

- **WHEN** a rendering option is invalid, such as a zero width, size 4097,
  density `NaN`, an empty palette, a malformed hex color, or shape `triangle`
- **THEN** the command reports the invalid option and exits unsuccessfully
- **AND** it does not create an output image

#### Scenario: Preserve an existing output file

- **WHEN** the requested output path already exists
- **THEN** the command exits unsuccessfully and reports the conflict
- **AND** the existing file remains unchanged

#### Scenario: Report an unavailable destination

- **WHEN** the output parent directory does not exist or cannot be written
- **THEN** the command exits unsuccessfully and identifies the output error
