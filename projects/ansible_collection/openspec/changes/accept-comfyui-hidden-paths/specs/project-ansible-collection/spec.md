## ADDED Requirements

### Requirement: ComfyUI absolute application and state paths

The ComfyUI tools role SHALL accept normalized absolute application and state
paths with ordinary components containing dots.

#### Scenario: A caller uses a hidden parent directory

- **WHEN** the application or state directory is under an absolute path with
  a component such as `.cache`
- **THEN** that component passes path validation

#### Scenario: A caller uses a parent traversal component

- **WHEN** an application or state path contains a `..` component
- **THEN** the role rejects the path before applying changes
