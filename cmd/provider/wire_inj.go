//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func NewProvider(ctx context.Context) (provider.Provider, func(), error) {
	wire.Build(bindings)
	return nil, nil, nil
}
