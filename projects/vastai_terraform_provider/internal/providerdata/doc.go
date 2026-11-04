// Package providerdata defines and validates dependencies shared by Terraform
// provider configuration with resources and data sources. The provider owns
// construction; consumers borrow the payload without owning its lifecycle.
package providerdata
