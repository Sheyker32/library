package domain

import (
	"time"
)

type Author struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Biography string    `db:"biography"`
	Books     []Book    `db:"books,omitempty"`
	CreatedAt time.Time `db:"created_at" swaggertype:"string" format:"date-time"`
}

type Book struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	AuthorID  int       `db:"author_id"`
	Author    *Author   `db:"author"`
	Available bool      `db:"available"`
	CreatedAt time.Time `db:"created_at" swaggertype:"string" format:"date-time"`
}

type User struct {
	ID          int          `db:"id"`
	Name        string       `db:"name"`
	Email       string       `db:"email"`
	CreatedAt   time.Time    `db:"created_at" swaggertype:"string" format:"date-time"`
	RentedBooks []BookRental `db:"rented_books"`
}

type BookRental struct {
	ID         int        `db:"id"`
	BookID     int        `db:"book_id"`
	UserID     int        `db:"user_id"`
	RentalDate time.Time  `db:"rental_date"`
	ReturnDate *time.Time `db:"return_date"`
	CreatedAt  time.Time  `db:"created_at" swaggertype:"string" format:"date-time"`
}

type UniqueBookRental struct {
	BookID int `db:"book_id"`
	UserID int `db:"user_id"`
}

type AuthorWithRentCount struct {
	Author    Author
	RentCount int
}
