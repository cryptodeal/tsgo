package tsgo

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	if ts.Doc != nil { // The spec has its own comment, which overrules the grouped comment.
		g.writeCommentGroup(s, ts.Doc, 0)
	} else if group.isGroupedDeclaration {
		g.writeCommentGroupIfNotNil(s, group.doc, 0)
	}

	st, isStruct := ts.Type.(*ast.StructType)
	if isStruct {
		s.WriteString(fmt.Sprintf("export interface %s", ts.Name.Name))
		if ts.TypeParams != nil {
			g.writeTypeParamsFields(s, ts.TypeParams.List)
		}
		s.WriteString(" {\n")
		g.ffi.StructHelpers[ts.Name.Name] = g.writeStructFields(s, st.Fields.List, 0)
		for _, helper := range g.ffi.StructHelpers[ts.Name.Name] {
			helper.args[0].OGGoType = ts.Name.Name
		}
		s.WriteByte('}')
	}

	id, isIdent := ts.Type.(*ast.Ident)
	if isIdent {
		g.ffi.TypeHelpers[ts.Name.Name] = getCGoIdent(id.Name)
		g.TSHelpers.EnumStructs[ts.Name.Name] = []*EnumField{}
		s.WriteString(fmt.Sprintf("export type %s = %s;\n", ts.Name.Name, getIdent(id.Name)))
	}

	if !isStruct && !isIdent {
		var tempSB = &strings.Builder{}
		g.writeCGoType(tempSB, ts.Type, 0, false)
		g.ffi.TypeHelpers[ts.Name.Name] = tempSB.String()
		s.WriteString(fmt.Sprintf("export type %s = ", ts.Name.Name))
		g.writeType(s, ts.Type, 0, true)
		s.WriteByte(';')
	}

	if ts.Comment != nil {
		s.WriteString(fmt.Sprintf(" // %s", ts.Comment.Text()))
	}
	s.WriteByte('\n')
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
			g.writeCommentGroup(s, vs.Doc, 0)
		} else if group.isGroupedDeclaration {
			g.writeCommentGroupIfNotNil(s, group.doc, 0)
		}

		hasExplicitValue := len(vs.Values) > i
		if hasExplicitValue {
			group.groupType = ""
		}

		s.WriteString(fmt.Sprintf("export const %s", name.Name))
		if vs.Type != nil {
			s.WriteString(": ")
			tempSB := &strings.Builder{}
			g.writeType(tempSB, vs.Type, 0, true)
			typeString := tempSB.String()
			s.WriteString(typeString)
			group.groupType = typeString
		} else if group.groupType != "" && !hasExplicitValue {
			s.WriteString(fmt.Sprintf(": %s", group.groupType))
		}

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

		s.WriteByte(';')
		if vs.Comment != nil {
			s.WriteString(fmt.Sprintf(" // %s", vs.Comment.Text()))
		} else {
			s.WriteByte('\n')
		}
	}
}

type InitStructParam struct {
	Name     string
	IsPtr    bool
	IsStruct bool
}

