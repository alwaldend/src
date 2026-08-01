resource "threexui_panel_subscription" "settings" {
  sub_enable      = true
  sub_json_enable = true
  sub_port        = 2096
  sub_path        = "/${local.base_path_trimmed}/sub/"
  sub_uri         = "https://${var.xui_address}/${local.base_path_trimmed}/sub/"
  sub_json_path   = "/${local.base_path_trimmed}/sub/json/"
  sub_json_uri    = "https://${var.xui_address}/${local.base_path_trimmed}/sub/json/"
  sub_listen      = "127.0.0.1"
}
