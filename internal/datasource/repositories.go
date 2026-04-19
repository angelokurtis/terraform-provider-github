package datasource

import (
	"context"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &Repository{}
	_ datasource.DataSourceWithConfigure = &Repository{}
)

type RepositoryModel struct {
	ID    types.String   `tfsdk:"id"`
	Repos []types.String `tfsdk:"repos"`
}

type Repository struct {
	githubClient *github.Client
}

func NewRepository() datasource.DataSource {
	return &Repository{}
}

func (r *Repository) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	githubClient, ok := req.ProviderData.(*github.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *provider.GitHubClient",
		)

		return
	}

	r.githubClient = githubClient
}

func (r *Repository) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_repositories"
}

func (r *Repository) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Static ID for this data source.",
			},
			"repos": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of repository names.",
			},
		},
	}
}

func (r *Repository) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var state RepositoryModel

	// List repositories for the authenticated user
	opt := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []types.String

	for {
		repos, resp, err := r.githubClient.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			res.Diagnostics.AddError(
				"Unable to list repositories",
				err.Error(),
			)

			return
		}

		for _, repo := range repos {
			if repo.Name != nil {
				allRepos = append(allRepos, types.StringValue(*repo.Name))
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	state.ID = types.StringValue("repositories")
	state.Repos = allRepos

	diags := res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}
