package tsgo

import (
	"go/ast"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UsedParams []string

func (g *PackageGenerator) writeCGoHeaders(cg *strings.Builder, gi *strings.Builder, ec *strings.Builder) {
	cg.WriteString("// ")
	g.writeFileCodegenHeader(cg)
	cg.WriteString("package main\n\n")
	cg.WriteString("/*\n")
	cg.WriteString("#include <stdlib.h>\n")
	// not needed afaik
	// cg.WriteString("#include <string.h>\n")
	cg.WriteString(ec.String())
	cg.WriteString("*/\n")
	cg.WriteString("import \"C\"\n\n")
	cg.WriteString("import (\n")
	cg.WriteString(gi.String())
	cg.WriteString(")\n\n")
}

func (g *PackageGenerator) writeCArrayHandler(cg *strings.Builder, ec *strings.Builder, t string, fmtr cases.Caser) string {
	sizeHandler := g.addCSizeHelper(ec, t)
	var arrTypeSB strings.Builder
	arrTypeSB.WriteByte('C')
	arrTypeSB.WriteString(fmtr.String(t))
	cg.WriteString("func ")
	cg.WriteString(arrTypeSB.String())
	cg.WriteString("(b []")
	cg.WriteString(t)
	cg.WriteString(") unsafe.Pointer {\n")
	g.writeIndent(cg, 1)
	cg.WriteString("p := C.malloc(C.size_t(len(b)) * C.")
	cg.WriteString(sizeHandler)
	cg.WriteString("())\n")
	g.writeIndent(cg, 1)
	cg.WriteString("sliceHeader := struct {\n")
	g.writeIndent(cg, 2)
	cg.WriteString("p   unsafe.Pointer\n")
	g.writeIndent(cg, 2)
	cg.WriteString("len int\n")
	g.writeIndent(cg, 2)
	cg.WriteString("cap int\n")
	g.writeIndent(cg, 1)
	cg.WriteString("}{p, len(b), len(b)}\n")
	g.writeIndent(cg, 1)
	cg.WriteString("s := *(*[]")
	cg.WriteString(t)
	cg.WriteString(")(unsafe.Pointer(&sliceHeader))\n")
	g.writeIndent(cg, 1)
	cg.WriteString("copy(s, b)\n")
	g.writeIndent(cg, 1)
	cg.WriteString("return p\n")
	cg.WriteString("}\n\n")
	return arrTypeSB.String()
}

func (g *PackageGenerator) addGoImport(s *strings.Builder, pkg string) {
	if _, ok := g.ffi.GoImports[pkg]; ok {
		return
	}
	g.writeIndent(s, 1)
	s.WriteByte('"')
	s.WriteString(pkg)
	s.WriteByte('"')
	s.WriteByte('\n')
	g.ffi.GoImports[pkg] = true
}

func (g *PackageGenerator) addCSizeHelper(s *strings.Builder, numType string) string {
	var fnNameSB strings.Builder
	fnNameSB.WriteString(numType)
	fnNameSB.WriteString("Size")
	if !g.ffi.FFIHelpers[fnNameSB.String()] {
		s.WriteString("static inline size_t ")
		s.WriteString(fnNameSB.String())
		s.WriteString("() {\n")
		g.writeIndent(s, 1)
		s.WriteString("return sizeof(")
		switch numType {
		case "int8":
			s.WriteString("int8_t")
		case "int16":
			s.WriteString("int16_t")
		case "int32":
			s.WriteString("int32_t")
		case "int64":
			s.WriteString("int64_t")
		case "uint8":
			s.WriteString("uint8_t")
		case "uint16":
			s.WriteString("uint16_t")
		case "uint32":
			s.WriteString("uint32_t")
		case "uint64":
			s.WriteString("uint64_t")
		case "float32":
			s.WriteString("float")
		case "float64":
			s.WriteString("double")
		}
		s.WriteString(");\n")
		s.WriteString("}\n\n")
	}
	return fnNameSB.String()
}

func (g *PackageGenerator) addDisposePtr(s *strings.Builder) {
	if !g.ffi.FFIHelpers["disposePtr"] {
		s.WriteString("//export disposePtr\n")
		s.WriteString("func disposePtr(ptr unsafe.Pointer, ctx unsafe.Pointer) {\n")
		g.writeIndent(s, 1)
		s.WriteString("ptr_num := uintptr(ptr)\n")
		g.writeIndent(s, 1)
		s.WriteString("if _, ok := ptrTrckr[ptr_num]; ok {\n")
		g.writeIndent(s, 2)
		s.WriteString("delete(ptrTrckr, ptr_num)\n")
		g.writeIndent(s, 2)
		s.WriteString("defer C.free(ptr)\n")
		g.writeIndent(s, 1)
		s.WriteString("}\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers["disposePtr"] = true
	}
}

func (g *PackageGenerator) addJSONEncoder(s *strings.Builder, gi *strings.Builder) {
	if !g.ffi.FFIHelpers["encodeJSON"] {
		g.addGoImport(gi, "encoding/json")
		s.WriteString("func encodeJSON(x interface{}{}) {\n")
		g.writeIndent(s, 1)
		s.WriteString("res, err := json.Marshal(x)\n")
		g.writeIndent(s, 1)
		s.WriteString("if err != nil {\n")
		g.writeIndent(s, 2)
		s.WriteString("fmt.Println(err)\n")
		g.writeIndent(s, 2)
		s.WriteString("panic(err)\n")
		g.writeIndent(s, 1)
		s.WriteString("}\n")
		g.writeIndent(s, 1)
		s.WriteString("return res\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers["encodeJSON"] = true
	}
}

func (g *PackageGenerator) addArraySize(s *strings.Builder) {
	if !g.ffi.FFIHelpers["ArraySize"] {
		s.WriteString("//export ArraySize\n")
		s.WriteString("func ArraySize(array unsafe.Pointer) C.size_t {\n")
		g.writeIndent(s, 1)
		s.WriteString("return ptrTrckr[uintptr(array)]\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers["ArraySize"] = true
	}
}

func (g *PackageGenerator) addPtrTrckr(s *strings.Builder) {
	if !g.ffi.FFIHelpers["ptrTrckr"] {
		g.ffi.FFIHelpers["ptrTrckr"] = true
		s.WriteString("var ptrTrckr = make(map[uintptr]C.size_t)\n\n")
	}
}

// TODO: parse to generate CGo code and/or Bun FFI Wrapper for specified functions
func (g *PackageGenerator) writeCGo(cg *strings.Builder, fd []*ast.FuncDecl, pkgName string) {
	var goImportsSB strings.Builder
	var embeddedCSB strings.Builder
	var goHelpersSB strings.Builder

	caser := cases.Title(language.AmericanEnglish)
	g.addGoImport(&goImportsSB, g.conf.Path)

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
			if type_str == "unsafe.Pointer" {
				g.addGoImport(&goImportsSB, "unsafe")
			}
			fn_str.WriteString(type_str)
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
		used_vars := UsedParams{}
		for _, param := range f.Type.Params.List {
			g.writeIndent(&fn_str, 1)
			var tempSB strings.Builder
			g.writeCGoType(&tempSB, param.Type, 0, true)
			type_str := tempSB.String()
			switch type_str {
			case "*C.char":
				parsedSB := strings.Builder{}
				parsedSB.WriteByte('_')
				parsedSB.WriteString(param.Names[0].Name)
				fn_str.WriteString(parsedSB.String())
				fn_str.WriteString(" := C.GoString(")
				fn_str.WriteString(param.Names[0].Name)
				fn_str.WriteString(")\n")
				used_vars = append(used_vars, parsedSB.String())
			default:
				used_vars = append(used_vars, param.Names[0].Name)
			}
		}
		g.writeIndent(&fn_str, 1)
		fn_str.WriteString("_returned_value := ")
		var tempResType strings.Builder
		g.writeCGoResType(&tempResType, &goImportsSB, &goHelpersSB, &embeddedCSB, caser, f.Type.Results.List[0].Type, 0, true)
		fn_str.WriteString(tempResType.String())

		fn_str.WriteByte('(')
		fn_str.WriteString(pkgName)
		fn_str.WriteByte('.')
		fn_str.WriteString(f.Name.Name)
		fn_str.WriteString("(")
		for i, param := range used_vars {
			fn_str.WriteString(param)
			if i < len(used_vars)-1 {
				fn_str.WriteString(", ")
			}
		}
		fn_str.WriteString("))\n")
		g.writeIndent(&fn_str, 1)
		fn_str.WriteString("return _returned_value\n")
		fn_str.WriteString("}\n\n")
	}

	g.writeCGoHeaders(cg, &goImportsSB, &embeddedCSB)

	cg.WriteString(goHelpersSB.String())

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
