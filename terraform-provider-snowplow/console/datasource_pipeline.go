package console

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewPipelineDataSource() datasource.DataSource {
	return &consoleDataSource{
		name:      "pipeline",
		schema:    &pipelineDataSourceSchema,
		populator: populatePipeline,
	}
}

var (
	pipelineDataSourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"cloud_provider": schema.StringAttribute{
				Computed: true,
			},
			"collector_endpoints": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
)

func populatePipeline(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configId types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &configId)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestedId := configId.ValueString()
	pipeline, err := client.GetPipeline(ctx, requestedId)

	if err != nil {
		resp.Diagnostics.AddError("error fetching pipeline", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &pipeline)...)
}
