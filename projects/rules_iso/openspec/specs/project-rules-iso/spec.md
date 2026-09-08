# Rules ISO

## Purpose

Download declared ISO images and expose explicit runnable flash targets. This
source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`. It records rule behavior; a flash
target's existence does not authorize writing any device.

Sources: [project description](../../../README.md),
[module extension](../../../pkg/bzl/iso_extension.bzl),
[flash repository](../../../pkg/bzl/iso_flash_repo.bzl), and
[flash wrapper](../../../pkg/bzl/iso_flash_binary.bzl).

## Requirements

### Requirement: Validate image declarations before download

An image declaration SHALL provide at least one download URL and an integrity
attribute; a custom downloaded filename SHALL end with `.iso`.

#### Scenario: A declaration supplies a non-ISO filename

- **WHEN** `downloaded_file_path` is nonempty and does not end with `.iso`
- **THEN** extension evaluation fails with a filename diagnostic.

### Requirement: Expose separate download and flash repositories

For an image named `<name>`, the extension SHALL declare an `http_file`
repository named `<name>` and a companion `<name>_flash` repository containing
the public `iso_flash` executable target.

#### Scenario: A caller only builds the downloaded image

- **WHEN** the image download target is built without running the flash executable
- **THEN** Bazel downloads the declared image without invoking the device-writing command.

### Requirement: Require an explicit block device for flashing

The flash wrapper SHALL require a device argument, reject non-block-device
paths, and invoke `dd` only after those checks pass.

#### Scenario: The flash wrapper receives no device argument

- **WHEN** the flash executable is run with no first argument
- **THEN** it prints usage and exits with status 2 before invoking `dd`.

#### Scenario: The supplied path is not a block device

- **WHEN** the first argument names a regular file or nonexistent path
- **THEN** the wrapper reports that the path is not a block device and exits with status 1.
