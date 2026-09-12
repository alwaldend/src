package lint

import (
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestCheckSourceRejectsDuplicateJSONKeys(t *testing.T) {
	for _, test := range []struct {
		name    string
		input   string
		pointer string
	}{
		{name: "record", input: `{"records":{"service":{},"service":{}}}`, pointer: "/records/service"},
		{name: "field", input: `{"records":{"service":{"A":{"name":"first","name":"second"}}}}`, pointer: "/records/service/A/name"},
		{name: "escaped equivalent key", input: `{"records":{"service":{},"ser\u0076ice":{}}}`, pointer: "/records/service"},
		{name: "array object", input: `{"records":[{"key":1,"key":2}]}`, pointer: "/records/0/key"},
		{name: "pointer escaping", input: `{"a/b~c":1,"a/b~c":2}`, pointer: "/a~1b~0c"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := CheckSource("projects/service/dnsconfig.json", []byte(test.input))
			if err == nil || !strings.Contains(err.Error(), "projects/service/dnsconfig.json: duplicate JSON key at "+test.pointer) {
				t.Fatalf("CheckSource() = %v, want file and pointer %s", err, test.pointer)
			}
		})
	}
}

func TestCheckSourceAcceptsIndependentObjectsAndEmptyRecords(t *testing.T) {
	for _, input := range []string{
		`{"records":{}}`,
		`{"records":{"first":{"A":{"name":"service"}},"second":{"A":{"name":"service"}}}}`,
		`[{"name":1},{"name":2}]`,
	} {
		if err := CheckSource("dnsconfig.json", []byte(input)); err != nil {
			t.Errorf("CheckSource(%s) = %v", input, err)
		}
	}
}

func TestCheckSourceRejectsMalformedOrTrailingJSON(t *testing.T) {
	for _, input := range []string{``, `{"records":`, `{"records":{}} {}`, `{"records":{}} trailing`, `{"records":[1,]}`} {
		if err := CheckSource("dnsconfig.json", []byte(input)); err == nil {
			t.Errorf("CheckSource(%q) accepted invalid JSON", input)
		}
	}
}

func TestDiscoverSourcesIncludesNestedWorkspacesAndUnknownFiles(t *testing.T) {
	workspace := fstest.MapFS{
		"infra/service/dnsconfig.json":              {Data: []byte(`{"records":{}}`)},
		"projects/nested/MODULE.bazel":              {Data: []byte(`module(name = "nested")`)},
		"projects/nested/dnsconfig.json":            {Data: []byte(`{"records":{}}`)},
		"projects/new/deeper/dnsconfig.json":        {Data: []byte(`{"records":{}}`)},
		"out/task/dnsconfig.json":                   {Data: []byte(`invalid ignored output`)},
		"projects/nested/out/task/dnsconfig.json":   {Data: []byte(`invalid ignored output`)},
		".git/dnsconfig.json":                       {Data: []byte(`invalid ignored metadata`)},
		"projects/nested/.terraform/dnsconfig.json": {Data: []byte(`invalid ignored cache`)},
		"node_modules/dependency/dnsconfig.json":    {Data: []byte(`invalid ignored dependency`)},
		"bazel-out/dnsconfig.json":                  {Data: []byte(`invalid ignored output`)},
	}
	sources, err := DiscoverSources(workspace)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"infra/service/dnsconfig.json",
		"projects/nested/dnsconfig.json",
		"projects/new/deeper/dnsconfig.json",
	}
	if !reflect.DeepEqual(sources, want) {
		t.Fatalf("DiscoverSources() = %#v, want %#v", sources, want)
	}
}

