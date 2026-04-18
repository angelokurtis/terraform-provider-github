package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Wire runs Wire to generate dependency injection code.
func (Generate) Wire(ctx context.Context) error {
	mg.CtxDeps(ctx, Generate.Templates)

	dirs, err := scanWireDirs("./")
	if err != nil {
		return err
	}

	// Run Wire codegen tool on each directory found.
	for _, dir := range dirs {
		if err = sh.RunV("go", "tool", "-modfile=magefiles/go.mod", "wire", dir); err != nil {
			return err
		}
	}

	return nil
}

// Templates generates a Go source file containing Wire provider sets.
func (Generate) Templates() error {
	providers, err := scanWireProviders()
	if err != nil {
		return err
	}

	tpl, err := template.New("wire_provs.tmpl").
		Funcs(sprig.FuncMap()).
		ParseFiles("hack/wire_provs.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	out, err := os.Create("cmd/app/wire_provs.go")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer out.Close()

	// Generate providers file from template.
	if err = tpl.Execute(out, map[string][]*ProviderSet{"Providers": providers}); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Run goimports to format and fix imports in the generated file.
	return sh.RunV("go", "tool", "-modfile=magefiles/go.mod", "goimports", "-w", "cmd/app/wire_provs.go")
}
