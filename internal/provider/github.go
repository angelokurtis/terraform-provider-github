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

func (g *GitHub) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "github"
	res.Version = "dev"
}

func (g *GitHub) Schema(ctx context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "GitHub API token used for authentication.",
			},
		},
	}
}

func (g *GitHub) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring GitHub client")

	config := new(Config)
	diags := req.Config.Get(ctx, config)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}
}

func (g *GitHub) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (g *GitHub) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
