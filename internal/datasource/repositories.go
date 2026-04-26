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
	User  types.String   `tfsdk:"user"`
	Org   types.String   `tfsdk:"org"`
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
			"Expected *github.Client",
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
			"user": schema.StringAttribute{
				Optional:    true,
				Description: "GitHub username to list repositories for.",
			},
			"org": schema.StringAttribute{
				Optional:    true,
				Description: "GitHub organization to list repositories for.",
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
	config := new(RepositoryModel)
	diags := req.Config.Get(ctx, config)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	if !config.User.IsNull() && !config.Org.IsNull() {
		res.Diagnostics.AddError(
			"Invalid configuration",
			"Only one of 'user' or 'org' can be specified.",
		)

		return
	}

	var allRepos []types.String

	listOpts := &github.ListOptions{PerPage: 100}

	switch {
	case !config.User.IsNull():
		opt := &github.RepositoryListByUserOptions{ListOptions: *listOpts}
		for {
			repos, resp, err := r.githubClient.Repositories.ListByUser(ctx, config.User.ValueString(), opt)
			if err != nil {
				res.Diagnostics.AddError("Unable to list user repositories", err.Error())
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
	case !config.Org.IsNull():
		opt := &github.RepositoryListByOrgOptions{ListOptions: *listOpts}
		for {
			repos, resp, err := r.githubClient.Repositories.ListByOrg(ctx, config.Org.ValueString(), opt)
			if err != nil {
				res.Diagnostics.AddError("Unable to list organization repositories", err.Error())
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
	default:
		opt := &github.RepositoryListByAuthenticatedUserOptions{
			ListOptions: *listOpts,
		}
		for {
			repos, resp, err := r.githubClient.Repositories.ListByAuthenticatedUser(ctx, opt)
			if err != nil {
				res.Diagnostics.AddError("Unable to list authenticated user repositories", err.Error())
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
	}

	state := new(RepositoryModel)
	state.User = config.User
	state.Org = config.Org
	state.Repos = allRepos
	diags = res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}
