package console

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewUserDataSource() datasource.DataSource {
	return &consoleDataSource{
		name:      "user",
		schema:    &userDataSourceSchema,
		populator: populateUser,
	}
}

var (
	userDataSourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Computed: true,
			},
			"last_name": schema.StringAttribute{
				Computed: true,
			},
			"job_title": schema.StringAttribute{
				Computed: true,
			},
			"last_login": schema.StringAttribute{
				Computed: true,
			},
			"permissions": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"organization_id": types.StringType,
						"capabilities": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"resource_type": types.StringType,
									"action":        types.StringType,
									"filters": types.ListType{
										ElemType: types.ObjectType{
											AttrTypes: map[string]attr.Type{
												"attribute": types.StringType,
												"value":     types.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
)

func populateUser(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var configId types.String

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &configId)...)

	if resp.Diagnostics.HasError() {
		return
	}

	requestedId := configId.ValueString()
	user, err := client.GetUser(ctx, requestedId)

	if err != nil {
		resp.Diagnostics.AddError("error fetching user", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &user)...)
}
