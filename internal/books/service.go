package books

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo  *Repository
	cache *redis.Client
}

func NewService(repo *Repository, cache *redis.Client) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) cacheKey(id uint) string {
	return fmt.Sprintf("book:%d", id)
}

func (s *Service) cacheListKey() string {
	return "books:all"
}

// CreateBook — هم دیتا ذخیره میشه، هم کش پاک میشه
func (s *Service) CreateBook(ctx context.Context, book *Book) error {
	if err := s.repo.CreateBook(ctx, book); err != nil {
		return err
	}
	// پاک‌سازی کش
	s.cache.Del(ctx, s.cacheListKey())
	return nil
}

func (s *Service) UpdateBook(ctx context.Context, book *Book) error {
	if err := s.repo.UpdateBook(ctx, book); err != nil {
		return err
	}
	s.cache.Del(ctx, s.cacheKey(book.ID))
	s.cache.Del(ctx, s.cacheListKey())
	return nil
}

func (s *Service) DeleteBook(ctx context.Context, id uint) error {
	if err := s.repo.DeleteBook(ctx, id); err != nil {
		return err
	}
	s.cache.Del(ctx, s.cacheKey(id))
	s.cache.Del(ctx, s.cacheListKey())
	return nil
}

func (s *Service) GetBookByID(ctx context.Context, id uint) (*Book, error) {
	key := s.cacheKey(id)
	val, err := s.cache.Get(ctx, key).Result()
	if err == nil {
		var b Book
		if json.Unmarshal([]byte(val), &b) == nil {
			return &b, nil
		}
	}

	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(book)
	s.cache.Set(ctx, key, data, 10*time.Minute)
	return book, nil
}

func (s *Service) ListBooks(ctx context.Context) ([]Book, error) {
	key := s.cacheListKey()
	val, err := s.cache.Get(ctx, key).Result()
	if err == nil {
		var books []Book
		if json.Unmarshal([]byte(val), &books) == nil {
			return books, nil
		}
	}

	books, err := s.repo.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(books)
	s.cache.Set(ctx, key, data, 5*time.Minute)
	return books, nil
}

func (s *Service) SearchBooks(ctx context.Context, q string) ([]Book, error) {
	return s.repo.SearchBooks(ctx, q)
}

func (s *Service) AddToFavorites(ctx context.Context, userID, bookID uint) error {
	return s.repo.AddToFavorites(ctx, userID, bookID)
}

func (s *Service) GetFavoritesByUser(ctx context.Context, userID int) ([]Book, error) {
	return s.repo.GetFavoritesByUser(ctx, userID)
}
