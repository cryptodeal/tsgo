package tsgo

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UsedParams []string

func (g *PackageGenerator) writeCGoHeaders(cg *strings.Builder, gi *strings.Builder, ec *strings.Builder, ci *strings.Builder) {
	g.writeFileCodegenHeader(cg)
	cg.WriteString("package main\n\n/*\n#include <stdlib.h>\n")
	cg.WriteString(ci.String())
	cg.WriteString(ec.String())
	cg.WriteString(fmt.Sprintf("*/\nimport %q\n\nimport(\n%s)\n\n", "C", gi.String()))
}

// TODO:
// * see if we can handle `complex64` and `complex128`?
// * perhaps do a better job of mapping (no default value??)
func (g *PackageGenerator) getNumCast(s string) string {
	// fmt.Println(s)
	switch s {
	case "C.int":
		return "int"
	case "C.int8_t":
		return "int8"
	case "C.int16_t":
		return "int16"
	case "C.int32_t":
		return "int32"
	case "C.int64_t":
		return "int64"
	case "C.uint":
		return "uint"
	case "C.uint8_t":
		return "uint8"
	case "C.uint16_t":
		return "uint16"
	case "C.uint32_t":
		return "uint32"
	case "C.uint64_t":
		return "uint64"
	case "C.float":
		return "float32"
	case "C.double":
		return "float64"
	case "C.uintptr_t":
		return "uintptr"
	}
	return ""
}

