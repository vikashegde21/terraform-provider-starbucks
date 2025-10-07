package main

import (
    "context"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &starbucksProvider{}

type starbucksProvider struct {
    version string
}

type starbucksProviderModel struct {
    APIKey   types.String `tfsdk:"api_key"`
    Endpoint types.String `tfsdk:"endpoint"`
    Region   types.String `tfsdk:"region"`
    Timeout  types.Int64  `tfsdk:"timeout"`
}

func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &starbucksProvider{
            version: version,
        }
    }
}

func (p *starbucksProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "starbucks"
    resp.Version = p.version
}

func (p *starbucksProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Terraform provider for managing Starbucks infrastructure, stores, employees, and operations.",
        Attributes: map[string]schema.Attribute{
            "api_key": schema.StringAttribute{
                Description: "API key for Starbucks Management API. Can also be set via STARBUCKS_API_KEY environment variable.",
                Optional:    true,
                Sensitive:   true,
            },
            "endpoint": schema.StringAttribute{
                Description: "API endpoint URL. Defaults to https://api.starbucks.com/v1",
                Optional:    true,
            },
            "region": schema.StringAttribute{
                Description: "Region for API calls (e.g., us-west-2, us-east-1)",
                Optional:    true,
            },
            "timeout": schema.Int64Attribute{
                Description: "API request timeout in seconds. Defaults to 30.",
                Optional:    true,
            },
        },
    }
}

func (p *starbucksProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config starbucksProviderModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...) 
    if resp.Diagnostics.HasError() {
        return
    }

    if config.APIKey.IsUnknown() {
        resp.Diagnostics.AddWarning(
            "Unable to create client",
            "Cannot use unknown value as api_key",
        )
        return
    }

    apiKey := os.Getenv("STARBUCKS_API_KEY")
    endpoint := "https://api.starbucks.com/v1"
    region := "us-west-2"
    timeout := int64(30)

    if !config.APIKey.IsNull() {
        apiKey = config.APIKey.ValueString()
    }
    if !config.Endpoint.IsNull() {
        endpoint = config.Endpoint.ValueString()
    }
    if !config.Region.IsNull() {
        region = config.Region.ValueString()
    }
    if !config.Timeout.IsNull() {
        timeout = config.Timeout.ValueInt64()
    }

    if apiKey == "" {
        resp.Diagnostics.AddError(
            "Missing API Key Configuration",
            "API key must be provided via api_key attribute or STARBUCKS_API_KEY environment variable",
        )
        return
    }

    client := NewStarbucksClient(apiKey, endpoint, region, timeout)
    resp.DataSourceData = client
    resp.ResourceData = client
}

func (p *starbucksProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewStoreResource,
        NewEmployeeResource,
        NewMenuItemResource,
        NewInventoryResource,
        NewPromotionResource,
    }
}

func (p *starbucksProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewStoreDataSource,
        NewStoresDataSource,
    }
}
