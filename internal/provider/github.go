package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &GitHub{}

type GitHub struct{}

func NewGitHub() *GitHub {
	return &GitHub{}
}

func (g *GitHub) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHub) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHub) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (g *GitHub) DataSources(ctx context.Context) []func() datasource.DataSource {
	//TODO implement me
	panic("implement me")
}

func (g *GitHub) Resources(ctx context.Context) []func() resource.Resource {
	//TODO implement me
	panic("implement me")
}
