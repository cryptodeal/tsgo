package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

// using this as a simplified testing ground for quick local dev
// test using `go run fl_test/parse_cpp.go` from root dir

type CPPArgDefault struct {
	Val *string
}

type TemplateArg struct {
	Identifier *string
}

type TemplateDeclArg struct {
	Identifier *string
	MetaType   *string
}

type TemplateDecl struct {
	Args *[]*TemplateDeclArg
}

type QualifiedIdentifier struct {
	Scope        *string
	Name         *string
	TemplateArgs *[]*TemplateArg
}

type CPPFriendFunc struct {
	QualifiedIdent *QualifiedIdentifier
	Args           *[]*CPPArg
}

type CPPFriend struct {
	Ident          *string
	QualifiedIdent *QualifiedIdentifier
	IsClass        bool
	Type           *CPPType
	FuncDecl       *CPPFriendFunc
}

type CPPType struct {
	FullType     *string
	Scope        *string
	Name         *string
	TemplateType *[]*TemplateArg
}

type CPPClass struct {
	Name         *string
	FieldDecl    *[]*CPPArg
	FriendDecl   *[]*CPPFriend
	Decl         *[]*CPPArg
	TemplateDecl *[]*TemplateMethod
}

type CPPArg struct {
	TypeQualifier *string
	Type          *string
	RefDecl       *string
	Ident         *string
	DefaultValue  *CPPArgDefault
}

type CPPMethod struct {
	Ident     *string
	Overloads []*[]*CPPArg
	Returns   *string
}

type TemplateMethod struct {
	TemplateDecl  *TemplateDecl
	Returns       *string
	PointerMethod bool
	Ident         *string
	Args          *[]*CPPArg
	TypeQualifier *string
}

type ParsedMethod struct {
	Ident   *string
	Args    *[]*CPPArg
	Returns *string
}

func (a *CPPArg) IsPtr() bool {
	return a.RefDecl != nil
}

func parseHeader(path string) map[string]*CPPMethod {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		err := errors.New("unable to get the path to cwd")
		panic(err)
	}

	fPath := filepath.Join(filepath.Dir(filename), "../", path)
	input, err := os.ReadFile(fPath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(cpp.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, input)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	n := tree.RootNode()
	methods := parseMethods(n, input)
	classFriends := parseClasses(n, input)
	for _, cf := range classFriends {
		fmt.Println(cf)
	}
	return methods

}

func parseClassFriend(n *sitter.Node, input []byte) *CPPFriend {
	new_friend := &CPPFriend{}
	child_count := int(n.ChildCount())
	for j := 0; j < child_count; j++ {
		grandchild := n.Child(j)
		tempType := grandchild.Type()
		switch tempType {
		case "type_identifier":
			new_friend.Ident = &tempType
			new_friend.IsClass = true
		case "declaration":
			great_grandchild_count := int(grandchild.ChildCount())
			for k := 0; k < great_grandchild_count; k++ {
				great_grandchild := grandchild.Child(k)
				temp_great_type := great_grandchild.Type()
				/* nested switch, is a bit ugly, but good perf */
				switch temp_great_type {
				case "qualified_identifier":
					qualID := &QualifiedIdentifier{}
					scope := great_grandchild.ChildByFieldName("scope")
					if scope != nil {
						tempScope := scope.Content(input)
						qualID.Scope = &tempScope
					}
					_name := great_grandchild.ChildByFieldName("name")
					if _name != nil {
						name := _name.ChildByFieldName("name")
						if name != nil {
							tempName := name.Content(input)
							qualID.Name = &tempName
						}
						arguments := _name.ChildByFieldName("arguments")
						arg_childs := int(arguments.ChildCount())
						template_args := &[]*TemplateArg{}
						for l := 0; l < arg_childs; l++ {
							arg := arguments.Child(l)
							argType := arg.Type()
							if argType == "type_descriptor" {
								temp_arg_type := arg.ChildByFieldName("type")
								if temp_arg_type != nil {
									parsed_temp_arg := temp_arg_type.Content(input)
									*template_args = append(*template_args, &TemplateArg{&parsed_temp_arg})
								}
							}
						}
						qualID.TemplateArgs = template_args
					}
					new_friend.QualifiedIdent = qualID
				case "function_declarator":
					decl := great_grandchild.ChildByFieldName("declarator")
					friend_func := &CPPFriendFunc{QualifiedIdent: &QualifiedIdentifier{}}
					if decl != nil {
						scope := decl.ChildByFieldName("scope")
						if scope != nil {
							tempScope := scope.Content(input)
							friend_func.QualifiedIdent.Scope = &tempScope
						}
						name := decl.ChildByFieldName("name")
						if name != nil {
							tempName := name.Content(input)
							friend_func.QualifiedIdent.Name = &tempName
						}
					}
					params := great_grandchild.ChildByFieldName("parameters")
					friend_func.Args = parseCPPArg(input, params)
					new_friend.FuncDecl = friend_func
				} /* end nested switch/case */
			} /* end outer switch/case */
		}
	}
	return new_friend
}

