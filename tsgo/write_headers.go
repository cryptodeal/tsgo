package tsgo

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
)

func (g *PackageGenerator) writeFileCodegenHeader(w *strings.Builder) {
	w.WriteString("// Code generated by tsgo. DO NOT EDIT.\n")
}

func (g *PackageGenerator) writeESLintIgnore(w *strings.Builder) {
	w.WriteString("/* eslint-disable */\n")
}

func (g *PackageGenerator) writeFileFrontmatter(w *strings.Builder) {
	if g.conf.Frontmatter != "" {
		w.WriteString(g.conf.Frontmatter)
	}
}

func (g *PackageGenerator) writeFileSourceHeader(w *strings.Builder, path string, file *ast.File) {
	w.WriteString("\n//////////\n// source: ")
	w.WriteString(fmt.Sprintf("%s\n", filepath.Base(path)))

	if file.Doc != nil {
		w.WriteString(fmt.Sprintf("/*\n%s*/\n", file.Doc.Text()))
	}
	w.WriteByte('\n')
}

func (g *PackageGenerator) writeFFIHeaders(w *strings.Builder) {
	w.WriteString("import { dlopen, FFIType, toArrayBuffer } from 'bun:ffi';\n")
}
