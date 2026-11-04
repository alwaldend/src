package providerdata

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	vastclient "git.alwaldend.com/alwaldend/src/internal/vastai_sdk/pkg/client"
)

// Data contains dependencies shared by provider configuration with its consumers.
// The provider constructs it; resources and data sources retain the same payload.
type Data struct {
	// Client supplies authenticated Vast.ai operations to configured consumers.
	Client *vastclient.Client
}

// From checks provider data before resources and data sources retain it.
func From(data any, diagnostics *diag.Diagnostics) *Data {
	if data == nil {
		return nil
	}
	configured, ok := data.(*Data)
	if !ok || configured == nil || configured.Client == nil {
		diagnostics.AddError("Invalid provider data", "Expected initialized Vast.ai provider data.")
		return nil
	}
	return configured
}
