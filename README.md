# lib-instance-gen-go

Library for starting new Go applications allowing for creating and updating common files.

## app

Create a skeleton structure for a new application

See godocs for specifics.

Create an 'init' go file in the root of your repository and use 'go generate'.  This will create the skeleton of your
application with boiler code providing the following:

 * config ingestion
 * logging (utilizes zerolog)

### Example File

Filename: app-init.go

```go
package main

//go:generate go run app-init.go

import instance_gen "github.com/skeletonkey/lib-instance-gen-go/app"

func main() {
	app := instance_gen.NewApp("rachio-next-run", "app")
	app.WithPackages("logger", "pushover", "rachio").
		WithGithubWorkflows("linter", "test").
		WithMakefile()
}
```