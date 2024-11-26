package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Database *Database `yaml:"database"`
	Redis    *Redis    `yaml:"redis"`
	Kafka    *Kafka    `yaml:"kafka"`
	App      *App      `yaml:"app"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

type Kafka struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Topic string `yaml:"topic"`
}

type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type ConfigProvider interface {
	GetConfig() *Config
}

type ConfigManager struct {
	config *Config
}

func NewConfigManager() *ConfigManager {
	cfg := &ConfigManager{}
	cfg.InitConfig()
	return cfg
}

func (cm *ConfigManager) InitConfig() {
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
	cm.config = &config
	log.Println("Config is initialized")
}

func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}