func (g *PackageGenerator) writeInitMethod(s *strings.Builder, cw *ClassWrapper, fmtr cases.Caser) {
	g.writeIndent(s, 1)
	s.WriteString(fmt.Sprintf("static init(struct: %s): _%s {\n", *cw.name, *cw.name))
	var constDestFields = []*StructAccessor{}
	var letDestFields = []*StructAccessor{}
	var usedArgs = []*InitStructParam{}

	for _, f := range cw.fieldAccessors {
		if f.returns[0].FFIType == "FFIType.cstring" || f.isHandleFn != nil {
			constDestFields = append(constDestFields, f)
		} else if !isParsingRequired(f.returns[0].FFIType) {
			constDestFields = append(constDestFields, f)
		} else {
			letDestFields = append(letDestFields, f)
		}
	}
	fieldCount := len(constDestFields)
	if fieldCount > 0 {
		g.writeIndent(s, 2)
		s.WriteString("const {")
		for i, c := range constDestFields {
			if i < fieldCount-1 {
				s.WriteString(fmt.Sprintf(" %s,", *c.name))
			} else {
				s.WriteString(fmt.Sprintf(" %s", *c.name))
			}
		}
		s.WriteString(" } = struct;\n")
	}

	fieldCount = len(letDestFields)
	if fieldCount > 0 {
		g.writeIndent(s, 2)
		s.WriteString("let {")
		for i, l := range letDestFields {
			if i < fieldCount-1 {
				s.WriteString(fmt.Sprintf(" %s,", *l.name))
			} else {
				s.WriteString(fmt.Sprintf(" %s", *l.name))
			}
		}
		s.WriteString(" } = struct;\n")
	}

	// parse fields that require `const` declaration (e.g. strings)
	for _, c := range constDestFields {
		if c.returns[0].FFIType == "FFIType.cstring" {
			g.writeIndent(s, 2)
			arg_name := fmt.Sprintf("_%s", *c.name)
			s.WriteString(fmt.Sprintf("const %s = Buffer.from(%s + '/%d', %q);\n", arg_name, *c.name, 0, "utf8"))
			var param = &InitStructParam{Name: arg_name, IsPtr: true}
			usedArgs = append(usedArgs, param)
		} else if c.isHandleFn != nil {
			g.writeIndent(s, 2)
			arg_name := fmt.Sprintf("_%s", *c.name)
			s.WriteString(fmt.Sprintf("const %s = _%s.init(%s);\n", arg_name, *c.isHandleFn, *c.name))
			var param = &InitStructParam{Name: arg_name, IsStruct: true, IsPtr: false}
			usedArgs = append(usedArgs, param)
		} else {
			g.writeIndent(s, 2)
			var param = &InitStructParam{Name: *c.name, IsPtr: false}
			usedArgs = append(usedArgs, param)
		}
	}

	// parse fields that don't require `const` declaration
	for _, l := range letDestFields {
		var param = &InitStructParam{Name: *l.name, IsPtr: false}
		if l.arrayType != nil && *l.arrayType != "" {
			g.writeIndent(s, 2)
			param.IsPtr = true
			usedArgs = append(usedArgs, param)
			var len_helper = &InitStructParam{Name: fmt.Sprintf("%s.length", *l.name), IsPtr: false}
			usedArgs = append(usedArgs, len_helper)
			s.WriteString(fmt.Sprintf("if (!(%s instanceof %sArray)) %s = new %sArray(%s);\n", *l.name, fmtr.String(*l.arrayType), *l.name, fmtr.String(*l.arrayType), *l.name))
		} else if l.isHandleFn != nil {
			g.writeIndent(s, 2)
			param.IsStruct = true
			param.IsPtr = true
			usedArgs = append(usedArgs, param)
			s.WriteString(fmt.Sprintf("if (!(%s instanceof _%s)) %s = _%s.init(%s);\n", *l.name, *l.isHandleFn, *l.name, *l.isHandleFn, *l.name))
		}
	}

	// write return fn
	g.writeIndent(s, 2)
	s.WriteString(fmt.Sprintf("return new _%s(_INIT_%s(", *cw.name, *cw.name))
	argCount := len(usedArgs)
	i := 0
	for _, v := range cw.fieldAccessors {
		for _, arg := range usedArgs {
			if !strings.Contains(arg.Name, *v.name) {
				continue
			} else {
				i++
				Fmt := ""
				if i < argCount-1 {
					Fmt = ", "
				}
				if arg.IsStruct {
					s.WriteString(fmt.Sprintf("%s.ptr%s", arg.Name, Fmt))
				} else if arg.IsPtr && v.returns[0].FFIType != "FFIType.cstring" {
					// increment twice, account for array length helper in `args`
					i++
					s.WriteString(fmt.Sprintf("ptr(%s), %s.length%s", arg.Name, arg.Name, Fmt))
				} else if arg.IsPtr {
					s.WriteString(fmt.Sprintf("ptr(%s)%s", arg.Name, Fmt))
				} else {
					s.WriteString(fmt.Sprintf("%s%s", arg.Name, Fmt))
				}
				break
			}
		}
	}
	s.WriteString("));\n")
	g.writeIndent(s, 1)
	s.WriteString("}\n\n")
}

