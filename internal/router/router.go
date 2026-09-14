package router

import (
	"net/http"

	"github.com/azulgautam79/go-sqlite-tasks/internal/scalar"
	"github.com/azulgautam79/go-sqlite-tasks/internal/task"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	taskHandler *task.Handler
}

func New(taskHandler *task.Handler) http.Handler {

	r := &Router{
		taskHandler: taskHandler,
	}

	return r.routes()
}

func (r *Router) routes() http.Handler {
	mux := http.NewServeMux()

	r.registerTaskRoutes(mux)

	// Serve static files
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /openapi.json", openAPISpec)
	mux.HandleFunc("GET /docs", scalar.Handler)

	return mux
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.json")
}
