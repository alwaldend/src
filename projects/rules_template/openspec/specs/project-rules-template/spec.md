# Rules template

## Purpose

Render Go text templates from declared data files through a command and a Bazel
toolchain rule. This source baseline was observed on 2026-09-08 at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`. The command operates on one template
and output per invocation; the rule's list attributes do not establish a
multi-template rendering guarantee.

Sources: [project description](../../../README.md),
[command flags](../../../main/go/cmd.go),
[templater](../../../main/go/templater.go),
[template functions](../../../main/go/templater_func_map.go),
and [Bazel action](../../../main/bzl/template_run_binary.bzl).

## Requirements

### Requirement: Load supported data formats into template context

The `template` command SHALL load data files in argument order into `DataFiles`,
decode `.toml`, `.json`, `.ndjson`, and `.yaml` data, and expose `.txt` files as
lines without structured decoding. `--extension` SHALL override extension
selection for the supplied data files.

#### Scenario: A caller renders a template using newline-delimited JSON

- **WHEN** a declared `.ndjson` input contains valid JSON on each data line
- **THEN** its template context contains a parsed sequence in `Data` and the input lines in `Lines`.

### Requirement: Report unsupported formats and rendering failures

The command SHALL return errors for unsupported data extensions, failed reads,
invalid structured data, invalid template syntax, template execution failures,
and failed output writes.

#### Scenario: A data file uses an unsupported extension

- **WHEN** a data path has an unrecognized extension and no supported override
- **THEN** rendering fails with an unsupported-extension diagnostic identifying the data path.

### Requirement: Expose the registered template function set

Template execution SHALL use the project's function map, including JSON
serialization, HTML escaping, path basename and dirname helpers, and first and
last element helpers.

#### Scenario: A template serializes structured data

- **WHEN** a valid template invokes `to_json` on JSON-serializable context data
- **THEN** the function returns the JSON representation for inclusion in the output.

### Requirement: Render declared Bazel inputs with the selected toolchain

`template_run_binary` SHALL invoke the registered templater with its `template`
subcommand, declare source and data files as inputs, declare output files, and
forward additional rule arguments.

#### Scenario: A target declares one template and one output

- **WHEN** Bazel executes the target with its registered template toolchain
- **THEN** the action passes the template, ordered data flags, and output path to the templater and exposes the declared output to downstream targets.
