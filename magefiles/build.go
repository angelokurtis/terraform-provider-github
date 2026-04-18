package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Prod builds the production version of the application
func (Build) Prod(ctx context.Context) error {
	mg.CtxDeps(ctx, Clean, Generate.Wire)

	env := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        runtime.GOOS,
		"GOARCH":      runtime.GOARCH,
	}
	ldFlags := []string{"-s", "-w"}
	args := []string{"build", "-ldflags", strings.Join(ldFlags, " "), "-o", "bin/terraform-provider-github", "./cmd/app/"}

	return sh.RunWithV(env, "go", args...)
}

// Dev builds the development version of the application
func (Build) Dev(ctx context.Context) error {
	mg.CtxDeps(ctx, Clean, Generate.Wire)

	args := []string{"build", "-o", "bin/terraform-provider-github", "./cmd/app/"}

	return sh.RunV("go", args...)
}

// Clean removes the "bin" directory to ensure a fresh build
func Clean() error {
	directory := "./bin/"

	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil
	}

	// Try to remove the directory
	if err := os.RemoveAll(directory); err != nil {
		return fmt.Errorf("failed to delete %s: %w", directory, err)
	}

	return nil
}
