## Purpose

Enter an ordered sequence of panel prompts in ComfyUI's native multiline box,
with blank lines available for visual spacing.

## ADDED Requirements

### Requirement: Native multiline prompt input

The node SHALL use the standard multiline STRING widget without a custom browser
editor. Line mode SHALL interpret each nonempty line as a panel and ignore blank
and whitespace-only lines, preserving prompt and filename order.

#### Scenario: Visually separated prompts

- **WHEN** two prompt lines have any number of empty lines between them
- **THEN** execution SHALL produce two ordered prompt values and paired filenames

### Requirement: Saved-input compatibility

The node SHALL retain the existing line, blank-line and JSON separator modes.
Blank-line mode SHALL preserve multiline paragraphs within a panel. JSON mode
SHALL retain string-array contents, including blank lines within a panel.
Invalid input SHALL fail without echoing submitted prompt text.

#### Scenario: Load a previous JSON list

- **WHEN** an existing workflow uses the JSON separator
- **THEN** its stored list SHALL remain accepted without silent text conversion
