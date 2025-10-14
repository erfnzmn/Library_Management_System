package loans
import (
	"context"
	"errors"
	"time"


	books "github.com/erfnzmn/Library_Management_System/internal/books"
	"gorm.io/gorm"
)

var (
	ErrBookNotFound     = errors.New("book not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrNoStockAvailable = errors.New("no copies available for this book")
	ErrLoanNotFound     = errors.New("loan not found")
)

type Service struct {
	repo *Repository
	bookRepo *books.Repository
	db *gorm.DB
}

func NewService(db *gorm.DB, loanRepo *Repository, bookRepo *books.Repository) *Service {
	return &Service{
		db:       db,
		repo:     loanRepo,
		bookRepo: bookRepo,
	}
}

// ReserveBook
func (s *Service) ReserveBook(ctx context.Context, userID, bookID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// book inf
		book, err := s.bookRepo.GetBookByID(ctx, bookID)
		if err != nil {
			return err
		}
		if book == nil {
			return ErrBookNotFound
		}

		// check stock
		if book.Stock <= 0 {
			return ErrNoStockAvailable
		}

		// status change stock --
		book.Stock--
		if book.Stock == 0 {
			book.ReservationStatus = "reserved"
		}
		if err := s.bookRepo.UpdateBook(ctx, book); err != nil {
			return err
		}

		// new reserve recoed
		loan := &Loan{
			UserID:     userID,
			BookID:     bookID,
			Status:     StatusReserved,
			IsActive:   true,
			ReservedAt: time.Now(),
		}
		if err := s.repo.CreateLoan(ctx, loan); err != nil {
			return err
		}

		return nil
	})
}

// ConfirmBorrow
func (s *Service) ConfirmBorrow(ctx context.Context, loanID uint) error {
	now := time.Now()
	loan, err := s.repo.GetLoanByID(ctx, loanID)
	if err != nil {
		return err
	}
	if loan == nil {
		return ErrLoanNotFound
	}
	loan.Status = StatusBorrowed
	loan.BorrowedAt = &now
	due := now.AddDate(0, 0, 14)
	loan.DueDate = &due
	return s.repo.UpdateLoan(ctx, loan)
}

// ReturnBook stock ++
func (s *Service) ReturnBook(ctx context.Context, loanID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		loan, err := s.repo.GetLoanByID(ctx, loanID)
		if err != nil {
			return err
		}
		if loan == nil {
			return ErrLoanNotFound
		}

		book, err := s.bookRepo.GetBookByID(ctx, loan.BookID)
		if err != nil {
			return err
		}
		if book == nil {
			return ErrBookNotFound
		}

		//status change stock ++
		book.Stock++
		if book.Stock > 0 {
			book.ReservationStatus = "available"
		}
		if err := s.bookRepo.UpdateBook(ctx, book); err != nil {
			return err
		}

		now := time.Now()
		loan.Status = StatusReturned
		loan.ReturnedAt = &now
		loan.IsActive = false

		if err := s.repo.UpdateLoan(ctx, loan); err != nil {
			return err
		}

		return nil
	})
}

// CancelReservation
func (s *Service) CancelReservation(ctx context.Context, loanID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		loan, err := s.repo.GetLoanByID(ctx, loanID)
		if err != nil {
			return err
		}
		if loan == nil {
			return ErrLoanNotFound
		}

		book, err := s.bookRepo.GetBookByID(ctx, loan.BookID)
		if err != nil {
			return err
		}
		if book == nil {
			return ErrBookNotFound
		}

		book.Stock++
		if book.Stock > 0 {
			book.ReservationStatus = "available"
		}
		if err := s.bookRepo.UpdateBook(ctx, book); err != nil {
			return err
		}

		now := time.Now()
		loan.Status = StatusCancelled
		loan.CancelledAt = &now
		loan.IsActive = false

		return s.repo.UpdateLoan(ctx, loan)
	})
}
