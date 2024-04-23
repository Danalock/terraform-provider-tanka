// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	// "encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TankaReleaseResource{}
var _ resource.ResourceWithImportState = &TankaReleaseResource{}

func NewTankaReleaseResource() resource.Resource {
	return &TankaReleaseResource{}
}

// TankaReleaseResource defines the resource implementation.
type TankaReleaseResource struct {
	client *Client
}

// TankaReleaseResourceModel describes the resource data model.
type TankaReleaseResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Namespace      types.String `tfsdk:"namespace"`
	Version        types.String `tfsdk:"version"`
	SourcePath     types.String `tfsdk:"source_path"`
	Config         types.String `tfsdk:"config"`
	ConfigOverride types.String `tfsdk:"config_override"`
	LastUpdated    types.String `tfsdk:"last_updated"`
}

func (r *TankaReleaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tanka_release" //req.ProviderTypeName +
}

func (r *TankaReleaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tanka release",

		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The Kubernetes namespace to install the release into. Defaults to `default`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "A version number for the Tanka package. Examples could be a git commit SHA, or a random value to force update on every run. This value is not passed to the tanka application, if version information needs to be available to tanka it should be set as a subkey in one of the config objects.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("0"),
			},
			"source_path": schema.StringAttribute{
				MarkdownDescription: "The location of the Tanka main file. Defaults to `tanka/environments/default`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("tanka/environments/default"),
			},
			"config": schema.StringAttribute{
				MarkdownDescription: "Configuration object in arbitrary JSON format. The data can be provided inline with jsonencode() or given as a file. Local file paths are prefixed with `file://` and remote sources with the correct protocol `http://` or `https://`. Remote sources must be publicly available. Defaults to the empty object.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("{}"),
			},
			"config_override": schema.StringAttribute{
				MarkdownDescription: "Configuration override object in arbitrary JSON format. The data can be provided inline with jsonencode() or given as a file. Local file paths are prefixed with `file://` and remote sources with the correct protocol `http://` or `https://`. Remote sources must be publicly available. Defaults to the empty object.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("{}"),
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "Timestamp updated on every apply operation.",
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the resource. Consists of the cluster endpoint suffixed with a six letter random string (underscore separated).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *TankaReleaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *TankaReleaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TankaReleaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.parseConfig(data.Config.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	config_override, err := r.client.parseConfig(data.ConfigOverride.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Apply(r.client.Endpoint, data.Namespace.ValueString(), config, config_override, data.SourcePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Apply Error", fmt.Sprintf("Unable to apply tanka package, got error: %s", err))
		return
	}

	data.Id = types.StringValue(r.client.Endpoint + "_" + randSeq(6))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TankaReleaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TankaReleaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TankaReleaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TankaReleaseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.parseConfig(data.Config.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	config_override, err := r.client.parseConfig(data.ConfigOverride.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Apply(r.client.Endpoint, data.Namespace.ValueString(), config, config_override, data.SourcePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Apply Error", fmt.Sprintf("Unable to apply tanka package, got error: %s", err))
		return
	}

	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "updated a resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TankaReleaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TankaReleaseResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.parseConfig(data.Config.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	config_override, err := r.client.parseConfig(data.ConfigOverride.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Delete(r.client.Endpoint, data.Namespace.ValueString(), config, config_override, data.SourcePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete Error", fmt.Sprintf("Unable to delete tanka package, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "deleted a resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TankaReleaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