func TestScanLoadsNestedUnknownAndEmptySources(t *testing.T) {
	workspace := fstest.MapFS{
		"infra/service/dnsconfig.json":       {Data: []byte(`{"records":{"service":{"A":{"name":"API","address":"192.0.2.1"},"AAAA":{"name":"api","address":"2001:db8::1"},"dsp":["global","dc1","global"]}}}`)},
		"projects/nested/MODULE.bazel":       {Data: []byte(`module(name = "nested")`)},
		"projects/nested/dnsconfig.json":     {Data: []byte(`{"records":{}}`)},
		"projects/new/deeper/dnsconfig.json": {Data: []byte(`{"records":{"new-record":{"CNAME":{"name":"New.Example.com.","target":"elsewhere.example.net."},"dsp":["all"]}}}`)},
		"out/task/dnsconfig.json":            {Data: []byte(`invalid ignored output`)},
	}
	report, err := Scan(workspace, "example.com")
	if err != nil {
		t.Fatal(err)
	}
	wantSources := []string{
		"infra/service/dnsconfig.json",
		"projects/nested/dnsconfig.json",
		"projects/new/deeper/dnsconfig.json",
	}
	if !reflect.DeepEqual(report.Sources, wantSources) {
		t.Fatalf("report sources = %v, want nested, unknown, and empty files %v", report.Sources, wantSources)
	}
	wantRecords := map[string]Record{
		"infra/service/dnsconfig.json/service/A": {
			Source: "infra/service/dnsconfig.json", Key: "service", Name: "api.example.com", Type: "A", Views: []string{"dc1", "global"},
		},
		"infra/service/dnsconfig.json/service/AAAA": {
			Source: "infra/service/dnsconfig.json", Key: "service", Name: "api.example.com", Type: "AAAA", Views: []string{"dc1", "global"},
		},
		"projects/new/deeper/dnsconfig.json/new-record/CNAME": {
			Source: "projects/new/deeper/dnsconfig.json", Key: "new-record", Name: "new.example.com", Type: "CNAME", Views: []string{"all"},
		},
	}
	gotRecords := map[string]Record{}
	for _, record := range report.Records {
		gotRecords[record.Source+"/"+record.Key+"/"+record.Type] = record
	}
	if len(report.Records) != len(wantRecords) || !reflect.DeepEqual(gotRecords, wantRecords) {
		t.Fatalf("report records = %#v, want %#v", report.Records, wantRecords)
	}
}

type recordingFS struct {
	fstest.MapFS
	reads []string
}

func (workspace *recordingFS) ReadFile(name string) ([]byte, error) {
	workspace.reads = append(workspace.reads, name)
	return fs.ReadFile(workspace.MapFS, name)
}

func TestScanReadsEveryFileEvenAfterFailures(t *testing.T) {
	workspace := &recordingFS{MapFS: fstest.MapFS{
		"first/dnsconfig.json":  {Data: []byte(`{"records":{"duplicate":{},"duplicate":{}}}`)},
		"second/dnsconfig.json": {Data: []byte(`{"records":{}}`)},
		"third/dnsconfig.json":  {Data: []byte(`{"records":{"third":{"A":{"name":"third"},"dsp":["all"]}}}`)},
	}}
	report, err := Scan(workspace, "example.com")
	if err == nil || !strings.Contains(err.Error(), "first/dnsconfig.json: duplicate JSON key") {
		t.Fatalf("duplicate source key was not rejected: %v", err)
	}
	wantSources := []string{"first/dnsconfig.json", "second/dnsconfig.json", "third/dnsconfig.json"}
	if !reflect.DeepEqual(workspace.reads, wantSources) {
		t.Fatalf("loaded files = %v, want every file %v", workspace.reads, wantSources)
	}
	if !reflect.DeepEqual(report.Sources, wantSources) {
		t.Fatalf("report sources = %v, want every file %v", report.Sources, wantSources)
	}
	if len(report.Records) != 1 || report.Records[0].Source != "third/dnsconfig.json" || report.Records[0].Name != "third.example.com" {
		t.Fatalf("valid records after a failed source were not scanned: %#v", report.Records)
	}
}

func TestScanCountsEmptySources(t *testing.T) {
	workspace := fstest.MapFS{
		"projects/empty/dnsconfig.json": {Data: []byte(`{"records":{}}`)},
	}
	report, err := Scan(workspace, "example.com")
	if err != nil {
		t.Fatalf("empty source rejected: %v", err)
	}
	if !reflect.DeepEqual(report.Sources, []string{"projects/empty/dnsconfig.json"}) || len(report.Records) != 0 {
		t.Fatalf("empty source was not counted: %#v", report)
	}
}
