package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-snyk/internal/cloudapi"
	"terraform-provider-snyk/internal/snykclient"
)

var (
	_ resource.Resource              = &SastResource{}
	_ resource.ResourceWithConfigure = &SastResource{}
)

type SastResource struct {
	client snykclient.Client
}

func NewSastResource() resource.Resource {
	return &SastResource{}
}

type sastResourceModel struct {
	Data dataResourceModel `tfsdk:"data"`
}

type dataResourceModel struct {
	Attributes attributesResourceModel `tfsdk:"attributes"`
	ID         types.String            `tfsdk:"id"`
	Type       types.String            `tfsdk:"type"`
}

type attributesResourceModel struct {
	AutofixEnabled types.Bool `tfsdk:"autofix_enabled"`
	SastEnabled    types.Bool `tfsdk:"sast_enabled"`
}

func (r *SastResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sast"
}

// Schema defines the schema for the resource.
func (r *SastResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"attributes": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"sast_enabled": schema.BoolAttribute{
								Required: true,
							},
							"autofix_enabled": schema.BoolAttribute{
								Computed: true,
							},
						},
					},
					"id": schema.StringAttribute{
						Required: true,
					},
					"type": schema.StringAttribute{
						Computed: true,
						Default:  stringdefault.StaticString("sast_settings"),
					},
				},
			},
		},
	}
}

func (r *SastResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*snykclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *snykclient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = *client
}

// Create creates the resource and sets the initial Terraform state.
func (r *SastResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sastResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request *cloudapi.SastRequest = &cloudapi.SastRequest{
		Data: cloudapi.SastDataRequest{
			Attributes: cloudapi.SastAttributesRequest{
				SastEnabled: plan.Data.Attributes.SastEnabled.ValueBool(),
			},
			Type: plan.Data.Type.ValueString(),
			ID:   plan.Data.ID.ValueString(),
		},
	}

	res, err := r.client.CloudapiClient.CreateRequest(ctx, plan.Data.ID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create request, got error: %s", err))
		return
	}

	plan = sastResourceModel{
		Data: dataResourceModel{
			Attributes: attributesResourceModel{
				AutofixEnabled: types.BoolValue(res.Data.Attributes.AutofixEnabled),
				SastEnabled:    types.BoolValue(res.Data.Attributes.SastEnabled),
			},
			ID:   types.StringValue(res.Data.ID),
			Type: types.StringValue(res.Data.Type),
		},
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *SastResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sastResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := uuid.Parse(state.Data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to parse data.id, got error: %s", err))
		return
	}

	res, err := r.client.CloudapiClient.GetRequest(ctx, state.Data.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Environment, got error: %s", err))
		return
	}

	state = sastResourceModel{
		Data: dataResourceModel{
			Attributes: attributesResourceModel{
				AutofixEnabled: types.BoolValue(res.Data.Attributes.AutofixEnabled),
				SastEnabled:    types.BoolValue(res.Data.Attributes.SastEnabled),
			},
			ID: types.StringValue(res.Data.ID),
			// Type: types.StringValue(res.Data.Type),
			Type: types.StringValue(res.Data.Type),
		},
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *SastResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var state sastResourceModel

	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request *cloudapi.SastRequest = &cloudapi.SastRequest{
		Data: cloudapi.SastDataRequest{
			Attributes: cloudapi.SastAttributesRequest{
				SastEnabled: state.Data.Attributes.SastEnabled.ValueBool(),
			},
			Type: state.Data.Type.ValueString(),
			ID:   state.Data.ID.ValueString(),
		},
	}

	res, err := r.client.CloudapiClient.CreateRequest(ctx, state.Data.ID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create request, got error: %s", err))
		return
	}

	state = sastResourceModel{
		Data: dataResourceModel{
			Attributes: attributesResourceModel{
				AutofixEnabled: types.BoolValue(res.Data.Attributes.AutofixEnabled),
				SastEnabled:    types.BoolValue(res.Data.Attributes.SastEnabled),
			},
			ID:   types.StringValue(res.Data.ID),
			Type: types.StringValue(res.Data.Type),
		},
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *SastResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sastResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
