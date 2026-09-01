package app

import (
	"database/sql"
	"net/http"

	configuration "github.com/azulgautam79/go-sqlite-tasks/internal/config"
	"github.com/azulgautam79/go-sqlite-tasks/router"
	"github.com/azulgautam79/go-sqlite-tasks/task"

	_ "modernc.org/sqlite"
)

func New(cfg configuration.Config) (http.Handler, func(), error) {
	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0
		)
	`)

	if err != nil {
		db.Close()
		return nil, nil, err
	}

	repository := task.NewRepository(db)

	service := task.NewService(repository)

	handler := task.NewHandler(service)

	r := router.New(handler)

	cleanup := func() {
		db.Close()
	}

	return r, cleanup, nil
}
