package books

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBook(book *Book) error {
	return r.db.Create(book).Error
}

func (r *Repository) GetAllBooks() ([]Book, error) {
	var books []Book
	err := r.db.Find(&books).Error
	return books, err
}

func (r *Repository) GetBookByID(id uint) (*Book, error) {
	var book Book
	err := r.db.First(&book, id).Error
	return &book, err
}

func (r *Repository) UpdateBook(book *Book) error {
	return r.db.Save(book).Error
}

func (r *Repository) DeleteBook(id uint) error {
	return r.db.Delete(&Book{}, id).Error
}
