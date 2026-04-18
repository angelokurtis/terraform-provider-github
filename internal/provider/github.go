package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &GitHub{}

type GitHub struct{}

func NewGitHub() *GitHub {
	return &GitHub{}
}

func (g *GitHub) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "hashicups"
	response.Version = "dev"
}

func (g *GitHub) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{}
}

func (g *GitHub) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring GitHub client")

	// TODO: implement me
}

func (g *GitHub) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (g *GitHub) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
