package tsgo

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

func (g *PackageGenerator) Generate() (string, error) {
	s := new(strings.Builder)
	filepaths := g.GoFiles
	has_func := false
	gen_decl := map[string][]*ast.GenDecl{}
	func_decl := []*ast.FuncDecl{}

	// iterate through pkg.Syntax to write collect `*ast.GenDecl` and `*ast.FuncDecl`
	for i, file := range g.pkg.Syntax {
		if g.conf.IsFileIgnored(filepaths[i]) {
			continue
		}
		gen_decl[filepaths[i]] = []*ast.GenDecl{}

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {

			// GenDecl can be an import, type, var, or const expression
			case *ast.GenDecl:
				if x.Tok == token.VAR || x.Tok == token.IMPORT {
					return false
				}
				gen_decl[filepaths[i]] = append(gen_decl[filepaths[i]], x)

				return false

			case *ast.FuncDecl:
				if g.conf.FFIBindings {
					if !has_func {
						has_func = true
					}
					func_decl = append(func_decl, x)
				}

				return false
			}
			return true
		})
	}

	// write headers for generated file for specific package
	g.writeFileCodegenHeader(s)
	if has_func {
		g.writeFFIHeaders(s)
	}
	g.writeFileFrontmatter(s)

	// iterate through pkg.Syntax to write the file
	for i, file := range g.pkg.Syntax {
		first := true
		for _, gd := range gen_decl[filepaths[i]] {
			if first {
				g.writeFileSourceHeader(s, filepaths[i], file)
				first = false
			}
			g.writeGroupDecl(s, gd)
		}
	}

	if g.conf.FFIBindings {
		path := g.conf.Path
		fmt.Println(path)
		fmt.Println(g.conf.Path)
		g.writeFFIConfig(s, func_decl, filepath.Dir(g.conf.Path))
	}

	return s.String(), nil
}
