## Purpose

Produce portable code images with the same palette and typography as the repository website.

## ADDED Requirements

### Requirement: Literal themed PNG capture

The command SHALL capture UTF-8 code as literal text in a PNG using the shared site theme, supporting light and dark modes and preserving whitespace with line numbers and wrapping long lines within the selected width.

#### Scenario: Source contains markup

- **WHEN** source contains HTML or script text
- **THEN** it appears literally and is not executed

#### Scenario: Dark capture of a long file

- **WHEN** dark mode is selected and the source exceeds the canvas width or viewport height
- **THEN** the complete source appears in a dark themed PNG with numbered original lines and wrapped continuation rows

### Requirement: Clear command failures

The command SHALL reject invalid options and unreadable source before producing an output image.

#### Scenario: Invalid scale

- **WHEN** scale is nonpositive or not finite
- **THEN** the command exits unsuccessfully with a diagnostic

### Requirement: Theme-coherent highlighting

The command SHALL highlight known source languages using theme variables, preserve multiline token scope across numbered wrapped lines, and render unknown extensions as plain text. Source SHALL never execute.

#### Scenario: Multiline comment

- **WHEN** a comment spans source lines
- **THEN** every comment line is muted in both themes and retains its original line number

#### Scenario: Explicit language

- **WHEN** the user supplies --language
- **THEN** the selected supported language overrides filename detection, or an unsupported language fails clearly
