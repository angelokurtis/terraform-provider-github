//go:build wireinject

package provider

import (
	"github.com/google/wire"
)

var providers = wire.NewSet(
	NewGitHub,
)
