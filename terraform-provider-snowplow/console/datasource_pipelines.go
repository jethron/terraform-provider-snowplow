package console

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewPipelinesDataSource() datasource.DataSource {
	return &consoleDataSource{
		name:      "pipelines",
		schema:    &pipelinesDataSourceSchema,
		populator: populatePipelines,
	}
}

var (
	pipelinesDataSourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"pipelines": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":             types.StringType,
						"name":           types.StringType,
						"cloud_provider": types.StringType,
						"collector_endpoints": types.ListType{
							ElemType: types.StringType,
						},
					},
				},
			},
		},
	}
)

func populatePipelines(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	pipelines, err := client.GetPipelines(ctx)

	if err != nil {
		resp.Diagnostics.AddError("error fetching pipelines", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipelines"), &pipelines)...)
}
