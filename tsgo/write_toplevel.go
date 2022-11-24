package tsgo

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

type groupContext struct {
	isGroupedDeclaration bool
	doc                  *ast.CommentGroup
	groupValue           string
	groupType            string
	iotaValue            int
	iotaOffset           int
}

func (g *PackageGenerator) writeGroupDecl(s *strings.Builder, decl *ast.GenDecl) {
	// This checks whether the declaration is a group declaration like:
	// const (
	// 	  X = 3
	//    Y = "abc"
	// )
	isGroupedDeclaration := len(decl.Specs) > 1

	if !isGroupedDeclaration {
		g.writeCommentGroupIfNotNil(s, decl.Doc, 0)
	}

	// We need a bit of state to handle syntax like
	// const (
	//   X SomeType = iota
	//   _
	//   Y
	//   Foo string = "Foo"
	//   _
	//   AlsoFoo
	// )
	group := &groupContext{
		isGroupedDeclaration: len(decl.Specs) > 1,
		doc:                  decl.Doc,
		groupType:            "",
		groupValue:           "",
		iotaValue:            -1,
	}

	for i, spec := range decl.Specs {
		isLast := i == len(decl.Specs)-1
		g.writeSpec(s, spec, group, isLast)
	}
}

func (g *PackageGenerator) writeSpec(s *strings.Builder, spec ast.Spec, group *groupContext, isLast bool) {
	// e.g. "type Foo struct {}" or "type Bar = string"
	ts, ok := spec.(*ast.TypeSpec)
	if ok && ts.Name.IsExported() {
		g.writeTypeSpec(s, ts, group)
	}

	// e.g. "const Foo = 123"
	vs, ok := spec.(*ast.ValueSpec)
	if ok {
		g.writeValueSpec(s, vs, group, isLast)

	}
}

// Writing of type specs, which are expressions like
// `type X struct { ... }`
// or
// `type Bar = string`
func (g *PackageGenerator) writeTypeSpec(s *strings.Builder, ts *ast.TypeSpec, group *groupContext) {
	// fmt.Println("name:", ts.Name.Name, "ts:", ts, "group:", group)

	if ts.Doc != nil { // The spec has its own comment, which overrules the grouped comment.
		g.writeCommentGroup(s, ts.Doc, 0)
	} else if group.isGroupedDeclaration {
		g.writeCommentGroupIfNotNil(s, group.doc, 0)
	}

	st, isStruct := ts.Type.(*ast.StructType)
	if isStruct {
		s.WriteString("export interface ")
		s.WriteString(ts.Name.Name)

		if ts.TypeParams != nil {
			g.writeTypeParamsFields(s, ts.TypeParams.List)
		}

		s.WriteString(" {\n")

		g.ffi.StructHelpers[ts.Name.Name] = g.writeStructFields(s, st.Fields.List, 0)
		for _, helper := range g.ffi.StructHelpers[ts.Name.Name] {
			helper.args[0].OGGoType = ts.Name.Name
		}
		s.WriteString("}")
	}

	id, isIdent := ts.Type.(*ast.Ident)
	if isIdent && g.IsEnumStruct(ts.Name.Name) {
		enumName := g.conf.EnumStructs[ts.Name.Name]
		// if names match, dev expects we overwrite the type as enum
		if enumName == "" {
			enumName = ts.Name.Name + "Enum"
		}
		if !strings.EqualFold(enumName, ts.Name.Name) {
			// keeps the original type
			s.WriteString("export type ")
			s.WriteString(ts.Name.Name)
			s.WriteString(" = ")
			s.WriteString(getIdent(id.Name))
			s.WriteString(";")
			s.WriteByte('\n')
		}
		s.WriteString("export enum ")
		s.WriteString(enumName)
		s.WriteString(" {")
	} else if isIdent {
		s.WriteString("export type ")
		s.WriteString(ts.Name.Name)
		s.WriteString(" = ")
		s.WriteString(getIdent(id.Name))
		s.WriteString(";")
	}

	if !isStruct && !isIdent {
		s.WriteString("export type ")
		s.WriteString(ts.Name.Name)
		s.WriteString(" = ")
		g.writeType(s, ts.Type, 0, true)
		s.WriteString(";")
	}

	if ts.Comment != nil {
		s.WriteString(" // " + ts.Comment.Text())
	}
	s.WriteString("\n")
}

