package main

import "github.com/skeletonkey/lib-core-go/config"

//nolint:gochecknoglobals // cfg holds the configuration data, which should only be retrieved and not changed
var cfg *{{ .ConfigName }}

func getConfig() *{{ .ConfigName }} {
	config.LoadConfig("{{ .ConfigName }}", &cfg)
	return cfg
}
