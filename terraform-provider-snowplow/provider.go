//
// Copyright (c) 2019-2023 Snowplow Analytics Ltd. All rights reserved.
//
// This program is licensed to you under the Apache License Version 2.0,
// and you may not use this file except in compliance with the Apache License Version 2.0.
// You may obtain a copy of the Apache License Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the Apache License Version 2.0 is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the Apache License Version 2.0 for the specific language governing permissions and limitations there under.
//

package main

import (
	"context"
	"errors"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	gt "github.com/snowplow/snowplow-golang-tracker/v2/tracker"

	"github.com/snowplow-devops/terraform-provider-snowplow/terraform-provider-snowplow/console"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &SnowplowProvider{}

// SnowplowProviderModel the struct made from the provider input options
type SnowplowProviderModel struct {
	CollectorURI       types.String `tfsdk:"collector_uri"`
	TrackerAppID       types.String `tfsdk:"tracker_app_id"`
	TrackerNamespace   types.String `tfsdk:"tracker_namespace"`
	TrackerPlatform    types.String `tfsdk:"tracker_platform"`
	EmitterRequestType types.String `tfsdk:"emitter_request_type"`
	EmitterProtocol    types.String `tfsdk:"emitter_protocol"`
	ConsoleAPIEndpoint types.String `tfsdk:"console_api_endpoint"`
	ConsoleAPIKeyID    types.String `tfsdk:"console_api_key_id"`
	ConsoleAPIKey      types.String `tfsdk:"console_api_key"`
	ConsoleOrgID       types.String `tfsdk:"console_organization_id"`
}

// SnowplowProvider defines the provider implementation.
type SnowplowProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type ResourceData struct {
	*SnowplowProviderModel
	*console.ApiClient
}

func (r ResourceData) GetApiClient() *console.ApiClient {
	return r.ApiClient
}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SnowplowProvider{version: version}
	}
}

func (p *SnowplowProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "snowplow"
	resp.Version = p.version
}

func (p *SnowplowProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for emitting Snowplow events",
		Attributes: map[string]schema.Attribute{
			"collector_uri": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "URI of your Snowplow Collector",
			},
			"tracker_app_id": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional application ID",
			},
			"tracker_namespace": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional namespace",
			},
			"tracker_platform": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Optional platform",
			},
			"emitter_request_type": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Whether to use GET or POST requests to emit events",
			},
			"emitter_protocol": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Whether to use HTTP or HTTPS to send events",
			},
			"console_api_endpoint": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "API endpoint hostname to use when interacting with the Console API",
			},
			"console_api_key_id": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Auth API v3 API Key ID to access the Console API with",
			},
			"console_api_key": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Auth API v2/v3 API Key to access the Console API with",
			},
			"console_organization_id": schema.StringAttribute{
				Optional:    true,
				Required:    false,
				Description: "Organization ID associated with the console_api_key credentials",
			},
		},
	}
}

func (p *SnowplowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SnowplowProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For sending events

	if data.EmitterRequestType.ValueString() == "" {
		data.EmitterRequestType = types.StringValue("POST")
	}
	if data.EmitterProtocol.ValueString() == "" {
		data.EmitterProtocol = types.StringValue("HTTPS")
	}
	if data.TrackerPlatform.ValueString() == "" {
		data.TrackerPlatform = types.StringValue("srv")
	}

	// Console interactions

	// Mostly only useful for testing with next.console
	if data.ConsoleAPIEndpoint.ValueString() == "" {
		data.ConsoleAPIEndpoint = types.StringValue("console.snowplowanalytics.com")
	}

	// As used by Snowtype: https://docs.snowplow.io/docs/collecting-data/code-generation/using-the-cli/#authenticating-with-the-console
	if data.ConsoleAPIKey.ValueString() == "" {
		data.ConsoleAPIKey = types.StringValue(os.Getenv("SNOWPLOW_CONSOLE_API_KEY"))
	}
	if data.ConsoleAPIKeyID.ValueString() == "" {
		data.ConsoleAPIKeyID = types.StringValue(os.Getenv("SNOWPLOW_CONSOLE_API_KEY_ID"))
	}

	// No precedent, but consistent convention
	if data.ConsoleOrgID.ValueString() == "" {
		data.ConsoleOrgID = types.StringValue(os.Getenv("SNOWPLOW_CONSOLE_ORGANIZATION_ID"))
	}

	// TODO: Fallback to snowplow-cli config: https://github.com/snowplow-product/snowplow-cli?tab=readme-ov-file#configuration

	// Console Auth
	var client *console.ApiClient

	if data.ConsoleAPIKey.ValueString() != "" {
		newClient, err := console.NewApiClient(
			ctx,
			version,
			data.ConsoleAPIEndpoint.ValueString(),
			data.ConsoleAPIKeyID.ValueString(),
			data.ConsoleAPIKey.ValueString(),
			data.ConsoleOrgID.ValueString(),
		)

		if err != nil {
			resp.Diagnostics.AddError("error authenticating with snowplow console api", err.Error())
		}

		client = newClient
	}

	rd := ResourceData{&data, client}
	resp.DataSourceData = rd
	resp.ResourceData = rd
}

