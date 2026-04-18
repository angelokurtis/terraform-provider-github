package term

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewGlamourTermRenderer,
	NewMarkdownRenderer,
)
