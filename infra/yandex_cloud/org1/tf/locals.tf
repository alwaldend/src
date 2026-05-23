locals {
  approles = jsondecode(file("${path.module}/../../../dc1/vault/tf/output/approles.json")).approles
}
