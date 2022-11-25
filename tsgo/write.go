package tsgo

import (
	"fmt"
	"regexp"
	"strings"

	"go/ast"
	"go/token"

	"github.com/fatih/structtag"
	"golang.org/x/text/cases"
)

var validJSNameRegexp = regexp.MustCompile(`(?m)^[\pL_][\pL\pN_]*$`)

func validJSName(n string) bool {
	return validJSNameRegexp.MatchString(n)
}

func getIdent(s string) string {
	switch s {
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"complex64", "complex128":
		return "number /* " + s + " */"
	}
	return s
}

func getByteSize(s string) int {
	switch s {
	case "float16", "int16", "uint16":
		return 2
	case "float32", "int32", "uint32":
		return 4
	case "float64", "int64", "uint64":
		return 8
	}
	return 1
}

// TODO:
// * see if we can handle `complex64` and `complex128`?
// * perhaps do a better job of mapping (no default value??)
func getFFIIdent(s string) string {
	// fmt.Println(s)
	switch s {
	case "bool":
		return "FFIType.bool"
	case "int", "C.int":
		return "FFIType.int"
	case "int8", "C.int8_t":
		return "FFIType.i8"
	case "int16", "C.int16_t":
		return "FFIType.i16"
	case "int32", "C.int32_t":
		return "FFIType.i32"
	case "int64", "C.int64_t":
		return "FFIType.i64_fast"
	case "uint", "C.uint":
		return "FFIType.u64_fast"
	case "uint8", "C.uint8_t":
		return "FFIType.u8"
	case "uint16", "C.uint16_t":
		return "FFIType.u16"
	case "uint32", "C.uint32_t":
		return "FFIType.u32"
	case "uint64", "C.uint64_t":
		return "FFIType.u64_fast"
	case "float32", "C.float":
		return "FFIType.f32"
	case "float64", "C.double":
		return "FFIType.f64"
	case "string", "*C.char":
		return "FFIType.cstring"
	}
	return "FFIType.ptr"
}

// TODO:
// * see if we can handle `complex64` and `complex128`?
// * perhaps do a better job of mapping (no default value??)
func (g *PackageGenerator) getCgoHandler(s string) string {
	// fmt.Println(s)
	switch s {
	case "int", "C.int":
		return "C.int"
	case "int8", "C.int8_t":
		return "C.int8_t"
	case "int16", "C.int16_t":
		return "C.int16_t"
	case "int32", "C.int32_t":
		return "C.int32_t"
	case "int64", "C.int64_t":
		return "C.int64_t"
	case "uint", "C.uint":
		return "C.uint"
	case "uint8", "C.uint8_t":
		return "C.uint8_t"
	case "uint16", "C.uint16_t":
		return "C.uint16_t"
	case "uint32", "C.uint32_t":
		return "C.uint32_t"
	case "uint64", "C.uint64_t":
		return "C.uint64_t"
	case "float32", "C.float":
		return "C.float"
	case "float64", "C.double":
		return "C.double"
	case "string", "*C.char":
		return "C.CString"
	case "uintptr", "C.uintptr_t":
		return "C.uintptr_t"
	}
	return "unsafe.Pointer"
}

func (g *PackageGenerator) getJSFromFFIType(t string) string {
	switch t {
	case "FFIType.bool":
		return "boolean"
	case "FFIType.int", "FFIType.i8", "FFIType.i16", "FFIType.i32", "FFIType.i64_fast", "FFIType.u8", "FFIType.u16", "FFIType.u32", "FFIType.u64_fast", "FFIType.f32", "FFIType.f64":
		return "number"
	case "FFIType.cstring":
		return "string"
	case "FFIType.ptr":
		return "any"
	}
	return t
}

func (g *PackageGenerator) getGoType(t string) string {
	switch t {
	case "*C.char":
		// not valid type
		return "string"
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
	case "double":
		return "float64"
	case "string":
		return "*C.char"
	}
	return "unsafe.Pointer"
}

// TODO:
// * see if we can handle `complex64` and `complex128`?
// * perhaps do a better job of mapping (no default value??)
func getCGoIdent(s string) string {
	// fmt.Println(s)
	switch s {
	case "bool":
		// not valid type
		return "C.bool"
	case "int":
		return "C.int"
	case "int8":
		return "C.int8_t"
	case "int16":
		return "C.int16_t"
	case "int32":
		return "C.int32_t"
	case "int64":
		return "C.int64_t"
	case "uint":
		return "C.uint"
	case "uint8":
		return "C.uint8_t"
	case "uint16":
		return "C.uint16_t"
	case "uint32":
		return "C.uint32_t"
	case "uint64":
		return "C.uint64_t"
	case "float32":
		return "C.float"
	case "float64":
		return "C.double"
	case "string":
		return "*C.char"
	}
	return "unsafe.Pointer"
}

