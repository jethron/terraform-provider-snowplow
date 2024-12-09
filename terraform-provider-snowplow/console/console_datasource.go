package console

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &consoleDataSource{}
var _ datasource.DataSourceWithConfigure = &consoleDataSource{}
var _ datasource.DataSourceWithConfigValidators = &consoleDataSource{}

type consoleDataSource struct {
	name       string
	schema     *schema.Schema
	client     *ApiClient
	populator  func(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse)
	validators []datasource.ConfigValidator
}

func (r *consoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.name
}

func (r *consoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = *r.schema
}

func (r *consoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerCtx, ok := req.ProviderData.(ApiClientProvider)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected ApiClientProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = providerCtx.GetApiClient()
}

func (r *consoleDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return r.validators
}

func (r *consoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("console api client not configured", "set provider console_api_key or use SNOWPLOW_CONSOLE_API_KEY")
		return
	}

	r.populator(r.client, ctx, req, resp)
}
