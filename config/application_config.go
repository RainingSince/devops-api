package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DbConfig struct {
	Hosts    string `yaml:"hosts"`
	DataBase string `yaml:"dataBase"`
}

type AuthConfig struct {
	SigningKey string `yaml:"signingKey"`
	IngorePath string `yaml:"ingorePath"`
}

type AppConfig struct {
	DbConfig   `yaml:"db"`
	Port       string `yaml:"port"`
	AuthConfig `yaml:"auth"`
}

var GlobabConfig *AppConfig

func LoadConfig(filePath string) (conf *AppConfig, err error) {
	conf = new(AppConfig)
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, conf)
	if err != nil {
		panic(err)
	}
	GlobabConfig = conf
	return conf, err
}