// Writing of value specs, which are exported const expressions like
// const SomeValue = 3
func (g *PackageGenerator) writeValueSpec(s *strings.Builder, vs *ast.ValueSpec, group *groupContext, isLast bool) {
	for i, name := range vs.Names {
		// fmt.Println("name:", name.Name, "vs:", vs, "group:", group)
		group.iotaValue = group.iotaValue + 1
		if name.Name == "_" {
			continue
		}
		if !name.IsExported() {
			continue
		}

		if vs.Doc != nil { // The spec has its own comment, which overrules the grouped comment.
			if group.isGroupedDeclaration {
				g.writeCommentGroup(s, vs.Doc, 1)
			} else {
				g.writeCommentGroup(s, vs.Doc, 0)
			}
		} else if group.isGroupedDeclaration {
			g.writeCommentGroupIfNotNil(s, group.doc, 1)
		}

		hasExplicitValue := len(vs.Values) > i
		if hasExplicitValue {
			group.groupType = ""
		}

		// TODO: really need to clean up this logic LOL
		if vs.Type != nil && group.isGroupedDeclaration {
			g.writeIndent(s, 1)
			s.WriteString(name.Name)
			tempSB := &strings.Builder{}
			g.writeType(tempSB, vs.Type, 0, true)
			typeString := tempSB.String()

			group.groupType = typeString
		} else if vs.Type != nil {
			g.writeIndent(s, 1)
			s.WriteString(name.Name)

			tempSB := &strings.Builder{}
			g.writeType(tempSB, vs.Type, 0, true)
			typeString := tempSB.String()

			group.groupType = typeString
		} else if group.isGroupedDeclaration {
			g.writeIndent(s, 1)
			s.WriteString(name.Name)
		} else if group.groupType != "" && !hasExplicitValue {
			if g.IsEnumStruct(group.groupType) {
				s.WriteString(name.Name)
			} else {
				s.WriteString("export const ")
				s.WriteString(name.Name)

				s.WriteString(": ")
				s.WriteString(group.groupType)
			}
		} else {
			if !group.isGroupedDeclaration {
				s.WriteString("export const ")
			} else {
				g.writeIndent(s, 1)
			}
			s.WriteString(name.Name)
		}
		// fmt.Println("name:", name.Name, "vs:", vs, "group:", group, "group.groupType:", group.groupType)

		s.WriteString(" = ")

		if hasExplicitValue {
			val := vs.Values[i]
			tempSB := &strings.Builder{}
			g.writeType(tempSB, val, 0, true)
			valueString := tempSB.String()

			if isProbablyIotaType(valueString) {
				group.iotaOffset = basicIotaOffsetValueParse(valueString)
				group.groupValue = "iota"
				valueString = fmt.Sprint(group.iotaValue + group.iotaOffset)
			} else {
				group.groupValue = valueString
			}
			s.WriteString(valueString)

		} else { // We must use the previous value or +1 in case of iota
			valueString := group.groupValue
			if group.groupValue == "iota" {
				valueString = fmt.Sprint(group.iotaValue + group.iotaOffset)
			}
			s.WriteString(valueString)
		}
		if !isLast && group.isGroupedDeclaration {
			s.WriteByte(',')
		} else if group.groupType != "" && !hasExplicitValue {
			s.WriteByte(';')
		}

		if vs.Comment != nil {
			s.WriteString(" // " + vs.Comment.Text())
		} else {
			s.WriteByte('\n')
		}

		if isLast && group.isGroupedDeclaration {
			s.WriteString("}\n")
		}
	}
}

func (g *PackageGenerator) writeFFIConfig(s *strings.Builder, fd []*ast.FuncDecl, path string) {
	s.WriteString("\n//////////\n")
	// source: misc.go
	s.WriteString("// Generated config for Bun FFI\n")
	s.WriteByte('\n')
	s.WriteString("export const {\n")
	g.writeIndent(s, 1)
	s.WriteString("symbols: {\n")

	count := len(g.ffi.FFIFuncs)
	visited := 0
	for k := range g.ffi.FFIFuncs {
		g.writeIndent(s, 2)
		if !g.ffi.FFIHelpers[k] {
			s.WriteByte('_')
		}
		s.WriteString(k)
		if visited == count-1 {
			s.WriteByte('\n')
			g.writeIndent(s, 1)
			s.WriteString("}\n")
		} else {
			s.WriteString(",\n")
		}
		visited++
	}

	s.WriteString("} = dlopen(import.meta.dir + '/")
	s.WriteString(path)
	s.WriteString("/gen_bindings")
	s.WriteString(".dylib', {\n")
	visited = 0
	for k, v := range g.ffi.FFIFuncs {
		g.writeIndent(s, 1)
		if !g.ffi.FFIHelpers[k] {
			s.WriteByte('_')
		}
		s.WriteString(k)
		s.WriteString(": {\n")
		argLen := len(v.args)
		resLen := len(v.returns)
		if len(v.args) > 0 {
			g.writeIndent(s, 2)
			s.WriteString("args: [")
			for i, arg := range v.args {
				s.WriteString(arg.FFIType)
				if i < argLen-1 {
					s.WriteString(", ")
				}
			}
			s.WriteByte(']')
			if resLen > 0 {
				s.WriteString(",\n")
			} else {
				s.WriteByte('\n')
			}
		}
		if resLen == 1 {
			g.writeIndent(s, 2)
			s.WriteString("returns: ")
			s.WriteString(v.returns[0].FFIType)
			s.WriteByte('\n')
		} else if resLen > 1 {
			var errStr strings.Builder
			errStr.WriteString("Function `")
			errStr.WriteString(k)
			errStr.WriteString("` has more than one return value, which is not supported by Bun FFI...\n")
			errStr.WriteString("Consider adjusting your `tsgo.yaml` config file to inject code before/after the function call in the CGo wrapper fn as you can coerce to a single return value in Go.\n")
			log.Fatalf("TSGo failed: %v", errStr.String())
		}

		g.writeIndent(s, 1)
		if visited == count-1 {
			s.WriteString("}\n")
		} else {
			s.WriteString("},\n")
		}
		visited++
	}
	s.WriteString("})\n")
}