func parseClassTemplateMethod(n *sitter.Node, input []byte) *TemplateMethod {
	template_method := &TemplateMethod{
		TemplateDecl: &TemplateDecl{},
	}
	params := n.ChildByFieldName("parameters")

	if params != nil {
		param_count := int(params.ChildCount())
		template_method.TemplateDecl.Args = &[]*TemplateDeclArg{}
		for i := 0; i < param_count; i++ {
			param := params.Child(i)
			paramType := param.Type()
			if paramType == "type_parameter_declaration" {
				param_split := strings.Split(param.Content(input), " ")
				DeclArg := &TemplateDeclArg{
					Identifier: &param_split[1],
					MetaType:   &param_split[0],
				}
				*template_method.TemplateDecl.Args = append(*template_method.TemplateDecl.Args, DeclArg)
			}
		}
	}

	childCount := int(n.ChildCount())
	for i := 0; i < childCount; i++ {
		childType := n.Child(i).Type()
		if childType == "declaration" {
			tempChild := n.Child(i)
			tempType := tempChild.ChildByFieldName("type")
			if tempType != nil {
				content := tempType.Content(input)
				template_method.Returns = &content
			}
			declarator := tempChild.ChildByFieldName("declarator")
			if declarator != nil {
				declType := declarator.Type()
				if declType == "pointer_declarator" {
					template_method.PointerMethod = true
					decl := declarator.ChildByFieldName("declarator")
					if decl != nil {
						declType := decl.Type()
						if declType == "function_declarator" {
							template_method.Args = parseCPPArg(input, decl.ChildByFieldName("parameters"))
							nameNode := decl.ChildByFieldName("name")
							if nameNode != nil {
								name := nameNode.Content(input)
								template_method.Ident = &name
							} else {
								fmt.Printf("unhandled node:\nnode view: %s\nnode type: %s\ncontent: %s\n", decl, declType, decl.Content(input))
							}
							decl_child_count := int(decl.ChildCount())
							for j := 0; j < decl_child_count; j++ {
								decl_child := decl.Child(j)
								decl_child_type := decl_child.Type()
								if decl_child_type == "type_qualifier" {
									type_qualifier := decl_child.Content(input)
									template_method.TypeQualifier = &type_qualifier
								}
							}
						}
					}
				} else if declType == "function_declarator" {
					decl := declarator.ChildByFieldName("declarator")
					if decl != nil {
						template_method.Args = parseCPPArg(input, decl.ChildByFieldName("parameters"))
						nameNode := decl.ChildByFieldName("name")
						if nameNode != nil {
							name := nameNode.Content(input)
							template_method.Ident = &name
						}
						decl_child_count := int(decl.ChildCount())
						for j := 0; j < decl_child_count; j++ {
							decl_child := decl.Child(j)
							decl_child_type := decl_child.Type()
							if decl_child_type == "type_qualifier" {
								type_qualifier := decl_child.Content(input)
								template_method.TypeQualifier = &type_qualifier
							} else {
								fmt.Printf("unhandled node:\nnode view: %s\nnode type: %s\ncontent: %s\n", decl, declType, decl.Content(input))
							}
						}
					}
				}
			}
		}
	}
	return template_method
}

