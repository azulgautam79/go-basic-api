package employees

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

type ServiceInterface interface {
	CreateEmployee(ctx context.Context, dto CreateEmployeeDTO) (*Employee, error)
	GetEmployees(ctx context.Context) ([]Employee, error)
	GetEmployeesWithFilter(ctx context.Context, params FilterParams) (*PaginatedResponse, error)
	GetEmployeeById(ctx context.Context, id uuid.UUID) (*Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, dto UpdateEmployeeDTO) (*Employee, error)
	InactiveEmployee(ctx context.Context, id uuid.UUID) error
	InactiveManyEmployees(ctx context.Context, dto BatchDeleteDTO) error
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
	DeleteManyEmployees(ctx context.Context, dto BatchDeleteDTO) error
}

// CreateEmployee creates a new employee record
func (s *Service) CreateEmployee(ctx context.Context, dto CreateEmployeeDTO) (*Employee, error) {
	return s.repository.Create(ctx, dto)
}

// GetEmployees fetches all active employees without pagination
func (s *Service) GetEmployees(ctx context.Context) ([]Employee, error) {
	return s.repository.FindAll(ctx)
}

// GetEmployeesWithFilter fetches employees with pagination, search, and filtering
func (s *Service) GetEmployeesWithFilter(ctx context.Context, params FilterParams) (*PaginatedResponse, error) {
	return s.repository.FindWithFilter(ctx, params)
}

// GetEmployeeById fetches a single employee by UUID
func (s *Service) GetEmployeeById(ctx context.Context, id uuid.UUID) (*Employee, error) {
	return s.repository.FindByID(ctx, id)
}

// UpdateEmployee modifies an existing employee's details
func (s *Service) UpdateEmployee(ctx context.Context, id uuid.UUID, dto UpdateEmployeeDTO) (*Employee, error) {
	return s.repository.Update(ctx, id, dto)
}

// InactiveEmployee marks an employee as inactive (is_active = false)
func (s *Service) InactiveEmployee(ctx context.Context, id uuid.UUID) error {
	return s.repository.ToggleActive(ctx, id, false)
}

// InactiveManyEmployees marks multiple employees as inactive
func (s *Service) InactiveManyEmployees(ctx context.Context, dto BatchDeleteDTO) error {
	return s.repository.ToggleActiveMany(ctx, dto.IDs, false)
}

// DeleteEmployee performs a soft delete on a single employee
func (s *Service) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return s.repository.SoftDelete(ctx, id)
}

// DeleteManyEmployees performs a soft delete on multiple employees
func (s *Service) DeleteManyEmployees(ctx context.Context, dto BatchDeleteDTO) error {
	return s.repository.SoftDeleteMany(ctx, dto.IDs)
}
