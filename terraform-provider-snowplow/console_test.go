package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const consoleProviderConfig = `
provider snowplow {}
`

func TestOrganizationDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: consoleProviderConfig + `
				data "snowplow_organization" "test" {}
				data "snowplow_users" "test" {}
				data "snowplow_pipelines" "test" {}

				output "output" {
					value = jsonencode([
						data.snowplow_organization.test,
						data.snowplow_users.test,
						data.snowplow_pipelines.test,
					])
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowplow_organization.test", "id"),
					resource.TestCheckOutput("output", "{}"),
				),
			},
		},
	})
}
