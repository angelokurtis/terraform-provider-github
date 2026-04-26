package datasource

import (
	"context"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &Tag{}
	_ datasource.DataSourceWithConfigure = &Tag{}
)

type Tag struct {
	githubClient *github.Client
}

func NewTag() datasource.DataSource {
	return &Tag{}
}

func (t *Tag) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	//TODO implement me
	panic("implement me")
}

func (t *Tag) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}
