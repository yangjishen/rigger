package rigger

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"

	pkgErr "github.com/pkg/errors"
	"html/template"
)

const GoCorePath = "/src/rigger"

func init() {
	Gopath = os.Getenv("GOPATH")
	if Gopath == "" {
		Gopath = build.Default.GOPATH
	}
}

var Gopath string

type core struct {
	debug bool
}

type data struct {
	AbsGenProjectPath string // the abs gen project path
	ProjectPath       string // the go import project path
	ProjectName       string // the project name which want to generated
	Quit              string
}

type templateSet struct {
	templateFilePath string
	templateFileName string
	genFilePath      string
}

type templateEngine struct {
	Templates []templateSet
	currDir   string
}

func New(debug bool) *core {
	return &core{
		debug: debug,
	}
}

func (c *core) Generate(path string) error {
	genAbsDir, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	projectName := filepath.Base(genAbsDir)
	goProjectPath := strings.TrimPrefix(genAbsDir, filepath.Join(Gopath, "src")+string(os.PathSeparator))
	d := data{
		AbsGenProjectPath: genAbsDir,
		ProjectPath:       goProjectPath,
		ProjectName:       projectName,
		Quit:              "quit",
	}
	if err := c.genFromTemplate(getTemplateSets(), d); err != nil {
		return err
	}

	if err := c.genFormStaticFle(d); err != nil {
		return err
	}
	return nil

}

func getTemplateSets() []templateSet {
	tt := templateEngine{}
	templatesFolder := filepath.Join(Gopath, GoCorePath, "template/")
	filepath.Walk(templatesFolder, tt.visit)
	return tt.Templates
}

func (c *core) genFromTemplate(templateSets []templateSet, d data) error {
	for _, tmpl := range templateSets {
		if err := c.tmplExec(tmpl, d); err != nil {
			return err
		}
	}
	return nil
}

func unescaped(x string) interface{} { return template.HTML(x) }

func (c *core) tmplExec(tmplSet templateSet, d data) error {
	tmpl := template.New(tmplSet.templateFileName)
	tmpl = tmpl.Funcs(template.FuncMap{"unescaped": unescaped})
	tmpl, err := tmpl.ParseFiles(tmplSet.templateFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}

	relateDir := filepath.Dir(tmplSet.genFilePath)

	distRelFilePath := filepath.Join(relateDir, filepath.Base(tmplSet.genFilePath))
	distAbsFilePath := filepath.Join(d.AbsGenProjectPath, distRelFilePath)

	c.debugPrintf("distRelFilePath:%s\n", distAbsFilePath)
	c.debugPrintf("distAbsFilePath:%s\n", distAbsFilePath)

	if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
		return pkgErr.WithStack(err)
	}

	dist, err := os.Create(distAbsFilePath)
	if err != nil {
		return pkgErr.WithStack(err)
	}
	defer func() {
		if err = dist.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Printf("Create %s \n", distRelFilePath)
	return tmpl.Execute(dist, d)

}

func (c *core) debugPrintf(format string, a ...interface{}) {
	if c.debug == true {
		fmt.Printf(format, a...)
	}
}

func (c *core) genFormStaticFle(d data) error {
	walkerFuc := func(path string, f os.FileInfo, err error) error {
		if f.Mode().IsRegular() == true {
			src, err := os.Open(path)
			if err != nil {
				return pkgErr.WithStack(err)
			}

			defer func() {
				if err := src.Close(); err != nil {
					panic(err)
				}
			}()

			basepath := filepath.Join(Gopath, GoCorePath, "static")
			distRelFilePath, err := filepath.Rel(basepath, path)
			if err != nil {
				return pkgErr.WithStack(err)
			}

			distAbsFilePath := filepath.Join(d.AbsGenProjectPath, distRelFilePath)

			if err := os.MkdirAll(filepath.Dir(distAbsFilePath), os.ModePerm); err != nil {
				return pkgErr.WithStack(err)
			}

			dist, err := os.Create(distAbsFilePath)
			if err != nil {
				return pkgErr.WithStack(err)
			}
			defer func() {
				if err := dist.Close(); err != nil {
					panic(err)
				}
			}()

			if _, err := io.Copy(dist, src); err != nil {
				return pkgErr.WithStack(err)
			}
			fmt.Printf("Create %s \n", distRelFilePath)

		}
		return nil
	}
	walkPath := filepath.Join(Gopath, GoCorePath, "static")
	return filepath.Walk(walkPath, walkerFuc)
}

func (templEngine *templateEngine) visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if ext := filepath.Ext(path); ext == ".tmpl" {
		templateFileName := filepath.Base(path)

		genFileBaseName := strings.TrimSuffix(templateFileName, ".tmpl") + ".go"
		genFileBasePath, err := filepath.Rel(filepath.Join(Gopath, GoCorePath, "template"), filepath.Join(filepath.Dir(path), genFileBaseName))
		if err != nil {
			return pkgErr.WithStack(err)
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(templEngine.currDir, genFileBasePath),
		}

		templEngine.Templates = append(templEngine.Templates, templ)

	} else if mode := f.Mode(); mode.IsRegular() {
		templateFileName := filepath.Base(path)

		basepath := filepath.Join(Gopath, GoCorePath, "template")
		targpath := filepath.Join(filepath.Dir(path), templateFileName)
		genFileBasePath, err := filepath.Rel(basepath, targpath)
		if err != nil {
			return pkgErr.WithStack(err)
		}

		templ := templateSet{
			templateFilePath: path,
			templateFileName: templateFileName,
			genFilePath:      filepath.Join(templEngine.currDir, genFileBasePath),
		}

		templEngine.Templates = append(templEngine.Templates, templ)
	}
	return nil
}
