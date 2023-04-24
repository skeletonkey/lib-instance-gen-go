// Package app is a library for creating a bare-bones application with boilerplate files taken
// care of. The long term goal is to be the basis of all Go applications allowing for quick propagation of updates,
// bug fixes, and new features.
//
//	app := instance_gen.NewApp("rachio-next-run", "app")
//	app.WithPackages("logger", "pushover", "rachio").
//		WithGithubWorkflows("linter", "test").
//		WithMakefile()
//
// Each generated file will be prepended with a 'warning' comment to not edit the file.
package app

import (
	"embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

//go:embed all:templates
var templatesFS embed.FS

const templateBaseDir = "templates"
const mkfilesSubDir = "Makefile"
const warning = "lib-instance-gen-go: File auto generated -- DO NOT EDIT!!!\n"

var templateExts = map[string]string{
	"go":     ".go.tpl",
	"mkfile": ".tpl",
	"yml":    ".yml.tpl",
}
var warnings = map[string]string{
	"go":     "// " + warning,
	"mkfile": "// " + warning,
	"yml":    "# " + warning,
}

// App struct containing necessary information for a new application
type App struct {
	binaryName string // name of the binary the 'make' will produce
	dir        string // subdirectory which will contain the program's source code
}

// NewApp returns the struct for a new applications which allows for generating boilerplate files.
//   - binaryName is used by the Makefile for the build command
//   - dir is the subdirectory that packages will be created in
func NewApp(binaryName string, dir string) App {
	return App{binaryName: binaryName, dir: dir}
}

// WithPackages takes a list of strings which results in creating a skeleton subdirectory for each.
// This sets up
//   - config.go - integration with github.com/skeletonkey/rachio-next-run/app/config
func (a App) WithPackages(packageNames ...string) App {
	for _, name := range packageNames {
		templateArgs := templateArgs{
			PackageName: name,
		}
		generateTemplate(generateTemplateArgs{
			fileType:       "go",
			outputName:     "config.go",
			outputSubDir:   path.Join(a.dir, name),
			templateName:   "config" + templateExts["go"],
			templateSubDir: "package",
			templateArgs:   templateArgs,
		})
	}
	return a
}

// WithGithubWorkflows sets up the specified workflows.
// Current supported work flows:
//   - linter - on pull request for all branches
//   - test - on pull request for all branches
func (a App) WithGithubWorkflows(flows ...string) App {
	for _, name := range flows {
		generateTemplate(generateTemplateArgs{
			fileType:       "yml",
			outputName:     name + ".yml",
			outputSubDir:   path.Join(".github", "workflows"),
			templateName:   name + templateExts["yml"],
			templateSubDir: "github_workflows",
		})
	}
	return a
}

// WithMakefile creates the basic Makefile with:
//   - build - runs fmt, vet and then builds the binary
//   - install - move binary to /usr/local/bin
//   - golib-latest - install go dependencies
//   - app-init - generate the boilerplate
func (a App) WithMakefile() App {
	nodes, err := templatesFS.ReadDir(path.Join(templateBaseDir, mkfilesSubDir))
	if err != nil {
		panic(fmt.Errorf("unable to read dir (%s): %s", mkfilesSubDir, err))
	}
	for _, node := range nodes {
		generateTemplate(generateTemplateArgs{
			fileType:       "mkfile",
			outputName:     node.Name()[:len(node.Name())-len(templateExts["mkfile"])],
			outputSubDir:   "",
			templateName:   node.Name(),
			templateSubDir: mkfilesSubDir,
			templateArgs:   templateArgs{BinaryName: a.binaryName},
		})
	}
	return a
}

type generateTemplateArgs struct {
	fileType       string       // type of file that the template is for the correct warning message
	outputName     string       // name of the template in its final form
	outputSubDir   string       // sub dir added to root dir for the final file
	templateName   string       // name of the template file
	templateSubDir string       // sub dir added to the template base dir to find template
	templateArgs   templateArgs // args that are fed to text/template
}
type templateArgs struct {
	BinaryName  string // name of the executable program
	PackageName string // name of the package
}

func generateTemplate(args generateTemplateArgs) {
	inputFileName := path.Join(templateBaseDir, args.templateSubDir, args.templateName)
	outputFileName := path.Join(args.outputSubDir, args.outputName)

	if args.outputSubDir != "" {
		err := os.MkdirAll(args.outputSubDir, 0755)
		if err != nil {
			panic(fmt.Errorf("unable to create directory structure (%s): %s", args.outputSubDir, err))
		}
	}
	f, err := os.Create(outputFileName)
	if err != nil {
		panic(fmt.Errorf("unable to create file (%s): %s", outputFileName, err))
	}
	_, err = f.WriteString(warnings[args.fileType])
	if err != nil {
		panic(fmt.Errorf("unable to write warning to file (%s): %s", outputFileName, err))
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(fmt.Errorf("error closing file (%s): %s", outputFileName, err))
		}
	}()
	if err != nil {
		panic(fmt.Errorf("unable to create new file (%s): %s", outputFileName, err))
	}
	temp := template.Must(template.ParseFS(templatesFS, inputFileName))
	err = temp.Execute(f, args.templateArgs)
	if err != nil {
		panic(fmt.Errorf("unable to execute template (%s): %s", inputFileName, err))
	}
}
