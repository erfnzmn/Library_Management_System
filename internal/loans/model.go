package loans

import (
	"time"
)

const (
	StatusReserved = "reserved"
	StatusBorrowed = "borrowed"
	StatusReturned = "returned"
	StatusCancelled = "cancelled"
)

type Loan struct {
	ID	uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	BookID    uint      `gorm:"not null;index" json:"book_id"`
	Status    string    `gorm:"type:enum('reserved','borrowed','returned','cancelled');not null;default:'reserved'" json:"status"`
	IsActive bool   `gorm:"not null;default:true" json:"is_active"`
	ReservedAt  time.Time  `gorm:"not null;autoCreateTime" json:"reserved_at"`
	BorrowedAt  *time.Time `json:"borrowed_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ReturnedAt  *time.Time `json:"returned_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	Notes       string     `gorm:"type:varchar(255)" json:"notes,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`	
}

func (Loan) TableName() string { return "loans"  }