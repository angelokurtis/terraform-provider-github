//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/hashicorp/terraform-plugin-framework/provider"

	intlprovider "github.com/angelokurtis/terraform-provider-github/internal/provider"
)

// bindings defines the Wire bindings for the application layer.
// It connects interfaces to their concrete implementations and
// aggregates provider sets for dependency injection.
var bindings = wire.NewSet(
	wire.Bind(new(provider.Provider), new(*intlprovider.GitHub)),

	// Include all application-level providers.
	providers,
)
