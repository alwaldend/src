package lint

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"testing/fstest"
)

func recordDocument(name, kind, view string) []byte {
	data, _ := json.Marshal(map[string]any{"records": map[string]any{
		"service": map[string]any{kind: map[string]string{"name": name}, "dsp": []string{view}},
	}})
	return data
}

func TestDifferentFilesCannotManageOneDomain(t *testing.T) {
	for _, test := range []struct {
		name     string
		first    string
		second   string
		typeName string
		view     string
	}{
		{name: "identical names", first: "service", second: "service", typeName: "A", view: "global"},
		{name: "different record types", first: "service", second: "service", typeName: "AAAA", view: "global"},
		{name: "different views", first: "service", second: "service", typeName: "A", view: "dc1"},
		{name: "different types and views", first: "service", second: "service", typeName: "CNAME", view: "dc1"},
		{name: "relative and absolute", first: "service", second: "service.example.com.", typeName: "TXT", view: "all"},
		{name: "mixed case", first: "SERVICE", second: "Service.Example.COM", typeName: "A", view: "global"},
		{name: "apex", first: "@", second: "EXAMPLE.COM.", typeName: "MX", view: "global"},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace := fstest.MapFS{
				"infra/first/dnsconfig.json":     {Data: recordDocument(test.first, "A", "global")},
				"projects/second/dnsconfig.json": {Data: recordDocument(test.second, test.typeName, test.view)},
			}
			_, err := Scan(workspace, "EXAMPLE.COM.")
			if err == nil {
				t.Fatal("duplicate domain across source files was accepted")
			}
			for _, expected := range []string{"is managed by multiple dnsconfig.json files", "infra/first/dnsconfig.json", "projects/second/dnsconfig.json", `record "service"`} {
				if !strings.Contains(err.Error(), expected) {
					t.Errorf("error %q does not identify %q", err, expected)
				}
			}
		})
	}
}

func TestSameFileCanManageSeveralValuesTypesAndViews(t *testing.T) {
	workspace := fstest.MapFS{
		"infra/service/dnsconfig.json": {Data: []byte(`{"records":{
			"first":{"A":{"name":"service","address":"192.0.2.1"},"AAAA":{"name":"service","address":"2001:db8::1"},"dsp":["global"]},
			"second":{"A":{"name":"SERVICE.EXAMPLE.COM.","address":"192.0.2.2"},"dsp":["dc1"]},
			"text":{"TXT":{"name":"service.example.com","content":"first"},"dsp":["all"]},
			"mail":{"MX":{"name":"service","priority":10,"target":"mail.example.com."},"dsp":["all"]}
		}}`)},
	}
	report, err := Scan(workspace, "example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Records) != 5 {
		t.Fatalf("record declarations = %d, want 5", len(report.Records))
	}
}

func TestDifferentFilesCanManageDistinctNamesInOneZone(t *testing.T) {
	workspace := fstest.MapFS{
		"infra/parent/dnsconfig.json": {Data: recordDocument("@", "A", "all")},
		"infra/child/dnsconfig.json":  {Data: recordDocument("service", "A", "all")},
		"infra/nested/dnsconfig.json": {Data: recordDocument("host.service", "A", "all")},
	}
	if _, err := Scan(workspace, "example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestScanRejectsInvalidNamesAndMetadata(t *testing.T) {
	for _, test := range []struct {
		name string
		data []byte
	}{
		{name: "absolute name outside zone", data: recordDocument("service.other.com.", "A", "global")},
		{name: "empty name", data: recordDocument("", "A", "global")},
		{name: "whitespace name", data: recordDocument("service example", "A", "global")},
		{name: "unsupported view", data: recordDocument("service", "A", "unknown")},
		{name: "missing records", data: []byte(`{}`)},
		{name: "missing views", data: []byte(`{"records":{"service":{"A":{"name":"service"}}}}`)},
		{name: "missing type", data: []byte(`{"records":{"service":{"dsp":["global"]}}}`)},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Scan(fstest.MapFS{"dnsconfig.json": {Data: test.data}}, "example.com")
			if err == nil || !strings.Contains(err.Error(), "dnsconfig.json") {
				t.Fatalf("invalid declaration accepted: %v", err)
			}
		})
	}
	if _, err := Scan(fstest.MapFS{}, "bad zone"); err == nil {
		t.Fatal("invalid DNS zone accepted")
	}
}

func TestWriteMarkdownGroupsAndSortsDomainsIncludingEmptyFiles(t *testing.T) {
	report := Report{
		Sources: []string{"projects/service/dnsconfig.json", "infra/empty/dnsconfig.json"},
		Records: []Record{
			{Source: "projects/service/dnsconfig.json", Name: "z.example.com", Type: "AAAA", Views: []string{"global"}},
			{Source: "projects/service/dnsconfig.json", Name: "a.example.com", Type: "TXT", Views: []string{"all"}},
			{Source: "projects/service/dnsconfig.json", Name: "z.example.com", Type: "A", Views: []string{"global", "dc1"}},
			{Source: "projects/service/dnsconfig.json", Name: "z.example.com", Type: "A", Views: []string{"dc1"}},
		},
	}
	var output bytes.Buffer
	if err := WriteMarkdown(&output, report); err != nil {
		t.Fatal(err)
	}
	want := "| Source | Domain | Types | Views |\n" +
		"| --- | --- | --- | --- |\n" +
		"| infra/empty/dnsconfig.json | - | - | - |\n" +
		"| projects/service/dnsconfig.json | a.example.com | TXT | all |\n" +
		"| projects/service/dnsconfig.json | z.example.com | A, AAAA | dc1, global |\n"
	if output.String() != want {
		t.Fatalf("table:\n%s\nwant:\n%s", output.String(), want)
	}
}
