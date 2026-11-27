package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sort"
	"strings"

	"git.alwaldend.com/alwaldend/src/tools/go/pkg/flagcompat"
	"github.com/spf13/cobra"
)

// lockArtifact decodes one Maven artifact from a rules_jvm_external lock.
type lockArtifact struct {
	// Version is the version pinned by the lock.
	Version string `json:"version"`
	// Shasums records checksums by artifact classifier.
	Shasums map[string]string `json:"shasums"`
}

// lockFile decodes the artifact table shared with the Bazel build.
type lockFile struct {
	// Artifacts maps Maven coordinates to their locked entries.
	Artifacts map[string]lockArtifact `json:"artifacts"`
}

// gradleDep is one library that Gradle should depend on.
type gradleDep struct {
	// group identifies the Maven group.
	group string
	// name identifies the artifact within the group.
	name string
	// version is the locked Gradle dependency version.
	version string
	// sha256 is the checksum written into verification metadata when present.
	sha256 string
	// isAAR selects the Android archive extension for verification.
	isAAR bool
}

// main selects the Gradle output and reports conversion failures at the command boundary.
func main() {
	var format, out string
	command := &cobra.Command{
		Use: "gradle_lock_gen <maven_lock.json>", Args: cobra.ExactArgs(1), SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := convert(cmd.Context(), args[0], format, out); err != nil {
				return fmt.Errorf("convert %q: %w", args[0], err)
			}
			return nil
		},
	}
	command.Flags().SetInterspersed(false)
	command.Flags().StringVar(&format, "format", "", "output format: toml or verification-xml")
	command.Flags().StringVar(&out, "out", "", "output file path")
	for _, name := range []string{"format", "out"} {
		if err := command.MarkFlagRequired(name); err != nil {
			fmt.Fprintf(os.Stderr, "configure --%s: %v\n", name, err)
			os.Exit(1)
		}
	}
	flagcompat.SetArgs(command, os.Args[1:])
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	err := command.ExecuteContext(ctx)
	stop()
	if err != nil {
		os.Exit(1)
	}
}

// convert reads one bounded Maven lock and writes the selected Gradle input.
func convert(ctx context.Context, lockPath, format, out string) (result error) {
	const maxLockSize = 32 << 20
	input, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("open lock: %w", err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close lock: %w", err))
		}
	}()
	lockData, err := io.ReadAll(io.LimitReader(input, maxLockSize+1))
	if err != nil {
		return fmt.Errorf("read lock: %w", err)
	}
	if len(lockData) > maxLockSize {
		return fmt.Errorf("lock exceeds %d bytes", maxLockSize)
	}
	var lock lockFile
	if err := json.Unmarshal(lockData, &lock); err != nil {
		return fmt.Errorf("parse lock file: %w", err)
	}

	deps := extractDeps(&lock)
	if len(deps) == 0 {
		return fmt.Errorf("no buildable artifacts found in lock file")
	}

	var content string
	switch format {
	case "toml":
		content = renderVersionsToml(deps)
	case "verification-xml":
		content = renderVerificationXML(deps)
	default:
		return fmt.Errorf("unknown format %q (want toml or verification-xml)", format)
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("generate %q: %w", out, err)
	}
	if err := os.WriteFile(out, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %q: %w", out, err)
	}
	return nil
}

// extractDeps converts the lock file's artifact map into a flat set of
// Gradle dependencies, one per group:name. It prefers the AAR variant when
// both a plain-jar and an AAR key exist for the same coordinate, because
// that is what an Android Gradle build needs.
func extractDeps(lock *lockFile) map[string]gradleDep {
	deps := make(map[string]gradleDep, len(lock.Artifacts))

	for key, art := range lock.Artifacts {
		parts := strings.SplitN(key, ":", 4)
		if len(parts) < 2 {
			continue
		}
		group, name := parts[0], parts[1]
		ga := group + ":" + name

		if len(parts) >= 4 && parts[2] == "jar" && (parts[3] == "sources" || parts[3] == "javadoc") {
			continue
		}

		sha := art.Shasums["jar"]
		isAAR := len(parts) >= 3 && parts[2] == "aar"

		existing, ok := deps[ga]
		if !ok {
			deps[ga] = gradleDep{group: group, name: name, version: art.Version, sha256: sha, isAAR: isAAR}
		} else if isAAR && !existing.isAAR && sha != "" {
			existing.version = art.Version
			existing.sha256 = sha
			existing.isAAR = true
			deps[ga] = existing
		} else if sha != "" && existing.sha256 == "" {
			existing.version = art.Version
			existing.sha256 = sha
			existing.isAAR = isAAR
			deps[ga] = existing
		}
	}

	return deps
}

// gradleAlias forms a Gradle version-catalog alias from Maven coordinates.
func gradleAlias(group, name string) string {
	return strings.NewReplacer(".", "-", ":", "-").Replace(group + ":" + name)
}

// renderVersionsToml lists dependencies by coordinate in a Gradle version catalog.
func renderVersionsToml(deps map[string]gradleDep) string {
	keys := make([]string, 0, len(deps))
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	buf.WriteString("# Generated by //projects/android_launcher/cmd/gradle_lock_gen. Do not edit.\n")
	buf.WriteString("# Source: projects/android_launcher/maven_lock.json (rules_jvm_external)\n\n")
	buf.WriteString("[libraries]\n")
	for _, ga := range keys {
		d := deps[ga]
		alias := gradleAlias(d.group, d.name)
		fmt.Fprintf(&buf, "%s = { module = %q, version = %q }\n", alias, ga, d.version)
	}
	return buf.String()
}

// renderVerificationXML emits Gradle checksum metadata for artifacts with locked digests.
func renderVerificationXML(deps map[string]gradleDep) string {
	keys := make([]string, 0, len(deps))
	for k := range deps {
		if deps[k].sha256 != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var buf strings.Builder
	buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buf.WriteString("<verification-metadata version=\"1.3\" xmlns=\"https://schema.gradle.org/dependency-verification\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"https://schema.gradle.org/dependency-verification https://schema.gradle.org/dependency-verification/dependency-verification-1.3.xsd\">\n")
	buf.WriteString("  <configuration>\n")
	buf.WriteString("    <verify-metadata>false</verify-metadata>\n")
	buf.WriteString("    <verify-signatures>false</verify-signatures>\n")
	buf.WriteString("  </configuration>\n")
	buf.WriteString("  <components>\n")
	for _, ga := range keys {
		d := deps[ga]
		ext := "jar"
		if d.isAAR {
			ext = "aar"
		}
		filename := d.name + "-" + d.version + "." + ext
		fmt.Fprintf(&buf, "    <component group=%q name=%q version=%q>\n", d.group, d.name, d.version)
		fmt.Fprintf(&buf, "      <artifact name=%q>\n", filename)
		fmt.Fprintf(&buf, "        <sha256 value=%q/>\n", d.sha256)
		fmt.Fprintf(&buf, "      </artifact>\n")
		fmt.Fprintf(&buf, "    </component>\n")
	}
	buf.WriteString("  </components>\n")
	buf.WriteString("</verification-metadata>\n")
	return buf.String()
}