func (g *PackageGenerator) writeAccessorClasses(s *strings.Builder, class_wrappers *[]*ClassWrapper, fmtr cases.Caser) {
	if len(*class_wrappers) > 0 {
		s.WriteString("const registry = new FinalizationRegistry((disp: { cb: (ptr: number) => void; ptr: number}) => {\n")
		g.writeIndent(s, 1)
		s.WriteString("const { cb, ptr } = disp;\n")
		g.writeIndent(s, 1)
		s.WriteString("return cb(ptr);\n")
		s.WriteString("});\n\n")
	}

	// Write the class wrappers
	struct_wrappers := map[string]bool{}
	for _, c := range *class_wrappers {
		if v, ok := g.ffi.StructHelpers[*c.name]; ok && len(v) > 0 {
			if !struct_wrappers[*c.name] {
				s.WriteString(fmt.Sprintf("export class _%s {\n", *c.name))
				g.writeIndent(s, 1)
				s.WriteString("private _ptr: number;\n\n")
				g.writeIndent(s, 1)
				s.WriteString("constructor(ptr: number) {\n")
				g.writeIndent(s, 2)
				s.WriteString("this._ptr = ptr;\n")
				g.writeIndent(s, 2)
				s.WriteString("registry.register(this, { cb: this._gc_dispose, ptr });\n")
				g.writeIndent(s, 1)
				s.WriteString("}\n\n")

				g.writeIndent(s, 1)
				s.WriteString("get ptr(): number {\n")
				g.writeIndent(s, 2)
				s.WriteString("return this._ptr;\n")
				g.writeIndent(s, 1)
				s.WriteString("}\n\n")

				// write struct field `getters`
				for _, f := range c.fieldAccessors {
					g.writeIndent(s, 1)
					s.WriteString(fmt.Sprintf("get %s(): ", *f.name))
					tempType := g.getJSFromFFIType(f.returns[0].FFIType)
					if f.isHandleFn != nil && f.returns[0].FFIType == "FFIType.ptr" {
						s.WriteString(fmt.Sprintf("_%s | undefined", *f.isHandleFn))
					} else if f.isHandleFn != nil {
						s.WriteString(fmt.Sprintf(*f.isHandleFn))
					} else if *f.arrayType != "" {
						s.WriteString(fmt.Sprintf("%sArray | undefined", fmtr.String(*f.arrayType)))
					} else {
						s.WriteString(tempType)
						if f.isOptional {
							s.WriteString(" | undefined")
						}
					}
					s.WriteString(" {\n")
					g.writeIndent(s, 2)
					if f.structType != nil {
						s.WriteString(fmt.Sprintf("return <%s>%s(this._ptr)", *f.structType, *f.fnName))
						if f.returns[0].FFIType == "FFIType.cstring" {
							s.WriteString(".toString()")
						}
						s.WriteString(";\n")
					} else if f.isHandleFn != nil {
						s.WriteString(fmt.Sprintf("const ptr = %s(this._ptr);\n", *f.fnName))
						g.writeIndent(s, 2)
						s.WriteString("if (!ptr) return undefined;\n")
						g.writeIndent(s, 2)
						s.WriteString(fmt.Sprintf("return new _%s(ptr);\n", *f.isHandleFn))
					} else if *f.arrayType != "" {
						s.WriteString(fmt.Sprintf("const ptr = %s(this._ptr);\n", *f.fnName))
						g.writeIndent(s, 2)
						s.WriteString("if (!ptr) return undefined;\n")
						g.writeIndent(s, 2)
						s.WriteString("// eslint-disable-next-line @typescript-eslint/ban-ts-comment\n")
						g.writeIndent(s, 2)
						s.WriteString("// @ts-ignore - overload toArrayBuffer params\n")
						g.writeIndent(s, 2)
						s.WriteString(fmt.Sprintf("return new %sArray(toArrayBuffer(ptr, 0, arraySize(ptr) * %d, genDisposePtr.native()));\n", fmtr.String(*f.arrayType), getByteSize(*f.arrayType)))
					} else if f.returns[0].FFIType == "FFIType.cstring" {
						s.WriteString(fmt.Sprintf("return %s(this._ptr)", *f.fnName))
						if tempType == "string" {
							s.WriteString(".toString()")
						}
						s.WriteString(";\n")
					} else {
						s.WriteString(fmt.Sprintf("return %s(this._ptr);\n", *f.fnName))
					}
					g.writeIndent(s, 1)
					s.WriteString("}\n\n")
				}

				// write static method to init new Go Struct
				g.writeInitMethod(s, c, fmtr)

				// write class method that frees `Handle` + CGo mem for struct @ GC
				g.writeIndent(s, 1)
				s.WriteString("public _gc_dispose(ptr: number): void {\n")
				g.writeIndent(s, 2)
				s.WriteString("return _DISPOSE_Struct(ptr);\n")
				g.writeIndent(s, 1)
				s.WriteString("}\n")
				s.WriteString("}\n\n")
				struct_wrappers[*c.name] = true
			}
		}
	}
}

