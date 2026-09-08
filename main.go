// @title Go SQLite Task API
// @version 1.0
// @description A simple REST API built with Go and SQLite.
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azulgautam79/go-sqlite-tasks/internal/app"
	configuration "github.com/azulgautam79/go-sqlite-tasks/internal/config"

	_ "github.com/azulgautam79/go-sqlite-tasks/docs"
)

func main() {
	cfg := configuration.Load()

	handler, cleanup, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer cleanup()

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server listening on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf(
			"Server forced to shutdown: %v",
			err,
		)
	}

	log.Println("Server stopped")
}
