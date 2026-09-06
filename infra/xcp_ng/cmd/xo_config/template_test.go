package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testImage(t *testing.T) (string, string) {
	t.Helper()
	data := append([]byte{'Q', 'F', 'I', 0xfb}, bytes.Repeat([]byte{0}, 100)...)
	path := filepath.Join(t.TempDir(), "image.qcow2")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	return path, hex.EncodeToString(hash[:])
}

func TestTemplateChecksumBeforeMutation(t *testing.T) {
	path, _ := testImage(t)
	// A nil RPC client proves that even inventory requests cannot precede the
	// checksum gate. No callback may acquire an upload capability either.
	err := importTemplate(nil, templateOptions{Image: path, Digest: strings.Repeat("0", 64), SR: "sr", Base: "base", Name: "name"}, func(string, *os.File) (string, error) { t.Fatal("upload called"); return "", nil }, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "SHA256 mismatch") {
		t.Fatalf("expected checksum failure: %v", err)
	}
}

func TestTemplateUploadBoundaryAndRedaction(t *testing.T) {
	path, digest := testImage(t)
	f, err := checkedImage(path, digest)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	origin, _ := url.Parse("https://xo.test/")
	for _, location := range []string{"https://evil.test/capability", "http://xo.test/capability", "//evil.test/capability", "https://user@xo.test/capability"} {
		if _, err := uploadTemplateDisk(nil, origin, location, f); err == nil {
			t.Errorf("accepted foreign upload URL %s", location)
		}
	}
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != "POST" || req.Header.Get("Authorization") != "" || req.Header.Get("Cookie") != "" {
			t.Error("unexpected upload authentication")
		}
		if req.ContentLength <= 104 {
			t.Error("multipart length missing")
		}
		reader, err := req.MultipartReader()
		if err != nil {
			t.Fatal(err)
		}
		part, err := reader.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		if part.FormName() != "file" {
			t.Fatal("wrong multipart field")
		}
		data, err := io.ReadAll(part)
		if err != nil || len(data) != 104 {
			t.Fatal("incorrect image upload")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"error":{"message":"private-capability"}}`)), Header: http.Header{}, Request: req}, nil
	})}
	_, err = uploadTemplateDisk(client, origin, "/upload/private-capability", f)
	if err == nil || strings.Contains(err.Error(), "private-capability") {
		t.Fatalf("unsafe error: %v", err)
	}
}

func TestTemplateImportAndRepeat(t *testing.T) {
	path, digest := testImage(t)
	options := templateOptions{Image: path, Digest: digest, SR: "sr", Base: "base", Name: "Fedora44"}
	description := "Managed by infra/xcp_ng import-template; qcow2 sha256=" + digest
	objects := map[string]templateObject{"base": {ID: "base", Type: "VM-template", Pool: "pool"}, "sr": {ID: "sr", Type: "SR", Pool: "pool"}}
	mutations := 0
	r := fakeRPC(t, func(method string, raw json.RawMessage) (any, any) {
		var params map[string]any
		if err := json.Unmarshal(raw, &params); err != nil {
			t.Fatal(err)
		}
		if method == "xo.getAllObjects" {
			return objects, nil
		}
		mutations++
		switch method {
		case "vm.create":
			if params["bootAfterCreate"] != false || len(params["VDIs"].([]any)) != 0 || len(params["VIFs"].([]any)) != 0 {
				t.Fatal("template VM must remain diskless and halted")
			}
			objects["vm"] = templateObject{ID: "vm", Type: "VM", Name: options.Name, Description: description, PowerState: "Halted"}
			return "vm", nil
		case "disk.import":
			if params["type"] != "qcow2" {
				t.Fatal("wrong import format")
			}
			return map[string]string{"$sendTo": "/upload/private-capability"}, nil
		case "vm.attachDisk":
			if params["position"] != "0" || params["bootable"] != true || params["vdi"] != "vdi" {
				t.Fatal("wrong boot disk attachment")
			}
			objects["vbd"] = templateObject{ID: "vbd", Type: "VBD", VM: "vm", VDI: "vdi", Bootable: true, Position: "0"}
			vm := objects["vm"]
			vm.VBDs = []string{"vbd"}
			objects["vm"] = vm
		case "vm.convertToTemplate":
			vm := objects["vm"]
			vm.Type = "VM-template"
			objects["vm"] = vm
		default:
			t.Errorf("unexpected mutation %s", method)
		}
		return nil, nil
	})
	upload := func(location string, f *os.File) (string, error) {
		objects["vdi"] = templateObject{ID: "vdi", Type: "VDI", Name: options.Name + " boot", Description: description, SR: "sr"}
		return "vdi", nil
	}
	var out bytes.Buffer
	if err := importTemplate(r, options, upload, &out); err != nil {
		t.Fatal(err)
	}
	if mutations != 4 || strings.Contains(out.String(), "private-capability") {
		t.Fatal("unexpected mutation count or capability leak")
	}
	if err := importTemplate(r, options, upload, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if mutations != 4 {
		t.Fatal("repeat created duplicate resources")
	}
	vm := objects["vm"]
	vm.Type = "VM"
	objects["vm"] = vm
	if err := importTemplate(r, options, upload, &bytes.Buffer{}); err == nil {
		t.Fatal("incomplete template was silently reused")
	}
	if mutations != 4 {
		t.Fatal("incomplete template caused mutations")
	}
}

func TestTemplateRejectsOtherBootDevices(t *testing.T) {
	options := templateOptions{SR: "sr", Name: "Fedora44"}
	vm := templateObject{ID: "vm", Type: "VM-template", Description: "marker", PowerState: "Halted", VBDs: []string{"vbd"}}
	objects := map[string]templateObject{"vbd": {Type: "VBD", VM: "vm", VDI: "vdi", Bootable: true, Position: "1"}, "vdi": {Type: "VDI", Name: "Fedora44 boot", Description: "marker", SR: "sr"}}
	if err := verifyTemplate(objects, vm, options, "marker"); err == nil {
		t.Fatal("nonzero boot device accepted")
	}
}

func TestUploadErrorClassification(t *testing.T) {
	for _, test := range []struct{ message, want string }{
		{"Unknown disk type, expected vhd private-token", "unsupported disk format"},
		{"QCOW2 file had an invalid header magic private-token", "rejected image header"},
		{"multipart boundary private-token", "multipart framing"},
		{"SR_BACKEND_FAILURE private-token", "storage backend"},
		{"host /import_raw_vdi/?token=private-token", "host disk import"},
		{"private-token https://server/secret", "unclassified server error"},
	} {
		raw, err := json.Marshal(map[string]any{"code": -32000, "message": test.message})
		if err != nil {
			t.Fatal(err)
		}
		actual := classifyUploadError(raw)
		if !strings.Contains(actual, test.want) || !strings.Contains(actual, "-32000") || strings.Contains(actual, "private-token") || strings.Contains(actual, "https://") {
			t.Fatalf("unsafe or incorrect classification %q", actual)
		}
	}
}

func TestVHDValidation(t *testing.T) {
	for _, valid := range []bool{false, true} {
		data := make([]byte, 1024)
		if valid {
			copy(data[512:], "conectix")
		}
		path := filepath.Join(t.TempDir(), "image.vhd")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(data)
		f, err := checkedImage(path, hex.EncodeToString(hash[:]), "vhd")
		if valid {
			if err != nil {
				t.Fatal(err)
			}
			position, _ := f.Seek(0, io.SeekCurrent)
			f.Close()
			if position != 0 {
				t.Fatal("VHD not rewound for upload")
			}
		} else if err == nil {
			f.Close()
			t.Fatal("invalid VHD footer accepted")
		}
	}
	if _, err := checkedImage("unused", strings.Repeat("0", 64), "raw"); err == nil {
		t.Fatal("unsupported format accepted")
	}
}

func TestVHDImportFormatAndMarker(t *testing.T) {
	data := make([]byte, 1024)
	copy(data[512:], "conectix")
	path := filepath.Join(t.TempDir(), "image.vhd")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	digest := hex.EncodeToString(hash[:])
	marker := "Managed by infra/xcp_ng import-template; vhd sha256=" + digest
	objects := map[string]templateObject{"base": {ID: "base", Type: "VM-template", Pool: "pool"}, "sr": {ID: "sr", Type: "SR", Pool: "pool"}}
	sawImport := false
	r := fakeRPC(t, func(method string, raw json.RawMessage) (any, any) {
		var params map[string]any
		if err := json.Unmarshal(raw, &params); err != nil {
			t.Fatal(err)
		}
		switch method {
		case "xo.getAllObjects":
			return objects, nil
		case "vm.create":
			if params["name_description"] != marker {
				t.Error("missing VHD provenance marker")
			}
			objects["vm"] = templateObject{ID: "vm", Type: "VM", PowerState: "Halted", Description: marker}
			return "vm", nil
		case "disk.import":
			sawImport = true
			if params["type"] != "vhd" || params["description"] != marker {
				t.Error("incorrect VHD import parameters")
			}
			return nil, map[string]any{"code": -32000, "message": "stop before upload"}
		default:
			t.Errorf("unexpected method %s", method)
			return nil, nil
		}
	})
	err := importTemplate(r, templateOptions{Image: path, Digest: digest, Format: "vhd", SR: "sr", Base: "base", Name: "Fedora44"}, func(string, *os.File) (string, error) { t.Fatal("unexpected upload"); return "", nil }, &bytes.Buffer{})
	if err == nil || !sawImport {
		t.Fatalf("VHD import not reached: %v", err)
	}
}

func TestTemplateExplicitResume(t *testing.T) {
	for _, variant := range []string{"valid", "wrong_id", "attached_disk", "wrong_marker", "not_requested"} {
		t.Run(variant, func(t *testing.T) {
			path, digest := testImage(t)
			oldDigest := strings.Repeat("1", 64)
			oldMarker := "Managed by infra/xcp_ng import-template; qcow2 sha256=" + oldDigest
			marker := "Managed by infra/xcp_ng import-template; qcow2 sha256=" + digest
			options := templateOptions{Image: path, Digest: digest, SR: "sr", Base: "base", Name: "Fedora44", ResumeID: "vm", PreviousDigest: oldDigest, PreviousFormat: "qcow2"}
			vm := templateObject{ID: "vm", Type: "VM", Name: "Fedora44", Description: oldMarker, PowerState: "Halted", Pool: "pool"}
			switch variant {
			case "wrong_id":
				options.ResumeID = "other"
			case "attached_disk":
				vm.VBDs = []string{"unknown-disk"}
			case "wrong_marker":
				vm.Description = "unrelated"
			case "not_requested":
				options.ResumeID = ""
			}
			objects := map[string]templateObject{"vm": vm, "base": {ID: "base", Type: "VM-template", Pool: "pool"}, "sr": {ID: "sr", Type: "SR", Pool: "pool"}}
			mutations := 0
			r := fakeRPC(t, func(method string, raw json.RawMessage) (any, any) {
				var params map[string]any
				if err := json.Unmarshal(raw, &params); err != nil {
					t.Fatal(err)
				}
				if method == "xo.getAllObjects" {
					return objects, nil
				}
				mutations++
				switch method {
				case "vm.set":
					if params["id"] != "vm" || params["name_description"] != marker {
						t.Fatal("wrong resumed VM update")
					}
					vm := objects["vm"]
					vm.Description = marker
					objects["vm"] = vm
				case "disk.import":
					return map[string]string{"$sendTo": "/upload"}, nil
				case "vm.attachDisk":
					if params["vm"] != "vm" {
						t.Fatal("wrong VM attachment")
					}
					objects["vbd"] = templateObject{ID: "vbd", Type: "VBD", VM: "vm", VDI: "vdi", Bootable: true, Position: "0"}
					vm := objects["vm"]
					vm.VBDs = []string{"vbd"}
					objects["vm"] = vm
				case "vm.convertToTemplate":
					vm := objects["vm"]
					vm.Type = "VM-template"
					objects["vm"] = vm
				default:
					t.Errorf("unexpected mutation while resuming: %s", method)
				}
				return nil, nil
			})
			err := importTemplate(r, options, func(string, *os.File) (string, error) {
				objects["vdi"] = templateObject{ID: "vdi", Type: "VDI", Name: "Fedora44 boot", Description: marker, SR: "sr"}
				return "vdi", nil
			}, &bytes.Buffer{})
			if variant == "valid" {
				if err != nil || mutations != 4 {
					t.Fatalf("resume failed: mutations=%d error=%v", mutations, err)
				}
			} else if err == nil || mutations != 0 {
				t.Fatalf("unsafe resume: mutations=%d error=%v", mutations, err)
			}
		})
	}
}
