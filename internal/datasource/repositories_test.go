package datasource_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	intlhttp "github.com/angelokurtis/terraform-provider-github/internal/http"
	"github.com/angelokurtis/terraform-provider-github/internal/provider"
)

func TestRepository_Read(t *testing.T) {
	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "should return all repositories with correct names",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {
  token = "dummy-token"
}
data "github_repositories" "test" {
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.github_repositories.test", "repos.#", "2"),
						resource.TestCheckResourceAttr(
							"data.github_repositories.test",
							"repos.0.name",
							"git-push-and-pray",
						),
						resource.TestCheckResourceAttr(
							"data.github_repositories.test",
							"repos.1.name",
							"taco-driven-development",
						),
					),
				},
			},
		},
	}
	hook := func(i *cassette.Interaction) error {
		for k, v := range i.Request.Headers {
			if strings.EqualFold(k, "Authorization") && len(v) > 0 {
				value := "dummy-token"
				if strings.HasPrefix(strings.ToLower(v[0]), "bearer ") {
					value = "Bearer dummy-token"
				}

				i.Request.Headers[k] = []string{value}
			}
		}

		return nil
	}

	r, err := recorder.New(
		filepath.Join("testdata", strings.ReplaceAll(t.Name(), "/", "_")),
		recorder.WithHook(hook, recorder.AfterCaptureHook),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := r.Stop(); err != nil {
			t.Error(err)
		}
	})

	httpClient := r.GetDefaultClient()
	httpClient.Transport = &intlhttp.RequestLogger{DefaultTransport: httpClient.Transport}
	ProtoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"github": providerserver.NewProtocol6WithError(provider.NewGitHub(httpClient)),
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: ProtoV6ProviderFactories,
				Steps:                    tc.steps,
			})
		})
	}
}
