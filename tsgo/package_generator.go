package tsgo

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func (g *PackageGenerator) Generate() (string, error) {
	s := new(strings.Builder)

	g.writeFileCodegenHeader(s)
	g.writeFileFrontmatter(s)

	filepaths := g.GoFiles

	for i, file := range g.pkg.Syntax {
		if g.conf.IsFileIgnored(filepaths[i]) {
			continue
		}

		first := true

		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {

			// GenDecl can be an import, type, var, or const expression
			case *ast.GenDecl:
				if x.Tok == token.VAR || x.Tok == token.IMPORT {
					return false
				}

				if first {
					g.writeFileSourceHeader(s, filepaths[i], file)
					first = false
				}

				g.writeGroupDecl(s, x)
				return false

			case *ast.FuncDecl:
				// TODO: enable generating Bun FFI wrapper + CGo (and possibly CGo methods w/ callback)
				fmt.Println("Case: *ast.FuncDecl - ", "x.Name:", x.Name.Name, "x.Body:", x.Body, "x.Type:", x.Type, "x.Recv", x.Recv)
			}
			return true

		})

	}

	return s.String(), nil
}
