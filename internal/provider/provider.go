// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure TankaProvider satisfies various provider interfaces.
var _ provider.Provider = &TankaProvider{}
var _ provider.ProviderWithFunctions = &TankaProvider{}

// TankaProvider defines the provider implementation.
type TankaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TankaProviderModel describes the provider data model.
type TankaProviderModel struct {
	Endpoint             types.String `tfsdk:"endpoint"`
	ClusterCaCertificate types.String `tfsdk:"cluster_ca_certificate"`
	Token                types.String `tfsdk:"token"`
}

func (p *TankaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tanka"
	resp.Version = p.version
}

func (p *TankaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The kubernetes cluster endpoint / the API server",
				Required:            true,
			},
			"cluster_ca_certificate": schema.StringAttribute{
				MarkdownDescription: "The certificate-authority for the cluster",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token for the user entry in kubeconfig",
				Required:            true,
			},
		},
	}
}

func (p *TankaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TankaProviderModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := data.Endpoint.ValueString()
	token := data.Token.ValueString()
	cluster_ca_certificate := data.ClusterCaCertificate.ValueString()

	client, err := NewClient(&endpoint, &token, &cluster_ca_certificate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Tanka API Client",
			"An unexpected error occurred when creating the Tanka API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Tanka Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *TankaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTankaReleaseResource,
	}
}

func (p *TankaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *TankaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TankaProvider{
			version: version,
		}
	}
}
