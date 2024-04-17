// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	// "encoding/json"
	"fmt"
	// "io"
	// "net/http"
	// "os"
	// "strings"
	// "time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	// "github.com/grafana/tanka/pkg/jsonnet"
	// "github.com/grafana/tanka/pkg/tanka"
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
			// "api_server": schema.StringAttribute{
			// 	Required: true,
			// },
			// "namespace": schema.StringAttribute{
			// 	Optional: true,
			// 	Default:  stringdefault.StaticString("default"),
			// },
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
			// "id": schema.StringAttribute{
			// 	Computed:            true,
			// 	MarkdownDescription: "Example identifier",
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.UseStateForUnknown(),
			// 	},
			// },
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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// name := fmt.Sprintf("%v", d.Get("name"))

	// api_server := fmt.Sprintf("%v", d.Get("api_server"))
	// namespace := fmt.Sprintf("%v", d.Get("namespace"))
	// source_path := fmt.Sprintf("%v", d.Get("source_path"))
	// config_inline := d.Get("config_inline").(map[string]any)
	// config_local := fmt.Sprintf("%v", d.Get("config_local"))

	// raw, err := json.Marshal(config_inline)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// ci := string(raw[:])

	// cl, err := getLocalConfig(config_local)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// err = tankaApply(api_server, namespace, ci, cl, source_path)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId(name)
	// d.Set("last_updated", time.Now().Format(time.RFC3339))

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

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

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *TankaReleaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
