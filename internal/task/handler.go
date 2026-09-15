package task

import (
	"database/sql"
	"encoding/json"
	"net/http"

	// "strconv"

	"github.com/google/uuid"
)

type Handler struct {
	service ServiceInterface
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

// func getTaskID(r *http.Request) (string, error) {
// 	parts := strings.Split(
// 		strings.Trim(r.URL.Path, "/"),
// 		"/",
// 	)

// 	idString := parts[len(parts)-1]

// 	return strconv.Atoi(idString)

// 	return strconv.Atoi(r.PathValue("id"))
// }

func getTaskID(r *http.Request) (string, error) {
	id := r.PathValue("id")
	_, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

// ! Get    /api/v1/tasks
// GetTasks godoc
// @Summary Get all tasks
// @Description Get all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} Task
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks [get]
func (h *Handler) GetTasks(
	w http.ResponseWriter,
	r *http.Request,
) {
	tasks, err := h.service.GetTasks()

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch tasks",
		})
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

// ! Post    /api/v1/tasks
// CreateTask godoc
// @Summary Create a task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task"
// @Success 201 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks [post]
func (h *Handler) CreateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateTaskRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
		return
	}

	if request.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
		return
	}

	task, err := h.service.CreateTask(request.Title)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		task,
	)
}

// ! Get by Id	/api/v1/tasks/${id}
// GetTask godoc
// @Summary Get a task
// @Description Get a task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{id} [get]
func (h *Handler) GetTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := getTaskID(r)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task ID",
		})
		return
	}

	task, err := h.service.GetTask(id)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch task",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// ! Put    /api/v1/tasks/${id}
// UpdateTask godoc
// @Summary Update a task
// @Description Update a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body UpdateTaskRequest true "Task"
// @Success 200 {object} Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{id} [put]
func (h *Handler) UpdateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := getTaskID(r)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task ID",
		})
		return
	}

	var request UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
		return
	}

	if request.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
		return
	}

	err = h.service.UpdateTask(
		id,
		request.Title,
		request.Completed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to update task",
		})
		return
	}

	task, err := h.service.GetTask(id)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch updated task",
		})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// ! Delete    /api/v1/tasks/${id}
// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := getTaskID(r)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task ID",
		})
		return
	}

	err = h.service.DeleteTask(id)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "task not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to delete task",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
