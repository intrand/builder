package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "" // to be filled in by goreleaser
	commit  = "" // to be filled in by goreleaser
	date    = "" // to be filled in by goreleaser
	builtBy = "" // to be filled in by goreleaser
	cmdname = filepath.Base(os.Args[0])
)

const magicLatest = "getLatest"

func main() {
	// get input
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// get config file (remotes)
	conf, err := getConf(*cmd_config)
	if err != nil {
		log.Fatalln(err)
	}

	// loop remotes
	for _, project := range conf.Projects {
		fmt.Print("\n")
		project, err = setProjectFromRemote(project)
		if err != nil {
			fmt.Println("Skipping due to error: ", err)
			continue
		}

		log.Println("Building " + project.Remote + " in path " + project.Path)

		repo, err := clone(project.Remote, project.Path)
		if err != nil {
			log.Fatalln(err)
		}

		project.Tag, err = checkoutRef(repo, project.Tag)
		if err != nil {
			log.Fatalln(err)
		}
		if project.RegistryTag == magicLatest {
			project.RegistryTag = project.Tag
		}

		err = build(project)
		if err == nil {
			log.Println("Success!")
		} else {
			log.Println(project.Remote+":", err)
			continue
		}
		fmt.Print("\n\n")
	}
}
