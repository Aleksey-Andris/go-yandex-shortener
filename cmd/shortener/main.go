package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/configs"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/delivery/handlers"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/logger"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/service"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/hashmapstorage"
	"github.com/Aleksey-Andris/go-yandex-shortener/internal/app/storage/postgresstorage"
)

func main() {
	initFlag()
	flag.Parse()
	configs.InitConfig(flagServAddr, flagBaseShortURL)

	if err := logger.Initialize(flagLogLevel); err != nil {
		log.Fatal(err)
	}
	var userStorage service.UserStorage
	var linkStorage service.LinkStorage
	var db *sqlx.DB
	var err error
	if flagConfigDB == "" {
		linkStorage, err = hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), flagFileStoragePath)
		if err != nil {
			logger.Log().Fatal(err.Error())
		}
		userStorage, err = hashmapstorage.NewLinkStorage(make(map[string]*domain.Link), make(map[int32][]*domain.Link), flagFileStoragePath)
		if err != nil {
			logger.Log().Fatal(err.Error())
		}
	} else {
		db, err = postgresstorage.NewPostgresDB(flagConfigDB)
		if err != nil {
			logger.Log().Fatal(err.Error())
		}
		linkStorage, err = postgresstorage.NewLinkStorage(db)
		if err != nil {
			logger.Log().Fatal(err.Error())
		}
		userStorage, err = postgresstorage.NewUserStorage(db)
		if err != nil {
			logger.Log().Fatal(err.Error())
		}
	}
	servises := handlers.NewServices(linkStorage, userStorage)
	handler := handlers.NewHandler(servises, flagBaseShortURL)
	router := handler.InitRouter()
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

	srv := &http.Server{
		Addr:    configs.AppConfig.ServAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	logger.Log().Debug("shutting down")

	context, gansel := context.WithTimeout(context.Background(), time.Second*10)
	defer gansel()

	if err := srv.Shutdown(context); err != nil {
		logger.Log().Error(err.Error())
	}

	handler.FlushMessagesDeleteNow()

	if err := linkStorage.Close(); err != nil {
		logger.Log().Error(err.Error())
	}
}
