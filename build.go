package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func build(project projectConfig) error {
	tag := project.Registry + "/" + project.RegistryOrg + "/" + project.RegistryRepo + ":" + project.RegistryTag                   // tag follows a particular formula
	args := "docker buildx build --progress plain --push --platform " + project.ArchsString + " --tag " + tag + " " + project.Path // build arguments
	if len(project.BuildArgs) > 0 {
		for _, arg := range project.BuildArgs {
			args = args + " --build-arg " + arg
		}
	}
	log.Println("Building with args:", args) // notify of what we're running to build so the user can debug easily

	cmd := exec.Command("sh", "-c", args) // set the command

	cmd.Start() // run the command

	cmd.Wait() // waits until everything has finished

	exitCode := cmd.ProcessState.ExitCode()

	if exitCode == 0 {
		return nil // assume no errors
	} else {
		return errors.New(fmt.Sprintf("build failed with code %d", exitCode))
	}
}