func parseClasses(n *sitter.Node, input []byte) map[string]*CPPClass {
	classes := map[string]*CPPClass{}

	q, err := sitter.NewQuery([]byte("(class_specifier) @class_def"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// fmt.Println(len(m.Captures))
		for _, c := range m.Captures {
			class_name := c.Node.ChildByFieldName("name").Content(input)
			classes[class_name] = &CPPClass{}
			class_body := c.Node.ChildByFieldName("body")
			class_friends := &[]*CPPFriend{}
			if class_body == nil {
				// TODO: should probably parse class def w/o body as well
				continue
			}

			child_count := int(class_body.ChildCount())
			for i := 0; i < child_count; i++ {
				temp_child := class_body.Child(i)
				if temp_child.Type() == "friend_declaration" { // WORKING
					new_friend := parseClassFriend(temp_child, input)
					*class_friends = append(*class_friends, new_friend)
				} else if temp_child.Type() == "field_declaration" { // TODO: parse class `field_declaration`
					// fmt.Println(temp_child.Content(input))
				} else if temp_child.Type() == "declaration" { // TODO: parse class `declaration`
					// fmt.Println(temp_child.Content(input))
				} else if temp_child.Type() == "template_declaration" { // TODO: parse class `template_declaration`
					if classes[class_name].TemplateDecl == nil {
						classes[class_name].TemplateDecl = &[]*TemplateMethod{}
					}
					temp_decl := parseClassTemplateMethod(temp_child, input)
					*classes[class_name].TemplateDecl = append(*classes[class_name].TemplateDecl, temp_decl)
				}
			}
			classes[class_name].FriendDecl = class_friends
		}
	}
	return classes
}

func parseMethods(n *sitter.Node, input []byte) map[string]*CPPMethod {
	methods := map[string]*CPPMethod{}
	q, err := sitter.NewQuery([]byte("(declaration type: (type_identifier) @type declarator: (function_declarator) @func)"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		res, body := splitMatches(m.Captures)
		parsed := parseCPPMethod(input, res.Node, body.Node)
		if v, ok := methods[*parsed.Ident]; ok {
			// encountered method previously (fn overloading)
			v.Overloads = append(v.Overloads, parsed.Args)
		} else {
			// first time having encountered this method, so create a new entry
			new_method := &CPPMethod{
				Ident:     parsed.Ident,
				Overloads: []*[]*CPPArg{parsed.Args},
				Returns:   parsed.Returns,
			}
			methods[*parsed.Ident] = new_method
		}
	}
	return methods
}

func parseCPPMethod(content []byte, r *sitter.Node, b *sitter.Node) *ParsedMethod {
	args := parseCPPArg(content, b.ChildByFieldName("parameters"))
	name := b.ChildByFieldName("declarator").Content(content)
	var returns *string
	if r != nil {
		tempReturns := r.Content(content)
		returns = &tempReturns
	}
	parsed := &ParsedMethod{
		Args:    args,
		Ident:   &name,
		Returns: returns,
	}
	return parsed
}

func parseCPPArg(content []byte, arg_list *sitter.Node) *[]*CPPArg {
	args := []*CPPArg{}
	if arg_list == nil {
		return &args
	}
	child_count := int(arg_list.ChildCount())

	for i := 0; i < child_count; i++ {
		scoped_arg := arg_list.Child(i)
		node_type := scoped_arg.Type()
		if node_type != "parameter_declaration" && node_type != "optional_parameter_declaration" {
			continue
		}

		var RefDecl string
		var Identifier string
		refDecl := scoped_arg.ChildByFieldName("declarator")
		refType := refDecl.Type()
		count := int(refDecl.ChildCount())
		if refType == "reference_declarator" {
			RefDecl = refDecl.Content(content)
			for j := 0; j < count; j++ {
				child := refDecl.Child(j)
				if child.Type() != "identifier" {
					continue
				}
				Identifier = child.Content(content)
			}
		} else if refType == "identifier" {
			Identifier = refDecl.Content(content)
		}

		var TypeQualifier string
		count = int(scoped_arg.ChildCount())
		for j := 0; j < count; j++ {
			child := scoped_arg.Child(j)
			if child.Type() != "type_qualifier" {
				continue
			}
			TypeQualifier = child.Content(content)
		}

		var DefaultValue CPPArgDefault
		if node_type == "optional_parameter_declaration" {
			tempDefault := scoped_arg.ChildByFieldName("default_value").Content(content)
			DefaultValue = CPPArgDefault{Val: &tempDefault}
		}

		argType := scoped_arg.ChildByFieldName("type").Content(content)
		temp_arg := &CPPArg{
			Ident:         &Identifier,
			Type:          &argType,
			TypeQualifier: &TypeQualifier,
			DefaultValue:  &DefaultValue,
		}

		if RefDecl != "" {
			temp_arg.RefDecl = &RefDecl
		}

		args = append(args, temp_arg)
	}
	return &args
}

func splitMatches(matched []sitter.QueryCapture) (sitter.QueryCapture, sitter.QueryCapture) {
	return matched[0], matched[1]
}

func main() {
	parseHeader("fl_test/TensorBase.h")
}
