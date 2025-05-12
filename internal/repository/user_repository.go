package repository

import (
	"context"
	"fmt"
	"library/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Userer interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int) error
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) Userer {
	return &UserRepository{db: db}
}

func (u UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2,$3) RETURNING id`
	err := u.db.QueryRowContext(ctx, query, user.ID, user.Email, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	err := u.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to user: %w", err)
	}

	var rentals []domain.BookRental
	queryActivRental := `SELECT id, book_id, user_id, rental_date, return_date, created_at
			  FROM book_rental
			  WHERE user_id = $1 AND return_date IS NULL`

	err = u.db.SelectContext(ctx, &rentals, queryActivRental, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list rentals: %w", err)
	}
	user.RentedBooks = append(user.RentedBooks, rentals...)

	return &user, nil
}

func (u UserRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT * FROM users`
	err := u.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	for _, value := range users {
		var rentals []domain.BookRental
		queryActivRental := `SELECT id, book_id, user_id, rental_date, return_date, created_at
			  FROM book_rental
			  WHERE user_id = $1`

		err = u.db.SelectContext(ctx, &rentals, queryActivRental, value.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to list rentals: %w", err)
		}
		value.RentedBooks = append(value.RentedBooks, rentals...)
	}

	return users, nil
}

func (u *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
