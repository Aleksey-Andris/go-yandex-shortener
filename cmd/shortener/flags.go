package main

import (
	"flag"
	"os"
)

var (
	defaultServAddr     = "localhost:8080"
	defaultBaseShortURL = "http://localhost:8080"
)

var (
	flagServAddr     string
	flagBaseShortURL string
)

func initFlag() {
	flag.StringVar(&flagServAddr, "a", defaultServAddr, "base server address")
	flag.StringVar(&flagBaseShortURL, "b", defaultBaseShortURL, "base address short URL")

	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		flagServAddr = envServAddr
	}

	if envBaseShortURL := os.Getenv("BASE_URL"); envBaseShortURL != "" {
		flagBaseShortURL = envBaseShortURL
	}
}
