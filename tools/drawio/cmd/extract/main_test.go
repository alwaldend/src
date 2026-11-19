package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func archive(t *testing.T, webapp entry, payload string) []byte {
	t.Helper()
	root := webapp
	for _, name := range []string{"webapp", "main", "src", "drawio"} {
		root = entry{Files: map[string]entry{name: root}}
	}
	header, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	headerSize := (len(header)+3)/4*4 + 8
	b := make([]byte, 8+headerSize)
	binary.LittleEndian.PutUint32(b, 4)
	binary.LittleEndian.PutUint32(b[4:], uint32(headerSize))
	binary.LittleEndian.PutUint32(b[8:], uint32(headerSize-4))
	binary.LittleEndian.PutUint32(b[12:], uint32(len(header)))
	copy(b[16:], header)
	return append(b, []byte(payload)...)
}

func TestExtract(t *testing.T) {
	b := archive(t, entry{Files: map[string]entry{
		"index.html": {Size: 5, Offset: "0"},
		"js":         {Files: map[string]entry{"app.js": {Size: 3, Offset: "5"}}},
	}}, "helloapp")
	out := t.TempDir()
	if err := extract(bytes.NewReader(b), int64(len(b)), out); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{"index.html": "hello", "js/app.js": "app"} {
		got, err := os.ReadFile(filepath.Join(out, name))
		if err != nil || string(got) != want {
			t.Fatalf("%s: got %q, err %v", name, got, err)
		}
	}
}

func TestRejectUnsafeEntries(t *testing.T) {
	for name, e := range map[string]entry{
		"../escape":   {Size: 1, Offset: "0"},
		"/absolute":   {Size: 1, Offset: "0"},
		"back\\slash": {Size: 1, Offset: "0"},
		"link":        {Link: "elsewhere"},
		"unpacked":    {Unpacked: true},
		"negative":    {Size: -1, Offset: "0"},
		"overflow":    {Size: 1, Offset: "9223372036854775807"},
		"truncated":   {Size: 2, Offset: "0"},
	} {
		t.Run(name, func(t *testing.T) {
			b := archive(t, entry{Files: map[string]entry{name: e}}, "a")
			if err := extract(bytes.NewReader(b), int64(len(b)), t.TempDir()); err == nil {
				t.Fatal("accepted invalid entry")
			}
		})
	}
}

func TestRejectMalformedHeader(t *testing.T) {
	valid := archive(t, entry{Files: map[string]entry{}}, "")
	for _, cut := range []int{0, 15, 16, len(valid) - 1} {
		b := valid[:cut]
		if err := extract(bytes.NewReader(b), int64(len(b)), t.TempDir()); err == nil {
			t.Fatalf("accepted %d-byte archive", cut)
		}
	}
	b := append([]byte(nil), valid...)
	binary.LittleEndian.PutUint32(b[12:], 1<<30)
	if err := extract(bytes.NewReader(b), int64(len(b)), t.TempDir()); err == nil {
		t.Fatal("accepted oversized header")
	}
}
