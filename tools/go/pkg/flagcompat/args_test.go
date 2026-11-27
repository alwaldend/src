package flagcompat

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestSetArgs executes commands with both flag spellings and checks the values
// and positional arguments delivered to the operation.
func TestSetArgs(t *testing.T) {
	verify := func(t *testing.T, args []string, wantPath string, wantCheck bool, wantOperands ...string) {
		t.Helper()
		var path string
		var check bool
		var operands []string
		command := &cobra.Command{
			Use: "fixture",
			Run: func(cmd *cobra.Command, args []string) {
				operands = args
			},
		}
		command.Flags().StringVar(&path, "path", "", "Input path")
		command.Flags().BoolVar(&check, "check", false, "Check mode")
		original := slices.Clone(args)
		SetArgs(command, args)
		if err := command.ExecuteContext(t.Context()); err != nil {
			t.Fatal(err)
		}
		if path != wantPath || check != wantCheck || !slices.Equal(operands, wantOperands) {
			t.Fatalf("got path %q, check %t, operands %q", path, check, operands)
		}
		if !slices.Equal(args, original) {
			t.Fatalf("input arguments changed from %q to %q", original, args)
		}
	}

	t.Run("separate flag-like value", func(t *testing.T) {
		verify(t, []string{"-path", "-help", "-check"}, "-help", true)
	})
	t.Run("double-dash flag and flag-like value", func(t *testing.T) {
		verify(t, []string{"--path", "-check", "-check=false"}, "-check", false)
	})
	t.Run("attached value and terminator", func(t *testing.T) {
		verify(t, []string{"-path=-check", "--", "-check"}, "-check", false, "-check")
	})
	t.Run("terminator as value", func(t *testing.T) {
		verify(t, []string{"-path", "--", "-check"}, "--", true)
	})
	t.Run("first operand stops parsing", func(t *testing.T) {
		verify(t, []string{"-", "-check", "--path=other"}, "", false, "-", "-check", "--path=other")
	})
	t.Run("empty arguments", func(t *testing.T) {
		verify(t, nil, "", false)
	})
}

// TestSingleDashHelp verifies that Cobra displays help for Go's help spelling.
func TestSingleDashHelp(t *testing.T) {
	var output bytes.Buffer
	command := &cobra.Command{
		Use: "fixture",
		Run: func(*cobra.Command, []string) {
			t.Fatal("help executed the operation")
		},
	}
	command.SetOut(&output)
	SetArgs(command, []string{"-help"})
	if err := command.ExecuteContext(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Usage:") {
		t.Fatalf("missing usage in help output %q", output.String())
	}
}
