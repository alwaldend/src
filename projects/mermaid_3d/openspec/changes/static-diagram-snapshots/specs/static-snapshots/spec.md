## Purpose

Provide themed static images of the existing 3D scenes for documentation and
ordinary image consumers while retaining links to interactive diagrams.

## ADDED Requirements

### Requirement: Exact-size themed images

The exporter SHALL encode the viewer canvas as PNG or WebP at explicit pixel
dimensions, using the requested light/dark appearance and configured 3D or front
camera. Defaults SHALL be 1280 by 960 pixels, light, and the configured 3D view.
Each dimension MUST be an integer from 128 through 8192 and the total MUST NOT
exceed 16,777,216 pixels. PNG SHALL preserve canvas pixels without loss.

#### Scenario: Export a diagram

- **WHEN** a supported self-contained viewer is exported with valid options
- **THEN** its scene and labels fit the exact requested dimensions with padding
- **AND** its embedded font and maintained palette are used without UI controls

#### Scenario: Reject invalid export options

- **WHEN** the format, dimensions, appearance, or camera is unsupported
- **THEN** export fails with a diagnostic and does not create an image

### Requirement: Self-contained offline input

Export SHALL use the pinned renderer environment and embedded inputs. Missing
fonts, external file or network dependencies, and readiness timeouts SHALL fail
explicitly without producing an image.

#### Scenario: Reject an incomplete viewer

- **WHEN** a viewer requires a remote font, image, or local dependency
- **THEN** export rejects that dependency and produces no image

#### Scenario: Await font and viewer readiness

- **WHEN** the viewer fails to become ready or its embedded font is absent
- **THEN** export fails within the configured timeout or with a font diagnostic

### Requirement: Static documentation preserves theme selection

Documentation image consumers SHALL support light/dark PNG or WebP siblings
and show static 3D scenes with links to their interactive views.

#### Scenario: Read a diagram without WebGL

- **WHEN** a reader opens documentation containing a 3D snapshot
- **THEN** the scene is visible as an ordinary image in the selected appearance
- **AND** its interactive view remains available through a link

### Requirement: Repeatable raster acceptance

Verification SHALL decode exported images independently, compare PNG pixels
with a freshly rendered canvas, assess WebP fidelity, and preserve images and
machine-readable results identifying the checked inputs.

#### Scenario: Verify framing and content

- **WHEN** raster acceptance runs on supported graph examples and terminal markers
- **THEN** it checks dimensions, appearance, fonts, nonblank content, label ink, and unclipped margins
- **AND** it checks exclusion of an injected full-screen UI overlay and preservation of a distant-camera scene
