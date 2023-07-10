package main

import (
	"flag"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/configs"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"net/http"
)

func main() {
	initFlag()
	flag.Parse()

	configs.InitConfig(flagServAddr, flagBaseShortURL)

	linkStorage := hashmapstorage.NewLinkStorage(make(map[string]domain.Link))
	linkService := service.NewLinkService(linkStorage)
	linkHandler := handlers.NewLinkHandler(linkService, flagBaseShortURL)

	if err := http.ListenAndServe(configs.AppConfig.ServAddr, linkHandler.InitRouter()); err != nil {
		panic(err)
	}
}
