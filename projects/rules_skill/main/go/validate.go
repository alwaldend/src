package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/goccy/go-yaml"
)

const maxSkillNameLength = 64

var (
	allowedSkillKeys = stringSet(
		"allowed-tools",
		"description",
		"license",
		"metadata",
		"name",
	)
	allowedOpenAIKeys    = stringSet("dependencies", "interface", "policy")
	allowedInterfaceKeys = stringSet(
		"brand_color",
		"default_prompt",
		"display_name",
		"icon_large",
		"icon_small",
		"short_description",
	)
	allowedToolKeys = stringSet(
		"description",
		"transport",
		"type",
		"url",
		"value",
	)
	frontmatterPattern = regexp.MustCompile(
		"(?s)^---\\n(.*?)\\n---(?:\\n|$)",
	)
	skillNamePattern  = regexp.MustCompile("^[a-z0-9-]+$")
	brandColorPattern = regexp.MustCompile("^#[0-9A-Fa-f]{6}$")
	fencePattern      = regexp.MustCompile(
		"^[ \\t]*(?:(?:[-+*]|\\d+[.)])[ \\t]+)?(`{3,}|~{3,})(.*)$",
	)
	todoPattern = regexp.MustCompile(
		"^[ ]{0,3}\\[TODO:[^\\n]*\\][ \\t]*$",
	)
)

