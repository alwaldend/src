package worker_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
	lua "github.com/yuin/gopher-lua"
)

// luaGenerator is the built schema-to-Lua command supplied by Bazel.
var luaGenerator = flag.String("lua-generator", "", "Lua declaration generator runfile")

// luaTypes is the generated, checked-in projection used by editors.
var luaTypes = flag.String("lua-types", "", "Lua plan declaration runfile")

// TestLuaAPIGeneration checks the command, maintained projection and reusable output.
// Failure cases are stale/missing declarations, unstable output, invalid Lua,
// failed output writes and a changed protocol producing stale editor types.
func TestLuaAPIGeneration(t *testing.T) {
	binary, err := runfiles.Rlocation(*luaGenerator)
	if err != nil {
		t.Fatal(err)
	}
	source, err := runfiles.Rlocation(*luaTypes)
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "plan.lua")
	command := exec.CommandContext(context.Background(), binary, "--out", output)
	if data, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generate Lua declarations: %v: %s", err, data)
	}
	generated, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	checked, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(generated, checked) {
		t.Fatal("editor plan declarations are stale; run generate_lua_api")
	}
	second, err := exec.CommandContext(context.Background(), binary).Output()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(generated, second) {
		t.Fatal("stdout and file generation differ")
	}
	state := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer state.Close()
	if _, err := state.Load(bytes.NewReader(generated), "plan.lua"); err != nil {
		t.Fatalf("generated declaration syntax: %v", err)
	}
	if data, err := exec.CommandContext(context.Background(), binary, "--out", filepath.Join(output, "missing.lua")).CombinedOutput(); err == nil {
		t.Fatalf("invalid output path succeeded: %s", data)
	}
	after, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(after, generated) {
		t.Fatal("failed generation changed existing output")
	}
	// This contract checks public editor capabilities, independently of generator internals.
	for _, definition := range []string{"WorkerDeletePostOptions", "WorkerUnrepostOptions", "WorkerUpdateProfile", "WorkerTimestamp", "ArchiveTweet"} {
		if !strings.Contains(string(generated), "---@class (exact) "+definition+"\n") {
			t.Fatalf("missing editor type %s", definition)
		}
	}
	for _, operation := range []string{"delete_post", "unrepost", "unlike", "unfollow", "create_post", "read_post", "update_profile"} {
		if !strings.Contains(string(generated), "function plan."+operation+"(") {
			t.Fatalf("missing plan method declaration %s", operation)
		}
	}
	if strings.Contains(string(generated), "WorkerPlan") || strings.Contains(string(generated), "WorkerOperation") {
		t.Fatal("obsolete returned-plan declarations remain")
	}
	digest := sha256.Sum256(generated)
	writeAsyncArtifact(t, "lua-api.json", map[string]any{"generated_sha256": hex.EncodeToString(digest[:]), "source_fresh": true, "stable_output": true, "lua_syntax_valid": true, "failed_write_preserved_source": true, "luals_diagnostics_run": false})
}
