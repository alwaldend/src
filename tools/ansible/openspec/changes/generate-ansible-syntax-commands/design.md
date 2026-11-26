## Context

A package can declare one default playbook command, named commands, or both.
Each command already owns the arguments needed to select its playbook and inventory.

## Decisions

Generate a native Bazel test for each entry in the argument map. Each test
invokes the packaged executable directly with `--syntax-check` and the command
arguments, including packaged extra-vars.

## Migration Plan

Remove the standalone declarations once the macro supplies their replacements.
For the appliance package, the existing multi-playbook check maps to one
generated check per playbook command.
