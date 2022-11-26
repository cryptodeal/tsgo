package tsgo

import (
	"fmt"
	"go/ast"
	"strings"
)

func (g *PackageGenerator) writeCommentGroupIfNotNil(s *strings.Builder, f *ast.CommentGroup, depth int) {
	if f != nil {
		g.writeCommentGroup(s, f, depth)
	}
}

func (g *PackageGenerator) writeCommentGroup(s *strings.Builder, f *ast.CommentGroup, depth int) {
	docLines := strings.Split(f.Text(), "\n")

	if depth != 0 {
		g.writeIndent(s, depth)
	}
	s.WriteString("/**\n")

	for _, c := range docLines {
		if len(strings.TrimSpace(c)) == 0 {
			continue
		}
		g.writeIndent(s, depth)
		s.WriteString(fmt.Sprintf(" * %s\n", strings.ReplaceAll(c, "*/", "*\\/")))
	}
	g.writeIndent(s, depth)
	s.WriteString(" */\n")
}
