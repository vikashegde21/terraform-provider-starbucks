package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &storeResource{}
var _ resource.ResourceWithImportState = &storeResource{}

type storeResource struct {
    client *StarbucksClient
}

type storeResourceModel struct {
    ID            types.String `tfsdk:"id"`
    Name          types.String `tfsdk:"name"`
    StoreNumber   types.String `tfsdk:"store_number"`
    Address       types.String `tfsdk:"address"`
    City          types.String `tfsdk:"city"`
    State         types.String `tfsdk:"state"`
    ZipCode       types.String `tfsdk:"zip_code"`
    Country       types.String `tfsdk:"country"`
    PhoneNumber   types.String `tfsdk:"phone_number"`
    Latitude      types.Float64 `tfsdk:"latitude"`
    Longitude     types.Float64 `tfsdk:"longitude"`
    OpeningHours  types.String `tfsdk:"opening_hours"`
    HasDriveThru  types.Bool   `tfsdk:"has_drive_thru"`
    HasWifi       types.Bool   `tfsdk:"has_wifi"`
    HasMobileOrder types.Bool  `tfsdk:"has_mobile_order"`
    Capacity      types.Int64  `tfsdk:"capacity"`
    StoreType     types.String `tfsdk:"store_type"`
    ManagerEmail  types.String `tfsdk:"manager_email"`
    Status        types.String `tfsdk:"status"`
}

func NewStoreResource() resource.Resource {
    return &storeResource{}
}

func (r *storeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_store"
}

func (r *storeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a Starbucks store location with full configuration options.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Unique identifier for the store",
                Computed:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                Description: "Store name/location description",
                Required:    true,
            },
            "store_number": schema.StringAttribute{
                Description: "Official Starbucks store number",
                Required:    true,
            },
            "address": schema.StringAttribute{
                Description: "Street address",
                Required:    true,
            },
            "city": schema.StringAttribute{
                Description: "City",
                Required:    true,
            },
            "state": schema.StringAttribute{
                Description: "State/Province",
                Required:    true,
            },
            "zip_code": schema.StringAttribute{
                Description: "ZIP/Postal code",
                Required:    true,
            },
            "country": schema.StringAttribute{
                Description: "Country code (ISO 3166-1 alpha-2)",
                Optional:    true,
            },
            "phone_number": schema.StringAttribute{
                Description: "Contact phone number",
                Required:    true,
            },
            "latitude": schema.Float64Attribute{
                Description: "Latitude coordinate",
                Optional:    true,
            },
            "longitude": schema.Float64Attribute{
                Description: "Longitude coordinate",
                Optional:    true,
            },
            "opening_hours": schema.StringAttribute{
                Description: "Store opening hours (e.g., 'Mon-Fri: 6AM-9PM, Sat-Sun: 7AM-8PM')",
                Optional:    true,
            },
            "has_drive_thru": schema.BoolAttribute{
                Description: "Whether store has drive-thru service",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(false),
            },
            "has_wifi": schema.BoolAttribute{
                Description: "Whether store offers WiFi",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(true),
            },
            "has_mobile_order": schema.BoolAttribute{
                Description: "Whether store supports mobile order & pay",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(true),
            },
            "capacity": schema.Int64Attribute{
                Description: "Maximum customer capacity",
                Optional:    true,
                Computed:    true,
                Default:     int64default.StaticInt64(50),
            },
            "store_type": schema.StringAttribute{
                Description: "Store type: standard, reserve, express, drive_thru_only",
                Optional:    true,
            },
            "manager_email": schema.StringAttribute{
                Description: "Store manager email",
                Optional:    true,
            },
            "status": schema.StringAttribute{
                Description: "Store status: active, temporarily_closed, permanently_closed",
                Computed:    true,
            },
        },
    }
}

func (r *storeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData),
        )
        return
    }

    r.client = client
}

func (r *storeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan storeResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := map[string]interface{}{
        "name":            plan.Name.ValueString(),
        "store_number":    plan.StoreNumber.ValueString(),
        "address":         plan.Address.ValueString(),
        "city":            plan.City.ValueString(),
        "state":           plan.State.ValueString(),
        "zip_code":        plan.ZipCode.ValueString(),
        "phone_number":    plan.PhoneNumber.ValueString(),
        "has_drive_thru":  plan.HasDriveThru.ValueBool(),
        "has_wifi":        plan.HasWifi.ValueBool(),
        "has_mobile_order": plan.HasMobileOrder.ValueBool(),
        "capacity":        plan.Capacity.ValueInt64(),
    }

    if !plan.Country.IsNull() {
        requestBody["country"] = plan.Country.ValueString()
    }
    if !plan.Latitude.IsNull() {
        requestBody["latitude"] = plan.Latitude.ValueFloat64()
    }
    if !plan.Longitude.IsNull() {
        requestBody["longitude"] = plan.Longitude.ValueFloat64()
    }
    if !plan.OpeningHours.IsNull() {
        requestBody["opening_hours"] = plan.OpeningHours.ValueString()
    }
    if !plan.StoreType.IsNull() {
        requestBody["store_type"] = plan.StoreType.ValueString()
    }
    if !plan.ManagerEmail.IsNull() {
        requestBody["manager_email"] = plan.ManagerEmail.ValueString()
    }

    respBody, err := r.client.DoRequest("POST", "/stores", requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create store: %s", err))
        return
    }

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
        return
    }

    if id, ok := result["id"].(string); ok {
        plan.ID = types.StringValue(id)
    }
    plan.Status = types.StringValue("active")

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *storeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state storeResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    respBody, err := r.client.DoRequest("GET", "/stores/"+state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read store: %s", err))
        return
    }

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
        return
    }

    // Update state with API response
    if val, ok := result["status"].(string); ok {
        state.Status = types.StringValue(val)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *storeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan storeResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := map[string]interface{}{
        "name":            plan.Name.ValueString(),
        "address":         plan.Address.ValueString(),
        "city":            plan.City.ValueString(),
        "state":           plan.State.ValueString(),
        "zip_code":        plan.ZipCode.ValueString(),
        "phone_number":    plan.PhoneNumber.ValueString(),
        "has_drive_thru":  plan.HasDriveThru.ValueBool(),
        "has_wifi":        plan.HasWifi.ValueBool(),
        "has_mobile_order": plan.HasMobileOrder.ValueBool(),
        "capacity":        plan.Capacity.ValueInt64(),
    }

    _, err := r.client.DoRequest("PUT", "/stores/"+plan.ID.ValueString(), requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update store: %s", err))
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *storeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state storeResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.DoRequest("DELETE", "/stores/"+state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete store: %s", err))
        return
    }
}

func (r *storeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, "id", req.ID)...)
}
