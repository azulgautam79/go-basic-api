package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azulgautam79/go-sqlite-tasks/internal/config"
	"github.com/azulgautam79/go-sqlite-tasks/internal/db"
	"github.com/azulgautam79/go-sqlite-tasks/internal/employees"
	"github.com/azulgautam79/go-sqlite-tasks/internal/router"
	"github.com/azulgautam79/go-sqlite-tasks/internal/task"
)

func main() {
	//! Config
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	//! Logger
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	logger.Info("database connected")

	//! Services
	//* Tasks
	taskRepository := task.NewRepository(db)
	taskService := task.NewService(taskRepository)
	taskHandler := task.NewHandler(taskService)
	//* Employees
	employeeRepository := employees.NewRepository(db)
	employeeService := employees.NewService(employeeRepository)
	employeeHandler := employees.NewHandler(employeeService)

	handler := router.New(
		taskHandler,
		employeeHandler,
	)

	//! Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		log.Printf("Server listening on port %s", cfg.Port)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {

			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutdown signal received")

	// Give active requests time to finish
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
