package main

import (
	handlers2 "github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hash_map_storage"
	"net/http"
)

func main() {
	linkStorage := hash_map_storage.NewLinkStorage()
	linkService := service.NewLinkService(linkStorage)
	linkHandler := handlers2.NewLinkHandler(linkService)

	if err := http.ListenAndServe(":8080", linkHandler.InitServeMux()); err != nil {
		panic(err)
	}
}
