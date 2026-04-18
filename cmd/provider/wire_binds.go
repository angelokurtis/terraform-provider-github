package main

import (
	terraformprovidergithub "github-personal/angelokurtis/terraform-provider-github"
	"github-personal/angelokurtis/terraform-provider-github/internal/term"
	"github.com/google/wire"
)

// bindings defines the Wire bindings for the application layer.
// It connects interfaces to their concrete implementations and
// aggregates provider sets for dependency injection.
var bindings = wire.NewSet(
	wire.Bind(new(Runner), new(*terraformprovidergithub.Runner)),
	wire.Bind(new(term.Renderer), new(*term.MarkdownRenderer)),

	// Include all application-level providers.
	providers,
)