func stringSet(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func unknownKeyErrors(
	value map[string]any,
	allowed map[string]struct{},
	context string,
) []string {
	var unknown []string
	for key := range value {
		if _, ok := allowed[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return []string{
		fmt.Sprintf(
			"%s has unexpected key(s): %s",
			context,
			strings.Join(unknown, ", "),
		),
	}
}

func readYAMLMapping(path string, context string) (map[string]any, []string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, []string{fmt.Sprintf("could not read %s: %v", context, err)}
	}
	var value map[string]any
	if err := yaml.Unmarshal(content, &value); err != nil {
		return nil, []string{
			fmt.Sprintf("%s contains invalid YAML: %v", context, err),
		}
	}
	if value == nil {
		return nil, []string{fmt.Sprintf("%s must be a YAML mapping", context)}
	}
	return value, nil
}

func validateSkill(
	skillPath string,
	openAIPath string,
	expectedName string,
	declaredFiles []string,
) []string {
	declared := make(map[string]struct{}, len(declaredFiles))
	for _, declaredFile := range declaredFiles {
		declared[filepath.Clean(declaredFile)] = struct{}{}
	}
	var errors []string
	if _, ok := declared[filepath.Clean(skillPath)]; !ok {
		errors = append(errors, "SKILL.md is not declared by skill_library")
	}
	if openAIPath != "" {
		if _, ok := declared[filepath.Clean(openAIPath)]; !ok {
			errors = append(
				errors,
				"agents/openai.yaml is not declared by skill_library",
			)
		}
	}

	skillName, markdownErrors := validateSkillMarkdown(skillPath, expectedName)
	errors = append(errors, markdownErrors...)
	if openAIPath == "" {
		return errors
	}
	if skillName == "" {
		return append(
			errors,
			"cannot validate agents/openai.yaml without a valid skill name",
		)
	}
	return append(
		errors,
		validateOpenAIYAML(
			openAIPath,
			skillName,
			filepath.Dir(skillPath),
			declared,
		)...,
	)
}

func validateSkillMarkdown(
	skillPath string,
	expectedName string,
) (string, []string) {
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return "", []string{fmt.Sprintf("could not read SKILL.md: %v", err)}
	}
	normalizedContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	normalizedContent = strings.ReplaceAll(normalizedContent, "\r", "\n")
	content = []byte(normalizedContent)
	match := frontmatterPattern.FindSubmatch(content)
	if match == nil {
		return "", []string{
			"SKILL.md has invalid or missing YAML frontmatter",
		}
	}

	var frontmatter map[string]any
	if err := yaml.Unmarshal(match[1], &frontmatter); err != nil {
		return "", []string{
			fmt.Sprintf("SKILL.md frontmatter contains invalid YAML: %v", err),
		}
	}
	if frontmatter == nil {
		return "", []string{"SKILL.md frontmatter must be a YAML mapping"}
	}

	errors := unknownKeyErrors(
		frontmatter,
		allowedSkillKeys,
		"SKILL.md frontmatter",
	)
	skillName, ok := frontmatter["name"].(string)
	skillName = strings.TrimSpace(skillName)
	if !ok || skillName == "" {
		errors = append(
			errors,
			"SKILL.md frontmatter name must be a non-empty string",
		)
		skillName = ""
	} else {
		errors = append(
			errors,
			validateSkillName(skillName, expectedName)...,
		)
	}

	description, ok := frontmatter["description"].(string)
	description = strings.TrimSpace(description)
	if !ok || description == "" {
		errors = append(
			errors,
			"SKILL.md frontmatter description must be a non-empty string",
		)
	} else {
		errors = append(errors, validateDescription(description)...)
	}

	body := strings.TrimSpace(string(content[len(match[0]):]))
	if body == "" {
		errors = append(errors, "SKILL.md instructions are empty")
	} else if hasTODOPlaceholder(body) {
		errors = append(
			errors,
			"SKILL.md contains an unfinished TODO placeholder",
		)
	}
	return skillName, errors
}

func validateSkillName(skillName string, expectedName string) []string {
	var errors []string
	if !skillNamePattern.MatchString(skillName) {
		errors = append(
			errors,
			"skill name must contain only lowercase letters, digits, and hyphens",
		)
	}
	if strings.HasPrefix(skillName, "-") ||
		strings.HasSuffix(skillName, "-") ||
		strings.Contains(skillName, "--") {
		errors = append(
			errors,
			"skill name cannot start or end with a hyphen or contain consecutive hyphens",
		)
	}
	if utf8.RuneCountInString(skillName) > maxSkillNameLength {
		errors = append(
			errors,
			fmt.Sprintf("skill name exceeds %d characters", maxSkillNameLength),
		)
	}
	if skillName != expectedName {
		errors = append(
			errors,
			fmt.Sprintf(
				"skill name %q does not match package directory %q",
				skillName,
				expectedName,
			),
		)
	}
	return errors
}

func validateDescription(description string) []string {
	var errors []string
	if strings.HasPrefix(description, "[TODO:") {
		errors = append(errors, "skill description contains a TODO placeholder")
	}
	if strings.ContainsAny(description, "<>") {
		errors = append(errors, "skill description cannot contain angle brackets")
	}
	if utf8.RuneCountInString(description) > 1024 {
		errors = append(errors, "skill description exceeds 1024 characters")
	}
	return errors
}

func hasTODOPlaceholder(body string) bool {
	var fenceMarker byte
	fenceLength := 0
	for _, line := range strings.Split(body, "\n") {
		fence := fencePattern.FindStringSubmatch(line)
		if fence != nil {
			marker := fence[1]
			if fenceMarker == 0 {
				fenceMarker = marker[0]
				fenceLength = len(marker)
			} else if marker[0] == fenceMarker &&
				len(marker) >= fenceLength &&
				strings.TrimSpace(fence[2]) == "" {
				fenceMarker = 0
				fenceLength = 0
			}
			continue
		}
		if fenceMarker == 0 && todoPattern.MatchString(line) {
			return true
		}
	}
	return false
}

func validateOpenAIYAML(
	openAIPath string,
	skillName string,
	skillDirectory string,
	declaredFiles map[string]struct{},
) []string {
	metadata, errors := readYAMLMapping(openAIPath, "agents/openai.yaml")
	if metadata == nil {
		return errors
	}
	errors = append(
		errors,
		unknownKeyErrors(metadata, allowedOpenAIKeys, "agents/openai.yaml")...,
	)
	if value, ok := metadata["interface"]; ok {
		errors = append(
			errors,
			validateInterface(
				value,
				skillName,
				skillDirectory,
				declaredFiles,
			)...,
		)
	}
	if value, ok := metadata["dependencies"]; ok {
		errors = append(errors, validateDependencies(value)...)
	}
	if value, ok := metadata["policy"]; ok {
		errors = append(errors, validatePolicy(value)...)
	}
	return errors
}

func validateInterface(
	value any,
	skillName string,
	skillDirectory string,
	declaredFiles map[string]struct{},
) []string {
	interfaceMap, ok := value.(map[string]any)
	if !ok {
		return []string{"interface must be a YAML mapping"}
	}
	errors := unknownKeyErrors(interfaceMap, allowedInterfaceKeys, "interface")

	if displayName, exists := interfaceMap["display_name"]; exists {
		if stringValue, ok := displayName.(string); !ok ||
			strings.TrimSpace(stringValue) == "" {
			errors = append(
				errors,
				"interface.display_name must be a non-empty string",
			)
		}
	}
	if description, exists := interfaceMap["short_description"]; exists {
		stringValue, ok := description.(string)
		length := utf8.RuneCountInString(stringValue)
		if !ok {
			errors = append(errors, "interface.short_description must be a string")
		} else if length < 25 || length > 64 {
			errors = append(
				errors,
				"interface.short_description must contain 25-64 characters",
			)
		}
	}
	if prompt, exists := interfaceMap["default_prompt"]; exists {
		stringValue, ok := prompt.(string)
		if !ok || strings.TrimSpace(stringValue) == "" {
			errors = append(
				errors,
				"interface.default_prompt must be a non-empty string",
			)
		} else if !mentionsSkill(stringValue, skillName) {
			errors = append(
				errors,
				fmt.Sprintf(
					"interface.default_prompt must mention $%s",
					skillName,
				),
			)
		}
	}
	if color, exists := interfaceMap["brand_color"]; exists {
		stringValue, ok := color.(string)
		if !ok || !brandColorPattern.MatchString(stringValue) {
			errors = append(
				errors,
				"interface.brand_color must be a six-digit hex color",
			)
		}
	}
	for _, iconField := range []string{"icon_small", "icon_large"} {
		if icon, exists := interfaceMap[iconField]; exists {
			errors = append(
				errors,
				validateIcon(
					skillDirectory,
					iconField,
					icon,
					declaredFiles,
				)...,
			)
		}
	}
	return errors
}

func mentionsSkill(prompt string, skillName string) bool {
	token := "$" + skillName
	searchFrom := 0
	for searchFrom < len(prompt) {
		index := strings.Index(prompt[searchFrom:], token)
		if index == -1 {
			return false
		}
		end := searchFrom + index + len(token)
		if end == len(prompt) {
			return true
		}
		next, _ := utf8.DecodeRuneInString(prompt[end:])
		if !unicode.IsLetter(next) &&
			!unicode.IsDigit(next) &&
			next != '-' &&
			next != '_' {
			return true
		}
		searchFrom = end
	}
	return false
}

func validateIcon(
	skillDirectory string,
	field string,
	value any,
	declaredFiles map[string]struct{},
) []string {
	iconPath, ok := value.(string)
	if !ok || iconPath == "" {
		return []string{
			fmt.Sprintf("interface.%s must be a non-empty string", field),
		}
	}
	if !strings.HasPrefix(iconPath, "./assets/") {
		return []string{
			fmt.Sprintf("interface.%s must point inside ./assets/", field),
		}
	}
	cleanPath := filepath.Clean(strings.TrimPrefix(iconPath, "./"))
	if filepath.IsAbs(cleanPath) ||
		!strings.HasPrefix(cleanPath, "assets"+string(filepath.Separator)) {
		return []string{
			fmt.Sprintf("interface.%s must be a safe relative path", field),
		}
	}
	fullPath := filepath.Clean(filepath.Join(skillDirectory, cleanPath))
	if _, ok := declaredFiles[fullPath]; !ok {
		return []string{
			fmt.Sprintf("interface.%s is not declared by skill_library", field),
		}
	}
	if stat, err := os.Stat(fullPath); err != nil || !stat.Mode().IsRegular() {
		return []string{
			fmt.Sprintf(
				"interface.%s does not reference an existing file",
				field,
			),
		}
	}
	return nil
}

func validateDependencies(value any) []string {
	dependencies, ok := value.(map[string]any)
	if !ok {
		return []string{"dependencies must be a YAML mapping"}
	}
	errors := unknownKeyErrors(dependencies, stringSet("tools"), "dependencies")
	tools, ok := dependencies["tools"].([]any)
	if !ok {
		return append(errors, "dependencies.tools must be a list")
	}
	for index, tool := range tools {
		toolMap, ok := tool.(map[string]any)
		context := fmt.Sprintf("dependencies.tools[%d]", index)
		if !ok {
			errors = append(errors, context+" must be a YAML mapping")
			continue
		}
		errors = append(
			errors,
			unknownKeyErrors(toolMap, allowedToolKeys, context)...,
		)
		if toolMap["type"] != "mcp" {
			errors = append(errors, context+".type must be 'mcp'")
		}
		dependencyName, ok := toolMap["value"].(string)
		if !ok || strings.TrimSpace(dependencyName) == "" {
			errors = append(
				errors,
				context+".value must be a non-empty string",
			)
		}
		for _, field := range []string{"description", "transport", "url"} {
			if fieldValue, exists := toolMap[field]; exists {
				if _, ok := fieldValue.(string); !ok {
					errors = append(errors, context+"."+field+" must be a string")
				}
			}
		}
	}
	return errors
}

func validatePolicy(value any) []string {
	policy, ok := value.(map[string]any)
	if !ok {
		return []string{"policy must be a YAML mapping"}
	}
	errors := unknownKeyErrors(
		policy,
		stringSet("allow_implicit_invocation"),
		"policy",
	)
	if implicit, exists := policy["allow_implicit_invocation"]; exists {
		if _, ok := implicit.(bool); !ok {
			errors = append(
				errors,
				"policy.allow_implicit_invocation must be a boolean",
			)
		}
	}
	return errors
}