func (g *PackageGenerator) writeNestedFieldExports(s *strings.Builder, v *StructAccessor, struct_exports map[string]bool, class_wrappers *[]*ClassWrapper, visited int, count int, isLast bool) {
	if v.isHandleFn != nil && !struct_exports[*v.name] {
		var classWrapper = &ClassWrapper{
			name:           v.isHandleFn,
			fieldAccessors: v.fieldAccessors,
			disposeHandle:  v.disposeHandle,
			structType:     v.structType,
			args:           v.args,
			returns:        v.returns,
		}
		*class_wrappers = append(*class_wrappers, classWrapper)
		g.writeIndent(s, 2)
		s.WriteString(fmt.Sprintf("_INIT_%s,\n", *v.isHandleFn))
		fieldCount := len(classWrapper.fieldAccessors)
		fieldsVisited := 0
		for _, fa := range classWrapper.fieldAccessors {
			g.writeIndent(s, 2)
			s.WriteString(*fa.fnName)
			if isLast && visited == count-1 && fieldsVisited == fieldCount-1 {
				s.WriteByte('\n')
			} else {
				s.WriteString(",\n")
			}
			if fa.isHandleFn != nil {
				g.writeNestedFieldExports(s, fa, struct_exports, class_wrappers, fieldsVisited, fieldCount, isLast)
			}
			fieldsVisited++
		}
		struct_exports[*v.name] = true
	}
}

func (g *PackageGenerator) writeAccessorFieldExports(s *strings.Builder, v *FFIFunc, struct_exports map[string]bool, class_wrappers *[]*ClassWrapper, visited int, count int, isDisposeWritten *bool) {
	if v.isHandleFn && !struct_exports[*v.name] {
		var classWrapper = &ClassWrapper{
			name:           v.name,
			fieldAccessors: v.fieldAccessors,
			disposeHandle:  v.disposeHandle,
			args:           v.args,
			returns:        v.returns,
		}
		*class_wrappers = append(*class_wrappers, classWrapper)
		// declare export for struct dispose fn
		if !*isDisposeWritten {
			g.writeIndent(s, 2)
			s.WriteString("_DISPOSE_Struct,\n")
			*isDisposeWritten = true
		}
		g.writeIndent(s, 2)
		s.WriteString(fmt.Sprintf("_INIT_%s,\n", *v.name))
		fieldCount := len(classWrapper.fieldAccessors)
		fieldsVisited := 0
		for _, fa := range classWrapper.fieldAccessors {
			g.writeIndent(s, 2)
			s.WriteString(*fa.fnName)
			if visited == count-1 && fieldsVisited == fieldCount-1 {
				s.WriteByte('\n')
			} else {
				s.WriteString(",\n")
			}
			if fa.isHandleFn != nil {
				g.writeNestedFieldExports(s, fa, struct_exports, class_wrappers, fieldsVisited, fieldCount, visited == count-1)
			}
			fieldsVisited++
		}
		struct_exports[*v.name] = true
	}
}

