package books

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListBooks(ctx context.Context) ([]Book, error) {
	var books []Book
	if err := r.db.WithContext(ctx).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) GetBookByID(ctx context.Context, id uint) (*Book, error) {
	var book Book
	if err := r.db.WithContext(ctx).First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *Repository) SearchBooks(ctx context.Context, q string) ([]Book, error) {
	var books []Book
	if err := r.db.WithContext(ctx).
		Where("title LIKE ? OR author LIKE ?", "%"+q+"%", "%"+q+"%").
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) AddToFavorites(ctx context.Context, userID, bookID int) error {
	fav := Favorite{UserID: userID, BookID: bookID}
	return r.db.WithContext(ctx).
		FirstOrCreate(&fav, Favorite{UserID: userID, BookID: bookID}).Error
}

func (r *Repository) GetFavoritesByUser(ctx context.Context, userID int) ([]Book, error) {
	var books []Book
	if err := r.db.WithContext(ctx).
		Table("books").
		Joins("JOIN favorites f ON f.book_id = books.id").
		Where("f.user_id = ?", userID).
		Order("f.created_at DESC").
		Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}
