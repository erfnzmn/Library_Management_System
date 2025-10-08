package books

import (
	"time"
)

type Book struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	Title             string `gorm:"type:varchar(200);not null" json:"title"`
	Author            string `gorm:"type:varchar(150);not null" json:"author"`
	ISBN              string `gorm:"type:varchar(20);uniqueIndex" json:"isbn"`
	Publisher         string `gorm:"type:varchar(150)" json:"publisher"`
	YearOfPublication int    `gorm:"type:int" json:"year_of_publication"`
	Edition           string `gorm:"type:varchar(50)" json:"edition"`
	Genre             string `gorm:"type:varchar(100)" json:"genre"`
	Language          string `gorm:"type:varchar(50);default:'fa'" json:"language"`
	Description       string `gorm:"type:text" json:"description"`

	ReservationStatus string `gorm:"type:enum('available','reserved');default:'available'" json:"reservation_status"`
	SellingStatus     string `gorm:"type:enum('available','sold_out');default:'available'" json:"selling_status"`

	CoverImage string  `gorm:"type:varchar(255)" json:"cover_image"`
	Tags       string  `gorm:"type:varchar(255)" json:"tags"`
	Price      float64 `gorm:"default:0" json:"price"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}
type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `json:"user_id"`
	BookID    uint      `json:'book_id`
	CreatedAt time.Time `json:"created_at"`
}

func (Book) TableName() string {
	return "books"
}
