package datasource

import (
	"context"
	"time"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &Repository{}
	_ datasource.DataSourceWithConfigure = &Repository{}
)

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
			"repos": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of repositories.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":           schema.StringAttribute{Computed: true},
						"full_name":      schema.StringAttribute{Computed: true},
						"description":    schema.StringAttribute{Computed: true},
						"default_branch": schema.StringAttribute{Computed: true},
						"created_at":     schema.StringAttribute{Computed: true},
						"pushed_at":      schema.StringAttribute{Computed: true},
						"updated_at":     schema.StringAttribute{Computed: true},
						"archived":       schema.BoolAttribute{Computed: true},
						"visibility":     schema.StringAttribute{Computed: true},
					},
				},
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

	var items []RepositoryItem

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
				item := NewRepositoryItem(repo)
				items = append(items, item)
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
				item := NewRepositoryItem(repo)
				items = append(items, item)
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
				item := NewRepositoryItem(repo)
				items = append(items, item)
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}

	repoValues, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":           types.StringType,
			"full_name":      types.StringType,
			"description":    types.StringType,
			"default_branch": types.StringType,
			"created_at":     types.StringType,
			"pushed_at":      types.StringType,
			"updated_at":     types.StringType,
			"archived":       types.BoolType,
			"visibility":     types.StringType,
		},
	}, items)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	state := new(RepositoryModel)
	state.User = config.User
	state.Org = config.Org
	state.Repos = repoValues

	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
}

type RepositoryModel struct {
	User  types.String `tfsdk:"user"`
	Org   types.String `tfsdk:"org"`
	Repos types.List   `tfsdk:"repos"`
}

type RepositoryItem struct {
	Name          types.String `tfsdk:"name"`
	FullName      types.String `tfsdk:"full_name"`
	Description   types.String `tfsdk:"description"`
	DefaultBranch types.String `tfsdk:"default_branch"`
	CreatedAt     types.String `tfsdk:"created_at"`
	PushedAt      types.String `tfsdk:"pushed_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	Archived      types.Bool   `tfsdk:"archived"`
	Visibility    types.String `tfsdk:"visibility"`
}

func NewRepositoryItem(repo *github.Repository) RepositoryItem {
	item := RepositoryItem{
		Name:          types.StringPointerValue(repo.Name),
		FullName:      types.StringPointerValue(repo.FullName),
		Description:   types.StringPointerValue(repo.Description),
		DefaultBranch: types.StringPointerValue(repo.DefaultBranch),
		Archived:      types.BoolValue(repo.GetArchived()),
		Visibility:    types.StringValue(repo.GetVisibility()),
	}

	if repo.CreatedAt != nil {
		item.CreatedAt = types.StringValue(repo.CreatedAt.Time.Format(time.RFC3339))
	} else {
		item.CreatedAt = types.StringNull()
	}

	if repo.PushedAt != nil {
		item.PushedAt = types.StringValue(repo.PushedAt.Time.Format(time.RFC3339))
	} else {
		item.PushedAt = types.StringNull()
	}

	if repo.UpdatedAt != nil {
		item.UpdatedAt = types.StringValue(repo.UpdatedAt.Time.Format(time.RFC3339))
	} else {
		item.UpdatedAt = types.StringNull()
	}

	return item
}
