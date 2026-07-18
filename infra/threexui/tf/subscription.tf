resource "threexui_panel_subscription" "settings" {
  sub_enable      = true
  sub_json_enable = true
  sub_port        = 2096
  sub_path        = "/sub/"
  sub_uri         = "${var.xui_url}/sub/"
  sub_json_path   = "/sub/json/"
  sub_json_uri    = "${var.xui_url}/sub/json/"
  sub_listen      = "127.0.0.1"
}
