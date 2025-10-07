package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type promotionResource struct { client *StarbucksClient }

type promotionResourceModel struct {
    ID          types.String `tfsdk:"id"`
    Name        types.String `tfsdk:"name"`
    Description types.String `tfsdk:"description"`
    StartDate   types.String `tfsdk:"start_date"`
    EndDate     types.String `tfsdk:"end_date"`
    Active      types.Bool   `tfsdk:"active"`
}

func NewPromotionResource() resource.Resource { return &promotionResource{} }

func (r *promotionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_promotion"
}

func (r *promotionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages promotional campaigns.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{Computed: true},
            "name": schema.StringAttribute{Required: true},
            "description": schema.StringAttribute{Optional: true},
            "start_date": schema.StringAttribute{Optional: true},
            "end_date": schema.StringAttribute{Optional: true},
            "active": schema.BoolAttribute{Optional: true, Computed: true},
        },
    }
}

func (r *promotionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil { return }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok { resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData)); return }
    r.client = client
}

func (r *promotionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan promotionResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }
    body := map[string]interface{}{"name": plan.Name.ValueString(), "description": func() interface{} { if plan.Description.IsNull() { return nil }; return plan.Description.ValueString() }()}
    respBody, err := r.client.DoRequest("POST", "/promotions", body)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create promotion: %s", err)); return }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }
    if id, ok := result["id"].(string); ok { plan.ID = types.StringValue(id) }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *promotionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state promotionResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }
    respBody, err := r.client.DoRequest("GET", "/promotions/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read promotion: %s", err)); return }
    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil { resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err)); return }
    if a, ok := result["active"].(bool); ok { state.Active = types.BoolValue(a) }
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *promotionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan promotionResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() { return }
    body := map[string]interface{}{"name": plan.Name.ValueString()}
    _, err := r.client.DoRequest("PUT", "/promotions/"+plan.ID.ValueString(), body)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update promotion: %s", err)); return }
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *promotionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state promotionResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() { return }
    _, err := r.client.DoRequest("DELETE", "/promotions/"+state.ID.ValueString(), nil)
    if err != nil { resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete promotion: %s", err)); return }
}
