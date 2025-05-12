package repository

import (
	"context"
	"library/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
)

type Rentaler interface {
	RentBook(ctx context.Context, bookID, userID int) error
	ReturnBook(ctx context.Context, bookId int) error
}

type RentalRepository struct {
	db *sqlx.DB
}

func NewRentalRepository(db *sqlx.DB) Rentaler {
	return &RentalRepository{db: db}
}

func (r RentalRepository) RentBook(ctx context.Context, bookID, userID int) error {
	uRental := domain.UniqueBookRental{
		BookID: bookID,
		UserID: userID,
	}
	queryBookRental := `INSERT INTO book_rental (book_id, user_id, rental_date) VALUES ($1, $2,$3)`
	_, err := r.db.ExecContext(ctx, queryBookRental, bookID, userID, time.Now())
	if err != nil {
		return err
	}

	queryUnique := `INSERT INTO unique_book_rental (book_id, user_id) VALUES(:book_id,:user_id)`
	_, err = r.db.NamedExecContext(ctx, queryUnique, uRental)
	if err != nil {
		return err
	}

	return nil
}
func (r RentalRepository) ReturnBook(ctx context.Context, bookId int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM unique_book_rental WHERE book_id=$1", bookId)
	if err != nil {
		return err
	}
	returnDate := time.Now()
	_, err = r.db.ExecContext(ctx, "UPDATE book_rental SET return_date = $1 WHERE book_id=$2 AND return_date IS NULL", returnDate, bookId)
	if err != nil {
		return err
	}

	return nil
}
