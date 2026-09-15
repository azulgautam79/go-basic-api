package router

import "net/http"

func (r *Router) registerEmployeesRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/employees", r.employeeHandler.CreateEmployee)
	mux.HandleFunc("GET /api/v1/employees", r.employeeHandler.GetEmployeesWithFilter)
	mux.HandleFunc("GET /api/v1/employees/all", r.employeeHandler.GetAllEmployees)
	mux.HandleFunc("GET /api/v1/employees/{id}", r.employeeHandler.GetEmployeeByID)
	mux.HandleFunc("PUT /api/v1/employees/{id}", r.employeeHandler.UpdateEmployee)
	mux.HandleFunc("PATCH /api/v1/employees/{id}/inactive", r.employeeHandler.InactiveEmployee)
	mux.HandleFunc("PATCH /api/v1/employees/inactive-many", r.employeeHandler.InactiveManyEmployees)
	mux.HandleFunc("DELETE /api/v1/employees/{id}", r.employeeHandler.DeleteEmployee)
	mux.HandleFunc("POST /api/v1/employees/batch-delete", r.employeeHandler.DeleteManyEmployees)
}
