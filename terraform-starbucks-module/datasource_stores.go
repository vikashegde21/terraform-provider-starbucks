package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type storesDataSource struct { client *StarbucksClient }

type storesDataSourceModel struct {
    Stores []map[string]interface{} `tfsdk:"stores"`
}

func NewStoresDataSource() datasource.DataSource { return &storesDataSource{} }

func (d *storesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_stores"
}

func (d *storesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "stores": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "id": schema.StringAttribute{Computed: true},
                        "name": schema.StringAttribute{Computed: true},
                        "city": schema.StringAttribute{Computed: true},
                    },
                },
            },
        },
    }
}

func (d *storesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok { resp.Diagnostics.AddError("Unexpected DataSource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData)); return }
    d.client = client
}

func (d *storesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    respBody, err := d.client.DoRequest("GET", "/stores", nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list stores: %s", err)); return }
    var result []map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }

    model := storesDataSourceModel{Stores: result}
    resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
