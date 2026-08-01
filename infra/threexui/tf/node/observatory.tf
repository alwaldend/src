resource "threexui_xray_observatory" "observatory" {
  observatory {
    tag                = "obs_ping"
    subject_selector   = ["out-*"]
    probe_url          = "https://www.google.com/generate_204"
    probe_interval     = "1m"
    enable_concurrency = true
  }

  # burst_observatory {
  #   tag              = "burst_ping"
  #   subject_selector = ["out-*"]
  #
  #   ping_config {
  #     destination     = "https://www.cloudflare.com/cdn-cgi/trace"
  #     interval        = "1m"
  #     connect_timeout = "5s"
  #     timeout         = "10s"
  #     samples         = 3
  #     sampling_count  = 2
  #     lazy            = true
  #   }
  # }
}
