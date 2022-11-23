package tsgo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Generator for one or more input packages, responsible for linking
// them together if necessary.
type TSGo struct {
	conf *Config

	packageGenerators map[string]*PackageGenerator
}

type FFIFunc struct {
	wrap_args      []string
	wrap_returns   []string
	args           []string
	returns        []string
	isHandleFn     bool
	name           *string
	fieldAccessors []*FFIFunc
}

type FFIState struct {
	GoImports     map[string]bool
	CImports      map[string]bool
	FFIHelpers    map[string]bool
	CHelpers      map[string]bool
	FFIFuncs      map[string]*FFIFunc
	StructHelpers map[string][]*FFIFunc
}

// Responsible for generating the code for an input package
type PackageGenerator struct {
	conf    *PackageConfig
	pkg     *packages.Package
	ffi     *FFIState
	GoFiles []string
}

func New(config *Config) *TSGo {
	return &TSGo{
		conf:              config,
		packageGenerators: make(map[string]*PackageGenerator),
	}
}

func (g *TSGo) SetTypeMapping(goType string, tsType string) {
	for _, p := range g.conf.Packages {
		p.TypeMappings[goType] = tsType
	}
}

func (g *TSGo) Generate() error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles,
	}, g.conf.PackageNames()...)
	if err != nil {
		return err
	}

	for i, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return fmt.Errorf("%+v", pkg.Errors)
		}

		if len(pkg.GoFiles) == 0 {
			return fmt.Errorf("no input go files for package index %d", i)
		}

		pkgConfig := g.conf.PackageConfig(pkg.ID)

		ffi := &FFIState{
			GoImports:     make(map[string]bool),
			CImports:      make(map[string]bool),
			FFIHelpers:    make(map[string]bool),
			CHelpers:      make(map[string]bool),
			FFIFuncs:      make(map[string]*FFIFunc),
			StructHelpers: make(map[string][]*FFIFunc),
		}

		pkgGen := &PackageGenerator{
			conf:    pkgConfig,
			GoFiles: pkg.GoFiles,
			ffi:     ffi,
			pkg:     pkg,
		}
		g.packageGenerators[pkg.PkgPath] = pkgGen
		code, err := pkgGen.Generate()
		if err != nil {
			return err
		}

		outPath := pkgGen.conf.ResolvedOutputPath(filepath.Dir(pkg.GoFiles[0]))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = ioutil.WriteFile(outPath, []byte(code), os.ModePerm)
		if err != nil {
			return nil
		}
		if pkgGen.conf.FFIBindings {
			dir := filepath.Dir(outPath)
			pkg_split := strings.Split(filepath.Dir(pkg.GoFiles[0]), "/")
			pkg_name := pkg_split[len(pkg_split)-1]

			var bindings_out strings.Builder
			bindings_out.WriteString(pkg_name)
			bindings_out.WriteString("/gen_bindings.dylib")
			bindings_out_path := filepath.Join(dir, bindings_out.String())

			var build_out strings.Builder
			build_out.WriteString(pkg_name)
			build_out_path := filepath.Join(dir, build_out.String())

			// builds command string to execute (used to compile bindings)
			cmd_str := []string{"build", "--buildmode", "c-shared", "-o", bindings_out_path, build_out_path}
			cmd := exec.Command("go", cmd_str...)
			err := cmd.Run()

			if err != nil {
				log.Fatal(err)
			}

		}
	}
	return nil
}
