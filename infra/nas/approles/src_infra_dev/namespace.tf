module "namespace" {
  source = "../../../../projects/tf_modules/truenas_dataset"
  parent = "default/encrypted/approles"
  config = {
    name   = "src_infra_dev"
    owner  = "src_infra_dev"
    access = { src_infra_xcp_ng = "TRAVERSE" }
    smb    = {}
  }
}
