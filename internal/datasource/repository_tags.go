package datasource

import (
	"context"

	"github.com/google/go-github/v84/github"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	t.githubClient = githubClient
}

func (t *RepositoryTag) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_repository_tags"
}

func (t *RepositoryTag) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"owner": schema.StringAttribute{
				Required:    true,
				Description: "Repository owner.",
			},
			"repo": schema.StringAttribute{
				Required:    true,
				Description: "Repository name.",
			},
			"tags": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Repository tags.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Computed: true},
						"sha":  schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (t *RepositoryTag) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	config := new(RepositoryTagModel)
	diags := req.Config.Get(ctx, config)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	var items []TagItem

	listOpts := &github.ListOptions{PerPage: 100}
	for {
		tags, resp, err := t.githubClient.Repositories.ListTags(ctx, config.Owner.ValueString(), config.Repo.ValueString(), listOpts)
		if err != nil {
			res.Diagnostics.AddError("Unable to list repository tags", err.Error())
			return
		}

		for _, tag := range tags {
			item := TagItem{
				Name: types.StringPointerValue(tag.Name),
				SHA:  types.StringPointerValue(tag.Commit.SHA),
			}
			items = append(items, item)
		}

		if resp.NextPage == 0 {
			break
		}

		listOpts.Page = resp.NextPage
	}

	tagValues, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
			"sha":  types.StringType,
		},
	}, items)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	state := &RepositoryTagModel{
		Owner: config.Owner,
		Repo:  config.Repo,
		Tags:  tagValues,
	}
	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
}

type RepositoryTagModel struct {
	Owner types.String `tfsdk:"owner"`
	Repo  types.String `tfsdk:"repo"`
	Tags  types.List   `tfsdk:"tags"`
}

type TagItem struct {
	Name types.String `tfsdk:"name"`
	SHA  types.String `tfsdk:"sha"`
}
