package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

// scanWireDirs finds directories with Go files using wire.Build.
func scanWireDirs(rootDir string) ([]string, error) {
	var dirs []string

	visited := make(map[string]bool)

	// Walk through the directory tree to find Wire Build calls.
	err := filepath.Walk(rootDir, func(filePath string, fileInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if fileInfo.IsDir() || filepath.Ext(filePath) != ".go" {
			return nil
		}

		fs := token.NewFileSet()

		node, parseErr := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
		if parseErr != nil {
			log.Printf("Error parsing file: %s, %v", filePath, parseErr)
			return nil
		}

		ast.Inspect(node, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if selector, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if selector.Sel.Name == "Build" {
						if ident, ok := selector.X.(*ast.Ident); ok && ident.Name == "wire" {
							dir := fmt.Sprintf("./%s/", filepath.Dir(filePath))
							if !visited[dir] {
								visited[dir] = true

								dirs = append(dirs, dir)
							}

							return false
						}
					}
				}
			}

			return true
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}

// ProviderSet represents a wire.ProviderSet global variable.
type ProviderSet struct {
	PkgPath string
	PkgName string
	Name    string
}

// scanWireProviders finds all global wire.ProviderSet variables in the project.
func scanWireProviders() ([]*ProviderSet, error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedDeps | packages.NeedName,
	}

	// Load all packages in the module
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	var result []*ProviderSet

	// Inspect all definitions in each package
	for _, pkg := range pkgs {
		for ident, obj := range pkg.TypesInfo.Defs {
			if obj == nil {
				continue
			}

			// Check if definition is a global variable
			if v, ok := obj.(*types.Var); ok && obj.Parent() == obj.Pkg().Scope() {
				// Check if variable is of type wire.ProviderSet
				if named, ok := v.Type().(*types.Named); ok &&
					named.Obj().Pkg() != nil &&
					named.Obj().Pkg().Path() == "github.com/google/wire" &&
					named.Obj().Name() == "ProviderSet" {
					result = append(result, &ProviderSet{
						PkgPath: pkg.PkgPath,
						PkgName: pkg.Name,
						Name:    ident.Name,
					})
				}
			}
		}
	}

	return result, nil
}
