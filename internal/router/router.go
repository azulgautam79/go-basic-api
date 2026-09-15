package router

import (
	"encoding/json"
	"net/http"

	"github.com/azulgautam79/go-sqlite-tasks/internal/employees"
	"github.com/azulgautam79/go-sqlite-tasks/internal/scalar"
	"github.com/azulgautam79/go-sqlite-tasks/internal/task"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	taskHandler     *task.Handler
	employeeHandler *employees.Handler
}

func New(
	taskHandler *task.Handler,
	employeeHandler *employees.Handler,
) http.Handler {

	r := &Router{
		taskHandler:     taskHandler,
		employeeHandler: employeeHandler,
	}

	return r.routes()
}

func (r *Router) routes() http.Handler {
	mux := http.NewServeMux()

	//! Home Route (Use GET /{$} so it matches exact root path only, preventing catch-all conflicts)
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"message": "API is running",
		})
	})

	//! Health
	mux.HandleFunc("GET /healthz", Healthz)

	r.registerTaskRoutes(mux)
	r.registerEmployeesRoutes(mux)

	// Serve static files (Specify GET method explicitly)
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", static))
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /openapi.json", openAPISpec)
	mux.HandleFunc("GET /docs", scalar.Handler)

	return mux
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/swagger.json")
}
