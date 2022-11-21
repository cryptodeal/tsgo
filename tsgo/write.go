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

// TODO:
// * see if we can handle `complex64` and `complex128`?
// * perhaps do a better job of mapping (no default value??)
func getFFIIdent(s string) string {
	// fmt.Println(s)
	switch s {
	case "bool":
		return "FFIType.bool"
	case "int":
		return "FFIType.int"
	case "int8":
		return "FFIType.i8"
	case "int16":
		return "FFIType.i16"
	case "int32":
		return "FFIType.i32"
	case "int64":
		return "FFIType.i64_fast"
	case "uint":
		return "FFIType.u64_fast"
	case "uint8":
		return "FFIType.u8"
	case "uint16":
		return "FFIType.u16"
	case "uint32":
		return "FFIType.u32"
	case "uint64":
		return "FFIType.u64_fast"
	case "float32":
		return "FFIType.f32"
	case "float64":
		return "FFIType.f64"
	case "string":
		return "FFIType.cstring"
	}
	return "FFIType.pointer"
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
	return "*C.void"
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
	return "*C.void"

}

func (g *PackageGenerator) writeIndent(s *strings.Builder, depth int) {
	for i := 0; i < depth; i++ {
		s.WriteString(g.conf.Indent)
	}
}

// TODO: `writeFFIType` needs a major overhaul as logic is copied directly from writeType
func (g *PackageGenerator) writeFFIType(s *strings.Builder, t ast.Expr, depth int, optionalParens bool) {
	switch t := t.(type) {
	case *ast.StarExpr:
		s.WriteString("FFIType.cstring")
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("FFIType.cstring")
			break
		}
		s.WriteString("FFIType.ptr")
	case *ast.StructType:
		s.WriteString("{\n")
		g.writeStructFields(s, t.Fields.List, depth+1)
		g.writeIndent(s, depth+1)
		s.WriteByte('}')
	case *ast.Ident:
		if t.String() == "any" {
			s.WriteString(getFFIIdent(g.conf.FallbackType))
		} else {
			s.WriteString(getFFIIdent(t.String()))
		}
	case *ast.SelectorExpr:
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

// used to add handlers for data types on an as needed basis (reduce code bloat)
func (g *PackageGenerator) writeCGoResType(s *strings.Builder, cg *strings.Builder, gh *strings.Builder, ec *strings.Builder, fmtr cases.Caser, t ast.Expr, depth int, optionalParens bool, pkgName string) {
	switch t := t.(type) {
	case *ast.StarExpr:
		g.addJSONEncoder(gh, cg)
		s.WriteString("encodeJSON")
	case *ast.ArrayType:
		fmt.Println("writeCGoResType - *ast.ArrayType")
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("C.CString")
			break
		} else if ok {
			g.addGoImport(cg, "unsafe")
			g.addPtrTrckr(gh)
			g.addDisposePtr(gh, cg)
			g.addCDisposeHelpers(pkgName)
			g.addArraySize(gh, cg)
			dat_type := g.getArrayType(t)
			handler := g.writeCArrayHandler(gh, ec, dat_type, fmtr)
			s.WriteString(handler)
		} else {
			fmt.Println("unknown ptr type; returning as unsafe.Pointer (void*)")
			s.WriteString("unsafe.Pointer")
		}
	case *ast.StructType:
		fmt.Println("writeCGoResType - *ast.StructType")
		s.WriteString("{\n")
		g.writeStructFields(s, t.Fields.List, depth+1)
		g.writeIndent(s, depth+1)
		s.WriteByte('}')
	case *ast.Ident:
		fmt.Println("writeCGoResType - *ast.Ident")
		if t.String() == "any" {
			s.WriteString(getCGoTypeHandler(g.conf.FallbackType))
		} else {
			s.WriteString(getCGoTypeHandler(t.String()))
		}
	case *ast.MapType:
		g.addJSONEncoder(gh, cg)
		s.WriteString("encodeJSON")
	case *ast.BasicLit:
		fmt.Println("writeCGoResType - *ast.BasicLit")
		s.WriteString(t.Value)
	case *ast.ParenExpr:
		fmt.Println("writeCGoResType - *ast.ParenExpr")
		s.WriteByte('(')
		g.writeType(s, t.X, depth, false)
		s.WriteByte(')')
	case *ast.BinaryExpr:
		fmt.Println("writeCGoResType - *ast.BinaryExpr")
		g.writeType(s, t.X, depth, false)
		s.WriteByte(' ')
		s.WriteString(t.Op.String())
		s.WriteByte(' ')
		g.writeType(s, t.Y, depth, false)
	case *ast.InterfaceType:
		fmt.Println("writeCGoResType - *ast.InterfaceType")
		g.writeInterfaceFields(s, t.Methods.List, depth+1)
	case *ast.CallExpr, *ast.FuncType, *ast.ChanType:
		fmt.Println("writeCGoResType - *ast.CallExpr, *ast.FuncType, *ast.ChanType")
		s.WriteString(g.conf.FallbackType)
	case *ast.UnaryExpr:
		fmt.Println("writeCGoResType - *ast.UnaryExpr")
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
		fmt.Println("writeCGoResType - *ast.IndexListExpr")
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
		fmt.Println("writeCGoResType - *ast.IndexExpr")
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
		fmt.Println("writeCGoResType - *ast.ArrayType")
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

func (g *PackageGenerator) writeType(s *strings.Builder, t ast.Expr, depth int, optionalParens bool) {
	switch t := t.(type) {
	case *ast.StarExpr:
		if optionalParens {
			s.WriteByte('(')
		}
		g.writeType(s, t.X, depth, false)
		s.WriteString(" | undefined")
		if optionalParens {
			s.WriteByte(')')
		}
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("string")
			break
		}
		g.writeType(s, t.Elt, depth, true)
		s.WriteString("[]")
	case *ast.StructType:
		s.WriteString("{\n")
		g.writeStructFields(s, t.Fields.List, depth+1)
		g.writeIndent(s, depth+1)
		s.WriteByte('}')
	case *ast.Ident:
		if t.String() == "any" {
			s.WriteString(getIdent(g.conf.FallbackType))
		} else {
			s.WriteString(getIdent(t.String()))
		}
	case *ast.SelectorExpr:
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

func (g *PackageGenerator) writeStructFields(s *strings.Builder, fields []*ast.Field, depth int) {
	for _, f := range fields {
		// fmt.Println(f.Type)
		optional := false
		required := false
		readonly := false

		var fieldName string
		if len(f.Names) != 0 && f.Names[0] != nil && len(f.Names[0].Name) != 0 {
			fieldName = f.Names[0].Name
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

		switch t := f.Type.(type) {
		case *ast.StarExpr:
			optional = !required
			f.Type = t.X
		}

		if optional {
			s.WriteByte('?')
		}

		s.WriteString(": ")

		if tstype == "" {
			g.writeType(s, f.Type, depth, false)
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

	}
}
