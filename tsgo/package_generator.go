package tsgo

import (
	"go/ast"
	"go/token"
	"strings"
)

func (g *PackageGenerator) Generate() (string, error) {
	s := new(strings.Builder)
	filepaths := g.GoFiles

	for i, file := range g.pkg.Syntax {
		if g.conf.IsFileIgnored(filepaths[i]) {
			continue
		}

		first := true
		has_func := false

		gen_decl := []*ast.GenDecl{}
		func_decl := []*ast.FuncDecl{}

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {

			// GenDecl can be an import, type, var, or const expression
			case *ast.GenDecl:
				if x.Tok == token.VAR || x.Tok == token.IMPORT {
					return false
				}
				gen_decl = append(gen_decl, x)

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
		g.writeFileCodegenHeader(s)
		if has_func {
			g.writeFFIHeaders(s)
		}
		g.writeFileFrontmatter(s)

		for _, gd := range gen_decl {
			if first {
				g.writeFileSourceHeader(s, filepaths[i], file)
				first = false
			}
			g.writeGroupDecl(s, gd)
		}

		if g.conf.FFIBindings && has_func {
			for _, fd := range func_decl {
				g.writeFuncDecl(s, fd)
			}
		}
	}

	return s.String(), nil
}
