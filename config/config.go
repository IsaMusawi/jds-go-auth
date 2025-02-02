package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName string `yaml:"app_name"`
	Port string `yaml:"port"`
	Environment string `yaml:"env"`
	EnvConf map[string]EnvConfig `yaml:",inline"`
	MainConfig EnvConfig
	
}

type EnvConfig struct {
	JwtSecret string `yaml:"jwt_secret"`
	ApiUrl ApiUrlConfig `yaml:"api_url"`
}

type ApiUrlConfig struct {
	Data string `yaml:"data"`
	CurrencyConversion string `yaml:"currency_conversion"`
}

func Init() Config {
	var (
		config Config
		exists bool
	)
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)  
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshal config file: %v", err)  
	}

	config.MainConfig, exists = config.EnvConf[config.Environment]
	if !exists {
		log.Fatalf("Environment not found") 
	}
	config.EnvConf = nil

	return config
}