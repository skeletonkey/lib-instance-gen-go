package {{ .PackageName }}

import "github.com/skeletonkey/rachio-next-run/app/config"

var cfg *{{ .PackageName }}

func getConfig() *{{ .PackageName }} {
	config.LoadConfig("{{ .PackageName }}", &cfg)
	return cfg
}
