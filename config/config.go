package config

import (
	"log"
	"os"
)

type Config struct {
	Database Database `yaml:"database"`
	Kafka    Kafka    `yaml:"kafka"`
	App      App      `yaml:"app"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Kafka struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	cfg *Config
)

func InitConfig() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
		}
	}(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalln(err)
	}
	log.Println("Config is initialized")
	cfg = &config
}

func GetConfig() *Config {
	return cfg
}
