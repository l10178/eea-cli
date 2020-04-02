package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Cloud  `yaml:"cloud"`
	GitLab `yaml:"gitlab"`
}

type Cloud struct {
	AccessKeyId     string `yaml:"access-key-id"`
	SecretAccessKey string `yaml:"secret-access-key"`
	Region          string `yaml:"region"`
	Registry        string `yaml:"registry"`
}

type GitLab struct {
	PrivateToken    string `yaml:"private-token"`
	GitRoot         string `yaml:"git-root"`
	ApiRoot         string `yaml:"api-root"`
	TagSleepSeconds string `yaml:"tag-sleep-seconds"`
}

func GetConfig() Config {
	configFile := os.Getenv("EEA_CFG")
	if configFile == "" {
		configFile = "./config.yaml"
	}

	yamlFile, err := ioutil.ReadFile(configFile)

	conf := Config{}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal config file error: %v", err)
	}
	return conf
}
