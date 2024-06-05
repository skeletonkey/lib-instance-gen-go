# lib-instance-gen-go

Library for starting new Go applications allowing for creating and updating common files.

## app

Create a skeleton structure for a new application

See godocs for specifics.

Create an 'init' go file in the root of your repository and use `go generate`.
This will create the skeleton of your application with boiler code providing the following:

 * config ingestion
 * logging (utilizes zerolog)

### Example File

Filename: app-init.go

```go
package main

//go:generate go run app-init.go

import instanceGen "github.com/skeletonkey/lib-instance-gen-go/app"

func main() {
	app := instanceGen.NewApp("rachio-next-run", "app")
	app.SetupApp(
		app.WithCGOEnabled(),
		app.WithConfig(),
		app.WithDependencies(
            "github.com/labstack/echo/v4",
        ),
		app.WithGithubWorkflows("changelog", "linter", "test"),
		app.WithMakefile(),
        app.WithGoVersion("1.22"),
        app.WithPackages("logger", "pushover", "rachio"),
	).Generate()

}
```
