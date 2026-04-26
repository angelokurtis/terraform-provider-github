package datasource

import (
	"context"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &RepositoryTag{}
	_ datasource.DataSourceWithConfigure = &RepositoryTag{}
)

type RepositoryTag struct {
	githubClient *github.Client
}

func NewRepositoryTag() datasource.DataSource {
	return &RepositoryTag{}
}

func (t *RepositoryTag) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	// TODO implement me
	panic("implement me")
}

func (t *RepositoryTag) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	// TODO implement me
	panic("implement me")
}

func (t *RepositoryTag) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	// TODO implement me
	panic("implement me")
}

func (t *RepositoryTag) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	// TODO implement me
	panic("implement me")
}
