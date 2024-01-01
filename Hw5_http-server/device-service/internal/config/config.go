package config

import (
	"os"
)

type Config struct {
	Address string `yaml:"address"`
}

func MustLoad() *Config {

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		panic("HOST is not set")
	}
	if port == "" {
		panic("PORT is not set")
	}

	address := host + ":" + port

	cfg := Config{address}

	return &cfg
}
