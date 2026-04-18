package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Config struct {
	Token types.String `tfsdk:"token"`
}
