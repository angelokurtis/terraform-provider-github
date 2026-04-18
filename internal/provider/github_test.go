package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const providerConfig = `
provider "github" {}
`

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"github": providerserver.NewProtocol6WithError(NewGitHub()),
}

func TestGitHub_Configure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig,
			},
		},
	})
}