func (p *SnowplowProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTrackSelfDescribingEventResource,
	}
}

func (p *SnowplowProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		console.NewOrganizationDataSource,
		console.NewUserDataSource,
		console.NewUsersDataSource,
		console.NewPipelineDataSource,
		console.NewPipelinesDataSource,
	}
}

// InitTracker takes a context and a channel of size 1 and returns
// a new Snowplow Tracker ready to create a resource
func InitTracker(ctx SnowplowProviderModel, ctxResource trackSelfDescribingEventResourceModel, trackerChan chan int) (*gt.Tracker, error) {
	var collectorUri, emitterRequestType, emitterProtocol, trackerNamespace, trackerAppId, trackerPlatform string

	if ctxResource.CollectorURI.ValueString() == "" {
		collectorUri = ctx.CollectorURI.ValueString()
	} else {
		collectorUri = ctxResource.CollectorURI.ValueString()
	}

	if collectorUri == "" {
		return nil, errors.New("URI of the Snowplow Collector is empty - this can be set either at the provider or resource level with the 'collector_uri' input")
	}

	if ctxResource.EmitterRequestType.ValueString() == "" {
		emitterRequestType = ctx.EmitterRequestType.ValueString()
	} else {
		emitterRequestType = ctxResource.EmitterRequestType.ValueString()
	}

	if ctxResource.EmitterProtocol.ValueString() == "" {
		emitterProtocol = ctx.EmitterProtocol.ValueString()
	} else {
		emitterProtocol = ctxResource.EmitterProtocol.ValueString()
	}

	if ctxResource.TrackerNamespace.IsNull() {
		trackerNamespace = ctx.TrackerNamespace.ValueString()
	} else {
		trackerNamespace = ctxResource.TrackerNamespace.ValueString()
	}

	if ctxResource.TrackerAppID.IsNull() {
		trackerAppId = ctx.TrackerAppID.ValueString()
	} else {
		trackerAppId = ctxResource.TrackerAppID.ValueString()
	}

	if ctxResource.TrackerPlatform.IsNull() {
		trackerPlatform = ctx.TrackerPlatform.ValueString()
	} else {
		trackerPlatform = ctxResource.TrackerPlatform.ValueString()
	}

	callback := func(s []gt.CallbackResult, f []gt.CallbackResult) {
		status := 0

		if len(s) == 1 {
			status = s[0].Status
		} else if len(f) == 1 {
			status = f[0].Status
		}

		trackerChan <- status
	}

	emitter := gt.InitEmitter(
		gt.RequireCollectorUri(collectorUri),
		gt.OptionRequestType(emitterRequestType),
		gt.OptionProtocol(emitterProtocol),
		gt.OptionCallback(callback),
		gt.OptionStorage(gt.InitStorageMemory()),
	)

	subject := gt.InitSubject()

	tracker := gt.InitTracker(
		gt.RequireEmitter(emitter),
		gt.OptionSubject(subject),
		gt.OptionNamespace(trackerNamespace),
		gt.OptionAppId(trackerAppId),
		gt.OptionPlatform(trackerPlatform),
		gt.OptionBase64Encode(true),
	)

	return tracker, nil
}
