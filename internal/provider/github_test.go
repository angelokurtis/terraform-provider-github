package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGitHub_Configure(t *testing.T) {
	tests := []struct {
		name  string
		steps []resource.TestStep
	}{
		{
			name: "missing token",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {}
data "github_repositories" "_" {}
`,
					ExpectError: regexp.MustCompile(`The argument "token" is required`),
				},
			},
		},
		{
			name: "null token",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {
  token = null
}
data "github_repositories" "_" {}
`,
					ExpectError: regexp.MustCompile(`Missing Configuration for Required Attribute`),
				},
			},
		},
		{
			name: "empty token",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {
  token = ""
}
data "github_repositories" "_" {}
`,
					ExpectError: regexp.MustCompile("Empty GitHub Token"),
				},
			},
		},
		{
			name: "valid token only",
			steps: []resource.TestStep{
				{
					Config: `
provider "github" {
  token = "dummy-token"
}

data "github_repositories" "test" {}
`,
				},
			},
		},
	}

	ProtoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"github": providerserver.NewProtocol6WithError(NewGitHub()),
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
