package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type inventoryResource struct { client *StarbucksClient }

type inventoryResourceModel struct {
    ID        types.String `tfsdk:"id"`
    StoreID   types.String `tfsdk:"store_id"`
    ItemSKU   types.String `tfsdk:"item_sku"`
    Quantity  types.Int64  `tfsdk:"quantity"`
    Threshold types.Int64  `tfsdk:"threshold"`
}

func NewInventoryResource() resource.Resource { return &inventoryResource{} }

func (r *inventoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (r *inventoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages inventory items for a store.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{Computed: true},
            "store_id": schema.StringAttribute{Required: true},
            "item_sku": schema.StringAttribute{Required: true},
            "quantity": schema.Int64Attribute{Required: true},
            "threshold": schema.Int64Attribute{Optional: true},
        },
    }
}

func (r *inventoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok { resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData)); return }
    r.client = client
}

func (r *inventoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan inventoryResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }

    body := map[string]interface{}{"store_id": plan.StoreID.ValueString(), "item_sku": plan.ItemSKU.ValueString(), "quantity": plan.Quantity.ValueInt64()}
    respBody, err := r.client.DoRequest("POST", "/inventory", body)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create inventory item: %s", err)); return }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }
    if id, ok := result["id"].(string); ok { plan.ID = types.StringValue(id) }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *inventoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state inventoryResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }
    respBody, err := r.client.DoRequest("GET", "/inventory/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read inventory item: %s", err)); return }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }
    if q, ok := result["quantity"].(float64); ok { state.Quantity = types.Int64Value(int64(q)) }
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *inventoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan inventoryResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }
    body := map[string]interface{}{"quantity": plan.Quantity.ValueInt64()}
    _, err := r.client.DoRequest("PUT", "/inventory/"+plan.ID.ValueString(), body)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update inventory item: %s", err)); return }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *inventoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state inventoryResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }
    _, err := r.client.DoRequest("DELETE", "/inventory/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete inventory item: %s", err)); return }
}
