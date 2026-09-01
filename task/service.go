package task

type Service struct {
	repository RepositoryInterface
}

type ServiceInterface interface {
	CreateTask(title string) (*Task, error)
	GetTasks() ([]Task, error)
	GetTask(id int) (*Task, error)
	UpdateTask(id int, title string, completed bool) error
	DeleteTask(id int) error
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateTask(title string) (*Task, error) {
	return s.repository.Create(title)
}

func (s *Service) GetTasks() ([]Task, error) {
	return s.repository.FindAll()
}

func (s *Service) GetTask(id int) (*Task, error) {
	return s.repository.FindByID(id)
}

func (s *Service) UpdateTask(
	id int,
	title string,
	completed bool,
) error {
	return s.repository.Update(
		id,
		title,
		completed,
	)
}

func (s *Service) DeleteTask(id int) error {
	return s.repository.Delete(id)
}
