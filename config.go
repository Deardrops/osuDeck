package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Config struct {
	OsuApiKey        string `yaml:"osu_api_key"`
	LocalOsuRootPath string `yaml:"local_osu_root_path"`
	OsuDbPath        string `yaml:"osu_db_path"`
	CollectionDbPath string `yaml:"collection_db_path"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Mirror           string `yaml:"mirror"`
}

const configFilePath = "./conf.yaml"

func (conf *Config) init() {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		conf.Mirror = "bloodcat"
		conf.output()
		return
	}
	buff, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) output() {
	buff, err := yaml.Marshal(&conf)
	if err != nil {
		log.Fatal(err)
	}
	workdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(workdir, configFilePath), buff, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
