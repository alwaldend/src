variables {
  zone = "example.test"
}

run "expands_each_type_and_deduplicates_destinations" {
  command = plan

  variables {
    document = {
      records = {
        host = {
          A    = { name = "host", address = "192.0.2.10", ttl = 1200 }
          AAAA = { name = "host", address = "2001:db8::10" }
          TXT  = { name = "host", content = "purpose=fixture; value=unchanged." }
          dsp  = ["all", "global", "dc1", "all"]
        }
        alias = {
          CNAME = { name = "alias", target = "host" }
          dsp   = ["global"]
        }
        nameserver = {
          NS  = { name = "delegated", address = "ns.external.test.", ttl = 3600 }
          dsp = ["dc1"]
        }
        mail = {
          MX  = { name = "@", target = "mail.external.test.", priority = 10, ttl = 10800 }
          dsp = ["global"]
        }
      }
    }
  }

  assert {
    condition = toset(keys(output.normalized_records)) == toset([
      "host/A/global", "host/A/dc1",
      "host/AAAA/global", "host/AAAA/dc1",
      "host/TXT/global", "host/TXT/dc1",
      "alias/CNAME/global", "nameserver/NS/dc1", "mail/MX/global",
    ])
    error_message = "Every supplied type must be emitted once per expanded destination."
  }

  assert {
    condition = alltrue([
      for key, record in output.normalized_records :
      record.key == key && !record.proxied && contains(["global", "dc1"], record.view)
    ])
    error_message = "Normalized records must retain their stable key, concrete view, and unproxied setting."
  }

  assert {
    condition = (
      output.normalized_records["host/A/global"].name == "host.example.test" &&
      output.normalized_records["host/A/global"].relative_name == "host" &&
      output.normalized_records["host/A/global"].value == "192.0.2.10" &&
      output.normalized_records["host/A/global"].ttl == 1200 &&
      output.normalized_records["host/AAAA/global"].value == "2001:db8::10" &&
      output.normalized_records["host/AAAA/global"].ttl == 600
    )
    error_message = "Address records must retain IP values and apply explicit or default TTLs."
  }

  assert {
    condition = (
      output.normalized_records["alias/CNAME/global"].value == "host.example.test" &&
      output.normalized_records["alias/CNAME/global"].ttl == 600 &&
      output.normalized_records["nameserver/NS/dc1"].value == "ns.external.test" &&
      output.normalized_records["nameserver/NS/dc1"].ttl == 3600
    )
    error_message = "Domain targets must expand relative names and remove absolute-name trailing dots."
  }

  assert {
    condition = (
      output.normalized_records["mail/MX/global"].name == "example.test" &&
      output.normalized_records["mail/MX/global"].relative_name == "@" &&
      output.normalized_records["mail/MX/global"].value == "mail.external.test" &&
      output.normalized_records["mail/MX/global"].priority == 10 &&
      output.normalized_records["mail/MX/global"].ttl == 10800 &&
      alltrue([for record in values(output.normalized_records) : record.priority == null if record.type != "MX"])
    )
    error_message = "MX priority must be preserved and remain absent from other record types."
  }

  assert {
    condition = (
      output.normalized_records["host/TXT/global"].value == "purpose=fixture; value=unchanged." &&
      output.normalized_records["host/TXT/global"].ttl == 300
    )
    error_message = "TXT content must remain byte-for-byte unchanged, including punctuation and its trailing dot."
  }
}

run "normalizes_apex_relative_and_fully_qualified_names" {
  command = plan

  variables {
    document = {
      records = {
        apex = {
          A   = { name = "example.test.", address = "192.0.2.1" }
          dsp = ["global"]
        }
        fqdn = {
          A   = { name = "host.example.test.", address = "192.0.2.2" }
          dsp = ["global"]
        }
        apex_alias = {
          CNAME = { name = "www", target = "@" }
          dsp   = ["global"]
        }
        fqdn_target = {
          CNAME = { name = "service", target = "host.example.test." }
          dsp   = ["dc1"]
        }
        relative_target = {
          NS  = { name = "delegated", address = "ns.delegated" }
          dsp = ["dc1"]
        }
        mail = {
          MX  = { name = "@", target = "mail", priority = 0 }
          dsp = ["global"]
        }
      }
    }
  }

  assert {
    condition = (
      output.normalized_records["apex/A/global"].name == "example.test" &&
      output.normalized_records["apex/A/global"].relative_name == "@" &&
      output.normalized_records["fqdn/A/global"].name == "host.example.test" &&
      output.normalized_records["fqdn/A/global"].relative_name == "host"
    )
    error_message = "An in-zone fully qualified owner must not receive the zone suffix twice."
  }

  assert {
    condition = (
      output.normalized_records["apex_alias/CNAME/global"].value == "example.test" &&
      output.normalized_records["fqdn_target/CNAME/dc1"].value == "host.example.test" &&
      output.normalized_records["relative_target/NS/dc1"].value == "ns.delegated.example.test" &&
      output.normalized_records["mail/MX/global"].value == "mail.example.test" &&
      output.normalized_records["mail/MX/global"].priority == 0
    )
    error_message = "All domain-valued types must share apex, relative, and absolute target normalization."
  }

  assert {
    condition = (
      output.normalized_records["apex/A/global"].ttl == 300 &&
      output.normalized_records["relative_target/NS/dc1"].ttl == 300 &&
      output.normalized_records["mail/MX/global"].ttl == 300
    )
    error_message = "A, NS, and MX must retain the effective legacy default TTL of 300 seconds."
  }
}

run "supports_separate_logical_records_in_one_rrset" {
  command = plan

  variables {
    document = {
      records = {
        primary = {
          A   = { name = "service", address = "192.0.2.10" }
          dsp = ["global"]
        }
        secondary = {
          A   = { name = "service", address = "192.0.2.11" }
          dsp = ["global"]
        }
      }
    }
  }

  assert {
    condition = (
      toset(keys(output.normalized_records)) == toset(["primary/A/global", "secondary/A/global"]) &&
      output.normalized_records["primary/A/global"].name == output.normalized_records["secondary/A/global"].name &&
      toset([for record in values(output.normalized_records) : record.value]) == toset(["192.0.2.10", "192.0.2.11"])
    )
    error_message = "Separate logical keys must preserve multiple values belonging to the same DNS RRset."
  }
}

run "accepts_an_empty_record_map" {
  command = plan

  variables {
    document = { records = {} }
  }

  assert {
    condition     = length(output.normalized_records) == 0
    error_message = "An empty owner document must produce no records."
  }
}
