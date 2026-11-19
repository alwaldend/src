package main

import (
	"strings"
	"testing"
)

func TestToolSubcommandRemoved(t *testing.T) {
	_, err := bazelArguments([]string{"tool", "run", "repo_delivery"})
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("tool subcommand error = %v, want unknown command", err)
	}
}
