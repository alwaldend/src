## Purpose

Let users enter an ordered sequence of panel prompts in separate text boxes
without managing delimiters or a separate editor.

## ADDED Requirements

### Requirement: Automatically growing prompt fields

The node SHALL display one multiline text box per prompt and one trailing empty
box. Filling the trailing box SHALL create another empty box without changing
keyboard focus. Empty boxes SHALL not generate images.

#### Scenario: Enter several panels

- **WHEN** the user enters two prompts in successive empty boxes
- **THEN** the node displays those prompts and a third empty box
- **AND** execution receives the two prompts in order

### Requirement: Lossless storage and compatibility

The node SHALL preserve multiline prompts through workflow save and reload,
including blank lines within a panel. Existing line and blank-line inputs SHALL
remain accepted. Invalid JSON input SHALL fail without echoing prompt content.

#### Scenario: Restore a saved sequence

- **WHEN** a workflow containing a multiline prompt with a blank line is reloaded
- **THEN** that prompt remains a single panel with its text unchanged

#### Scenario: Use an existing workflow

- **WHEN** a workflow uses the existing line or blank-line separator
- **THEN** its prompts appear as separate text boxes in the original order

#### Scenario: Reject malformed data

- **WHEN** JSON mode receives malformed JSON or values other than strings
- **THEN** execution reports an input error without disclosing the submitted text
