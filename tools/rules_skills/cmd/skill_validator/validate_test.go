package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validSkill contains a complete Markdown skill used as the validation baseline.
const validSkill = `---
name: example-skill
description: Validate an example skill.
---

# Example skill

Follow these instructions.
`

// validOpenAIYAML supplies valid optional metadata for the example skill.
const validOpenAIYAML = `interface:
  display_name: "Example Skill"
  short_description: "Validate an example repository skill"
  default_prompt: "Use $example-skill to validate this skill."
`

// writeTestSkill writes the requested fixtures under a temporary skill directory.
func writeTestSkill(t *testing.T, skill string, metadata string) (string, string) {
	t.Helper()
	skillDirectory := filepath.Join(t.TempDir(), "example-skill")
	if err := os.MkdirAll(filepath.Join(skillDirectory, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	skillPath := filepath.Join(skillDirectory, "SKILL.md")
	if err := os.WriteFile(skillPath, []byte(skill), 0o600); err != nil {
		t.Fatal(err)
	}
	metadataPath := filepath.Join(skillDirectory, "agents", "openai.yaml")
	if err := os.WriteFile(metadataPath, []byte(metadata), 0o600); err != nil {
		t.Fatal(err)
	}
	return skillPath, metadataPath
}

// TestValidSkillAndMetadata checks a complete skill with optional OpenAI metadata.
func TestValidSkillAndMetadata(t *testing.T) {
	skillPath, metadataPath := writeTestSkill(
		t,
		validSkill,
		validOpenAIYAML,
	)
	if errors := validateSkill(
		skillPath,
		metadataPath,
		"example-skill",
		[]string{skillPath, metadataPath},
	); len(errors) != 0 {
		t.Fatalf("validateSkill() errors = %q, want no errors", errors)
	}
}

// TestCompatibilityFrontmatterIsAllowed covers the optional compatibility field.
func TestCompatibilityFrontmatterIsAllowed(t *testing.T) {
	skill := `---
name: example-skill
description: Validate an example skill.
compatibility: Requires the example CLI.
---

# Example skill

Follow these instructions.
`
	skillPath, metadataPath := writeTestSkill(t, skill, validOpenAIYAML)
	if errors := validateSkill(
		skillPath,
		metadataPath,
		"example-skill",
		[]string{skillPath, metadataPath},
	); len(errors) != 0 {
		t.Fatalf("validateSkill() errors = %q, want no errors", errors)
	}
}

// TestNameMustMatchPackage checks that skill names match their owning directory.
func TestNameMustMatchPackage(t *testing.T) {
	skill := strings.Replace(
		validSkill,
		"name: example-skill",
		"name: another-skill",
		1,
	)
	skillPath, _ := writeTestSkill(t, skill, validOpenAIYAML)
	errors := validateSkill(skillPath, "", "example-skill", []string{skillPath})
	want := `skill name "another-skill" does not match package directory "example-skill"`
	if !containsError(errors, want) {
		t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
	}
}

// TestDefaultPromptMustNameSkill requires the prompt to mention its skill by name.
func TestDefaultPromptMustNameSkill(t *testing.T) {
	metadata := strings.Replace(
		validOpenAIYAML,
		"$example-skill",
		"the skill",
		1,
	)
	skillPath, metadataPath := writeTestSkill(t, validSkill, metadata)
	errors := validateSkill(
		skillPath,
		metadataPath,
		"example-skill",
		[]string{skillPath, metadataPath},
	)
	want := "interface.default_prompt must mention $example-skill"
	if !containsError(errors, want) {
		t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
	}
}

// TestDefaultPromptRequiresSkillNameBoundary rejects references to a longer skill name.
func TestDefaultPromptRequiresSkillNameBoundary(t *testing.T) {
	metadata := strings.Replace(
		validOpenAIYAML,
		"$example-skill",
		"$example-skill-extra",
		1,
	)
	skillPath, metadataPath := writeTestSkill(t, validSkill, metadata)
	errors := validateSkill(
		skillPath,
		metadataPath,
		"example-skill",
		[]string{skillPath, metadataPath},
	)
	want := "interface.default_prompt must mention $example-skill"
	if !containsError(errors, want) {
		t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
	}
}

// TestReferencedIconMustBeDeclared checks that existing icons are included in the declared inputs.
func TestReferencedIconMustBeDeclared(t *testing.T) {
	metadata := validOpenAIYAML + `  icon_small: "./assets/icon.svg"
`
	skillPath, metadataPath := writeTestSkill(t, validSkill, metadata)
	iconPath := filepath.Join(filepath.Dir(skillPath), "assets", "icon.svg")
	if err := os.MkdirAll(filepath.Dir(iconPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(iconPath, []byte("<svg/>\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	errors := validateSkill(
		skillPath,
		metadataPath,
		"example-skill",
		[]string{skillPath, metadataPath},
	)
	want := "interface.icon_small is not declared by skill_library"
	if !containsError(errors, want) {
		t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
	}
}

// TestUnfinishedTODOIsRejected covers an instruction placeholder outside a code fence.
func TestUnfinishedTODOIsRejected(t *testing.T) {
	skillPath, _ := writeTestSkill(
		t,
		validSkill+"\n[TODO: finish these instructions]\n",
		validOpenAIYAML,
	)
	errors := validateSkill(skillPath, "", "example-skill", []string{skillPath})
	want := "SKILL.md contains an unfinished TODO placeholder"
	if !containsError(errors, want) {
		t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
	}
}

// TestTODOInsideFenceIsAllowed accepts an example placeholder in fenced text.
func TestTODOInsideFenceIsAllowed(t *testing.T) {
	skillPath, _ := writeTestSkill(
		t,
		validSkill+"\n```text\n[TODO: example only]\n```\n",
		validOpenAIYAML,
	)
	if errors := validateSkill(
		skillPath,
		"",
		"example-skill",
		[]string{skillPath},
	); len(errors) != 0 {
		t.Fatalf("validateSkill() errors = %q, want no errors", errors)
	}
}

// TestLineEndingsAreNormalized checks that frontmatter and instructions behave consistently across line endings.
func TestLineEndingsAreNormalized(t *testing.T) {
	for _, test := range []struct {
		// name labels the line-ending convention.
		name string
		// newline separates the fixture lines.
		newline string
	}{
		{name: "CRLF", newline: "\r\n"},
		{name: "CR", newline: "\r"},
	} {
		t.Run(test.name, func(t *testing.T) {
			skill := strings.ReplaceAll(validSkill, "\n", test.newline)
			skillPath, _ := writeTestSkill(t, skill, validOpenAIYAML)
			if errors := validateSkill(
				skillPath,
				"",
				"example-skill",
				[]string{skillPath},
			); len(errors) != 0 {
				t.Fatalf("validateSkill() errors = %q, want no errors", errors)
			}

			skill += test.newline + "[TODO: finish instructions]" + test.newline
			skillPath, _ = writeTestSkill(t, skill, validOpenAIYAML)
			errors := validateSkill(
				skillPath,
				"",
				"example-skill",
				[]string{skillPath},
			)
			want := "SKILL.md contains an unfinished TODO placeholder"
			if !containsError(errors, want) {
				t.Fatalf("validateSkill() errors = %q, want %q", errors, want)
			}
		})
	}
}

// TestMalformedFrontmatterIsRejected exercises a YAML parser failure.
func TestMalformedFrontmatterIsRejected(t *testing.T) {
	skillPath, _ := writeTestSkill(t, "---\nname: [\n---\n", validOpenAIYAML)
	errors := validateSkill(skillPath, "", "example-skill", []string{skillPath})
	if len(errors) == 0 {
		t.Fatal("validateSkill() returned no errors for malformed YAML")
	}
}

// TestDuplicateFrontmatterKeyIsRejected checks duplicate mapping key diagnostics.
func TestDuplicateFrontmatterKeyIsRejected(t *testing.T) {
	skill := strings.Replace(
		validSkill,
		"name: example-skill",
		"name: example-skill\nname: duplicate",
		1,
	)
	skillPath, _ := writeTestSkill(t, skill, validOpenAIYAML)
	errors := validateSkill(skillPath, "", "example-skill", []string{skillPath})
	if len(errors) == 0 {
		t.Fatal("validateSkill() returned no errors for a duplicate YAML key")
	}
}

// TestRunWritesReport checks the report file created on a successful invocation.
func TestRunWritesReport(t *testing.T) {
	skillPath, metadataPath := writeTestSkill(
		t,
		validSkill,
		validOpenAIYAML,
	)
	outputPath := filepath.Join(t.TempDir(), "validation.txt")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(t.Context(),
		[]string{
			"--skill",
			skillPath,
			"--openai-yaml",
			metadataPath,
			"--expected-name",
			"example-skill",
			"--output",
			outputPath,
			"--file",
			skillPath,
			"--file",
			metadataPath,
		},
		&stdout,
		&stderr,
	)
	if exitCode != 0 {
		t.Fatalf("run() = %d, stderr = %q", exitCode, stderr.String())
	}
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("validation report was not written: %v", err)
	}
}

// TestRunDoesNotWriteReportAfterValidationFailure checks exit status and output for an invalid skill.
func TestRunDoesNotWriteReportAfterValidationFailure(t *testing.T) {
	skillPath, _ := writeTestSkill(t, validSkill, validOpenAIYAML)
	outputPath := filepath.Join(t.TempDir(), "validation.txt")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(t.Context(),
		[]string{
			"--skill",
			skillPath,
			"--expected-name",
			"another-skill",
			"--output",
			outputPath,
			"--file",
			skillPath,
		},
		&stdout,
		&stderr,
	)
	if exitCode != 1 {
		t.Fatalf("run() = %d, stderr = %q, want 1", exitCode, stderr.String())
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("validation report exists after failure or stat failed: %v", err)
	}
}

// TestRunReportsFirstMissingOptionDeterministically checks the command’s required-option diagnostic and exit code.
func TestRunReportsFirstMissingOptionDeterministically(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(t.Context(), nil, &stdout, &stderr)
	if exitCode != 2 {
		t.Fatalf("run() = %d, want 2", exitCode)
	}
	want := "ERROR: --skill is required\n"
	if stderr.String() != want {
		t.Fatalf("run() stderr = %q, want %q", stderr.String(), want)
	}
}

// containsError looks up one exact message in a list of validation diagnostics.
func containsError(errors []string, want string) bool {
	for _, validationError := range errors {
		if validationError == want {
			return true
		}
	}
	return false
}
