package tsgo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

type CPPArgDefault struct {
	Val *string
}

type CPPArg struct {
	TypeQualifier *string
	Type          *string
	RefDecl       *string
	Identifier    *string
	DefaultValue  *CPPArgDefault
}

type CPPMethod struct {
	Identifier *string
	Overloads  []*[]*CPPArg
	Returns    *string
}

type ParsedMethod struct {
	Identifier *string
	Args       *[]*CPPArg
	Returns    *string
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
	parseClasses(n, input)
	return methods

}

func parseClasses(n *sitter.Node, input []byte) {
	//classes := map[string]*CPPMethod{}
	q, err := sitter.NewQuery([]byte("(class_specifier name: (type_identifier) @class_name)"), cpp.GetLanguage())
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
		for _, c := range m.Captures {
			fmt.Println(c.Node.Content(input))
		}
	}
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
		if v, ok := methods[*parsed.Identifier]; ok {
			// encountered method previously (fn overloading)
			v.Overloads = append(v.Overloads, parsed.Args)
		} else {
			// first time having encountered this method, so create a new entry
			new_method := &CPPMethod{
				Identifier: parsed.Identifier,
				Overloads:  []*[]*CPPArg{parsed.Args},
				Returns:    parsed.Returns,
			}
			methods[*parsed.Identifier] = new_method
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
		Args:       args,
		Identifier: &name,
		Returns:    returns,
	}
	return parsed
}

func parseCPPArg(content []byte, arg_list *sitter.Node) *[]*CPPArg {
	args := []*CPPArg{}
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
			Identifier:    &Identifier,
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
