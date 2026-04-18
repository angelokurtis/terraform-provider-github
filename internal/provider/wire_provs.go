//go:build wireinject

package provider

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewGitHub,
)
