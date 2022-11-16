package tsgo

import (
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type UsedParams []string

// TODO: parse to generate CGo code and/or Bun FFI Wrapper for specified functions
func (g *PackageGenerator) writeCGo(cg *strings.Builder, fd []*ast.FuncDecl, pkgName string) {
	cg.WriteString("package main\n\n")
	cg.WriteString("/*\n")
	cg.WriteString("#include <stdlib.h>\n")
	cg.WriteString("#include <string.h>\n")
	cg.WriteString("*/\n")
	cg.WriteString("import \"C\"\n\n")
	cg.WriteString("import (\n")

	has_str_param := false

	var fn_str strings.Builder
	for _, f := range fd {
		fn_str.WriteString("//export _")
		fn_str.WriteString(f.Name.Name)
		fn_str.WriteString("\n func _")
		fn_str.WriteString(f.Name.Name)
		fn_str.WriteString(" (")
		for i, param := range f.Type.Params.List {
			fn_str.WriteString(param.Names[0].Name)
			fn_str.WriteByte(' ')
			var tempSB strings.Builder
			g.writeCGoType(&tempSB, param.Type, 0, true)
			type_str := tempSB.String()
			fn_str.WriteString(type_str)
			if !has_str_param && type_str == "*C.char" {
				has_str_param = true
				g.writeIndent(cg, 1)
				cg.WriteString("\"unsafe\"\n")
			}
			if i < len(f.Type.Params.List)-1 {
				fn_str.WriteString(", ")
			}
		}
		var resSB strings.Builder
		g.writeCGoType(&resSB, f.Type.Results.List[0].Type, 0, true)
		res_type := resSB.String()
		fn_str.WriteString(") ")
		fn_str.WriteString(res_type)
		fn_str.WriteString(" {\n")
		g.writeIndent(&fn_str, 1)
		used_vars := UsedParams{}
		for _, param := range f.Type.Params.List {
			var tempSB strings.Builder
			g.writeCGoType(&tempSB, param.Type, 0, true)
			type_str := tempSB.String()
			switch type_str {
			case "*C.char":
				parsedSB := strings.Builder{}
				parsedSB.WriteByte('_')
				parsedSB.WriteString(param.Names[0].Name)
				fn_str.WriteString(parsedSB.String())
				fn_str.WriteString(" = C.GoString(")
				fn_str.WriteString(param.Names[0].Name)
				fn_str.WriteString(")\n")
				fn_str.WriteString("defer C.free(unsafe.Pointer(")
				fn_str.WriteString(parsedSB.String())
				fn_str.WriteString("))\n")
				used_vars = append(used_vars, parsedSB.String())
			default:
				used_vars = append(used_vars, param.Names[0].Name)
			}

		}
		fn_str.WriteString("return ")
		fn_str.WriteString(f.Name.Name)
		fn_str.WriteString("(")
		for i, param := range used_vars {
			fn_str.WriteString(param)
			if i < len(used_vars)-1 {
				fn_str.WriteString(", ")
			}
		}
		fn_str.WriteString(");\n")
		fn_str.WriteString("}\n\n")
	}

	cg.WriteString(")\n\n")
	cg.WriteString(fn_str.String())
	cg.WriteString("func main() {} // Required but ignored")

	var outPath strings.Builder
	outPath.WriteString(filepath.Dir(g.pkg.GoFiles[0]))
	outPath.WriteByte('/')
	outPath.WriteString(pkgName)
	outPath.WriteString("/gen_bindings.go")

	err := os.MkdirAll(filepath.Dir(outPath.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
	err = ioutil.WriteFile(outPath.String(), []byte(cg.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
}
