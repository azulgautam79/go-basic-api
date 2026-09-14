package task

type Service struct {
	repository RepositoryInterface
}

type ServiceInterface interface {
	CreateTask(title string) (*Task, error)
	GetTasks() ([]Task, error)
	GetTask(id string) (*Task, error)
	UpdateTask(id string, title string, completed bool) error
	DeleteTask(id string) error
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

//! Create Task 
func (s *Service) CreateTask(title string) (*Task, error) {
	return s.repository.Create(title)
}

func (s *Service) GetTasks() ([]Task, error) {
	return s.repository.FindAll()
}

func (s *Service) GetTask(id string) (*Task, error) {
	return s.repository.FindByID(id)
}

func (s *Service) UpdateTask(
	id string,
	title string,
	completed bool,
) error {
	return s.repository.Update(
		id,
		title,
		completed,
	)
}

func (s *Service) DeleteTask(id string) error {
	return s.repository.Delete(id)
}
