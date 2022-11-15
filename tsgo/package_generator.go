package tsgo

import (
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
				if first {
					g.writeFileSourceHeader(s, filepaths[i], file)
					first = false
				}
				g.writeFuncDecl(s, x)
				return false

			case *ast.FuncLit:
				if first {
					g.writeFileSourceHeader(s, filepaths[i], file)
					first = false
				}
				g.writeFuncLit(s, x)
				return false

			case *ast.FuncType:
				if first {
					g.writeFileSourceHeader(s, filepaths[i], file)
					first = false
				}
				g.writeFuncType(s, x)
				return false
			}

			return true

		})

	}

	return s.String(), nil
}
