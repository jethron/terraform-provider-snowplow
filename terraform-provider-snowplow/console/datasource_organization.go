package console

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewOrganizationDataSource() datasource.DataSource {
	return &consoleDataSource{
		name:      "organization",
		schema:    &organizationDataSourceSchema,
		populator: populateOrganization,
	}
}

var (
	organizationDataSourceSchema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"domain": schema.StringAttribute{
				Computed: true,
			},
			"tier": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"esso_domain": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"features": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
			},
			"source": schema.ObjectAttribute{
				Computed: true,
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					"name": types.StringType,
					"metadata": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"databricks_organization_id":  types.Int64Type,
							"account_locator":             types.StringType,
							"account_locator_with_region": types.StringType,
						},
					},
				},
			},
			"packages": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.ObjectType{},
			},
			"cloud": schema.ObjectAttribute{
				Computed: true,
				Optional: true,
				AttributeTypes: map[string]attr.Type{
					"provider": types.StringType,
					"accounts": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"provider":                 types.StringType,
								"account_id":               types.StringType,
								"iam_permissions_boundary": types.StringType,
								"subscription_id":          types.StringType,
								"subscription_name":        types.StringType,
								"tenant_id":                types.StringType,
								"project":                  types.StringType,
							},
						},
					},
				},
			},
		},
	}
)

func populateOrganization(client *ApiClient, ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	organizations, err := client.GetOrganizations(ctx)

	if err != nil {
		resp.Diagnostics.AddError("error fetching organizations", err.Error())
		return
	}

	if len(organizations) != 1 {
		resp.Diagnostics.AddError("expected single console organization", fmt.Sprintf("received %d", len(organizations)))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &organizations[0])...)
}
