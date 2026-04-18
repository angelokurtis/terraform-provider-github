// Code generated. DO NOT EDIT.

//go:build wireinject

package main

import (
	"github.com/google/wire"

	provider "github.com/angelokurtis/terraform-provider-github/internal/provider"
)

// providers groups together all constructors needed to build the
// application's core dependencies. It tells Wire how to instantiate key
// components so they can be injected wherever required.
var providers = wire.NewSet(
	provider.Providers,
)
