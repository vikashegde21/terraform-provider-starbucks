package main

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &employeeResource{}

type employeeResource struct {
    client *StarbucksClient
}

type employeeResourceModel struct {
    ID               types.String  `tfsdk:"id"`
    EmployeeNumber   types.String  `tfsdk:"employee_number"`
    FirstName        types.String  `tfsdk:"first_name"`
    LastName         types.String  `tfsdk:"last_name"`
    Email            types.String  `tfsdk:"email"`
    PhoneNumber      types.String  `tfsdk:"phone_number"`
    StoreID          types.String  `tfsdk:"store_id"`
    Position         types.String  `tfsdk:"position"`
    HireDate         types.String  `tfsdk:"hire_date"`
    HourlyRate       types.Float64 `tfsdk:"hourly_rate"`
    IsBarista        types.Bool    `tfsdk:"is_barista"`
    IsShiftSupervisor types.Bool   `tfsdk:"is_shift_supervisor"`
    IsCertified      types.Bool    `tfsdk:"is_certified"`
    AvailableHours   types.String  `tfsdk:"available_hours"`
    EmploymentType   types.String  `tfsdk:"employment_type"`
    Status           types.String  `tfsdk:"status"`
}

func NewEmployeeResource() resource.Resource {
    return &employeeResource{}
}

func (r *employeeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_employee"
}

func (r *employeeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a Starbucks employee (partner) with full employee lifecycle.",
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Description: "Unique identifier",
                Computed:    true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "employee_number": schema.StringAttribute{
                Description: "Unique employee/partner number",
                Required:    true,
            },
            "first_name": schema.StringAttribute{
                Description: "First name",
                Required:    true,
            },
            "last_name": schema.StringAttribute{
                Description: "Last name",
                Required:    true,
            },
            "email": schema.StringAttribute{
                Description: "Email address",
                Required:    true,
            },
            "phone_number": schema.StringAttribute{
                Description: "Contact phone number",
                Optional:    true,
            },
            "store_id": schema.StringAttribute{
                Description: "ID of the assigned store",
                Required:    true,
            },
            "position": schema.StringAttribute{
                Description: "Job position: barista, shift_supervisor, store_manager, assistant_manager",
                Required:    true,
            },
            "hire_date": schema.StringAttribute{
                Description: "Hire date (YYYY-MM-DD format)",
                Required:    true,
            },
            "hourly_rate": schema.Float64Attribute{
                Description: "Hourly pay rate (USD)",
                Optional:    true,
                Sensitive:   true,
            },
            "is_barista": schema.BoolAttribute{
                Description: "Whether employee is a certified barista",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(true),
            },
            "is_shift_supervisor": schema.BoolAttribute{
                Description: "Whether employee is a shift supervisor",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(false),
            },
            "is_certified": schema.BoolAttribute{
                Description: "Whether employee completed all certifications",
                Optional:    true,
                Computed:    true,
                Default:     booldefault.StaticBool(false),
            },
            "available_hours": schema.StringAttribute{
                Description: "Available working hours (e.g., 'Mon-Fri: 9AM-5PM')",
                Optional:    true,
            },
            "employment_type": schema.StringAttribute{
                Description: "Employment type: full_time, part_time, seasonal",
                Optional:    true,
            },
            "status": schema.StringAttribute{
                Description: "Employment status",
                Computed:    true,
            },
        },
    }
}

func (r *employeeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    client, ok := req.ProviderData.(*StarbucksClient)
    if !ok {
        resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *StarbucksClient, got: %T", req.ProviderData))
        return
    }
    r.client = client
}

func (r *employeeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan employeeResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := map[string]interface{}{
        "employee_number":     plan.EmployeeNumber.ValueString(),
        "first_name":          plan.FirstName.ValueString(),
        "last_name":           plan.LastName.ValueString(),
        "email":               plan.Email.ValueString(),
        "store_id":            plan.StoreID.ValueString(),
        "position":            plan.Position.ValueString(),
        "hire_date":           plan.HireDate.ValueString(),
        "is_barista":          plan.IsBarista.ValueBool(),
        "is_shift_supervisor": plan.IsShiftSupervisor.ValueBool(),
        "is_certified":        plan.IsCertified.ValueBool(),
    }

    respBody, err := r.client.DoRequest("POST", "/employees", requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create employee: %s", err))
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

func (r *employeeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state employeeResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    respBody, err := r.client.DoRequest("GET", "/employees/"+state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read employee: %s", err))
        return
    }

    var result map[string]interface{}
    if err := json.Unmarshal(respBody, &result); err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse response: %s", err))
        return
    }

    if val, ok := result["status"].(string); ok {
        state.Status = types.StringValue(val)
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *employeeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan employeeResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    requestBody := map[string]interface{}{
        "first_name": plan.FirstName.ValueString(),
        "last_name":  plan.LastName.ValueString(),
        "email":      plan.Email.ValueString(),
        "phone_number": func() interface{} {
            if plan.PhoneNumber.IsNull() { return nil }
            return plan.PhoneNumber.ValueString()
        }(),
        "position": plan.Position.ValueString(),
    }

    _, err := r.client.DoRequest("PUT", "/employees/"+plan.ID.ValueString(), requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update employee: %s", err))
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *employeeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state employeeResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.DoRequest("DELETE", "/employees/"+state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete employee: %s", err))
        return
    }
}
