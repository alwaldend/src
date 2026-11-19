package dns_records_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

func TestTargetedImportSkipsUnrelatedProviderAuthentication(t *testing.T) {
	for _, fixture := range []struct {
		name      string
		seedState bool
	}{
		{"empty_state", false},
		{"existing_service_state", true},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			targetedImportFixture(t, fixture.seedState)
		})
	}
}

func targetedImportFixture(t *testing.T, seedState bool) {
	t.Helper()
	workspace := t.TempDir()

	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		http.Error(w, "service API must not be contacted", http.StatusServiceUnavailable)
	}))
	defer server.Close()
	defer func() {
		if got := requests.Load(); got != 0 {
			t.Errorf("unrelated Proxmox API received %d requests, want zero", got)
		}
	}()

	writeFixture := func(name, content string) {
		t.Helper()
		path := filepath.Join(workspace, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	writeFixture("main.tf", `
terraform {
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = "= 3.0.2-rc07"
    }
  }
}

variable "dns_import_ids" {
  type    = map(string)
  default = {}
}

provider "proxmox" {
  pm_minimum_permission_check = false
}

resource "proxmox_pool" "service" {
  poolid = "service-fixture"
}

data "proxmox_ha_groups" "service" {
  group_name = "service-fixture"
}

resource "proxmox_pool" "discovered" {
  for_each = toset(data.proxmox_ha_groups.service.nodes)
  poolid   = each.key
}

import {
  for_each = toset(data.proxmox_ha_groups.service.nodes)
  to       = proxmox_pool.discovered[each.key]
  id       = "pools/${each.key}"
}

module "dns" {
  source = "./dns"
}

import {
  for_each = var.dns_import_ids
  to       = module.dns.terraform_data.records[each.key]
  id       = each.value
}
`)
	writeFixture("dns/main.tf", `
resource "terraform_data" "records" {
  for_each = toset(["record/A/global"])
}
`)
	const importID = "dns-import-fixture"
	variables, err := json.Marshal(map[string]any{
		"dns_import_ids": map[string]string{
			"record/A/global": importID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	writeFixture("fixture.tfvars.json", string(variables))

	command := func(endpoint bool, args ...string) ([]byte, error) {
		var environment []string
		if endpoint {
			environment = []string{"PM_API_URL=" + server.URL + "/api2/json"}
		}
		return runTerraform(workspace, environment, args...)
	}
	run := func(args ...string) []byte {
		t.Helper()
		output, err := command(true, args...)
		if err != nil {
			t.Fatalf("terraform %v: %v\n%s", args, err, output)
		}
		return output
	}
	run("init", "-backend=false", "-input=false", "-no-color")

	// Seed only local fixture state; no Proxmox resource is ever created.
	const serviceAttributes = `{"id":"pools/service-fixture","poolid":"service-fixture","comment":"preserve original service bytes"}`
	if seedState {
		writeFixture("terraform.tfstate", `{
  "version": 4,
  "terraform_version": "1.14.8",
  "serial": 1,
  "lineage": "87e53e7d-28cb-4216-807f-fd92ee8c74ea",
  "outputs": {},
  "resources": [{
    "mode": "managed",
    "type": "proxmox_pool",
    "name": "service",
    "provider": "provider[\"registry.terraform.io/telmate/proxmox\"]",
    "instances": [{
      "schema_version": 0,
      "attributes": `+serviceAttributes+`,
      "sensitive_attributes": []
    }]
  }]
}`)
	}

	// Match the owner provider block: omitting pm_api_url makes validation
	// consume the environment default. A Terraform variable would be unknown
	// during validation and conceal the missing-endpoint failure.
	output, err := command(false, "plan", "-input=false", "-no-color",
		"-target=module.dns", "-var-file=fixture.tfvars.json")
	if err == nil || !strings.Contains(string(output), "endpoint for the Proxmox Virtual Environment API") {
		t.Fatalf("targeted plan without endpoint must fail provider validation: %v\n%s", err, output)
	}

	// Only PM_API_URL is supplied; no provider credentials are inherited.
	run("plan", "-input=false", "-no-color", "-target=module.dns",
		"-var-file=fixture.tfvars.json", "-out=import.tfplan")
	planJSON := run("show", "-json", "import.tfplan")
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		// This fixture contains only public built-in resource values and a
		// loopback URL, so its plan is safe evidence for import gate reviews.
		name := "targeted_import_plan.json"
		if !seedState {
			name = "targeted_import_plan_empty_state.json"
		}
		path := filepath.Join(directory, name)
		if err := os.WriteFile(path, planJSON, 0o600); err != nil {
			t.Fatal(err)
		}
		t.Logf("public fixture plan JSON: %s", path)
	}
	var plan struct {
		ResourceChanges []struct {
			Address string `json:"address"`
			Change  struct {
				Actions   []string `json:"actions"`
				Importing *struct {
					ID string `json:"id"`
				} `json:"importing"`
			} `json:"change"`
		} `json:"resource_changes"`
	}
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.ResourceChanges) != 1 {
		t.Fatalf("targeted plan changes = %+v, want exactly one imported DNS fixture", plan.ResourceChanges)
	}
	change := plan.ResourceChanges[0]
	if change.Address != `module.dns.terraform_data.records["record/A/global"]` ||
		len(change.Change.Actions) != 1 || change.Change.Actions[0] != "no-op" ||
		change.Change.Importing == nil || change.Change.Importing.ID != importID {
		t.Fatalf("targeted plan change = %+v, want exact no-op import %q", change, importID)
	}

	// The saved plan carries its own target selection into apply.
	run("apply", "-input=false", "-no-color", "import.tfplan")
	if seedState {
		stateJSON, err := os.ReadFile(filepath.Join(workspace, "terraform.tfstate"))
		if err != nil {
			t.Fatal(err)
		}
		var state struct {
			Resources []struct {
				Type      string `json:"type"`
				Name      string `json:"name"`
				Instances []struct {
					Attributes json.RawMessage `json:"attributes"`
				} `json:"instances"`
			} `json:"resources"`
		}
		if err := json.Unmarshal(stateJSON, &state); err != nil {
			t.Fatal(err)
		}
		var preserved bool
		for _, resource := range state.Resources {
			if resource.Type != "proxmox_pool" || resource.Name != "service" {
				continue
			}
			if len(resource.Instances) != 1 {
				t.Fatalf("service instances = %d, want one", len(resource.Instances))
			}
			var compact bytes.Buffer
			if err := json.Compact(&compact, resource.Instances[0].Attributes); err != nil {
				t.Fatal(err)
			}
			if compact.String() != serviceAttributes {
				t.Fatalf("unrelated service attributes changed: %s", compact.String())
			}
			preserved = true
		}
		if !preserved {
			t.Fatal("targeted apply removed the unrelated service from state")
		}
	}
	run("plan", "-input=false", "-no-color", "-detailed-exitcode",
		"-target=module.dns", "-var-file=fixture.tfvars.json")

	// With the same source and empty credential environment, a full plan
	// reaches Proxmox Configure and fails before making an API request.
	output, err = command(true, "plan", "-input=false", "-no-color", "-var-file=fixture.tfvars.json")
	if err == nil || !strings.Contains(string(output), "your API TokenID username should contain a !") {
		t.Fatalf("untargeted plan must fail Proxmox authentication: %v\n%s", err, output)
	}
}
