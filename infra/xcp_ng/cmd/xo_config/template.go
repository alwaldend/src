package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type (
	templateOptions struct{ Image, Digest, SR, Base, Name, Format, ResumeID, PreviousDigest, PreviousFormat string }
	templateObject  struct {
		ID          string   `json:"id"`
		Type        string   `json:"type"`
		Name        string   `json:"name_label"`
		Description string   `json:"name_description"`
		PowerState  string   `json:"power_state"`
		Pool        string   `json:"$pool"`
		SR          string   `json:"$SR"`
		VBDs        []string `json:"$VBDs"`
		VIFs        []string `json:"VIFs"`
		VM          string   `json:"VM"`
		VDI         string   `json:"VDI"`
		CD          bool     `json:"is_cd_drive"`
		Bootable    bool     `json:"bootable"`
		Position    any      `json:"position"`
	}
)

func templateInventory(r *rpc) (map[string]templateObject, error) {
	var objects map[string]templateObject
	err := r.call("xo.getAllObjects", map[string]any{}, &objects)
	return objects, err
}

func checkedImage(path, digest string, formats ...string) (*os.File, error) {
	format := "qcow2"
	if len(formats) != 0 {
		format = formats[0]
	}
	if format != "qcow2" && format != "vhd" {
		return nil, errors.New("XO_TEMPLATE_FORMAT must be qcow2 or vhd")
	}
	decoded, err := hex.DecodeString(digest)
	if err != nil || len(decoded) != sha256.Size {
		return nil, errors.New("XO_TEMPLATE_SHA256 must be a SHA256 digest")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("opening template image failed")
	}
	valid := false
	defer func() {
		if !valid {
			f.Close()
		}
	}()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return nil, errors.New("template image must be a regular file")
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, errors.New("reading template image failed")
	}
	if !bytes.Equal(hash.Sum(nil), decoded) {
		return nil, errors.New("template image SHA256 mismatch; no resources created")
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("seeking template image failed")
	}
	if format == "qcow2" {
		var magic [4]byte
		if _, err := io.ReadFull(f, magic[:]); err != nil || magic != [4]byte{'Q', 'F', 'I', 0xfb} {
			return nil, errors.New("template image is not QCOW2")
		}
	} else {
		if info.Size() < 512 {
			return nil, errors.New("template image has no VHD footer")
		}
		if _, err := f.Seek(-512, io.SeekEnd); err != nil {
			return nil, errors.New("seeking VHD footer failed")
		}
		var footer [512]byte
		if _, err := io.ReadFull(f, footer[:]); err != nil || string(footer[:8]) != "conectix" {
			return nil, errors.New("template image has an invalid VHD footer signature")
		}
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("seeking template image failed")
	}
	valid = true
	return f, nil
}

func uploadURL(origin *url.URL, location string) (*url.URL, error) {
	u, err := origin.Parse(location)
	if err != nil || u.Scheme != "https" || u.Host != origin.Host || u.User != nil || u.Fragment != "" || location == "" {
		return nil, errors.New("XO returned an untrusted disk upload destination")
	}
	return u, nil
}