func getCGoTypeHandler(s string) string {
	switch s {
	case "bool":
		// not valid type
		return "C.bool"
	case "int":
		return "C.int"
	case "int8":
		return "C.int8_t"
	case "int16":
		return "C.int16_t"
	case "int32":
		return "C.int32_t"
	case "int64":
		return "C.int64_t"
	case "uint":
		return "C.uint"
	case "uint8":
		return "C.uint8_t"
	case "uint16":
		return "C.uint16_t"
	case "uint32":
		return "C.uint32_t"
	case "uint64":
		return "C.uint64_t"
	case "float32":
		return "C.float"
	case "float64":
		return "C.double"
	case "string":
		return "C.CString"
	}
	return "unsafe.Pointer"

}

func (g *PackageGenerator) writeIndent(s *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		s.WriteString(g.conf.Indent)
	}
}

// used to add handlers for data types on an as needed basis (reduce code bloat)
func (g *PackageGenerator) writeCGoResType(s *strings.Builder, cg *strings.Builder, gh *strings.Builder, ec *strings.Builder, ci *strings.Builder, fmtr cases.Caser, t ast.Expr, depth int, optionalParens bool, pkgName string) {
	switch t := t.(type) {
	case *ast.StarExpr:
		// fmt.Println("writeCGoResType - *ast.StarExpr", t)
		struct_name := g.getStructName(t.X)
		if g.IsWrappedEnum(struct_name) {
			g.addCImport(ci, "stdint.h", false)
			g.addGoImport(cg, "runtime/cgo")
			s.WriteString("C.hackyHandle(C.uintptr_t(cgo.NewHandle(")
		} else {
			g.addJSONEncoder(gh, cg)
			s.WriteString("encodeJSON")
		}
	case *ast.ArrayType:
		// fmt.Println("writeCGoResType - *ast.ArrayType", t)
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("C.CString")
			break
		} else if ok {
			g.addGoImport(cg, "unsafe")
			g.addPtrTrckr(gh)
			g.addDisposePtr(gh, cg)
			g.addCDisposeHelpers(ci, pkgName)
			g.addArraySize(gh, cg)
			dat_type := g.getArrayType(t)
			handler := g.writeCArrayHandler(gh, ec, dat_type, fmtr)
			s.WriteString(handler)
		} else {
			// fmt.Println("unknown ptr type; returning as unsafe.Pointer (void*)")
			s.WriteString("unsafe.Pointer")
		}
	case *ast.StructType:
		s.WriteString(g.getStructName(t))
	case *ast.Ident:
		// fmt.Println("writeCGoResType - *ast.Ident", t)
		if t.String() == "any" {
			s.WriteString(getCGoTypeHandler(g.conf.FallbackType))
		} else {
			s.WriteString(getCGoTypeHandler(t.String()))
		}
	case *ast.MapType:
		// fmt.Println("writeCGoResType - *ast.MapType", t)
		g.addJSONEncoder(gh, cg)
		s.WriteString("encodeJSON")
	case *ast.BasicLit:
		// fmt.Println("writeCGoResType - *ast.BasicLit", t)
		s.WriteString(t.Value)
	case *ast.ParenExpr:
		// fmt.Println("writeCGoResType - *ast.ParenExpr", t)
		s.WriteByte('(')
		g.writeType(s, t.X, depth, false)
		s.WriteByte(')')
	case *ast.BinaryExpr:
		// fmt.Println("writeCGoResType - *ast.BinaryExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte(' ')
		s.WriteString(t.Op.String())
		s.WriteByte(' ')
		g.writeType(s, t.Y, depth, false)
	case *ast.InterfaceType:
		// fmt.Println("writeCGoResType - *ast.InterfaceType", t)
		g.writeInterfaceFields(s, t.Methods.List, depth+1)
	case *ast.CallExpr, *ast.FuncType, *ast.ChanType:
		// fmt.Println("writeCGoResType - *ast.CallExpr, *ast.FuncType, *ast.ChanType", t)
		s.WriteString(g.conf.FallbackType)
	case *ast.UnaryExpr:
		// fmt.Println("writeCGoResType - *ast.UnaryExpr", t)
		if t.Op == token.TILDE {
			// We just ignore the tilde token, in Typescript extended types are
			// put into the generic typing itself, which we can't support yet.
			g.writeType(s, t.X, depth, false)
		} else {
			err := fmt.Errorf("unhandled unary expr: %v\n %T", t, t)
			fmt.Println(err)
			panic(err)
		}
	case *ast.IndexListExpr:
		// fmt.Println("writeCGoResType - *ast.IndexListExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		for i, index := range t.Indices {
			g.writeType(s, index, depth, false)
			if i != len(t.Indices)-1 {
				s.WriteString(", ")
			}
		}
		s.WriteByte('>')
	case *ast.IndexExpr:
		// fmt.Println("writeCGoResType - *ast.IndexExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		g.writeType(s, t.Index, depth, false)
		s.WriteByte('>')
	default:
		err := fmt.Errorf("unhandled: %s\n %T", t, t)
		fmt.Println(err)
		panic(err)
	}
}

