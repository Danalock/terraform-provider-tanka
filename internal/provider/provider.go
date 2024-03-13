// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

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

// provider "helm" {
//   kubernetes {
//     host                   = data.aws_eks_cluster.this.endpoint
//     cluster_ca_certificate = base64decode(data.aws_eks_cluster.this.certificate_authority[0].data)
//     token                  = data.aws_eks_cluster_auth.tools.token
//   }
// }

func (p *TankaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tanka"
	resp.Version = p.version
}

func (p *TankaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The kubernetes cluster endpoint",
				Required:            true,
			},
			"cluster_certificate_authority": schema.StringAttribute{
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

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
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
