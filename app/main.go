// Package app is a library for creating a bare-bones application with boilerplate files taken care of.
// The long-term goal is to be the basis of all Go applications allowing for quick propagation of updates,
// bug fixes, and new features.
//
// Reference README.md for examples.
//
//	app := instanceGen.NewApp("rachio-next-run", "app")
//	app.SetupApp(
//		app.WithGithubWorkflows("linter", "test"),
//		app.WithGoVersion("1.23"),
//		app.WithMakefile(),
//		app.WithPackages("logger", "pushover", "rachio"),
//	).Generate()
//
// Each generated file will be prepended with a 'warning' comment to not edit the file.
package app

import (
	"embed"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

//go:embed all:templates
var templatesFS embed.FS

const (
	cgoEnabled         = "CGOEnabled"
	codeOwnersFileName = "CODEOWNERS"
	dependencies       = "dependencies"
	goModFile          = "go.mod"
	goVersion          = "GoVersion"
	mkfilesSubDir      = "Makefile"
	newDirPermission   = 0750
	newFilePermission  = 0644
	templateBaseDir    = "templates"
	warning            = "lib-instance-gen-go: File auto generated -- DO NOT EDIT!!!\n"
)

type setupOp func(App) error

var templateExts = map[string]string{
	"go":     ".go.tpl",
	"mkfile": ".tpl",
	"toml":   ".tpl",
	"yml":    ".yml.tpl",
}
var warnings = map[string]string{
	codeOwnersFileName: "# " + warning,
	"go":               "// " + warning,
	"mkfile":           "// " + warning,
	"toml":             "# " + warning,
	"yml":              "# " + warning,
}

// App struct containing necessary information for a new application.
type App struct {
	binaryName string         // name of the binary the 'make' will produce
	dir        string         // subdirectory which will contain the program's source code
	ops        []setupOp      // list of operations to perform during Generate method
	settings   map[string]any // misc settings
}

func noOp() setupOp {
	return func(_ App) error { return nil }
}

// NewApp returns the struct for a new applications which allows for generating boilerplate files.
//   - binaryName is used by the Makefile for the build command
//   - dir is the subdirectory that packages will be created in
func NewApp(binaryName string, dir string) App {
	return App{binaryName: binaryName, dir: dir, settings: make(map[string]any)}
}

// SetupApp takes a list of With* functions that will be applied to the Application.
func (a App) SetupApp(ops ...setupOp) App {
	a.ops = ops
	return a
}

// Generate will apply all the settings and create the boilerplate files.
func (a App) Generate() {
	for _, op := range a.ops {
		err := op(a)
		if err != nil {
			panic(err)
		}
	}
}

// WithCodeOwners creates a CODEOWNERS file with the provided codeowners configuration.
// This is provided for simple CODEOWNERS use cases to centralize all project configurations in the "init.go" file.
// Each string provided will be written on a single line, which provides flexibility.
// However, once the file becomes complicated, it will be best to create the file manually.
func (App) WithCodeOwners(codeOwners ...string) setupOp {
	if len(codeOwners) == 0 {
		return noOp()
	}

	return func(_ App) error {
		file, err := os.OpenFile(codeOwnersFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, newFilePermission)
		if err != nil {
			return fmt.Errorf("unable to write file (%s): %s", codeOwnersFileName, err)
		}
		defer func() {
			_ = file.Close()
		}()

		warning, found := warnings[codeOwnersFileName]
		if !found {
			return fmt.Errorf("unable to find a 'warnings' entry for %s", codeOwnersFileName)
		}
		if _, err := file.WriteString(warning); err != nil {
			return fmt.Errorf("unable to write warning to file (%s): %s", codeOwnersFileName, err)
		}
		if _, err := file.WriteString(strings.Join(codeOwners, "\n")); err != nil {
			return fmt.Errorf("unable to write codeOwners (%s) to file (%s): %s", codeOwners, codeOwnersFileName, err)
		}

		return nil
	}
}

// WithPackages takes a list of strings which results in creating a skeleton subdirectory for each.
// Foreach package listed the following will be created:
//   - config.go - template to use github.com/skeletonkey/lib-core-go/config module
func (App) WithPackages(packageNames ...string) setupOp {
	return func(a App) error {
		for _, name := range packageNames {
			packageName := name
			templateArgs := templateArgs{
				PackageName: packageName,
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
		return nil
	}
}

// WithCGOEnabled will add CGO_ENABLED=1 to the build statement
func (a App) WithCGOEnabled() setupOp {
	a.settings[cgoEnabled] = true
	return noOp()
}

// WithConfig adds a config file for the main app. Config
func (App) WithConfig() setupOp {
	return func(a App) error {
		templateArgs := templateArgs{
			ConfigName: a.dir,
		}
		generateTemplate(generateTemplateArgs{
			fileType:       "go",
			outputName:     "config.go",
			outputSubDir:   a.dir,
			templateName:   "config" + templateExts["go"],
			templateSubDir: "",
			templateArgs:   templateArgs,
		})
		return nil
	}
}

// WithDependencies received a list of strings that are Go libraries that should only be updated with 'make golib-latest'
func (a App) WithDependencies(deps ...string) setupOp {
	a.settings[dependencies] = deps
	return noOp()
}

// WithGithubWorkflows sets up the specified workflows.
// Current supported work flows:
//   - linter - on pull request for all branches
//   - test - on pull request for all branches
func (App) WithGithubWorkflows(flows ...string) setupOp {
	return func(a App) error {
		tmplArgs := templateArgs{}
		if ver, ok := a.settings[goVersion]; !ok {
			panic(fmt.Errorf("no %s provided - please call WithGoVersion", goVersion))
		} else {
			tmplArgs.GoVersion = ver.(string)
		}

		linterPresent := false
		for _, name := range flows {
			if name == "linter" {
				linterPresent = true
			}
			generateTemplate(generateTemplateArgs{
				fileType:       "yml",
				outputName:     name + ".yml",
				outputSubDir:   path.Join(".github", "workflows"),
				templateArgs:   tmplArgs,
				templateName:   name + templateExts["yml"],
				templateSubDir: "github_workflows",
			})
		}

		if linterPresent {
			generateTemplate(generateTemplateArgs{
				fileType:       "toml",
				outputName:     ".golangci.toml",
				outputSubDir:   "",
				templateArgs:   templateArgs{},
				templateName:   ".golangci.toml",
				templateSubDir: "",
			})
		}

		return nil
	}
}

// WithGoVersion provide the current version of Go to use for GitHub actions
// and the go.mod file
func (a App) WithGoVersion(ver string) setupOp {
	a.settings[goVersion] = ver

	return func(_ App) error {
		_, err := os.Stat(goModFile)
		if err == nil { // we have a go mod file, and we can replace the version
			data, err := os.ReadFile(goModFile)
			if err != nil {
				return fmt.Errorf("unable to read go.mod file (%s): %s", goModFile, err)
			}

			pattern := regexp.MustCompile(`(?m)$\s*go \d+\.\d+(\.\d+)?\s*$`)
			newData := pattern.ReplaceAll(data, []byte("\n\ngo "+ver+"\n"))

			err = os.WriteFile(goModFile, newData, newFilePermission)
			if err != nil {
				return fmt.Errorf("unable to write go.mod file (%s): %s", goModFile, err)
			}
		}
		return nil
	}
}

// WithMakefile creates the basic Makefile with:
//   - build - runs fmt, vet and then builds the binary
//   - install - move binary to /usr/local/bin
//   - golib-latest - install go dependencies
//   - app-init - generate the boilerplate
//
// Method accepts a list of strings which will be used to create "include" statements.
// Each string will be concatenated with "Makefile.". This allows for custom "make" commands
// for a project. These customer make files will not be generated nor effected by app-init.
func (App) WithMakefile(makeExt ...string) setupOp {
	return func(a App) error {
		nodes, err := templatesFS.ReadDir(path.Join(templateBaseDir, mkfilesSubDir))
		if err != nil {
			return fmt.Errorf("unable to read dir (%s): %s", mkfilesSubDir, err)
		}

		tmplArgs := templateArgs{BinaryName: a.binaryName}
		if yes, ok := a.settings[cgoEnabled]; ok && yes.(bool) {
			tmplArgs.BuildEnvArgs = "CGO_ENABLED=1 "
		}

		if deps, ok := a.settings[dependencies]; ok {
			depString := ""
			for _, dep := range deps.([]string) {
				depString = fmt.Sprintf("%sgo get -u %s@latest\n\t", depString, dep)
			}
			tmplArgs.Dependencies = depString
		}

		for _, ext := range makeExt {
			tmplArgs.Includes = fmt.Sprintf("%sinclude Makefile.%s\n", tmplArgs.Includes, ext)
		}

		for _, node := range nodes {
			generateTemplate(generateTemplateArgs{
				fileType:       "mkfile",
				outputName:     node.Name()[:len(node.Name())-len(templateExts["mkfile"])],
				outputSubDir:   "",
				templateName:   node.Name(),
				templateSubDir: mkfilesSubDir,
				templateArgs:   tmplArgs,
			})
		}
		return nil
	}
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
	BinaryName   string // name of the executable program
	BuildEnvArgs string // any env args that are needed when building the app
	ConfigName   string // name of the config element for the main program
	Dependencies string // see WithDependencies
	GoVersion    string // see WithGoVersion
	Includes     string // custom Makefile.{name} to be included - specified by WithMakefile
	PackageName  string // name of the package
}

func generateTemplate(args generateTemplateArgs) {
	inputFileName := path.Join(templateBaseDir, args.templateSubDir, args.templateName)
	outputFileName := path.Join(args.outputSubDir, args.outputName)

	if args.outputSubDir != "" {
		err := os.MkdirAll(args.outputSubDir, newDirPermission)
		if err != nil {
			panic(fmt.Errorf("unable to create directory structure (%s): %s", args.outputSubDir, err))
		}
	}
	//nolint:gosec // G304 -- This is safe as the file is being opened for write/create/truncation - no reading
	f, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, newFilePermission)
	if err != nil {
		panic(fmt.Errorf("unable to create file (%s): %s", outputFileName, err))
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(fmt.Errorf("error closing file (%s): %s", outputFileName, err))
		}
	}()
	_, err = f.WriteString(warnings[args.fileType])
	if err != nil {
		panic(fmt.Errorf("unable to write warning to file (%s): %s", outputFileName, err))
	}
	temp := template.Must(template.ParseFS(templatesFS, inputFileName))
	err = temp.Execute(f, args.templateArgs)
	if err != nil {
		panic(fmt.Errorf("unable to execute template (%s): %s", inputFileName, err))
	}
}
