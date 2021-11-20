package main

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type projectConfig struct {
	Path          string
	Remote        string   `yaml:"remote"`
	Protocol      string   `yaml:"protocol"`
	Scm           string   `yaml:"scm"`
	Org           string   `yaml:"organization"`
	Repo          string   `yaml:"repository"`
	Tag           string   `yaml:"tag"`
	Registry      string   `yaml:"registry"`
	RegistryOrg   string   `yaml:"registry_organization"`
	RegistryRepo  string   `yaml:"registry_repository"`
	RegistryTag   string   `yaml:"registry_tag"`
	BuildArgs     []string `yaml:"build_args"`
	Architectures []string `yaml:"architectures"`
	ArchsString   string
}

type config struct {
	Projects []projectConfig `yaml:"projects"`
}

func getConf(path string) (config, error) {
	conf := config{}

	// read file to mem
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	// parse yaml to golang
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func setProjectFromRemote(project projectConfig) (projectConfig, error) {
	if strings.HasPrefix(project.Remote, "http") {
		splitSlash := strings.Split(project.Remote, "/")
		project.Protocol = splitSlash[0]
		project.Scm = splitSlash[2]
		project.Org = splitSlash[3]
		project.Repo = splitSlash[4]
		if strings.HasSuffix(project.Repo, ".git") {
			project.Repo = strings.ReplaceAll(project.Repo, ".git", "")
		}
		if len(project.Tag) < 1 {
			project.Tag = magicLatest
		}
		project.Path = os.TempDir() + "/" + cmdname + "/" + project.Org + "/" + project.Repo
		if len(project.RegistryOrg) == 0 {
			project.RegistryOrg = project.Org
		}
		if len(project.RegistryRepo) == 0 {
			project.RegistryRepo = project.Repo
		}
		if len(project.RegistryTag) == 0 {
			project.RegistryTag = project.Tag
		}
		if len(project.Architectures) == 0 {
			project.Architectures = []string{"linux/amd64"}
		}
		if len(project.Architectures) > 0 {
			project.ArchsString = strings.Join(project.Architectures, ",")
		}
	}

	err := verifyProjectValuesNotEmpty(project)
	if err != nil {
		return project, err
	}

	return project, nil
}

func verifyProjectValuesNotEmpty(project projectConfig) error {
	vc := reflect.ValueOf(project)
	for i := 0; i < vc.NumField(); i++ {
		key := vc.Type().Field(i).Name
		value := vc.Field(i).String()
		if key == "Tag" && value == magicLatest {
			continue
		}
		if len(value) < 1 {
			return errors.New(key + " must not be empty")
		}
	}

	return nil
}
