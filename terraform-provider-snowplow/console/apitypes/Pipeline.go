package apitypes

type Pipeline struct {
	ID                 string   `json:"id" tfsdk:"id"`
	Name               string   `json:"name" tfsdk:"name"`
	CloudProvider      string   `json:"cloudProvider" tfsdk:"cloud_provider"`
	CollectorEndpoints []string `json:"collectorEndpoints" tfsdk:"collector_endpoints"`
}
