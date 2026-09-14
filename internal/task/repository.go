package task

import "database/sql"

type RepositoryInterface interface {
	Create(title string) (*Task, error)
	FindAll() ([]Task, error)
	FindByID(id string) (*Task, error)
	Update(id string, title string, completed bool) error
	Delete(id string) error
}

type Repository struct {
	db *sql.DB
}

var _ RepositoryInterface = (*Repository)(nil)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// ! Create Task
func (r *Repository) Create(title string) (*Task, error) {
	var task Task

	err := r.db.QueryRow(
		`
		INSERT INTO tasks(title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed
		`,
		title,
		false,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ! Find All
func (r *Repository) FindAll() ([]Task, error) {
	rows, err := r.db.Query(`
	SELECT id, title, completed
	FROM tasks
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// ! Find By Id
func (r *Repository) FindByID(id string) (*Task, error) {
	var task Task

	err := r.db.QueryRow(
		`
		SELECT id, title, completed
		FROM tasks
		WHERE id = $1
		`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
	)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// ! Update
func (r *Repository) Update(
	id string,
	title string,
	completed bool,
) error {
	result, err := r.db.Exec(
		`
		UPDATE tasks
		SET title = $1, completed = $2
		WHERE id = $3
		`,
		title,
		completed,
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ! Delete
func (r *Repository) Delete(id string) error {
	result, err := r.db.Exec(
		"DELETE FROM tasks WHERE id = $1",
		id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
