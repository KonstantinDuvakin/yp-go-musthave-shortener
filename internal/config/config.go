package config

import (
	"flag"
	"os"
)

type Config struct {
	BaseUrl    string
	ServerAddr string
}

func NewConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&config.BaseUrl, "b", "http://localhost:8080", "The base url for short links.")

	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		config.ServerAddr = envServerAddr
	}

	if envBaseUrl := os.Getenv("BASE_URL"); envBaseUrl != "" {
		config.BaseUrl = envBaseUrl
	}

	return config
}
