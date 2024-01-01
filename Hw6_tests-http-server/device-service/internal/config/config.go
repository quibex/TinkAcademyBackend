package config

import (
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()

	err := validate.Var(port, "numeric")
	if err != nil {
		panic(err)
	}

	err = validate.Var(host, "hostname_rfc1123")
	if err != nil {
		panic(err)
	}

	address := host + ":" + port

	cfg := Config{address}

	return &cfg
}
