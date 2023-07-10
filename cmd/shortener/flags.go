package main

import (
	"flag"
)

var (
	defaultServAddr     = "localhost:8080"
	defaultBaseShortUrl = "http://localhost:8080"
)

var (
	flagServAddr     string
	flagBaseShortUrl string
)

func initFlag() {
	flag.StringVar(&flagServAddr, "a", defaultServAddr, "base server address")
	flag.StringVar(&flagBaseShortUrl, "b", defaultBaseShortUrl, "base address short url")
}
