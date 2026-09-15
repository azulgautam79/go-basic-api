package employees

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RepositoryInterface interface {
	Create(ctx context.Context, dto CreateEmployeeDTO) (*Employee, error)
	FindAll(ctx context.Context) ([]Employee, error)
	FindWithFilter(ctx context.Context, params FilterParams) (*PaginatedResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateEmployeeDTO) (*Employee, error)
	ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error
	ToggleActiveMany(ctx context.Context, ids []uuid.UUID, isActive bool) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	SoftDeleteMany(ctx context.Context, ids []uuid.UUID) error
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

//! Create inserts a new employee into the database
func (r *Repository) Create(ctx context.Context, dto CreateEmployeeDTO) (*Employee, error) {
	query := `
		INSERT INTO employees (first_name, last_name, email, gender, is_active)
		VALUES ($1, $2, $3, $4, COALESCE($5, TRUE))
		RETURNING id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at`

	var emp Employee
	err := r.db.QueryRowContext(ctx, query,
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Gender,
		dto.IsActive,
	).Scan(
		&emp.ID,
		&emp.FirstName,
		&emp.LastName,
		&emp.Email,
		&emp.Gender,
		&emp.IsActive,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

//! FindAll returns all active, non-deleted employees without pagination
func (r *Repository) FindAll(ctx context.Context) ([]Employee, error) {
	query := `
		SELECT id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at
		FROM employees
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(
			&emp.ID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Email,
			&emp.Gender,
			&emp.IsActive,
			&emp.CreatedAt,
			&emp.UpdatedAt,
			&emp.DeletedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}

//! FindWithFilter applies dynamic search, filters, and pagination
func (r *Repository) FindWithFilter(ctx context.Context, params FilterParams) (*PaginatedResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	whereClause := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if params.Search != "" {
		whereClause = append(whereClause, fmt.Sprintf("(first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	if params.Gender != "" {
		whereClause = append(whereClause, fmt.Sprintf("gender = $%d", argIdx))
		args = append(args, params.Gender)
		argIdx++
	}

	if params.IsActive != nil {
		whereClause = append(whereClause, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *params.IsActive)
		argIdx++
	}

	whereStmt := strings.Join(whereClause, " AND ")

	// 1. Get Total Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM employees WHERE %s", whereStmt)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// 2. Fetch Paginated Records
	dataQuery := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at
		FROM employees
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereStmt, argIdx, argIdx+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []Employee{}
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(
			&emp.ID,
			&emp.FirstName,
			&emp.LastName,
			&emp.Email,
			&emp.Gender,
			&emp.IsActive,
			&emp.CreatedAt,
			&emp.UpdatedAt,
			&emp.DeletedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &PaginatedResponse{
		Data:       employees,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

//! FindByID retrieves a single non-deleted employee
func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	query := `
		SELECT id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at
		FROM employees
		WHERE id = $1 AND deleted_at IS NULL`

	var emp Employee
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&emp.ID,
		&emp.FirstName,
		&emp.LastName,
		&emp.Email,
		&emp.Gender,
		&emp.IsActive,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee with id %s not found", id)
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}

//! Update modifies specific non-nil fields in the employee record using COALESCE
func (r *Repository) Update(ctx context.Context, id uuid.UUID, dto UpdateEmployeeDTO) (*Employee, error) {
	query := `
		UPDATE employees
		SET first_name = COALESCE($1, first_name),
		    last_name = COALESCE($2, last_name),
		    email = COALESCE($3, email),
		    gender = COALESCE($4, gender),
		    is_active = COALESCE($5, is_active),
		    updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, first_name, last_name, email, gender, is_active, created_at, updated_at, deleted_at`

	var emp Employee
	err := r.db.QueryRowContext(ctx, query,
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Gender,
		dto.IsActive,
		id,
	).Scan(
		&emp.ID,
		&emp.FirstName,
		&emp.LastName,
		&emp.Email,
		&emp.Gender,
		&emp.IsActive,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee with id %s not found or deleted", id)
	} else if err != nil {
		return nil, err
	}

	return &emp, nil
}

//! ToggleActive sets the is_active status of an employee
func (r *Repository) ToggleActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `UPDATE employees SET is_active = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, isActive, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("employee not found or deleted")
	}

	return nil
}

//! ToggleActiveMany sets the is_active status for a slice of employee IDs
func (r *Repository) ToggleActiveMany(ctx context.Context, ids []uuid.UUID, isActive bool) error {
	query := `UPDATE employees SET is_active = $1, updated_at = NOW() WHERE id = ANY($2) AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, isActive, pq.Array(ids))
	return err
}

//! SoftDelete sets deleted_at timestamp for a single employee
func (r *Repository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE employees SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("employee not found or already deleted")
	}

	return nil
}

//! SoftDeleteMany performs batch soft deletion for a list of UUIDs
func (r *Repository) SoftDeleteMany(ctx context.Context, ids []uuid.UUID) error {
	query := `UPDATE employees SET deleted_at = NOW() WHERE id = ANY($1) AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, pq.Array(ids))
	return err
}