func (g *PackageGenerator) writeNestedFieldConfig(s *strings.Builder, v *StructAccessor, struct_config map[string]bool, k string, visited int, count int, resLen int, isLast bool) {
	if v.isHandleFn != nil && !struct_config[*v.isHandleFn] {
		g.writeIndent(s, 1)
		s.WriteString(fmt.Sprintf("_INIT_%s: {\n", *v.isHandleFn))
		g.writeIndent(s, 2)
		s.WriteString("args: [")
		argLen := len(v.fieldAccessors)
		for i, fa := range v.fieldAccessors {
			if fa.isHandleFn != nil {
				s.WriteString("FFIType.u64_fast")
			} else {
				s.WriteString(fa.args[0].FFIType)
				if fa.arrayType != nil {
					s.WriteString(", FFIType.u64_fast")
				}
			}
			if i < argLen-1 {
				s.WriteString(", ")
			}
		}
		s.WriteString("],\n")
		g.writeIndent(s, 2)
		s.WriteString("returns: FFIType.ptr\n")
		g.writeIndent(s, 1)
		s.WriteString("},\n")
		// write config for struct field accessors
		fieldCount := len(v.fieldAccessors)
		fieldsVisited := 0
		for _, fa := range v.fieldAccessors {
			g.writeIndent(s, 1)
			s.WriteString(fmt.Sprintf("%s: {\n", *fa.fnName))
			fieldArgLen := len(fa.args)
			fieldResLen := len(fa.returns)
			if fieldArgLen > 0 {
				g.writeIndent(s, 2)
				s.WriteString("args: [")
				for i, arg := range fa.args {
					s.WriteString(arg.FFIType)
					if i < fieldArgLen-1 {
						s.WriteString(", ")
					}
				}
				s.WriteByte(']')
				if fieldResLen > 0 {
					s.WriteString(",\n")
				} else {
					s.WriteByte('\n')
				}
			}
			if resLen == 1 {
				g.writeIndent(s, 2)
				s.WriteString(fmt.Sprintf("returns: %s\n", fa.returns[0].FFIType))
			} else if resLen > 1 {
				var errStr strings.Builder
				errStr.WriteString("Function `")
				errStr.WriteString(k)
				errStr.WriteString("` has more than one return value, which is not supported by Bun FFI...\n")
				errStr.WriteString("Consider adjusting your `tsgo.yaml` config file to inject code before/after the function call in the CGo wrapper fn as you can coerce to a single return value in Go.\n")
				log.Fatalf("TSGo failed: %v", errStr.String())
			}

			g.writeIndent(s, 1)
			if isLast && visited == count-1 && fieldsVisited == fieldCount-1 {
				s.WriteString("}\n")
			} else {
				s.WriteString("},\n")
			}
			if fa.isHandleFn != nil {
				g.writeNestedFieldConfig(s, fa, struct_config, k, fieldsVisited, fieldCount, fieldResLen, isLast)
			}
			fieldsVisited++
		}
		struct_config[*v.name] = true
	}
}

