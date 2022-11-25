package tsgo

import (
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
	cg.WriteString("package main\n\n")
	cg.WriteString("/*\n")
	cg.WriteString("#include <stdlib.h>\n")
	cg.WriteString(ci.String())

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
	if !g.ffi.FFIHelpers[arrTypeSB.String()] {
		cg.WriteString("func ")
		cg.WriteString(arrTypeSB.String())
		cg.WriteString("(b []")
		cg.WriteString(t)
		cg.WriteString(") unsafe.Pointer {\n")
		g.writeIndent(cg, 1)
		cg.WriteString("arr_len := len(b)\n")
		g.writeIndent(cg, 1)
		cg.WriteString("p := C.malloc(C.size_t(arr_len) * C.")
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
		cg.WriteString("}{p, arr_len, arr_len}\n")
		g.writeIndent(cg, 1)
		cg.WriteString("s := *(*[]")
		cg.WriteString(t)
		cg.WriteString(")(unsafe.Pointer(&sliceHeader))\n")
		g.writeIndent(cg, 1)
		cg.WriteString("copy(s, b)\n")
		g.writeIndent(cg, 1)
		cg.WriteString("ptrTrckr[p] = C.size_t(arr_len)\n")
		g.writeIndent(cg, 1)
		cg.WriteString("return p\n")
		cg.WriteString("}\n\n")
		g.ffi.FFIHelpers[arrTypeSB.String()] = true

	}
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

func (g *PackageGenerator) addCImport(s *strings.Builder, pkg string, isLocal bool) {
	if _, ok := g.ffi.CImports[pkg]; ok {
		return
	}
	g.writeIndent(s, 1)
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
		s.WriteByte('\n')
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

		var headersPath strings.Builder
		headersPath.WriteString(filepath.Dir(g.pkg.GoFiles[0]))
		headersPath.WriteByte('/')
		headersPath.WriteString(pkgName)
		headersPath.WriteString("/helpers.h")
		err := os.MkdirAll(filepath.Dir(headersPath.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}
		err = os.WriteFile(headersPath.String(), []byte(cHelpersHeaders.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}

		g.writeFileCodegenHeader(&cHelpers)
		cHelpers.WriteString("#include \"_cgo_export.h\"\n\n")
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

		var helpersPath strings.Builder
		helpersPath.WriteString(filepath.Dir(g.pkg.GoFiles[0]))
		helpersPath.WriteByte('/')
		helpersPath.WriteString(pkgName)
		helpersPath.WriteString("/helpers.c")

		err = os.MkdirAll(filepath.Dir(helpersPath.String()), os.ModePerm)
		if err != nil {
			log.Fatalf("TSGo failed: %v", err)
		}
		err = os.WriteFile(helpersPath.String(), []byte(cHelpers.String()), os.ModePerm)
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
			OGGoType:    "uint64",
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
		parsedSB := strings.Builder{}
		parsedSB.WriteByte('_')
		parsedSB.WriteString(f.Name)
		s.WriteString(parsedSB.String())
		s.WriteString(" := C.GoString(")
		s.WriteString(f.Name)
		s.WriteString(")\n")
		*usedVars = append(*usedVars, parsedSB.String())
	case "unsafe.Pointer":
		parsedSB := strings.Builder{}
		parsedSB.WriteByte('_')
		parsedSB.WriteString(f.Name)
		g.addGoImport(gi, "unsafe")
		arr_dat_type := g.getArrayType(f.ASTField.Type)
		s.WriteString(parsedSB.String())
		s.WriteString(" := unsafe.Slice((*")
		s.WriteString(arr_dat_type)
		s.WriteString(")(")
		s.WriteString(f.Name)
		s.WriteString("), ")
		s.WriteString(f.Name)
		s.WriteString("_len)\n")
		*usedVars = append(*usedVars, parsedSB.String())
	default:
		*usedVars = append(*usedVars, f.Name)
	}
}

// bulk of parsing the function is done here
func (g *PackageGenerator) isResHandle(t ast.Expr) (bool, string) {
	isHandle := false
	structName := ""
	switch t := t.(type) {
	case *ast.StarExpr:
		struct_name := g.getStructName(t.X)
		if g.IsWrappedEnum(struct_name) {
			isHandle = true
			structName = struct_name
		}
	case *ast.Ident:
		struct_name := g.getStructName(t)
		if g.IsWrappedEnum(struct_name) {
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

func (g *PackageGenerator) parseFn(f *ast.FuncDecl) *FFIFunc {
	var ffi_func = &FFIFunc{
		args:       []*ArgHelpers{},
		returns:    []*ResHelpers{},
		isHandleFn: false,
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
				CGoWrapType: "uint64",
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
	}
	return ffi_func
}

func (g *PackageGenerator) writeCGoFieldAccessor(gi *strings.Builder, gh *strings.Builder, ec *strings.Builder, ci *strings.Builder, fmtr cases.Caser, f *StructAccessor, pkgName string, structName string) string {
	used_args := UsedParams{}
	var fnSB strings.Builder
	fnSB.WriteString("//export _")
	fnSB.WriteString(*f.fnName)
	fnSB.WriteByte('\n')
	fnSB.WriteString("func _")
	fnSB.WriteString(*f.fnName)
	fnSB.WriteByte('(')
	// iterate through fn params, generating cgo function decl line
	argLen := len(f.args)
	if argLen > 0 {
		for i, arg := range f.args {
			fnSB.WriteString(arg.Name)
			fnSB.WriteByte(' ')
			fnSB.WriteString(arg.CGoWrapType)
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
		for i, arg := range f.args {
			if i == 0 {
				g.writeIndent(&fnSB, 1)
				fnSB.WriteString("h := cgo.Handle(handle)\n")
				g.writeIndent(&fnSB, 1)
				fnSB.WriteString("s := h.Value().(")
				fnSB.WriteString(pkgName)
				fnSB.WriteByte('.')
				fnSB.WriteString(structName)
				fnSB.WriteString(")\n")
			} else {
				// if `arg.ASTField == nil`, it's a helper arg like `len`, which isn't passed to wrapped Go func
				if arg.ASTField != nil {
					g.addArgHandler(&fnSB, gi, arg, &used_args)
				}
			}
		}
	}
	// write returned value (or intermediary, if necessary)

	tempResType := g.getCgoHandler(f.returns[0].CGoWrapType)
	g.writeIndent(&fnSB, 1)

	fnSB.WriteString("_returned_value := ")
	fnSB.WriteString(tempResType)

	fnSB.WriteByte('(')
	fnSB.WriteString(g.getGoType(f.returns[0].CGoWrapType))
	fnSB.WriteByte('(')
	if f.isStarExpr {
		fnSB.WriteByte('*')
	}
	fnSB.WriteString("s.")
	fnSB.WriteString(*f.name)
	fnSB.WriteString("))\n")

	// TODO: need to improve API so this code is simplified/handles more edge cases
	if tempResType == "C.CString" {
		g.writeIndent(&fnSB, 1)
		fnSB.WriteString("defer C.free(unsafe.Pointer(_returned_value))\n")
	}
	g.writeIndent(&fnSB, 1)
	fnSB.WriteString("return _returned_value\n")

	fnSB.WriteString("}\n\n")
	return fnSB.String()
}

func (g *PackageGenerator) writeCGoFn(gi *strings.Builder, gh *strings.Builder, ec *strings.Builder, ci *strings.Builder, fmtr cases.Caser, f *FFIFunc, name string, pkgName string) string {
	used_args := UsedParams{}
	var fnSB strings.Builder
	fnSB.WriteString("//export _")
	fnSB.WriteString(name)
	fnSB.WriteByte('\n')
	fnSB.WriteString("func _")
	fnSB.WriteString(name)
	fnSB.WriteByte('(')
	// iterate through fn params, generating cgo function decl line
	argLen := len(f.args)
	if argLen > 0 {
		for i, arg := range f.args {
			fnSB.WriteString(arg.Name)
			fnSB.WriteByte(' ')
			fnSB.WriteString(arg.CGoWrapType)
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
	// write returned value (or intermediary, if necessary)
	g.writeIndent(&fnSB, 1)
	if f.isHandleFn {
		fnSB.WriteString("_returned_value := ")
		fnSB.WriteString("C.hackyHandle(C.uintptr_t(cgo.NewHandle(")
	} else {
		g.writeCGoResType(&tempResType, gi, gh, ec, ci, fmtr, *f.returns[0].ASTType, 0, true, pkgName)
		if tempResType.String() == "encodeJSON" {
			fnSB.WriteString("_temp_res_val := ")
		} else {
			fnSB.WriteString("_returned_value := ")
		}
		fnSB.WriteString(tempResType.String())
	}
	fnSB.WriteByte('(')
	fnSB.WriteString(pkgName)
	fnSB.WriteByte('.')
	fnSB.WriteString(name)
	fnSB.WriteString("(")

	// iterate through params (and converted params), writing args passed to function call
	for i, param := range used_args {
		fnSB.WriteString(param)
		if i < len(used_args)-1 {
			fnSB.WriteString(", ")
		}
	}
	if f.isHandleFn {
		fnSB.WriteString(")))")
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
	g.writeIndent(&fnSB, 1)
	fnSB.WriteString("return _returned_value\n")

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
		if func_data.isHandleFn {
			for _, field := range func_data.fieldAccessors {
				var accessorSB strings.Builder
				accessorSB.WriteString("GET_")
				accessorSB.WriteString(*func_data.name)
				accessorSB.WriteString("_")
				accessorSB.WriteString(*field.name)
				name := accessorSB.String()
				field.fnName = &name
				fn_str.WriteString(g.writeCGoFieldAccessor(&goImportsSB, &goHelpersSB, &embeddedCSB, &cImportsSB, caser, field, pkgName, *func_data.name))
			}
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

	// write the generated Cgo code to file
	var outPath strings.Builder
	outPath.WriteString(filepath.Dir(g.pkg.GoFiles[0]))
	outPath.WriteByte('/')
	outPath.WriteString(pkgName)
	outPath.WriteString("/gen_bindings.go")

	err := os.MkdirAll(filepath.Dir(outPath.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
	err = os.WriteFile(outPath.String(), []byte(cg.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("TSGo failed: %v", err)
	}
}
