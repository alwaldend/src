locals {
  approles = jsondecode(file("${path.module}/../../vault/tf/output/approles.json")).approles
}
