package main

import (
	"flag"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/configs"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/logger"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"log"
	"net/http"
)

func main() {
	initFlag()
	flag.Parse()
	configs.InitConfig(flagServAddr, flagBaseShortURL)

	if err := logger.Initialize(flagLogLevel); err != nil {
		log.Fatal(err)
	}

	linkStorage, err := hashmapstorage.NewLinkStorage(make(map[string]domain.Link), flagFileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(){
		if err := linkStorage.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	
	linkService := service.NewLinkService(linkStorage, linkStorage.GetSequense())
	linkHandler := handlers.NewLinkHandler(linkService, flagBaseShortURL)

	if err := http.ListenAndServe(configs.AppConfig.ServAddr, linkHandler.InitRouter()); err != nil {
		log.Fatal(err)
	}
}
