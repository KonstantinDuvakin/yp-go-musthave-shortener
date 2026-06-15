package config

import "flag"

type Config struct {
	BaseUrl    string
	ServerAddr string
}

func CreateConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.ServerAddr, "a", "localhost:8080", "The address to listen on for HTTP requests.")
	flag.StringVar(&config.BaseUrl, "b", "http://localhost:8080", "The base url for short links.")

	flag.Parse()

	return config
}
