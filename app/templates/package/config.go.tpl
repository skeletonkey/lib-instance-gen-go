package {{ .PackageName }}

import "github.com/skeletonkey/lib-core-go/config"

//nolint:gochecknoglobals // cfg holds the configuration data, which should only be retrieved and not changed
var cfg *{{ .PackageName }}

func getConfig() *{{ .PackageName }} {
	config.LoadConfig("{{ .PackageName }}", &cfg)
	return cfg
}
