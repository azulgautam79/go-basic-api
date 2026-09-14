package task

import (
	"strconv"
	"testing"
)

type fakeRepository struct {
	tasks []Task
}

//! Create Task
func (f *fakeRepository) Create(title string) (*Task, error) {
	task := Task{
		ID:        strconv.Itoa(len(f.tasks) + 1),
		Title:     title,
		Completed: false,
	}

	f.tasks = append(f.tasks, task)

	return &task, nil
}

//! Find All
func (f *fakeRepository) FindAll() ([]Task, error) {
	return f.tasks, nil
}

//! Find By id
func (f *fakeRepository) FindByID(id string) (*Task, error) {
	for _, task := range f.tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, nil
}

//! Update
func (f *fakeRepository) Update(
	id string,
	title string,
	completed bool,
) error {
	for i := range f.tasks {
		if f.tasks[i].ID == id {
			f.tasks[i].Title = title
			f.tasks[i].Completed = completed
			return nil
		}
	}

	return nil
}

//! Delete
func (f *fakeRepository) Delete(id string) error {
	for i := range f.tasks {
		if f.tasks[i].ID == id {
			f.tasks = append(
				f.tasks[:i],
				f.tasks[i+1:]...,
			)

			return nil
		}
	}

	return nil
}

// func TestServiceCreateTask(t *testing.T) {
// 	repository := &fakeRepository{}
// 	service := NewService(repository)

// 	task, err := service.CreateTask("Learn Go")

// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if task.Title != "Learn Go" {
// 		t.Errorf(
// 			"expected title %q, got %q",
// 			"Learn Go",
// 			task.Title,
// 		)
// 	}

// 	if task.Completed {
// 		t.Error("expected task to be incomplete")
// 	}
// }

//! Table Driver test
//* Create Task
func TestServiceCreateTask(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{
			name:  "simple task",
			title: "Learn Go",
		},
		{
			name:  "long task title",
			title: "Learn Go REST API with SQLite",
		},
		{
			name:  "task with spaces",
			title: "Build my first Go API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeRepository{}
			service := NewService(repository)

			task, err := service.CreateTask(tt.title)

			if err != nil {
				t.Fatalf(
					"expected no error, got %v",
					err,
				)
			}

			if task.Title != tt.title {
				t.Errorf(
					"expected title %q, got %q",
					tt.title,
					task.Title,
				)
			}

			if task.Completed {
				t.Error(
					"expected task to be incomplete",
				)
			}
		})
	}
}

//* Get Tasks
func TestServiceGetTasks(t *testing.T) {
	repository := &fakeRepository{
		tasks: []Task{
			{
				ID:        "1",
				Title:     "Learn Go",
				Completed: false,
			},
			{
				ID:        "2",
				Title:     "Build API",
				Completed: true,
			},
		},
	}

	service := NewService(repository)

	tasks, err := service.GetTasks()

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if len(tasks) != 2 {
		t.Fatalf(
			"expected 2 tasks, got %d",
			len(tasks),
		)
	}

	if tasks[0].Title != "Learn Go" {
		t.Errorf(
			"expected first task to be %q, got %q",
			"Learn Go",
			tasks[0].Title,
		)
	}
}

//* Get Task by Id
func TestServiceGetTask(t *testing.T) {
	repository := &fakeRepository{
		tasks: []Task{
			{
				ID:        "1",
				Title:     "Learn Go",
				Completed: false,
			},
		},
	}

	service := NewService(repository)

	task, err := service.GetTask("1")

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID != "1" {
		t.Errorf(
			"expected ID 1, got %s",
			task.ID,
		)
	}

	if task.Title != "Learn Go" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Go",
			task.Title,
		)
	}
}

//* Update Task
func TestServiceUpdateTask(t *testing.T) {
	repository := &fakeRepository{
		tasks: []Task{
			{
				ID:        "1",
				Title:     "Learn Go",
				Completed: false,
			},
		},
	}

	service := NewService(repository)

	err := service.UpdateTask(
		"1",
		"Learn Go REST API",
		true,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	task, err := service.GetTask("1")

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if task.Title != "Learn Go REST API" {
		t.Errorf(
			"expected updated title, got %q",
			task.Title,
		)
	}

	if !task.Completed {
		t.Error(
			"expected task to be completed",
		)
	}
}

//* Delete Task
func TestServiceDeleteTask(t *testing.T) {
	repository := &fakeRepository{
		tasks: []Task{
			{
				ID:    "1",
				Title: "Learn Go",
			},
		},
	}

	service := NewService(repository)

	err := service.DeleteTask("1")

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	tasks, err := service.GetTasks()

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if len(tasks) != 0 {
		t.Errorf(
			"expected 0 tasks, got %d",
			len(tasks),
		)
	}
}
