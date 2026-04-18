package main

import "github.com/magefile/mage/mg"

// Aliases maps command names to their corresponding actions.
// It allows using shorthand names for tasks, e.g., "gen" for Generate and "lint" for Lint.All.
var Aliases = map[string]any{
	"build": Build.Prod,
	"gen":   Generate.Wire,
	"lint":  Lint.All,
}

// Build defines the set of tasks related to the application's build process.
type Build mg.Namespace

// Lint defines the set of tasks related to linting the codebase.
type Lint mg.Namespace

// Generate defines the set of tasks related to code generation.
type Generate mg.Namespace
