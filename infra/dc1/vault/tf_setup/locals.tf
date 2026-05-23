locals {
  dns = jsondecode(file("${path.module}/../dnsconfig.json"))
}
