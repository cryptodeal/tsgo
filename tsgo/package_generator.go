package tsgo

import (
	"go/ast"
	"go/token"
	"strings"
)

func (g *PackageGenerator) Generate() (string, error) {
	cg := new(strings.Builder)
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
		path := strings.Split(g.conf.Path, "/")
		g.writeFFIConfig(s, func_decl, path[len(path)-1])
		g.writeCGo(cg, func_decl, path[len(path)-1])
	}

	return s.String(), nil
}
