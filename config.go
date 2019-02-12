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
	LocalOsuRootPath string `yaml:"local_osu_root_path"`
	OsuDbPath        string `yaml:"osu_db_path"`
	CollectionDbPath string `yaml:"collection_db_path"`
	OsuApiKey        string `yaml:"osu_api_key"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Mirror           string `yaml:"mirror"`
	FilePath         string
}

func newConf(paths ...string) Config {
	conf := Config{}
	var confPath string
	if len(paths) == 0 {
		confPath = "./conf.yaml"
	} else {
		confPath = paths[0]
	}
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		conf.Mirror = "bloodcat"
		conf.FilePath = confPath
		conf.output()
		return conf
	}
	buff, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(buff, &conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.FilePath = confPath
	return conf
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
	err = ioutil.WriteFile(path.Join(workdir, conf.FilePath), buff, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
