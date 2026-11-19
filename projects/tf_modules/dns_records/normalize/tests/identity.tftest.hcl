variables {
  zone = "example.test"
}

run "original_record" {
  command = plan

  variables {
    document = {
      records = {
        service = {
          A   = { name = "service", address = "192.0.2.10", ttl = 300 }
          dsp = ["all"]
        }
      }
    }
  }
}

run "changes_values_without_changing_identity" {
  command = plan

  variables {
    document = {
      records = {
        service = {
          A   = { name = "renamed-service", address = "192.0.2.11", ttl = 1200 }
          dsp = ["dc1", "global"]
        }
      }
    }
  }

  assert {
    condition     = keys(output.normalized_records) == keys(run.original_record.normalized_records)
    error_message = "Changing DNS name, address, TTL, or equivalent destination spelling must not change logical record identities."
  }

  assert {
    condition = alltrue([
      for record in values(output.normalized_records) :
      record.name == "renamed-service.example.test" && record.value == "192.0.2.11" && record.ttl == 1200
    ])
    error_message = "Stable resource identities must still carry the changed record attributes."
  }
}

run "removes_a_destination_without_renaming_existing_records" {
  command = plan

  variables {
    document = {
      records = {
        service = {
          A   = { name = "service", address = "192.0.2.10", ttl = 300 }
          dsp = ["global"]
        }
      }
    }
  }

  assert {
    condition = (
      keys(output.normalized_records) == ["service/A/global"] &&
      output.normalized_records["service/A/global"] == run.original_record.normalized_records["service/A/global"]
    )
    error_message = "The global record must retain its identity and content when the dc1 destination is removed."
  }
}
