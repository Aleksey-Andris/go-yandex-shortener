package main

import (
	handlers2 "github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"net/http"
)

func main() {
	linkStorage := hashmapstorage.NewLinkStorage()
	linkService := service.NewLinkService(linkStorage)
	linkHandler := handlers2.NewLinkHandler(linkService)

	if err := http.ListenAndServe(":8080", linkHandler.InitServeMux()); err != nil {
		panic(err)
	}
}
