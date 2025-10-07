package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type storeDataSource struct { client *StarbucksClient }

type storeDataSourceModel struct {
    ID    types.String `tfsdk:"id"`
    Name  types.String `tfsdk:"name"`
    City  types.String `tfsdk:"city"`
    State types.String `tfsdk:"state"`
}

func NewStoreDataSource() datasource.DataSource { return &storeDataSource{} }

func (d *storeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_store"
}

func (d *storeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{Required: true},
            "name": schema.StringAttribute{Computed: true},
            "city": schema.StringAttribute{Computed: true},
            "state": schema.StringAttribute{Computed: true},
        },
    }
}

func (d *storeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok { resp.Diagnostics.AddError("Unexpected DataSource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData)); return }
    d.client = client
}

func (d *storeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state storeDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &state)...) 
    if resp.Diagnostics.HasError() { return }

    respBody, err := d.client.DoRequest("GET", "/stores/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store: %s", err)); return }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }
    if v, ok := result["name"].(string); ok { state.Name = types.StringValue(v) }
    if v, ok := result["city"].(string); ok { state.City = types.StringValue(v) }
    if v, ok := result["state"].(string); ok { state.State = types.StringValue(v) }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
