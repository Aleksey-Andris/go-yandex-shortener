package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/configs"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/logger"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/postgresstorage"
	"github.com/jmoiron/sqlx"
)

func main() {
	initFlag()
	flag.Parse()
	configs.InitConfig(flagServAddr, flagBaseShortURL)

	if err := logger.Initialize(flagLogLevel); err != nil {
		log.Fatal(err)
	}

	var linkStorage service.LinkStorage
	var db *sqlx.DB
	var err error
	if flagConfigDB == "" {
		linkStorage, err = hashmapstorage.NewLinkStorage(make(map[string]domain.Link), flagFileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = postgresstorage.NewPostgresDB(flagConfigDB)
		if err != nil {
			log.Fatal(err)
		}
		linkStorage, err = postgresstorage.NewLinkStorage(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer func() {
		if err := linkStorage.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	linkService := service.NewLinkService(linkStorage)
	linkHandler := handlers.NewLinkHandler(linkService, flagBaseShortURL)

	router := linkHandler.InitRouter()
	router.Get("/ping", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if db == nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := db.Ping(); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}))

	if err := http.ListenAndServe(configs.AppConfig.ServAddr, router); err != nil {
		log.Fatal(err)
	}
}
