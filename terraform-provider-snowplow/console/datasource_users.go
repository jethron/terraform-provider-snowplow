package console

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewUsersDataSource() datasource.DataSource {
	return &consoleDataSource{
		name:      "users",
		schema:    &usersDataSourceSchema,
		populator: populateUsers,
	}
}

var (
	usersDataSourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"users": schema.ListAttribute{
				Computed: true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":              types.StringType,
						"email":           types.StringType,
						"organization_id": types.StringType,
						"first_name":      types.StringType,
						"last_name":       types.StringType,
						"job_title":       types.StringType,
						"last_login":      types.StringType,
						"permissions": types.ListType{
							ElemType: types.ObjectType{
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
				},
			},
		},
	}
)

func populateUsers(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	users, err := client.GetUsers(ctx)

	if err != nil {
		resp.Diagnostics.AddError("error fetching users", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("users"), &users)...)
}