// TODO: `writeFFIType` needs a major overhaul as logic is copied directly from writeType
func (g *PackageGenerator) writeCGoType(s *strings.Builder, t ast.Expr, depth int, optionalParens bool) {
	switch t := t.(type) {
	case *ast.StarExpr:
		s.WriteString("*C.char")
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("*C.char")
			break
		}
		s.WriteString("unsafe.Pointer")
	case *ast.StructType:

		s.WriteString("{\n")
		g.writeStructFields(s, t.Fields.List, depth+1)
		g.writeIndent(s, depth+1)
		s.WriteByte('}')
	case *ast.Ident:
		if t.String() == "any" {
			s.WriteString(getCGoIdent(g.conf.FallbackType))
		} else {
			s.WriteString(getCGoIdent(t.String()))
		}

	case *ast.SelectorExpr:
		// fmt.Println("writeType - *ast.SelectorExpr", t)
		// e.g. `time.Time`
		longType := fmt.Sprintf("%s.%s", t.X, t.Sel)
		mappedTsType, ok := g.conf.FFITypeMappings[longType]
		if ok {
			s.WriteString(mappedTsType)
		} else { // For unknown types we use the fallback type
			s.WriteString(g.conf.FFIFallbackType)
			s.WriteString(" /* ")
			s.WriteString(longType)
			s.WriteString(" */")
		}
	case *ast.MapType:
		s.WriteString("{ [key: ")
		g.writeType(s, t.Key, depth, false)
		s.WriteString("]: ")
		g.writeType(s, t.Value, depth, false)
		s.WriteByte('}')
	case *ast.BasicLit:
		s.WriteString(t.Value)
	case *ast.ParenExpr:
		s.WriteByte('(')
		g.writeType(s, t.X, depth, false)
		s.WriteByte(')')
	case *ast.BinaryExpr:
		g.writeType(s, t.X, depth, false)
		s.WriteByte(' ')
		s.WriteString(t.Op.String())
		s.WriteByte(' ')
		g.writeType(s, t.Y, depth, false)
	case *ast.InterfaceType:
		g.writeInterfaceFields(s, t.Methods.List, depth+1)
	case *ast.CallExpr, *ast.FuncType, *ast.ChanType:
		s.WriteString(g.conf.FallbackType)
	case *ast.UnaryExpr:
		if t.Op == token.TILDE {
			// We just ignore the tilde token, in Typescript extended types are
			// put into the generic typing itself, which we can't support yet.
			g.writeType(s, t.X, depth, false)
		} else {
			err := fmt.Errorf("unhandled unary expr: %v\n %T", t, t)
			fmt.Println(err)
			panic(err)
		}
	case *ast.IndexListExpr:
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		for i, index := range t.Indices {
			g.writeType(s, index, depth, false)
			if i != len(t.Indices)-1 {
				s.WriteString(", ")
			}
		}
		s.WriteByte('>')
	case *ast.IndexExpr:
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		g.writeType(s, t.Index, depth, false)
		s.WriteByte('>')
	default:
		err := fmt.Errorf("unhandled: %s\n %T", t, t)
		fmt.Println(err)
		panic(err)
	}
}

func (g *PackageGenerator) getArrayType(t ast.Expr) string {
	switch t := t.(type) {
	case *ast.ArrayType:
		// fmt.Println("writeCGoResType - *ast.ArrayType")
		if v, ok := t.Elt.(*ast.Ident); ok {
			return v.String()
		}
		err := fmt.Errorf("unhandled: no ident found in `getArrayType`")
		fmt.Println(err)
		panic(err)
	default:
		err := fmt.Errorf("unhandled: %s\n %T", t, t)
		fmt.Println(err)
		panic(err)
	}
}

