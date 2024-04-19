// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	// "encoding/json"
	"time"
	"fmt"

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
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Namespace    types.String `tfsdk:"namespace"`
	Version      types.String `tfsdk:"version"`
	SourcePath   types.String `tfsdk:"source_path"`
	ConfigInline types.Map    `tfsdk:"config_inline"`
	ConfigLocal  types.String `tfsdk:"config_local"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}

func (r *TankaReleaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tanka_release" //req.ProviderTypeName +
}

func (r *TankaReleaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tanka release",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tanka release",
				Required:            true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The cluster endpoint / apiServer",
				Required:            true,
			},
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The cluster namespace to apply against",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("default"),
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The release version",
				Optional:            true,
			},
			"source_path": schema.StringAttribute{
				MarkdownDescription: "Path to Tanka source",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("tanka/environments/default"),
			},
			// Config inline could be allowed to have nested values - see https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes/map-nested
			"config_inline": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				// Default:  "{}",
			},
			"config_local": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("{}"),
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
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





	// config_inline_map, _ := data.ConfigInline.ToMapValue(ctx)

	// raw, err := json.Marshal(config_inline_map)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to marshal json data, got error: %s", err))
	// 	return
	// }
	// ci := string(raw[:])
	ci := ""

	cl, err := r.client.getLocalConfig(data.ConfigLocal.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Apply(data.Endpoint.ValueString(), data.Namespace.ValueString(), ci, cl, data.SourcePath.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Apply Error", fmt.Sprintf("Unable to apply tanka package, got error: %s", err))
		return
	}

	data.Id = data.Name
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

	// config_inline_map, _ := data.ConfigInline.ToMapValue(ctx)

	// raw, err := json.Marshal(config_inline_map)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to marshal json data, got error: %s", err))
	// 	return
	// }
	// ci := string(raw[:])
	ci := ""

	cl, err := r.client.getLocalConfig(data.ConfigLocal.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Apply(data.Endpoint.ValueString(), data.Namespace.ValueString(), ci, cl, data.SourcePath.ValueString())
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

	// config_inline_map, _ := data.ConfigInline.ToMapValue(ctx)

	// raw, err := json.Marshal(config_inline_map)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to marshal json data, got error: %s", err))
	// 	return
	// }
	// ci := string(raw[:])
	ci := ""

	cl, err := r.client.getLocalConfig(data.ConfigLocal.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse json data, got error: %s", err))
		return
	}

	err = r.client.Delete(data.Endpoint.ValueString(), data.Namespace.ValueString(), ci, cl, data.SourcePath.ValueString())
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
