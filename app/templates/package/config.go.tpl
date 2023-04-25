package {{ .PackageName }}

import "github.com/skeletonkey/lib-core-go/config"

var cfg *{{ .PackageName }}

func getConfig() *{{ .PackageName }} {
	config.LoadConfig("{{ .PackageName }}", &cfg)
	return cfg
}
