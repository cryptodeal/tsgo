package tsgo

import (
	"go/ast"
	"path/filepath"
	"strings"
)

func (g *PackageGenerator) writeFileCodegenHeader(w *strings.Builder) {
	w.WriteString("// Code generated by tsgo. DO NOT EDIT.\n")
}

func (g *PackageGenerator) writeFileFrontmatter(w *strings.Builder) {
	if g.conf.Frontmatter != "" {
		w.WriteString(g.conf.Frontmatter)
	}
}

func (g *PackageGenerator) writeFileSourceHeader(w *strings.Builder, path string, file *ast.File) {
	w.WriteString("\n//////////\n// source: ")
	w.WriteString(filepath.Base(path))
	w.WriteString("\n")

	if file.Doc != nil {
		w.WriteString("/*\n")
		w.WriteString(file.Doc.Text())
		w.WriteString("*/\n")
	}
	w.WriteString("\n")
}

func (g *PackageGenerator) writeFFIHeaders(w *strings.Builder) {
	w.WriteString("import { dlopen, FFIType, toArrayBuffer } from 'bun:ffi';\n")
}
