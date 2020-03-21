package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Cloud `yaml:"cloud"`
}

type Cloud struct {
	AccessKeyId     string `yaml:"access-key-id"`
	SecretAccessKey string `yaml:"secret-access-key"`
	Region          string `yaml:"region"`
	Registry        string `yaml:"registry"`
}

func GetConfig() (Config, error) {
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
	return conf, err
}
