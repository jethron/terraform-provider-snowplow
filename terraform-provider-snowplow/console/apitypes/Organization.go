package apitypes

type Organization struct {
	ID         string   `json:"id" tfsdk:"id"`
	Name       *string  `json:"name" tfsdk:"name"`
	Domain     *string  `json:"domain" tfsdk:"domain"`
	Tier       string   `json:"tier" tfsdk:"tier"`
	Tags       []string `json:"tags" tfsdk:"tags"`
	ESSODomain *string  `json:"essoDomain" tfsdk:"esso_domain"`
	Features   []string `json:"features" tfsdk:"features"`
	Source     *struct {
		Name     *string `json:"name" tfsdk:"name"`
		Metadata struct {
			DatabricksOrganizationID int64  `json:"databricksOrganizationId" tfsdk:"databricks_organization_id"`
			AccountLocator           string `json:"accountLocator" tfsdk:"account_locator"`
			AccountLocatorWithRegion string `json:"accountLocatorWithRegion" tfsdk:"account_locator_with_region"`
		} `json:"metadata" tfsdk:"metadata"`
	} `json:"source" tfsdk:"source"`
	Packages []struct{} `json:"packages" tfsdk:"packages"`
	Cloud    *struct {
		Provider string `json:"provider" tfsdk:"provider"`
		Accounts []struct {
			Provider               string  `json:"provider" tfsdk:"provider"`
			AccountID              *string `json:"accountId" tfsdk:"account_id"`
			IAMPermissionsBoundary *string `json:"iamPermissionsBoundary" tfsdk:"iam_permissions_boundary"`
			SubscriptionID         *string `json:"subscriptionId" tfsdk:"subscription_id"`
			SubscriptionName       *string `json:"subscriptionName" tfsdk:"subscription_name"`
			TenantID               *string `json:"tenantId" tfsdk:"tenant_id"`
			Project                *string `json:"project" tfsdk:"project"`
		} `json:"accounts" tfsdk:"accounts"`
	} `json:"cloud" tfsdk:"cloud"`
}
