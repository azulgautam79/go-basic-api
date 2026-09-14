package task

type Task struct {
	ID        string `json:"id" example:"8a7572fe-411b-4f4a-b560-ea39948062d0"`
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
