package users

type UserService interface {
	Register(user User) (User, error)
	GetUser(id int) (User, error)
	GetUsers() ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int) error
}
type userService struct {
	repo userrepository
}

func NewUserservice(repo userrepository) UserService {
	return &userService{repo: repo}
}
func (s *userService) Register(user User) (User, error) {
	return s.repo.create(user)
}

func (s *userService) GetUser(id int) (User, error) {
	return s.repo.getbyid(id)
}
func (s *userService) GetUsers() ([]User, error) {
	return s.repo.getAll()
}
func (s *userService) UpdateUser(user User) (User, error) {
	return s.repo.update(user)
}
func (s *userService) DeleteUser(id int) error {
	return s.repo.delete(id)
}
