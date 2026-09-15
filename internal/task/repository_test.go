package task

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// ! Test DB
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(`
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)

	if err != nil {
		t.Fatalf("failed to create tasks table: %v", err)
	}

	return db
}

// ! Create Task
func TestRepositoryCreate(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	task, err := repository.Create("Learn Go")

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID == "0" {
		t.Error("expected task ID to be generated")
	}

	if task.Title != "Learn Go" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Go",
			task.Title,
		)
	}

	if task.Completed {
		t.Error("expected new task to be incomplete")
	}
}

// ! Find All
func TestRepositoryFindAll(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	_, err := repository.Create("Learn Go")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	_, err = repository.Create("Build REST API")
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	tasks, err := repository.FindAll()

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
			"expected first task %q, got %q",
			"Learn Go",
			tasks[0].Title,
		)
	}

	if tasks[1].Title != "Build REST API" {
		t.Errorf(
			"expected second task %q, got %q",
			"Build REST API",
			tasks[1].Title,
		)
	}
}

// ! Find By Id
func TestRepositoryFindByID(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	createdTask, err := repository.Create("Learn Go")

	if err != nil {
		t.Fatalf(
			"failed to create task: %v",
			err,
		)
	}

	task, err := repository.FindByID(createdTask.ID)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID != createdTask.ID {
		t.Errorf(
			"expected ID %s, got %s",
			createdTask.ID,
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

	if task.Completed {
		t.Error("expected task to be incomplete")
	}
}

// ! Find By Id Not Found
func TestRepositoryFindByIDNotFound(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	task, err := repository.FindByID("999")

	if task != nil {
		t.Error("expected nil task")
	}

	if err != sql.ErrNoRows {
		t.Errorf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

// ! Update
func TestRepositoryUpdate(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	createdTask, err := repository.Create("Learn Go")

	if err != nil {
		t.Fatalf(
			"failed to create task: %v",
			err,
		)
	}

	err = repository.Update(
		createdTask.ID,
		"Learn Go REST API",
		true,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	task, err := repository.FindByID(createdTask.ID)

	if err != nil {
		t.Fatalf(
			"failed to find updated task: %v",
			err,
		)
	}

	if task.Title != "Learn Go REST API" {
		t.Errorf(
			"expected title %q, got %q",
			"Learn Go REST API",
			task.Title,
		)
	}

	if !task.Completed {
		t.Error("expected task to be completed")
	}
}

// ! Update Not Found of task that doesn't exist
func TestRepositoryUpdateNotFound(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	err := repository.Update(
		"999",
		"Does not exist",
		true,
	)

	if err != sql.ErrNoRows {
		t.Errorf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}

// ! Delete
func TestRepositoryDelete(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	createdTask, err := repository.Create("Learn Go")

	if err != nil {
		t.Fatalf(
			"failed to create task: %v",
			err,
		)
	}

	err = repository.Delete(createdTask.ID)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	_, err = repository.FindByID(createdTask.ID)

	if err != sql.ErrNoRows {
		t.Errorf(
			"expected sql.ErrNoRows after delete, got %v",
			err,
		)
	}
}

// ! Delete nonexistent task
func TestRepositoryDeleteNotFound(t *testing.T) {
	db := setupTestDB(t)

	repository := NewRepository(db)

	err := repository.Delete("999")

	if err != sql.ErrNoRows {
		t.Errorf(
			"expected sql.ErrNoRows, got %v",
			err,
		)
	}
}
