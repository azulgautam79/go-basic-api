package handler

import (
	"net/http"
	"sync"

	"github.com/azulgautam79/go-sqlite-tasks/internal/app"
	configuration "github.com/azulgautam79/go-sqlite-tasks/internal/config"
)

var (
	handler http.Handler
	once    sync.Once
	err     error
)

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		cfg := configuration.Load()

		handler, _, err = app.New(cfg)
	})

	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)

		return
	}

	handler.ServeHTTP(w, r)
}