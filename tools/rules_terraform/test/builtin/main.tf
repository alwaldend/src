resource "terraform_data" "example" {
  input = "packaged Terraform"
}

output "value" {
  value = terraform_data.example.output
}
