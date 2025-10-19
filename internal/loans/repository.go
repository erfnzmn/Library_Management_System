package loans

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateLoan(ctx context.Context, loan *Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *Repository) GetLoanByID(ctx context.Context, id uint) (*Loan, error) {
	var loan Loan
	if err := r.db.WithContext(ctx).First(&loan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &loan, nil
}

func (r *Repository) GetActiveLoansByBook(ctx context.Context, bookID uint) ([]Loan, error) {
	var loans []Loan
	if err := r.db.WithContext(ctx).
		Where("book_id = ? AND is_active = TRUE", bookID).
		Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *Repository) UpdateLoan(ctx context.Context, loan *Loan) error {
	return r.db.WithContext(ctx).Save(loan).Error
}

func (r *Repository) GetLoansByUser(ctx context.Context, userID uint) ([]Loan, error) {
	var loans []Loan
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}
