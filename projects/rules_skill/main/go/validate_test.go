package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validSkill = `---
name: example-skill
description: Validate an example skill.
---

# Example skill

Follow these instructions.
`

const validOpenAIYAML = `interface:
  display_name: "Example Skill"
  short_description: "Validate an example repository skill"
  default_prompt: "Use $example-skill to validate this skill."
`

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

func TestLineEndingsAreNormalized(t *testing.T) {
	for _, test := range []struct {
		name    string
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

func TestMalformedFrontmatterIsRejected(t *testing.T) {
	skillPath, _ := writeTestSkill(t, "---\nname: [\n---\n", validOpenAIYAML)
	errors := validateSkill(skillPath, "", "example-skill", []string{skillPath})
	if len(errors) == 0 {
		t.Fatal("validateSkill() returned no errors for malformed YAML")
	}
}

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

func TestRunWritesReport(t *testing.T) {
	skillPath, metadataPath := writeTestSkill(
		t,
		validSkill,
		validOpenAIYAML,
	)
	outputPath := filepath.Join(t.TempDir(), "validation.txt")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(
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

func TestRunDoesNotWriteReportAfterValidationFailure(t *testing.T) {
	skillPath, _ := writeTestSkill(t, validSkill, validOpenAIYAML)
	outputPath := filepath.Join(t.TempDir(), "validation.txt")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(
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

func TestRunReportsFirstMissingOptionDeterministically(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run(nil, &stdout, &stderr)
	if exitCode != 2 {
		t.Fatalf("run() = %d, want 2", exitCode)
	}
	want := "ERROR: --skill is required\n"
	if stderr.String() != want {
		t.Fatalf("run() stderr = %q, want %q", stderr.String(), want)
	}
}

func containsError(errors []string, want string) bool {
	for _, validationError := range errors {
		if validationError == want {
			return true
		}
	}
	return false
}