func (g *PackageGenerator) writeCArrayHandler(cg *strings.Builder, ec *strings.Builder, t string, fmtr cases.Caser) string {
	sizeHandler := g.addCSizeHelper(ec, t)
	arrType := fmt.Sprintf("C%s", fmtr.String(t))
	if !g.ffi.FFIHelpers[arrType] {
		cg.WriteString(fmt.Sprintf("func %s(b []%s) unsafe.Pointer {\n", arrType, t))
		g.writeIndent(cg, 1)
		cg.WriteString("arr_len := len(b)\n")
		g.writeIndent(cg, 1)
		cg.WriteString(fmt.Sprintf("p := C.malloc(C.size_t(arr_len) * C.%s())\n", sizeHandler))
		g.writeIndent(cg, 1)
		cg.WriteString("sliceHeader := struct {\n")
		g.writeIndent(cg, 2)
		cg.WriteString("p   unsafe.Pointer\n")
		g.writeIndent(cg, 2)
		cg.WriteString("len int\n")
		g.writeIndent(cg, 2)
		cg.WriteString("cap int\n")
		g.writeIndent(cg, 1)
		cg.WriteString("}{p, arr_len, arr_len}\n")
		g.writeIndent(cg, 1)
		cg.WriteString(fmt.Sprintf("s := *(*[]%s)(unsafe.Pointer(&sliceHeader))\n", t))
		g.writeIndent(cg, 1)
		cg.WriteString("copy(s, b)\n")
		g.writeIndent(cg, 1)
		cg.WriteString("ptrTrckr[p] = C.size_t(arr_len)\n")
		g.writeIndent(cg, 1)
		cg.WriteString("return p\n}\n\n")
		g.ffi.FFIHelpers[arrType] = true
	}
	return arrType
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

func (g *PackageGenerator) addCImport(s *strings.Builder, pkg string, isLocal bool) {
	if _, ok := g.ffi.CImports[pkg]; ok {
		return
	}
	s.WriteString("#include ")
	if isLocal {
		s.WriteByte('"')
	} else {
		s.WriteByte('<')
	}
	s.WriteString(pkg)
	if isLocal {
		s.WriteString("\"\n")
	} else {
		s.WriteString(">\n")
	}
	g.ffi.CImports[pkg] = true
}

func (g *PackageGenerator) addCSizeHelper(s *strings.Builder, numType string) string {
	var fnNameSB strings.Builder
	fnNameSB.WriteString(numType)
	fnNameSB.WriteString("Size")
	if !g.ffi.CHelpers[fnNameSB.String()] {
		s.WriteString(fmt.Sprintf("\nstatic inline size_t %s() {\n", fnNameSB.String()))
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
		s.WriteString("}\n")
		g.ffi.CHelpers[fnNameSB.String()] = true
	}
	return fnNameSB.String()
}

func (g *PackageGenerator) addDisposePtr(s *strings.Builder, gi *strings.Builder) {
	if !g.ffi.FFIHelpers["dispose"] {
		s.WriteString("//export dispose\n")
		s.WriteString("func dispose(ptr unsafe.Pointer, ctx unsafe.Pointer) {\n")
		g.writeIndent(s, 1)
		s.WriteString("if _, ok := ptrTrckr[ptr]; ok {\n")
		g.writeIndent(s, 2)
		s.WriteString("delete(ptrTrckr, ptr)\n")
		g.writeIndent(s, 2)
		s.WriteString("defer C.free(ptr)\n")
		g.writeIndent(s, 1)
		s.WriteString("} else {\n")
		g.writeIndent(s, 2)
		g.addGoImport(gi, "fmt")
		s.WriteString("panic(fmt.Sprintf(\"Error: `%#v` not found in ptrTrckr\", ptr))\n")
		g.writeIndent(s, 1)
		s.WriteString("}\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers["dispose"] = true
	}

	if !g.ffi.FFIHelpers["genDisposePtr"] {
		name := "genDisposePtr"
		var res_helper = &ResHelpers{
			FFIType:     "FFIType.ptr",
			CGoWrapType: "unsafe.Pointer",
		}
		var ffi_func = &FFIFunc{
			args:           []*ArgHelpers{},
			returns:        []*ResHelpers{res_helper},
			isHandleFn:     false,
			name:           &name,
			fieldAccessors: []*StructAccessor{},
		}
		s.WriteString("//export genDisposePtr\n")
		s.WriteString("func genDisposePtr() unsafe.Pointer {\n")
		g.writeIndent(s, 1)
		s.WriteString("return C.disposePtr\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers["genDisposePtr"] = true
		g.ffi.FFIFuncs[name] = ffi_func
	}
}

func (g *PackageGenerator) addCDisposeHelpers(ci *strings.Builder, pkgName string) {
	if !g.ffi.CHelpers["helpers.h"] && !g.ffi.CHelpers["helpers.c"] && !g.ffi.CImports["helpers.h"] {
		ci.WriteString("#include \"helpers.h\"\n")
		g.ffi.CImports["helpers.h"] = true

		var cHelpersHeaders strings.Builder
		var cHelpers strings.Builder

		g.writeFileCodegenHeader(&cHelpersHeaders)
		cHelpersHeaders.WriteString("#include <stdint.h>\n")
		cHelpersHeaders.WriteString("void disposePtr(void *, void *);\n")
		cHelpersHeaders.WriteString("void *hackyHandle(uintptr_t);\n")

		filePathDir := filepath.Dir(g.pkg.GoFiles[0])
		headersPath := fmt.Sprintf("%s/%s/helpers.h", filePathDir, pkgName)
		err := os.MkdirAll(filepath.Dir(headersPath), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}
		err = os.WriteFile(headersPath, []byte(cHelpersHeaders.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}

		g.writeFileCodegenHeader(&cHelpers)
		cHelpers.WriteString(fmt.Sprintf("#include %q\n\n", "_cgo_export.h"))
		cHelpers.WriteString("void disposePtr(void *ptr, void *ctx)\n")
		cHelpers.WriteString("{\n")
		g.writeIndent(&cHelpers, 1)
		cHelpers.WriteString("dispose(ptr, ctx);\n")
		cHelpers.WriteString("}\n")

		cHelpers.WriteString("void *hackyHandle(uintptr_t ptr)\n")
		cHelpers.WriteString("{\n")
		g.writeIndent(&cHelpers, 1)
		cHelpers.WriteString("return (void *)ptr;\n")
		cHelpers.WriteString("}\n")

		helpersPath := fmt.Sprintf("%s/%s/helpers.c", filePathDir, pkgName)

		err = os.MkdirAll(filepath.Dir(helpersPath), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}
		err = os.WriteFile(helpersPath, []byte(cHelpers.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}
	}
}

func (g *PackageGenerator) addJSONEncoder(s *strings.Builder, gi *strings.Builder) {
	if !g.ffi.FFIHelpers["encodeJSON"] {
		g.addGoImport(gi, "encoding/json")
		g.addGoImport(gi, "fmt")
		s.WriteString("func encodeJSON(x interface{}) []byte {\n")
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

func (g *PackageGenerator) addArraySize(s *strings.Builder, gi *strings.Builder) {
	name := "arraySize"
	if !g.ffi.FFIHelpers[name] {
		var arg_helper = &ArgHelpers{
			FFIType:     "FFIType.ptr",
			CGoWrapType: "unsafe.Pointer",
			OGGoType:    "unsafe.Pointer",
			Name:        name,
		}
		var res_helper = &ResHelpers{
			FFIType:     "FFIType.u64_fast",
			CGoWrapType: "C.size_t",
			OGGoType:    "C.uint64_t",
		}
		var ffi_func = &FFIFunc{
			args:       []*ArgHelpers{arg_helper},
			returns:    []*ResHelpers{res_helper},
			isHandleFn: false,
			name:       &name,
		}
		s.WriteString("//export arraySize\n")
		s.WriteString("func arraySize(ptr unsafe.Pointer) C.size_t {\n")
		g.writeIndent(s, 1)
		s.WriteString("if val, ok := ptrTrckr[ptr]; ok {\n")
		g.writeIndent(s, 2)
		s.WriteString("return val\n")
		g.writeIndent(s, 1)
		s.WriteString("}\n")
		g.writeIndent(s, 1)
		g.addGoImport(gi, "fmt")
		s.WriteString("panic(fmt.Sprintf(\"Error: `%#v` not found in ptrTrckr\", ptr))\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers[name] = true
		g.ffi.FFIFuncs[name] = ffi_func
	}
}

func (g *PackageGenerator) addPtrTrckr(s *strings.Builder) {
	if !g.ffi.FFIHelpers["ptrTrckr"] {
		s.WriteString("var ptrTrckr = make(map[unsafe.Pointer]C.size_t)\n\n")
		g.ffi.FFIHelpers["ptrTrckr"] = true
	}
}

func (g *PackageGenerator) addArgHandler(s *strings.Builder, gi *strings.Builder, f *ArgHelpers, usedVars *UsedParams) {
	g.writeIndent(s, 1)
	var tempSB strings.Builder
	g.writeCGoType(&tempSB, f.ASTField.Type, 0, true)
	type_str := tempSB.String()
	switch type_str {
	case "*C.char":
		parsedSB := fmt.Sprintf("_%s", f.Name)
		s.WriteString(fmt.Sprintf("%s := C.GoString(%s)\n", parsedSB, f.Name))
		*usedVars = append(*usedVars, parsedSB)
	case "unsafe.Pointer":
		parsedSB := fmt.Sprintf("_%s", f.Name)
		g.addGoImport(gi, "unsafe")
		arr_dat_type := g.getArrayType(f.ASTField.Type)
		s.WriteString(fmt.Sprintf("%s := unsafe.Slice((*%s)(%s), %s_len)\n", parsedSB, arr_dat_type, f.Name, f.Name))
		*usedVars = append(*usedVars, parsedSB)
	default:
		*usedVars = append(*usedVars, f.Name)
	}
}

func isStruct(t string) bool {
	if strings.Contains(t, "C.") || strings.Contains(t, "map[") {
		return false
	}
	switch t {
	case "bool":
		return false
	case "string":
		return false
	case "int", "int8", "int16", "int32", "int64":
		return false
	case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return false
	case "byte":
		return false
	case "rune":
		return false
	case "float32", "float64":
		return false
	case "complex64", "complex128":
		return false
	default:
		return true
	}
}

// bulk of parsing the function is done here
func (g *PackageGenerator) isResHandle(t ast.Expr) (bool, string) {
	isHandle := false
	structName := ""
	switch t := t.(type) {
	case *ast.StarExpr:
		return g.isResHandle(t.X)

	case *ast.Ident:
		struct_name := g.getStructName(t)
		if isStruct(struct_name) {
			isHandle = true
			structName = struct_name
		}
	}
	return isHandle, structName
}

// bulk of parsing the function is done here
func (g *PackageGenerator) isTypedArray(t ast.Expr) (bool, string) {
	isArray := false
	var dType string
	switch t := t.(type) {
	case *ast.ArrayType:
		isArray = true
		dType = g.getArrayType(t)
	}
	return isArray, dType
}

func (g *PackageGenerator) parseAccessors(fields *[]*StructAccessor, name string) {
	if _, ok := g.ffi.StructHelpers[name]; ok && len(*fields) > 0 && !g.ffi.ParsedStructs[name] {
		g.ffi.ParsedStructs[name] = true
		for _, fa := range *fields {
			if fa.isHandleFn != nil && !g.ffi.ParsedStructs[*fa.isHandleFn] {
				fa.fieldAccessors = g.ffi.StructHelpers[*fa.isHandleFn]
				g.parseAccessors(&fa.fieldAccessors, *fa.isHandleFn)
			}
		}
	}
}

func (g *PackageGenerator) parseFn(f *ast.FuncDecl) *FFIFunc {
	isStarExpr := false
	switch f.Type.Results.List[0].Type.(type) {
	case *ast.StarExpr:
		isStarExpr = true
	}

	var ffi_func = &FFIFunc{
		args:       []*ArgHelpers{},
		returns:    []*ResHelpers{},
		isHandleFn: false,
		isStarExpr: isStarExpr,
	}

	for _, param := range f.Type.Params.List {
		var tempSB strings.Builder
		g.writeCGoType(&tempSB, param.Type, 0, true)
		var arg_helper = &ArgHelpers{
			FFIType:     getFFIIdent(tempSB.String()),
			CGoWrapType: tempSB.String(),
			OGGoType:    tempSB.String(),
			Name:        param.Names[0].Name,
			ASTField:    param,
		}
		isArray, dType := g.isTypedArray(param.Type)
		if isArray {
			var tempTypeSB strings.Builder
			tempTypeSB.WriteString("[]")
			tempTypeSB.WriteString(dType)
			arg_helper.OGGoType = tempTypeSB.String()
		}
		ffi_func.args = append(ffi_func.args, arg_helper)
		if isArray {
			var lenName strings.Builder
			lenName.WriteString(param.Names[0].Name)
			lenName.WriteString("_len")
			var len_helper = &ArgHelpers{
				FFIType:     "FFIType.u64_fast",
				CGoWrapType: "C.uint64_t",
				OGGoType:    "C.size_t",
				Name:        lenName.String(),
			}
			ffi_func.args = append(ffi_func.args, len_helper)
		}
	}

	for i, res := range f.Type.Results.List {
		if i == 0 {
			isHandle, structName := g.isResHandle(res.Type)
			if isHandle {
				ffi_func.isHandleFn = true
				ffi_func.name = &structName
			}
		}
		var tempSB strings.Builder
		g.writeCGoType(&tempSB, res.Type, 0, true)

		var res_helper = &ResHelpers{
			CGoWrapType: tempSB.String(),
			OGGoType:    tempSB.String(),
			FFIType:     getFFIIdent(tempSB.String()),
			ASTType:     &res.Type,
		}
		ffi_func.returns = append(ffi_func.returns, res_helper)
	}

	if ffi_func.isHandleFn {
		ffi_func.fieldAccessors = g.ffi.StructHelpers[*ffi_func.name]
		g.parseAccessors(&ffi_func.fieldAccessors, *ffi_func.name)
		var ptr_arg = &ArgHelpers{
			Name:        "handle",
			FFIType:     "FFIType.ptr",
			CGoWrapType: "C.uintptr_t",
			OGGoType:    "unsafe.Pointer",
		}
		disposeFnName := fmt.Sprintf("_dispose_%s", *ffi_func.name)
		ffi_func.disposeHandle = &DisposeStructFunc{
			args:   []*ArgHelpers{ptr_arg},
			fnName: disposeFnName,
			name:   *ffi_func.name,
		}
	}
	return ffi_func
}

// TODO: think 1 fn can handle disposing all structs
func (g *PackageGenerator) writeDisposeStruct(t *DisposeStructFunc) string {
	var disposeSB strings.Builder
	if !g.ffi.FFIHelpers["_DISPOSE_Struct"] {
		disposeSB.WriteString("//export _DISPOSE_Struct\n")
		disposeSB.WriteString(fmt.Sprintf("func _DISPOSE_Struct(%s %s) {\n", t.args[0].Name, t.args[0].CGoWrapType))
		g.writeIndent(&disposeSB, 1)
		disposeSB.WriteString(fmt.Sprintf("h := cgo.Handle(%s)\n", t.args[0].Name))
		g.writeIndent(&disposeSB, 1)
		disposeSB.WriteString(fmt.Sprintf("fmt.Println(\"deleted handle @ uintptr:\", %s)\n", t.args[0].Name))
		g.writeIndent(&disposeSB, 1)
		disposeSB.WriteString("h.Delete()\n")
		disposeSB.WriteString("}\n\n")
		g.ffi.FFIHelpers["_DISPOSE_Struct"] = true
	}
	return disposeSB.String()
}

func (g *PackageGenerator) writeCGoFieldAccessor(gi *strings.Builder, gh *strings.Builder, ec *strings.Builder, ci *strings.Builder, fmtr cases.Caser, f *StructAccessor, pkgName string, structName string) string {
	used_args := UsedParams{}
	var fnSB strings.Builder
	fnSB.WriteString(fmt.Sprintf("//export %s\n", *f.fnName))
	fnSB.WriteString(fmt.Sprintf("func %s(", *f.fnName))
	// iterate through fn params, generating cgo function decl line
	argLen := len(f.args)
	if argLen > 0 {
		for i, arg := range f.args {
			fnSB.WriteString(fmt.Sprintf("%s %s", arg.Name, arg.CGoWrapType))
			if i < argLen-1 {
				fnSB.WriteString(", ")
			}
		}
	}
	fnSB.WriteString(") ")
	// write return type (if any)
	if len(f.returns) > 0 {
		if f.isHandleFn != nil {
			g.addGoImport(gi, "unsafe")
			g.addGoImport(gi, "runtime/cgo")
			g.addCDisposeHelpers(ci, pkgName)
			g.addCImport(ci, "stdint.h", false)
			fnSB.WriteString("unsafe.Pointer")
		} else {
			fnSB.WriteString(f.returns[0].CGoWrapType)
		}
		fnSB.WriteByte(' ')
	}
	fnSB.WriteString("{\n")
	// gen necessary type coercions (CGo C types -> Go types)
	if argLen > 0 {
		for i, arg := range f.args {
			if i == 0 {
				g.writeIndent(&fnSB, 1)
				fnSB.WriteString("h := cgo.Handle(handle)\n")
				g.writeIndent(&fnSB, 1)
				fnSB.WriteString(fmt.Sprintf("s := h.Value().(%s.%s)\n", pkgName, structName))
			} else {
				// if `arg.ASTField == nil`, it's a helper arg like `len`, which isn't passed to wrapped Go func
				if arg.ASTField != nil {
					g.addArgHandler(&fnSB, gi, arg, &used_args)
				}
			}
		}
	}

	// return `nil` if no value @ field
	if f.isOptional || *f.arrayType != "" {
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString(fmt.Sprintf("if s.%s == nil {\n", *f.name))
		g.writeIndent(&fnSB, 2)
		fnSB.WriteString("return nil\n")
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("}\n")
	}

	// write returned value (or intermediary, if necessary)
	tempResType := g.getCgoHandler(f.returns[0].CGoWrapType)
	g.writeIndent(&fnSB, 1)
	if f.isHandleFn != nil {
		fnSB.WriteString("return C.hackyHandle(C.uintptr_t(cgo.NewHandle")
	} else {
		fnSB.WriteString(fmt.Sprintf("_returned_value := %s(", tempResType))
		if *f.arrayType != "" {
			fnSB.WriteString(g.writeCArrayHandler(gh, ec, *f.arrayType, fmtr))
		} else {
			fnSB.WriteString(g.getGoType(f.returns[0].CGoWrapType))
		}
	}

	fnSB.WriteByte('(')
	if f.isStarExpr {
		fnSB.WriteByte('*')
	}
	fnSB.WriteString(fmt.Sprintf("s.%s", *f.name))
	if f.isHandleFn != nil {
		fnSB.WriteString(")))\n")
	} else {
		fnSB.WriteString("))\n")
		// TODO: need to improve API so this code is simplified/handles more edge cases
		if tempResType == "C.CString" {
			g.writeIndent(&fnSB, 1)
			fnSB.WriteString("defer C.free(unsafe.Pointer(_returned_value))\n")
		}
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("return _returned_value\n")
	}

	fnSB.WriteString("}\n\n")

	if f.isHandleFn != nil {
		for _, fa := range f.fieldAccessors {
			name := fmt.Sprintf("_GET_%s_%s", *f.isHandleFn, *fa.name)
			fa.fnName = &name
			fnSB.WriteString(g.writeCGoFieldAccessor(gi, gh, ec, ci, fmtr, fa, pkgName, *f.isHandleFn))
		}
		fnSB.WriteString(g.writeDisposeStruct(f.disposeHandle))
		g.ffi.GoWrappedStructs[*f.isHandleFn] = true
		// write wrapper to create new struct
		g.writeInitStructMethod(&fnSB, *f.isHandleFn, pkgName, f.fieldAccessors)
	}

	return fnSB.String()
}

func (g *PackageGenerator) writeInitStructMethod(s *strings.Builder, name string, pkgName string, fieldAccessors []*StructAccessor) {
	fnName := fmt.Sprintf("_INIT_%s", name)
	if !g.ffi.FFIHelpers[fnName] {
		// write wrapper to create new struct
		const alphaArgs = "abcdefghijklmnopqrstuvwxyz"
		s.WriteString(fmt.Sprintf("//export _INIT_%s\n", name))
		s.WriteString(fmt.Sprintf("func _INIT_%s(", name))
		argLen := len(fieldAccessors)
		for i, arg := range fieldAccessors {
			if arg.returns[0].CGoWrapType == "unsafe.Pointer" && arg.isHandleFn != nil {
				s.WriteString(fmt.Sprintf("%s %s", string(alphaArgs[i]), "C.uintptr_t"))
			} else {
				s.WriteString(fmt.Sprintf("%s %s", string(alphaArgs[i]), arg.returns[0].CGoWrapType))
				if arg.returns[0].CGoWrapType == "unsafe.Pointer" && arg.arrayType != nil {
					s.WriteString(fmt.Sprintf(", %s_len C.uint64_t", string(alphaArgs[i])))
				}
			}
			if i < argLen-1 {
				s.WriteString(", ")
			}
		}
		s.WriteString(") unsafe.Pointer {\n")
		//TODO: parse args (casting types as need be) and return Handle for new struct
		var usedArgs = []string{}
		for i, arg := range fieldAccessors {
			if arg.arrayType != nil && g.isTypedArrayHelper(*arg.arrayType) {
				usedName := fmt.Sprintf("_%s", string(alphaArgs[i]))
				usedArgs = append(usedArgs, usedName)
				g.writeIndent(s, 1)
				s.WriteString(fmt.Sprintf("%s := unsafe.Slice((*%s)(%s), %s_len)\n", usedName, *arg.arrayType, string(alphaArgs[i]), string(alphaArgs[i])))
			} else if arg.isHandleFn != nil {
				usedName := fmt.Sprintf("_%s", string(alphaArgs[i]))
				usedArgs = append(usedArgs, usedName)
				g.writeIndent(s, 1)
				s.WriteString(fmt.Sprintf("%s_h := cgo.Handle(%s)\n", string(alphaArgs[i]), string(alphaArgs[i])))
				g.writeIndent(s, 1)
				s.WriteString(fmt.Sprintf("%s := %s_h.Value().(%s.%s)\n", usedName, string(alphaArgs[i]), pkgName, *arg.isHandleFn))
			} else if arg.returns[0].CGoWrapType == "*C.char" {
				isStructType, structType := g.isResHandle(*arg.returns[0].ASTType)
				usedName := fmt.Sprintf("_%s", string(alphaArgs[i]))
				usedArgs = append(usedArgs, usedName)
				g.writeIndent(s, 1)
				if !isStructType {
					s.WriteString(fmt.Sprintf("%s := C.GoString(%s)\n", usedName, string(alphaArgs[i])))
				} else {
					s.WriteString(fmt.Sprintf("%s := %s.%s(C.GoString(%s))\n", usedName, pkgName, structType, string(alphaArgs[i])))

				}
			} else {
				usedArgs = append(usedArgs, string(alphaArgs[i]))
			}
		}
		g.writeIndent(s, 1)
		s.WriteString(fmt.Sprintf("res := %s.%s{", pkgName, name))
		for i, arg := range fieldAccessors {
			s.WriteString(fmt.Sprintf("%s: ", *arg.name))
			if arg.isStarExpr {
				s.WriteString(fmt.Sprintf("&%s", usedArgs[i]))
			} else {
				num_type := g.getNumCast(arg.returns[0].CGoWrapType)
				if num_type != "" {
					s.WriteString(fmt.Sprintf("%s(%s)", num_type, usedArgs[i]))
				} else {
					s.WriteString(usedArgs[i])
				}
			}
			if i < argLen-1 {
				s.WriteString(", ")
			}
		}
		s.WriteString("}\n")
		g.writeIndent(s, 1)
		s.WriteString("return C.hackyHandle(C.uintptr_t(cgo.NewHandle(res)))\n")
		s.WriteString("}\n\n")
		g.ffi.FFIHelpers[fnName] = true
	}
}

func (g *PackageGenerator) writeCGoFn(gi *strings.Builder, gh *strings.Builder, ec *strings.Builder, ci *strings.Builder, fmtr cases.Caser, f *FFIFunc, name string, pkgName string) string {
	used_args := UsedParams{}
	var fnSB strings.Builder

	fnSB.WriteString(fmt.Sprintf("//export _%s\n", name))
	fnSB.WriteString(fmt.Sprintf("func _%s(", name))
	// iterate through fn params, generating cgo function decl line
	argLen := len(f.args)
	if argLen > 0 {
		for i, arg := range f.args {
			fnSB.WriteString(fmt.Sprintf("%s %s", arg.Name, arg.CGoWrapType))
			if i < argLen-1 {
				fnSB.WriteString(", ")
			}
		}
	}
	fnSB.WriteString(") ")

	// write return type (if any)
	if len(f.returns) > 0 {
		if f.isHandleFn {
			g.addGoImport(gi, "unsafe")
			g.addGoImport(gi, "runtime/cgo")
			g.addCDisposeHelpers(ci, pkgName)
			g.addCImport(ci, "stdint.h", false)
			fnSB.WriteString("unsafe.Pointer")
		} else {
			fnSB.WriteString(f.returns[0].CGoWrapType)
		}
		fnSB.WriteByte(' ')
	}
	fnSB.WriteString("{\n")

	// gen necessary type coercions (CGo C types -> Go types)
	if argLen > 0 {
		for _, arg := range f.args {
			// if `arg.ASTField == nil`, it's a helper arg like `len`, which isn't passed to wrapped Go func
			if arg.ASTField != nil {
				g.addArgHandler(&fnSB, gi, arg, &used_args)
			}
		}
	}

	var tempResType strings.Builder
	g.writeCGoResType(&tempResType, gi, gh, ec, ci, fmtr, *f.returns[0].ASTType, 0, true, pkgName)
	// write returned value (or intermediary, if necessary)
	g.writeIndent(&fnSB, 1)
	if f.isHandleFn {
		fnSB.WriteString("return C.hackyHandle(C.uintptr_t(cgo.NewHandle")
	} else {
		if tempResType.String() == "encodeJSON" {
			fnSB.WriteString("_temp_res_val := ")
		} else {
			fnSB.WriteString("_returned_value := ")
		}
		fnSB.WriteString(tempResType.String())
	}
	fnSB.WriteByte('(')
	if f.isStarExpr {
		fnSB.WriteByte('*')
	}
	fnSB.WriteString(fmt.Sprintf("%s.%s(", pkgName, name))

	// iterate through params (and converted params), writing args passed to function call
	for i, param := range used_args {
		fnSB.WriteString(param)
		if i < len(used_args)-1 {
			fnSB.WriteString(", ")
		}
	}
	if f.isHandleFn {
		fnSB.WriteString("))")
	}
	fnSB.WriteString("))\n")

	// TODO: need to improve API so this code is simplified/handles more edge cases
	if tempResType.String() == "encodeJSON" {
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("_returned_value := C.CString(string(_temp_res_val))\n")
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("defer C.free(unsafe.Pointer(_returned_value))\n")
	} else if tempResType.String() == "C.CString" {
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("defer C.free(unsafe.Pointer(_returned_value))\n")
	}
	if !f.isHandleFn {
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("return _returned_value\n")
	}

	fnSB.WriteString("}\n\n")
	return fnSB.String()
}

// TODO: parse to generate CGo code and/or Bun FFI Wrapper for specified functions
func (g *PackageGenerator) writeCGo(cg *strings.Builder, fd []*ast.FuncDecl, pkgName string) {
	var goImportsSB strings.Builder
	var embeddedCSB strings.Builder
	var goHelpersSB strings.Builder
	var cImportsSB strings.Builder

	caser := cases.Title(language.AmericanEnglish)
	g.addGoImport(&goImportsSB, g.conf.Path)
	// `C` is always required import for CGo
	g.ffi.GoImports["C"] = true

	// writes all functions to single string builder
	var fn_str strings.Builder

	// iterate through all function declarations, parsing into `*FFIFunc` helper struct & writing CGo/C helpers
	for _, f := range fd {
		func_data := g.parseFn(f)
		/*
			fmt.Println("test_func_parser:", func_data)
			if func_data.name != nil {
				fmt.Println("name: ", *func_data.name)
			}
			if func_data.fieldAccessors != nil {
				for _, field := range func_data.fieldAccessors {
					fmt.Println("field: ", *field.name)

					for _, a := range field.args {
						fmt.Println("arg: ", a)
					}

					for _, r := range field.returns {
						fmt.Println("returns: ", r)
					}
				}
			}
		*/
		g.ffi.FFIFuncs[f.Name.Name] = func_data

		fn_str.WriteString(g.writeCGoFn(&goImportsSB, &goHelpersSB, &embeddedCSB, &cImportsSB, caser, func_data, f.Name.Name, pkgName))
		if func_data.isHandleFn && !g.ffi.GoWrappedStructs[*func_data.name] {
			for _, field := range func_data.fieldAccessors {
				name := fmt.Sprintf("_GET_%s_%s", *func_data.name, *field.name)
				field.fnName = &name
				fn_str.WriteString(g.writeCGoFieldAccessor(&goImportsSB, &goHelpersSB, &embeddedCSB, &cImportsSB, caser, field, pkgName, *func_data.name))
			}
			fn_str.WriteString(g.writeDisposeStruct(func_data.disposeHandle))
			g.ffi.GoWrappedStructs[*func_data.name] = true

			// write wrapper to create new struct
			g.writeInitStructMethod(&fn_str, *func_data.name, pkgName, func_data.fieldAccessors)
		}
	}

	// write headers, embedded C imports/logic, and Go imports
	g.writeCGoHeaders(cg, &goImportsSB, &embeddedCSB, &cImportsSB)
	// writes Go helper Fns (e.g. encodeJSON, ArraySize, CFloat32, etc.)
	cg.WriteString(goHelpersSB.String())

	// writes all of the wrapper functions for FFI (generated above)
	cg.WriteString(fn_str.String())
	// writes required `func main()` to appease compiler
	cg.WriteString("func main() {} // Required but ignored")

	// write generated CGo wrapper bindings to file
	outPath := fmt.Sprintf("%s/%s/gen_bindings.go", filepath.Dir(g.pkg.GoFiles[0]), pkgName)

	err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
	err = os.WriteFile(outPath, []byte(cg.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
}
