package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rvecwxqz/apod/internal/config"
	apodfetcher "github.com/rvecwxqz/apod/internal/fetcher/apod-fetcher"
	"github.com/rvecwxqz/apod/internal/http-server/handlers/get"
	getall "github.com/rvecwxqz/apod/internal/http-server/handlers/get-all"
	minio "github.com/rvecwxqz/apod/internal/storage/minio"
	"github.com/rvecwxqz/apod/internal/storage/postgresql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	fmt.Println(cfg)

	storage, err := postgresql.New(context.Background(), cfg.DataBaseDSN)
	if err != nil {
		log.Fatal(fmt.Errorf("create storage error: %w", err))
	}
	binaryStorage, err := minio.NewProvider(
		cfg.MinioUser,
		cfg.MinioPass,
		cfg.MinioEndpoint,
		cfg.MinioBucket,
		cfg.ServerAddress,
		cfg.MinioPort,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("creating binary storage error: %w", err))
	}
	apodfetcher.New(
		ctx,
		cfg.APIKey,
		cfg.WorkerInterval,
		cfg.WorkerRetries,
		storage,
		binaryStorage,
	)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/get", get.New(storage, binaryStorage))
		r.Post("/get_all", getall.New(storage, binaryStorage))
	})

	serv := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}
	go func() {
		if err = serv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println(fmt.Sprintf("Server started. Port: %v", cfg.ServerPort))

	<-ctx.Done()
	storage.Stop()

	log.Println("Stopped")

}
