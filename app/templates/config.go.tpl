package main

import "github.com/skeletonkey/lib-core-go/config"

var cfg *{{ .ConfigName }}

func getConfig() *{{ .ConfigName }} {
	config.LoadConfig("{{ .ConfigName }}", &cfg)
	return cfg
}