func (g *PackageGenerator) getStructName(t ast.Expr) string {
	if v, ok := t.(*ast.Ident); ok {
		// fmt.Println(v)
		return v.String()
	}
	return ""
}

func (g *PackageGenerator) writeType(s *strings.Builder, t ast.Expr, depth int, optionalParens bool) {
	switch t := t.(type) {
	case *ast.StarExpr:
		// fmt.Println("writeType - *ast.StarExpr", t)
		if optionalParens {
			s.WriteByte('(')
		}
		g.writeType(s, t.X, depth, false)
		s.WriteString(" | undefined")
		if optionalParens {
			s.WriteByte(')')
		}
	case *ast.ArrayType:
		// fmt.Println("writeType - *ast.ArrayType", t)
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("string")
			break
		}
		g.writeType(s, t.Elt, depth, true)
		s.WriteString("[]")
	case *ast.StructType:
		// fmt.Println("writeType - *ast.StructType", t)
		s.WriteString("{\n")
		g.writeStructFields(s, t.Fields.List, depth+1)
		g.writeIndent(s, depth+1)
		s.WriteByte('}')
	case *ast.Ident:
		// fmt.Println("writeType - *ast.Ident", t)
		if t.String() == "any" {
			s.WriteString(getIdent(g.conf.FallbackType))
		} else {
			s.WriteString(getIdent(t.String()))
		}
	case *ast.SelectorExpr:
		// fmt.Println("writeType - *ast.SelectorExpr", t)
		// e.g. `time.Time`
		longType := fmt.Sprintf("%s.%s", t.X, t.Sel)
		mappedTsType, ok := g.conf.TypeMappings[longType]
		if ok {
			s.WriteString(mappedTsType)
		} else { // For unknown types we use the fallback type
			s.WriteString(g.conf.FallbackType)
			s.WriteString(" /* ")
			s.WriteString(longType)
			s.WriteString(" */")
		}
	case *ast.MapType:
		// fmt.Println("writeType - *ast.MapType", t)
		s.WriteString("{ [key: ")
		g.writeType(s, t.Key, depth, false)
		s.WriteString("]: ")
		g.writeType(s, t.Value, depth, false)
		s.WriteByte('}')
	case *ast.BasicLit:
		// fmt.Println("writeType - *ast.BasicLit", t)
		s.WriteString(t.Value)
	case *ast.ParenExpr:
		// fmt.Println("writeType - *ast.ParenExpr", t)
		s.WriteByte('(')
		g.writeType(s, t.X, depth, false)
		s.WriteByte(')')
	case *ast.BinaryExpr:
		// fmt.Println("writeType - *ast.BinaryExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte(' ')
		s.WriteString(t.Op.String())
		s.WriteByte(' ')
		g.writeType(s, t.Y, depth, false)
	case *ast.InterfaceType:
		// fmt.Println("writeType - *ast.InterfaceType", t)
		g.writeInterfaceFields(s, t.Methods.List, depth+1)
	case *ast.CallExpr, *ast.FuncType, *ast.ChanType:
		// fmt.Println("writeType - *ast.CallExpr, *ast.FuncType, *ast.ChanType", t)
		s.WriteString(g.conf.FallbackType)
	case *ast.UnaryExpr:
		// fmt.Println("writeType - *ast.UnaryExpr", t)
		if t.Op == token.TILDE {
			// We just ignore the tilde token, in Typescript extended types are
			// put into the generic typing itself, which we can't support yet.
			g.writeType(s, t.X, depth, false)
		} else {
			err := fmt.Errorf("unhandled unary expr: %v\n %T", t, t)
			fmt.Println(err)
			panic(err)
		}
	case *ast.IndexListExpr:
		// fmt.Println("writeType - *ast.IndexListExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		for i, index := range t.Indices {
			g.writeType(s, index, depth, false)
			if i != len(t.Indices)-1 {
				s.WriteString(", ")
			}
		}
		s.WriteByte('>')
	case *ast.IndexExpr:
		// fmt.Println("writeType - *ast.IndexExpr", t)
		g.writeType(s, t.X, depth, false)
		s.WriteByte('<')
		g.writeType(s, t.Index, depth, false)
		s.WriteByte('>')
	default:
		err := fmt.Errorf("unhandled: %s\n %T", t, t)
		fmt.Println(err)
		panic(err)
	}
}

func (g *PackageGenerator) writeTypeParamsFields(s *strings.Builder, fields []*ast.Field) {
	s.WriteByte('<')
	for i, f := range fields {
		for j, ident := range f.Names {
			s.WriteString(ident.Name)
			s.WriteString(" extends ")
			g.writeType(s, f.Type, 0, true)

			if i != len(fields)-1 || j != len(f.Names)-1 {
				s.WriteString(", ")
			}
		}
	}
	s.WriteByte('>')
}

