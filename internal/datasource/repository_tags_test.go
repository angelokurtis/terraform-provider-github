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

func TestRepositoryTags_Read(t *testing.T) {
	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "should return repository tags with name and sha",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {
  token = "dummy-token"
}
data "github_repository_tags" "test" {
  owner = "kubernetes"
  repo  = "kubernetes"
}
`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.github_repository_tags.test", "tags.#"),
						resource.TestCheckResourceAttrSet("data.github_repository_tags.test", "tags.0.name"),
						resource.TestCheckResourceAttrSet("data.github_repository_tags.test", "tags.0.sha"),
						resource.TestCheckTypeSetElemNestedAttrs(
							"data.github_repository_tags.test",
							"tags.*",
							map[string]string{
								"name": "v1.36.0",
								"sha":  "ecf6decece6a6de25a57aad9ba90b6ce580f6f78",
							},
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
	protoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"github": providerserver.NewProtocol6WithError(provider.NewGitHub(httpClient)),
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: protoV6ProviderFactories,
				Steps:                    tc.steps,
			})
		})
	}
}
