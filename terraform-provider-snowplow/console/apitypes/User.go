package apitypes

type User struct {
	ID             string  `json:"id" tfsdk:"id"`
	Email          string  `json:"email" tfsdk:"email"`
	OrganizationID string  `json:"organizationId" tfsdk:"organization_id"`
	FirstName      *string `json:"firstName" tfsdk:"first_name"`
	LastName       *string `json:"lastName" tfsdk:"last_name"`
	JobTitle       *string `json:"jobTitle" tfsdk:"job_title"`
	LastLogin      *string `json:"lastLogin" tfsdk:"last_login"`
	Permissions    []struct {
		OrganizationID string `json:"organizationId" tfsdk:"organization_id"`
		Capabilities   []struct {
			ResourceType string `json:"resourceType" tfsdk:"resource_type"`
			Action       string `json:"action" tfsdk:"action"`
			Filters      []struct {
				Attribute string `json:"attribute" tfsdk:"attribute"`
				Value     string `json:"value" tfsdk:"value"`
			} `json:"filters" tfsdk:"filters"`
		} `json:"capabilities" tfsdk:"capabilities"`
	} `json:"permissions" tfsdk:"permissions"`
}
