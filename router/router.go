package router

import (
	"net/http"

	"github.com/azulgautam79/go-sqlite-tasks/scalar"
	"github.com/azulgautam79/go-sqlite-tasks/task"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(taskHandler *task.Handler) http.Handler {
	mux := http.NewServeMux()

	//! Tasks
	mux.HandleFunc("GET /api/tasks", taskHandler.GetTasks)
	mux.HandleFunc("POST /api/tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /api/tasks/{id}", taskHandler.GetTask)
	mux.HandleFunc("PUT /api/tasks/{id}", taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", taskHandler.DeleteTask)

	//* Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /openapi.json", openAPISpec)

	//? Scalar
	mux.HandleFunc("GET /docs", scalar.Handler)
	return mux
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.json")
}
