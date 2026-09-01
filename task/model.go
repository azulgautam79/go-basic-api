package task

type Task struct {
	ID        int    `json:"id" example:"1"`
	Title     string `json:"title" example:"Learn Go"`
	Completed bool   `json:"completed" example:"false"`
}

type CreateTaskRequest struct {
	Title string `json:"title" example:"Learn Go"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title" example:"Learn Go REST API"`
	Completed bool   `json:"completed" example:"true"`
}
