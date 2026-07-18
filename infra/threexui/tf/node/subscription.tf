resource "threexui_panel_subscription" "settings" {
  sub_enable      = true
  sub_json_enable = true
  sub_port        = 2096
  sub_path        = "/${var.xui_base_path}/sub/"
  sub_uri         = "https://${var.xui_address}/${var.xui_base_path}/sub/"
  sub_json_path   = "/${var.xui_base_path}/sub/json/"
  sub_json_uri    = "https://${var.xui_address}/${var.xui_base_path}/sub/json/"
  sub_listen      = "127.0.0.1"
}
