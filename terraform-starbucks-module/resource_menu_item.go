package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type menuItemResource struct {
    client *StarbucksClient
}

type menuItemResourceModel struct {
    ID          types.String  `tfsdk:"id"`
    Name        types.String  `tfsdk:"name"`
    Category    types.String  `tfsdk:"category"`
    Size        types.String  `tfsdk:"size"`
    Price       types.Float64 `tfsdk:"price"`
    Calories    types.Int64   `tfsdk:"calories"`
    Description types.String  `tfsdk:"description"`
    IsAvailable types.Bool    `tfsdk:"is_available"`
    IsSeasonal  types.Bool    `tfsdk:"is_seasonal"`
}

func NewMenuItemResource() resource.Resource { return &menuItemResource{} }

func (r *menuItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_menu_item"
}

func (r *menuItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a Starbucks menu item.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{Computed: true},
            "name": schema.StringAttribute{Required: true},
            "category": schema.StringAttribute{Optional: true},
            "size": schema.StringAttribute{Optional: true},
            "price": schema.Float64Attribute{Optional: true},
            "calories": schema.Int64Attribute{Optional: true},
            "description": schema.StringAttribute{Optional: true},
            "is_available": schema.BoolAttribute{Optional: true, Computed: true},
            "is_seasonal": schema.BoolAttribute{Optional: true, Computed: true},
        },
    }
}

func (r *menuItemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok {
        resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData))
        return
    }
    r.client = client
}

func (r *menuItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan menuItemResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }

    requestBody := map[string]interface{}{
        "name": plan.Name.ValueString(),
        "category": func() interface{} { if plan.Category.IsNull() { return nil }; return plan.Category.ValueString() }(),
        "size": func() interface{} { if plan.Size.IsNull() { return nil }; return plan.Size.ValueString() }(),
        "price": func() interface{} { if plan.Price.IsNull() { return nil }; return plan.Price.ValueFloat64() }(),
    }

    respBody, err := r.client.DoRequest("POST", "/menu_items", requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create menu item: %s", err))
        return
    }

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
        return
    }
    if id, ok := result["id"].(string); ok { plan.ID = types.StringValue(id) }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *menuItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state menuItemResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }

    respBody, err := r.client.DoRequest("GET", "/menu_items/"+state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read menu item: %s", err))
        return
    }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
        return
    }
    if v, ok := result["name"].(string); ok { state.Name = types.StringValue(v) }
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *menuItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan menuItemResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }

    body := map[string]interface{}{"name": plan.Name.ValueString()}
    _, err := r.client.DoRequest("PUT", "/menu_items/"+plan.ID.ValueString(), body)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update menu item: %s", err)); return }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *menuItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state menuItemResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }
    _, err := r.client.DoRequest("DELETE", "/menu_items/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete menu item: %s", err)); return }
}