func (g *PackageGenerator) writeInterfaceFields(s *strings.Builder, fields []*ast.Field, depth int) {
	if len(fields) == 0 { // Type without any fields (probably only has methods)
		s.WriteString(g.conf.FallbackType)
		return
	}
	s.WriteByte('\n')
	for _, f := range fields {
		if _, isFunc := f.Type.(*ast.FuncType); isFunc {
			continue
		}
		g.writeCommentGroupIfNotNil(s, f.Doc, depth+1)
		g.writeIndent(s, depth+1)
		g.writeType(s, f.Type, depth, false)

		if f.Comment != nil {
			s.WriteString(" // ")
			s.WriteString(f.Comment.Text())
		}
	}
}

func (g *PackageGenerator) writeStructFields(s *strings.Builder, fields []*ast.Field, depth int) []*StructAccessor {
	struct_fields := []*StructAccessor{}
	for _, f := range fields {
		// fmt.Println(f.Type)
		optional := false
		required := false
		readonly := false
		_, dType := g.isTypedArray(f.Type)

		var ptr_arg = &ArgHelpers{
			Name:        "handle",
			FFIType:     "FFIType.ptr",
			CGoWrapType: "C.uintptr_t",
			OGGoType:    "unsafe.Pointer",
			ASTField:    f,
		}

		var field_func = &StructAccessor{
			args:      []*ArgHelpers{ptr_arg},
			returns:   []*ResHelpers{},
			arrayType: &dType,
		}

		var fieldName string
		if len(f.Names) != 0 && f.Names[0] != nil && len(f.Names[0].Name) != 0 {
			fieldName = f.Names[0].Name
			field_func.name = &fieldName
		}
		if len(fieldName) == 0 || 'A' > fieldName[0] || fieldName[0] > 'Z' {
			continue
		}

		var name string
		var tstype string
		if f.Tag != nil {
			tags, err := structtag.Parse(f.Tag.Value[1 : len(f.Tag.Value)-1])
			if err != nil {
				panic(err)
			}

			jsonTag, err := tags.Get("json")
			if err == nil {
				name = jsonTag.Name
				if name == "-" {
					continue
				}

				optional = jsonTag.HasOption("omitempty")
			}
			tstypeTag, err := tags.Get("tstype")
			if err == nil {
				tstype = tstypeTag.Name
				if tstype == "-" {
					continue
				}
				required = tstypeTag.HasOption("required")
				readonly = tstypeTag.HasOption("readonly")
			}
		}

		if len(name) == 0 {
			name = fieldName
		}

		g.writeCommentGroupIfNotNil(s, f.Doc, depth+1)

		g.writeIndent(s, depth+1)
		quoted := !validJSName(name)
		if quoted {
			s.WriteByte('\'')
		}
		if readonly {
			s.WriteString("readonly ")
		}
		s.WriteString(name)
		if quoted {
			s.WriteByte('\'')
		}

		isStarExpr := false

		switch t := f.Type.(type) {
		case *ast.StarExpr:
			isStarExpr = true
			optional = !required
			f.Type = t.X
		}

		isHandleFn, structName := g.isResHandle(f.Type)

		if optional {
			s.WriteByte('?')
		}

		s.WriteString(": ")

		if tstype == "" {
			g.writeType(s, f.Type, depth, false)
			var tempSB strings.Builder
			g.writeCGoType(&tempSB, f.Type, depth, false)
			cgoType := tempSB.String()

			longType := fmt.Sprintf("%s", f.Type)
			// fmt.Println(longType)
			if val, ok := g.ffi.TypeHelpers[longType]; ok {
				cgoType = val
			}
			var res_helper = &ResHelpers{
				FFIType:     getFFIIdent(cgoType),
				CGoWrapType: cgoType,
				OGGoType:    cgoType,
				ASTType:     &f.Type,
			}
			if isHandleFn {
				field_func.isHandleFn = &structName
			}
			field_func.isOptional = optional
			field_func.isStarExpr = isStarExpr

			field_func.returns = append(field_func.returns, res_helper)
		} else {
			s.WriteString(tstype)
		}
		s.WriteByte(';')

		if f.Comment != nil {
			// Line comment is present, that means a comment after the field.
			s.WriteString(" // ")
			s.WriteString(f.Comment.Text())
		} else {
			s.WriteByte('\n')
		}

		struct_fields = append(struct_fields, field_func)
	}
	return struct_fields
}
