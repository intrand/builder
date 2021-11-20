package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// global
	app        = kingpin.New("builder", "builds things using docker buildx").Author("intrand")
	cmd_config = app.Flag("config", "config").Envar("builder_config").Default("config.yml").String()
	cmd_user   = app.Flag("user", "git remote user").Envar("builder_user").String()
	cmd_token  = app.Flag("token", "git remote token").Envar("builder_token").String()
)
