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

// Responsible for generating the code for an input package
type PackageGenerator struct {
	conf    *PackageConfig
	pkg     *packages.Package
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

		pkgGen := &PackageGenerator{
			conf:    pkgConfig,
			GoFiles: pkg.GoFiles,
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
			// used to build the bindings
			// TODO: clean up logic
			dir := filepath.Dir(outPath)
			var bindings_out strings.Builder
			var build_out strings.Builder

			pkg_split := strings.Split(filepath.Dir(pkg.GoFiles[0]), "/")
			pkg_name := pkg_split[len(pkg_split)-1]
			bindings_out.WriteString(pkg_name)
			bindings_out.WriteString("/gen_bindings.dylib")
			build_out.WriteString(pkg_name)

			build_out.WriteString("/gen_bindings.go")

			bindings_out_path := filepath.Join(dir, bindings_out.String())
			build_out_path := filepath.Join(dir, build_out.String())
			fmt.Println("bindings_out_path:", bindings_out_path)
			fmt.Println("build_out_path:", build_out_path)

			// TODO: write CGo wrappers to build_out_path

			// builds command string to execute (used to compile bindings)
			var cmd_str strings.Builder
			cmd_str.WriteString("go build --buildmode c-shared -o ")
			cmd_str.WriteString(bindings_out_path)
			cmd_str.WriteByte(' ')
			cmd_str.WriteString(build_out_path)

			cmd := exec.Command(cmd_str.String())
			err := cmd.Run()

			if err != nil {
				log.Fatal(err)
			}

		}
	}
	return nil
}