// The multipart body streams the verified file. Neither the upload capability
// nor the potentially secret-bearing server error body is exposed to callers.
func uploadTemplateDisk(client *http.Client, origin *url.URL, location string, f *os.File) (string, error) {
	u, err := uploadURL(origin, location)
	if err != nil {
		return "", err
	}
	var framing bytes.Buffer
	w := multipart.NewWriter(&framing)
	if _, err := w.CreateFormFile("file", "image"); err != nil {
		return "", errors.New("creating multipart upload failed")
	}
	header := bytes.Clone(framing.Bytes())
	framing.Reset()
	if err := w.Close(); err != nil {
		return "", errors.New("closing multipart framing failed")
	}
	info, err := f.Stat()
	if err != nil {
		return "", errors.New("reading image size failed")
	}
	req, err := http.NewRequest(http.MethodPost, u.String(), io.MultiReader(bytes.NewReader(header), f, bytes.NewReader(framing.Bytes())))
	if err != nil {
		return "", errors.New("creating disk upload request failed")
	}
	req.ContentLength = int64(len(header)+framing.Len()) + info.Size()
	req.Header.Set("Content-Type", w.FormDataContentType())
	response, err := client.Do(req)
	if err != nil {
		return "", errors.New("disk upload failed; inspect the marked incomplete template before retrying")
	}
	defer response.Body.Close()
	var result struct {
		Result string          `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	decodeErr := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result)
	if response.StatusCode != http.StatusOK || decodeErr != nil || (len(result.Error) > 0 && string(result.Error) != "null") || result.Result == "" {
		return "", fmt.Errorf("disk upload failed (HTTP %d; %s); inspect the marked incomplete template", response.StatusCode, classifyUploadError(result.Error))
	}
	return result.Result, nil
}

// Server errors can contain passwords, signed URLs, or request payloads. Emit
// only a numeric JSON-RPC code and a fixed classification, never source text.
func classifyUploadError(raw json.RawMessage) string {
	var rpcError struct {
		Code    *int64 `json:"code"`
		Message string `json:"message"`
	}
	if json.Unmarshal(raw, &rpcError) != nil {
		return "unclassified response; details suppressed"
	}
	message := strings.ToLower(rpcError.Message)
	classification := "unclassified server error"
	switch {
	case strings.Contains(message, "unknown disk type"), strings.Contains(message, "unsupported format"), strings.Contains(message, "format not supported"), strings.Contains(message, "invalid format"):
		classification = "server reported unsupported disk format"
	case strings.Contains(message, "invalid header"), strings.Contains(message, "file is too small"), strings.Contains(message, "invalid magic"):
		classification = "server rejected image header"
	case strings.Contains(message, "multipart"), strings.Contains(message, "boundary"), strings.Contains(message, "content-length"), strings.Contains(message, "stream ended unexpectedly"):
		classification = "server reported multipart framing or length error"
	case strings.Contains(message, "sr_backend_failure"), strings.Contains(message, "sr_not_available"), strings.Contains(message, "vdi_io_error"), strings.Contains(message, "no space left"):
		classification = "server reported storage backend failure"
	case strings.Contains(message, "permission_denied"), strings.Contains(message, "session_authentication_failed"), strings.Contains(message, "unauthorized"):
		classification = "server reported authorization failure"
	case strings.Contains(message, "import_raw_vdi"):
		classification = "server reported host disk import failure"
	}
	if rpcError.Code != nil {
		return fmt.Sprintf("JSON-RPC code %d; %s; details suppressed", *rpcError.Code, classification)
	}
	return classification + "; details suppressed"
}

func verifyTemplate(objects map[string]templateObject, vm templateObject, options templateOptions, description string) error {
	if vm.Type != "VM-template" || vm.Description != description || vm.PowerState != "Halted" || len(vm.VIFs) != 0 || len(vm.VBDs) != 1 {
		return errors.New("existing template is incomplete or does not match the managed image; inspect before retrying")
	}
	vbd, ok := objects[vm.VBDs[0]]
	if !ok || vbd.Type != "VBD" || vbd.VM != vm.ID || vbd.CD || !vbd.Bootable || fmt.Sprint(vbd.Position) != "0" {
		return errors.New("template must have exactly one bootable disk at device 0")
	}
	vdi, ok := objects[vbd.VDI]
	if !ok || vdi.Type != "VDI" || vdi.SR != options.SR || vdi.Description != description || vdi.Name != options.Name+" boot" {
		return errors.New("template boot disk does not match the requested SR and image marker")
	}
	return nil
}

func waitTemplateInventory(r *rpc, accept func(map[string]templateObject) bool) (map[string]templateObject, error) {
	for attempt := 0; attempt < 20; attempt++ {
		objects, err := templateInventory(r)
		if err != nil {
			return nil, err
		}
		if accept(objects) {
			return objects, nil
		}
		time.Sleep(250 * time.Millisecond)
	}
	return nil, errors.New("XO inventory did not reach the template postcondition; inspect the marked incomplete template")
}

func importTemplate(r *rpc, options templateOptions, upload func(string, *os.File) (string, error), out io.Writer) error {
	if options.Format == "" {
		options.Format = "qcow2"
	}
	if options.Name == "" || options.SR == "" || options.Base == "" {
		return errors.New("XO_TEMPLATE_NAME, XO_TEMPLATE_SR and XO_BASE_TEMPLATE_ID are required")
	}
	f, err := checkedImage(options.Image, options.Digest, options.Format)
	if err != nil {
		return err
	}
	defer f.Close()
	description := "Managed by infra/xcp_ng import-template; " + options.Format + " sha256=" + options.Digest
	objects, err := templateInventory(r)
	if err != nil {
		return err
	}
	var existing *templateObject
	for _, obj := range objects {
		if (obj.Type == "VM" || obj.Type == "VM-template") && obj.Name == options.Name {
			if existing != nil {
				return errors.New("multiple VMs/templates match XO_TEMPLATE_NAME")
			}
			copy := obj
			existing = &copy
		}
	}
	var vmID string
	if existing != nil {
		if err := verifyTemplate(objects, *existing, options, description); err == nil {
			return json.NewEncoder(out).Encode(map[string]string{"template_id": existing.ID, "sha256": options.Digest, "status": "verified existing"})
		}
		previousDigest, digestErr := hex.DecodeString(options.PreviousDigest)
		previousMarker := "Managed by infra/xcp_ng import-template; " + options.PreviousFormat + " sha256=" + options.PreviousDigest
		if options.ResumeID == "" || existing.ID != options.ResumeID || existing.Type != "VM" || existing.PowerState != "Halted" || len(existing.VBDs) != 0 || len(existing.VIFs) != 0 || existing.Description != previousMarker || digestErr != nil || len(previousDigest) != sha256.Size || (options.PreviousFormat != "qcow2" && options.PreviousFormat != "vhd") {
			return errors.New("incomplete template requires an exact resume ID and previous image marker on a halted VM with zero VBDs and VIFs")
		}
		vmID = existing.ID
	} else if options.ResumeID != "" {
		return errors.New("requested template resume ID does not match an existing named VM")
	}
	for _, obj := range objects {
		if obj.Type == "VDI" && obj.Name == options.Name+" boot" {
			return errors.New("a matching boot disk exists without a completed template; inspect before retrying")
		}
	}
	base, ok := objects[options.Base]
	if !ok || base.Type != "VM-template" {
		return errors.New("XO_BASE_TEMPLATE_ID is not an existing template")
	}
	sr, ok := objects[options.SR]
	if !ok || sr.Type != "SR" || sr.Pool == "" || sr.Pool != base.Pool {
		return errors.New("template base and destination SR must belong to the same pool")
	}
	if existing != nil && existing.Pool != sr.Pool {
		return errors.New("resumed template VM and destination SR must belong to the same pool")
	}
	for _, id := range base.VBDs {
		vbd, ok := objects[id]
		if !ok || !vbd.CD || vbd.VDI != "" {
			return errors.New("base template must have no disks or inserted media")
		}
	}
	// Create the marked shell first: any interrupted import leaves evidence that
	// blocks accidental duplicate disks on a subsequent invocation.
	resuming := vmID != ""
	if resuming {
		if err := r.call("vm.set", map[string]string{"id": vmID, "name_description": description}, nil); err != nil {
			return err
		}
	} else {
		if err := r.call("vm.create", map[string]any{"template": options.Base, "name_label": options.Name, "name_description": description, "CPUs": 2, "memory": int64(4 << 30), "VDIs": []any{}, "VIFs": []any{}, "bootAfterCreate": false}, &vmID); err != nil {
			return err
		}
	}
	if vmID == "" {
		return errors.New("vm.create returned no ID; inspect template inventory")
	}
	objects, err = waitTemplateInventory(r, func(objects map[string]templateObject) bool {
		vm, ok := objects[vmID]
		return ok && vm.Description == description
	})
	if err != nil {
		return err
	}
	vm := objects[vmID]
	if vm.Type != "VM" || vm.PowerState != "Halted" || vm.Description != description || len(vm.VIFs) != 0 {
		return errors.New("new template VM is not the expected halted diskless VM")
	}
	if resuming && len(vm.VBDs) != 0 {
		return errors.New("resumed VM acquired a VBD; refusing disk removal or import")
	}
	deletionErrorVerified := false
	for _, id := range vm.VBDs {
		vbd, ok := objects[id]
		if !ok || vbd.VM != vmID || !vbd.CD || vbd.VDI != "" {
			return errors.New("new template VM has an unexpected disk; no disk was removed")
		}
		deleteErr := r.call("vbd.delete", map[string]string{"id": id}, nil)
		// Some XO versions delete the empty VBD then fail resolving its absent
		// VDI. Accept only the verified deletion, and report this API defect.
		objects, err = waitTemplateInventory(r, func(objects map[string]templateObject) bool { _, present := objects[id]; return !present })
		if err != nil {
			return err
		}
		if deleteErr != nil {
			deletionErrorVerified = true
		}
	}
	var capability struct {
		SendTo string `json:"$sendTo"`
	}
	if err := r.call("disk.import", map[string]any{"sr": options.SR, "type": options.Format, "name": options.Name + " boot", "description": description}, &capability); err != nil {
		return err
	}
	vdiID, err := upload(capability.SendTo, f)
	if err != nil {
		return err
	}
	if err := r.call("vm.attachDisk", map[string]any{"vm": vmID, "vdi": vdiID, "position": "0", "mode": "RW", "bootable": true}, nil); err != nil {
		return err
	}
	if err := r.call("vm.convertToTemplate", map[string]string{"id": vmID}, nil); err != nil {
		return err
	}
	objects, err = waitTemplateInventory(r, func(objects map[string]templateObject) bool {
		return verifyTemplate(objects, objects[vmID], options, description) == nil
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(map[string]any{"template_id": vmID, "disk_id": vdiID, "sha256": options.Digest, "status": "created and verified", "empty_cd_delete_error_verified": deletionErrorVerified})
}

func importTemplateEnv(r *rpc, imagePath string, out io.Writer) error {
	origin, err := url.Parse(os.Getenv("XOA_URL"))
	if err != nil {
		return errors.New("invalid XOA_URL")
	}
	origin.Scheme = "https"
	origin.Path = "/"
	origin.RawQuery = ""
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: os.Getenv("XOA_INSECURE") == "true"}
	defer transport.CloseIdleConnections()
	client := &http.Client{Transport: transport, Timeout: 30 * time.Minute, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	options := templateOptions{Image: imagePath, Digest: os.Getenv("XO_TEMPLATE_SHA256"), SR: os.Getenv("XO_TEMPLATE_SR"), Base: os.Getenv("XO_BASE_TEMPLATE_ID"), Name: os.Getenv("XO_TEMPLATE_NAME"), Format: os.Getenv("XO_TEMPLATE_FORMAT"), ResumeID: os.Getenv("XO_TEMPLATE_RESUME_ID"), PreviousDigest: os.Getenv("XO_TEMPLATE_PREVIOUS_SHA256"), PreviousFormat: os.Getenv("XO_TEMPLATE_PREVIOUS_FORMAT")}
	return importTemplate(r, options, func(location string, f *os.File) (string, error) {
		return uploadTemplateDisk(client, origin, location, f)
	}, out)
}
