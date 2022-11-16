package tsgo

import (
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: parse to generate CGo code and/or Bun FFI Wrapper for specified functions
func (g *PackageGenerator) writeCGo(cg *strings.Builder, fd []*ast.FuncDecl, pkgName string) {
	cg.WriteString("package ")
	cg.WriteString(pkgName)
	cg.WriteString("_Gen_TSGo\n\n")
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
				cg.WriteString("\"encoding/json\"\n")
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
		fn_str.WriteString("return _")
		fn_str.WriteString(f.Name.Name)
		fn_str.WriteString("(")
		for i, param := range f.Type.Params.List {
			fn_str.WriteString(param.Names[0].Name)
			if i < len(f.Type.Params.List)-1 {
				fn_str.WriteString(", ")
			}
		}
		fn_str.WriteString(");\n")
		fn_str.WriteString("}\n\n")
	}

	cg.WriteString(")\n")
	cg.WriteString(fn_str.String())

	var outPath strings.Builder
	outPath.WriteString(filepath.Dir(g.conf.OutputPath))
	outPath.WriteByte('/')
	outPath.WriteString(pkgName)
	outPath.WriteString("/_ffi_bindings.go")
	err := ioutil.WriteFile(outPath.String(), []byte(cg.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
}
