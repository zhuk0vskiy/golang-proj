package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

const configPath = "config/config.yml"

type Config struct {
	Logger     LoggerConfig     `yaml:"logger"`
	HTTP       HTTPConfig       `yaml:"http"`
	Database   DatabaseConfig   `yaml:"database"`
	JwtKey     string           `yaml:"JwtKey"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
}

type PrometheusConfig struct {
	MetricHost string `yaml:"metricHost"`
	MetricPort int `yaml:"metricPort"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type HTTPConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig
	MongoDB  MongoDBConfig
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Driver   string `yaml:"driver"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Bucket   string `yaml:"bucket"`
}

func NewConfig() (*Config, error) {
	var err error
	var config Config
	backendType := os.Getenv("BACKEND_TYPE")

	viper.SetConfigFile(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	if backendType == "docker" {
		backendPort := os.Getenv("BACKEND_PORT")

		postgresPort, err := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
		if err != nil {
			return nil, err
		}

		config.HTTP = HTTPConfig{Port: backendPort}

		config.Database.Postgres.Host = os.Getenv("POSTGRESQL_HOST")
		config.Database.Postgres.Port = postgresPort
		config.Database.Postgres.User = os.Getenv("POSTGRESQL_USERNAME")
		config.Database.Postgres.Password = os.Getenv("POSTGRESQL_PASSWORD")
		config.Database.Postgres.Database = os.Getenv("POSTGRESQL_DATABASE")
	}

	return &config, nil

}
