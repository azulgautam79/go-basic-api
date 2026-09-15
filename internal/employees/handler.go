package employees

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

// Helpers

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// ! Post    /api/v1/employees
// CreateEmployee godoc
// @Summary      Create a new employee
// @Description  Creates an employee record in the system
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        request body CreateEmployeeDTO true "Create Employee Request"
// @Success      201  {object}  Employee
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/employees [post]
func (h *Handler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var dto CreateEmployeeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	emp, err := h.service.CreateEmployee(r.Context(), dto)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, emp)
}

// ! Get    /api/v1/employees/all
// GetAllEmployees godoc
// @Summary      Get all employees (Unpaginated)
// @Description  Fetches a list of all active employees without pagination
// @Tags         employees
// @Produce      json
// @Success      200  {array}   Employee
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/employees/all [get]
func (h *Handler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.GetEmployees(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, employees)
}

// ! Get    /api/v1/employees
// GetEmployeesWithFilter godoc
// @Summary      Get filtered and paginated employees
// @Description  Retrieves employees with support for search, gender/active filters, and pagination
// @Tags         employees
// @Produce      json
// @Param        page      query     int     false  "Page number (default 1)"
// @Param        limit     query     int     false  "Page size limit (default 10)"
// @Param        search    query     string  false  "Search term for name or email"
// @Param        gender    query     string  false  "Filter by gender" Enums(male, female, other, prefer_not_to_say)
// @Param        is_active query     bool    false  "Filter by active status"
// @Success      200       {object}  PaginatedResponse
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/employees [get]
func (h *Handler) GetEmployeesWithFilter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	search := query.Get("search")
	genderStr := query.Get("gender")
	isActiveStr := query.Get("is_active")

	var isActive *bool
	if isActiveStr != "" {
		val, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &val
		}
	}

	params := FilterParams{
		Page:     page,
		Limit:    limit,
		Search:   search,
		Gender:   Gender(genderStr),
		IsActive: isActive,
	}

	result, err := h.service.GetEmployeesWithFilter(r.Context(), params)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// ! Get    /api/v1/employees/{id}
// GetEmployeeByID godoc
// @Summary      Get employee by ID
// @Description  Fetches details of a single employee by their UUID
// @Tags         employees
// @Produce      json
// @Param        id   path      string  true  "Employee UUID"
// @Success      200  {object}  Employee
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/employees/{id} [get]
func (h *Handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid employee UUID")
		return
	}

	emp, err := h.service.GetEmployeeById(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emp)
}

// ! Put    /api/v1/employees/{id}
// UpdateEmployee godoc
// @Summary      Update employee details
// @Description  Updates select fields of an employee by UUID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id       path      string             true  "Employee UUID"
// @Param        request  body      UpdateEmployeeDTO  true  "Update Payload"
// @Success      200      {object}  Employee
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/employees/{id} [put]
func (h *Handler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid employee UUID")
		return
	}

	var dto UpdateEmployeeDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	emp, err := h.service.UpdateEmployee(r.Context(), id, dto)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, emp)
}

// ! Patch    /api/v1/employees/{id}/inactive
// InactiveEmployee godoc
// @Summary      Inactivate an employee
// @Description  Sets an employee's is_active field to false
// @Tags         employees
// @Produce      json
// @Param        id   path      string  true  "Employee UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/employees/{id}/inactive [patch]
func (h *Handler) InactiveEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid employee UUID")
		return
	}

	if err := h.service.InactiveEmployee(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "employee marked as inactive"})
}

// ! Patch    /api/v1/employees/inactive-many
// InactiveManyEmployees godoc
// @Summary      Inactivate multiple employees
// @Description  Sets is_active field to false for a list of UUIDs
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        request  body      BatchDeleteDTO  true  "Batch Inactivate Request"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/employees/inactive-many [patch]
func (h *Handler) InactiveManyEmployees(w http.ResponseWriter, r *http.Request) {
	var dto BatchDeleteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil || len(dto.IDs) == 0 {
		respondError(w, http.StatusBadRequest, "invalid payload or empty ID list")
		return
	}

	if err := h.service.InactiveManyEmployees(r.Context(), dto); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "employees marked as inactive"})
}

// ! Delete    /api/v1/employees/{id}
// DeleteEmployee godoc
// @Summary      Soft delete an employee
// @Description  Performs a soft delete on an employee record by setting deleted_at
// @Tags         employees
// @Produce      json
// @Param        id   path      string  true  "Employee UUID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/employees/{id} [delete]
func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid employee UUID")
		return
	}

	if err := h.service.DeleteEmployee(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "employee deleted successfully"})
}

// ! Delete    /api/v1/employees/batch-delete
// DeleteManyEmployees godoc
// @Summary      Batch soft delete employees
// @Description  Soft deletes multiple employee records in a single request
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        request  body      BatchDeleteDTO  true  "Batch Delete Payload"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/employees/batch-delete [post]
func (h *Handler) DeleteManyEmployees(w http.ResponseWriter, r *http.Request) {
	var dto BatchDeleteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil || len(dto.IDs) == 0 {
		respondError(w, http.StatusBadRequest, "invalid payload or empty ID list")
		return
	}

	if err := h.service.DeleteManyEmployees(r.Context(), dto); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "employees batch deleted successfully"})
}
