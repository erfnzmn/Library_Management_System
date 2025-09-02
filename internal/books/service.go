package books

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBook(book *Book) error {
	return s.repo.CreateBook(book)
}

func (s *Service) GetAllBooks() ([]Book, error) {
	return s.repo.GetAllBooks()
}

func (s *Service) GetBookByID(id uint) (*Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *Service) UpdateBook(book *Book) error {
	return s.repo.UpdateBook(book)
}

func (s *Service) DeleteBook(id uint) error {
	return s.repo.DeleteBook(id)
}
