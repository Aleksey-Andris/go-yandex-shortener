package main

import (
	"flag"
	"os"
)

var (
	defaultServAddr            = "localhost:8080"
	defaultBaseShortURL        = "http://localhost:8080"
	defaultLogLevel            = "info"
	defaitflagFileStoragePath  = "/tmp/short-url-db.json"
	defaultFlagFileStoragePath = "host=localhost port=5432 user=postgres password=shortener dbname=postgres sslmode=disable"
)

var (
	flagServAddr        string
	flagBaseShortURL    string
	flagLogLevel        string
	flagFileStoragePath string
	flagConfigDB        string
)

func initFlag() {
	flag.StringVar(&flagServAddr, "a", defaultServAddr, "base server address")
	flag.StringVar(&flagBaseShortURL, "b", defaultBaseShortURL, "base address short URL")
	flag.StringVar(&flagLogLevel, "l", defaultLogLevel, "log level")
	flag.StringVar(&flagFileStoragePath, "f", defaitflagFileStoragePath, "file storage path")
	flag.StringVar(&flagConfigDB, "d", defaultFlagFileStoragePath, "file storage path")

	if envServAddr := os.Getenv("SERVER_ADDRESS"); envServAddr != "" {
		flagServAddr = envServAddr
	}
	if envBaseShortURL := os.Getenv("BASE_URL"); envBaseShortURL != "" {
		flagBaseShortURL = envBaseShortURL
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		flagFileStoragePath = envFileStoragePath
	}
	if envConfigDB := os.Getenv("DATABASE_DSN"); envConfigDB != "" {
		flagConfigDB = envConfigDB
	}
}