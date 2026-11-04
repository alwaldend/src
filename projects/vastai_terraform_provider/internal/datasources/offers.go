package datasources

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"git.alwaldend.com/alwaldend/src/projects/vastai_terraform_provider/internal/providerdata"
)

// offersDataSource translates a marketplace query into typed Terraform results.
type offersDataSource struct {
	// providerData contains shared dependencies supplied by provider configuration.
	providerData *providerdata.Data
}

// offersModel contains the query and a typed result list.
type offersModel struct {
	// Query is a JSON object in Vast.ai filter syntax.
	Query types.String `tfsdk:"query"`
	// Offers preserves the API's result ordering, including an empty list.
	Offers []offerModel `tfsdk:"offers"`
}

// offerModel contains the supported marketplace result attributes.
type offerModel struct {
	// ID selects the offer when creating an instance.
	ID types.Int64 `tfsdk:"id"`
	// GPUName names the GPU model.
	GPUName types.String `tfsdk:"gpu_name"`
	// GPUs is the available GPU count.
	GPUs types.Int64 `tfsdk:"num_gpus"`
	// GPURAM is memory per GPU in MiB.
	GPURAM types.Float64 `tfsdk:"gpu_ram_mb"`
	// HourlyPrice is the total hourly offer price in USD reported by Vast.ai.
	HourlyPrice types.Float64 `tfsdk:"hourly_price"`
}

// Metadata identifies the data source under the provider's prefix.
func (d *offersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_offers"
}

// Schema defines the JSON filter input and typed offer results.
func (d *offersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Search Vast.ai GPU offers without renting. Results are volatile; an offer can disappear before apply.", Attributes: map[string]schema.Attribute{
		"query": schema.StringAttribute{Required: true, Description: "JSON object with Vast.ai filter operators, order, type, and limit. Use jsonencode; set type to on-demand for this provider's rental surface."},
		"offers": schema.ListNestedAttribute{Computed: true, Description: "Offers in API order; may be empty.", NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"id":           schema.Int64Attribute{Computed: true, Description: "Offer ID accepted by vastai_instance.offer_id."},
			"gpu_name":     schema.StringAttribute{Computed: true, Description: "GPU model name."},
			"num_gpus":     schema.Int64Attribute{Computed: true, Description: "Number of GPUs."},
			"gpu_ram_mb":   schema.Float64Attribute{Computed: true, Description: "Memory per GPU in MiB."},
			"hourly_price": schema.Float64Attribute{Computed: true, Description: "API-reported total hourly price in USD, affected by search parameters."},
		}}},
	}}
}

// Configure receives the shared provider client.
func (d *offersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.providerData = providerdata.From(req.ProviderData, &resp.Diagnostics)
}

// Read validates the query before searching and maps only non-secret result fields.
func (d *offersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data offersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var query map[string]any
	if err := json.Unmarshal([]byte(data.Query.ValueString()), &query); err != nil || query == nil {
		resp.Diagnostics.AddError("Invalid offer query", "query must be a valid JSON object, produced with jsonencode.")
		return
	}
	offers, err := d.providerData.Client.Search(ctx, query)
	if err != nil {
		resp.Diagnostics.AddError("Cannot search offers", err.Error())
		return
	}
	data.Offers = make([]offerModel, 0, len(offers))
	for _, offer := range offers {
		data.Offers = append(data.Offers, offerModel{ID: types.Int64Value(offer.ID), GPUName: types.StringValue(offer.GPUName), GPUs: types.Int64Value(offer.GPUs), GPURAM: types.Float64Value(offer.GPURAM), HourlyPrice: types.Float64Value(offer.HourlyPrice)})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// NewOffers constructs an unconfigured offer data source for provider registration.
func NewOffers() datasource.DataSource {
	return &offersDataSource{}
}
