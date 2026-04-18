//go:build wireinject

package github

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewClient,
)
