package router

import "net/http"

func (r *Router) registerTaskRoutes(mux *http.ServeMux) {

	mux.HandleFunc("GET /api/v1/tasks", r.taskHandler.GetTasks)
	mux.HandleFunc("POST /api/v1/tasks", r.taskHandler.CreateTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", r.taskHandler.GetTask)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", r.taskHandler.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", r.taskHandler.DeleteTask)
}
