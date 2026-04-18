package main

import "github.com/magefile/mage/sh"

// All runs linting on the entire codebase.
func (Lint) All() error {
	if err := sh.RunV("go", "tool", "-modfile=magefiles/go.mod", "golangci-lint", "run"); err != nil {
		return err
	}

	return nil
}

// Changes runs linting only on modified files.
func (Lint) Changes() error {
	if err := sh.RunV("go", "tool", "-modfile=magefiles/go.mod", "golangci-lint", "run", "--new"); err != nil {
		return err
	}

	return nil
}