func (g *PackageGenerator) writeAccessorFieldConfig(s *strings.Builder, v *FFIFunc, struct_config map[string]bool, k string, visited int, count int, resLen int, isDisposeWritten *bool) {
	if v.isHandleFn && !struct_config[*v.name] {
		// write config for struct dispose fn
		if !*isDisposeWritten {
			g.writeIndent(s, 1)
			s.WriteString("_DISPOSE_Struct: {\n")
			g.writeIndent(s, 2)
			s.WriteString(fmt.Sprintf("args: [%s]\n", v.disposeHandle.args[0].FFIType))
			g.writeIndent(s, 1)
			s.WriteString("},\n")
			*isDisposeWritten = true
		}
		g.writeIndent(s, 1)
		s.WriteString(fmt.Sprintf("_INIT_%s: {\n", *v.name))
		g.writeIndent(s, 2)
		s.WriteString("args: [")
		argLen := len(v.fieldAccessors)
		for i, fa := range v.fieldAccessors {
			s.WriteString(fa.args[0].FFIType)
			if fa.arrayType != nil {
				s.WriteString(", FFIType.u64_fast")
			}
			if i < argLen-1 {
				s.WriteString(", ")
			}
		}
		s.WriteString("],\n")
		g.writeIndent(s, 2)
		s.WriteString("returns: FFIType.ptr\n")
		g.writeIndent(s, 1)
		s.WriteString("},\n")

		// write config for struct field accessors
		fieldCount := len(v.fieldAccessors)
		fieldsVisited := 0
		for _, fa := range v.fieldAccessors {
			g.writeIndent(s, 1)
			s.WriteString(fmt.Sprintf("%s: {\n", *fa.fnName))
			fieldArgLen := len(fa.args)
			fieldResLen := len(fa.returns)
			if fieldArgLen > 0 {
				g.writeIndent(s, 2)
				s.WriteString("args: [")
				for i, arg := range fa.args {
					s.WriteString(arg.FFIType)
					if i < fieldArgLen-1 {
						s.WriteString(", ")
					}
				}
				s.WriteByte(']')
				if fieldResLen > 0 {
					s.WriteString(",\n")
				} else {
					s.WriteByte('\n')
				}
			}
			if resLen == 1 {
				g.writeIndent(s, 2)
				s.WriteString(fmt.Sprintf("returns: %s\n", fa.returns[0].FFIType))
			} else if resLen > 1 {
				var errStr strings.Builder
				errStr.WriteString("Function `")
				errStr.WriteString(k)
				errStr.WriteString("` has more than one return value, which is not supported by Bun FFI...\n")
				errStr.WriteString("Consider adjusting your `tsgo.yaml` config file to inject code before/after the function call in the CGo wrapper fn as you can coerce to a single return value in Go.\n")
				log.Fatalf("TSGo failed: %v", errStr.String())
			}

			g.writeIndent(s, 1)
			if visited == count-1 && fieldsVisited == fieldCount-1 {
				s.WriteString("}\n")
			} else {
				s.WriteString("},\n")
			}
			if fa.isHandleFn != nil {
				g.writeNestedFieldConfig(s, fa, struct_config, k, fieldsVisited, fieldCount, fieldResLen, visited == count-1)
			}
			fieldsVisited++
		}
		struct_config[*v.name] = true
	}
}

func (g *PackageGenerator) writeFFIConfig(s *strings.Builder, fd []*ast.FuncDecl, path string) {
	var class_wrappers = []*ClassWrapper{}
	caser := cases.Title(language.AmericanEnglish)

	s.WriteString("\n//////////\n")
	// source: misc.go
	s.WriteString("// Generated config for Bun FFI\n\n")
	s.WriteString("export const {\n")
	g.writeIndent(s, 1)
	s.WriteString("symbols: {\n")

	count := len(g.ffi.FFIFuncs)
	visited := 0
	struct_exports := map[string]bool{}
	disposeWritten := false
	for k, v := range g.ffi.FFIFuncs {
		g.writeIndent(s, 2)
		if !g.ffi.FFIHelpers[k] {
			s.WriteByte('_')
		}
		s.WriteString(k)
		if visited == count-1 && !v.isHandleFn {
			s.WriteByte('\n')
		} else {
			s.WriteString(",\n")
		}
		g.writeAccessorFieldExports(s, v, struct_exports, &class_wrappers, visited, count, &disposeWritten)
		visited++
	}

	g.writeIndent(s, 1)
	s.WriteString(fmt.Sprintf("}\n} = dlopen(import.meta.dir + '/%s/gen_bindings.dylib', {\n", path))
	visited = 0
	struct_config := map[string]bool{}
	isDisposeWritten := false
	for k, v := range g.ffi.FFIFuncs {
		g.writeIndent(s, 1)
		if !g.ffi.FFIHelpers[k] {
			s.WriteByte('_')
		}
		s.WriteString(fmt.Sprintf("%s: {\n", k))
		argLen := len(v.args)
		resLen := len(v.returns)
		if argLen > 0 {
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
			if v.isHandleFn {
				s.WriteString("FFIType.ptr")
			} else {
				s.WriteString(v.returns[0].FFIType)
			}
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
		if visited == count-1 && !v.isHandleFn {
			s.WriteString("}\n")
		} else {
			s.WriteString("},\n")
		}
		g.writeAccessorFieldConfig(s, v, struct_config, k, visited, count, resLen, &isDisposeWritten)
		visited++
	}
	s.WriteString("})\n\n")

	g.writeAccessorClasses(s, &class_wrappers, caser)
}
