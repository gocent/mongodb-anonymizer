package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"sync"
)

type Config struct {
	DB uri
}

type uri struct {
	SourceURI   url.URL  `env:"SOURCE_DB_URI,required"`
	SourceName  string   `env:"SOURCE_DB_NAME,required"`
	TargetURI   url.URL  `env:"TARGET_DB_URI,required"`
	TargetName  string   `env:"TARGET_DB_NAME,required"`
	Collections []string `env:"COLLECTIONS,required"`
}

var (
	once     sync.Once
	instance *Config
	err      error
)

func GetEnv() *Config {
	once.Do(func() {
		instance, err = getInstance()
	})
	if err != nil {
		log.Fatal("Unable to parse env files")
	}
	return instance
}

func getInstance() (*Config, error) {
	environment := os.Getenv("MONGODB_ANONYMIZER_ENV")
	if environment == "" {
		environment = "development"
	}
	var err error

	switch environment {
	case "prod":
		err = godotenv.Load(".env")
	case "staging":
		err = godotenv.Load(".env." + environment)
	default:
		err = godotenv.Load(".env." + environment + ".local")
	}
	if err != nil {
		log.Printf("Unable to load environment file for %s", environment)
		return nil, err
	}
	cfg := Config{}
	err = env.Parse(&cfg.DB)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return &cfg, err
}
