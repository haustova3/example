package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceConf struct {
		Host       string `yaml:"host"`
		GRPCPort   string `yaml:"grpc_port"`
		HTTPPort   string `yaml:"http_port"`
		MetricPort string `yaml:"metric_port"`
	} `yaml:"service"`

	NotificationConf struct {
		MaxCount int `yaml:"max_count"`
		Timer    int `yaml:"timer"`
	} `yaml:"notification"`

	ProductsConf struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"products"`

	UsersConf struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"users"`

	KafkaConf struct {
		OrderTopic string `yaml:"order_topic"`
		Brokers    string `yaml:"brokers"`
	} `yaml:"kafka"`

	JaegerConf struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"jaeger"`

	DBConf struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"postgres"`
}

func LoadConfig(filename string) (*Config, error) {
	f, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("load config failed: %w", err)
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	config := &Config{}
	config.NotificationConf.MaxCount = 100
	config.NotificationConf.Timer = 300
	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
