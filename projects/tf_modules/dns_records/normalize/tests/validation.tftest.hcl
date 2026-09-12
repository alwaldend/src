variables {
  zone = "example.test"
}

run "rejects_missing_records" {
  command = plan

  variables {
    document = {}
  }

  expect_failures = [var.document]
}

run "rejects_unknown_record_type" {
  command = plan

  variables {
    document = { records = { bad = { CAA = { name = "@", value = "0 issue example.test" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_record_without_any_type" {
  command = plan

  variables {
    document = { records = { bad = { dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_unknown_type_field" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1", proxied = true }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_missing_address" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_type_lists" {
  command = plan

  variables {
    document = { records = { bad = { A = [{ name = "host", address = "192.0.2.1" }], dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_unknown_destination" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1" }, dsp = ["unknown"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_missing_destinations" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1" } } } }
  }

  expect_failures = [var.document]
}

run "rejects_empty_destinations" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1" }, dsp = [] } } }
  }

  expect_failures = [var.document]
}

run "rejects_scalar_destination" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1" }, dsp = "global" } } }
  }

  expect_failures = [var.document]
}

run "rejects_absolute_owner_outside_zone" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host.external.test.", address = "192.0.2.1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_invalid_owner_name" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "invalid..host", address = "192.0.2.1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_empty_owner_name" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "", address = "192.0.2.1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_zero_ttl" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1", ttl = 0 }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_fractional_ttl" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1", ttl = 300.5 }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_nonnumeric_ttl" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.1", ttl = "invalid" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_ipv6_in_a_record" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "2001:db8::1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_ipv4_in_aaaa_record" {
  command = plan

  variables {
    document = { records = { bad = { AAAA = { name = "host", address = "192.0.2.1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_invalid_ipv4" {
  command = plan

  variables {
    document = { records = { bad = { A = { name = "host", address = "192.0.2.999" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_empty_domain_target" {
  command = plan

  variables {
    document = { records = { bad = { CNAME = { name = "alias", target = "" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_missing_mx_priority" {
  command = plan

  variables {
    document = { records = { bad = { MX = { name = "@", target = "mail.example.test." }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_out_of_range_mx_priority" {
  command = plan

  variables {
    document = { records = { bad = { MX = { name = "@", target = "mail.example.test.", priority = 65536 }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_invalid_logical_key" {
  command = plan

  variables {
    document = { records = { "nested/key" = { A = { name = "host", address = "192.0.2.1" }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}

run "rejects_null_txt_content" {
  command = plan

  variables {
    document = { records = { bad = { TXT = { name = "@", content = null }, dsp = ["global"] } } }
  }

  expect_failures = [var.document]
}
