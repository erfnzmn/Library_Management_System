package books

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListBooks(ctx context.Context) ([]Book, error) {
	return s.repo.ListBooks(ctx)
}

func (s *Service) GetBook(ctx context.Context, id uint) (*Book, error) {
	return s.repo.GetBookByID(ctx, id)
}

func (s *Service) SearchBooks(ctx context.Context, q string) ([]Book, error) {
	return s.repo.SearchBooks(ctx, q)
}

func (s *Service) AddFavorite(ctx context.Context, userID, bookID int) error {
	return s.repo.AddToFavorites(ctx, userID, bookID)
}

func (s *Service) GetFavorites(ctx context.Context, userID int) ([]Book, error) {
	return s.repo.GetFavoritesByUser(ctx, userID)
}
