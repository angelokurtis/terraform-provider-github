package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Model struct {
	Token types.String `tfsdk:"token"`
}
